package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
	"github.com/tranvuongduy2003/go-mvc/pkg/jwt"
)

// MiddlewareManager manages all application middleware
type MiddlewareManager struct {
	logger      *logger.Logger
	rateLimiter *RateLimiter
	config      MiddlewareConfig
}

// MiddlewareConfig represents middleware configuration
type MiddlewareConfig struct {
	// Rate limiting
	RateLimit struct {
		Enabled bool
		RPS     int
		Burst   int
	}

	// CORS
	CORS struct {
		Enabled        bool
		AllowedOrigins []string
		AllowedMethods []string
		AllowedHeaders []string
	}

	// Security
	Security struct {
		Enabled           bool
		FrameOptions      string
		ReferrerPolicy    string
		CSP               string
		PermissionsPolicy string
	}

	// Logging
	Logging struct {
		Enabled         bool
		SkipPaths       []string
		LogRequestBody  bool
		LogResponseBody bool
		MaxBodySize     int64
	}

	// Metrics
	Metrics struct {
		Enabled bool
		Path    string
	}

	// Timeout
	Timeout struct {
		Enabled  bool
		Duration time.Duration
	}

	// API Key Authentication
	APIAuth struct {
		Enabled    bool
		ValidKeys  map[string]string
		HeaderName string
	}

	// IP Whitelist
	IPWhitelist struct {
		Enabled    bool
		AllowedIPs []string
	}

	// Request Size Limit
	SizeLimit struct {
		Enabled bool
		MaxSize int64
	}
}

// DefaultMiddlewareConfig returns default middleware configuration
func DefaultMiddlewareConfig() MiddlewareConfig {
	return MiddlewareConfig{
		RateLimit: struct {
			Enabled bool
			RPS     int
			Burst   int
		}{
			Enabled: true,
			RPS:     100,
			Burst:   200,
		},
		CORS: struct {
			Enabled        bool
			AllowedOrigins []string
			AllowedMethods []string
			AllowedHeaders []string
		}{
			Enabled:        true,
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Origin", "Content-Type", "Authorization", "X-Request-ID", "X-API-Key"},
		},
		Security: struct {
			Enabled           bool
			FrameOptions      string
			ReferrerPolicy    string
			CSP               string
			PermissionsPolicy string
		}{
			Enabled:           true,
			FrameOptions:      "DENY",
			ReferrerPolicy:    "strict-origin-when-cross-origin",
			CSP:               "default-src 'self'",
			PermissionsPolicy: "camera=(), microphone=(), geolocation=()",
		},
		Logging: struct {
			Enabled         bool
			SkipPaths       []string
			LogRequestBody  bool
			LogResponseBody bool
			MaxBodySize     int64
		}{
			Enabled:         true,
			SkipPaths:       []string{"/health", "/metrics"},
			LogRequestBody:  false,
			LogResponseBody: false,
			MaxBodySize:     1024 * 1024, // 1MB
		},
		Metrics: struct {
			Enabled bool
			Path    string
		}{
			Enabled: true,
			Path:    "/metrics",
		},
		Timeout: struct {
			Enabled  bool
			Duration time.Duration
		}{
			Enabled:  true,
			Duration: 30 * time.Second,
		},
		APIAuth: struct {
			Enabled    bool
			ValidKeys  map[string]string
			HeaderName string
		}{
			Enabled:    false,
			ValidKeys:  map[string]string{},
			HeaderName: "X-API-Key",
		},
		IPWhitelist: struct {
			Enabled    bool
			AllowedIPs []string
		}{
			Enabled:    false,
			AllowedIPs: []string{},
		},
		SizeLimit: struct {
			Enabled bool
			MaxSize int64
		}{
			Enabled: true,
			MaxSize: 10 * 1024 * 1024, // 10MB
		},
	}
}

// NewMiddlewareManager creates a new middleware manager
func NewMiddlewareManager(
	logger *logger.Logger,
	config MiddlewareConfig,
	jwtService jwt.JWTService,
) *MiddlewareManager {
	var rateLimiter *RateLimiter
	if config.RateLimit.Enabled {
		rateLimiter = NewRateLimiter(config.RateLimit.RPS, config.RateLimit.Burst)
		rateLimiter.CleanupOldLimiters()
	}

	return &MiddlewareManager{
		logger:      logger,
		rateLimiter: rateLimiter,
		config:      config,
	}
}

// SetupMiddleware configures all middleware for the Gin engine
func (mm *MiddlewareManager) SetupMiddleware(r *gin.Engine) {
	// Request ID (should be first)
	r.Use(RequestIDMiddleware())

	// Security headers
	if mm.config.Security.Enabled {
		securityConfig := SecurityConfig{
			FrameOptions:      mm.config.Security.FrameOptions,
			ReferrerPolicy:    mm.config.Security.ReferrerPolicy,
			CSP:               mm.config.Security.CSP,
			PermissionsPolicy: mm.config.Security.PermissionsPolicy,
			HSTSMaxAge:        "max-age=31536000; includeSubDomains",
			CustomHeaders:     map[string]string{"X-API-Version": "v1"},
		}
		r.Use(SecureHeaders(securityConfig))
	}

	// CORS
	if mm.config.CORS.Enabled {
		if len(mm.config.CORS.AllowedOrigins) == 1 && mm.config.CORS.AllowedOrigins[0] == "*" {
			r.Use(DevCORSMiddleware())
		} else {
			r.Use(ProductionCORSMiddleware(mm.config.CORS.AllowedOrigins))
		}
	}

	// Recovery
	r.Use(DefaultRecoveryMiddleware(mm.logger))

	// Logging
	if mm.config.Logging.Enabled {
		loggerConfig := LoggerConfig{
			Logger:          mm.logger,
			SkipPaths:       mm.config.Logging.SkipPaths,
			LogRequestBody:  mm.config.Logging.LogRequestBody,
			LogResponseBody: mm.config.Logging.LogResponseBody,
			MaxBodySize:     mm.config.Logging.MaxBodySize,
		}
		r.Use(LoggerMiddleware(loggerConfig))
	}

	// Metrics
	if mm.config.Metrics.Enabled {
		r.Use(PrometheusMiddleware())
		r.Use(BusinessMetricsMiddleware())

		// Metrics endpoint
		r.GET(mm.config.Metrics.Path, gin.WrapH(promhttp.Handler()))
	}

	// Rate limiting
	if mm.config.RateLimit.Enabled && mm.rateLimiter != nil {
		r.Use(mm.rateLimiter.RateLimitMiddleware())
	}

	// Timeout
	if mm.config.Timeout.Enabled {
		r.Use(TimeoutMiddleware(mm.config.Timeout.Duration))
	}

	// Request size limit
	if mm.config.SizeLimit.Enabled {
		r.Use(SizeLimit(mm.config.SizeLimit.MaxSize))
	}

	// Content type validation for POST/PUT/PATCH
	allowedContentTypes := []string{
		"application/json",
		"application/x-www-form-urlencoded",
		"multipart/form-data",
	}
	r.Use(ContentTypeMiddleware(allowedContentTypes))

	// IP whitelist (if enabled)
	if mm.config.IPWhitelist.Enabled && len(mm.config.IPWhitelist.AllowedIPs) > 0 {
		r.Use(IPWhitelistMiddleware(mm.config.IPWhitelist.AllowedIPs))
	}

	// API key authentication (if enabled)
	if mm.config.APIAuth.Enabled && len(mm.config.APIAuth.ValidKeys) > 0 {
		r.Use(APIKeyAuthMiddleware(mm.config.APIAuth.ValidKeys))
	}

	// Health check (should be after rate limiting but before auth)
	r.Use(HealthCheckMiddleware())

	// 404 and 405 handlers
	r.NoRoute(NoRouteMiddleware())
	r.NoMethod(NoMethodMiddleware())
}

// SetupDevelopmentMiddleware configures middleware for development environment
func (mm *MiddlewareManager) SetupDevelopmentMiddleware(r *gin.Engine) {
	// Request ID
	r.Use(RequestIDMiddleware())

	// Security headers
	if mm.config.Security.Enabled {
		securityConfig := SecurityConfig{
			FrameOptions:      mm.config.Security.FrameOptions,
			ReferrerPolicy:    mm.config.Security.ReferrerPolicy,
			CSP:               mm.config.Security.CSP,
			PermissionsPolicy: mm.config.Security.PermissionsPolicy,
			HSTSMaxAge:        "max-age=31536000; includeSubDomains",
			CustomHeaders:     map[string]string{"X-API-Version": "v1"},
		}
		r.Use(SecureHeaders(securityConfig))
	}

	// Development CORS (allows all origins)
	r.Use(DevCORSMiddleware())

	// Development recovery (with detailed error info)
	r.Use(DevelopmentRecoveryMiddleware(mm.logger))

	// Detailed logging
	loggerConfig := LoggerConfig{
		Logger:          mm.logger,
		SkipPaths:       []string{"/health"},
		LogRequestBody:  true,
		LogResponseBody: true,
		MaxBodySize:     1024 * 1024, // 1MB
	}
	r.Use(LoggerMiddleware(loggerConfig))

	// Metrics
	r.Use(PrometheusMiddleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Lenient rate limiting
	r.Use(GlobalRateLimitMiddleware(1000, 2000))

	// Health check
	r.Use(HealthCheckMiddleware())

	// 404 and 405 handlers
	r.NoRoute(NoRouteMiddleware())
	r.NoMethod(NoMethodMiddleware())
}

// SetupProductionMiddleware configures middleware for production environment
func (mm *MiddlewareManager) SetupProductionMiddleware(r *gin.Engine, allowedOrigins []string) {
	// Request ID
	r.Use(RequestIDMiddleware())

	// Production security headers
	securityConfig := DefaultSecurityConfig()
	r.Use(SecureHeaders(securityConfig))

	// Production CORS
	r.Use(ProductionCORSMiddleware(allowedOrigins))

	// Production recovery (no detailed error info)
	r.Use(ProductionRecoveryMiddleware(mm.logger))

	// Production logging (no request/response bodies)
	loggerConfig := LoggerConfig{
		Logger:          mm.logger,
		SkipPaths:       []string{"/health", "/metrics"},
		LogRequestBody:  false,
		LogResponseBody: false,
		MaxBodySize:     1024 * 1024, // 1MB
	}
	r.Use(LoggerMiddleware(loggerConfig))

	// Metrics (on separate path)
	r.Use(PrometheusMiddleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Strict rate limiting
	if mm.rateLimiter != nil {
		r.Use(mm.rateLimiter.RateLimitMiddleware())
	}

	// Request timeout
	r.Use(TimeoutMiddleware(30 * time.Second))

	// Request size limit
	r.Use(SizeLimit(10 * 1024 * 1024)) // 10MB

	// Content type validation
	allowedContentTypes := []string{"application/json"}
	r.Use(ContentTypeMiddleware(allowedContentTypes))

	// Health check
	r.Use(HealthCheckMiddleware())

	// 404 and 405 handlers
	r.NoRoute(NoRouteMiddleware())
	r.NoMethod(NoMethodMiddleware())
}
