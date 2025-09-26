package queries

import (
	"context"

	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/user"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/repositories"
	apperrors "github.com/tranvuongduy2003/go-mvc/pkg/errors"
)

// GetUserByIDQuery represents the query to get a user by ID
type GetUserByIDQuery struct {
	ID string `json:"id" validate:"required"`
}

// GetUserByIDQueryHandler handles the GetUserByIDQuery
type GetUserByIDQueryHandler struct {
	userRepo repositories.UserRepository
}

// NewGetUserByIDQueryHandler creates a new GetUserByIDQueryHandler
func NewGetUserByIDQueryHandler(userRepo repositories.UserRepository) *GetUserByIDQueryHandler {
	return &GetUserByIDQueryHandler{
		userRepo: userRepo,
	}
}

// Handle executes the GetUserByIDQuery
func (h *GetUserByIDQueryHandler) Handle(ctx context.Context, query GetUserByIDQuery) (*user.User, error) {
	user, err := h.userRepo.GetByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.NewNotFoundError("user not found")
	}
	return user, nil
}
