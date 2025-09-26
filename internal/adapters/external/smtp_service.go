package external

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/tranvuongduy2003/go-mvc/internal/shared/config"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
)

// SMTPService handles SMTP email sending
type SMTPService struct {
	host     string
	port     string
	username string
	password string
	from     string
	useTLS   bool
	logger   *logger.Logger
}

// NewSMTPService creates a new SMTP email service
func NewSMTPService(cfg *config.SMTPConfig, logger *logger.Logger) *SMTPService {
	return &SMTPService{
		host:     cfg.Host,
		port:     cfg.Port,
		username: cfg.Username,
		password: cfg.Password,
		from:     cfg.From,
		useTLS:   cfg.UseTLS,
		logger:   logger,
	}
}

// SendEmail sends an email using SMTP
func (s *SMTPService) SendEmail(ctx context.Context, to []string, subject, body string) error {
	s.logger.Infof("🚀 SMTP DEBUG: SendEmail called - to: %v, subject: %s", to, subject)
	s.logger.Infof("🚀 SMTP DEBUG: SMTP Config - %s:%s (TLS: %t)", s.host, s.port, s.useTLS)

	// Create the email message
	message := s.buildMessage(to, subject, body)
	s.logger.Infof("🚀 SMTP DEBUG: Built message length: %d", len(message))

	// Set up authentication (if required)
	var auth smtp.Auth
	if s.username != "" && s.password != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
		s.logger.Infof("🚀 SMTP DEBUG: Using auth with username: %s", s.username)
	} else {
		s.logger.Infof("🚀 SMTP DEBUG: No auth - username/password empty")
	}

	// Send the email
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	s.logger.Infof("🚀 SMTP DEBUG: Attempting to send to: %s", addr)
	err := smtp.SendMail(addr, auth, s.from, to, []byte(message))
	if err != nil {
		s.logger.Errorf("❌ SMTP ERROR: Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Infof("✅ SMTP SUCCESS: Email sent successfully to: %v", to)
	return nil
}

// SendVerificationEmail sends an email verification message
func (s *SMTPService) SendVerificationEmail(ctx context.Context, to, firstName, verificationToken string) error {
	subject := "Verify Your Email Address"
	body := fmt.Sprintf(`
Hello %s,

Please click the link below to verify your email address:

http://localhost:8080/api/v1/auth/verify-email?token=%s

If you didn't create an account, please ignore this email.

Best regards,
The Team
`, firstName, verificationToken)

	return s.SendEmail(ctx, []string{to}, subject, body)
}

// SendPasswordResetEmail sends a password reset email
func (s *SMTPService) SendPasswordResetEmail(ctx context.Context, to, firstName, resetToken string) error {
	s.logger.Infof("🔥 SMTP DEBUG: SendPasswordResetEmail called with to=%s, firstName=%s, token=%s", to, firstName, resetToken)
	s.logger.Infof("🔥 SMTP DEBUG: Config - host=%s, port=%s, from=%s, useTLS=%t", s.host, s.port, s.from, s.useTLS)

	subject := "Password Reset Request"
	body := fmt.Sprintf(`
Hello %s,

You requested a password reset. Please click the link below to reset your password:

http://localhost:8080/api/v1/auth/confirm-reset?token=%s

This link will expire in 1 hour. If you didn't request this, please ignore this email.

Best regards,
The Team
`, firstName, resetToken)

	return s.SendEmail(ctx, []string{to}, subject, body)
}

// buildMessage constructs the email message with proper headers
func (s *SMTPService) buildMessage(to []string, subject, body string) string {
	message := fmt.Sprintf("To: %s\r\n", to[0])
	if len(to) > 1 {
		for _, addr := range to[1:] {
			message += fmt.Sprintf("Cc: %s\r\n", addr)
		}
	}
	message += fmt.Sprintf("From: %s\r\n", s.from)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "Content-Type: text/plain; charset=UTF-8\r\n"
	message += fmt.Sprintf("\r\n%s", body)
	return message
}
