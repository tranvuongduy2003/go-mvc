package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/shared/valueobject"
)

// Repository defines the interface for user data access
type Repository interface {
	// Basic CRUD operations
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email valueobject.Email) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	List(ctx context.Context, filters ListFilters) ([]*User, error)
	Count(ctx context.Context, filters ListFilters) (int64, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	ExistsByEmail(ctx context.Context, email valueobject.Email) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// Advanced operations
	GetActiveUsers(ctx context.Context, limit, offset int) ([]*User, error)
	GetUsersByRole(ctx context.Context, role Role, limit, offset int) ([]*User, error)
	GetUsersCreatedAfter(ctx context.Context, timestamp int64) ([]*User, error)
	GetUsersWithProfile(ctx context.Context, limit, offset int) ([]*User, error)

	// Bulk operations
	CreateBatch(ctx context.Context, users []*User) error
	UpdateBatch(ctx context.Context, users []*User) error
	DeleteBatch(ctx context.Context, ids []uuid.UUID) error

	// Search operations
	SearchByName(ctx context.Context, query string, limit, offset int) ([]*User, error)
	SearchByEmail(ctx context.Context, query string, limit, offset int) ([]*User, error)

	// Profile operations
	CreateProfile(ctx context.Context, profile *Profile) error
	GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*Profile, error)
	UpdateProfile(ctx context.Context, profile *Profile) error
	DeleteProfile(ctx context.Context, userID uuid.UUID) error
}

// ListFilters represents filters for listing users
type ListFilters struct {
	// Pagination
	Limit  int
	Offset int

	// Sorting
	SortBy    string // field to sort by
	SortOrder string // "asc" or "desc"

	// Filters
	Role      *Role
	IsActive  *bool
	IsDeleted *bool
	Email     *string
	Username  *string
	FirstName *string
	LastName  *string

	// Date range filters
	CreatedAfter  *int64
	CreatedBefore *int64
	UpdatedAfter  *int64
	UpdatedBefore *int64

	// Search
	Search string // general search across multiple fields
}
