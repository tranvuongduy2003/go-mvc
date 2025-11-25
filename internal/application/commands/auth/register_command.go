package commands

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/ports/services"
)

// RegisterCommand represents the register command
type RegisterCommand struct {
	Email    string `validate:"required,email"`
	Name     string `validate:"required,min=2,max=100"`
	Phone    string `validate:"omitempty"`
	Password string `validate:"required,min=8"`
}

// RegisterCommandHandler handles the RegisterCommand
type RegisterCommandHandler struct {
	authService services.AuthService
}

// NewRegisterCommandHandler creates a new RegisterCommandHandler
func NewRegisterCommandHandler(authService services.AuthService) *RegisterCommandHandler {
	return &RegisterCommandHandler{
		authService: authService,
	}
}

// Handle executes the RegisterCommand
func (h *RegisterCommandHandler) Handle(ctx context.Context, cmd RegisterCommand) (*dto.RegisterResponse, error) {
	// Create register request
	registerReq := &services.RegisterRequest{
		Email:    cmd.Email,
		Name:     cmd.Name,
		Phone:    cmd.Phone,
		Password: cmd.Password,
	}

	// Register user
	authenticatedUser, err := h.authService.Register(ctx, registerReq)
	if err != nil {
		return nil, err
	}

	// Convert to response DTO
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
