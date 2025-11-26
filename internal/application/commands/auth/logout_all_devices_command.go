package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
)

// LogoutAllDevicesCommand represents the logout from all devices command
type LogoutAllDevicesCommand struct {
	UserID string `validate:"required"`
}

// LogoutAllDevicesCommandHandler handles the LogoutAllDevicesCommand
type LogoutAllDevicesCommandHandler struct {
	tokenService contracts.TokenManagementService
}

// NewLogoutAllDevicesCommandHandler creates a new LogoutAllDevicesCommandHandler
func NewLogoutAllDevicesCommandHandler(tokenService contracts.TokenManagementService) *LogoutAllDevicesCommandHandler {
	return &LogoutAllDevicesCommandHandler{
		tokenService: tokenService,
	}
}

// Handle executes the LogoutAllDevicesCommand
func (h *LogoutAllDevicesCommandHandler) Handle(ctx context.Context, cmd LogoutAllDevicesCommand) (*dto.StatusResponse, error) {
	// Logout user from all devices
	err := h.tokenService.LogoutAll(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	// Return success response
	return &dto.StatusResponse{
		Status:  "success",
		Message: "Logged out from all devices successfully",
	}, nil
}
