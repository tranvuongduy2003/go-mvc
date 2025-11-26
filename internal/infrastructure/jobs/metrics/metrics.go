package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/job"
)

// JobMetricsCollector implements the JobMetrics interface using Prometheus
type JobMetricsCollector struct {
	// Counters
	jobsEnqueued  *prometheus.CounterVec
	jobsProcessed *prometheus.CounterVec
	jobsRetries   *prometheus.CounterVec

	// Gauges
	queueSize     *prometheus.GaugeVec
	activeWorkers prometheus.Gauge

	// Histograms
	jobDuration *prometheus.HistogramVec

	// Summary for processing times
	processingSummary *prometheus.SummaryVec
}

// NewJobMetricsCollector creates a new job metrics collector
func NewJobMetricsCollector() *JobMetricsCollector {
	return &JobMetricsCollector{
		jobsEnqueued: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "go_mvc",
				Subsystem: "jobs",
				Name:      "enqueued_total",
				Help:      "Total number of jobs enqueued",
			},
			[]string{"job_type", "queue", "priority"},
		),

		jobsProcessed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "go_mvc",
				Subsystem: "jobs",
				Name:      "processed_total",
				Help:      "Total number of jobs processed",
			},
			[]string{"job_type", "queue", "status"},
		),

		jobsRetries: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "go_mvc",
				Subsystem: "jobs",
				Name:      "retries_total",
				Help:      "Total number of job retries",
			},
			[]string{"job_type", "queue"},
		),

		queueSize: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "go_mvc",
				Subsystem: "jobs",
				Name:      "queue_size",
				Help:      "Current size of job queues",
			},
			[]string{"queue", "priority"},
		),

		activeWorkers: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "go_mvc",
				Subsystem: "jobs",
				Name:      "active_workers",
				Help:      "Number of active job workers",
			},
		),

		jobDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "go_mvc",
				Subsystem: "jobs",
				Name:      "duration_seconds",
				Help:      "Time taken to process jobs",
				Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15), // 1ms to ~32s
			},
			[]string{"job_type", "queue"},
		),

		processingSummary: promauto.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace: "go_mvc",
				Subsystem: "jobs",
				Name:      "processing_duration_seconds",
				Help:      "Summary of job processing duration",
				Objectives: map[float64]float64{
					0.5:  0.05,  // 50th percentile with 5% error
					0.9:  0.01,  // 90th percentile with 1% error
					0.95: 0.005, // 95th percentile with 0.5% error
					0.99: 0.001, // 99th percentile with 0.1% error
				},
			},
			[]string{"job_type", "queue"},
		),
	}
}

// IncrementJobsEnqueued increments the counter for enqueued jobs
func (m *JobMetricsCollector) IncrementJobsEnqueued(jobType string) {
	m.jobsEnqueued.WithLabelValues(jobType, "default", "normal").Inc()
}

// IncrementJobsEnqueuedWithLabels increments the counter with custom labels
func (m *JobMetricsCollector) IncrementJobsEnqueuedWithLabels(jobType, queue string, priority job.JobPriority) {
	priorityStr := m.priorityToString(priority)
	m.jobsEnqueued.WithLabelValues(jobType, queue, priorityStr).Inc()
}

// IncrementJobsProcessed increments the counter for processed jobs
func (m *JobMetricsCollector) IncrementJobsProcessed(jobType string, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	m.jobsProcessed.WithLabelValues(jobType, "default", status).Inc()
}

// IncrementJobsProcessedWithLabels increments the counter with custom labels
func (m *JobMetricsCollector) IncrementJobsProcessedWithLabels(jobType, queue string, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	m.jobsProcessed.WithLabelValues(jobType, queue, status).Inc()
}

// ObserveJobDuration records the time taken to process a job
func (m *JobMetricsCollector) ObserveJobDuration(jobType string, duration time.Duration) {
	seconds := duration.Seconds()
	m.jobDuration.WithLabelValues(jobType, "default").Observe(seconds)
	m.processingSummary.WithLabelValues(jobType, "default").Observe(seconds)
}

// ObserveJobDurationWithLabels records job duration with custom labels
func (m *JobMetricsCollector) ObserveJobDurationWithLabels(jobType, queue string, duration time.Duration) {
	seconds := duration.Seconds()
	m.jobDuration.WithLabelValues(jobType, queue).Observe(seconds)
	m.processingSummary.WithLabelValues(jobType, queue).Observe(seconds)
}

// SetQueueSize sets the current size of a queue
func (m *JobMetricsCollector) SetQueueSize(queue string, size int64) {
	m.queueSize.WithLabelValues(queue, "all").Set(float64(size))
}

// SetQueueSizeWithPriority sets the queue size for a specific priority
func (m *JobMetricsCollector) SetQueueSizeWithPriority(queue string, priority job.JobPriority, size int64) {
	priorityStr := m.priorityToString(priority)
	m.queueSize.WithLabelValues(queue, priorityStr).Set(float64(size))
}

// IncrementJobRetries increments retry counter
func (m *JobMetricsCollector) IncrementJobRetries(jobType string) {
	m.jobsRetries.WithLabelValues(jobType, "default").Inc()
}

// IncrementJobRetriesWithLabels increments retry counter with custom labels
func (m *JobMetricsCollector) IncrementJobRetriesWithLabels(jobType, queue string) {
	m.jobsRetries.WithLabelValues(jobType, queue).Inc()
}

// SetActiveWorkers sets the number of active workers
func (m *JobMetricsCollector) SetActiveWorkers(count int) {
	m.activeWorkers.Set(float64(count))
}

// Helper methods

func (m *JobMetricsCollector) priorityToString(priority job.JobPriority) string {
	switch priority {
	case job.PriorityLow:
		return "low"
	case job.PriorityNormal:
		return "normal"
	case job.PriorityHigh:
		return "high"
	case job.PriorityCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// JobMetricsMiddleware wraps job handlers to automatically collect metrics
type JobMetricsMiddleware struct {
	next    job.JobHandler
	metrics *JobMetricsCollector
}

// NewJobMetricsMiddleware creates a new metrics middleware
func NewJobMetricsMiddleware(handler job.JobHandler, metrics *JobMetricsCollector) *JobMetricsMiddleware {
	return &JobMetricsMiddleware{
		next:    handler,
		metrics: metrics,
	}
}

// Execute executes the job and collects metrics
func (m *JobMetricsMiddleware) Execute(ctx context.Context, job job.Job) error {
	start := time.Now()

	// Execute the job
	err := m.next.Execute(ctx, job)

	// Record metrics
	duration := time.Since(start)
	m.metrics.ObserveJobDuration(job.GetType(), duration)

	success := err == nil
	m.metrics.IncrementJobsProcessed(job.GetType(), success)

	if !success {
		m.metrics.IncrementJobRetries(job.GetType())
	}

	return err
}

// GetJobType returns the job type this handler can process
func (m *JobMetricsMiddleware) GetJobType() string {
	return m.next.GetJobType()
}

// JobQueueMetricsCollector periodically collects queue metrics
type JobQueueMetricsCollector struct {
	queue   job.JobQueue
	metrics *JobMetricsCollector
	queues  []string
}

// NewJobQueueMetricsCollector creates a new queue metrics collector
func NewJobQueueMetricsCollector(queue job.JobQueue, metrics *JobMetricsCollector, queues []string) *JobQueueMetricsCollector {
	if len(queues) == 0 {
		queues = []string{"default"}
	}

	return &JobQueueMetricsCollector{
		queue:   queue,
		metrics: metrics,
		queues:  queues,
	}
}

// CollectQueueMetrics collects current queue size metrics
func (c *JobQueueMetricsCollector) CollectQueueMetrics(ctx context.Context) error {
	for _, queueName := range c.queues {
		size, err := c.queue.GetQueueSize(ctx, queueName)
		if err != nil {
			// Log error but continue with other queues
			continue
		}

		c.metrics.SetQueueSize(queueName, size)
	}

	return nil
}

// WorkerPoolMetricsCollector collects worker pool metrics
type WorkerPoolMetricsCollector struct {
	metrics *JobMetricsCollector
}

// NewWorkerPoolMetricsCollector creates a new worker pool metrics collector
func NewWorkerPoolMetricsCollector(metrics *JobMetricsCollector) *WorkerPoolMetricsCollector {
	return &WorkerPoolMetricsCollector{
		metrics: metrics,
	}
}

// UpdateWorkerCount updates the active worker count metric
func (c *WorkerPoolMetricsCollector) UpdateWorkerCount(count int) {
	c.metrics.SetActiveWorkers(count)
}

// JobTypeMetrics provides job-type specific metrics
type JobTypeMetrics struct {
	JobType        string  `json:"job_type"`
	TotalEnqueued  int64   `json:"total_enqueued"`
	TotalProcessed int64   `json:"total_processed"`
	TotalRetries   int64   `json:"total_retries"`
	SuccessRate    float64 `json:"success_rate"`
	AvgDuration    float64 `json:"avg_duration_seconds"`
}

// QueueMetrics provides queue-specific metrics
type QueueMetrics struct {
	Queue       string `json:"queue"`
	CurrentSize int64  `json:"current_size"`
	Priority    string `json:"priority"`
}

// WorkerMetrics provides worker-specific metrics
type WorkerMetrics struct {
	ActiveWorkers int `json:"active_workers"`
}

// OverallMetrics provides overall system metrics
type OverallMetrics struct {
	JobTypes []JobTypeMetrics `json:"job_types"`
	Queues   []QueueMetrics   `json:"queues"`
	Workers  WorkerMetrics    `json:"workers"`
}

// MetricsReporter provides methods to generate metric reports
type MetricsReporter struct {
	collector *JobMetricsCollector
}

// NewMetricsReporter creates a new metrics reporter
func NewMetricsReporter(collector *JobMetricsCollector) *MetricsReporter {
	return &MetricsReporter{
		collector: collector,
	}
}

// GetOverallMetrics returns overall system metrics
func (r *MetricsReporter) GetOverallMetrics() (*OverallMetrics, error) {
	// This is a simplified implementation
	// In practice, you would query Prometheus or maintain internal counters
	return &OverallMetrics{
		JobTypes: []JobTypeMetrics{},
		Queues:   []QueueMetrics{},
		Workers:  WorkerMetrics{},
	}, nil
}

// Custom metrics for specific business needs

// BusinessMetricsCollector provides business-specific metrics
type BusinessMetricsCollector struct {
	emailJobsCounter    *prometheus.CounterVec
	fileProcessingGauge *prometheus.GaugeVec
	cleanupJobsDuration *prometheus.HistogramVec
}

// NewBusinessMetricsCollector creates a new business metrics collector
func NewBusinessMetricsCollector() *BusinessMetricsCollector {
	return &BusinessMetricsCollector{
		emailJobsCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "go_mvc",
				Subsystem: "business",
				Name:      "email_jobs_total",
				Help:      "Total number of email jobs processed",
			},
			[]string{"email_type", "status"},
		),

		fileProcessingGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "go_mvc",
				Subsystem: "business",
				Name:      "file_processing_active",
				Help:      "Number of active file processing jobs",
			},
			[]string{"file_type"},
		),

		cleanupJobsDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "go_mvc",
				Subsystem: "business",
				Name:      "cleanup_duration_seconds",
				Help:      "Time taken for cleanup jobs",
				Buckets:   prometheus.ExponentialBuckets(0.1, 2, 12), // 100ms to ~6 minutes
			},
			[]string{"cleanup_type"},
		),
	}
}

// RecordEmailJob records email job metrics
func (b *BusinessMetricsCollector) RecordEmailJob(emailType string, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	b.emailJobsCounter.WithLabelValues(emailType, status).Inc()
}

// SetActiveFileProcessingJobs sets the number of active file processing jobs
func (b *BusinessMetricsCollector) SetActiveFileProcessingJobs(fileType string, count int) {
	b.fileProcessingGauge.WithLabelValues(fileType).Set(float64(count))
}

// RecordCleanupDuration records cleanup job duration
func (b *BusinessMetricsCollector) RecordCleanupDuration(cleanupType string, duration time.Duration) {
	b.cleanupJobsDuration.WithLabelValues(cleanupType).Observe(duration.Seconds())
}

// MetricsConfiguration holds metrics system configuration
type MetricsConfiguration struct {
	Enabled            bool          `yaml:"enabled" env:"METRICS_ENABLED" env-default:"true"`
	CollectionInterval time.Duration `yaml:"collection_interval" env:"METRICS_COLLECTION_INTERVAL" env-default:"30s"`
	RetentionPeriod    time.Duration `yaml:"retention_period" env:"METRICS_RETENTION_PERIOD" env-default:"24h"`
}

// MetricsServer provides HTTP endpoints for metrics
type MetricsServer struct {
	collector *JobMetricsCollector
	reporter  *MetricsReporter
}

// NewMetricsServer creates a new metrics server
func NewMetricsServer(collector *JobMetricsCollector) *MetricsServer {
	return &MetricsServer{
		collector: collector,
		reporter:  NewMetricsReporter(collector),
	}
}

// GetMetrics returns current metrics as JSON
func (s *MetricsServer) GetMetrics() (*OverallMetrics, error) {
	return s.reporter.GetOverallMetrics()
}
