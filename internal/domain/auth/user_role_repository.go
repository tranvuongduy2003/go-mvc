package auth

import (
	"context"
	"time"

	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
)

type UserRoleRepository interface {
	AssignRoleToUser(ctx context.Context, userID, roleID string, assignedBy *string, expiresAt *time.Time) error

	RevokeRoleFromUser(ctx context.Context, userID, roleID string) error

	GetUserRole(ctx context.Context, userID, roleID string) (*UserRole, error)

	GetUserRoleByID(ctx context.Context, id string) (*UserRole, error)

	GetUserRoles(ctx context.Context, userID string) ([]*UserRole, error)

	GetActiveUserRoles(ctx context.Context, userID string) ([]*UserRole, error)

	GetRoleUsers(ctx context.Context, roleID string) ([]*UserRole, error)

	GetActiveRoleUsers(ctx context.Context, roleID string) ([]*UserRole, error)

	List(ctx context.Context, params ListUserRolesParams) ([]*UserRole, *pagination.Pagination, error)

	UpdateUserRole(ctx context.Context, userRole *UserRole) error

	ActivateUserRole(ctx context.Context, id string) error

	DeactivateUserRole(ctx context.Context, id string) error

	SetExpiration(ctx context.Context, userID, roleID string, expiresAt *time.Time) error

	IsUserRoleExpired(ctx context.Context, userID, roleID string) (bool, error)

	GetExpiredUserRoles(ctx context.Context) ([]*UserRole, error)

	CleanupExpiredRoles(ctx context.Context) (int64, error)

	UserHasRole(ctx context.Context, userID, roleID string) (bool, error)

	UserHasRoleName(ctx context.Context, userID, roleName string) (bool, error)

	CountUsersByRole(ctx context.Context, roleID string) (int64, error)

	CountRolesByUser(ctx context.Context, userID string) (int64, error)

	Exists(ctx context.Context, userID, roleID string) (bool, error)
}

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
