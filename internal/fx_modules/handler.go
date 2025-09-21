package fxmodules

import (
	"go.uber.org/fx"

	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/tracing"
)

// HandlerModule provides handler layer dependencies
var HandlerModule = fx.Module("handler",
	fx.Provide(),
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
	Logger *logger.Logger
}
