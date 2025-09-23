package fxmodules

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"

	postgresRepos "github.com/tranvuongduy2003/go-mvc/internal/adapters/persistence/postgres/repositories"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/repositories"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/config"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/database"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/security"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/tracing"
)

// InfrastructureModule provides infrastructure dependencies
var InfrastructureModule = fx.Module("infrastructure",
	fx.Provide(
		NewConfig,
		NewLogger,
		NewDatabaseManager,
		NewDatabase,
		NewPasswordHasher,
		NewTracingService,
		NewUserRepository,
	),
)

// NewConfig provides application configuration
func NewConfig() (*config.AppConfig, error) {
	return config.LoadConfig("configs/development.yaml")
}

// NewLogger provides application logger
func NewLogger(cfg *config.AppConfig) (*logger.Logger, error) {
	return logger.NewLogger(cfg.Logger)
}

// NewDatabaseManager provides database manager
func NewDatabaseManager(cfg *config.AppConfig, log *logger.Logger) (*database.Manager, error) {
	return database.NewManager(cfg.Database, log)
}

// DatabaseParams holds parameters for database provider
type DatabaseParams struct {
	fx.In
	Manager *database.Manager
}

// NewDatabase provides primary database connection
func NewDatabase(params DatabaseParams) *gorm.DB {
	return params.Manager.Primary()
}

// NewPasswordHasher provides password hasher
func NewPasswordHasher() *security.PasswordHasher {
	return security.NewPasswordHasher(12) // Default cost of 12
}

// NewTracingService provides tracing service
func NewTracingService(cfg *config.AppConfig) (*tracing.TracingService, error) {
	return tracing.NewTracingService(cfg)
}

// NewUserRepository provides user repository
func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return postgresRepos.NewUserRepository(db)
}

// InfrastructureLifecycle handles infrastructure lifecycle
func InfrastructureLifecycle(
	lc fx.Lifecycle,
	manager *database.Manager,
	tracingService *tracing.TracingService,
	logger *logger.Logger,
) {
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
}
