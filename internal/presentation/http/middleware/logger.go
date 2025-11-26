package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/logger"
)

type LoggerConfig struct {
	Logger          *logger.Logger
	SkipPaths       []string
	LogRequestBody  bool
	LogResponseBody bool
	MaxBodySize     int64
}

func DefaultLoggerConfig(logger *logger.Logger) LoggerConfig {
	return LoggerConfig{
		Logger:          logger,
		SkipPaths:       []string{"/health", "/metrics", "/favicon.ico"},
		LogRequestBody:  false,
		LogResponseBody: false,
		MaxBodySize:     1024 * 1024, // 1MB
	}
}

func RequestIDMiddleware() gin.HandlerFunc {
	return requestid.New()
}

func LoggerMiddleware(config LoggerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, path := range config.SkipPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		var requestBody []byte
		if config.LogRequestBody && c.Request.Body != nil {
			requestBody, _ = io.ReadAll(io.LimitReader(c.Request.Body, config.MaxBodySize))
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

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

		c.Next()

		latency := time.Since(start)

		requestID := requestid.Get(c)

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

		if raw != "" {
			fields = append(fields, zap.String("query", raw))
		}

		if config.LogRequestBody && len(requestBody) > 0 {
			fields = append(fields, zap.String("request_body", string(requestBody)))
		}

		if config.LogResponseBody && len(responseBody) > 0 {
			fields = append(fields, zap.String("response_body", string(responseBody)))
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

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

type bodyLogWriter struct {
	gin.ResponseWriter
	body    *bytes.Buffer
	maxSize int64
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
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

func ErrorLogMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			logger.Error("Request Error",
				zap.String("request_id", requestid.Get(c)),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
				zap.Error(err.Err),
				zap.Uint64("type", uint64(err.Type)),
				zap.Any("meta", err.Meta),
			)
		}
	}
}
