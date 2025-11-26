package handlers

import (
	"go.uber.org/fx"

	authCommands "github.com/tranvuongduy2003/go-mvc/internal/application/commands/auth"
	authQueries "github.com/tranvuongduy2003/go-mvc/internal/application/queries/auth"
	appservices "github.com/tranvuongduy2003/go-mvc/internal/application/services"
	userValidators "github.com/tranvuongduy2003/go-mvc/internal/application/validators/user"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/logger"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/tracing"
	v1 "github.com/tranvuongduy2003/go-mvc/internal/presentation/http/handlers/v1"
)

var HandlerModule = fx.Module("handler",
	fx.Provide(
		NewUserHandler,
		NewAuthHandler,
	),
)

type HandlerParams struct {
	fx.In
	Logger  *logger.Logger
	Tracing *tracing.TracingService
}

type AuthHandlerParams struct {
	fx.In
	LoginHandler                *authCommands.LoginCommandHandler
	RegisterHandler             *authCommands.RegisterCommandHandler
	RefreshTokenHandler         *authCommands.RefreshTokenCommandHandler
	ChangePasswordHandler       *authCommands.ChangePasswordCommandHandler
	ResetPasswordHandler        *authCommands.ResetPasswordCommandHandler
	ConfirmPasswordResetHandler *authCommands.ConfirmPasswordResetCommandHandler
	VerifyEmailHandler          *authCommands.VerifyEmailCommandHandler
	ResendVerificationHandler   *authCommands.ResendVerificationEmailCommandHandler
	LogoutHandler               *authCommands.LogoutCommandHandler
	LogoutAllDevicesHandler     *authCommands.LogoutAllDevicesCommandHandler
	GetUserProfileHandler       *authQueries.GetUserProfileQueryHandler
	GetUserPermissionsHandler   *authQueries.GetUserPermissionsQueryHandler
}

func NewUserHandler(userService *appservices.UserService, userValidator userValidators.IUserValidator) *v1.UserHandler {
	return v1.NewUserHandler(userService, userValidator)
}

func NewAuthHandler(params AuthHandlerParams) *v1.AuthHandler {
	return v1.NewAuthHandler(
		params.LoginHandler,
		params.RegisterHandler,
		params.RefreshTokenHandler,
		params.ChangePasswordHandler,
		params.ResetPasswordHandler,
		params.ConfirmPasswordResetHandler,
		params.VerifyEmailHandler,
		params.ResendVerificationHandler,
		params.LogoutHandler,
		params.LogoutAllDevicesHandler,
		params.GetUserProfileHandler,
		params.GetUserPermissionsHandler,
	)
}
