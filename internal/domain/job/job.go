package job

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrJobNotFound         = errors.New("job not found")
	ErrInvalidJobType      = errors.New("invalid job type")
	ErrInvalidPayload      = errors.New("invalid job payload")
	ErrMaxRetriesReached   = errors.New("maximum retries reached")
	ErrJobAlreadyProcessed = errors.New("job already processed")
)

// BaseJob implements the Job interface with common functionality
type BaseJob struct {
	id          uuid.UUID
	jobType     string
	payload     JobPayload
	priority    JobPriority
	status      JobStatus
	maxRetries  int
	retryCount  int
	createdAt   time.Time
	scheduledAt *time.Time
	processedAt *time.Time
	error       error
}

// NewBaseJob creates a new base job
func NewBaseJob(jobType string, payload JobPayload) *BaseJob {
	return &BaseJob{
		id:         uuid.New(),
		jobType:    jobType,
		payload:    payload,
		priority:   PriorityNormal,
		status:     JobStatusPending,
		maxRetries: 3,
		retryCount: 0,
		createdAt:  time.Now(),
	}
}

// NewBaseJobWithOptions creates a new base job with options
func NewBaseJobWithOptions(jobType string, payload JobPayload, opts JobOptions) *BaseJob {
	job := &BaseJob{
		id:          uuid.New(),
		jobType:     jobType,
		payload:     payload,
		priority:    opts.Priority,
		status:      JobStatusPending,
		maxRetries:  opts.MaxRetries,
		retryCount:  0,
		createdAt:   time.Now(),
		scheduledAt: opts.ScheduledAt,
	}

	if job.maxRetries == 0 {
		job.maxRetries = 3
	}

	if job.priority == 0 {
		job.priority = PriorityNormal
	}

	return job
}

// GetID returns the unique identifier of the job
func (j *BaseJob) GetID() uuid.UUID {
	return j.id
}

// GetType returns the type/name of the job
func (j *BaseJob) GetType() string {
	return j.jobType
}

// GetPayload returns the job's data payload
func (j *BaseJob) GetPayload() JobPayload {
	return j.payload
}

// GetPriority returns the job's priority
func (j *BaseJob) GetPriority() JobPriority {
	return j.priority
}

// GetStatus returns the current status of the job
func (j *BaseJob) GetStatus() JobStatus {
	return j.status
}

// SetStatus updates the job's status
func (j *BaseJob) SetStatus(status JobStatus) {
	j.status = status
}

// GetMaxRetries returns the maximum number of retry attempts
func (j *BaseJob) GetMaxRetries() int {
	return j.maxRetries
}

// GetRetryCount returns the current retry count
func (j *BaseJob) GetRetryCount() int {
	return j.retryCount
}

// IncrementRetryCount increments the retry counter
func (j *BaseJob) IncrementRetryCount() {
	j.retryCount++
}

// GetCreatedAt returns when the job was created
func (j *BaseJob) GetCreatedAt() time.Time {
	return j.createdAt
}

// GetScheduledAt returns when the job should be executed
func (j *BaseJob) GetScheduledAt() *time.Time {
	return j.scheduledAt
}

// SetScheduledAt sets when the job should be executed
func (j *BaseJob) SetScheduledAt(at *time.Time) {
	j.scheduledAt = at
}

// GetProcessedAt returns when the job was processed
func (j *BaseJob) GetProcessedAt() *time.Time {
	return j.processedAt
}

// SetProcessedAt sets when the job was processed
func (j *BaseJob) SetProcessedAt(at *time.Time) {
	j.processedAt = at
}

// GetError returns the last error if job failed
func (j *BaseJob) GetError() error {
	return j.error
}

// SetError sets the error for failed jobs
func (j *BaseJob) SetError(err error) {
	j.error = err
}

// CanRetry checks if the job can be retried
func (j *BaseJob) CanRetry() bool {
	return j.retryCount < j.maxRetries
}

// IsExpired checks if the job has expired (for scheduled jobs)
func (j *BaseJob) IsExpired(ttl time.Duration) bool {
	if j.scheduledAt != nil {
		return time.Now().After(j.scheduledAt.Add(ttl))
	}
	return time.Now().After(j.createdAt.Add(ttl))
}

// Validate validates the job data
func (j *BaseJob) Validate() error {
	if j.jobType == "" {
		return ErrInvalidJobType
	}

	if j.payload == nil {
		return ErrInvalidPayload
	}

	return nil
}

// ToMap converts the job to a map for serialization
func (j *BaseJob) ToMap() map[string]interface{} {
	data := map[string]interface{}{
		"id":         j.id.String(),
		"type":       j.jobType,
		"payload":    j.payload,
		"priority":   j.priority,
		"status":     j.status,
		"maxRetries": j.maxRetries,
		"retryCount": j.retryCount,
		"createdAt":  j.createdAt.Unix(),
	}

	if j.scheduledAt != nil {
		data["scheduledAt"] = j.scheduledAt.Unix()
	}

	if j.processedAt != nil {
		data["processedAt"] = j.processedAt.Unix()
	}

	if j.error != nil {
		data["error"] = j.error.Error()
	}

	return data
}

// FromMap creates a job from map data
func FromMap(data map[string]interface{}) (*BaseJob, error) {
	id, err := uuid.Parse(data["id"].(string))
	if err != nil {
		return nil, err
	}

	job := &BaseJob{
		id:         id,
		jobType:    data["type"].(string),
		payload:    data["payload"].(JobPayload),
		priority:   JobPriority(data["priority"].(int)),
		status:     JobStatus(data["status"].(string)),
		maxRetries: data["maxRetries"].(int),
		retryCount: data["retryCount"].(int),
		createdAt:  time.Unix(data["createdAt"].(int64), 0),
	}

	if scheduledAt, exists := data["scheduledAt"]; exists && scheduledAt != nil {
		t := time.Unix(scheduledAt.(int64), 0)
		job.scheduledAt = &t
	}

	if processedAt, exists := data["processedAt"]; exists && processedAt != nil {
		t := time.Unix(processedAt.(int64), 0)
		job.processedAt = &t
	}

	if errorStr, exists := data["error"]; exists && errorStr != nil {
		job.error = errors.New(errorStr.(string))
	}

	return job, nil
}
