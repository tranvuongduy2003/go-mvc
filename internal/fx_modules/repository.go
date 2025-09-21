package fxmodules

import (
	"go.uber.org/fx"
	"gorm.io/gorm"

	postgresRepos "github.com/tranvuongduy2003/go-mvc/internal/adapters/persistence/postgres/repositories"
	"github.com/tranvuongduy2003/go-mvc/internal/adapters/repositories"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/rbac"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/user"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/tracing"
)

// RepositoryModule provides repository layer dependencies
var RepositoryModule = fx.Module("repository",
	fx.Provide(
		NewUserRepository,
		NewRoleRepository,
		NewPermissionRepository,
		NewUserRoleRepository,
		NewRolePermissionRepository,
		fx.Annotate(
			NewUserRepository,
			fx.As(new(user.Repository)),
		),
		fx.Annotate(
			NewRoleRepository,
			fx.As(new(rbac.RoleRepository)),
		),
		fx.Annotate(
			NewPermissionRepository,
			fx.As(new(rbac.PermissionRepository)),
		),
		fx.Annotate(
			NewUserRoleRepository,
			fx.As(new(rbac.UserRoleRepository)),
		),
		fx.Annotate(
			NewRolePermissionRepository,
			fx.As(new(rbac.RolePermissionRepository)),
		),
	),
)

// RepositoryParams holds parameters for repository providers
type RepositoryParams struct {
	fx.In
	DB      *gorm.DB
	Logger  *logger.Logger
	Tracing *tracing.TracingService
}

// NewUserRepository provides user repository implementation
func NewUserRepository(params RepositoryParams) *repositories.UserRepository {
	return repositories.NewUserRepository(params.DB, params.Logger, params.Tracing)
}

// NewRoleRepository provides role repository implementation
func NewRoleRepository(params RepositoryParams) rbac.RoleRepository {
	return postgresRepos.NewRoleRepository(params.DB)
}

// NewPermissionRepository provides permission repository implementation
func NewPermissionRepository(params RepositoryParams) rbac.PermissionRepository {
	return postgresRepos.NewPermissionRepository(params.DB)
}

// NewUserRoleRepository provides user role repository implementation
func NewUserRoleRepository(params RepositoryParams) rbac.UserRoleRepository {
	return postgresRepos.NewUserRoleRepository(params.DB)
}

// NewRolePermissionRepository provides role permission repository implementation
func NewRolePermissionRepository(params RepositoryParams) rbac.RolePermissionRepository {
	return postgresRepos.NewRolePermissionRepository(params.DB)
}
