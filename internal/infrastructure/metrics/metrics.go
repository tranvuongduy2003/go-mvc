package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/config"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/logger"
)

type Metrics struct {
	httpRequestsTotal     *prometheus.CounterVec
	httpRequestDuration   *prometheus.HistogramVec
	httpRequestSize       *prometheus.HistogramVec
	httpResponseSize      *prometheus.HistogramVec
	databaseConnections   *prometheus.GaugeVec
	databaseQueries       *prometheus.CounterVec
	databaseQueryDuration *prometheus.HistogramVec
	cacheHits             *prometheus.CounterVec
	activeConnections     prometheus.Gauge
	uptimeSeconds         prometheus.Gauge
	memoryUsage           prometheus.Gauge
	cpuUsage              prometheus.Gauge
}

func NewMetrics() *Metrics {
	metrics := &Metrics{
		httpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status_code"},
		),
		httpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
		httpRequestSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_size_bytes",
				Help:    "HTTP request size in bytes",
				Buckets: prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"method", "endpoint"},
		),
		httpResponseSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "HTTP response size in bytes",
				Buckets: prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"method", "endpoint"},
		),
		databaseConnections: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "database_connections",
				Help: "Number of database connections",
			},
			[]string{"database", "state"},
		),
		databaseQueries: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "database_queries_total",
				Help: "Total number of database queries",
			},
			[]string{"database", "table", "operation"},
		),
		databaseQueryDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "database_query_duration_seconds",
				Help:    "Database query duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"database", "table", "operation"},
		),
		cacheHits: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_requests_total",
				Help: "Total number of cache requests",
			},
			[]string{"cache", "result"},
		),
		activeConnections: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "active_connections",
				Help: "Number of active connections",
			},
		),
		uptimeSeconds: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "uptime_seconds",
				Help: "Application uptime in seconds",
			},
		),
		memoryUsage: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "memory_usage_bytes",
				Help: "Memory usage in bytes",
			},
		),
		cpuUsage: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "cpu_usage_percent",
				Help: "CPU usage percentage",
			},
		),
	}

	prometheus.MustRegister(
		metrics.httpRequestsTotal,
		metrics.httpRequestDuration,
		metrics.httpRequestSize,
		metrics.httpResponseSize,
		metrics.databaseConnections,
		metrics.databaseQueries,
		metrics.databaseQueryDuration,
		metrics.cacheHits,
		metrics.activeConnections,
		metrics.uptimeSeconds,
		metrics.memoryUsage,
		metrics.cpuUsage,
	)

	return metrics
}

func (m *Metrics) RecordHTTPRequest(method, endpoint, statusCode string, duration time.Duration, requestSize, responseSize int64) {
	m.httpRequestsTotal.WithLabelValues(method, endpoint, statusCode).Inc()
	m.httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
	m.httpRequestSize.WithLabelValues(method, endpoint).Observe(float64(requestSize))
	m.httpResponseSize.WithLabelValues(method, endpoint).Observe(float64(responseSize))
}

func (m *Metrics) RecordDatabaseQuery(database, table, operation string, duration time.Duration) {
	m.databaseQueries.WithLabelValues(database, table, operation).Inc()
	m.databaseQueryDuration.WithLabelValues(database, table, operation).Observe(duration.Seconds())
}

func (m *Metrics) SetDatabaseConnections(database, state string, count float64) {
	m.databaseConnections.WithLabelValues(database, state).Set(count)
}

func (m *Metrics) RecordCacheHit(cache, result string) {
	m.cacheHits.WithLabelValues(cache, result).Inc()
}

func (m *Metrics) SetActiveConnections(count float64) {
	m.activeConnections.Set(count)
}

func (m *Metrics) SetUptime(seconds float64) {
	m.uptimeSeconds.Set(seconds)
}

func (m *Metrics) SetMemoryUsage(bytes float64) {
	m.memoryUsage.Set(bytes)
}

func (m *Metrics) SetCPUUsage(percent float64) {
	m.cpuUsage.Set(percent)
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.Handler()
}

func (m *Metrics) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		statusCode := string(rune(c.Writer.Status()))

		requestSize := c.Request.ContentLength
		if requestSize < 0 {
			requestSize = 0
		}
		responseSize := int64(c.Writer.Size())

		m.RecordHTTPRequest(
			c.Request.Method,
			c.FullPath(),
			statusCode,
			duration,
			requestSize,
			responseSize,
		)
	}
}

type Manager struct {
	metrics *Metrics
	config  config.Metrics
	logger  *logger.Logger
	server  *http.Server
}

func NewManager(cfg config.Metrics, log *logger.Logger) *Manager {
	return &Manager{
		metrics: NewMetrics(),
		config:  cfg,
		logger:  log,
	}
}

func (m *Manager) Start() error {
	if !m.config.Enabled {
		m.logger.Info("Metrics collection is disabled")
		return nil
	}

	mux := http.NewServeMux()
	mux.Handle(m.config.Path, m.metrics.Handler())

	m.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", m.config.Port),
		Handler: mux,
	}

	go func() {
		m.logger.Infof("Starting metrics server on port %d", m.config.Port)
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			m.logger.Errorf("Failed to start metrics server: %v", err)
		}
	}()

	return nil
}

func (m *Manager) Stop() error {
	if m.server != nil {
		m.logger.Info("Stopping metrics server")
		return m.server.Close()
	}
	return nil
}

func (m *Manager) Metrics() *Metrics {
	return m.metrics
}

func (m *Manager) GinMiddleware() gin.HandlerFunc {
	return m.metrics.GinMiddleware()
}
