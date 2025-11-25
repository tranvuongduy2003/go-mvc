package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	domainjobs "github.com/tranvuongduy2003/go-mvc/internal/domain/jobs"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/ports/jobs"
)

const (
	// Redis key prefixes
	keyPrefixQueue      = "job:queue:"
	keyPrefixDelayed    = "job:delayed:"
	keyPrefixProcessing = "job:processing:"
	keyPrefixFailed     = "job:failed:"
	keyPrefixCompleted  = "job:completed:"
	keyPrefixJob        = "job:data:"
	keyPrefixStats      = "job:stats:"
	keyPrefixLock       = "job:lock:"

	// Default values
	defaultTimeout           = 30 * time.Second
	defaultVisibilityTimeout = 10 * time.Minute
	defaultLockTimeout       = 5 * time.Minute
)

// RedisJobQueue implements the JobQueue interface using Redis
type RedisJobQueue struct {
	client            redis.UniversalClient
	defaultTimeout    time.Duration
	visibilityTimeout time.Duration
	lockTimeout       time.Duration
}

// NewRedisJobQueue creates a new Redis job queue
func NewRedisJobQueue(client redis.UniversalClient) *RedisJobQueue {
	return &RedisJobQueue{
		client:            client,
		defaultTimeout:    defaultTimeout,
		visibilityTimeout: defaultVisibilityTimeout,
		lockTimeout:       defaultLockTimeout,
	}
}

// Enqueue adds a job to the queue
func (r *RedisJobQueue) Enqueue(ctx context.Context, job jobs.Job) error {
	return r.enqueueJob(ctx, job, false)
}

// EnqueueDelayed adds a job to be executed after a delay
func (r *RedisJobQueue) EnqueueDelayed(ctx context.Context, job jobs.Job, delay time.Duration) error {
	scheduledAt := time.Now().Add(delay)
	job.SetScheduledAt(&scheduledAt)
	return r.enqueueDelayedJob(ctx, job, scheduledAt)
}

// EnqueueAt adds a job to be executed at a specific time
func (r *RedisJobQueue) EnqueueAt(ctx context.Context, job jobs.Job, at time.Time) error {
	job.SetScheduledAt(&at)
	return r.enqueueDelayedJob(ctx, job, at)
}

// Dequeue retrieves the next job from the queue
func (r *RedisJobQueue) Dequeue(ctx context.Context, queues ...string) (jobs.Job, error) {
	if len(queues) == 0 {
		queues = []string{"default"}
	}

	// Process delayed jobs first
	for _, queue := range queues {
		r.processDelayedJobs(ctx, queue)
	}

	// Try to dequeue from each priority level
	priorities := []jobs.JobPriority{
		jobs.PriorityCritical,
		jobs.PriorityHigh,
		jobs.PriorityNormal,
		jobs.PriorityLow,
	}

	for _, priority := range priorities {
		for _, queue := range queues {
			job, err := r.dequeueFromQueue(ctx, queue, priority)
			if err != nil {
				continue
			}
			if job != nil {
				return job, nil
			}
		}
	}

	return nil, nil
}

// AckJob acknowledges that a job has been processed successfully
func (r *RedisJobQueue) AckJob(ctx context.Context, job jobs.Job) error {
	pipe := r.client.Pipeline()

	jobID := job.GetID().String()

	// Remove from processing queue
	pipe.LRem(ctx, r.getProcessingKey("default"), 1, jobID)

	// Update job status
	job.SetStatus(jobs.JobStatusCompleted)
	now := time.Now()
	job.SetProcessedAt(&now)

	// Store completed job data
	jobData, err := r.serializeJob(job)
	if err != nil {
		return fmt.Errorf("failed to serialize job: %w", err)
	}

	pipe.HSet(ctx, r.getJobKey(jobID), jobData)
	pipe.LPush(ctx, r.getCompletedKey("default"), jobID)
	pipe.Expire(ctx, r.getCompletedKey("default"), 24*time.Hour)

	// Update stats
	pipe.HIncrBy(ctx, r.getStatsKey("default"), "completed", 1)
	pipe.HIncrBy(ctx, r.getStatsKey("default"), fmt.Sprintf("completed:%s", job.GetType()), 1)

	_, err = pipe.Exec(ctx)
	return err
}

// NackJob marks a job as failed and potentially retry
func (r *RedisJobQueue) NackJob(ctx context.Context, job jobs.Job, err error) error {
	jobID := job.GetID().String()

	// Remove from processing queue
	r.client.LRem(ctx, r.getProcessingKey("default"), 1, jobID)

	// Update job with error
	job.SetError(err)
	job.IncrementRetryCount()

	if job.CanRetry() {
		// Requeue for retry with exponential backoff
		job.SetStatus(jobs.JobStatusRetrying)
		delay := r.calculateRetryDelay(job.GetRetryCount())
		return r.EnqueueDelayed(ctx, job, delay)
	} else {
		// Mark as permanently failed
		job.SetStatus(jobs.JobStatusFailed)
		now := time.Now()
		job.SetProcessedAt(&now)

		// Store failed job data
		jobData, serialErr := r.serializeJob(job)
		if serialErr != nil {
			return fmt.Errorf("failed to serialize failed job: %w", serialErr)
		}

		pipe := r.client.Pipeline()
		pipe.HSet(ctx, r.getJobKey(jobID), jobData)
		pipe.LPush(ctx, r.getFailedKey("default"), jobID)
		pipe.Expire(ctx, r.getFailedKey("default"), 7*24*time.Hour)

		// Update stats
		pipe.HIncrBy(ctx, r.getStatsKey("default"), "failed", 1)
		pipe.HIncrBy(ctx, r.getStatsKey("default"), fmt.Sprintf("failed:%s", job.GetType()), 1)

		_, execErr := pipe.Exec(ctx)
		return execErr
	}
}

// GetQueueSize returns the number of jobs in a specific queue
func (r *RedisJobQueue) GetQueueSize(ctx context.Context, queue string) (int64, error) {
	var total int64

	priorities := []jobs.JobPriority{
		jobs.PriorityCritical,
		jobs.PriorityHigh,
		jobs.PriorityNormal,
		jobs.PriorityLow,
	}

	for _, priority := range priorities {
		size, err := r.client.LLen(ctx, r.getQueueKey(queue, priority)).Result()
		if err != nil {
			return 0, fmt.Errorf("failed to get queue size: %w", err)
		}
		total += size
	}

	// Add delayed jobs
	delayedSize, err := r.client.ZCard(ctx, r.getDelayedKey(queue)).Result()
	if err != nil {
		return total, fmt.Errorf("failed to get delayed queue size: %w", err)
	}

	return total + delayedSize, nil
}

// GetPendingJobs returns jobs with pending status
func (r *RedisJobQueue) GetPendingJobs(ctx context.Context, queue string, limit int) ([]jobs.Job, error) {
	var allJobs []jobs.Job

	priorities := []jobs.JobPriority{
		jobs.PriorityCritical,
		jobs.PriorityHigh,
		jobs.PriorityNormal,
		jobs.PriorityLow,
	}

	remainingLimit := limit
	for _, priority := range priorities {
		if remainingLimit <= 0 {
			break
		}

		jobIDs, err := r.client.LRange(ctx, r.getQueueKey(queue, priority), 0, int64(remainingLimit-1)).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get pending job IDs: %w", err)
		}

		for _, jobID := range jobIDs {
			job, err := r.getJobByID(ctx, jobID)
			if err != nil {
				continue
			}
			if job != nil {
				allJobs = append(allJobs, job)
				remainingLimit--
			}
		}
	}

	return allJobs, nil
}

// GetFailedJobs returns jobs with failed status
func (r *RedisJobQueue) GetFailedJobs(ctx context.Context, queue string, limit int) ([]jobs.Job, error) {
	jobIDs, err := r.client.LRange(ctx, r.getFailedKey(queue), 0, int64(limit-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get failed job IDs: %w", err)
	}

	var failedJobs []jobs.Job
	for _, jobID := range jobIDs {
		job, err := r.getJobByID(ctx, jobID)
		if err != nil {
			continue
		}
		if job != nil {
			failedJobs = append(failedJobs, job)
		}
	}

	return failedJobs, nil
}

// DeleteJob removes a job from the queue permanently
func (r *RedisJobQueue) DeleteJob(ctx context.Context, jobID uuid.UUID) error {
	jobIDStr := jobID.String()

	pipe := r.client.Pipeline()

	// Remove from all possible queues
	priorities := []jobs.JobPriority{
		jobs.PriorityCritical,
		jobs.PriorityHigh,
		jobs.PriorityNormal,
		jobs.PriorityLow,
	}

	for _, priority := range priorities {
		pipe.LRem(ctx, r.getQueueKey("default", priority), 0, jobIDStr)
	}

	pipe.LRem(ctx, r.getProcessingKey("default"), 0, jobIDStr)
	pipe.LRem(ctx, r.getFailedKey("default"), 0, jobIDStr)
	pipe.LRem(ctx, r.getCompletedKey("default"), 0, jobIDStr)
	pipe.ZRem(ctx, r.getDelayedKey("default"), jobIDStr)
	pipe.Del(ctx, r.getJobKey(jobIDStr))

	_, err := pipe.Exec(ctx)
	return err
}

// RetryJob requeues a failed job for retry
func (r *RedisJobQueue) RetryJob(ctx context.Context, jobID uuid.UUID) error {
	jobIDStr := jobID.String()

	job, err := r.getJobByID(ctx, jobIDStr)
	if err != nil {
		return fmt.Errorf("failed to get job for retry: %w", err)
	}

	if job == nil {
		return fmt.Errorf("job not found: %s", jobIDStr)
	}

	if job.GetStatus() != jobs.JobStatusFailed {
		return fmt.Errorf("job is not in failed status: %s", job.GetStatus())
	}

	// Reset job for retry
	job.SetStatus(jobs.JobStatusPending)
	job.SetError(nil)
	job.SetProcessedAt(nil)

	// Remove from failed queue
	r.client.LRem(ctx, r.getFailedKey("default"), 1, jobIDStr)

	// Requeue the job
	return r.enqueueJob(ctx, job, true)
}

// Private helper methods

func (r *RedisJobQueue) enqueueJob(ctx context.Context, job jobs.Job, isRetry bool) error {
	jobID := job.GetID().String()

	// Serialize job data
	jobData, err := r.serializeJob(job)
	if err != nil {
		return fmt.Errorf("failed to serialize job: %w", err)
	}

	pipe := r.client.Pipeline()

	// Store job data
	pipe.HSet(ctx, r.getJobKey(jobID), jobData)
	pipe.Expire(ctx, r.getJobKey(jobID), 7*24*time.Hour)

	// Add to appropriate priority queue
	queueKey := r.getQueueKey("default", job.GetPriority())
	pipe.LPush(ctx, queueKey, jobID)

	// Update stats
	if !isRetry {
		pipe.HIncrBy(ctx, r.getStatsKey("default"), "enqueued", 1)
		pipe.HIncrBy(ctx, r.getStatsKey("default"), fmt.Sprintf("enqueued:%s", job.GetType()), 1)
	}

	_, err = pipe.Exec(ctx)
	return err
}

func (r *RedisJobQueue) enqueueDelayedJob(ctx context.Context, job jobs.Job, at time.Time) error {
	jobID := job.GetID().String()

	// Serialize job data
	jobData, err := r.serializeJob(job)
	if err != nil {
		return fmt.Errorf("failed to serialize delayed job: %w", err)
	}

	pipe := r.client.Pipeline()

	// Store job data
	pipe.HSet(ctx, r.getJobKey(jobID), jobData)
	pipe.Expire(ctx, r.getJobKey(jobID), 7*24*time.Hour)

	// Add to delayed jobs sorted set
	score := float64(at.Unix())
	pipe.ZAdd(ctx, r.getDelayedKey("default"), redis.Z{
		Score:  score,
		Member: jobID,
	})

	// Update stats
	pipe.HIncrBy(ctx, r.getStatsKey("default"), "scheduled", 1)
	pipe.HIncrBy(ctx, r.getStatsKey("default"), fmt.Sprintf("scheduled:%s", job.GetType()), 1)

	_, err = pipe.Exec(ctx)
	return err
}

func (r *RedisJobQueue) dequeueFromQueue(ctx context.Context, queue string, priority jobs.JobPriority) (jobs.Job, error) {
	queueKey := r.getQueueKey(queue, priority)
	processingKey := r.getProcessingKey(queue)

	// Use BRPOPLPUSH for atomic move from queue to processing
	jobID, err := r.client.BRPopLPush(ctx, queueKey, processingKey, 1*time.Second).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to dequeue job: %w", err)
	}

	// Get job data
	job, err := r.getJobByID(ctx, jobID)
	if err != nil {
		r.client.LRem(ctx, processingKey, 1, jobID)
		return nil, fmt.Errorf("failed to get job data: %w", err)
	}

	if job == nil {
		r.client.LRem(ctx, processingKey, 1, jobID)
		return nil, nil
	}

	// Update job status
	job.SetStatus(jobs.JobStatusProcessing)

	// Set visibility timeout for the job
	r.client.Expire(ctx, r.getJobKey(jobID), r.visibilityTimeout)

	return job, nil
}

func (r *RedisJobQueue) processDelayedJobs(ctx context.Context, queue string) error {
	now := float64(time.Now().Unix())

	// Get jobs that should be executed now
	result, err := r.client.ZRangeByScoreWithScores(ctx, r.getDelayedKey(queue), &redis.ZRangeBy{
		Min: "0",
		Max: fmt.Sprintf("%.0f", now),
	}).Result()

	if err != nil {
		return fmt.Errorf("failed to get delayed jobs: %w", err)
	}

	for _, z := range result {
		jobID := z.Member.(string)

		// Move job from delayed to main queue
		pipe := r.client.Pipeline()
		pipe.ZRem(ctx, r.getDelayedKey(queue), jobID)

		// Get job to determine priority
		job, err := r.getJobByID(ctx, jobID)
		if err != nil {
			continue
		}

		if job != nil {
			job.SetStatus(jobs.JobStatusPending)
			pipe.LPush(ctx, r.getQueueKey(queue, job.GetPriority()), jobID)
		}

		pipe.Exec(ctx)
	}

	return nil
}

func (r *RedisJobQueue) getJobByID(ctx context.Context, jobID string) (jobs.Job, error) {
	jobData, err := r.client.HGetAll(ctx, r.getJobKey(jobID)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get job data: %w", err)
	}

	if len(jobData) == 0 {
		return nil, nil
	}

	return r.deserializeJob(jobData)
}

func (r *RedisJobQueue) serializeJob(job jobs.Job) (map[string]interface{}, error) {
	// Convert job to BaseJob for serialization
	if baseJob, ok := job.(*domainjobs.BaseJob); ok {
		return baseJob.ToMap(), nil
	}

	// For other job types, create a map manually
	data := map[string]interface{}{
		"id":         job.GetID().String(),
		"type":       job.GetType(),
		"payload":    job.GetPayload(),
		"priority":   int(job.GetPriority()),
		"status":     string(job.GetStatus()),
		"maxRetries": job.GetMaxRetries(),
		"retryCount": job.GetRetryCount(),
		"createdAt":  job.GetCreatedAt().Unix(),
	}

	if scheduledAt := job.GetScheduledAt(); scheduledAt != nil {
		data["scheduledAt"] = scheduledAt.Unix()
	}

	if processedAt := job.GetProcessedAt(); processedAt != nil {
		data["processedAt"] = processedAt.Unix()
	}

	if err := job.GetError(); err != nil {
		data["error"] = err.Error()
	}

	return data, nil
}

func (r *RedisJobQueue) deserializeJob(data map[string]string) (jobs.Job, error) {
	// Convert string map to interface map
	jobData := make(map[string]interface{})
	for k, v := range data {
		switch k {
		case "priority", "maxRetries", "retryCount":
			if intVal, err := strconv.Atoi(v); err == nil {
				jobData[k] = intVal
			}
		case "createdAt", "scheduledAt", "processedAt":
			if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
				jobData[k] = intVal
			}
		case "payload":
			var payload jobs.JobPayload
			if err := json.Unmarshal([]byte(v), &payload); err == nil {
				jobData[k] = payload
			}
		default:
			jobData[k] = v
		}
	}

	return domainjobs.FromMap(jobData)
}

func (r *RedisJobQueue) calculateRetryDelay(retryCount int) time.Duration {
	// Exponential backoff: 2^retryCount seconds, max 5 minutes
	delay := time.Duration(1<<uint(retryCount)) * time.Second
	if delay > 5*time.Minute {
		delay = 5 * time.Minute
	}
	return delay
}

// Redis key generation methods

func (r *RedisJobQueue) getQueueKey(queue string, priority jobs.JobPriority) string {
	return fmt.Sprintf("%s%s:%d", keyPrefixQueue, queue, priority)
}

func (r *RedisJobQueue) getDelayedKey(queue string) string {
	return fmt.Sprintf("%s%s", keyPrefixDelayed, queue)
}

func (r *RedisJobQueue) getProcessingKey(queue string) string {
	return fmt.Sprintf("%s%s", keyPrefixProcessing, queue)
}

func (r *RedisJobQueue) getFailedKey(queue string) string {
	return fmt.Sprintf("%s%s", keyPrefixFailed, queue)
}

func (r *RedisJobQueue) getCompletedKey(queue string) string {
	return fmt.Sprintf("%s%s", keyPrefixCompleted, queue)
}

func (r *RedisJobQueue) getJobKey(jobID string) string {
	return fmt.Sprintf("%s%s", keyPrefixJob, jobID)
}

func (r *RedisJobQueue) getStatsKey(queue string) string {
	return fmt.Sprintf("%s%s", keyPrefixStats, queue)
}

func (r *RedisJobQueue) getLockKey(key string) string {
	return fmt.Sprintf("%s%s", keyPrefixLock, key)
}
