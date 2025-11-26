package external

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/logger"
)

// Compile-time check to ensure FileStorageService implements the port interface
var _ contracts.FileStorageService = (*FileStorageService)(nil)

// FileStorageService handles MinIO S3 file storage
// Implements contracts.FileStorageService port interface
type FileStorageService struct {
	client     *minio.Client
	bucketName string
	cdnURL     string
	logger     *logger.Logger
}

// UploadResult represents a file upload result
type UploadResult struct {
	FileKey string `json:"file_key"`
	CDNUrl  string `json:"cdn_url"`
	Size    int64  `json:"size"`
}

// FileStorageConfig holds MinIO configuration
type FileStorageConfig struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	BucketName      string `yaml:"bucket_name"`
	CDNUrl          string `yaml:"cdn_url"`
	UseSSL          bool   `yaml:"use_ssl"`
}

// NewFileStorageService creates a new MinIO file storage service
func NewFileStorageService(cfg *FileStorageConfig, logger *logger.Logger) (*FileStorageService, error) {
	// Initialize MinIO client
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	service := &FileStorageService{
		client:     minioClient,
		bucketName: cfg.BucketName,
		cdnURL:     cfg.CDNUrl,
		logger:     logger,
	}

	// Ensure bucket exists
	if err := service.ensureBucketExists(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ensure bucket exists: %w", err)
	}

	return service, nil
}

// ensureBucketExists creates bucket if it doesn't exist
func (s *FileStorageService) ensureBucketExists(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		s.logger.Infof("Creating bucket: %s", s.bucketName)
		err := s.client.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}

		// Set bucket policy for public read access
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`, s.bucketName)

		err = s.client.SetBucketPolicy(ctx, s.bucketName, policy)
		if err != nil {
			s.logger.Warnf("Failed to set bucket policy: %v", err)
		}
	}

	return nil
}

// Upload implements contracts.FileStorageService.Upload
// Uploads a file to storage and returns the file key and CDN URL
func (s *FileStorageService) Upload(ctx context.Context, file io.Reader, filename string, contentType string, size int64) (fileKey string, cdnURL string, err error) {
	// Validate file type for images
	if !s.isValidImageType(contentType) {
		return "", "", fmt.Errorf("invalid file type: %s", contentType)
	}

	// Generate unique file key
	fileExt := filepath.Ext(filename)
	fileKey = fmt.Sprintf("uploads/%s%s", uuid.New().String(), fileExt)

	// Upload file
	info, err := s.client.PutObject(ctx, s.bucketName, fileKey, file, size, minio.PutObjectOptions{
		ContentType: contentType,
		UserMetadata: map[string]string{
			"uploaded-at": time.Now().Format(time.RFC3339),
		},
	})
	if err != nil {
		s.logger.Errorf("Failed to upload file: %v", err)
		return "", "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Generate CDN URL
	cdnURL = fmt.Sprintf("%s/%s/%s", s.cdnURL, s.bucketName, fileKey)

	s.logger.Infof("File uploaded successfully: key=%s, size=%d", fileKey, info.Size)

	return fileKey, cdnURL, nil
}

// Delete implements contracts.FileStorageService.Delete
// Removes a file from storage
func (s *FileStorageService) Delete(ctx context.Context, fileKey string) error {
	s.logger.Infof("Deleting file: %s", fileKey)

	err := s.client.RemoveObject(ctx, s.bucketName, fileKey, minio.RemoveObjectOptions{})
	if err != nil {
		s.logger.Errorf("Failed to delete file %s: %v", fileKey, err)
		return fmt.Errorf("failed to delete file: %w", err)
	}

	s.logger.Infof("File deleted successfully: %s", fileKey)
	return nil
}

// GetURL implements contracts.FileStorageService.GetURL
// Returns a presigned/public URL for a file
func (s *FileStorageService) GetURL(ctx context.Context, fileKey string) (string, error) {
	// For public buckets, return the CDN URL directly
	return fmt.Sprintf("%s/%s/%s", s.cdnURL, s.bucketName, fileKey), nil
}

// Exists implements contracts.FileStorageService.Exists
// Checks if a file exists in storage
func (s *FileStorageService) Exists(ctx context.Context, fileKey string) (bool, error) {
	_, err := s.client.StatObject(ctx, s.bucketName, fileKey, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}
	return true, nil
}

// UploadAvatar uploads user avatar image (legacy method for backward compatibility)
func (s *FileStorageService) UploadAvatar(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (*UploadResult, error) {
	// Validate file type
	if !s.isValidImageType(header.Header.Get("Content-Type")) {
		return nil, fmt.Errorf("invalid file type: %s", header.Header.Get("Content-Type"))
	}

	// Generate unique file key
	fileExt := filepath.Ext(header.Filename)
	fileKey := fmt.Sprintf("avatars/%s/%s%s", userID, uuid.New().String(), fileExt)

	// Upload file
	info, err := s.client.PutObject(ctx, s.bucketName, fileKey, file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
		UserMetadata: map[string]string{
			"user-id":     userID,
			"upload-type": "avatar",
			"uploaded-at": time.Now().Format(time.RFC3339),
		},
	})
	if err != nil {
		s.logger.Errorf("Failed to upload avatar for user %s: %v", userID, err)
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Generate CDN URL
	cdnURL := fmt.Sprintf("%s/%s/%s", s.cdnURL, s.bucketName, fileKey)

	s.logger.Infof("Avatar uploaded successfully for user %s: key=%s, size=%d", userID, fileKey, info.Size)

	return &UploadResult{
		FileKey: fileKey,
		CDNUrl:  cdnURL,
		Size:    info.Size,
	}, nil
}

// DeleteFile deletes a file from MinIO (legacy method - use Delete instead)
func (s *FileStorageService) DeleteFile(ctx context.Context, fileKey string) error {
	return s.Delete(ctx, fileKey)
}

// GetFileInfo gets file information
func (s *FileStorageService) GetFileInfo(ctx context.Context, fileKey string) (*minio.ObjectInfo, error) {
	info, err := s.client.StatObject(ctx, s.bucketName, fileKey, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	return &info, nil
}

// isValidImageType checks if the content type is a valid image format
func (s *FileStorageService) isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	contentType = strings.ToLower(contentType)
	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}
