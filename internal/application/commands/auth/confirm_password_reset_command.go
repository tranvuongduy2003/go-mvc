package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/services"
)

// ConfirmPasswordResetCommand represents the password reset confirmation command
type ConfirmPasswordResetCommand struct {
	Token       string `validate:"required"`
	NewPassword string `validate:"required,min=8"`
}

// ConfirmPasswordResetCommandHandler handles the ConfirmPasswordResetCommand
type ConfirmPasswordResetCommandHandler struct {
	authService services.AuthService
}

// NewConfirmPasswordResetCommandHandler creates a new ConfirmPasswordResetCommandHandler
func NewConfirmPasswordResetCommandHandler(authService services.AuthService) *ConfirmPasswordResetCommandHandler {
	return &ConfirmPasswordResetCommandHandler{
		authService: authService,
	}
}

// Handle executes the ConfirmPasswordResetCommand
func (h *ConfirmPasswordResetCommandHandler) Handle(ctx context.Context, cmd ConfirmPasswordResetCommand) (*dto.StatusResponse, error) {
	// Confirm password reset
	err := h.authService.ConfirmPasswordReset(ctx, cmd.Token, cmd.NewPassword)
	if err != nil {
		return nil, err
	}

	// Return success response
	return &dto.StatusResponse{
		Status:  "success",
		Message: "Password has been reset successfully",
	}, nil
}
