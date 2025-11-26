package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
)

// LogoutCommand represents the logout command
type LogoutCommand struct {
	UserID string `validate:"required"`
}

// LogoutCommandHandler handles the LogoutCommand
type LogoutCommandHandler struct {
	tokenService contracts.TokenManagementService
}

// NewLogoutCommandHandler creates a new LogoutCommandHandler
func NewLogoutCommandHandler(tokenService contracts.TokenManagementService) *LogoutCommandHandler {
	return &LogoutCommandHandler{
		tokenService: tokenService,
	}
}

// Handle executes the LogoutCommand
func (h *LogoutCommandHandler) Handle(ctx context.Context, cmd LogoutCommand) (*dto.StatusResponse, error) {
	// Logout user by invalidating tokens
	err := h.tokenService.Logout(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	// Return success response
	return &dto.StatusResponse{
		Status:  "success",
		Message: "Logged out successfully",
	}, nil
}
