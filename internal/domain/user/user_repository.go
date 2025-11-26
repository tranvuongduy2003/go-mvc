package user

import (
	"context"

	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error

	GetByID(ctx context.Context, id string) (*User, error)

	GetByEmail(ctx context.Context, email string) (*User, error)

	Update(ctx context.Context, user *User) error

	Delete(ctx context.Context, id string) error

	List(ctx context.Context, params ListUsersParams) ([]*User, *pagination.Pagination, error)

	Exists(ctx context.Context, id string) (bool, error)

	ExistsByEmail(ctx context.Context, email string) (bool, error)

	Count(ctx context.Context) (int64, error)
}

type ListUsersParams struct {
	Page     int
	Limit    int
	Search   string
	SortBy   string
	SortDir  string
	IsActive *bool
}
