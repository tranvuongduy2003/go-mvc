package contracts

import (
	"context"
	"io"
)

type FileStorageService interface {
	Upload(ctx context.Context, file io.Reader, filename string, contentType string, size int64) (fileKey string, cdnURL string, err error)

	Delete(ctx context.Context, fileKey string) error

	GetURL(ctx context.Context, fileKey string) (string, error)

	Exists(ctx context.Context, fileKey string) (bool, error)
}

type EmailService interface {
	SendVerificationEmail(ctx context.Context, to, token string) error

	SendPasswordResetEmail(ctx context.Context, to, token string) error

	SendWelcomeEmail(ctx context.Context, to, name string) error
}

type SMSService interface {
	SendVerificationCode(ctx context.Context, phoneNumber, code string) error

	SendPasswordResetCode(ctx context.Context, phoneNumber, code string) error
}

type PushNotificationService interface {
	SendNotification(ctx context.Context, deviceToken, title, body string, data map[string]interface{}) error

	SendToMultipleDevices(ctx context.Context, deviceTokens []string, title, body string, data map[string]interface{}) error
}
