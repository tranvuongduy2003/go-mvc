package scheduler

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/job"
)

const (
	keyPrefixScheduled = "job:scheduled:"
	keyPrefixRecurring = "job:recurring:"
)

// SimpleScheduler implements the Scheduler interface using Redis
type SimpleScheduler struct {
	client   redis.UniversalClient
	queue    job.JobQueue
	running  bool
	shutdown chan struct{}
	wg       sync.WaitGroup
	mu       sync.RWMutex

	// Scheduled and recurring jobs tracking
	scheduledJobs map[uuid.UUID]*ScheduledJobInfo
	recurringJobs map[uuid.UUID]*RecurringJobInfo
}

// ScheduledJobInfo holds information about a scheduled job
type ScheduledJobInfo struct {
	Job         job.Job   `json:"job"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

// RecurringJobInfo holds information about a recurring job
type RecurringJobInfo struct {
	Job       job.Job       `json:"job"`
	Interval  time.Duration `json:"interval"`
	LastRun   time.Time     `json:"last_run"`
	NextRun   time.Time     `json:"next_run"`
	CreatedAt time.Time     `json:"created_at"`
}

// NewSimpleScheduler creates a new simple scheduler
func NewSimpleScheduler(client redis.UniversalClient, queue job.JobQueue) *SimpleScheduler {
	return &SimpleScheduler{
		client:        client,
		queue:         queue,
		shutdown:      make(chan struct{}),
		scheduledJobs: make(map[uuid.UUID]*ScheduledJobInfo),
		recurringJobs: make(map[uuid.UUID]*RecurringJobInfo),
	}
}

// Schedule adds a job to be executed at a specific time
func (ss *SimpleScheduler) Schedule(ctx context.Context, job job.Job, at time.Time) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	jobID := job.GetID()

	// Store scheduled job info
	info := &ScheduledJobInfo{
		Job:         job,
		ScheduledAt: at,
	}

	ss.scheduledJobs[jobID] = info

	// Store in Redis for persistence
	if err := ss.storeScheduledJob(ctx, jobID, info); err != nil {
		delete(ss.scheduledJobs, jobID)
		return fmt.Errorf("failed to store scheduled job: %w", err)
	}

	return nil
}

// ScheduleRecurring adds a recurring job with interval
func (ss *SimpleScheduler) ScheduleRecurring(ctx context.Context, job job.Job, cronExpr string) error {
	// Parse simple interval from cronExpr (simplified implementation)
	interval, err := ss.parseInterval(cronExpr)
	if err != nil {
		return fmt.Errorf("invalid cron expression: %w", err)
	}

	ss.mu.Lock()
	defer ss.mu.Unlock()

	jobID := job.GetID()
	now := time.Now()

	// Store recurring job info
	info := &RecurringJobInfo{
		Job:       job,
		Interval:  interval,
		LastRun:   now,
		NextRun:   now.Add(interval),
		CreatedAt: now,
	}

	ss.recurringJobs[jobID] = info

	// Store in Redis for persistence
	if err := ss.storeRecurringJob(ctx, jobID, info); err != nil {
		delete(ss.recurringJobs, jobID)
		return fmt.Errorf("failed to store recurring job: %w", err)
	}

	return nil
}

// Cancel cancels a scheduled job
func (ss *SimpleScheduler) Cancel(ctx context.Context, jobID uuid.UUID) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	// Check if it's a scheduled job
	if _, exists := ss.scheduledJobs[jobID]; exists {
		delete(ss.scheduledJobs, jobID)
		return ss.removeScheduledJob(ctx, jobID)
	}

	// Check if it's a recurring job
	if _, exists := ss.recurringJobs[jobID]; exists {
		delete(ss.recurringJobs, jobID)
		return ss.removeRecurringJob(ctx, jobID)
	}

	return fmt.Errorf("scheduled job not found: %s", jobID)
}

// Start begins the scheduler
func (ss *SimpleScheduler) Start(ctx context.Context) error {
	ss.mu.Lock()
	if ss.running {
		ss.mu.Unlock()
		return fmt.Errorf("scheduler is already running")
	}
	ss.running = true
	ss.mu.Unlock()

	// Load scheduled jobs from Redis
	if err := ss.loadJobs(ctx); err != nil {
		ss.mu.Lock()
		ss.running = false
		ss.mu.Unlock()
		return fmt.Errorf("failed to load jobs: %w", err)
	}

	// Start the main scheduler loop
	ss.wg.Add(1)
	go ss.run(ctx)

	return nil
}

// Stop stops the scheduler
func (ss *SimpleScheduler) Stop(ctx context.Context) error {
	ss.mu.Lock()
	if !ss.running {
		ss.mu.Unlock()
		return fmt.Errorf("scheduler is not running")
	}
	ss.running = false
	ss.mu.Unlock()

	// Signal shutdown
	close(ss.shutdown)

	// Wait for main loop to finish
	done := make(chan struct{})
	go func() {
		ss.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("scheduler shutdown timeout")
	}
}

// GetScheduledJobs returns all scheduled jobs
func (ss *SimpleScheduler) GetScheduledJobs(ctx context.Context) ([]job.Job, error) {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	allJobs := make([]job.Job, 0, len(ss.scheduledJobs)+len(ss.recurringJobs))

	for _, info := range ss.scheduledJobs {
		allJobs = append(allJobs, info.Job)
	}

	for _, info := range ss.recurringJobs {
		allJobs = append(allJobs, info.Job)
	}

	return allJobs, nil
}

// run is the main scheduler loop
func (ss *SimpleScheduler) run(ctx context.Context) {
	defer ss.wg.Done()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ss.shutdown:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			ss.processJobs(ctx)
		}
	}
}

// processJobs processes both scheduled and recurring jobs that are due
func (ss *SimpleScheduler) processJobs(ctx context.Context) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	now := time.Now()

	// Process one-time scheduled jobs
	var scheduledJobsToProcess []uuid.UUID
	for jobID, info := range ss.scheduledJobs {
		if now.After(info.ScheduledAt) {
			scheduledJobsToProcess = append(scheduledJobsToProcess, jobID)
		}
	}

	for _, jobID := range scheduledJobsToProcess {
		info := ss.scheduledJobs[jobID]

		// Enqueue the job
		if err := ss.queue.Enqueue(ctx, info.Job); err != nil {
			continue
		}

		// Remove from scheduled jobs
		delete(ss.scheduledJobs, jobID)

		// Remove from Redis
		go func(id uuid.UUID) {
			ss.removeScheduledJob(context.Background(), id)
		}(jobID)
	}

	// Process recurring jobs
	for jobID, info := range ss.recurringJobs {
		if now.After(info.NextRun) {
			// Clone and enqueue the job
			clonedJob := ss.cloneJob(info.Job)
			if err := ss.queue.Enqueue(ctx, clonedJob); err != nil {
				continue
			}

			// Update next run time
			info.LastRun = now
			info.NextRun = now.Add(info.Interval)

			// Update in Redis
			go func(id uuid.UUID, jobInfo *RecurringJobInfo) {
				ss.storeRecurringJob(context.Background(), id, jobInfo)
			}(jobID, info)
		}
	}
}

// cloneJob creates a copy of a job for recurring execution
func (ss *SimpleScheduler) cloneJob(clonedJob job.Job) job.Job {
	// Create a new job with the same type and payload but new ID
	newID := uuid.New()

	payload := make(job.JobPayload)
	for k, v := range clonedJob.GetPayload() {
		payload[k] = v
	}

	return &BasicScheduledJob{
		id:         newID,
		jobType:    clonedJob.GetType(),
		payload:    payload,
		priority:   clonedJob.GetPriority(),
		status:     job.JobStatusPending,
		maxRetries: 3,
		createdAt:  time.Now(),
	}
}

// parseInterval parses a simple interval from a string
// Supports formats like "5m", "1h", "30s", "1h30m"
func (ss *SimpleScheduler) parseInterval(cronExpr string) (time.Duration, error) {
	// For simplicity, treat cronExpr as Go duration
	// In practice, you'd implement proper cron parsing
	duration, err := time.ParseDuration(cronExpr)
	if err != nil {
		return 0, fmt.Errorf("invalid duration format: %w", err)
	}

	if duration <= 0 {
		return 0, fmt.Errorf("duration must be positive")
	}

	return duration, nil
}

// BasicScheduledJob is a simple implementation for scheduled jobs
type BasicScheduledJob struct {
	id          uuid.UUID
	jobType     string
	payload     job.JobPayload
	priority    job.JobPriority
	status      job.JobStatus
	maxRetries  int
	retryCount  int
	createdAt   time.Time
	scheduledAt *time.Time
	processedAt *time.Time
	error       error
}

func (b *BasicScheduledJob) GetID() uuid.UUID               { return b.id }
func (b *BasicScheduledJob) GetType() string                { return b.jobType }
func (b *BasicScheduledJob) GetPayload() job.JobPayload     { return b.payload }
func (b *BasicScheduledJob) GetPriority() job.JobPriority   { return b.priority }
func (b *BasicScheduledJob) GetStatus() job.JobStatus       { return b.status }
func (b *BasicScheduledJob) SetStatus(status job.JobStatus) { b.status = status }
func (b *BasicScheduledJob) GetMaxRetries() int             { return b.maxRetries }
func (b *BasicScheduledJob) GetRetryCount() int             { return b.retryCount }
func (b *BasicScheduledJob) IncrementRetryCount()           { b.retryCount++ }
func (b *BasicScheduledJob) CanRetry() bool                 { return b.retryCount < b.maxRetries }
func (b *BasicScheduledJob) GetCreatedAt() time.Time        { return b.createdAt }
func (b *BasicScheduledJob) GetScheduledAt() *time.Time     { return b.scheduledAt }
func (b *BasicScheduledJob) SetScheduledAt(at *time.Time)   { b.scheduledAt = at }
func (b *BasicScheduledJob) GetProcessedAt() *time.Time     { return b.processedAt }
func (b *BasicScheduledJob) SetProcessedAt(at *time.Time)   { b.processedAt = at }
func (b *BasicScheduledJob) GetError() error                { return b.error }
func (b *BasicScheduledJob) SetError(err error)             { b.error = err }

// Redis storage methods

func (ss *SimpleScheduler) storeScheduledJob(ctx context.Context, jobID uuid.UUID, info *ScheduledJobInfo) error {
	key := ss.getScheduledJobKey(jobID)

	data := map[string]interface{}{
		"job_id":       jobID.String(),
		"job_type":     info.Job.GetType(),
		"scheduled_at": info.ScheduledAt.Unix(),
	}

	return ss.client.HSet(ctx, key, data).Err()
}

func (ss *SimpleScheduler) storeRecurringJob(ctx context.Context, jobID uuid.UUID, info *RecurringJobInfo) error {
	key := ss.getRecurringJobKey(jobID)

	data := map[string]interface{}{
		"job_id":     jobID.String(),
		"job_type":   info.Job.GetType(),
		"interval":   info.Interval.String(),
		"last_run":   info.LastRun.Unix(),
		"next_run":   info.NextRun.Unix(),
		"created_at": info.CreatedAt.Unix(),
	}

	return ss.client.HSet(ctx, key, data).Err()
}

func (ss *SimpleScheduler) removeScheduledJob(ctx context.Context, jobID uuid.UUID) error {
	key := ss.getScheduledJobKey(jobID)
	return ss.client.Del(ctx, key).Err()
}

func (ss *SimpleScheduler) removeRecurringJob(ctx context.Context, jobID uuid.UUID) error {
	key := ss.getRecurringJobKey(jobID)
	return ss.client.Del(ctx, key).Err()
}

func (ss *SimpleScheduler) loadJobs(ctx context.Context) error {
	// Load scheduled jobs
	scheduledPattern := ss.getScheduledJobKey(uuid.UUID{}) + "*"
	scheduledKeys, err := ss.client.Keys(ctx, scheduledPattern).Result()
	if err != nil {
		return err
	}

	for _, key := range scheduledKeys {
		data, err := ss.client.HGetAll(ctx, key).Result()
		if err != nil {
			continue
		}

		info, err := ss.parseScheduledJobData(data)
		if err != nil {
			continue
		}

		ss.scheduledJobs[info.Job.GetID()] = info
	}

	// Load recurring jobs
	recurringPattern := ss.getRecurringJobKey(uuid.UUID{}) + "*"
	recurringKeys, err := ss.client.Keys(ctx, recurringPattern).Result()
	if err != nil {
		return err
	}

	for _, key := range recurringKeys {
		data, err := ss.client.HGetAll(ctx, key).Result()
		if err != nil {
			continue
		}

		info, err := ss.parseRecurringJobData(data)
		if err != nil {
			continue
		}

		ss.recurringJobs[info.Job.GetID()] = info
	}

	return nil
}

func (ss *SimpleScheduler) parseScheduledJobData(data map[string]string) (*ScheduledJobInfo, error) {
	jobIDStr, exists := data["job_id"]
	if !exists {
		return nil, fmt.Errorf("missing job_id")
	}

	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid job_id: %w", err)
	}

	jobType, exists := data["job_type"]
	if !exists {
		return nil, fmt.Errorf("missing job_type")
	}

	job := &BasicScheduledJob{
		id:         jobID,
		jobType:    jobType,
		payload:    make(job.JobPayload),
		priority:   job.PriorityNormal,
		status:     job.JobStatusPending,
		maxRetries: 3,
		createdAt:  time.Now(),
	}

	info := &ScheduledJobInfo{Job: job}

	if scheduledAtStr, exists := data["scheduled_at"]; exists {
		if timestamp, err := time.Parse("1136239445", scheduledAtStr); err == nil {
			info.ScheduledAt = timestamp
		}
	}

	return info, nil
}

func (ss *SimpleScheduler) parseRecurringJobData(data map[string]string) (*RecurringJobInfo, error) {
	jobIDStr, exists := data["job_id"]
	if !exists {
		return nil, fmt.Errorf("missing job_id")
	}

	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid job_id: %w", err)
	}

	jobType, exists := data["job_type"]
	if !exists {
		return nil, fmt.Errorf("missing job_type")
	}

	intervalStr, exists := data["interval"]
	if !exists {
		return nil, fmt.Errorf("missing interval")
	}

	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return nil, fmt.Errorf("invalid interval: %w", err)
	}

	job := &BasicScheduledJob{
		id:         jobID,
		jobType:    jobType,
		payload:    make(job.JobPayload),
		priority:   job.PriorityNormal,
		status:     job.JobStatusPending,
		maxRetries: 3,
		createdAt:  time.Now(),
	}

	info := &RecurringJobInfo{
		Job:      job,
		Interval: interval,
	}

	// Parse timestamps
	if lastRunStr, exists := data["last_run"]; exists {
		if timestamp, err := time.Parse("1136239445", lastRunStr); err == nil {
			info.LastRun = timestamp
		}
	}

	if nextRunStr, exists := data["next_run"]; exists {
		if timestamp, err := time.Parse("1136239445", nextRunStr); err == nil {
			info.NextRun = timestamp
		}
	}

	if createdAtStr, exists := data["created_at"]; exists {
		if timestamp, err := time.Parse("1136239445", createdAtStr); err == nil {
			info.CreatedAt = timestamp
		}
	}

	return info, nil
}

func (ss *SimpleScheduler) getScheduledJobKey(jobID uuid.UUID) string {
	if jobID == (uuid.UUID{}) {
		return keyPrefixScheduled
	}
	return fmt.Sprintf("%s%s", keyPrefixScheduled, jobID.String())
}

func (ss *SimpleScheduler) getRecurringJobKey(jobID uuid.UUID) string {
	if jobID == (uuid.UUID{}) {
		return keyPrefixRecurring
	}
	return fmt.Sprintf("%s%s", keyPrefixRecurring, jobID.String())
}

// Helper functions

// ScheduleBatch schedules multiple jobs at once
func (ss *SimpleScheduler) ScheduleBatch(ctx context.Context, jobs []job.Job, at time.Time) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	for _, job := range jobs {
		jobID := job.GetID()
		info := &ScheduledJobInfo{
			Job:         job,
			ScheduledAt: at,
		}

		if err := ss.storeScheduledJob(ctx, jobID, info); err != nil {
			return fmt.Errorf("failed to store batch job %s: %w", jobID, err)
		}

		ss.scheduledJobs[jobID] = info
	}

	return nil
}

// GetJobsByScheduleTime returns jobs scheduled within a time range
func (ss *SimpleScheduler) GetJobsByScheduleTime(start, end time.Time) []job.Job {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	var filteredJobs []job.Job

	for _, info := range ss.scheduledJobs {
		if info.ScheduledAt.After(start) && info.ScheduledAt.Before(end) {
			filteredJobs = append(filteredJobs, info.Job)
		}
	}

	// Sort by scheduled time
	sort.Slice(filteredJobs, func(i, j int) bool {
		return ss.scheduledJobs[filteredJobs[i].GetID()].ScheduledAt.Before(
			ss.scheduledJobs[filteredJobs[j].GetID()].ScheduledAt)
	})

	return filteredJobs
}
