package fxmodules

import (
	"go.uber.org/fx"

	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/rbac"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/user"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/security"
)

// DomainModule provides domain layer dependencies
var DomainModule = fx.Module("domain",
	fx.Provide(
		NewUserDomainService,
		NewRBACService,
	),
)

// DomainParams holds parameters for domain service providers
type DomainParams struct {
	fx.In
	Repository     user.Repository
	PasswordHasher *security.PasswordHasher
	Logger         *logger.Logger
}

// RBACParams holds parameters for RBAC service
type RBACParams struct {
	fx.In
	RoleRepository           rbac.RoleRepository
	PermissionRepository     rbac.PermissionRepository
	UserRoleRepository       rbac.UserRoleRepository
	RolePermissionRepository rbac.RolePermissionRepository
}

// NewUserDomainService provides user domain service
func NewUserDomainService(params DomainParams) *user.Service {
	return user.NewService(params.Repository, params.PasswordHasher, params.Logger)
}

// NewRBACService provides RBAC service
func NewRBACService(params RBACParams) rbac.RBACService {
	return rbac.NewRBACService(
		params.RoleRepository,
		params.PermissionRepository,
		params.UserRoleRepository,
		params.RolePermissionRepository,
	)
}
