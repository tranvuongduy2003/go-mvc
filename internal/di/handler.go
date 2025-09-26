package di

import (
	"go.uber.org/fx"

	// "github.com/tranvuongduy2003/go-mvc/internal/adapters/cache" // Commented out - not used currently
	authCommands "github.com/tranvuongduy2003/go-mvc/internal/application/commands/auth"
	authQueries "github.com/tranvuongduy2003/go-mvc/internal/application/queries/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/application/services"
	userValidators "github.com/tranvuongduy2003/go-mvc/internal/application/validators/user"
	v1 "github.com/tranvuongduy2003/go-mvc/internal/handlers/http/rest/v1"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/tracing"
)

// HandlerModule provides handler layer dependencies
var HandlerModule = fx.Module("handler",
	fx.Provide(
		NewUserHandler,
		NewAuthHandler,
	),
)

// HandlerParams holds parameters for handler providers
type HandlerParams struct {
	fx.In
	Logger  *logger.Logger
	Tracing *tracing.TracingService
}

// AuthHandlerParams holds parameters for auth handler
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

// NewUserHandler provides UserHandler
func NewUserHandler(userService *services.UserService, userValidator userValidators.IUserValidator) *v1.UserHandler {
	return v1.NewUserHandler(userService, userValidator)
}

// NewAuthHandler provides AuthHandler
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
