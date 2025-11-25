package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	di "github.com/tranvuongduy2003/go-mvc/internal/di"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/logger"
)

func main() {
	fx.New(
		// Infrastructure modules
		di.InfrastructureModule,

		// Domain layer - temporarily disabled
		di.DomainModule,

		// Application layer - temporarily disabled
		di.ApplicationModule,

		// HTTP handlers - temporarily disabled
		di.HandlerModule,

		// Server
		di.ServerModule,

		// Lifecycle hooks - temporarily disabled due to zap.Logger dependencies
		fx.Invoke(di.InfrastructureLifecycle),
		fx.Invoke(di.SetupMiddleware),
		fx.Invoke(di.RegisterRoutes), // Routes after middleware
		fx.Invoke(di.HTTPServerLifecycle),

		// Logger configuration
		fx.WithLogger(func(customLogger *logger.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: customLogger.Logger}
		}),
	).Run()
}
