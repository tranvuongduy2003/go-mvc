package jobs

import (
	"fmt"

	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/jobs"
)

const (
	// Job types
	JobTypeEmail          = "email"
	JobTypeEmailTemplate  = "email_template"
	JobTypeFileProcessing = "file_processing"
	JobTypeImageResize    = "image_resize"
	JobTypeDataCleanup    = "data_cleanup"
	JobTypeUserCleanup    = "user_cleanup"
	JobTypeBackup         = "backup"
	JobTypeExport         = "export"
	JobTypeNotification   = "notification"
	JobTypeAnalytics      = "analytics"
)

// EmailJob represents an email sending job
type EmailJob struct {
	*BaseJob
}

// NewEmailJob creates a new email job
func NewEmailJob(to, subject, body string) *EmailJob {
	payload := jobs.JobPayload{
		"to":      to,
		"subject": subject,
		"body":    body,
		"type":    "text",
	}

	baseJob := NewBaseJob(JobTypeEmail, payload)
	baseJob.priority = jobs.PriorityNormal

	return &EmailJob{
		BaseJob: baseJob,
	}
}

// NewEmailTemplateJob creates a new email template job
func NewEmailTemplateJob(to, template string, data map[string]interface{}) *EmailJob {
	payload := jobs.JobPayload{
		"to":       to,
		"template": template,
		"data":     data,
		"type":     "template",
	}

	baseJob := NewBaseJob(JobTypeEmailTemplate, payload)
	baseJob.priority = jobs.PriorityNormal

	return &EmailJob{
		BaseJob: baseJob,
	}
}

// Validate validates the email job
func (e *EmailJob) Validate() error {
	if err := e.BaseJob.Validate(); err != nil {
		return err
	}

	payload := e.GetPayload()

	// Check required fields
	if to, exists := payload["to"]; !exists || to == "" {
		return fmt.Errorf("email recipient is required")
	}

	emailType, exists := payload["type"]
	if !exists {
		return fmt.Errorf("email type is required")
	}

	switch emailType {
	case "text":
		if subject, exists := payload["subject"]; !exists || subject == "" {
			return fmt.Errorf("email subject is required")
		}
		if body, exists := payload["body"]; !exists || body == "" {
			return fmt.Errorf("email body is required")
		}
	case "template":
		if template, exists := payload["template"]; !exists || template == "" {
			return fmt.Errorf("email template is required")
		}
	default:
		return fmt.Errorf("invalid email type: %s", emailType)
	}

	return nil
}

// FileProcessingJob represents a file processing job
type FileProcessingJob struct {
	*BaseJob
}

// NewFileProcessingJob creates a new file processing job
func NewFileProcessingJob(filePath, operation string, params map[string]interface{}) *FileProcessingJob {
	payload := jobs.JobPayload{
		"filePath":  filePath,
		"operation": operation,
		"params":    params,
	}

	baseJob := NewBaseJob(JobTypeFileProcessing, payload)
	baseJob.priority = jobs.PriorityNormal

	return &FileProcessingJob{
		BaseJob: baseJob,
	}
}

// NewImageResizeJob creates a new image resize job
func NewImageResizeJob(imagePath string, width, height int, quality int) *FileProcessingJob {
	params := map[string]interface{}{
		"width":   width,
		"height":  height,
		"quality": quality,
	}

	payload := jobs.JobPayload{
		"filePath":  imagePath,
		"operation": "resize",
		"params":    params,
	}

	baseJob := NewBaseJob(JobTypeImageResize, payload)
	baseJob.priority = jobs.PriorityHigh

	return &FileProcessingJob{
		BaseJob: baseJob,
	}
}

// Validate validates the file processing job
func (f *FileProcessingJob) Validate() error {
	if err := f.BaseJob.Validate(); err != nil {
		return err
	}

	payload := f.GetPayload()

	if filePath, exists := payload["filePath"]; !exists || filePath == "" {
		return fmt.Errorf("file path is required")
	}

	if operation, exists := payload["operation"]; !exists || operation == "" {
		return fmt.Errorf("file operation is required")
	}

	return nil
}

// DataCleanupJob represents a data cleanup job
type DataCleanupJob struct {
	*BaseJob
}

// NewDataCleanupJob creates a new data cleanup job
func NewDataCleanupJob(table string, condition map[string]interface{}, olderThan string) *DataCleanupJob {
	payload := jobs.JobPayload{
		"table":     table,
		"condition": condition,
		"olderThan": olderThan,
	}

	baseJob := NewBaseJob(JobTypeDataCleanup, payload)
	baseJob.priority = jobs.PriorityLow

	return &DataCleanupJob{
		BaseJob: baseJob,
	}
}

// NewUserCleanupJob creates a new user cleanup job
func NewUserCleanupJob(userID string, actions []string) *DataCleanupJob {
	payload := jobs.JobPayload{
		"userID":  userID,
		"actions": actions,
		"type":    "user",
	}

	baseJob := NewBaseJob(JobTypeUserCleanup, payload)
	baseJob.priority = jobs.PriorityHigh

	return &DataCleanupJob{
		BaseJob: baseJob,
	}
}

// Validate validates the data cleanup job
func (d *DataCleanupJob) Validate() error {
	if err := d.BaseJob.Validate(); err != nil {
		return err
	}

	payload := d.GetPayload()

	if cleanupType, exists := payload["type"]; exists && cleanupType == "user" {
		if userID, exists := payload["userID"]; !exists || userID == "" {
			return fmt.Errorf("user ID is required for user cleanup")
		}
		if actions, exists := payload["actions"]; !exists || len(actions.([]string)) == 0 {
			return fmt.Errorf("cleanup actions are required")
		}
	} else {
		if table, exists := payload["table"]; !exists || table == "" {
			return fmt.Errorf("table name is required for data cleanup")
		}
	}

	return nil
}

// NotificationJob represents a push notification job
type NotificationJob struct {
	*BaseJob
}

// NewNotificationJob creates a new notification job
func NewNotificationJob(userID, title, message string, data map[string]interface{}) *NotificationJob {
	payload := jobs.JobPayload{
		"userID":  userID,
		"title":   title,
		"message": message,
		"data":    data,
	}

	baseJob := NewBaseJob(JobTypeNotification, payload)
	baseJob.priority = jobs.PriorityHigh

	return &NotificationJob{
		BaseJob: baseJob,
	}
}

// Validate validates the notification job
func (n *NotificationJob) Validate() error {
	if err := n.BaseJob.Validate(); err != nil {
		return err
	}

	payload := n.GetPayload()

	if userID, exists := payload["userID"]; !exists || userID == "" {
		return fmt.Errorf("user ID is required")
	}

	if title, exists := payload["title"]; !exists || title == "" {
		return fmt.Errorf("notification title is required")
	}

	if message, exists := payload["message"]; !exists || message == "" {
		return fmt.Errorf("notification message is required")
	}

	return nil
}

// JobFactory creates jobs of different types
type JobFactory struct{}

// NewJobFactory creates a new job factory
func NewJobFactory() *JobFactory {
	return &JobFactory{}
}

// CreateJob creates a job based on type and payload
func (f *JobFactory) CreateJob(jobType string, payload jobs.JobPayload) (jobs.Job, error) {
	switch jobType {
	case JobTypeEmail:
		job := &EmailJob{
			BaseJob: &BaseJob{
				jobType: jobType,
				payload: payload,
			},
		}
		return job, job.Validate()

	case JobTypeEmailTemplate:
		job := &EmailJob{
			BaseJob: &BaseJob{
				jobType: jobType,
				payload: payload,
			},
		}
		return job, job.Validate()

	case JobTypeFileProcessing, JobTypeImageResize:
		job := &FileProcessingJob{
			BaseJob: &BaseJob{
				jobType: jobType,
				payload: payload,
			},
		}
		return job, job.Validate()

	case JobTypeDataCleanup, JobTypeUserCleanup:
		job := &DataCleanupJob{
			BaseJob: &BaseJob{
				jobType: jobType,
				payload: payload,
			},
		}
		return job, job.Validate()

	case JobTypeNotification:
		job := &NotificationJob{
			BaseJob: &BaseJob{
				jobType: jobType,
				payload: payload,
			},
		}
		return job, job.Validate()

	default:
		// For unknown types, create a generic BaseJob
		job := NewBaseJob(jobType, payload)
		return job, job.Validate()
	}
}

// CreateJobWithOptions creates a job with options
func (f *JobFactory) CreateJobWithOptions(jobType string, payload jobs.JobPayload, opts jobs.JobOptions) (jobs.Job, error) {
	job, err := f.CreateJob(jobType, payload)
	if err != nil {
		return nil, err
	}

	// Apply options
	if baseJob, ok := job.(*BaseJob); ok {
		if opts.Priority != 0 {
			baseJob.priority = opts.Priority
		}
		if opts.MaxRetries > 0 {
			baseJob.maxRetries = opts.MaxRetries
		}
		if opts.ScheduledAt != nil {
			baseJob.scheduledAt = opts.ScheduledAt
		}
	}

	return job, nil
}
