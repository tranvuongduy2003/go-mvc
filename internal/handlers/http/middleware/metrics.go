package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
	)

	httpRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: []float64{100, 1000, 10000, 100000, 1000000},
		},
		[]string{"method", "endpoint"},
	)

	httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: []float64{100, 1000, 10000, 100000, 1000000},
		},
		[]string{"method", "endpoint"},
	)

	// Application metrics
	activeConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)

	databaseConnections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "database_connections",
			Help: "Number of database connections",
		},
		[]string{"state"}, // open, in_use, idle
	)

	cacheOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_operations_total",
			Help: "Total number of cache operations",
		},
		[]string{"operation", "result"}, // get/set/delete, hit/miss/error
	)
)

// PrometheusMiddleware collects HTTP metrics for Prometheus
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Increment in-flight requests
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		// Record request size
		if c.Request.ContentLength > 0 {
			httpRequestSize.WithLabelValues(
				c.Request.Method,
				c.FullPath(),
			).Observe(float64(c.Request.ContentLength))
		}

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get response size
		responseSize := c.Writer.Size()
		if responseSize > 0 {
			httpResponseSize.WithLabelValues(
				c.Request.Method,
				c.FullPath(),
			).Observe(float64(responseSize))
		}

		// Record metrics
		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			strconv.Itoa(c.Writer.Status()),
		).Inc()

		httpRequestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
		).Observe(duration)
	}
}

// RecordDatabaseConnection records database connection metrics
func RecordDatabaseConnection(state string, count float64) {
	databaseConnections.WithLabelValues(state).Set(count)
}

// RecordCacheOperation records cache operation metrics
func RecordCacheOperation(operation, result string) {
	cacheOperations.WithLabelValues(operation, result).Inc()
}

// RecordActiveConnections records active connection count
func RecordActiveConnections(count float64) {
	activeConnections.Set(count)
}

// CustomMetrics allows applications to register custom metrics
type CustomMetrics struct {
	Registry *prometheus.Registry
}

// NewCustomMetrics creates a new custom metrics registry
func NewCustomMetrics() *CustomMetrics {
	return &CustomMetrics{
		Registry: prometheus.NewRegistry(),
	}
}

// MustRegister registers metrics with the custom registry
func (cm *CustomMetrics) MustRegister(collectors ...prometheus.Collector) {
	cm.Registry.MustRegister(collectors...)
}

// BusinessMetricsMiddleware collects business-specific metrics
func BusinessMetricsMiddleware() gin.HandlerFunc {
	userActions := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_actions_total",
			Help: "Total number of user actions",
		},
		[]string{"action", "user_id", "status"},
	)

	return func(c *gin.Context) {
		c.Next()

		// Extract user ID from context if available
		userID := "anonymous"
		if uid, exists := c.Get("user_id"); exists {
			if id, ok := uid.(string); ok {
				userID = id
			}
		}

		// Record user action based on endpoint
		action := "unknown"
		switch c.FullPath() {
		case "/api/v1/users":
			if c.Request.Method == "POST" {
				action = "user_create"
			} else if c.Request.Method == "GET" {
				action = "user_list"
			}
		case "/api/v1/users/:id":
			switch c.Request.Method {
			case "GET":
				action = "user_get"
			case "PUT", "PATCH":
				action = "user_update"
			case "DELETE":
				action = "user_delete"
			}
		case "/api/v1/auth/login":
			action = "auth_login"
		case "/api/v1/auth/register":
			action = "auth_register"
		case "/api/v1/auth/logout":
			action = "auth_logout"
		}

		if action != "unknown" {
			status := "success"
			if c.Writer.Status() >= 400 {
				status = "error"
			}

			userActions.WithLabelValues(action, userID, status).Inc()
		}
	}
}
