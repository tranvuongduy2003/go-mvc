package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	fxmodules "github.com/tranvuongduy2003/go-mvc/internal/fx_modules"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
)

func main() {
	fx.New(
		// Infrastructure modules
		fxmodules.InfrastructureModule,

		// Domain layer - temporarily disabled
		fxmodules.DomainModule,

		// Application layer - temporarily disabled
		fxmodules.ApplicationModule,

		// HTTP handlers - temporarily disabled
		fxmodules.HandlerModule,

		// Server
		fxmodules.ServerModule,

		// Lifecycle hooks - temporarily disabled due to zap.Logger dependencies
		fx.Invoke(fxmodules.InfrastructureLifecycle),
		fx.Invoke(fxmodules.SetupMiddleware),
		fx.Invoke(fxmodules.RegisterRoutes), // Routes after middleware
		fx.Invoke(fxmodules.HTTPServerLifecycle),

		// Logger configuration
		fx.WithLogger(func(customLogger *logger.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: customLogger.Logger}
		}),
	).Run()
}
