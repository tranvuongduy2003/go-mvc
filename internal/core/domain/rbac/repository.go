package rbac

import (
	"context"

	"github.com/google/uuid"
)

// RoleRepository defines role repository interface
type RoleRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, role *Role) error
	GetByID(ctx context.Context, id uuid.UUID) (*Role, error)
	GetByName(ctx context.Context, name string) (*Role, error)
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*Role, error)
	Count(ctx context.Context) (int64, error)

	// Role-specific operations
	GetRolesWithPermissions(ctx context.Context, roleIDs []uuid.UUID) ([]*Role, error)
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*Role, error)
	SearchByName(ctx context.Context, name string, offset, limit int) ([]*Role, error)
	GetActiveRoles(ctx context.Context) ([]*Role, error)

	// Permission management
	AddPermissionToRole(ctx context.Context, roleID, permissionID, grantedBy uuid.UUID) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID, removedBy uuid.UUID) error
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*Permission, error)
}

// PermissionRepository defines permission repository interface
type PermissionRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, permission *Permission) error
	GetByID(ctx context.Context, id uuid.UUID) (*Permission, error)
	GetByName(ctx context.Context, name string) (*Permission, error)
	GetByResourceAndAction(ctx context.Context, resource, action string) (*Permission, error)
	Update(ctx context.Context, permission *Permission) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*Permission, error)
	Count(ctx context.Context) (int64, error)

	// Permission-specific operations
	GetByResource(ctx context.Context, resource string) ([]*Permission, error)
	GetByAction(ctx context.Context, action string) ([]*Permission, error)
	SearchByName(ctx context.Context, name string, offset, limit int) ([]*Permission, error)
	GetActivePermissions(ctx context.Context) ([]*Permission, error)
	GetPermissionsByIDs(ctx context.Context, permissionIDs []uuid.UUID) ([]*Permission, error)
}

// UserRoleRepository defines user role repository interface
type UserRoleRepository interface {
	// User-Role assignment operations
	AssignRoleToUser(ctx context.Context, userRole *UserRole) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID, removedBy uuid.UUID) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*UserRole, error)
	GetActiveUserRoles(ctx context.Context, userID uuid.UUID) ([]*UserRole, error)
	GetRoleUsers(ctx context.Context, roleID uuid.UUID) ([]*UserRole, error)

	// Validation and checking
	HasRole(ctx context.Context, userID, roleID uuid.UUID) (bool, error)
	IsRoleExpired(ctx context.Context, userID, roleID uuid.UUID) (bool, error)
	GetExpiredRoles(ctx context.Context) ([]*UserRole, error)

	// Bulk operations
	AssignMultipleRolesToUser(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID, assignedBy uuid.UUID) error
	RemoveAllRolesFromUser(ctx context.Context, userID uuid.UUID, removedBy uuid.UUID) error
	GetUsersWithRole(ctx context.Context, roleID uuid.UUID) ([]uuid.UUID, error)

	// List and search
	List(ctx context.Context, offset, limit int) ([]*UserRole, error)
	Count(ctx context.Context) (int64, error)
	GetUserRoleHistory(ctx context.Context, userID uuid.UUID) ([]*UserRole, error)
}

// RolePermissionRepository defines role permission repository interface
type RolePermissionRepository interface {
	// Role-Permission assignment operations
	GrantPermissionToRole(ctx context.Context, rolePermission *RolePermission) error
	RevokePermissionFromRole(ctx context.Context, roleID, permissionID uuid.UUID, revokedBy uuid.UUID) error
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*RolePermission, error)
	GetActiveRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*RolePermission, error)
	GetPermissionRoles(ctx context.Context, permissionID uuid.UUID) ([]*RolePermission, error)

	// Validation and checking
	HasPermission(ctx context.Context, roleID, permissionID uuid.UUID) (bool, error)

	// Bulk operations
	GrantMultiplePermissionsToRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID, grantedBy uuid.UUID) error
	RevokeAllPermissionsFromRole(ctx context.Context, roleID uuid.UUID, revokedBy uuid.UUID) error
	GetRolesWithPermission(ctx context.Context, permissionID uuid.UUID) ([]uuid.UUID, error)

	// List and search
	List(ctx context.Context, offset, limit int) ([]*RolePermission, error)
	Count(ctx context.Context) (int64, error)
	GetRolePermissionHistory(ctx context.Context, roleID uuid.UUID) ([]*RolePermission, error)
}
