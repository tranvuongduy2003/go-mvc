package commands

import (
	"context"
	"errors"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/user"
	apperrors "github.com/tranvuongduy2003/go-mvc/pkg/errors"
)

type CreateUserCommand struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Phone    string `json:"phone" validate:"omitempty"`
	Password string `json:"password" validate:"required,min=8"`
}

func (c CreateUserCommand) Validate() error {
	if c.Email == "" {
		return errors.New("email is required")
	}
	if c.Name == "" {
		return errors.New("name is required")
	}
	if len(c.Name) < 2 {
		return errors.New("name must be at least 2 characters")
	}
	if c.Password == "" {
		return errors.New("password is required")
	}
	if len(c.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

type CreateUserCommandHandler struct {
	userRepo user.UserRepository
}

func NewCreateUserCommandHandler(userRepo user.UserRepository) *CreateUserCommandHandler {
	return &CreateUserCommandHandler{
		userRepo: userRepo,
	}
}

func (h *CreateUserCommandHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*user.User, error) {
	existingUser, err := h.userRepo.GetByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to check existing user", err)
	}
	if existingUser != nil {
		return nil, apperrors.NewConflictError("User with email "+cmd.Email+" already exists", nil)
	}

	newUser, err := user.NewUser(cmd.Email, cmd.Name, cmd.Phone, cmd.Password)
	if err != nil {
		return nil, apperrors.NewValidationError(err.Error(), err)
	}

	if err := h.userRepo.Create(ctx, newUser); err != nil {
		return nil, apperrors.NewInternalError("Failed to create user", err)
	}

	return newUser, nil
}
