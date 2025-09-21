package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
)

// LoggerConfig represents logger middleware configuration
type LoggerConfig struct {
	Logger          *logger.Logger
	SkipPaths       []string
	LogRequestBody  bool
	LogResponseBody bool
	MaxBodySize     int64
}

// DefaultLoggerConfig returns default logger configuration
func DefaultLoggerConfig(logger *logger.Logger) LoggerConfig {
	return LoggerConfig{
		Logger:          logger,
		SkipPaths:       []string{"/health", "/metrics", "/favicon.ico"},
		LogRequestBody:  false,
		LogResponseBody: false,
		MaxBodySize:     1024 * 1024, // 1MB
	}
}

// RequestIDMiddleware adds request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return requestid.New()
}

// LoggerMiddleware creates a structured logging middleware
func LoggerMiddleware(config LoggerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip logging for certain paths
		for _, path := range config.SkipPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Read request body if enabled
		var requestBody []byte
		if config.LogRequestBody && c.Request.Body != nil {
			requestBody, _ = io.ReadAll(io.LimitReader(c.Request.Body, config.MaxBodySize))
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Custom response writer to capture response body
		var responseBody []byte
		if config.LogResponseBody {
			blw := &bodyLogWriter{
				ResponseWriter: c.Writer,
				body:           &bytes.Buffer{},
				maxSize:        config.MaxBodySize,
			}
			c.Writer = blw
			defer func() {
				responseBody = blw.body.Bytes()
			}()
		}

		// Process request
		c.Next()

		// Calculate request duration
		latency := time.Since(start)

		// Get request ID
		requestID := requestid.Get(c)

		// Build log fields
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.Int("response_size", c.Writer.Size()),
		}

		// Add query parameters if present
		if raw != "" {
			fields = append(fields, zap.String("query", raw))
		}

		// Add request body if enabled
		if config.LogRequestBody && len(requestBody) > 0 {
			fields = append(fields, zap.String("request_body", string(requestBody)))
		}

		// Add response body if enabled
		if config.LogResponseBody && len(responseBody) > 0 {
			fields = append(fields, zap.String("response_body", string(responseBody)))
		}

		// Add error if present
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		// Log based on status code
		status := c.Writer.Status()
		switch {
		case status >= 500:
			config.Logger.Error("HTTP Request - Server Error", fields...)
		case status >= 400:
			config.Logger.Warn("HTTP Request - Client Error", fields...)
		case status >= 300:
			config.Logger.Info("HTTP Request - Redirect", fields...)
		default:
			config.Logger.Info("HTTP Request - Success", fields...)
		}
	}
}

// bodyLogWriter captures response body for logging
type bodyLogWriter struct {
	gin.ResponseWriter
	body    *bytes.Buffer
	maxSize int64
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	// Capture response body up to maxSize
	if w.body.Len() < int(w.maxSize) {
		remaining := int(w.maxSize) - w.body.Len()
		if len(b) <= remaining {
			w.body.Write(b)
		} else {
			w.body.Write(b[:remaining])
		}
	}
	return w.ResponseWriter.Write(b)
}

// AccessLogMiddleware creates a simple access log middleware
func AccessLogMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info("HTTP Access Log",
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.String("ip", param.ClientIP),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("user_agent", param.Request.UserAgent()),
		)
		return ""
	})
}

// ErrorLogMiddleware logs errors with structured logging
func ErrorLogMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Log errors if any
		for _, err := range c.Errors {
			logger.Error("Request Error",
				zap.String("request_id", requestid.Get(c)),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
				zap.Error(err.Err),
				zap.String("type", string(err.Type)),
				zap.Any("meta", err.Meta),
			)
		}
	}
}
