package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
)

type ChangePasswordCommand struct {
	UserID      string `validate:"required"`
	OldPassword string `validate:"required"`
	NewPassword string `validate:"required,min=8"`
}

type ChangePasswordCommandHandler struct {
	passwordService contracts.PasswordManagementService
}

func NewChangePasswordCommandHandler(passwordService contracts.PasswordManagementService) *ChangePasswordCommandHandler {
	return &ChangePasswordCommandHandler{
		passwordService: passwordService,
	}
}

func (h *ChangePasswordCommandHandler) Handle(ctx context.Context, cmd ChangePasswordCommand) (*dto.StatusResponse, error) {
	err := h.passwordService.ChangePassword(ctx, cmd.UserID, cmd.OldPassword, cmd.NewPassword)
	if err != nil {
		return nil, err
	}

	return &dto.StatusResponse{
		Status:  "success",
		Message: "Password changed successfully",
	}, nil
}
