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

var AuthModule = fx.Module("auth",
	fx.Provide(
		NewAuthService,
		NewTokenManagementService,
		NewPasswordManagementService,
		NewEmailVerificationService,
		NewAuthorizationService,
		NewSMTPService,

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

		NewGetUserProfileQueryHandler,
		NewGetUserPermissionsQueryHandler,
	),
)

func NewLoginCommandHandler(authService contracts.AuthService) *authCommands.LoginCommandHandler {
	return authCommands.NewLoginCommandHandler(authService)
}

func NewRegisterCommandHandler(authService contracts.AuthService) *authCommands.RegisterCommandHandler {
	return authCommands.NewRegisterCommandHandler(authService)
}

func NewRefreshTokenCommandHandler(authService contracts.AuthService) *authCommands.RefreshTokenCommandHandler {
	return authCommands.NewRefreshTokenCommandHandler(authService)
}

func NewChangePasswordCommandHandler(passwordService contracts.PasswordManagementService) *authCommands.ChangePasswordCommandHandler {
	return authCommands.NewChangePasswordCommandHandler(passwordService)
}

func NewResetPasswordCommandHandler(passwordService contracts.PasswordManagementService) *authCommands.ResetPasswordCommandHandler {
	return authCommands.NewResetPasswordCommandHandler(passwordService)
}

func NewConfirmPasswordResetCommandHandler(passwordService contracts.PasswordManagementService) *authCommands.ConfirmPasswordResetCommandHandler {
	return authCommands.NewConfirmPasswordResetCommandHandler(passwordService)
}

func NewVerifyEmailCommandHandler(emailVerificationService contracts.EmailVerificationService) *authCommands.VerifyEmailCommandHandler {
	return authCommands.NewVerifyEmailCommandHandler(emailVerificationService)
}

func NewResendVerificationEmailCommandHandler(emailVerificationService contracts.EmailVerificationService) *authCommands.ResendVerificationEmailCommandHandler {
	return authCommands.NewResendVerificationEmailCommandHandler(emailVerificationService)
}

func NewLogoutCommandHandler(tokenService contracts.TokenManagementService) *authCommands.LogoutCommandHandler {
	return authCommands.NewLogoutCommandHandler(tokenService)
}

func NewLogoutAllDevicesCommandHandler(tokenService contracts.TokenManagementService) *authCommands.LogoutAllDevicesCommandHandler {
	return authCommands.NewLogoutAllDevicesCommandHandler(tokenService)
}

func NewGetUserProfileQueryHandler(
	userRepo user.UserRepository,
	roleRepo auth.RoleRepository,
	authorizationService contracts.AuthorizationService,
) *authQueries.GetUserProfileQueryHandler {
	return authQueries.NewGetUserProfileQueryHandler(userRepo, roleRepo, authorizationService)
}

func NewGetUserPermissionsQueryHandler(authorizationService contracts.AuthorizationService) *authQueries.GetUserPermissionsQueryHandler {
	return authQueries.NewGetUserPermissionsQueryHandler(authorizationService)
}

type AuthServiceParams struct {
	fx.In
	UserRepo       user.UserRepository
	JWTService     jwt.JWTService
	PasswordHasher *security.PasswordHasher
	CacheService   *cache.Service
	SMTPService    *external.SMTPService
	Logger         *logger.Logger
}

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

type AuthorizationServiceParams struct {
	fx.In
	UserRepo           user.UserRepository
	RoleRepo           auth.RoleRepository
	PermissionRepo     auth.PermissionRepository
	UserRoleRepo       auth.UserRoleRepository
	RolePermissionRepo auth.RolePermissionRepository
}

func NewAuthorizationService(params AuthorizationServiceParams) contracts.AuthorizationService {
	return appServices.NewAuthorizationService(
		params.UserRepo,
		params.RoleRepo,
		params.PermissionRepo,
		params.UserRoleRepo,
		params.RolePermissionRepo,
	)
}

func NewSMTPService(cfg *config.AppConfig, logger *logger.Logger) *external.SMTPService {
	return external.NewSMTPService(&cfg.External.EmailService.SMTP, logger)
}
