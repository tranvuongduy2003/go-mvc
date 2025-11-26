package auth

import (
	"context"

	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
)

// RoleRepository defines the interface for role data access
type RoleRepository interface {
	// Create saves a new role to the repository
	Create(ctx context.Context, r *Role) error

	// GetByID retrieves a role by its ID
	GetByID(ctx context.Context, id string) (*Role, error)

	// GetByName retrieves a role by its name
	GetByName(ctx context.Context, name string) (*Role, error)

	// Update saves changes to an existing role
	Update(ctx context.Context, r *Role) error

	// Delete removes a role from the repository
	Delete(ctx context.Context, id string) error

	// List retrieves a paginated list of roles
	List(ctx context.Context, params ListRolesParams) ([]*Role, *pagination.Pagination, error)

	// GetActiveRoles retrieves all active roles
	GetActiveRoles(ctx context.Context) ([]*Role, error)

	// Exists checks if a role exists by ID
	Exists(ctx context.Context, id string) (bool, error)

	// ExistsByName checks if a role exists by name
	ExistsByName(ctx context.Context, name string) (bool, error)

	// Count returns the total number of roles
	Count(ctx context.Context) (int64, error)

	// Activate activates a role
	Activate(ctx context.Context, id string) error

	// Deactivate deactivates a role
	Deactivate(ctx context.Context, id string) error

	// GetRolesByUserID retrieves all roles assigned to a user
	GetRolesByUserID(ctx context.Context, userID string) ([]*Role, error)

	// GetActiveRolesByUserID retrieves all active roles assigned to a user
	GetActiveRolesByUserID(ctx context.Context, userID string) ([]*Role, error)
}

// ListRolesParams represents parameters for listing roles
type ListRolesParams struct {
	Page     int
	Limit    int
	Search   string
	SortBy   string
	SortDir  string
	IsActive *bool
}
