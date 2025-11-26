package auth

import (
	"context"

	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
)

type PermissionRepository interface {
	Create(ctx context.Context, p *Permission) error

	GetByID(ctx context.Context, id string) (*Permission, error)

	GetByName(ctx context.Context, name string) (*Permission, error)

	GetByResourceAndAction(ctx context.Context, resource, action string) (*Permission, error)

	Update(ctx context.Context, p *Permission) error

	Delete(ctx context.Context, id string) error

	List(ctx context.Context, params ListPermissionsParams) ([]*Permission, *pagination.Pagination, error)

	GetActivePermissions(ctx context.Context) ([]*Permission, error)

	GetPermissionsByResource(ctx context.Context, resource string) ([]*Permission, error)

	GetPermissionsByAction(ctx context.Context, action string) ([]*Permission, error)

	Exists(ctx context.Context, id string) (bool, error)

	ExistsByName(ctx context.Context, name string) (bool, error)

	ExistsByResourceAndAction(ctx context.Context, resource, action string) (bool, error)

	Count(ctx context.Context) (int64, error)

	Activate(ctx context.Context, id string) error

	Deactivate(ctx context.Context, id string) error

	GetPermissionsByUserID(ctx context.Context, userID string) ([]*Permission, error)

	GetActivePermissionsByUserID(ctx context.Context, userID string) ([]*Permission, error)

	GetPermissionsByRoleID(ctx context.Context, roleID string) ([]*Permission, error)

	GetActivePermissionsByRoleID(ctx context.Context, roleID string) ([]*Permission, error)

	UserHasPermission(ctx context.Context, userID, resource, action string) (bool, error)

	UserHasPermissionByName(ctx context.Context, userID, permissionName string) (bool, error)
}

type ListPermissionsParams struct {
	Page     int
	Limit    int
	Search   string
	SortBy   string
	SortDir  string
	Resource string
	Action   string
	IsActive *bool
}
