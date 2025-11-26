package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/tranvuongduy2003/go-mvc/internal/application"
	"github.com/tranvuongduy2003/go-mvc/internal/domain"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/logger"
	"github.com/tranvuongduy2003/go-mvc/internal/presentation"
	"github.com/tranvuongduy2003/go-mvc/internal/presentation/http/handlers"
)

func main() {
	fx.New(
		// Infrastructure modules
		infrastructure.InfrastructureModule,

		// Domain layer - temporarily disabled
		domain.DomainModule,

		// Application layer - temporarily disabled
		application.ApplicationModule,

		// HTTP handlers - temporarily disabled
		handlers.HandlerModule,

		// Server
		presentation.ServerModule,

		// Lifecycle hooks - temporarily disabled due to zap.Logger dependencies
		fx.Invoke(infrastructure.InfrastructureLifecycle),
		fx.Invoke(presentation.SetupMiddleware),
		fx.Invoke(presentation.RegisterRoutes), // Routes after middleware
		fx.Invoke(presentation.HTTPServerLifecycle),

		// Logger configuration
		fx.WithLogger(func(customLogger *logger.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: customLogger.Logger}
		}),
	).Run()
}
