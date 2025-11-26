package modules

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/job"
	jobsmetrics "github.com/tranvuongduy2003/go-mvc/internal/infrastructure/jobs/metrics"
	redisqueue "github.com/tranvuongduy2003/go-mvc/internal/infrastructure/jobs/redis"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/jobs/scheduler"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/jobs/worker"
)

// JobModuleParams defines the parameters required by the jobs module
type JobModuleParams struct {
	fx.In

	RedisClient redis.UniversalClient
}

// JobModuleResult defines what the jobs module provides
type JobModuleResult struct {
	fx.Out

	JobQueue   job.JobQueue
	WorkerPool *worker.WorkerPool
	Scheduler  job.Scheduler
	JobService job.BackgroundJobService
	JobMetrics job.JobMetrics
}

// JobModule provides background job processing functionality
var JobModule = fx.Module("jobs",
	// Provide core job components
	fx.Provide(
		NewJobQueue,
		NewScheduler,
		NewWorkerPool,
		NewJobMetricsCollector,
		NewBackgroundJobService,
	),

	// Register lifecycle hooks
	fx.Invoke(RegisterJobsLifecycle),
)

// NewJobQueue creates a new Redis-based job queue
func NewJobQueue(params JobModuleParams) job.JobQueue {
	return redisqueue.NewRedisJobQueue(params.RedisClient)
}

// NewScheduler creates a new scheduler
func NewScheduler(params JobModuleParams, queue job.JobQueue) job.Scheduler {
	return scheduler.NewSimpleScheduler(params.RedisClient, queue)
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(queue job.JobQueue, metrics job.JobMetrics) *worker.WorkerPool {
	// For now, we'll create the pool without metrics integration
	// The individual workers can use metrics when they're created
	return worker.NewWorkerPool(queue, 3)
}

// NewJobMetricsCollector creates a new job metrics collector
func NewJobMetricsCollector() job.JobMetrics {
	return jobsmetrics.NewJobMetricsCollector()
}

// BackgroundJobServiceImpl implements the BackgroundJobService interface
type BackgroundJobServiceImpl struct {
	queue     job.JobQueue
	scheduler job.Scheduler
	factory   *JobFactory
	metrics   job.JobMetrics
}

// JobFactory creates jobs from type and payload
type JobFactory struct{}

// NewJobFactory creates a new job factory
func NewJobFactory() *JobFactory {
	return &JobFactory{}
}

// CreateJob creates a job based on type and payload
func (f *JobFactory) CreateJob(jobType string, payload job.JobPayload, opts *job.JobOptions) (job.Job, error) {
	// Use the domain job factory
	domainFactory := job.NewJobFactory()
	if opts != nil {
		return domainFactory.CreateJobWithOptions(jobType, payload, *opts)
	}
	return domainFactory.CreateJob(jobType, payload)
}

// NewBackgroundJobService creates a new background job service
func NewBackgroundJobService(queue job.JobQueue, schedulerSvc job.Scheduler, metricsCollector job.JobMetrics) job.BackgroundJobService {
	return &BackgroundJobServiceImpl{
		queue:     queue,
		scheduler: schedulerSvc,
		factory:   NewJobFactory(),
		metrics:   metricsCollector,
	}
}

// SubmitJob submits a job for background processing
func (s *BackgroundJobServiceImpl) SubmitJob(ctx context.Context, jobType string, payload job.JobPayload) (uuid.UUID, error) {
	job, err := s.factory.CreateJob(jobType, payload, nil)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create job: %w", err)
	}

	if err := s.queue.Enqueue(ctx, job); err != nil {
		return uuid.Nil, fmt.Errorf("failed to enqueue job: %w", err)
	}

	// Record metrics
	s.metrics.IncrementJobsEnqueued(jobType)

	return job.GetID(), nil
}

// SubmitJobWithOptions submits a job with specific options
func (s *BackgroundJobServiceImpl) SubmitJobWithOptions(ctx context.Context, jobType string, payload job.JobPayload, opts job.JobOptions) (uuid.UUID, error) {
	job, err := s.factory.CreateJob(jobType, payload, &opts)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create job: %w", err)
	}

	// Handle scheduling options
	if opts.ScheduledAt != nil {
		if err := s.scheduler.Schedule(ctx, job, *opts.ScheduledAt); err != nil {
			return uuid.Nil, fmt.Errorf("failed to schedule job: %w", err)
		}
	} else if opts.Delay != nil {
		if err := s.queue.EnqueueDelayed(ctx, job, *opts.Delay); err != nil {
			return uuid.Nil, fmt.Errorf("failed to enqueue delayed job: %w", err)
		}
	} else {
		if err := s.queue.Enqueue(ctx, job); err != nil {
			return uuid.Nil, fmt.Errorf("failed to enqueue job: %w", err)
		}
	}

	// Record metrics
	s.metrics.IncrementJobsEnqueued(jobType)

	return job.GetID(), nil
}

// GetJob retrieves a job by ID (simplified implementation)
func (s *BackgroundJobServiceImpl) GetJob(ctx context.Context, jobID uuid.UUID) (job.Job, error) {
	// This would require additional Redis operations to look up jobs by ID
	// For now, return an error indicating this needs implementation
	return nil, fmt.Errorf("job lookup by ID not implemented")
}

// CancelJob cancels a pending job
func (s *BackgroundJobServiceImpl) CancelJob(ctx context.Context, jobID uuid.UUID) error {
	// Try to cancel from scheduler first
	if err := s.scheduler.Cancel(ctx, jobID); err == nil {
		return nil
	}

	// If not found in scheduler, try to delete from queue
	return s.queue.DeleteJob(ctx, jobID)
}

// RetryJob retries a failed job
func (s *BackgroundJobServiceImpl) RetryJob(ctx context.Context, jobID uuid.UUID) error {
	return s.queue.RetryJob(ctx, jobID)
}

// GetJobStatus returns the current status of a job (simplified)
func (s *BackgroundJobServiceImpl) GetJobStatus(ctx context.Context, jobID uuid.UUID) (job.JobStatus, error) {
	// This would require additional implementation to track job status
	return "", fmt.Errorf("job status lookup not implemented")
}

// JobsLifecycleParams defines parameters for lifecycle management
type JobsLifecycleParams struct {
	fx.In

	Lifecycle  fx.Lifecycle
	WorkerPool *worker.WorkerPool
	Scheduler  job.Scheduler
}

// RegisterJobsLifecycle registers lifecycle hooks for jobs components
func RegisterJobsLifecycle(params JobsLifecycleParams) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Start scheduler first
			if err := params.Scheduler.Start(ctx); err != nil {
				return fmt.Errorf("failed to start scheduler: %w", err)
			}

			// Then start worker pool
			if err := params.WorkerPool.Start(ctx); err != nil {
				params.Scheduler.Stop(ctx) // Cleanup scheduler if worker pool fails
				return fmt.Errorf("failed to start worker pool: %w", err)
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Stop worker pool first
			workerErr := params.WorkerPool.Stop(ctx)

			// Then stop scheduler
			schedulerErr := params.Scheduler.Stop(ctx)

			// Return any errors
			if workerErr != nil {
				return fmt.Errorf("failed to stop worker pool: %w", workerErr)
			}
			if schedulerErr != nil {
				return fmt.Errorf("failed to stop scheduler: %w", schedulerErr)
			}

			return nil
		},
	})
}

// JobHandlerRegistry manages job handlers
type JobHandlerRegistry struct {
	handlers map[string]job.JobHandler
}

// NewJobHandlerRegistry creates a new job handler registry
func NewJobHandlerRegistry() *JobHandlerRegistry {
	return &JobHandlerRegistry{
		handlers: make(map[string]job.JobHandler),
	}
}

// Register registers a job handler
func (r *JobHandlerRegistry) Register(handler job.JobHandler) {
	r.handlers[handler.GetJobType()] = handler
}

// Get retrieves a job handler by job type
func (r *JobHandlerRegistry) Get(jobType string) (job.JobHandler, bool) {
	handler, exists := r.handlers[jobType]
	return handler, exists
}

// RegisterHandlersWithWorkerPool registers all handlers with the worker pool
func RegisterHandlersWithWorkerPool(registry *JobHandlerRegistry, pool *worker.WorkerPool) {
	for _, handler := range registry.handlers {
		pool.RegisterHandler(handler)
	}
}

// JobConfiguration holds job system configuration
type JobConfiguration struct {
	WorkerCount       int           `yaml:"worker_count" env:"JOB_WORKER_COUNT" env-default:"3"`
	RedisKeyPrefix    string        `yaml:"redis_key_prefix" env:"JOB_REDIS_KEY_PREFIX" env-default:"job:"`
	ProcessingTimeout time.Duration `yaml:"processing_timeout" env:"JOB_PROCESSING_TIMEOUT" env-default:"10m"`
	RetryDelay        time.Duration `yaml:"retry_delay" env:"JOB_RETRY_DELAY" env-default:"30s"`
	MaxRetries        int           `yaml:"max_retries" env:"JOB_MAX_RETRIES" env-default:"3"`
}

// ConfigurableJobModule provides a configurable jobs module
func ConfigurableJobModule(config JobConfiguration) fx.Option {
	return fx.Module("configurable-jobs",
		fx.Supply(config),
		fx.Provide(
			NewConfigurableJobQueue,
			NewConfigurableScheduler,
			NewConfigurableWorkerPool,
			NewJobHandlerRegistry,
		),
		fx.Invoke(RegisterJobsLifecycle),
	)
}

// NewConfigurableJobQueue creates a configurable job queue
func NewConfigurableJobQueue(params JobModuleParams, config JobConfiguration) job.JobQueue {
	queue := redisqueue.NewRedisJobQueue(params.RedisClient)
	// Apply configuration here if the queue supports it
	return queue
}

// NewConfigurableScheduler creates a configurable scheduler
func NewConfigurableScheduler(params JobModuleParams, queue job.JobQueue, config JobConfiguration) job.Scheduler {
	return scheduler.NewSimpleScheduler(params.RedisClient, queue)
}

// NewConfigurableWorkerPool creates a configurable worker pool
func NewConfigurableWorkerPool(queue job.JobQueue, config JobConfiguration) *worker.WorkerPool {
	return worker.NewWorkerPool(queue, config.WorkerCount)
}
