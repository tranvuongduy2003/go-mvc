package infrastructure

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/messaging"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/user"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/cache"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/config"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/database"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/external"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/logger"
	natsAdapter "github.com/tranvuongduy2003/go-mvc/internal/infrastructure/messaging/nats"
	postgresRepos "github.com/tranvuongduy2003/go-mvc/internal/infrastructure/persistence/postgres/repositories"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/security"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/tracing"
)

// InfrastructureModule provides infrastructure dependencies
var InfrastructureModule = fx.Module("infrastructure",
	fx.Provide(
		NewConfig,
		NewLogger,
		NewDatabaseManager,
		NewDatabase,
		NewRedisClient,
		NewPasswordHasher,
		NewTokenGenerator,
		NewCacheService,
		NewTracingService,
		NewFileStorageService,
		NewMessageBroker,
		NewEventBus,
		NewUserRepository,
		NewRoleRepository,
		NewPermissionRepository,
		NewUserRoleRepository,
		NewRolePermissionRepository,
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
func NewUserRepository(db *gorm.DB) user.UserRepository {
	return postgresRepos.NewUserRepository(db)
}

// NewRoleRepository provides role repository
func NewRoleRepository(db *gorm.DB) auth.RoleRepository {
	return postgresRepos.NewRoleRepository(db)
}

// NewPermissionRepository provides permission repository
func NewPermissionRepository(db *gorm.DB) auth.PermissionRepository {
	return postgresRepos.NewPermissionRepository(db)
}

// NewUserRoleRepository provides user role repository
func NewUserRoleRepository(db *gorm.DB) auth.UserRoleRepository {
	return postgresRepos.NewUserRoleRepository(db)
}

// NewRolePermissionRepository provides role permission repository
func NewRolePermissionRepository(db *gorm.DB) auth.RolePermissionRepository {
	return postgresRepos.NewRolePermissionRepository(db)
}

// NewTokenGenerator provides token generator
func NewTokenGenerator() *security.TokenGenerator {
	return security.NewTokenGenerator()
}

// CacheServiceParams holds parameters for cache service
type CacheServiceParams struct {
	fx.In
	Config *config.AppConfig
	Logger *logger.Logger
}

// NewRedisClient provides a shared Redis client
func NewRedisClient(config *config.AppConfig) redis.UniversalClient {
	// For now, create a Redis client - can be configured based on config later
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,  // default DB
	})
}

// NewCacheService provides cache service
func NewCacheService(rdb redis.UniversalClient, logger *logger.Logger) *cache.Service {
	// Convert UniversalClient to concrete *redis.Client
	client, ok := rdb.(*redis.Client)
	if !ok {
		// If it's not a *redis.Client, create a new one with same options
		client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password
			DB:       0,  // default DB
		})
	}
	return cache.NewCacheService(client, logger)
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

// NewFileStorageService provides file storage service as interface
// Returns the port interface, not the concrete implementation
// This follows Dependency Inversion Principle
func NewFileStorageService(cfg *config.AppConfig, logger *logger.Logger) (contracts.FileStorageService, error) {
	fileStorageConfig := &external.FileStorageConfig{
		Endpoint:        cfg.External.FileStorage.Endpoint,
		AccessKeyID:     cfg.External.FileStorage.AccessKeyID,
		SecretAccessKey: cfg.External.FileStorage.SecretAccessKey,
		BucketName:      cfg.External.FileStorage.BucketName,
		CDNUrl:          cfg.External.FileStorage.CDNUrl,
		UseSSL:          cfg.External.FileStorage.UseSSL,
	}

	// Return as interface type, not concrete type
	return external.NewFileStorageService(fileStorageConfig, logger)
}

// NewMessageBroker provides NATS message broker
func NewMessageBroker(cfg *config.AppConfig, logger *logger.Logger) (messaging.MessageBroker, error) {
	broker := natsAdapter.NewNATSBroker(cfg.Messaging.NATS, logger.Logger)

	if err := broker.Connect(); err != nil {
		return nil, err
	}

	return broker, nil
}

// NewEventBus provides NATS event bus
func NewEventBus(broker messaging.MessageBroker, logger *logger.Logger) messaging.EventBus {
	if natsBroker, ok := broker.(*natsAdapter.NATSBroker); ok {
		return natsAdapter.NewNATSEventBus(natsBroker, logger.Logger)
	}

	// Fallback - should not happen in normal circumstances
	natsBroker := natsAdapter.NewNATSBroker(config.NATSConfig{
		URL: "nats://localhost:4222",
	}, logger.Logger)

	return natsAdapter.NewNATSEventBus(natsBroker, logger.Logger)
}
