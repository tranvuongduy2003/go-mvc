package main

import (
	"context"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	fxmodules "github.com/tranvuongduy2003/go-mvc/internal/fx_modules"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/config"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/database"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/tracing"
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
		fx.Invoke(func(lc fx.Lifecycle, manager *database.Manager, tracingService *tracing.TracingService, logger *logger.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error {
					logger.Info("Infrastructure started successfully")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					logger.Info("Shutting down infrastructure...")

					// Shutdown tracing
					if err := tracingService.Shutdown(ctx); err != nil {
						logger.Error("Failed to shutdown tracing", zap.Error(err))
					}

					// Close database connections
					if err := manager.Close(); err != nil {
						logger.Error("Failed to close database connections", zap.Error(err))
					}

					logger.Info("Infrastructure shutdown complete")
					return nil
				},
			})
		}),
		fx.Invoke(func(lc fx.Lifecycle, server *http.Server, config *config.AppConfig, logger *logger.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logger.Info("Starting HTTP server",
						zap.String("addr", server.Addr),
						zap.String("environment", config.App.Environment),
					)

					go func() {
						if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
							logger.Fatal("Failed to start HTTP server", zap.Error(err))
						}
					}()

					return nil
				},
				OnStop: func(ctx context.Context) error {
					logger.Info("Shutting down HTTP server...")
					return server.Shutdown(ctx)
				},
			})
		}),

		// Logger configuration
		fx.WithLogger(func(customLogger *logger.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: customLogger.Logger}
		}),
	).Run()
}
