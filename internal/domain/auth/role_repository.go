package auth

import (
	"context"

	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
)

type RoleRepository interface {
	Create(ctx context.Context, r *Role) error

	GetByID(ctx context.Context, id string) (*Role, error)

	GetByName(ctx context.Context, name string) (*Role, error)

	Update(ctx context.Context, r *Role) error

	Delete(ctx context.Context, id string) error

	List(ctx context.Context, params ListRolesParams) ([]*Role, *pagination.Pagination, error)

	GetActiveRoles(ctx context.Context) ([]*Role, error)

	Exists(ctx context.Context, id string) (bool, error)

	ExistsByName(ctx context.Context, name string) (bool, error)

	Count(ctx context.Context) (int64, error)

	Activate(ctx context.Context, id string) error

	Deactivate(ctx context.Context, id string) error

	GetRolesByUserID(ctx context.Context, userID string) ([]*Role, error)

	GetActiveRolesByUserID(ctx context.Context, userID string) ([]*Role, error)
}

type ListRolesParams struct {
	Page     int
	Limit    int
	Search   string
	SortBy   string
	SortDir  string
	IsActive *bool
}
