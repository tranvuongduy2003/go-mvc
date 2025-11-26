package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
)

// LoginCommand represents the login command
type LoginCommand struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

// LoginCommandHandler handles the LoginCommand
type LoginCommandHandler struct {
	authService contracts.AuthService
}

// NewLoginCommandHandler creates a new LoginCommandHandler
func NewLoginCommandHandler(authService contracts.AuthService) *LoginCommandHandler {
	return &LoginCommandHandler{
		authService: authService,
	}
}

// Handle executes the LoginCommand
func (h *LoginCommandHandler) Handle(ctx context.Context, cmd LoginCommand) (*dto.LoginResponse, error) {
	// Create login credentials
	credentials := &contracts.LoginCredentials{
		Email:    cmd.Email,
		Password: cmd.Password,
	}

	// Authenticate user
	authenticatedUser, err := h.authService.Login(ctx, credentials)
	if err != nil {
		return nil, err
	}

	// Convert to response DTO
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
