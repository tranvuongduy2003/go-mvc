package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSConfig represents CORS configuration
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig returns default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"X-Request-ID",
			"X-API-Key",
			"Accept",
			"Accept-Encoding",
			"Accept-Language",
			"Cache-Control",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"X-Request-ID",
			"X-RateLimit-Limit",
			"X-RateLimit-Remaining",
			"X-RateLimit-Reset",
		},
		AllowCredentials: false,
		MaxAge:           300, // 5 minutes
	}
}

// ProductionCORSConfig returns production-safe CORS configuration
func ProductionCORSConfig(allowedOrigins []string) CORSConfig {
	return CORSConfig{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"X-Request-ID",
			"X-API-Key",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"X-Request-ID",
			"X-RateLimit-Limit",
			"X-RateLimit-Remaining",
			"X-RateLimit-Reset",
		},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}
}

// CORSMiddleware creates CORS middleware with custom configuration
func CORSMiddleware(config CORSConfig) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowOrigins:     config.AllowOrigins,
		AllowMethods:     config.AllowMethods,
		AllowHeaders:     config.AllowHeaders,
		ExposeHeaders:    config.ExposeHeaders,
		AllowCredentials: config.AllowCredentials,
		MaxAge:           time.Duration(config.MaxAge) * time.Second,
	}

	return cors.New(corsConfig)
}

// DefaultCORSMiddleware creates CORS middleware with default configuration
func DefaultCORSMiddleware() gin.HandlerFunc {
	return CORSMiddleware(DefaultCORSConfig())
}

// ProductionCORSMiddleware creates CORS middleware for production
func ProductionCORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return CORSMiddleware(ProductionCORSConfig(allowedOrigins))
}

// DevCORSMiddleware creates CORS middleware for development (allows all origins)
func DevCORSMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowHeaders = []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
		"X-Requested-With",
		"X-Request-ID",
		"X-API-Key",
		"Accept",
		"Accept-Encoding",
		"Accept-Language",
		"Cache-Control",
	}
	config.ExposeHeaders = []string{
		"Content-Length",
		"X-Request-ID",
		"X-RateLimit-Limit",
		"X-RateLimit-Remaining",
		"X-RateLimit-Reset",
	}

	return cors.New(config)
}
