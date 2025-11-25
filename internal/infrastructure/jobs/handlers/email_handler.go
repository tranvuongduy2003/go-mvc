package handlers

import (
	"context"
	"fmt"
	"time"

	domainjobs "github.com/tranvuongduy2003/go-mvc/internal/domain/jobs"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/ports/jobs"
)

// EmailJobHandler handles email sending jobs
type EmailJobHandler struct {
	emailService EmailService
	metrics      jobs.JobMetrics
}

// EmailService defines the interface for sending emails
type EmailService interface {
	SendEmail(ctx context.Context, to, subject, body string) error
	SendBulkEmail(ctx context.Context, recipients []string, subject, body string) error
	SendTemplateEmail(ctx context.Context, to, templateID string, data map[string]interface{}) error
}

// NewEmailJobHandler creates a new email job handler
func NewEmailJobHandler(emailService EmailService, metrics jobs.JobMetrics) *EmailJobHandler {
	return &EmailJobHandler{
		emailService: emailService,
		metrics:      metrics,
	}
}

// Execute processes an email job
func (h *EmailJobHandler) Execute(ctx context.Context, job jobs.Job) error {
	// Start timing
	start := time.Now()
	defer func() {
		if h.metrics != nil {
			h.metrics.ObserveJobDuration(job.GetType(), time.Since(start))
		}
	}()

	// Cast to EmailJob
	emailJob, ok := job.(*domainjobs.EmailJob)
	if !ok {
		err := fmt.Errorf("expected EmailJob, got %T", job)
		if h.metrics != nil {
			h.metrics.IncrementJobsProcessed(job.GetType(), false)
		}
		return err
	}

	// Extract email data
	payload := emailJob.GetPayload()
	emailType, _ := payload["email_type"].(string)
	to, _ := payload["to"].(string)
	subject, _ := payload["subject"].(string)
	body, _ := payload["body"].(string)

	var err error

	// Handle different email types
	switch emailType {
	case "welcome":
		err = h.handleWelcomeEmail(ctx, to, payload)
	case "password_reset":
		err = h.handlePasswordResetEmail(ctx, to, payload)
	case "notification":
		err = h.handleNotificationEmail(ctx, to, subject, body)
	case "bulk":
		err = h.handleBulkEmail(ctx, payload)
	case "template":
		err = h.handleTemplateEmail(ctx, to, payload)
	default:
		// Generic email
		err = h.emailService.SendEmail(ctx, to, subject, body)
	}

	// Record metrics
	if h.metrics != nil {
		success := err == nil
		h.metrics.IncrementJobsProcessed(job.GetType(), success)

		// Custom business metrics
		if businessMetrics, ok := h.metrics.(*BusinessEmailMetrics); ok {
			businessMetrics.RecordEmailJob(emailType, success)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// GetJobType returns the job type this handler processes
func (h *EmailJobHandler) GetJobType() string {
	return "email"
}

// Private helper methods for different email types

func (h *EmailJobHandler) handleWelcomeEmail(ctx context.Context, to string, payload jobs.JobPayload) error {
	username, _ := payload["username"].(string)

	subject := "Welcome to Go-MVC!"
	body := fmt.Sprintf(`
		<h1>Welcome %s!</h1>
		<p>Thank you for signing up for Go-MVC. We're excited to have you on board!</p>
		<p>Get started by exploring our features and documentation.</p>
		<p>Best regards,<br>The Go-MVC Team</p>
	`, username)

	return h.emailService.SendEmail(ctx, to, subject, body)
}

func (h *EmailJobHandler) handlePasswordResetEmail(ctx context.Context, to string, payload jobs.JobPayload) error {
	token, _ := payload["reset_token"].(string)
	resetURL, _ := payload["reset_url"].(string)

	if resetURL == "" {
		resetURL = fmt.Sprintf("https://example.com/reset-password?token=%s", token)
	}

	subject := "Password Reset Request"
	body := fmt.Sprintf(`
		<h1>Password Reset</h1>
		<p>You requested a password reset for your account.</p>
		<p>Click the link below to reset your password:</p>
		<p><a href="%s">Reset Password</a></p>
		<p>This link will expire in 1 hour.</p>
		<p>If you didn't request this, please ignore this email.</p>
	`, resetURL)

	return h.emailService.SendEmail(ctx, to, subject, body)
}

func (h *EmailJobHandler) handleNotificationEmail(ctx context.Context, to, subject, body string) error {
	// Add notification wrapper
	wrappedBody := fmt.Sprintf(`
		<div style="padding: 20px; border-left: 4px solid #007cba;">
			<h2>📢 Notification</h2>
			%s
			<hr>
			<p style="font-size: 12px; color: #666;">
				This is an automated notification from Go-MVC.
			</p>
		</div>
	`, body)

	return h.emailService.SendEmail(ctx, to, subject, wrappedBody)
}

func (h *EmailJobHandler) handleBulkEmail(ctx context.Context, payload jobs.JobPayload) error {
	recipients, _ := payload["recipients"].([]string)
	subject, _ := payload["subject"].(string)
	body, _ := payload["body"].(string)

	if len(recipients) == 0 {
		return fmt.Errorf("no recipients specified for bulk email")
	}

	return h.emailService.SendBulkEmail(ctx, recipients, subject, body)
}

func (h *EmailJobHandler) handleTemplateEmail(ctx context.Context, to string, payload jobs.JobPayload) error {
	templateID, _ := payload["template_id"].(string)
	templateData, _ := payload["template_data"].(map[string]interface{})

	if templateID == "" {
		return fmt.Errorf("template_id is required for template emails")
	}

	return h.emailService.SendTemplateEmail(ctx, to, templateID, templateData)
}

// EmailJobFactory creates email jobs with proper validation
type EmailJobFactory struct{}

// NewEmailJobFactory creates a new email job factory
func NewEmailJobFactory() *EmailJobFactory {
	return &EmailJobFactory{}
}

// CreateWelcomeEmailJob creates a welcome email job
func (f *EmailJobFactory) CreateWelcomeEmailJob(to, username string) (jobs.Job, error) {
	payload := jobs.JobPayload{
		"email_type": "welcome",
		"to":         to,
		"username":   username,
	}

	factory := domainjobs.NewJobFactory()
	return factory.CreateJob("email", payload)
}

// CreatePasswordResetEmailJob creates a password reset email job
func (f *EmailJobFactory) CreatePasswordResetEmailJob(to, resetToken, resetURL string) (jobs.Job, error) {
	payload := jobs.JobPayload{
		"email_type":  "password_reset",
		"to":          to,
		"reset_token": resetToken,
		"reset_url":   resetURL,
	}

	factory := domainjobs.NewJobFactory()
	return factory.CreateJob("email", payload)
}

// CreateNotificationEmailJob creates a notification email job
func (f *EmailJobFactory) CreateNotificationEmailJob(to, subject, body string, priority jobs.JobPriority) (jobs.Job, error) {
	payload := jobs.JobPayload{
		"email_type": "notification",
		"to":         to,
		"subject":    subject,
		"body":       body,
	}

	opts := jobs.JobOptions{
		Priority: priority,
		Queue:    "notifications",
	}

	factory := domainjobs.NewJobFactory()
	return factory.CreateJobWithOptions("email", payload, opts)
}

// CreateBulkEmailJob creates a bulk email job
func (f *EmailJobFactory) CreateBulkEmailJob(recipients []string, subject, body string) (jobs.Job, error) {
	payload := jobs.JobPayload{
		"email_type": "bulk",
		"recipients": recipients,
		"subject":    subject,
		"body":       body,
	}

	opts := jobs.JobOptions{
		Priority: jobs.PriorityLow, // Bulk emails have lower priority
		Queue:    "bulk",
	}

	factory := domainjobs.NewJobFactory()
	return factory.CreateJobWithOptions("email", payload, opts)
}

// Mock email service for demonstration
type MockEmailService struct{}

// NewMockEmailService creates a new mock email service
func NewMockEmailService() *MockEmailService {
	return &MockEmailService{}
}

// SendEmail simulates sending a single email
func (m *MockEmailService) SendEmail(ctx context.Context, to, subject, body string) error {
	// Simulate processing time
	time.Sleep(100 * time.Millisecond)

	// Simulate occasional failures
	if to == "fail@example.com" {
		return fmt.Errorf("failed to send email to %s", to)
	}

	// In a real implementation, this would integrate with an email service
	// like SendGrid, Amazon SES, or SMTP
	fmt.Printf("📧 Email sent to %s: %s\n", to, subject)
	return nil
}

// SendBulkEmail simulates sending bulk emails
func (m *MockEmailService) SendBulkEmail(ctx context.Context, recipients []string, subject, body string) error {
	// Simulate longer processing time for bulk emails
	time.Sleep(time.Duration(len(recipients)) * 50 * time.Millisecond)

	for _, recipient := range recipients {
		if err := m.SendEmail(ctx, recipient, subject, body); err != nil {
			return err
		}
	}

	return nil
}

// SendTemplateEmail simulates sending template-based emails
func (m *MockEmailService) SendTemplateEmail(ctx context.Context, to, templateID string, data map[string]interface{}) error {
	// Simulate template processing time
	time.Sleep(150 * time.Millisecond)

	fmt.Printf("📧 Template email sent to %s using template %s\n", to, templateID)
	return nil
}

// BusinessEmailMetrics extends the basic metrics with email-specific metrics
type BusinessEmailMetrics struct {
	jobs.JobMetrics
	emailCounts map[string]int64
}

// NewBusinessEmailMetrics creates enhanced email metrics
func NewBusinessEmailMetrics(baseMetrics jobs.JobMetrics) *BusinessEmailMetrics {
	return &BusinessEmailMetrics{
		JobMetrics:  baseMetrics,
		emailCounts: make(map[string]int64),
	}
}

// RecordEmailJob records email-specific metrics
func (m *BusinessEmailMetrics) RecordEmailJob(emailType string, success bool) {
	if success {
		m.emailCounts[emailType]++
	}
}
