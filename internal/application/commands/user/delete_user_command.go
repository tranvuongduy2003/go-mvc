package commands

import (
	"context"
	"errors"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/user"
)

type DeleteUserCommand struct {
	ID string `json:"id" validate:"required"`
}

type DeleteUserCommandHandler struct {
	userRepo user.UserRepository
}

func NewDeleteUserCommandHandler(userRepo user.UserRepository) *DeleteUserCommandHandler {
	return &DeleteUserCommandHandler{
		userRepo: userRepo,
	}
}

func (h *DeleteUserCommandHandler) Handle(ctx context.Context, cmd DeleteUserCommand) error {
	existingUser, err := h.userRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	if err := h.userRepo.Delete(ctx, cmd.ID); err != nil {
		return err
	}

	return nil
}
