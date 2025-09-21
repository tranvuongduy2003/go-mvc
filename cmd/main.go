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

		// Repository layer
		fxmodules.RepositoryModule,

		// Domain layer - temporarily disabled
		// fxmodules.DomainModule,

		// Application layer - temporarily disabled
		// fxmodules.ApplicationModule,

		// HTTP handlers - temporarily disabled
		// fxmodules.HandlerModule,

		// Server
		// fxmodules.ServerModule,

		// Lifecycle hooks
		fx.Invoke(fxmodules.InfrastructureLifecycle),
		// fx.Invoke(fxmodules.HTTPServerLifecycle),

		// Logger configuration
		fx.WithLogger(func(customLogger *logger.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: customLogger.Logger}
		}),
	).Run()
}
