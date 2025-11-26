package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
)

type LogoutAllDevicesCommand struct {
	UserID string `validate:"required"`
}

type LogoutAllDevicesCommandHandler struct {
	tokenService contracts.TokenManagementService
}

func NewLogoutAllDevicesCommandHandler(tokenService contracts.TokenManagementService) *LogoutAllDevicesCommandHandler {
	return &LogoutAllDevicesCommandHandler{
		tokenService: tokenService,
	}
}

func (h *LogoutAllDevicesCommandHandler) Handle(ctx context.Context, cmd LogoutAllDevicesCommand) (*dto.StatusResponse, error) {
	err := h.tokenService.LogoutAll(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	return &dto.StatusResponse{
		Status:  "success",
		Message: "Logged out from all devices successfully",
	}, nil
}
