package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/services"
)

// LogoutAllDevicesCommand represents the logout from all devices command
type LogoutAllDevicesCommand struct {
	UserID string `validate:"required"`
}

// LogoutAllDevicesCommandHandler handles the LogoutAllDevicesCommand
type LogoutAllDevicesCommandHandler struct {
	authService services.AuthService
}

// NewLogoutAllDevicesCommandHandler creates a new LogoutAllDevicesCommandHandler
func NewLogoutAllDevicesCommandHandler(authService services.AuthService) *LogoutAllDevicesCommandHandler {
	return &LogoutAllDevicesCommandHandler{
		authService: authService,
	}
}

// Handle executes the LogoutAllDevicesCommand
func (h *LogoutAllDevicesCommandHandler) Handle(ctx context.Context, cmd LogoutAllDevicesCommand) (*dto.StatusResponse, error) {
	// Logout user from all devices
	err := h.authService.LogoutAll(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	// Return success response
	return &dto.StatusResponse{
		Status:  "success",
		Message: "Logged out from all devices successfully",
	}, nil
}
