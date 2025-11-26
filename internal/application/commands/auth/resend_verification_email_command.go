package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
)

type ResendVerificationEmailCommand struct {
	Email string `validate:"required,email"`
}

type ResendVerificationEmailCommandHandler struct {
	emailVerificationService contracts.EmailVerificationService
}

func NewResendVerificationEmailCommandHandler(emailVerificationService contracts.EmailVerificationService) *ResendVerificationEmailCommandHandler {
	return &ResendVerificationEmailCommandHandler{
		emailVerificationService: emailVerificationService,
	}
}

func (h *ResendVerificationEmailCommandHandler) Handle(ctx context.Context, cmd ResendVerificationEmailCommand) (*dto.StatusResponse, error) {
	err := h.emailVerificationService.ResendVerificationEmail(ctx, cmd.Email)
	if err != nil {
		return nil, err
	}

	return &dto.StatusResponse{
		Status:  "success",
		Message: "Verification email has been sent",
	}, nil
}
