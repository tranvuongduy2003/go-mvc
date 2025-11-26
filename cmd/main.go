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
		infrastructure.InfrastructureModule,

		domain.DomainModule,

		application.ApplicationModule,

		handlers.HandlerModule,

		presentation.ServerModule,

		fx.Invoke(infrastructure.InfrastructureLifecycle),
		fx.Invoke(presentation.SetupMiddleware),
		fx.Invoke(presentation.RegisterRoutes), // Routes after middleware
		fx.Invoke(presentation.HTTPServerLifecycle),

		fx.WithLogger(func(customLogger *logger.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: customLogger.Logger}
		}),
	).Run()
}
