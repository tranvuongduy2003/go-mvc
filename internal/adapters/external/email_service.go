package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
)

// EmailService handles external email service integration
type EmailService struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	logger     *logger.Logger
}

// EmailRequest represents an email request
type EmailRequest struct {
	To       []string          `json:"to"`
	Subject  string            `json:"subject"`
	Body     string            `json:"body"`
	IsHTML   bool              `json:"is_html"`
	From     string            `json:"from,omitempty"`
	Template string            `json:"template,omitempty"`
	Data     map[string]string `json:"data,omitempty"`
}

// EmailResponse represents an email service response
type EmailResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// NewEmailService creates a new email service
func NewEmailService(apiKey, baseURL string, logger *logger.Logger) *EmailService {
	return &EmailService{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// SendEmail sends an email through external service
func (s *EmailService) SendEmail(ctx context.Context, req *EmailRequest) (*EmailResponse, error) {
	s.logger.Infof("Sending email to: %v, subject: %s", req.To, req.Subject)

	// Marshal request body
	body, err := json.Marshal(req)
	if err != nil {
		s.logger.Errorf("Failed to marshal email request: %v", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/send", bytes.NewBuffer(body))
	if err != nil {
		s.logger.Errorf("Failed to create HTTP request: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	// Send request
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		s.logger.Errorf("Failed to send email request: %v", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Errorf("Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		s.logger.Errorf("Email service returned error: status=%d, body=%s", resp.StatusCode, string(respBody))
		return nil, fmt.Errorf("email service error: status=%d", resp.StatusCode)
	}

	// Parse response
	var emailResp EmailResponse
	if err := json.Unmarshal(respBody, &emailResp); err != nil {
		s.logger.Errorf("Failed to unmarshal email response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	s.logger.Infof("Email sent successfully: message_id=%s", emailResp.MessageID)
	return &emailResp, nil
}

// SendWelcomeEmail sends a welcome email to new users
func (s *EmailService) SendWelcomeEmail(ctx context.Context, to, firstName string) error {
	req := &EmailRequest{
		To:       []string{to},
		Subject:  "Welcome to Our Platform!",
		Template: "welcome",
		Data: map[string]string{
			"first_name": firstName,
		},
	}

	_, err := s.SendEmail(ctx, req)
	return err
}

// SendPasswordResetEmail sends a password reset email
func (s *EmailService) SendPasswordResetEmail(ctx context.Context, to, resetToken string) error {
	req := &EmailRequest{
		To:       []string{to},
		Subject:  "Password Reset Request",
		Template: "password_reset",
		Data: map[string]string{
			"reset_token": resetToken,
		},
	}

	_, err := s.SendEmail(ctx, req)
	return err
}
