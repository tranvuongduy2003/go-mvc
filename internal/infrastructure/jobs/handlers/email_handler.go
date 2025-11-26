package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/job"
)

type EmailJobHandler struct {
	emailService EmailService
	metrics      job.JobMetrics
}

type EmailService interface {
	SendEmail(ctx context.Context, to, subject, body string) error
	SendBulkEmail(ctx context.Context, recipients []string, subject, body string) error
	SendTemplateEmail(ctx context.Context, to, templateID string, data map[string]interface{}) error
}

func NewEmailJobHandler(emailService EmailService, metrics job.JobMetrics) *EmailJobHandler {
	return &EmailJobHandler{
		emailService: emailService,
		metrics:      metrics,
	}
}

func (h *EmailJobHandler) Execute(ctx context.Context, excutedJob job.Job) error {
	start := time.Now()
	defer func() {
		if h.metrics != nil {
			h.metrics.ObserveJobDuration(excutedJob.GetType(), time.Since(start))
		}
	}()

	emailJob, ok := excutedJob.(*job.EmailJob)
	if !ok {
		err := fmt.Errorf("expected EmailJob, got %T", excutedJob)
		if h.metrics != nil {
			h.metrics.IncrementJobsProcessed(excutedJob.GetType(), false)
		}
		return err
	}

	payload := emailJob.GetPayload()
	emailType, _ := payload["email_type"].(string)
	to, _ := payload["to"].(string)
	subject, _ := payload["subject"].(string)
	body, _ := payload["body"].(string)

	var err error

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
		err = h.emailService.SendEmail(ctx, to, subject, body)
	}

	if h.metrics != nil {
		success := err == nil
		h.metrics.IncrementJobsProcessed(excutedJob.GetType(), success)

		if businessMetrics, ok := h.metrics.(*BusinessEmailMetrics); ok {
			businessMetrics.RecordEmailJob(emailType, success)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (h *EmailJobHandler) GetJobType() string {
	return "email"
}

func (h *EmailJobHandler) handleWelcomeEmail(ctx context.Context, to string, payload job.JobPayload) error {
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

func (h *EmailJobHandler) handlePasswordResetEmail(ctx context.Context, to string, payload job.JobPayload) error {
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

func (h *EmailJobHandler) handleBulkEmail(ctx context.Context, payload job.JobPayload) error {
	recipients, _ := payload["recipients"].([]string)
	subject, _ := payload["subject"].(string)
	body, _ := payload["body"].(string)

	if len(recipients) == 0 {
		return fmt.Errorf("no recipients specified for bulk email")
	}

	return h.emailService.SendBulkEmail(ctx, recipients, subject, body)
}

func (h *EmailJobHandler) handleTemplateEmail(ctx context.Context, to string, payload job.JobPayload) error {
	templateID, _ := payload["template_id"].(string)
	templateData, _ := payload["template_data"].(map[string]interface{})

	if templateID == "" {
		return fmt.Errorf("template_id is required for template emails")
	}

	return h.emailService.SendTemplateEmail(ctx, to, templateID, templateData)
}

type EmailJobFactory struct{}

func NewEmailJobFactory() *EmailJobFactory {
	return &EmailJobFactory{}
}

func (f *EmailJobFactory) CreateWelcomeEmailJob(to, username string) (job.Job, error) {
	payload := job.JobPayload{
		"email_type": "welcome",
		"to":         to,
		"username":   username,
	}

	factory := job.NewJobFactory()
	return factory.CreateJob("email", payload)
}

func (f *EmailJobFactory) CreatePasswordResetEmailJob(to, resetToken, resetURL string) (job.Job, error) {
	payload := job.JobPayload{
		"email_type":  "password_reset",
		"to":          to,
		"reset_token": resetToken,
		"reset_url":   resetURL,
	}

	factory := job.NewJobFactory()
	return factory.CreateJob("email", payload)
}

func (f *EmailJobFactory) CreateNotificationEmailJob(to, subject, body string, priority job.JobPriority) (job.Job, error) {
	payload := job.JobPayload{
		"email_type": "notification",
		"to":         to,
		"subject":    subject,
		"body":       body,
	}

	opts := job.JobOptions{
		Priority: priority,
		Queue:    "notifications",
	}

	factory := job.NewJobFactory()
	return factory.CreateJobWithOptions("email", payload, opts)
}

func (f *EmailJobFactory) CreateBulkEmailJob(recipients []string, subject, body string) (job.Job, error) {
	payload := job.JobPayload{
		"email_type": "bulk",
		"recipients": recipients,
		"subject":    subject,
		"body":       body,
	}

	opts := job.JobOptions{
		Priority: job.PriorityLow, // Bulk emails have lower priority
		Queue:    "bulk",
	}

	factory := job.NewJobFactory()
	return factory.CreateJobWithOptions("email", payload, opts)
}

type MockEmailService struct{}

func NewMockEmailService() *MockEmailService {
	return &MockEmailService{}
}

func (m *MockEmailService) SendEmail(ctx context.Context, to, subject, body string) error {
	time.Sleep(100 * time.Millisecond)

	if to == "fail@example.com" {
		return fmt.Errorf("failed to send email to %s", to)
	}

	fmt.Printf("📧 Email sent to %s: %s\n", to, subject)
	return nil
}

func (m *MockEmailService) SendBulkEmail(ctx context.Context, recipients []string, subject, body string) error {
	time.Sleep(time.Duration(len(recipients)) * 50 * time.Millisecond)

	for _, recipient := range recipients {
		if err := m.SendEmail(ctx, recipient, subject, body); err != nil {
			return err
		}
	}

	return nil
}

func (m *MockEmailService) SendTemplateEmail(ctx context.Context, to, templateID string, data map[string]interface{}) error {
	time.Sleep(150 * time.Millisecond)

	fmt.Printf("📧 Template email sent to %s using template %s\n", to, templateID)
	return nil
}

type BusinessEmailMetrics struct {
	job.JobMetrics
	emailCounts map[string]int64
}

func NewBusinessEmailMetrics(baseMetrics job.JobMetrics) *BusinessEmailMetrics {
	return &BusinessEmailMetrics{
		JobMetrics:  baseMetrics,
		emailCounts: make(map[string]int64),
	}
}

func (m *BusinessEmailMetrics) RecordEmailJob(emailType string, success bool) {
	if success {
		m.emailCounts[emailType]++
	}
}
