package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
)

type RefreshTokenCommand struct {
	RefreshToken string `validate:"required"`
}

type RefreshTokenCommandHandler struct {
	authService contracts.AuthService
}

func NewRefreshTokenCommandHandler(authService contracts.AuthService) *RefreshTokenCommandHandler {
	return &RefreshTokenCommandHandler{
		authService: authService,
	}
}

func (h *RefreshTokenCommandHandler) Handle(ctx context.Context, cmd RefreshTokenCommand) (*dto.RefreshTokenResponse, error) {
	tokens, err := h.authService.RefreshToken(ctx, cmd.RefreshToken)
	if err != nil {
		return nil, err
	}

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
