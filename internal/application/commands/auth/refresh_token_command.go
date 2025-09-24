package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/services"
)

// RefreshTokenCommand represents the refresh token command
type RefreshTokenCommand struct {
	RefreshToken string `validate:"required"`
}

// RefreshTokenCommandHandler handles the RefreshTokenCommand
type RefreshTokenCommandHandler struct {
	authService services.AuthService
}

// NewRefreshTokenCommandHandler creates a new RefreshTokenCommandHandler
func NewRefreshTokenCommandHandler(authService services.AuthService) *RefreshTokenCommandHandler {
	return &RefreshTokenCommandHandler{
		authService: authService,
	}
}

// Handle executes the RefreshTokenCommand
func (h *RefreshTokenCommandHandler) Handle(ctx context.Context, cmd RefreshTokenCommand) (*dto.RefreshTokenResponse, error) {
	// Refresh token
	tokens, err := h.authService.RefreshToken(ctx, cmd.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Convert to response DTO
	return &dto.RefreshTokenResponse{
		Tokens: dto.TokensDTO{
			AccessToken:           tokens.AccessToken,
			RefreshToken:          tokens.RefreshToken,
			AccessTokenExpiresAt:  tokens.AccessTokenExpiresAt,
			RefreshTokenExpiresAt: tokens.RefreshTokenExpiresAt,
			TokenType:             tokens.TokenType,
		},
	}, nil
}
