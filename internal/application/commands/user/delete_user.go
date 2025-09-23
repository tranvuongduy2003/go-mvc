package commands

import (
	"context"
	"errors"

	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/repositories"
)

// DeleteUserCommand represents the command to delete a user
type DeleteUserCommand struct {
	ID string `json:"id" validate:"required"`
}

// DeleteUserCommandHandler handles the DeleteUserCommand
type DeleteUserCommandHandler struct {
	userRepo repositories.UserRepository
}

// NewDeleteUserCommandHandler creates a new DeleteUserCommandHandler
func NewDeleteUserCommandHandler(userRepo repositories.UserRepository) *DeleteUserCommandHandler {
	return &DeleteUserCommandHandler{
		userRepo: userRepo,
	}
}

// Handle executes the DeleteUserCommand
func (h *DeleteUserCommandHandler) Handle(ctx context.Context, cmd DeleteUserCommand) error {
	// Check if user exists
	existingUser, err := h.userRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	// Delete user (soft delete)
	if err := h.userRepo.Delete(ctx, cmd.ID); err != nil {
		return err
	}

	return nil
}
