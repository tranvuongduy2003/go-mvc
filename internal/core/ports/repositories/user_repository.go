package repositories

import (
	"context"

	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/user"
	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create saves a new user to the repository
	Create(ctx context.Context, user *user.User) error

	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id string) (*user.User, error)

	// GetByEmail retrieves a user by their email address
	GetByEmail(ctx context.Context, email string) (*user.User, error)

	// Update saves changes to an existing user
	Update(ctx context.Context, user *user.User) error

	// Delete removes a user from the repository
	Delete(ctx context.Context, id string) error

	// List retrieves a paginated list of users
	List(ctx context.Context, params ListUsersParams) ([]*user.User, *pagination.Pagination, error)

	// Exists checks if a user exists by ID
	Exists(ctx context.Context, id string) (bool, error)

	// ExistsByEmail checks if a user exists by email
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)
}

// ListUsersParams represents parameters for listing users
type ListUsersParams struct {
	Page     int
	Limit    int
	Search   string
	SortBy   string
	SortDir  string
	IsActive *bool
}
