package modules

import (
	"go.uber.org/fx"

	"github.com/tranvuongduy2003/go-mvc/internal/adapters/cache"
	authCommands "github.com/tranvuongduy2003/go-mvc/internal/application/commands/auth"
	authQueries "github.com/tranvuongduy2003/go-mvc/internal/application/queries/auth"
	appServices "github.com/tranvuongduy2003/go-mvc/internal/application/services"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/repositories"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/services"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/security"
	"github.com/tranvuongduy2003/go-mvc/pkg/jwt"
)

// AuthModule provides authentication and authorization dependencies
var AuthModule = fx.Module("auth",
	fx.Provide(
		// Command Handlers
		NewLoginCommandHandler,
		NewRegisterCommandHandler,
		NewRefreshTokenCommandHandler,
		NewChangePasswordCommandHandler,
		NewResetPasswordCommandHandler,
		NewConfirmPasswordResetCommandHandler,
		NewVerifyEmailCommandHandler,
		NewResendVerificationEmailCommandHandler,
		NewLogoutCommandHandler,
		NewLogoutAllDevicesCommandHandler,

		// Query Handlers
		NewGetUserProfileQueryHandler,
		NewGetUserPermissionsQueryHandler,

		// Services
		NewAuthService,
		NewAuthorizationService,
	),
)

// Command Handler Providers

// NewLoginCommandHandler provides LoginCommandHandler
func NewLoginCommandHandler(authService services.AuthService) *authCommands.LoginCommandHandler {
	return authCommands.NewLoginCommandHandler(authService)
}

// NewRegisterCommandHandler provides RegisterCommandHandler
func NewRegisterCommandHandler(authService services.AuthService) *authCommands.RegisterCommandHandler {
	return authCommands.NewRegisterCommandHandler(authService)
}

// NewRefreshTokenCommandHandler provides RefreshTokenCommandHandler
func NewRefreshTokenCommandHandler(authService services.AuthService) *authCommands.RefreshTokenCommandHandler {
	return authCommands.NewRefreshTokenCommandHandler(authService)
}

// NewChangePasswordCommandHandler provides ChangePasswordCommandHandler
func NewChangePasswordCommandHandler(authService services.AuthService) *authCommands.ChangePasswordCommandHandler {
	return authCommands.NewChangePasswordCommandHandler(authService)
}

// NewResetPasswordCommandHandler provides ResetPasswordCommandHandler
func NewResetPasswordCommandHandler(authService services.AuthService) *authCommands.ResetPasswordCommandHandler {
	return authCommands.NewResetPasswordCommandHandler(authService)
}

// NewConfirmPasswordResetCommandHandler provides ConfirmPasswordResetCommandHandler
func NewConfirmPasswordResetCommandHandler(authService services.AuthService) *authCommands.ConfirmPasswordResetCommandHandler {
	return authCommands.NewConfirmPasswordResetCommandHandler(authService)
}

// NewVerifyEmailCommandHandler provides VerifyEmailCommandHandler
func NewVerifyEmailCommandHandler(authService services.AuthService) *authCommands.VerifyEmailCommandHandler {
	return authCommands.NewVerifyEmailCommandHandler(authService)
}

// NewResendVerificationEmailCommandHandler provides ResendVerificationEmailCommandHandler
func NewResendVerificationEmailCommandHandler(authService services.AuthService) *authCommands.ResendVerificationEmailCommandHandler {
	return authCommands.NewResendVerificationEmailCommandHandler(authService)
}

// NewLogoutCommandHandler provides LogoutCommandHandler
func NewLogoutCommandHandler(authService services.AuthService) *authCommands.LogoutCommandHandler {
	return authCommands.NewLogoutCommandHandler(authService)
}

// NewLogoutAllDevicesCommandHandler provides LogoutAllDevicesCommandHandler
func NewLogoutAllDevicesCommandHandler(authService services.AuthService) *authCommands.LogoutAllDevicesCommandHandler {
	return authCommands.NewLogoutAllDevicesCommandHandler(authService)
}

// Query Handler Providers

// NewGetUserProfileQueryHandler provides GetUserProfileQueryHandler
func NewGetUserProfileQueryHandler(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	authorizationService services.AuthorizationService,
) *authQueries.GetUserProfileQueryHandler {
	return authQueries.NewGetUserProfileQueryHandler(userRepo, roleRepo, authorizationService)
}

// NewGetUserPermissionsQueryHandler provides GetUserPermissionsQueryHandler
func NewGetUserPermissionsQueryHandler(authorizationService services.AuthorizationService) *authQueries.GetUserPermissionsQueryHandler {
	return authQueries.NewGetUserPermissionsQueryHandler(authorizationService)
}

// Service Providers

// AuthServiceParams holds parameters for AuthService
type AuthServiceParams struct {
	fx.In
	UserRepo       repositories.UserRepository
	JWTService     jwt.JWTService
	PasswordHasher *security.PasswordHasher
	CacheService   *cache.Service
}

// NewAuthService provides AuthService
func NewAuthService(params AuthServiceParams) services.AuthService {
	return appServices.NewAuthService(
		params.UserRepo,
		params.JWTService,
		params.PasswordHasher,
		params.CacheService,
	)
}

// AuthorizationServiceParams holds parameters for AuthorizationService
type AuthorizationServiceParams struct {
	fx.In
	UserRepo           repositories.UserRepository
	RoleRepo           repositories.RoleRepository
	PermissionRepo     repositories.PermissionRepository
	UserRoleRepo       repositories.UserRoleRepository
	RolePermissionRepo repositories.RolePermissionRepository
}

// NewAuthorizationService provides AuthorizationService
func NewAuthorizationService(params AuthorizationServiceParams) services.AuthorizationService {
	return appServices.NewAuthorizationService(
		params.UserRepo,
		params.RoleRepo,
		params.PermissionRepo,
		params.UserRoleRepo,
		params.RolePermissionRepo,
	)
}
