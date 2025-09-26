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
