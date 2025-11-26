package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
)

type ResetPasswordCommand struct {
	Email string `validate:"required,email"`
}

type ResetPasswordCommandHandler struct {
	passwordService contracts.PasswordManagementService
}

func NewResetPasswordCommandHandler(passwordService contracts.PasswordManagementService) *ResetPasswordCommandHandler {
	return &ResetPasswordCommandHandler{
		passwordService: passwordService,
	}
}

func (h *ResetPasswordCommandHandler) Handle(ctx context.Context, cmd ResetPasswordCommand) (*dto.StatusResponse, error) {
	err := h.passwordService.ResetPassword(ctx, cmd.Email)
	if err != nil {
		return nil, err
	}

	return &dto.StatusResponse{
		Status:  "success",
		Message: "Password reset instructions have been sent to your email",
	}, nil
}
