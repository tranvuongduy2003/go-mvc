package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
)

// ChangePasswordCommand represents the change password command
type ChangePasswordCommand struct {
	UserID      string `validate:"required"`
	OldPassword string `validate:"required"`
	NewPassword string `validate:"required,min=8"`
}

// ChangePasswordCommandHandler handles the ChangePasswordCommand
type ChangePasswordCommandHandler struct {
	passwordService contracts.PasswordManagementService
}

// NewChangePasswordCommandHandler creates a new ChangePasswordCommandHandler
func NewChangePasswordCommandHandler(passwordService contracts.PasswordManagementService) *ChangePasswordCommandHandler {
	return &ChangePasswordCommandHandler{
		passwordService: passwordService,
	}
}

// Handle executes the ChangePasswordCommand
func (h *ChangePasswordCommandHandler) Handle(ctx context.Context, cmd ChangePasswordCommand) (*dto.StatusResponse, error) {
	// Change password
	err := h.passwordService.ChangePassword(ctx, cmd.UserID, cmd.OldPassword, cmd.NewPassword)
	if err != nil {
		return nil, err
	}

	// Return success response
	return &dto.StatusResponse{
		Status:  "success",
		Message: "Password changed successfully",
	}, nil
}
