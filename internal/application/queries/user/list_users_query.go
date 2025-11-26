package user

import (
	"context"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/user"
	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
)

type ListUsersQuery struct {
	Page     int    `json:"page" validate:"min=1"`
	Limit    int    `json:"limit" validate:"min=1,max=100"`
	Search   string `json:"search"`
	SortBy   string `json:"sort_by"`
	SortDir  string `json:"sort_dir"`
	IsActive *bool  `json:"is_active"`
}

type ListUsersQueryHandler struct {
	userRepo user.UserRepository
}

func NewListUsersQueryHandler(userRepo user.UserRepository) *ListUsersQueryHandler {
	return &ListUsersQueryHandler{
		userRepo: userRepo,
	}
}

func (h *ListUsersQueryHandler) Handle(ctx context.Context, query ListUsersQuery) ([]*user.User, *pagination.Pagination, error) {
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

	params := user.ListUsersParams{
		Page:     query.Page,
		Limit:    query.Limit,
		Search:   query.Search,
		SortBy:   query.SortBy,
		SortDir:  query.SortDir,
		IsActive: query.IsActive,
	}

	return h.userRepo.List(ctx, params)
}
