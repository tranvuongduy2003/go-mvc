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

type JobModuleParams struct {
	fx.In

	RedisClient redis.UniversalClient
}

type JobModuleResult struct {
	fx.Out

	JobQueue   job.JobQueue
	WorkerPool *worker.WorkerPool
	Scheduler  job.Scheduler
	JobService job.BackgroundJobService
	JobMetrics job.JobMetrics
}

var JobModule = fx.Module("jobs",
	fx.Provide(
		NewJobQueue,
		NewScheduler,
		NewWorkerPool,
		NewJobMetricsCollector,
		NewBackgroundJobService,
	),

	fx.Invoke(RegisterJobsLifecycle),
)

func NewJobQueue(params JobModuleParams) job.JobQueue {
	return redisqueue.NewRedisJobQueue(params.RedisClient)
}

func NewScheduler(params JobModuleParams, queue job.JobQueue) job.Scheduler {
	return scheduler.NewSimpleScheduler(params.RedisClient, queue)
}

func NewWorkerPool(queue job.JobQueue, metrics job.JobMetrics) *worker.WorkerPool {
	return worker.NewWorkerPool(queue, 3)
}

func NewJobMetricsCollector() job.JobMetrics {
	return jobsmetrics.NewJobMetricsCollector()
}

type BackgroundJobServiceImpl struct {
	queue     job.JobQueue
	scheduler job.Scheduler
	factory   *JobFactory
	metrics   job.JobMetrics
}

type JobFactory struct{}

func NewJobFactory() *JobFactory {
	return &JobFactory{}
}

func (f *JobFactory) CreateJob(jobType string, payload job.JobPayload, opts *job.JobOptions) (job.Job, error) {
	domainFactory := job.NewJobFactory()
	if opts != nil {
		return domainFactory.CreateJobWithOptions(jobType, payload, *opts)
	}
	return domainFactory.CreateJob(jobType, payload)
}

func NewBackgroundJobService(queue job.JobQueue, schedulerSvc job.Scheduler, metricsCollector job.JobMetrics) job.BackgroundJobService {
	return &BackgroundJobServiceImpl{
		queue:     queue,
		scheduler: schedulerSvc,
		factory:   NewJobFactory(),
		metrics:   metricsCollector,
	}
}

func (s *BackgroundJobServiceImpl) SubmitJob(ctx context.Context, jobType string, payload job.JobPayload) (uuid.UUID, error) {
	job, err := s.factory.CreateJob(jobType, payload, nil)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create job: %w", err)
	}

	if err := s.queue.Enqueue(ctx, job); err != nil {
		return uuid.Nil, fmt.Errorf("failed to enqueue job: %w", err)
	}

	s.metrics.IncrementJobsEnqueued(jobType)

	return job.GetID(), nil
}

func (s *BackgroundJobServiceImpl) SubmitJobWithOptions(ctx context.Context, jobType string, payload job.JobPayload, opts job.JobOptions) (uuid.UUID, error) {
	job, err := s.factory.CreateJob(jobType, payload, &opts)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create job: %w", err)
	}

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

	s.metrics.IncrementJobsEnqueued(jobType)

	return job.GetID(), nil
}

func (s *BackgroundJobServiceImpl) GetJob(ctx context.Context, jobID uuid.UUID) (job.Job, error) {
	return nil, fmt.Errorf("job lookup by ID not implemented")
}

func (s *BackgroundJobServiceImpl) CancelJob(ctx context.Context, jobID uuid.UUID) error {
	if err := s.scheduler.Cancel(ctx, jobID); err == nil {
		return nil
	}

	return s.queue.DeleteJob(ctx, jobID)
}

func (s *BackgroundJobServiceImpl) RetryJob(ctx context.Context, jobID uuid.UUID) error {
	return s.queue.RetryJob(ctx, jobID)
}

func (s *BackgroundJobServiceImpl) GetJobStatus(ctx context.Context, jobID uuid.UUID) (job.JobStatus, error) {
	return "", fmt.Errorf("job status lookup not implemented")
}

type JobsLifecycleParams struct {
	fx.In

	Lifecycle  fx.Lifecycle
	WorkerPool *worker.WorkerPool
	Scheduler  job.Scheduler
}

func RegisterJobsLifecycle(params JobsLifecycleParams) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := params.Scheduler.Start(ctx); err != nil {
				return fmt.Errorf("failed to start scheduler: %w", err)
			}

			if err := params.WorkerPool.Start(ctx); err != nil {
				params.Scheduler.Stop(ctx) // Cleanup scheduler if worker pool fails
				return fmt.Errorf("failed to start worker pool: %w", err)
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			workerErr := params.WorkerPool.Stop(ctx)

			schedulerErr := params.Scheduler.Stop(ctx)

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

type JobHandlerRegistry struct {
	handlers map[string]job.JobHandler
}

func NewJobHandlerRegistry() *JobHandlerRegistry {
	return &JobHandlerRegistry{
		handlers: make(map[string]job.JobHandler),
	}
}

func (r *JobHandlerRegistry) Register(handler job.JobHandler) {
	r.handlers[handler.GetJobType()] = handler
}

func (r *JobHandlerRegistry) Get(jobType string) (job.JobHandler, bool) {
	handler, exists := r.handlers[jobType]
	return handler, exists
}

func RegisterHandlersWithWorkerPool(registry *JobHandlerRegistry, pool *worker.WorkerPool) {
	for _, handler := range registry.handlers {
		pool.RegisterHandler(handler)
	}
}

type JobConfiguration struct {
	WorkerCount       int           `yaml:"worker_count" env:"JOB_WORKER_COUNT" env-default:"3"`
	RedisKeyPrefix    string        `yaml:"redis_key_prefix" env:"JOB_REDIS_KEY_PREFIX" env-default:"job:"`
	ProcessingTimeout time.Duration `yaml:"processing_timeout" env:"JOB_PROCESSING_TIMEOUT" env-default:"10m"`
	RetryDelay        time.Duration `yaml:"retry_delay" env:"JOB_RETRY_DELAY" env-default:"30s"`
	MaxRetries        int           `yaml:"max_retries" env:"JOB_MAX_RETRIES" env-default:"3"`
}

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

func NewConfigurableJobQueue(params JobModuleParams, config JobConfiguration) job.JobQueue {
	queue := redisqueue.NewRedisJobQueue(params.RedisClient)
	return queue
}

func NewConfigurableScheduler(params JobModuleParams, queue job.JobQueue, config JobConfiguration) job.Scheduler {
	return scheduler.NewSimpleScheduler(params.RedisClient, queue)
}

func NewConfigurableWorkerPool(queue job.JobQueue, config JobConfiguration) *worker.WorkerPool {
	return worker.NewWorkerPool(queue, config.WorkerCount)
}
