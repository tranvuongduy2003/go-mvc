package user

import (
	"context"

	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create saves a new user to the repository
	Create(ctx context.Context, user *User) error

	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id string) (*User, error)

	// GetByEmail retrieves a user by their email address
	GetByEmail(ctx context.Context, email string) (*User, error)

	// Update saves changes to an existing user
	Update(ctx context.Context, user *User) error

	// Delete removes a user from the repository
	Delete(ctx context.Context, id string) error

	// List retrieves a paginated list of users
	List(ctx context.Context, params ListUsersParams) ([]*User, *pagination.Pagination, error)

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
