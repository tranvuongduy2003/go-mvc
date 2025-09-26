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
