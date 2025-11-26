package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
)

type VerifyEmailCommand struct {
	Token string `validate:"required"`
}

type VerifyEmailCommandHandler struct {
	emailVerificationService contracts.EmailVerificationService
}

func NewVerifyEmailCommandHandler(emailVerificationService contracts.EmailVerificationService) *VerifyEmailCommandHandler {
	return &VerifyEmailCommandHandler{
		emailVerificationService: emailVerificationService,
	}
}

func (h *VerifyEmailCommandHandler) Handle(ctx context.Context, cmd VerifyEmailCommand) (*dto.StatusResponse, error) {
	err := h.emailVerificationService.VerifyEmail(ctx, cmd.Token)
	if err != nil {
		return nil, err
	}

	return &dto.StatusResponse{
		Status:  "success",
		Message: "Email verified successfully",
	}, nil
}
