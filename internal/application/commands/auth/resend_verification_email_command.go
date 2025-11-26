package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
)

// ResendVerificationEmailCommand represents the resend verification email command
type ResendVerificationEmailCommand struct {
	Email string `validate:"required,email"`
}

// ResendVerificationEmailCommandHandler handles the ResendVerificationEmailCommand
type ResendVerificationEmailCommandHandler struct {
	emailVerificationService contracts.EmailVerificationService
}

// NewResendVerificationEmailCommandHandler creates a new ResendVerificationEmailCommandHandler
func NewResendVerificationEmailCommandHandler(emailVerificationService contracts.EmailVerificationService) *ResendVerificationEmailCommandHandler {
	return &ResendVerificationEmailCommandHandler{
		emailVerificationService: emailVerificationService,
	}
}

// Handle executes the ResendVerificationEmailCommand
func (h *ResendVerificationEmailCommandHandler) Handle(ctx context.Context, cmd ResendVerificationEmailCommand) (*dto.StatusResponse, error) {
	err := h.emailVerificationService.ResendVerificationEmail(ctx, cmd.Email)
	if err != nil {
		return nil, err
	}

	// Return success response
	return &dto.StatusResponse{
		Status:  "success",
		Message: "Verification email has been sent",
	}, nil
}
