package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/job"
)

type JobMetricsCollector struct {
	jobsEnqueued  *prometheus.CounterVec
	jobsProcessed *prometheus.CounterVec
	jobsRetries   *prometheus.CounterVec

	queueSize     *prometheus.GaugeVec
	activeWorkers prometheus.Gauge

	jobDuration *prometheus.HistogramVec

	processingSummary *prometheus.SummaryVec
}

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

func (m *JobMetricsCollector) IncrementJobsEnqueued(jobType string) {
	m.jobsEnqueued.WithLabelValues(jobType, "default", "normal").Inc()
}

func (m *JobMetricsCollector) IncrementJobsEnqueuedWithLabels(jobType, queue string, priority job.JobPriority) {
	priorityStr := m.priorityToString(priority)
	m.jobsEnqueued.WithLabelValues(jobType, queue, priorityStr).Inc()
}

func (m *JobMetricsCollector) IncrementJobsProcessed(jobType string, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	m.jobsProcessed.WithLabelValues(jobType, "default", status).Inc()
}

func (m *JobMetricsCollector) IncrementJobsProcessedWithLabels(jobType, queue string, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	m.jobsProcessed.WithLabelValues(jobType, queue, status).Inc()
}

func (m *JobMetricsCollector) ObserveJobDuration(jobType string, duration time.Duration) {
	seconds := duration.Seconds()
	m.jobDuration.WithLabelValues(jobType, "default").Observe(seconds)
	m.processingSummary.WithLabelValues(jobType, "default").Observe(seconds)
}

func (m *JobMetricsCollector) ObserveJobDurationWithLabels(jobType, queue string, duration time.Duration) {
	seconds := duration.Seconds()
	m.jobDuration.WithLabelValues(jobType, queue).Observe(seconds)
	m.processingSummary.WithLabelValues(jobType, queue).Observe(seconds)
}

func (m *JobMetricsCollector) SetQueueSize(queue string, size int64) {
	m.queueSize.WithLabelValues(queue, "all").Set(float64(size))
}

func (m *JobMetricsCollector) SetQueueSizeWithPriority(queue string, priority job.JobPriority, size int64) {
	priorityStr := m.priorityToString(priority)
	m.queueSize.WithLabelValues(queue, priorityStr).Set(float64(size))
}

func (m *JobMetricsCollector) IncrementJobRetries(jobType string) {
	m.jobsRetries.WithLabelValues(jobType, "default").Inc()
}

func (m *JobMetricsCollector) IncrementJobRetriesWithLabels(jobType, queue string) {
	m.jobsRetries.WithLabelValues(jobType, queue).Inc()
}

func (m *JobMetricsCollector) SetActiveWorkers(count int) {
	m.activeWorkers.Set(float64(count))
}

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

type JobMetricsMiddleware struct {
	next    job.JobHandler
	metrics *JobMetricsCollector
}

func NewJobMetricsMiddleware(handler job.JobHandler, metrics *JobMetricsCollector) *JobMetricsMiddleware {
	return &JobMetricsMiddleware{
		next:    handler,
		metrics: metrics,
	}
}

func (m *JobMetricsMiddleware) Execute(ctx context.Context, job job.Job) error {
	start := time.Now()

	err := m.next.Execute(ctx, job)

	duration := time.Since(start)
	m.metrics.ObserveJobDuration(job.GetType(), duration)

	success := err == nil
	m.metrics.IncrementJobsProcessed(job.GetType(), success)

	if !success {
		m.metrics.IncrementJobRetries(job.GetType())
	}

	return err
}

func (m *JobMetricsMiddleware) GetJobType() string {
	return m.next.GetJobType()
}

type JobQueueMetricsCollector struct {
	queue   job.JobQueue
	metrics *JobMetricsCollector
	queues  []string
}

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

func (c *JobQueueMetricsCollector) CollectQueueMetrics(ctx context.Context) error {
	for _, queueName := range c.queues {
		size, err := c.queue.GetQueueSize(ctx, queueName)
		if err != nil {
			continue
		}

		c.metrics.SetQueueSize(queueName, size)
	}

	return nil
}

type WorkerPoolMetricsCollector struct {
	metrics *JobMetricsCollector
}

func NewWorkerPoolMetricsCollector(metrics *JobMetricsCollector) *WorkerPoolMetricsCollector {
	return &WorkerPoolMetricsCollector{
		metrics: metrics,
	}
}

func (c *WorkerPoolMetricsCollector) UpdateWorkerCount(count int) {
	c.metrics.SetActiveWorkers(count)
}

type JobTypeMetrics struct {
	JobType        string  `json:"job_type"`
	TotalEnqueued  int64   `json:"total_enqueued"`
	TotalProcessed int64   `json:"total_processed"`
	TotalRetries   int64   `json:"total_retries"`
	SuccessRate    float64 `json:"success_rate"`
	AvgDuration    float64 `json:"avg_duration_seconds"`
}

type QueueMetrics struct {
	Queue       string `json:"queue"`
	CurrentSize int64  `json:"current_size"`
	Priority    string `json:"priority"`
}

type WorkerMetrics struct {
	ActiveWorkers int `json:"active_workers"`
}

type OverallMetrics struct {
	JobTypes []JobTypeMetrics `json:"job_types"`
	Queues   []QueueMetrics   `json:"queues"`
	Workers  WorkerMetrics    `json:"workers"`
}

type MetricsReporter struct {
	collector *JobMetricsCollector
}

func NewMetricsReporter(collector *JobMetricsCollector) *MetricsReporter {
	return &MetricsReporter{
		collector: collector,
	}
}

func (r *MetricsReporter) GetOverallMetrics() (*OverallMetrics, error) {
	return &OverallMetrics{
		JobTypes: []JobTypeMetrics{},
		Queues:   []QueueMetrics{},
		Workers:  WorkerMetrics{},
	}, nil
}

type BusinessMetricsCollector struct {
	emailJobsCounter    *prometheus.CounterVec
	fileProcessingGauge *prometheus.GaugeVec
	cleanupJobsDuration *prometheus.HistogramVec
}

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

func (b *BusinessMetricsCollector) RecordEmailJob(emailType string, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	b.emailJobsCounter.WithLabelValues(emailType, status).Inc()
}

func (b *BusinessMetricsCollector) SetActiveFileProcessingJobs(fileType string, count int) {
	b.fileProcessingGauge.WithLabelValues(fileType).Set(float64(count))
}

func (b *BusinessMetricsCollector) RecordCleanupDuration(cleanupType string, duration time.Duration) {
	b.cleanupJobsDuration.WithLabelValues(cleanupType).Observe(duration.Seconds())
}

type MetricsConfiguration struct {
	Enabled            bool          `yaml:"enabled" env:"METRICS_ENABLED" env-default:"true"`
	CollectionInterval time.Duration `yaml:"collection_interval" env:"METRICS_COLLECTION_INTERVAL" env-default:"30s"`
	RetentionPeriod    time.Duration `yaml:"retention_period" env:"METRICS_RETENTION_PERIOD" env-default:"24h"`
}

type MetricsServer struct {
	collector *JobMetricsCollector
	reporter  *MetricsReporter
}

func NewMetricsServer(collector *JobMetricsCollector) *MetricsServer {
	return &MetricsServer{
		collector: collector,
		reporter:  NewMetricsReporter(collector),
	}
}

func (s *MetricsServer) GetMetrics() (*OverallMetrics, error) {
	return s.reporter.GetOverallMetrics()
}
