package job

import (
	"fmt"
	"regexp"
)

type JobValidator struct {
	emailRegex *regexp.Regexp
}

func NewJobValidator() *JobValidator {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return &JobValidator{
		emailRegex: emailRegex,
	}
}

func (v *JobValidator) ValidateJob(job Job) error {
	if err := v.validateBasicJob(job); err != nil {
		return err
	}

	switch job.GetType() {
	case JobTypeEmail, JobTypeEmailTemplate:
		return v.validateEmailJob(job)
	case JobTypeFileProcessing, JobTypeImageResize:
		return v.validateFileProcessingJob(job)
	case JobTypeDataCleanup, JobTypeUserCleanup:
		return v.validateDataCleanupJob(job)
	case JobTypeNotification:
		return v.validateNotificationJob(job)
	default:
		return nil
	}
}

func (v *JobValidator) validateBasicJob(job Job) error {
	if job.GetType() == "" {
		return fmt.Errorf("job type is required")
	}

	if job.GetPayload() == nil {
		return fmt.Errorf("job payload is required")
	}

	if job.GetMaxRetries() < 0 {
		return fmt.Errorf("max retries cannot be negative")
	}

	if job.GetRetryCount() < 0 {
		return fmt.Errorf("retry count cannot be negative")
	}

	if job.GetRetryCount() > job.GetMaxRetries() {
		return fmt.Errorf("retry count cannot exceed max retries")
	}

	return nil
}

func (v *JobValidator) validateEmailJob(job Job) error {
	payload := job.GetPayload()

	to, exists := payload["to"]
	if !exists {
		return fmt.Errorf("email recipient is required")
	}

	toStr, ok := to.(string)
	if !ok || toStr == "" {
		return fmt.Errorf("email recipient must be a non-empty string")
	}

	if !v.emailRegex.MatchString(toStr) {
		return fmt.Errorf("invalid email address: %s", toStr)
	}

	emailType, exists := payload["type"]
	if !exists {
		return fmt.Errorf("email type is required")
	}

	switch emailType {
	case "text":
		if subject, exists := payload["subject"]; !exists || subject == "" {
			return fmt.Errorf("email subject is required for text emails")
		}
		if body, exists := payload["body"]; !exists || body == "" {
			return fmt.Errorf("email body is required for text emails")
		}
	case "template":
		if template, exists := payload["template"]; !exists || template == "" {
			return fmt.Errorf("email template is required for template emails")
		}
	default:
		return fmt.Errorf("invalid email type: %s", emailType)
	}

	return nil
}

func (v *JobValidator) validateFileProcessingJob(job Job) error {
	payload := job.GetPayload()

	filePath, exists := payload["filePath"]
	if !exists {
		return fmt.Errorf("file path is required")
	}

	if filePathStr, ok := filePath.(string); !ok || filePathStr == "" {
		return fmt.Errorf("file path must be a non-empty string")
	}

	operation, exists := payload["operation"]
	if !exists {
		return fmt.Errorf("file operation is required")
	}

	if operationStr, ok := operation.(string); !ok || operationStr == "" {
		return fmt.Errorf("file operation must be a non-empty string")
	}

	switch operation {
	case "resize":
		return v.validateImageResizeParams(payload)
	case "convert", "compress", "thumbnail":
		return nil
	default:
		return fmt.Errorf("unsupported file operation: %s", operation)
	}
}

func (v *JobValidator) validateImageResizeParams(payload JobPayload) error {
	params, exists := payload["params"]
	if !exists {
		return fmt.Errorf("resize parameters are required")
	}

	paramsMap, ok := params.(map[string]interface{})
	if !ok {
		return fmt.Errorf("resize parameters must be a map")
	}

	width, exists := paramsMap["width"]
	if !exists {
		return fmt.Errorf("width is required for image resize")
	}

	if widthInt, ok := width.(int); !ok || widthInt <= 0 {
		return fmt.Errorf("width must be a positive integer")
	}

	height, exists := paramsMap["height"]
	if !exists {
		return fmt.Errorf("height is required for image resize")
	}

	if heightInt, ok := height.(int); !ok || heightInt <= 0 {
		return fmt.Errorf("height must be a positive integer")
	}

	if quality, exists := paramsMap["quality"]; exists {
		if qualityInt, ok := quality.(int); !ok || qualityInt < 1 || qualityInt > 100 {
			return fmt.Errorf("quality must be an integer between 1 and 100")
		}
	}

	return nil
}

func (v *JobValidator) validateDataCleanupJob(job Job) error {
	payload := job.GetPayload()

	if cleanupType, exists := payload["type"]; exists && cleanupType == "user" {
		return v.validateUserCleanupJob(payload)
	}

	table, exists := payload["table"]
	if !exists {
		return fmt.Errorf("table name is required for data cleanup")
	}

	if tableStr, ok := table.(string); !ok || tableStr == "" {
		return fmt.Errorf("table name must be a non-empty string")
	}

	return nil
}

func (v *JobValidator) validateUserCleanupJob(payload JobPayload) error {
	userID, exists := payload["userID"]
	if !exists {
		return fmt.Errorf("user ID is required for user cleanup")
	}

	if userIDStr, ok := userID.(string); !ok || userIDStr == "" {
		return fmt.Errorf("user ID must be a non-empty string")
	}

	actions, exists := payload["actions"]
	if !exists {
		return fmt.Errorf("cleanup actions are required")
	}

	actionsSlice, ok := actions.([]string)
	if !ok {
		return fmt.Errorf("cleanup actions must be a slice of strings")
	}

	if len(actionsSlice) == 0 {
		return fmt.Errorf("at least one cleanup action is required")
	}

	validActions := map[string]bool{
		"delete_files":    true,
		"anonymize_data":  true,
		"remove_sessions": true,
		"cleanup_cache":   true,
		"revoke_tokens":   true,
	}

	for _, action := range actionsSlice {
		if !validActions[action] {
			return fmt.Errorf("invalid cleanup action: %s", action)
		}
	}

	return nil
}

func (v *JobValidator) validateNotificationJob(job Job) error {
	payload := job.GetPayload()

	userID, exists := payload["userID"]
	if !exists {
		return fmt.Errorf("user ID is required for notifications")
	}

	if userIDStr, ok := userID.(string); !ok || userIDStr == "" {
		return fmt.Errorf("user ID must be a non-empty string")
	}

	title, exists := payload["title"]
	if !exists {
		return fmt.Errorf("notification title is required")
	}

	if titleStr, ok := title.(string); !ok || titleStr == "" {
		return fmt.Errorf("notification title must be a non-empty string")
	}

	message, exists := payload["message"]
	if !exists {
		return fmt.Errorf("notification message is required")
	}

	if messageStr, ok := message.(string); !ok || messageStr == "" {
		return fmt.Errorf("notification message must be a non-empty string")
	}

	return nil
}

func (v *JobValidator) ValidateJobOptions(opts JobOptions) error {
	if opts.MaxRetries < 0 {
		return fmt.Errorf("max retries cannot be negative")
	}

	if opts.Priority < PriorityLow || opts.Priority > PriorityCritical {
		return fmt.Errorf("invalid job priority: %d", opts.Priority)
	}

	return nil
}
