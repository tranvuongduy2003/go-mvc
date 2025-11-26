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

func (j *BaseJob) GetID() uuid.UUID {
	return j.id
}

func (j *BaseJob) GetType() string {
	return j.jobType
}

func (j *BaseJob) GetPayload() JobPayload {
	return j.payload
}

func (j *BaseJob) GetPriority() JobPriority {
	return j.priority
}

func (j *BaseJob) GetStatus() JobStatus {
	return j.status
}

func (j *BaseJob) SetStatus(status JobStatus) {
	j.status = status
}

func (j *BaseJob) GetMaxRetries() int {
	return j.maxRetries
}

func (j *BaseJob) GetRetryCount() int {
	return j.retryCount
}

func (j *BaseJob) IncrementRetryCount() {
	j.retryCount++
}

func (j *BaseJob) GetCreatedAt() time.Time {
	return j.createdAt
}

func (j *BaseJob) GetScheduledAt() *time.Time {
	return j.scheduledAt
}

func (j *BaseJob) SetScheduledAt(at *time.Time) {
	j.scheduledAt = at
}

func (j *BaseJob) GetProcessedAt() *time.Time {
	return j.processedAt
}

func (j *BaseJob) SetProcessedAt(at *time.Time) {
	j.processedAt = at
}

func (j *BaseJob) GetError() error {
	return j.error
}

func (j *BaseJob) SetError(err error) {
	j.error = err
}

func (j *BaseJob) CanRetry() bool {
	return j.retryCount < j.maxRetries
}

func (j *BaseJob) IsExpired(ttl time.Duration) bool {
	if j.scheduledAt != nil {
		return time.Now().After(j.scheduledAt.Add(ttl))
	}
	return time.Now().After(j.createdAt.Add(ttl))
}

func (j *BaseJob) Validate() error {
	if j.jobType == "" {
		return ErrInvalidJobType
	}

	if j.payload == nil {
		return ErrInvalidPayload
	}

	return nil
}

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
