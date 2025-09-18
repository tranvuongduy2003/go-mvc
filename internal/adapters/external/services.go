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

// SMSService handles external SMS service integration
type SMSService struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	logger     *logger.Logger
}

// SMSRequest represents an SMS request
type SMSRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
	From    string `json:"from,omitempty"`
}

// SMSResponse represents an SMS service response
type SMSResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// NewSMSService creates a new SMS service
func NewSMSService(apiKey, baseURL string, logger *logger.Logger) *SMSService {
	return &SMSService{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// SendSMS sends an SMS through external service
func (s *SMSService) SendSMS(ctx context.Context, req *SMSRequest) (*SMSResponse, error) {
	s.logger.Infof("Sending SMS to: %s", req.To)

	// Marshal request body
	body, err := json.Marshal(req)
	if err != nil {
		s.logger.Errorf("Failed to marshal SMS request: %v", err)
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
		s.logger.Errorf("Failed to send SMS request: %v", err)
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
		s.logger.Errorf("SMS service returned error: status=%d, body=%s", resp.StatusCode, string(respBody))
		return nil, fmt.Errorf("SMS service error: status=%d", resp.StatusCode)
	}

	// Parse response
	var smsResp SMSResponse
	if err := json.Unmarshal(respBody, &smsResp); err != nil {
		s.logger.Errorf("Failed to unmarshal SMS response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	s.logger.Infof("SMS sent successfully: message_id=%s", smsResp.MessageID)
	return &smsResp, nil
}

// PushNotificationService handles push notification service integration
type PushNotificationService struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	logger     *logger.Logger
}

// PushRequest represents a push notification request
type PushRequest struct {
	Tokens   []string          `json:"tokens"`
	Title    string            `json:"title"`
	Body     string            `json:"body"`
	Data     map[string]string `json:"data,omitempty"`
	Priority string            `json:"priority,omitempty"`
}

// PushResponse represents a push notification response
type PushResponse struct {
	SuccessCount int      `json:"success_count"`
	FailureCount int      `json:"failure_count"`
	Results      []string `json:"results"`
}

// NewPushNotificationService creates a new push notification service
func NewPushNotificationService(apiKey, baseURL string, logger *logger.Logger) *PushNotificationService {
	return &PushNotificationService{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// SendPushNotification sends a push notification through external service
func (s *PushNotificationService) SendPushNotification(ctx context.Context, req *PushRequest) (*PushResponse, error) {
	s.logger.Infof("Sending push notification to %d devices: %s", len(req.Tokens), req.Title)

	// Marshal request body
	body, err := json.Marshal(req)
	if err != nil {
		s.logger.Errorf("Failed to marshal push notification request: %v", err)
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
		s.logger.Errorf("Failed to send push notification request: %v", err)
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
		s.logger.Errorf("Push notification service returned error: status=%d, body=%s", resp.StatusCode, string(respBody))
		return nil, fmt.Errorf("push notification service error: status=%d", resp.StatusCode)
	}

	// Parse response
	var pushResp PushResponse
	if err := json.Unmarshal(respBody, &pushResp); err != nil {
		s.logger.Errorf("Failed to unmarshal push notification response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	s.logger.Infof("Push notification sent: success=%d, failure=%d", pushResp.SuccessCount, pushResp.FailureCount)
	return &pushResp, nil
}

// FileStorageService handles external file storage service integration
type FileStorageService struct {
	apiKey     string
	baseURL    string
	bucketName string
	httpClient *http.Client
	logger     *logger.Logger
}

// UploadRequest represents a file upload request
type UploadRequest struct {
	FileName    string            `json:"file_name"`
	ContentType string            `json:"content_type"`
	Data        []byte            `json:"data"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// UploadResponse represents a file upload response
type UploadResponse struct {
	FileID   string `json:"file_id"`
	URL      string `json:"url"`
	FileName string `json:"file_name"`
	Size     int64  `json:"size"`
}

// NewFileStorageService creates a new file storage service
func NewFileStorageService(apiKey, baseURL, bucketName string, logger *logger.Logger) *FileStorageService {
	return &FileStorageService{
		apiKey:     apiKey,
		baseURL:    baseURL,
		bucketName: bucketName,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute, // Longer timeout for file uploads
		},
		logger: logger,
	}
}

// UploadFile uploads a file to external storage service
func (s *FileStorageService) UploadFile(ctx context.Context, req *UploadRequest) (*UploadResponse, error) {
	s.logger.Infof("Uploading file: %s, size: %d bytes", req.FileName, len(req.Data))

	// Marshal request body
	body, err := json.Marshal(req)
	if err != nil {
		s.logger.Errorf("Failed to marshal upload request: %v", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/upload", bytes.NewBuffer(body))
	if err != nil {
		s.logger.Errorf("Failed to create HTTP request: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)
	httpReq.Header.Set("X-Bucket-Name", s.bucketName)

	// Send request
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		s.logger.Errorf("Failed to send upload request: %v", err)
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
		s.logger.Errorf("File storage service returned error: status=%d, body=%s", resp.StatusCode, string(respBody))
		return nil, fmt.Errorf("file storage service error: status=%d", resp.StatusCode)
	}

	// Parse response
	var uploadResp UploadResponse
	if err := json.Unmarshal(respBody, &uploadResp); err != nil {
		s.logger.Errorf("Failed to unmarshal upload response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	s.logger.Infof("File uploaded successfully: file_id=%s, url=%s", uploadResp.FileID, uploadResp.URL)
	return &uploadResp, nil
}

// DeleteFile deletes a file from external storage service
func (s *FileStorageService) DeleteFile(ctx context.Context, fileID string) error {
	s.logger.Infof("Deleting file: %s", fileID)

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", s.baseURL+"/files/"+fileID, nil)
	if err != nil {
		s.logger.Errorf("Failed to create HTTP request: %v", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)
	httpReq.Header.Set("X-Bucket-Name", s.bucketName)

	// Send request
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		s.logger.Errorf("Failed to send delete request: %v", err)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		s.logger.Errorf("File storage service returned error: status=%d, body=%s", resp.StatusCode, string(respBody))
		return fmt.Errorf("file storage service error: status=%d", resp.StatusCode)
	}

	s.logger.Infof("File deleted successfully: %s", fileID)
	return nil
}
