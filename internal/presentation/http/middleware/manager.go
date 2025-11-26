package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/logger"
	"github.com/tranvuongduy2003/go-mvc/pkg/jwt"
)

type MiddlewareManager struct {
	logger      *logger.Logger
	rateLimiter *RateLimiter
	config      MiddlewareConfig
}

type MiddlewareConfig struct {
	ServiceName string

	RateLimit struct {
		Enabled bool
		RPS     int
		Burst   int
	}

	CORS struct {
		Enabled        bool
		AllowedOrigins []string
		AllowedMethods []string
		AllowedHeaders []string
	}

	Security struct {
		Enabled           bool
		FrameOptions      string
		ReferrerPolicy    string
		CSP               string
		PermissionsPolicy string
	}

	Logging struct {
		Enabled         bool
		SkipPaths       []string
		LogRequestBody  bool
		LogResponseBody bool
		MaxBodySize     int64
	}

	Metrics struct {
		Enabled bool
		Path    string
	}

	Timeout struct {
		Enabled  bool
		Duration time.Duration
	}

	APIAuth struct {
		Enabled    bool
		ValidKeys  map[string]string
		HeaderName string
	}

	IPWhitelist struct {
		Enabled    bool
		AllowedIPs []string
	}

	SizeLimit struct {
		Enabled bool
		MaxSize int64
	}
}

func DefaultMiddlewareConfig() MiddlewareConfig {
	return MiddlewareConfig{
		ServiceName: "go-mvc",
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

func (mm *MiddlewareManager) SetupMiddleware(r *gin.Engine) {
	r.Use(RequestIDMiddleware())

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

	if mm.config.CORS.Enabled {
		if len(mm.config.CORS.AllowedOrigins) == 1 && mm.config.CORS.AllowedOrigins[0] == "*" {
			r.Use(DevCORSMiddleware())
		} else {
			r.Use(ProductionCORSMiddleware(mm.config.CORS.AllowedOrigins))
		}
	}

	r.Use(DefaultRecoveryMiddleware(mm.logger))

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

	if mm.config.Metrics.Enabled {
		r.Use(PrometheusMiddleware())
		r.Use(BusinessMetricsMiddleware())

		r.GET(mm.config.Metrics.Path, gin.WrapH(promhttp.Handler()))
	}

	if mm.config.RateLimit.Enabled && mm.rateLimiter != nil {
		r.Use(mm.rateLimiter.RateLimitMiddleware())
	}

	if mm.config.Timeout.Enabled {
		r.Use(TimeoutMiddleware(mm.config.Timeout.Duration))
	}

	if mm.config.SizeLimit.Enabled {
		r.Use(SizeLimit(mm.config.SizeLimit.MaxSize))
	}

	allowedContentTypes := []string{
		"application/json",
		"application/x-www-form-urlencoded",
		"multipart/form-data",
	}
	r.Use(ContentTypeMiddleware(allowedContentTypes))

	if mm.config.IPWhitelist.Enabled && len(mm.config.IPWhitelist.AllowedIPs) > 0 {
		r.Use(IPWhitelistMiddleware(mm.config.IPWhitelist.AllowedIPs))
	}

	if mm.config.APIAuth.Enabled && len(mm.config.APIAuth.ValidKeys) > 0 {
		r.Use(APIKeyAuthMiddleware(mm.config.APIAuth.ValidKeys))
	}

	r.Use(HealthCheckMiddleware())

	r.NoRoute(NoRouteMiddleware())
	r.NoMethod(NoMethodMiddleware())
}

func (mm *MiddlewareManager) SetupDevelopmentMiddleware(r *gin.Engine) {
	r.Use(RequestIDMiddleware())

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

	r.Use(DevCORSMiddleware())

	r.Use(CustomTracingMiddleware(mm.config.ServiceName))

	r.Use(DevelopmentRecoveryMiddleware(mm.logger))

	loggerConfig := LoggerConfig{
		Logger:          mm.logger,
		SkipPaths:       []string{"/health"},
		LogRequestBody:  true,
		LogResponseBody: true,
		MaxBodySize:     1024 * 1024, // 1MB
	}
	r.Use(LoggerMiddleware(loggerConfig))

	r.Use(PrometheusMiddleware())
	r.Use(BusinessMetricsMiddleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Use(GlobalRateLimitMiddleware(1000, 2000))

	r.Use(HealthCheckMiddleware())

	r.NoRoute(NoRouteMiddleware())
	r.NoMethod(NoMethodMiddleware())
}

func (mm *MiddlewareManager) SetupProductionMiddleware(r *gin.Engine, allowedOrigins []string) {
	r.Use(RequestIDMiddleware())

	securityConfig := DefaultSecurityConfig()
	r.Use(SecureHeaders(securityConfig))

	r.Use(ProductionCORSMiddleware(allowedOrigins))

	r.Use(CustomTracingMiddleware(mm.config.ServiceName))

	r.Use(ProductionRecoveryMiddleware(mm.logger))

	loggerConfig := LoggerConfig{
		Logger:          mm.logger,
		SkipPaths:       []string{"/health", "/metrics"},
		LogRequestBody:  false,
		LogResponseBody: false,
		MaxBodySize:     1024 * 1024, // 1MB
	}
	r.Use(LoggerMiddleware(loggerConfig))

	r.Use(PrometheusMiddleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	if mm.rateLimiter != nil {
		r.Use(mm.rateLimiter.RateLimitMiddleware())
	}

	r.Use(TimeoutMiddleware(30 * time.Second))

	r.Use(SizeLimit(10 * 1024 * 1024)) // 10MB

	allowedContentTypes := []string{"application/json"}
	r.Use(ContentTypeMiddleware(allowedContentTypes))

	r.Use(HealthCheckMiddleware())

	r.NoRoute(NoRouteMiddleware())
	r.NoMethod(NoMethodMiddleware())
}
