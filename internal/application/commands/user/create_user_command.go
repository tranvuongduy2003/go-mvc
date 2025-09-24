package commands

import (
	"context"
	"errors"

	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/user"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/repositories"
)

// CreateUserCommand represents the command to create a new user
type CreateUserCommand struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Phone    string `json:"phone" validate:"omitempty"`
	Password string `json:"password" validate:"required,min=8"`
}

// CreateUserCommandHandler handles the CreateUserCommand
type CreateUserCommandHandler struct {
	userRepo repositories.UserRepository
}

// NewCreateUserCommandHandler creates a new CreateUserCommandHandler
func NewCreateUserCommandHandler(userRepo repositories.UserRepository) *CreateUserCommandHandler {
	return &CreateUserCommandHandler{
		userRepo: userRepo,
	}
}

// Handle executes the CreateUserCommand
func (h *CreateUserCommandHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*user.User, error) {
	// Check if user already exists
	existingUser, err := h.userRepo.GetByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create new user
	newUser, err := user.NewUser(cmd.Email, cmd.Name, cmd.Phone, cmd.Password)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := h.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
