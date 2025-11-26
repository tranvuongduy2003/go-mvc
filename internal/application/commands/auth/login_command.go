package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
)

type LoginCommand struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type LoginCommandHandler struct {
	authService contracts.AuthService
}

func NewLoginCommandHandler(authService contracts.AuthService) *LoginCommandHandler {
	return &LoginCommandHandler{
		authService: authService,
	}
}

func (h *LoginCommandHandler) Handle(ctx context.Context, cmd LoginCommand) (*dto.LoginResponse, error) {
	credentials := &contracts.LoginCredentials{
		Email:    cmd.Email,
		Password: cmd.Password,
	}

	authenticatedUser, err := h.authService.Login(ctx, credentials)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		User: dto.ToAuthUserDTO(authenticatedUser.User),
		Tokens: dto.TokensDTO{
			AccessToken:           authenticatedUser.Tokens.AccessToken,
			RefreshToken:          authenticatedUser.Tokens.RefreshToken,
			AccessTokenExpiresAt:  authenticatedUser.Tokens.AccessTokenExpiresAt,
			RefreshTokenExpiresAt: authenticatedUser.Tokens.RefreshTokenExpiresAt,
			TokenType:             authenticatedUser.Tokens.TokenType,
		},
	}, nil
}
