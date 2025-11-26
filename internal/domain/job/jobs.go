package job

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// JobStatus represents the current status of a job
type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
	JobStatusRetrying   JobStatus = "retrying"
)

// JobPriority defines the priority level of a job
type JobPriority int

const (
	PriorityLow      JobPriority = 0
	PriorityNormal   JobPriority = 1
	PriorityHigh     JobPriority = 2
	PriorityCritical JobPriority = 3
)

// JobPayload represents the data payload for a job
type JobPayload map[string]interface{}

// Job represents a background job that can be executed
type Job interface {
	// GetID returns the unique identifier of the job
	GetID() uuid.UUID

	// GetType returns the type/name of the job
	GetType() string

	// GetPayload returns the job's data payload
	GetPayload() JobPayload

	// GetPriority returns the job's priority
	GetPriority() JobPriority

	// GetStatus returns the current status of the job
	GetStatus() JobStatus

	// SetStatus updates the job's status
	SetStatus(status JobStatus)

	// GetMaxRetries returns the maximum number of retry attempts
	GetMaxRetries() int

	// GetRetryCount returns the current retry count
	GetRetryCount() int

	// IncrementRetryCount increments the retry counter
	IncrementRetryCount()

	// CanRetry checks if the job can be retried
	CanRetry() bool

	// GetCreatedAt returns when the job was created
	GetCreatedAt() time.Time

	// GetScheduledAt returns when the job should be executed
	GetScheduledAt() *time.Time

	// SetScheduledAt sets when the job should be executed
	SetScheduledAt(at *time.Time)

	// GetProcessedAt returns when the job was processed
	GetProcessedAt() *time.Time

	// SetProcessedAt sets when the job was processed
	SetProcessedAt(at *time.Time)

	// GetError returns the last error if job failed
	GetError() error

	// SetError sets the error for failed jobs
	SetError(err error)
}

// JobHandler defines how to execute a specific job type
type JobHandler interface {
	// Execute processes the job with the given context
	Execute(ctx context.Context, job Job) error

	// GetJobType returns the job type this handler can process
	GetJobType() string
}

// JobQueue defines the interface for job queue operations
type JobQueue interface {
	// Enqueue adds a job to the queue
	Enqueue(ctx context.Context, job Job) error

	// EnqueueDelayed adds a job to be executed at a specific time
	EnqueueDelayed(ctx context.Context, job Job, delay time.Duration) error

	// EnqueueAt adds a job to be executed at a specific time
	EnqueueAt(ctx context.Context, job Job, at time.Time) error

	// Dequeue retrieves the next job from the queue
	Dequeue(ctx context.Context, queues ...string) (Job, error)

	// AckJob acknowledges that a job has been processed successfully
	AckJob(ctx context.Context, job Job) error

	// NackJob marks a job as failed and potentially retry
	NackJob(ctx context.Context, job Job, err error) error

	// GetQueueSize returns the number of jobs in a specific queue
	GetQueueSize(ctx context.Context, queue string) (int64, error)

	// GetPendingJobs returns jobs with pending status
	GetPendingJobs(ctx context.Context, queue string, limit int) ([]Job, error)

	// GetFailedJobs returns jobs with failed status
	GetFailedJobs(ctx context.Context, queue string, limit int) ([]Job, error)

	// DeleteJob removes a job from the queue permanently
	DeleteJob(ctx context.Context, jobID uuid.UUID) error

	// RetryJob requeues a failed job for retry
	RetryJob(ctx context.Context, jobID uuid.UUID) error
}

// Worker defines the interface for job workers
type Worker interface {
	// Start begins processing jobs from the queue
	Start(ctx context.Context) error

	// Stop gracefully stops the worker
	Stop(ctx context.Context) error

	// RegisterHandler registers a job handler for a specific job type
	RegisterHandler(handler JobHandler)

	// GetWorkerID returns the unique identifier of the worker
	GetWorkerID() string

	// IsRunning returns whether the worker is currently running
	IsRunning() bool
}

// WorkerPool manages multiple workers
type WorkerPool interface {
	// Start starts all workers in the pool
	Start(ctx context.Context) error

	// Stop stops all workers gracefully
	Stop(ctx context.Context) error

	// AddWorker adds a worker to the pool
	AddWorker(worker Worker)

	// RemoveWorker removes a worker from the pool
	RemoveWorker(workerID string)

	// GetWorkerCount returns the number of active workers
	GetWorkerCount() int

	// GetStats returns worker pool statistics
	GetStats() WorkerPoolStats
}

// WorkerPoolStats provides statistics about the worker pool
type WorkerPoolStats struct {
	ActiveWorkers      int   `json:"active_workers"`
	TotalJobsProcessed int64 `json:"total_jobs_processed"`
	SuccessfulJobs     int64 `json:"successful_jobs"`
	FailedJobs         int64 `json:"failed_jobs"`
}

// Scheduler handles scheduled and recurring jobs
type Scheduler interface {
	// Schedule adds a job to be executed at a specific time
	Schedule(ctx context.Context, job Job, at time.Time) error

	// ScheduleRecurring adds a recurring job with cron expression
	ScheduleRecurring(ctx context.Context, job Job, cronExpr string) error

	// Cancel cancels a scheduled job
	Cancel(ctx context.Context, jobID uuid.UUID) error

	// Start begins the scheduler
	Start(ctx context.Context) error

	// Stop stops the scheduler
	Stop(ctx context.Context) error

	// GetScheduledJobs returns all scheduled jobs
	GetScheduledJobs(ctx context.Context) ([]Job, error)
}

// JobMetrics provides metrics about job processing
type JobMetrics interface {
	// IncrementJobsEnqueued increments the counter for enqueued jobs
	IncrementJobsEnqueued(jobType string)

	// IncrementJobsProcessed increments the counter for processed jobs
	IncrementJobsProcessed(jobType string, success bool)

	// ObserveJobDuration records the time taken to process a job
	ObserveJobDuration(jobType string, duration time.Duration)

	// SetQueueSize sets the current size of a queue
	SetQueueSize(queue string, size int64)

	// IncrementJobRetries increments retry counter
	IncrementJobRetries(jobType string)
}

// BackgroundJobService is the main service for managing background jobs
type BackgroundJobService interface {
	// SubmitJob submits a job for background processing
	SubmitJob(ctx context.Context, jobType string, payload JobPayload) (uuid.UUID, error)

	// SubmitJobWithOptions submits a job with specific options
	SubmitJobWithOptions(ctx context.Context, jobType string, payload JobPayload, opts JobOptions) (uuid.UUID, error)

	// GetJob retrieves a job by ID
	GetJob(ctx context.Context, jobID uuid.UUID) (Job, error)

	// CancelJob cancels a pending job
	CancelJob(ctx context.Context, jobID uuid.UUID) error

	// RetryJob retries a failed job
	RetryJob(ctx context.Context, jobID uuid.UUID) error

	// GetJobStatus returns the current status of a job
	GetJobStatus(ctx context.Context, jobID uuid.UUID) (JobStatus, error)
}

// JobOptions provides options when submitting jobs
type JobOptions struct {
	Priority    JobPriority    `json:"priority"`
	MaxRetries  int            `json:"max_retries"`
	Delay       *time.Duration `json:"delay"`
	ScheduledAt *time.Time     `json:"scheduled_at"`
	Queue       string         `json:"queue"`
}
