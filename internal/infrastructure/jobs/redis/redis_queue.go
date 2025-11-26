package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/job"
)

const (
	keyPrefixQueue      = "job:queue:"
	keyPrefixDelayed    = "job:delayed:"
	keyPrefixProcessing = "job:processing:"
	keyPrefixFailed     = "job:failed:"
	keyPrefixCompleted  = "job:completed:"
	keyPrefixJob        = "job:data:"
	keyPrefixStats      = "job:stats:"
	keyPrefixLock       = "job:lock:"

	defaultTimeout           = 30 * time.Second
	defaultVisibilityTimeout = 10 * time.Minute
	defaultLockTimeout       = 5 * time.Minute
)

type RedisJobQueue struct {
	client            redis.UniversalClient
	defaultTimeout    time.Duration
	visibilityTimeout time.Duration
	lockTimeout       time.Duration
}

func NewRedisJobQueue(client redis.UniversalClient) *RedisJobQueue {
	return &RedisJobQueue{
		client:            client,
		defaultTimeout:    defaultTimeout,
		visibilityTimeout: defaultVisibilityTimeout,
		lockTimeout:       defaultLockTimeout,
	}
}

func (r *RedisJobQueue) Enqueue(ctx context.Context, job job.Job) error {
	return r.enqueueJob(ctx, job, false)
}

func (r *RedisJobQueue) EnqueueDelayed(ctx context.Context, job job.Job, delay time.Duration) error {
	scheduledAt := time.Now().Add(delay)
	job.SetScheduledAt(&scheduledAt)
	return r.enqueueDelayedJob(ctx, job, scheduledAt)
}

func (r *RedisJobQueue) EnqueueAt(ctx context.Context, job job.Job, at time.Time) error {
	job.SetScheduledAt(&at)
	return r.enqueueDelayedJob(ctx, job, at)
}

func (r *RedisJobQueue) Dequeue(ctx context.Context, queues ...string) (job.Job, error) {
	if len(queues) == 0 {
		queues = []string{"default"}
	}

	for _, queue := range queues {
		r.processDelayedJobs(ctx, queue)
	}

	priorities := []job.JobPriority{
		job.PriorityCritical,
		job.PriorityHigh,
		job.PriorityNormal,
		job.PriorityLow,
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

func (r *RedisJobQueue) AckJob(ctx context.Context, ackedJob job.Job) error {
	pipe := r.client.Pipeline()

	jobID := ackedJob.GetID().String()

	pipe.LRem(ctx, r.getProcessingKey("default"), 1, jobID)

	ackedJob.SetStatus(job.JobStatusCompleted)
	now := time.Now()
	ackedJob.SetProcessedAt(&now)

	jobData, err := r.serializeJob(ackedJob)
	if err != nil {
		return fmt.Errorf("failed to serialize job: %w", err)
	}

	pipe.HSet(ctx, r.getJobKey(jobID), jobData)
	pipe.LPush(ctx, r.getCompletedKey("default"), jobID)
	pipe.Expire(ctx, r.getCompletedKey("default"), 24*time.Hour)

	pipe.HIncrBy(ctx, r.getStatsKey("default"), "completed", 1)
	pipe.HIncrBy(ctx, r.getStatsKey("default"), fmt.Sprintf("completed:%s", ackedJob.GetType()), 1)

	_, err = pipe.Exec(ctx)
	return err
}

func (r *RedisJobQueue) NackJob(ctx context.Context, nackedJob job.Job, err error) error {
	jobID := nackedJob.GetID().String()

	r.client.LRem(ctx, r.getProcessingKey("default"), 1, jobID)

	nackedJob.SetError(err)
	nackedJob.IncrementRetryCount()

	if nackedJob.CanRetry() {
		nackedJob.SetStatus(job.JobStatusRetrying)
		delay := r.calculateRetryDelay(nackedJob.GetRetryCount())
		return r.EnqueueDelayed(ctx, nackedJob, delay)
	} else {
		nackedJob.SetStatus(job.JobStatusFailed)
		now := time.Now()
		nackedJob.SetProcessedAt(&now)

		jobData, serialErr := r.serializeJob(nackedJob)
		if serialErr != nil {
			return fmt.Errorf("failed to serialize failed job: %w", serialErr)
		}

		pipe := r.client.Pipeline()
		pipe.HSet(ctx, r.getJobKey(jobID), jobData)
		pipe.LPush(ctx, r.getFailedKey("default"), jobID)
		pipe.Expire(ctx, r.getFailedKey("default"), 7*24*time.Hour)

		pipe.HIncrBy(ctx, r.getStatsKey("default"), "failed", 1)
		pipe.HIncrBy(ctx, r.getStatsKey("default"), fmt.Sprintf("failed:%s", nackedJob.GetType()), 1)

		_, execErr := pipe.Exec(ctx)
		return execErr
	}
}

func (r *RedisJobQueue) GetQueueSize(ctx context.Context, queue string) (int64, error) {
	var total int64

	priorities := []job.JobPriority{
		job.PriorityCritical,
		job.PriorityHigh,
		job.PriorityNormal,
		job.PriorityLow,
	}

	for _, priority := range priorities {
		size, err := r.client.LLen(ctx, r.getQueueKey(queue, priority)).Result()
		if err != nil {
			return 0, fmt.Errorf("failed to get queue size: %w", err)
		}
		total += size
	}

	delayedSize, err := r.client.ZCard(ctx, r.getDelayedKey(queue)).Result()
	if err != nil {
		return total, fmt.Errorf("failed to get delayed queue size: %w", err)
	}

	return total + delayedSize, nil
}

func (r *RedisJobQueue) GetPendingJobs(ctx context.Context, queue string, limit int) ([]job.Job, error) {
	var allJobs []job.Job

	priorities := []job.JobPriority{
		job.PriorityCritical,
		job.PriorityHigh,
		job.PriorityNormal,
		job.PriorityLow,
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

func (r *RedisJobQueue) GetFailedJobs(ctx context.Context, queue string, limit int) ([]job.Job, error) {
	jobIDs, err := r.client.LRange(ctx, r.getFailedKey(queue), 0, int64(limit-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get failed job IDs: %w", err)
	}

	var failedJobs []job.Job
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

func (r *RedisJobQueue) DeleteJob(ctx context.Context, jobID uuid.UUID) error {
	jobIDStr := jobID.String()

	pipe := r.client.Pipeline()

	priorities := []job.JobPriority{
		job.PriorityCritical,
		job.PriorityHigh,
		job.PriorityNormal,
		job.PriorityLow,
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

func (r *RedisJobQueue) RetryJob(ctx context.Context, jobID uuid.UUID) error {
	jobIDStr := jobID.String()

	retrievedJob, err := r.getJobByID(ctx, jobIDStr)
	if err != nil {
		return fmt.Errorf("failed to get job for retry: %w", err)
	}

	if retrievedJob == nil {
		return fmt.Errorf("job not found: %s", jobIDStr)
	}

	if retrievedJob.GetStatus() != job.JobStatusFailed {
		return fmt.Errorf("job is not in failed status: %s", retrievedJob.GetStatus())
	}

	retrievedJob.SetStatus(job.JobStatusPending)
	retrievedJob.SetError(nil)
	retrievedJob.SetProcessedAt(nil)

	r.client.LRem(ctx, r.getFailedKey("default"), 1, jobIDStr)

	return r.enqueueJob(ctx, retrievedJob, true)
}

func (r *RedisJobQueue) enqueueJob(ctx context.Context, job job.Job, isRetry bool) error {
	jobID := job.GetID().String()

	jobData, err := r.serializeJob(job)
	if err != nil {
		return fmt.Errorf("failed to serialize job: %w", err)
	}

	pipe := r.client.Pipeline()

	pipe.HSet(ctx, r.getJobKey(jobID), jobData)
	pipe.Expire(ctx, r.getJobKey(jobID), 7*24*time.Hour)

	queueKey := r.getQueueKey("default", job.GetPriority())
	pipe.LPush(ctx, queueKey, jobID)

	if !isRetry {
		pipe.HIncrBy(ctx, r.getStatsKey("default"), "enqueued", 1)
		pipe.HIncrBy(ctx, r.getStatsKey("default"), fmt.Sprintf("enqueued:%s", job.GetType()), 1)
	}

	_, err = pipe.Exec(ctx)
	return err
}

func (r *RedisJobQueue) enqueueDelayedJob(ctx context.Context, job job.Job, at time.Time) error {
	jobID := job.GetID().String()

	jobData, err := r.serializeJob(job)
	if err != nil {
		return fmt.Errorf("failed to serialize delayed job: %w", err)
	}

	pipe := r.client.Pipeline()

	pipe.HSet(ctx, r.getJobKey(jobID), jobData)
	pipe.Expire(ctx, r.getJobKey(jobID), 7*24*time.Hour)

	score := float64(at.Unix())
	pipe.ZAdd(ctx, r.getDelayedKey("default"), redis.Z{
		Score:  score,
		Member: jobID,
	})

	pipe.HIncrBy(ctx, r.getStatsKey("default"), "scheduled", 1)
	pipe.HIncrBy(ctx, r.getStatsKey("default"), fmt.Sprintf("scheduled:%s", job.GetType()), 1)

	_, err = pipe.Exec(ctx)
	return err
}

func (r *RedisJobQueue) dequeueFromQueue(ctx context.Context, queue string, priority job.JobPriority) (job.Job, error) {
	queueKey := r.getQueueKey(queue, priority)
	processingKey := r.getProcessingKey(queue)

	jobID, err := r.client.BRPopLPush(ctx, queueKey, processingKey, 1*time.Second).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to dequeue job: %w", err)
	}

	retrievedJob, err := r.getJobByID(ctx, jobID)
	if err != nil {
		r.client.LRem(ctx, processingKey, 1, jobID)
		return nil, fmt.Errorf("failed to get job data: %w", err)
	}

	if retrievedJob == nil {
		r.client.LRem(ctx, processingKey, 1, jobID)
		return nil, nil
	}

	retrievedJob.SetStatus(job.JobStatusProcessing)

	r.client.Expire(ctx, r.getJobKey(jobID), r.visibilityTimeout)

	return retrievedJob, nil
}

func (r *RedisJobQueue) processDelayedJobs(ctx context.Context, queue string) error {
	now := float64(time.Now().Unix())

	result, err := r.client.ZRangeByScoreWithScores(ctx, r.getDelayedKey(queue), &redis.ZRangeBy{
		Min: "0",
		Max: fmt.Sprintf("%.0f", now),
	}).Result()

	if err != nil {
		return fmt.Errorf("failed to get delayed jobs: %w", err)
	}

	for _, z := range result {
		jobID := z.Member.(string)

		pipe := r.client.Pipeline()
		pipe.ZRem(ctx, r.getDelayedKey(queue), jobID)

		retrievedJob, err := r.getJobByID(ctx, jobID)
		if err != nil {
			continue
		}

		if retrievedJob != nil {
			retrievedJob.SetStatus(job.JobStatusPending)
			pipe.LPush(ctx, r.getQueueKey(queue, retrievedJob.GetPriority()), jobID)
		}

		pipe.Exec(ctx)
	}

	return nil
}

func (r *RedisJobQueue) getJobByID(ctx context.Context, jobID string) (job.Job, error) {
	jobData, err := r.client.HGetAll(ctx, r.getJobKey(jobID)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get job data: %w", err)
	}

	if len(jobData) == 0 {
		return nil, nil
	}

	return r.deserializeJob(jobData)
}

func (r *RedisJobQueue) serializeJob(serializedJob job.Job) (map[string]interface{}, error) {
	if baseJob, ok := serializedJob.(*job.BaseJob); ok {
		return baseJob.ToMap(), nil
	}

	data := map[string]interface{}{
		"id":         serializedJob.GetID().String(),
		"type":       serializedJob.GetType(),
		"payload":    serializedJob.GetPayload(),
		"priority":   int(serializedJob.GetPriority()),
		"status":     string(serializedJob.GetStatus()),
		"maxRetries": serializedJob.GetMaxRetries(),
		"retryCount": serializedJob.GetRetryCount(),
		"createdAt":  serializedJob.GetCreatedAt().Unix(),
	}

	if scheduledAt := serializedJob.GetScheduledAt(); scheduledAt != nil {
		data["scheduledAt"] = scheduledAt.Unix()
	}

	if processedAt := serializedJob.GetProcessedAt(); processedAt != nil {
		data["processedAt"] = processedAt.Unix()
	}

	if err := serializedJob.GetError(); err != nil {
		data["error"] = err.Error()
	}

	return data, nil
}

func (r *RedisJobQueue) deserializeJob(data map[string]string) (job.Job, error) {
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
			var payload job.JobPayload
			if err := json.Unmarshal([]byte(v), &payload); err == nil {
				jobData[k] = payload
			}
		default:
			jobData[k] = v
		}
	}

	return job.FromMap(jobData)
}

func (r *RedisJobQueue) calculateRetryDelay(retryCount int) time.Duration {
	delay := time.Duration(1<<uint(retryCount)) * time.Second
	if delay > 5*time.Minute {
		delay = 5 * time.Minute
	}
	return delay
}

func (r *RedisJobQueue) getQueueKey(queue string, priority job.JobPriority) string {
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
