package auth

import (
	"context"
	"time"

	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
)

type RolePermission struct {
	ID           string
	RoleID       string
	PermissionID string
	GrantedBy    *string
	GrantedAt    time.Time
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Version      int64
}

type RolePermissionRepository interface {
	GrantPermissionToRole(ctx context.Context, roleID, permissionID string, grantedBy *string) error

	RevokePermissionFromRole(ctx context.Context, roleID, permissionID string) error

	GetRolePermission(ctx context.Context, roleID, permissionID string) (*RolePermission, error)

	GetRolePermissionByID(ctx context.Context, id string) (*RolePermission, error)

	GetRolePermissions(ctx context.Context, roleID string) ([]*RolePermission, error)

	GetActiveRolePermissions(ctx context.Context, roleID string) ([]*RolePermission, error)

	GetPermissionRoles(ctx context.Context, permissionID string) ([]*RolePermission, error)

	GetActivePermissionRoles(ctx context.Context, permissionID string) ([]*RolePermission, error)

	List(ctx context.Context, params ListRolePermissionsParams) ([]*RolePermission, *pagination.Pagination, error)

	UpdateRolePermission(ctx context.Context, rolePermission *RolePermission) error

	ActivateRolePermission(ctx context.Context, id string) error

	DeactivateRolePermission(ctx context.Context, id string) error

	RoleHasPermission(ctx context.Context, roleID, permissionID string) (bool, error)

	RoleHasPermissionByName(ctx context.Context, roleID, permissionName string) (bool, error)

	RoleHasResourceAction(ctx context.Context, roleID, resource, action string) (bool, error)

	CountPermissionsByRole(ctx context.Context, roleID string) (int64, error)

	CountRolesByPermission(ctx context.Context, permissionID string) (int64, error)

	Exists(ctx context.Context, roleID, permissionID string) (bool, error)

	BulkGrantPermissions(ctx context.Context, roleID string, permissionIDs []string, grantedBy *string) error

	BulkRevokePermissions(ctx context.Context, roleID string, permissionIDs []string) error

	SyncRolePermissions(ctx context.Context, roleID string, permissionIDs []string, grantedBy *string) error

	GetRolePermissionsByResource(ctx context.Context, roleID, resource string) ([]*RolePermission, error)

	GetRolePermissionsByAction(ctx context.Context, roleID, action string) ([]*RolePermission, error)
}

type ListRolePermissionsParams struct {
	Page         int
	Limit        int
	RoleID       string
	PermissionID string
	Resource     string
	Action       string
	GrantedBy    string
	IsActive     *bool
	SortBy       string
	SortDir      string
}
