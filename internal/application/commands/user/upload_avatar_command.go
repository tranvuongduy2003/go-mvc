package commands

import (
	"context"
	"mime/multipart"

	"github.com/tranvuongduy2003/go-mvc/internal/adapters/external"
	userDto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/user"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/shared/events"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/messaging"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/repositories"
	apperrors "github.com/tranvuongduy2003/go-mvc/pkg/errors"
)

// UploadAvatarCommand represents a command to upload user avatar
type UploadAvatarCommand struct {
	UserID string
	File   multipart.File
	Header *multipart.FileHeader
}

// UploadAvatarCommandHandler handles uploading user avatar
type UploadAvatarCommandHandler struct {
	userRepo           repositories.UserRepository
	fileStorageService *external.FileStorageService
	eventBus           messaging.EventBus
}

// NewUploadAvatarCommandHandler creates a new UploadAvatarCommandHandler
func NewUploadAvatarCommandHandler(
	userRepo repositories.UserRepository,
	fileStorageService *external.FileStorageService,
	eventBus messaging.EventBus,
) *UploadAvatarCommandHandler {
	return &UploadAvatarCommandHandler{
		userRepo:           userRepo,
		fileStorageService: fileStorageService,
		eventBus:           eventBus,
	}
}

// Handle executes the upload avatar command
func (h *UploadAvatarCommandHandler) Handle(ctx context.Context, cmd UploadAvatarCommand) (userDto.UserResponse, error) {
	// Get user
	user, err := h.userRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return userDto.UserResponse{}, apperrors.NewInternalError("Failed to get user", err)
	}
	if user == nil {
		return userDto.UserResponse{}, apperrors.NewNotFoundError("User not found")
	}

	// Delete old avatar if exists
	if !user.Avatar().IsEmpty() {
		_ = h.fileStorageService.DeleteFile(ctx, user.Avatar().FileKey()) // Ignore error, just log it
	}

	// Upload new avatar
	uploadResult, err := h.fileStorageService.UploadAvatar(ctx, cmd.UserID, cmd.File, cmd.Header)
	if err != nil {
		return userDto.UserResponse{}, apperrors.NewInternalError("Failed to upload avatar", err)
	}

	// Update user avatar
	user.UpdateAvatar(uploadResult.FileKey, uploadResult.CDNUrl)

	// Save updated user
	if err := h.userRepo.Update(ctx, user); err != nil {
		// If save failed, try to delete uploaded file
		_ = h.fileStorageService.DeleteFile(ctx, uploadResult.FileKey)
		return userDto.UserResponse{}, apperrors.NewInternalError("Failed to update user", err)
	}

	// Publish avatar uploaded event
	avatarEvent := events.NewUserAvatarUploadedEvent(
		user.ID(),
		uploadResult.CDNUrl,
		uploadResult.FileKey,
		1, // Version - could be incremented from user version
	)

	if err := h.eventBus.PublishEvent(ctx, avatarEvent); err != nil {
		// Log error but don't fail the operation
		// The avatar was successfully uploaded and saved
		// Event publishing failure is not critical
	}

	// Convert to DTO
	return userDto.UserResponse{
		ID:        user.ID(),
		Email:     user.Email(),
		Name:      user.Name(),
		Phone:     user.Phone(),
		AvatarURL: user.Avatar().CDNUrl(),
		IsActive:  user.IsActive(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}, nil
}
