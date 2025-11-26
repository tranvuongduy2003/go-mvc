package commands

import (
	"context"
	"errors"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/user"
)

type UpdateUserCommand struct {
	ID    string `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Phone string `json:"phone" validate:"omitempty"`
}

type UpdateUserCommandHandler struct {
	userRepo user.UserRepository
}

func NewUpdateUserCommandHandler(userRepo user.UserRepository) *UpdateUserCommandHandler {
	return &UpdateUserCommandHandler{
		userRepo: userRepo,
	}
}

func (h *UpdateUserCommandHandler) Handle(ctx context.Context, cmd UpdateUserCommand) (*user.User, error) {
	existingUser, err := h.userRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, errors.New("user not found")
	}

	if err := existingUser.UpdateProfile(cmd.Name, cmd.Phone); err != nil {
		return nil, err
	}

	if err := h.userRepo.Update(ctx, existingUser); err != nil {
		return nil, err
	}

	return existingUser, nil
}
