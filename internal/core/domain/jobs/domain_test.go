package jobs

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	jobsports "github.com/tranvuongduy2003/go-mvc/internal/core/ports/jobs"
)

func TestBaseJob(t *testing.T) {
	t.Run("NewBaseJob creates job with correct defaults", func(t *testing.T) {
		jobType := "test_job"
		payload := jobsports.JobPayload{
			"key1": "value1",
			"key2": 42,
		}

		baseJob := NewBaseJob(jobType, payload)

		assert.NotEqual(t, uuid.Nil, baseJob.GetID())
		assert.Equal(t, jobType, baseJob.GetType())
		assert.Equal(t, payload, baseJob.GetPayload())
		assert.Equal(t, jobsports.PriorityNormal, baseJob.GetPriority())
		assert.Equal(t, jobsports.JobStatusPending, baseJob.GetStatus())
		assert.WithinDuration(t, time.Now(), baseJob.GetCreatedAt(), time.Second)
		assert.Equal(t, 3, baseJob.GetMaxRetries())
		assert.Equal(t, 0, baseJob.GetRetryCount())
	})

	t.Run("CanRetry returns correct value", func(t *testing.T) {
		baseJob := NewBaseJob("test", jobsports.JobPayload{})

		// Should be able to retry initially
		assert.True(t, baseJob.CanRetry())

		// Increment retries up to max
		for i := 0; i < 3; i++ {
			baseJob.IncrementRetryCount()
		}

		// Should not be able to retry after max retries reached
		assert.False(t, baseJob.CanRetry())
	})
}

func TestEmailJob(t *testing.T) {
	t.Run("NewEmailJob creates valid email job", func(t *testing.T) {
		to := "test@example.com"
		subject := "Test Subject"
		body := "Test Body"

		emailJob := NewEmailJob(to, subject, body)

		assert.Equal(t, "email", emailJob.GetType())
		assert.Equal(t, jobsports.PriorityNormal, emailJob.GetPriority())

		payload := emailJob.GetPayload()
		assert.Equal(t, to, payload["to"])
		assert.Equal(t, subject, payload["subject"])
		assert.Equal(t, body, payload["body"])
	})
}

func TestFileProcessingJob(t *testing.T) {
	t.Run("NewFileProcessingJob creates valid job", func(t *testing.T) {
		filePath := "/path/to/file.jpg"
		operation := "resize"
		params := map[string]interface{}{
			"width":  800,
			"height": 600,
		}

		fileJob := NewFileProcessingJob(filePath, operation, params)

		assert.Equal(t, "file_processing", fileJob.GetType())
		assert.Equal(t, jobsports.PriorityNormal, fileJob.GetPriority())

		payload := fileJob.GetPayload()
		assert.Equal(t, filePath, payload["filePath"])
		assert.Equal(t, operation, payload["operation"])
		assert.Equal(t, params, payload["params"])
	})
}

func TestDataCleanupJob(t *testing.T) {
	t.Run("NewDataCleanupJob creates valid job", func(t *testing.T) {
		table := "users"
		condition := map[string]interface{}{
			"status": "inactive",
		}
		olderThan := "30d"

		cleanupJob := NewDataCleanupJob(table, condition, olderThan)

		assert.Equal(t, "data_cleanup", cleanupJob.GetType())
		assert.Equal(t, jobsports.PriorityLow, cleanupJob.GetPriority())

		payload := cleanupJob.GetPayload()
		assert.Equal(t, table, payload["table"])
		assert.Equal(t, condition, payload["condition"])
		assert.Equal(t, olderThan, payload["olderThan"])
	})
}

func TestJobFactory(t *testing.T) {
	t.Run("CreateJob creates email job correctly", func(t *testing.T) {
		factory := NewJobFactory()
		payload := jobsports.JobPayload{
			"to":      "test@example.com",
			"type":    "text",
			"subject": "Test",
			"body":    "Body",
		}

		job, err := factory.CreateJob("email", payload)

		assert.NoError(t, err)
		assert.NotNil(t, job)
		assert.Equal(t, "email", job.GetType())
	})

	t.Run("CreateJob creates file processing job correctly", func(t *testing.T) {
		factory := NewJobFactory()
		payload := jobsports.JobPayload{
			"filePath":  "/path/to/file.jpg",
			"operation": "resize",
			"params": map[string]interface{}{
				"width":  800,
				"height": 600,
			},
		}

		job, err := factory.CreateJob("file_processing", payload)

		assert.NoError(t, err)
		assert.NotNil(t, job)
		assert.Equal(t, "file_processing", job.GetType())
	})

	t.Run("CreateJob with invalid type creates generic job", func(t *testing.T) {
		factory := NewJobFactory()
		payload := jobsports.JobPayload{
			"data": "test",
		}

		job, err := factory.CreateJob("invalid_type", payload)

		assert.NoError(t, err)
		assert.NotNil(t, job)
		assert.Equal(t, "invalid_type", job.GetType())
	})
}
