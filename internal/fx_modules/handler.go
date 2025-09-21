package fxmodules

import (
	"go.uber.org/fx"

	"github.com/tranvuongduy2003/go-mvc/internal/application/services"
	"github.com/tranvuongduy2003/go-mvc/internal/application/validators"
	"github.com/tranvuongduy2003/go-mvc/internal/handlers"
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
	UserService   *services.UserApplicationService
	UserValidator *validators.UserValidator
	Logger        *logger.Logger
	Tracing       *tracing.TracingService
}

// AuthHandlerParams holds parameters for auth handler
type AuthHandlerParams struct {
	fx.In
	UserService   *services.UserApplicationService
	UserValidator *validators.UserValidator
	Logger        *logger.Logger
}

// NewUserHandler provides user HTTP handler
func NewUserHandler(params HandlerParams) *handlers.UserHandler {
	return handlers.NewUserHandler(
		params.UserService,
		params.UserValidator,
		params.Logger,
		params.Tracing,
	)
}

// NewAuthHandler provides auth HTTP handler
func NewAuthHandler(params AuthHandlerParams) *handlers.AuthHandler {
	return handlers.NewAuthHandler(
		params.UserService,
		params.UserValidator,
		params.Logger,
	)
}
