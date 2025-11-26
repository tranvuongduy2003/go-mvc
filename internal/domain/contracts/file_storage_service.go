package contracts

import (
	"context"
	"io"
)

// FileStorageService defines the interface for file storage operations
// This abstraction allows the domain/application layer to remain independent
// of specific storage implementations (S3, Azure Blob, GCS, etc.)
type FileStorageService interface {
	// Upload uploads a file to storage and returns the file key and CDN URL
	Upload(ctx context.Context, file io.Reader, filename string, contentType string, size int64) (fileKey string, cdnURL string, err error)

	// Delete removes a file from storage
	Delete(ctx context.Context, fileKey string) error

	// GetURL returns a presigned/public URL for a file
	GetURL(ctx context.Context, fileKey string) (string, error)

	// Exists checks if a file exists in storage
	Exists(ctx context.Context, fileKey string) (bool, error)
}

// EmailService defines the interface for email operations
// Separates email concerns from authentication logic
type EmailService interface {
	// SendVerificationEmail sends email verification link
	SendVerificationEmail(ctx context.Context, to, token string) error

	// SendPasswordResetEmail sends password reset link
	SendPasswordResetEmail(ctx context.Context, to, token string) error

	// SendWelcomeEmail sends welcome email to new users
	SendWelcomeEmail(ctx context.Context, to, name string) error
}

// SMSService defines the interface for SMS operations
type SMSService interface {
	// SendVerificationCode sends SMS verification code
	SendVerificationCode(ctx context.Context, phoneNumber, code string) error

	// SendPasswordResetCode sends password reset code via SMS
	SendPasswordResetCode(ctx context.Context, phoneNumber, code string) error
}

// PushNotificationService defines the interface for push notification operations
type PushNotificationService interface {
	// SendNotification sends a push notification to a device
	SendNotification(ctx context.Context, deviceToken, title, body string, data map[string]interface{}) error

	// SendToMultipleDevices sends notification to multiple devices
	SendToMultipleDevices(ctx context.Context, deviceTokens []string, title, body string, data map[string]interface{}) error
}
