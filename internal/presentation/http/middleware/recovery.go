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

	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/logger"
)

type RecoveryConfig struct {
	Logger           *logger.Logger
	EnableStackTrace bool
	EnablePanic      bool
}

func CustomRecoveryMiddleware(config RecoveryConfig) gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, recovered interface{}) {
		requestID := requestid.Get(c)

		err := fmt.Sprintf("Panic recovered: %v", recovered)

		var stackTrace string
		if config.EnableStackTrace {
			stackTrace = getStackTrace()
		}

		config.Logger.Error("Panic Recovered",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("error", err),
			zap.String("stack_trace", stackTrace),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "Internal server error",
			"message":    "An unexpected error occurred",
			"code":       "INTERNAL_SERVER_ERROR",
			"request_id": requestID,
		})
	})
}

func DefaultRecoveryMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return CustomRecoveryMiddleware(RecoveryConfig{
		Logger:           logger,
		EnableStackTrace: true,
		EnablePanic:      false,
	})
}

func ProductionRecoveryMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return CustomRecoveryMiddleware(RecoveryConfig{
		Logger:           logger,
		EnableStackTrace: false,
		EnablePanic:      false,
	})
}

func DevelopmentRecoveryMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				requestID := requestid.Get(c)

				stackTrace := getStackTrace()

				logger.Error("Panic Recovered (Development)",
					zap.String("request_id", requestID),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("ip", c.ClientIP()),
					zap.String("user_agent", c.Request.UserAgent()),
					zap.Any("panic", recovered),
					zap.String("stack_trace", stackTrace),
				)

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

func getStackTrace() string {
	buf := make([]byte, 1024*4)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()

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
