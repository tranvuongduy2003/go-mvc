package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/services"
)

// ResendVerificationEmailCommand represents the resend verification email command
type ResendVerificationEmailCommand struct {
	Email string `validate:"required,email"`
}

// ResendVerificationEmailCommandHandler handles the ResendVerificationEmailCommand
type ResendVerificationEmailCommandHandler struct {
	authService services.AuthService
}

// NewResendVerificationEmailCommandHandler creates a new ResendVerificationEmailCommandHandler
func NewResendVerificationEmailCommandHandler(authService services.AuthService) *ResendVerificationEmailCommandHandler {
	return &ResendVerificationEmailCommandHandler{
		authService: authService,
	}
}

// Handle executes the ResendVerificationEmailCommand
func (h *ResendVerificationEmailCommandHandler) Handle(ctx context.Context, cmd ResendVerificationEmailCommand) (*dto.StatusResponse, error) {
	// Resend verification email
	err := h.authService.ResendVerificationEmail(ctx, cmd.Email)
	if err != nil {
		return nil, err
	}

	// Return success response
	return &dto.StatusResponse{
		Status:  "success",
		Message: "Verification email has been sent",
	}, nil
}
