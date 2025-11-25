package queries

import (
	"context"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/ports/repositories"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/user"
	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
)

// ListUsersQuery represents the query to list users with pagination
type ListUsersQuery struct {
	Page     int    `json:"page" validate:"min=1"`
	Limit    int    `json:"limit" validate:"min=1,max=100"`
	Search   string `json:"search"`
	SortBy   string `json:"sort_by"`
	SortDir  string `json:"sort_dir"`
	IsActive *bool  `json:"is_active"`
}

// ListUsersQueryHandler handles the ListUsersQuery
type ListUsersQueryHandler struct {
	userRepo repositories.UserRepository
}

// NewListUsersQueryHandler creates a new ListUsersQueryHandler
func NewListUsersQueryHandler(userRepo repositories.UserRepository) *ListUsersQueryHandler {
	return &ListUsersQueryHandler{
		userRepo: userRepo,
	}
}

// Handle executes the ListUsersQuery
func (h *ListUsersQueryHandler) Handle(ctx context.Context, query ListUsersQuery) ([]*user.User, *pagination.Pagination, error) {
	// Set default values
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 10
	}
	if query.SortBy == "" {
		query.SortBy = "created_at"
	}
	if query.SortDir == "" {
		query.SortDir = "desc"
	}

	params := repositories.ListUsersParams{
		Page:     query.Page,
		Limit:    query.Limit,
		Search:   query.Search,
		SortBy:   query.SortBy,
		SortDir:  query.SortDir,
		IsActive: query.IsActive,
	}

	return h.userRepo.List(ctx, params)
}
