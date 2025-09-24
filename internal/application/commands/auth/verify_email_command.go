package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/services"
)

// VerifyEmailCommand represents the email verification command
type VerifyEmailCommand struct {
	Token string `validate:"required"`
}

// VerifyEmailCommandHandler handles the VerifyEmailCommand
type VerifyEmailCommandHandler struct {
	authService services.AuthService
}

// NewVerifyEmailCommandHandler creates a new VerifyEmailCommandHandler
func NewVerifyEmailCommandHandler(authService services.AuthService) *VerifyEmailCommandHandler {
	return &VerifyEmailCommandHandler{
		authService: authService,
	}
}

// Handle executes the VerifyEmailCommand
func (h *VerifyEmailCommandHandler) Handle(ctx context.Context, cmd VerifyEmailCommand) (*dto.StatusResponse, error) {
	// Verify email
	err := h.authService.VerifyEmail(ctx, cmd.Token)
	if err != nil {
		return nil, err
	}

	// Return success response
	return &dto.StatusResponse{
		Status:  "success",
		Message: "Email verified successfully",
	}, nil
}
