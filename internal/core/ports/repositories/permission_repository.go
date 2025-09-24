package repositories

import (
	"context"

	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/permission"
	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
)

// PermissionRepository defines the interface for permission data access
type PermissionRepository interface {
	// Create saves a new permission to the repository
	Create(ctx context.Context, p *permission.Permission) error

	// GetByID retrieves a permission by its ID
	GetByID(ctx context.Context, id string) (*permission.Permission, error)

	// GetByName retrieves a permission by its name
	GetByName(ctx context.Context, name string) (*permission.Permission, error)

	// GetByResourceAndAction retrieves a permission by resource and action
	GetByResourceAndAction(ctx context.Context, resource, action string) (*permission.Permission, error)

	// Update saves changes to an existing permission
	Update(ctx context.Context, p *permission.Permission) error

	// Delete removes a permission from the repository
	Delete(ctx context.Context, id string) error

	// List retrieves a paginated list of permissions
	List(ctx context.Context, params ListPermissionsParams) ([]*permission.Permission, *pagination.Pagination, error)

	// GetActivePermissions retrieves all active permissions
	GetActivePermissions(ctx context.Context) ([]*permission.Permission, error)

	// GetPermissionsByResource retrieves all permissions for a specific resource
	GetPermissionsByResource(ctx context.Context, resource string) ([]*permission.Permission, error)

	// GetPermissionsByAction retrieves all permissions for a specific action
	GetPermissionsByAction(ctx context.Context, action string) ([]*permission.Permission, error)

	// Exists checks if a permission exists by ID
	Exists(ctx context.Context, id string) (bool, error)

	// ExistsByName checks if a permission exists by name
	ExistsByName(ctx context.Context, name string) (bool, error)

	// ExistsByResourceAndAction checks if a permission exists by resource and action
	ExistsByResourceAndAction(ctx context.Context, resource, action string) (bool, error)

	// Count returns the total number of permissions
	Count(ctx context.Context) (int64, error)

	// Activate activates a permission
	Activate(ctx context.Context, id string) error

	// Deactivate deactivates a permission
	Deactivate(ctx context.Context, id string) error

	// GetPermissionsByUserID retrieves all permissions for a user (through roles)
	GetPermissionsByUserID(ctx context.Context, userID string) ([]*permission.Permission, error)

	// GetActivePermissionsByUserID retrieves all active permissions for a user (through roles)
	GetActivePermissionsByUserID(ctx context.Context, userID string) ([]*permission.Permission, error)

	// GetPermissionsByRoleID retrieves all permissions assigned to a role
	GetPermissionsByRoleID(ctx context.Context, roleID string) ([]*permission.Permission, error)

	// GetActivePermissionsByRoleID retrieves all active permissions assigned to a role
	GetActivePermissionsByRoleID(ctx context.Context, roleID string) ([]*permission.Permission, error)

	// UserHasPermission checks if a user has a specific permission
	UserHasPermission(ctx context.Context, userID, resource, action string) (bool, error)

	// UserHasPermissionByName checks if a user has a specific permission by name
	UserHasPermissionByName(ctx context.Context, userID, permissionName string) (bool, error)
}

// ListPermissionsParams represents parameters for listing permissions
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
