package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
)

type ConfirmPasswordResetCommand struct {
	Token       string `validate:"required"`
	NewPassword string `validate:"required,min=8"`
}

type ConfirmPasswordResetCommandHandler struct {
	passwordService contracts.PasswordManagementService
}

func NewConfirmPasswordResetCommandHandler(passwordService contracts.PasswordManagementService) *ConfirmPasswordResetCommandHandler {
	return &ConfirmPasswordResetCommandHandler{
		passwordService: passwordService,
	}
}

func (h *ConfirmPasswordResetCommandHandler) Handle(ctx context.Context, cmd ConfirmPasswordResetCommand) (*dto.StatusResponse, error) {
	err := h.passwordService.ConfirmPasswordReset(ctx, cmd.Token, cmd.NewPassword)
	if err != nil {
		return nil, err
	}

	return &dto.StatusResponse{
		Status:  "success",
		Message: "Password has been reset successfully",
	}, nil
}
