package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
)

type LogoutCommand struct {
	UserID string `validate:"required"`
}

type LogoutCommandHandler struct {
	tokenService contracts.TokenManagementService
}

func NewLogoutCommandHandler(tokenService contracts.TokenManagementService) *LogoutCommandHandler {
	return &LogoutCommandHandler{
		tokenService: tokenService,
	}
}

func (h *LogoutCommandHandler) Handle(ctx context.Context, cmd LogoutCommand) (*dto.StatusResponse, error) {
	err := h.tokenService.Logout(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	return &dto.StatusResponse{
		Status:  "success",
		Message: "Logged out successfully",
	}, nil
}
