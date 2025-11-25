package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/logger"
)

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
