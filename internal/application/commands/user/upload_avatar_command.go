package commands

import (
	"context"
	"mime/multipart"

	userDto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/user"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/messaging"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/shared/events"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/user"
	apperrors "github.com/tranvuongduy2003/go-mvc/pkg/errors"
)

type UploadAvatarCommand struct {
	UserID string
	File   multipart.File
	Header *multipart.FileHeader
}

type UploadAvatarCommandHandler struct {
	userRepo           user.UserRepository
	fileStorageService contracts.FileStorageService // Changed to use port interface
	eventBus           messaging.EventBus
}

func NewUploadAvatarCommandHandler(
	userRepo user.UserRepository,
	fileStorageService contracts.FileStorageService, // Changed parameter type
	eventBus messaging.EventBus,
) *UploadAvatarCommandHandler {
	return &UploadAvatarCommandHandler{
		userRepo:           userRepo,
		fileStorageService: fileStorageService,
		eventBus:           eventBus,
	}
}

func (h *UploadAvatarCommandHandler) Handle(ctx context.Context, cmd UploadAvatarCommand) (userDto.UserResponse, error) {
	user, err := h.userRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return userDto.UserResponse{}, apperrors.NewInternalError("Failed to get user", err)
	}
	if user == nil {
		return userDto.UserResponse{}, apperrors.NewNotFoundError("User not found")
	}

	if !user.Avatar().IsEmpty() {
		_ = h.fileStorageService.Delete(ctx, user.Avatar().FileKey()) // Ignore error, just log it
	}

	fileKey, cdnURL, err := h.fileStorageService.Upload(
		ctx,
		cmd.File,
		cmd.Header.Filename,
		cmd.Header.Header.Get("Content-Type"),
		cmd.Header.Size,
	)
	if err != nil {
		return userDto.UserResponse{}, apperrors.NewInternalError("Failed to upload avatar", err)
	}

	user.UpdateAvatar(fileKey, cdnURL)

	if err := h.userRepo.Update(ctx, user); err != nil {
		_ = h.fileStorageService.Delete(ctx, fileKey)
		return userDto.UserResponse{}, apperrors.NewInternalError("Failed to update user", err)
	}

	avatarEvent := events.NewUserAvatarUploadedEvent(
		user.ID(),
		cdnURL,
		fileKey,
		int(user.Version()), // Use actual user version
	)

	if err := h.eventBus.PublishEvent(ctx, avatarEvent); err != nil {
	}

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
