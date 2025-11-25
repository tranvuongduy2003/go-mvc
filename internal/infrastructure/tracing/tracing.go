package tracing

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type TracingService struct {
	tracer   trace.Tracer
	provider *sdktrace.TracerProvider
}

func NewTracingService(cfg *config.AppConfig) (*TracingService, error) {
	// Skip tracing if disabled
	if !cfg.Tracing.Enabled {
		log.Println("Tracing is disabled")
		provider := sdktrace.NewTracerProvider()
		return &TracingService{
			tracer:   provider.Tracer(cfg.App.Name),
			provider: provider,
		}, nil
	}

	// Create resource with service information
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.App.Name),
			semconv.ServiceVersion(cfg.App.Version),
			semconv.DeploymentEnvironment(cfg.App.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create OTLP HTTP exporter
	var exporter sdktrace.SpanExporter
	endpoint := cfg.Tracing.Endpoint
	if endpoint == "" {
		endpoint = "http://localhost:4318"
	}

	exporter, err = otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		log.Printf("Warning: Failed to create OTLP exporter: %v. Tracing will be disabled.", err)
		// Create a no-op tracer provider
		provider := sdktrace.NewTracerProvider()
		return &TracingService{
			tracer:   provider.Tracer(cfg.App.Name),
			provider: provider,
		}, nil
	}

	// Determine sampling rate
	sampleRate := cfg.Tracing.SampleRate
	if sampleRate <= 0 {
		sampleRate = 1.0 // Default to 100% sampling
	}

	// Create sampler
	var sampler sdktrace.Sampler
	if sampleRate >= 1.0 {
		sampler = sdktrace.AlwaysSample()
	} else {
		sampler = sdktrace.TraceIDRatioBased(sampleRate)
	}

	// Create tracer provider
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	// Set global tracer provider
	otel.SetTracerProvider(provider)

	// Set global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Create tracer
	tracer := provider.Tracer(cfg.App.Name)

	return &TracingService{
		tracer:   tracer,
		provider: provider,
	}, nil
}

// StartSpan starts a new span with the given name
func (t *TracingService) StartSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, spanName, opts...)
}

// AddSpanAttributes adds attributes to the current span
func (t *TracingService) AddSpanAttributes(span trace.Span, attrs ...attribute.KeyValue) {
	span.SetAttributes(attrs...)
}

// AddSpanEvent adds an event to the current span
func (t *TracingService) AddSpanEvent(span trace.Span, name string, attrs ...attribute.KeyValue) {
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// RecordError records an error on the span
func (t *TracingService) RecordError(span trace.Span, err error) {
	span.RecordError(err)
}

// SetSpanStatus sets the status of the span
func (t *TracingService) SetSpanStatus(span trace.Span, code codes.Code, description string) {
	span.SetStatus(code, description)
}

// StartHTTPSpan starts a span for HTTP requests
func (t *TracingService) StartHTTPSpan(ctx context.Context, method, path string) (context.Context, trace.Span) {
	spanName := fmt.Sprintf("%s %s", method, path)
	ctx, span := t.tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer))

	span.SetAttributes(
		attribute.String("http.method", method),
		attribute.String("http.route", path),
	)

	return ctx, span
}

// StartDatabaseSpan starts a span for database operations
func (t *TracingService) StartDatabaseSpan(ctx context.Context, operation, tableName string) (context.Context, trace.Span) {
	spanName := fmt.Sprintf("db.%s", operation)
	ctx, span := t.tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindClient))

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.operation", operation),
		attribute.String("db.table.name", tableName),
	)

	return ctx, span
}

// StartServiceSpan starts a span for service operations
func (t *TracingService) StartServiceSpan(ctx context.Context, serviceName, operation string) (context.Context, trace.Span) {
	spanName := fmt.Sprintf("%s.%s", serviceName, operation)
	ctx, span := t.tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindInternal))

	span.SetAttributes(
		attribute.String("service.name", serviceName),
		attribute.String("service.operation", operation),
	)

	return ctx, span
}

// StartExternalServiceSpan starts a span for external service calls
func (t *TracingService) StartExternalServiceSpan(ctx context.Context, serviceName, operation string) (context.Context, trace.Span) {
	spanName := fmt.Sprintf("external.%s.%s", serviceName, operation)
	ctx, span := t.tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindClient))

	span.SetAttributes(
		attribute.String("external.service.name", serviceName),
		attribute.String("external.service.operation", operation),
	)

	return ctx, span
}

// Middleware creates a tracing middleware for HTTP handlers
func (t *TracingService) Middleware() func(ctx context.Context, next func(context.Context) error) error {
	return func(ctx context.Context, next func(context.Context) error) error {
		start := time.Now()

		err := next(ctx)

		duration := time.Since(start)

		// Add timing information if there's an active span
		if span := trace.SpanFromContext(ctx); span.IsRecording() {
			span.SetAttributes(
				attribute.Int64("duration_ms", duration.Milliseconds()),
			)

			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "")
			}
		}

		return err
	}
}

// Shutdown shuts down the tracer provider
func (t *TracingService) Shutdown(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}

// GetTracer returns the tracer instance
func (t *TracingService) GetTracer() trace.Tracer {
	return t.tracer
}

// Common attribute keys
var (
	UserIDKey        = attribute.Key("user.id")
	UserEmailKey     = attribute.Key("user.email")
	RequestIDKey     = attribute.Key("request.id")
	CorrelationIDKey = attribute.Key("correlation.id")
)
