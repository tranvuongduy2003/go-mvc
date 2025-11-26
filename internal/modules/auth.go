package modules

import (
	"go.uber.org/fx"

	authCommands "github.com/tranvuongduy2003/go-mvc/internal/application/commands/auth"
	authQueries "github.com/tranvuongduy2003/go-mvc/internal/application/queries/auth"
	appServices "github.com/tranvuongduy2003/go-mvc/internal/application/services"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/user"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/cache"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/config"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/external"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/logger"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/security"
	"github.com/tranvuongduy2003/go-mvc/pkg/jwt"
)

// AuthModule provides authentication and authorization dependencies
var AuthModule = fx.Module("auth",
	fx.Provide(
		// Services - provide concrete service as all interface types
		NewAuthService,
		NewTokenManagementService,
		NewPasswordManagementService,
		NewEmailVerificationService,
		NewAuthorizationService,
		NewSMTPService,

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
	),
)

// Command Handler Providers

// NewLoginCommandHandler provides LoginCommandHandler
func NewLoginCommandHandler(authService contracts.AuthService) *authCommands.LoginCommandHandler {
	return authCommands.NewLoginCommandHandler(authService)
}

// NewRegisterCommandHandler provides RegisterCommandHandler
func NewRegisterCommandHandler(authService contracts.AuthService) *authCommands.RegisterCommandHandler {
	return authCommands.NewRegisterCommandHandler(authService)
}

// NewRefreshTokenCommandHandler provides RefreshTokenCommandHandler
func NewRefreshTokenCommandHandler(authService contracts.AuthService) *authCommands.RefreshTokenCommandHandler {
	return authCommands.NewRefreshTokenCommandHandler(authService)
}

// NewChangePasswordCommandHandler provides ChangePasswordCommandHandler
func NewChangePasswordCommandHandler(passwordService contracts.PasswordManagementService) *authCommands.ChangePasswordCommandHandler {
	return authCommands.NewChangePasswordCommandHandler(passwordService)
}

// NewResetPasswordCommandHandler provides ResetPasswordCommandHandler
func NewResetPasswordCommandHandler(passwordService contracts.PasswordManagementService) *authCommands.ResetPasswordCommandHandler {
	return authCommands.NewResetPasswordCommandHandler(passwordService)
}

// NewConfirmPasswordResetCommandHandler provides ConfirmPasswordResetCommandHandler
func NewConfirmPasswordResetCommandHandler(passwordService contracts.PasswordManagementService) *authCommands.ConfirmPasswordResetCommandHandler {
	return authCommands.NewConfirmPasswordResetCommandHandler(passwordService)
}

// NewVerifyEmailCommandHandler provides VerifyEmailCommandHandler
func NewVerifyEmailCommandHandler(emailVerificationService contracts.EmailVerificationService) *authCommands.VerifyEmailCommandHandler {
	return authCommands.NewVerifyEmailCommandHandler(emailVerificationService)
}

// NewResendVerificationEmailCommandHandler provides ResendVerificationEmailCommandHandler
func NewResendVerificationEmailCommandHandler(emailVerificationService contracts.EmailVerificationService) *authCommands.ResendVerificationEmailCommandHandler {
	return authCommands.NewResendVerificationEmailCommandHandler(emailVerificationService)
}

// NewLogoutCommandHandler provides LogoutCommandHandler
func NewLogoutCommandHandler(tokenService contracts.TokenManagementService) *authCommands.LogoutCommandHandler {
	return authCommands.NewLogoutCommandHandler(tokenService)
}

// NewLogoutAllDevicesCommandHandler provides LogoutAllDevicesCommandHandler
func NewLogoutAllDevicesCommandHandler(tokenService contracts.TokenManagementService) *authCommands.LogoutAllDevicesCommandHandler {
	return authCommands.NewLogoutAllDevicesCommandHandler(tokenService)
}

// Query Handler Providers

// NewGetUserProfileQueryHandler provides GetUserProfileQueryHandler
func NewGetUserProfileQueryHandler(
	userRepo user.UserRepository,
	roleRepo auth.RoleRepository,
	authorizationService contracts.AuthorizationService,
) *authQueries.GetUserProfileQueryHandler {
	return authQueries.NewGetUserProfileQueryHandler(userRepo, roleRepo, authorizationService)
}

// NewGetUserPermissionsQueryHandler provides GetUserPermissionsQueryHandler
func NewGetUserPermissionsQueryHandler(authorizationService contracts.AuthorizationService) *authQueries.GetUserPermissionsQueryHandler {
	return authQueries.NewGetUserPermissionsQueryHandler(authorizationService)
}

// Service Providers

// AuthServiceParams holds parameters for AuthService
type AuthServiceParams struct {
	fx.In
	UserRepo       user.UserRepository
	JWTService     jwt.JWTService
	PasswordHasher *security.PasswordHasher
	CacheService   *cache.Service
	SMTPService    *external.SMTPService
	Logger         *logger.Logger
}

// NewAuthService provides AuthService interface
func NewAuthService(params AuthServiceParams) contracts.AuthService {
	return appServices.NewAuthService(
		params.UserRepo,
		params.JWTService,
		params.PasswordHasher,
		params.CacheService,
		params.SMTPService,
		params.Logger,
	)
}

// NewTokenManagementService provides TokenManagementService
// The concrete auth service implements all split interfaces
func NewTokenManagementService(params AuthServiceParams) contracts.TokenManagementService {
	return appServices.NewAuthService(
		params.UserRepo,
		params.JWTService,
		params.PasswordHasher,
		params.CacheService,
		params.SMTPService,
		params.Logger,
	)
}

// NewPasswordManagementService provides PasswordManagementService
// The concrete auth service implements all split interfaces
func NewPasswordManagementService(params AuthServiceParams) contracts.PasswordManagementService {
	return appServices.NewAuthService(
		params.UserRepo,
		params.JWTService,
		params.PasswordHasher,
		params.CacheService,
		params.SMTPService,
		params.Logger,
	)
}

// NewEmailVerificationService provides EmailVerificationService
// The concrete auth service implements all split interfaces
func NewEmailVerificationService(params AuthServiceParams) contracts.EmailVerificationService {
	return appServices.NewAuthService(
		params.UserRepo,
		params.JWTService,
		params.PasswordHasher,
		params.CacheService,
		params.SMTPService,
		params.Logger,
	)
}

// AuthorizationServiceParams holds parameters for AuthorizationService
type AuthorizationServiceParams struct {
	fx.In
	UserRepo           user.UserRepository
	RoleRepo           auth.RoleRepository
	PermissionRepo     auth.PermissionRepository
	UserRoleRepo       auth.UserRoleRepository
	RolePermissionRepo auth.RolePermissionRepository
}

// NewAuthorizationService provides AuthorizationService
func NewAuthorizationService(params AuthorizationServiceParams) contracts.AuthorizationService {
	return appServices.NewAuthorizationService(
		params.UserRepo,
		params.RoleRepo,
		params.PermissionRepo,
		params.UserRoleRepo,
		params.RolePermissionRepo,
	)
}

// NewSMTPService provides SMTPService
func NewSMTPService(cfg *config.AppConfig, logger *logger.Logger) *external.SMTPService {
	return external.NewSMTPService(&cfg.External.EmailService.SMTP, logger)
}
