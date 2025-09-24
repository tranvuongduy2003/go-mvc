package repositories

import (
	"context"
	"time"

	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
)

// RolePermission represents a role-permission assignment
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

// RolePermissionRepository defines the interface for role-permission relationship data access
type RolePermissionRepository interface {
	// GrantPermissionToRole grants a permission to a role
	GrantPermissionToRole(ctx context.Context, roleID, permissionID string, grantedBy *string) error

	// RevokePermissionFromRole revokes a permission from a role
	RevokePermissionFromRole(ctx context.Context, roleID, permissionID string) error

	// GetRolePermission retrieves a specific role-permission assignment
	GetRolePermission(ctx context.Context, roleID, permissionID string) (*RolePermission, error)

	// GetRolePermissionByID retrieves a role-permission assignment by ID
	GetRolePermissionByID(ctx context.Context, id string) (*RolePermission, error)

	// GetRolePermissions retrieves all permission assignments for a role
	GetRolePermissions(ctx context.Context, roleID string) ([]*RolePermission, error)

	// GetActiveRolePermissions retrieves all active permission assignments for a role
	GetActiveRolePermissions(ctx context.Context, roleID string) ([]*RolePermission, error)

	// GetPermissionRoles retrieves all role assignments for a permission
	GetPermissionRoles(ctx context.Context, permissionID string) ([]*RolePermission, error)

	// GetActivePermissionRoles retrieves all active role assignments for a permission
	GetActivePermissionRoles(ctx context.Context, permissionID string) ([]*RolePermission, error)

	// List retrieves a paginated list of role-permission assignments
	List(ctx context.Context, params ListRolePermissionsParams) ([]*RolePermission, *pagination.Pagination, error)

	// UpdateRolePermission updates a role-permission assignment
	UpdateRolePermission(ctx context.Context, rolePermission *RolePermission) error

	// ActivateRolePermission activates a role-permission assignment
	ActivateRolePermission(ctx context.Context, id string) error

	// DeactivateRolePermission deactivates a role-permission assignment
	DeactivateRolePermission(ctx context.Context, id string) error

	// RoleHasPermission checks if a role has a specific permission
	RoleHasPermission(ctx context.Context, roleID, permissionID string) (bool, error)

	// RoleHasPermissionByName checks if a role has a specific permission by name
	RoleHasPermissionByName(ctx context.Context, roleID, permissionName string) (bool, error)

	// RoleHasResourceAction checks if a role has permission for specific resource and action
	RoleHasResourceAction(ctx context.Context, roleID, resource, action string) (bool, error)

	// CountPermissionsByRole counts permissions assigned to a specific role
	CountPermissionsByRole(ctx context.Context, roleID string) (int64, error)

	// CountRolesByPermission counts roles assigned to a specific permission
	CountRolesByPermission(ctx context.Context, permissionID string) (int64, error)

	// Exists checks if a role-permission assignment exists
	Exists(ctx context.Context, roleID, permissionID string) (bool, error)

	// BulkGrantPermissions grants multiple permissions to a role in a single transaction
	BulkGrantPermissions(ctx context.Context, roleID string, permissionIDs []string, grantedBy *string) error

	// BulkRevokePermissions revokes multiple permissions from a role in a single transaction
	BulkRevokePermissions(ctx context.Context, roleID string, permissionIDs []string) error

	// SyncRolePermissions synchronizes role permissions (grants missing, revokes extra)
	SyncRolePermissions(ctx context.Context, roleID string, permissionIDs []string, grantedBy *string) error

	// GetRolePermissionsByResource gets all role-permission assignments for a specific resource
	GetRolePermissionsByResource(ctx context.Context, roleID, resource string) ([]*RolePermission, error)

	// GetRolePermissionsByAction gets all role-permission assignments for a specific action
	GetRolePermissionsByAction(ctx context.Context, roleID, action string) ([]*RolePermission, error)
}

// ListRolePermissionsParams represents parameters for listing role-permission assignments
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
