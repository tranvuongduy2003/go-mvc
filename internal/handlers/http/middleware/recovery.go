package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
)

// RecoveryConfig represents recovery middleware configuration
type RecoveryConfig struct {
	Logger           *logger.Logger
	EnableStackTrace bool
	EnablePanic      bool
}

// CustomRecoveryMiddleware creates a recovery middleware with structured logging
func CustomRecoveryMiddleware(config RecoveryConfig) gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, recovered interface{}) {
		// Get request ID
		requestID := requestid.Get(c)

		// Build error message
		err := fmt.Sprintf("Panic recovered: %v", recovered)

		// Get stack trace if enabled
		var stackTrace string
		if config.EnableStackTrace {
			stackTrace = getStackTrace()
		}

		// Log the panic with structured logging
		config.Logger.Error("Panic Recovered",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("error", err),
			zap.String("stack_trace", stackTrace),
		)

		// Return error response
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "Internal server error",
			"message":    "An unexpected error occurred",
			"code":       "INTERNAL_SERVER_ERROR",
			"request_id": requestID,
		})
	})
}

// DefaultRecoveryMiddleware creates a recovery middleware with default configuration
func DefaultRecoveryMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return CustomRecoveryMiddleware(RecoveryConfig{
		Logger:           logger,
		EnableStackTrace: true,
		EnablePanic:      false,
	})
}

// ProductionRecoveryMiddleware creates a recovery middleware for production
func ProductionRecoveryMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return CustomRecoveryMiddleware(RecoveryConfig{
		Logger:           logger,
		EnableStackTrace: false,
		EnablePanic:      false,
	})
}

// DevelopmentRecoveryMiddleware creates a recovery middleware for development
func DevelopmentRecoveryMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				// Get request ID
				requestID := requestid.Get(c)

				// Get stack trace
				stackTrace := getStackTrace()

				// Log the panic
				logger.Error("Panic Recovered (Development)",
					zap.String("request_id", requestID),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("ip", c.ClientIP()),
					zap.String("user_agent", c.Request.UserAgent()),
					zap.Any("panic", recovered),
					zap.String("stack_trace", stackTrace),
				)

				// Return detailed error response for development
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":       "Internal server error",
					"message":     "An unexpected error occurred",
					"code":        "INTERNAL_SERVER_ERROR",
					"request_id":  requestID,
					"panic":       recovered,
					"stack_trace": stackTrace,
				})
			}
		}()

		c.Next()
	}
}

// getStackTrace returns formatted stack trace
func getStackTrace() string {
	buf := make([]byte, 1024*4)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// TimeoutMiddleware adds request timeout
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		// Check if context was cancelled due to timeout
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error":      "Request timeout",
				"message":    "Request processing exceeded maximum allowed time",
				"code":       "REQUEST_TIMEOUT",
				"timeout":    timeout.String(),
				"request_id": requestid.Get(c),
			})
			c.Abort()
		}
	}
}

// HealthCheckMiddleware provides basic health check functionality
func HealthCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": "now",
				"service":   "go-mvc",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// MaintenanceMiddleware returns maintenance mode response
func MaintenanceMiddleware(isMaintenanceMode func() bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isMaintenanceMode() {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "Service unavailable",
				"message": "The service is currently under maintenance",
				"code":    "MAINTENANCE_MODE",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// NoRouteMiddleware handles 404 responses
func NoRouteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":      "Route not found",
			"message":    "The requested route does not exist",
			"code":       "ROUTE_NOT_FOUND",
			"path":       c.Request.URL.Path,
			"method":     c.Request.Method,
			"request_id": requestid.Get(c),
		})
	}
}

// NoMethodMiddleware handles 405 responses
func NoMethodMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error":      "Method not allowed",
			"message":    "The HTTP method is not allowed for this route",
			"code":       "METHOD_NOT_ALLOWED",
			"path":       c.Request.URL.Path,
			"method":     c.Request.Method,
			"request_id": requestid.Get(c),
		})
	}
}
