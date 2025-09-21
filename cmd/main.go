package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	fxmodules "github.com/tranvuongduy2003/go-mvc/internal/fx_modules"
)

func main() {
	fx.New(
		// Infrastructure modules
		fxmodules.InfrastructureModule,

		// Repository layer
		fxmodules.RepositoryModule,

		// Domain layer
		fxmodules.DomainModule,

		// Application layer
		fxmodules.ApplicationModule,

		// Handler layer
		fxmodules.HandlerModule,

		// HTTP Server
		fxmodules.ServerModule,

		// Lifecycle hooks
		fx.Invoke(fxmodules.InfrastructureLifecycle),
		fx.Invoke(fxmodules.HTTPServerLifecycle),

		// Logger configuration
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
	).Run()
}
