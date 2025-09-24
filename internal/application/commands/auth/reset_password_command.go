package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/services"
)

// ResetPasswordCommand represents the reset password initiation command
type ResetPasswordCommand struct {
	Email string `validate:"required,email"`
}

// ResetPasswordCommandHandler handles the ResetPasswordCommand
type ResetPasswordCommandHandler struct {
	authService services.AuthService
}

// NewResetPasswordCommandHandler creates a new ResetPasswordCommandHandler
func NewResetPasswordCommandHandler(authService services.AuthService) *ResetPasswordCommandHandler {
	return &ResetPasswordCommandHandler{
		authService: authService,
	}
}

// Handle executes the ResetPasswordCommand
func (h *ResetPasswordCommandHandler) Handle(ctx context.Context, cmd ResetPasswordCommand) (*dto.StatusResponse, error) {
	// Initiate password reset
	err := h.authService.ResetPassword(ctx, cmd.Email)
	if err != nil {
		return nil, err
	}

	// Return success response
	return &dto.StatusResponse{
		Status:  "success",
		Message: "Password reset instructions have been sent to your email",
	}, nil
}
