package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TracingMiddleware creates a Gin middleware for distributed tracing
func TracingMiddleware(serviceName string) gin.HandlerFunc {
	return otelgin.Middleware(serviceName)
}

// CustomTracingMiddleware creates a custom tracing middleware with additional attributes
func CustomTracingMiddleware(serviceName string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Use the otelgin middleware first
		otelgin.Middleware(serviceName)(c)

		// Add custom attributes to the span
		span := trace.SpanFromContext(c.Request.Context())
		if span.IsRecording() {
			span.SetAttributes(
				attribute.String("http.client_ip", c.ClientIP()),
				attribute.String("http.user_agent", c.GetHeader("User-Agent")),
				attribute.String("http.request_id", c.GetHeader("X-Request-ID")),
			)

			// Add user information if available
			if userID, exists := c.Get("user_id"); exists {
				if uid, ok := userID.(string); ok {
					span.SetAttributes(attribute.String("user.id", uid))
				}
			}

			if userEmail, exists := c.Get("user_email"); exists {
				if email, ok := userEmail.(string); ok {
					span.SetAttributes(attribute.String("user.email", email))
				}
			}
		}

		c.Next()

		// Add response attributes after processing
		if span.IsRecording() {
			span.SetAttributes(
				attribute.Int("http.status_code", c.Writer.Status()),
				attribute.Int("http.response_size", c.Writer.Size()),
			)
		}
	})
}

// TraceContext extracts tracing context from Gin context
func TraceContext(c *gin.Context) context.Context {
	return c.Request.Context()
}

// WithSpan creates a new span in the context
func WithSpan(ctx context.Context, spanName string, fn func(ctx context.Context, span trace.Span) error) error {
	span := trace.SpanFromContext(ctx)
	tracer := span.TracerProvider().Tracer("gin-middleware")

	ctx, newSpan := tracer.Start(ctx, spanName)
	defer newSpan.End()

	err := fn(ctx, newSpan)
	if err != nil {
		newSpan.RecordError(err)
		newSpan.SetStatus(codes.Error, err.Error())
	}

	return err
}
