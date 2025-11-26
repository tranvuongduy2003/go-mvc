package job

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
	JobStatusRetrying   JobStatus = "retrying"
)

type JobPriority int

const (
	PriorityLow      JobPriority = 0
	PriorityNormal   JobPriority = 1
	PriorityHigh     JobPriority = 2
	PriorityCritical JobPriority = 3
)

type JobPayload map[string]interface{}

type Job interface {
	GetID() uuid.UUID

	GetType() string

	GetPayload() JobPayload

	GetPriority() JobPriority

	GetStatus() JobStatus

	SetStatus(status JobStatus)

	GetMaxRetries() int

	GetRetryCount() int

	IncrementRetryCount()

	CanRetry() bool

	GetCreatedAt() time.Time

	GetScheduledAt() *time.Time

	SetScheduledAt(at *time.Time)

	GetProcessedAt() *time.Time

	SetProcessedAt(at *time.Time)

	GetError() error

	SetError(err error)
}

type JobHandler interface {
	Execute(ctx context.Context, job Job) error

	GetJobType() string
}

type JobQueue interface {
	Enqueue(ctx context.Context, job Job) error

	EnqueueDelayed(ctx context.Context, job Job, delay time.Duration) error

	EnqueueAt(ctx context.Context, job Job, at time.Time) error

	Dequeue(ctx context.Context, queues ...string) (Job, error)

	AckJob(ctx context.Context, job Job) error

	NackJob(ctx context.Context, job Job, err error) error

	GetQueueSize(ctx context.Context, queue string) (int64, error)

	GetPendingJobs(ctx context.Context, queue string, limit int) ([]Job, error)

	GetFailedJobs(ctx context.Context, queue string, limit int) ([]Job, error)

	DeleteJob(ctx context.Context, jobID uuid.UUID) error

	RetryJob(ctx context.Context, jobID uuid.UUID) error
}

type Worker interface {
	Start(ctx context.Context) error

	Stop(ctx context.Context) error

	RegisterHandler(handler JobHandler)

	GetWorkerID() string

	IsRunning() bool
}

type WorkerPool interface {
	Start(ctx context.Context) error

	Stop(ctx context.Context) error

	AddWorker(worker Worker)

	RemoveWorker(workerID string)

	GetWorkerCount() int

	GetStats() WorkerPoolStats
}

type WorkerPoolStats struct {
	ActiveWorkers      int   `json:"active_workers"`
	TotalJobsProcessed int64 `json:"total_jobs_processed"`
	SuccessfulJobs     int64 `json:"successful_jobs"`
	FailedJobs         int64 `json:"failed_jobs"`
}

type Scheduler interface {
	Schedule(ctx context.Context, job Job, at time.Time) error

	ScheduleRecurring(ctx context.Context, job Job, cronExpr string) error

	Cancel(ctx context.Context, jobID uuid.UUID) error

	Start(ctx context.Context) error

	Stop(ctx context.Context) error

	GetScheduledJobs(ctx context.Context) ([]Job, error)
}

type JobMetrics interface {
	IncrementJobsEnqueued(jobType string)

	IncrementJobsProcessed(jobType string, success bool)

	ObserveJobDuration(jobType string, duration time.Duration)

	SetQueueSize(queue string, size int64)

	IncrementJobRetries(jobType string)
}

type BackgroundJobService interface {
	SubmitJob(ctx context.Context, jobType string, payload JobPayload) (uuid.UUID, error)

	SubmitJobWithOptions(ctx context.Context, jobType string, payload JobPayload, opts JobOptions) (uuid.UUID, error)

	GetJob(ctx context.Context, jobID uuid.UUID) (Job, error)

	CancelJob(ctx context.Context, jobID uuid.UUID) error

	RetryJob(ctx context.Context, jobID uuid.UUID) error

	GetJobStatus(ctx context.Context, jobID uuid.UUID) (JobStatus, error)
}

type JobOptions struct {
	Priority    JobPriority    `json:"priority"`
	MaxRetries  int            `json:"max_retries"`
	Delay       *time.Duration `json:"delay"`
	ScheduledAt *time.Time     `json:"scheduled_at"`
	Queue       string         `json:"queue"`
}
