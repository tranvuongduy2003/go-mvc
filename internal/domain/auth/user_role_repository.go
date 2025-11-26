package auth

import (
	"context"
	"time"

	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
)

// UserRoleRepository defines the interface for user-role relationship data access
type UserRoleRepository interface {
	// AssignRoleToUser assigns a role to a user
	AssignRoleToUser(ctx context.Context, userID, roleID string, assignedBy *string, expiresAt *time.Time) error

	// RevokeRoleFromUser revokes a role from a user
	RevokeRoleFromUser(ctx context.Context, userID, roleID string) error

	// GetUserRole retrieves a specific user-role assignment
	GetUserRole(ctx context.Context, userID, roleID string) (*UserRole, error)

	// GetUserRoleByID retrieves a user-role assignment by ID
	GetUserRoleByID(ctx context.Context, id string) (*UserRole, error)

	// GetUserRoles retrieves all role assignments for a user
	GetUserRoles(ctx context.Context, userID string) ([]*UserRole, error)

	// GetActiveUserRoles retrieves all active role assignments for a user
	GetActiveUserRoles(ctx context.Context, userID string) ([]*UserRole, error)

	// GetRoleUsers retrieves all user assignments for a role
	GetRoleUsers(ctx context.Context, roleID string) ([]*UserRole, error)

	// GetActiveRoleUsers retrieves all active user assignments for a role
	GetActiveRoleUsers(ctx context.Context, roleID string) ([]*UserRole, error)

	// List retrieves a paginated list of user-role assignments
	List(ctx context.Context, params ListUserRolesParams) ([]*UserRole, *pagination.Pagination, error)

	// UpdateUserRole updates a user-role assignment
	UpdateUserRole(ctx context.Context, userRole *UserRole) error

	// ActivateUserRole activates a user-role assignment
	ActivateUserRole(ctx context.Context, id string) error

	// DeactivateUserRole deactivates a user-role assignment
	DeactivateUserRole(ctx context.Context, id string) error

	// SetExpiration sets or updates expiration for a user-role assignment
	SetExpiration(ctx context.Context, userID, roleID string, expiresAt *time.Time) error

	// IsUserRoleExpired checks if a user-role assignment is expired
	IsUserRoleExpired(ctx context.Context, userID, roleID string) (bool, error)

	// GetExpiredUserRoles retrieves all expired user-role assignments
	GetExpiredUserRoles(ctx context.Context) ([]*UserRole, error)

	// CleanupExpiredRoles deactivates all expired user-role assignments
	CleanupExpiredRoles(ctx context.Context) (int64, error)

	// UserHasRole checks if a user currently has a specific role (active and not expired)
	UserHasRole(ctx context.Context, userID, roleID string) (bool, error)

	// UserHasRoleName checks if a user currently has a specific role by name
	UserHasRoleName(ctx context.Context, userID, roleName string) (bool, error)

	// CountUsersByRole counts users assigned to a specific role
	CountUsersByRole(ctx context.Context, roleID string) (int64, error)

	// CountRolesByUser counts roles assigned to a specific user
	CountRolesByUser(ctx context.Context, userID string) (int64, error)

	// Exists checks if a user-role assignment exists
	Exists(ctx context.Context, userID, roleID string) (bool, error)
}

// ListUserRolesParams represents parameters for listing user-role assignments
type ListUserRolesParams struct {
	Page       int
	Limit      int
	UserID     string
	RoleID     string
	AssignedBy string
	IsActive   *bool
	IsExpired  *bool
	SortBy     string
	SortDir    string
}
