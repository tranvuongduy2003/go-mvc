package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
)

type RegisterCommand struct {
	Email    string `validate:"required,email"`
	Name     string `validate:"required,min=2,max=100"`
	Phone    string `validate:"omitempty"`
	Password string `validate:"required,min=8"`
}

type RegisterCommandHandler struct {
	authService contracts.AuthService
}

func NewRegisterCommandHandler(authService contracts.AuthService) *RegisterCommandHandler {
	return &RegisterCommandHandler{
		authService: authService,
	}
}

func (h *RegisterCommandHandler) Handle(ctx context.Context, cmd RegisterCommand) (*dto.RegisterResponse, error) {
	registerReq := &contracts.RegisterRequest{
		Email:    cmd.Email,
		Name:     cmd.Name,
		Phone:    cmd.Phone,
		Password: cmd.Password,
	}

	authenticatedUser, err := h.authService.Register(ctx, registerReq)
	if err != nil {
		return nil, err
	}

	return &dto.RegisterResponse{
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
