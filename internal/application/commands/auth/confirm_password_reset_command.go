package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/ports/services"
)

// ConfirmPasswordResetCommand represents the password reset confirmation command
type ConfirmPasswordResetCommand struct {
	Token       string `validate:"required"`
	NewPassword string `validate:"required,min=8"`
}

// ConfirmPasswordResetCommandHandler handles the ConfirmPasswordResetCommand
type ConfirmPasswordResetCommandHandler struct {
	passwordService services.PasswordManagementService
}

// NewConfirmPasswordResetCommandHandler creates a new ConfirmPasswordResetCommandHandler
func NewConfirmPasswordResetCommandHandler(passwordService services.PasswordManagementService) *ConfirmPasswordResetCommandHandler {
	return &ConfirmPasswordResetCommandHandler{
		passwordService: passwordService,
	}
}

// Handle executes the ConfirmPasswordResetCommand
func (h *ConfirmPasswordResetCommandHandler) Handle(ctx context.Context, cmd ConfirmPasswordResetCommand) (*dto.StatusResponse, error) {
	// Confirm password reset
	err := h.passwordService.ConfirmPasswordReset(ctx, cmd.Token, cmd.NewPassword)
	if err != nil {
		return nil, err
	}

	// Return success response
	return &dto.StatusResponse{
		Status:  "success",
		Message: "Password has been reset successfully",
	}, nil
}
