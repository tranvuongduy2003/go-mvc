package commands

import (
	"context"
	"errors"

	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/user"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/repositories"
)

// UpdateUserCommand represents the command to update an existing user
type UpdateUserCommand struct {
	ID    string `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Phone string `json:"phone" validate:"omitempty"`
}

// UpdateUserCommandHandler handles the UpdateUserCommand
type UpdateUserCommandHandler struct {
	userRepo repositories.UserRepository
}

// NewUpdateUserCommandHandler creates a new UpdateUserCommandHandler
func NewUpdateUserCommandHandler(userRepo repositories.UserRepository) *UpdateUserCommandHandler {
	return &UpdateUserCommandHandler{
		userRepo: userRepo,
	}
}

// Handle executes the UpdateUserCommand
func (h *UpdateUserCommandHandler) Handle(ctx context.Context, cmd UpdateUserCommand) (*user.User, error) {
	// Get existing user
	existingUser, err := h.userRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, errors.New("user not found")
	}

	// Update user profile
	if err := existingUser.UpdateProfile(cmd.Name, cmd.Phone); err != nil {
		return nil, err
	}

	// Save to repository
	if err := h.userRepo.Update(ctx, existingUser); err != nil {
		return nil, err
	}

	return existingUser, nil
}
