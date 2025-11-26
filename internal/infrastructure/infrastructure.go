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

func NewConfig() (*config.AppConfig, error) {
	return config.LoadConfig("configs/development.yaml")
}

func NewLogger(cfg *config.AppConfig) (*logger.Logger, error) {
	return logger.NewLogger(cfg.Logger)
}

func NewDatabaseManager(cfg *config.AppConfig, log *logger.Logger) (*database.Manager, error) {
	return database.NewManager(cfg.Database, log)
}

type DatabaseParams struct {
	fx.In
	Manager *database.Manager
}

func NewDatabase(params DatabaseParams) *gorm.DB {
	return params.Manager.Primary()
}

func NewPasswordHasher() *security.PasswordHasher {
	return security.NewPasswordHasher(12) // Default cost of 12
}

func NewTracingService(cfg *config.AppConfig) (*tracing.TracingService, error) {
	return tracing.NewTracingService(cfg)
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return postgresRepos.NewUserRepository(db)
}

func NewRoleRepository(db *gorm.DB) auth.RoleRepository {
	return postgresRepos.NewRoleRepository(db)
}

func NewPermissionRepository(db *gorm.DB) auth.PermissionRepository {
	return postgresRepos.NewPermissionRepository(db)
}

func NewUserRoleRepository(db *gorm.DB) auth.UserRoleRepository {
	return postgresRepos.NewUserRoleRepository(db)
}

func NewRolePermissionRepository(db *gorm.DB) auth.RolePermissionRepository {
	return postgresRepos.NewRolePermissionRepository(db)
}

func NewTokenGenerator() *security.TokenGenerator {
	return security.NewTokenGenerator()
}

type CacheServiceParams struct {
	fx.In
	Config *config.AppConfig
	Logger *logger.Logger
}

func NewRedisClient(config *config.AppConfig) redis.UniversalClient {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,  // default DB
	})
}

func NewCacheService(rdb redis.UniversalClient, logger *logger.Logger) *cache.Service {
	client, ok := rdb.(*redis.Client)
	if !ok {
		client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password
			DB:       0,  // default DB
		})
	}
	return cache.NewCacheService(client, logger)
}

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

			if err := tracingService.Shutdown(ctx); err != nil {
				logger.Error("Failed to shutdown tracing", zap.Error(err))
			}

			if err := manager.Close(); err != nil {
				logger.Error("Failed to close database connections", zap.Error(err))
			}

			logger.Info("Infrastructure shutdown complete")
			return nil
		},
	})
}

func NewFileStorageService(cfg *config.AppConfig, logger *logger.Logger) (contracts.FileStorageService, error) {
	fileStorageConfig := &external.FileStorageConfig{
		Endpoint:        cfg.External.FileStorage.Endpoint,
		AccessKeyID:     cfg.External.FileStorage.AccessKeyID,
		SecretAccessKey: cfg.External.FileStorage.SecretAccessKey,
		BucketName:      cfg.External.FileStorage.BucketName,
		CDNUrl:          cfg.External.FileStorage.CDNUrl,
		UseSSL:          cfg.External.FileStorage.UseSSL,
	}

	return external.NewFileStorageService(fileStorageConfig, logger)
}

func NewMessageBroker(cfg *config.AppConfig, logger *logger.Logger) (messaging.MessageBroker, error) {
	broker := natsAdapter.NewNATSBroker(cfg.Messaging.NATS, logger.Logger)

	if err := broker.Connect(); err != nil {
		return nil, err
	}

	return broker, nil
}

func NewEventBus(broker messaging.MessageBroker, logger *logger.Logger) messaging.EventBus {
	if natsBroker, ok := broker.(*natsAdapter.NATSBroker); ok {
		return natsAdapter.NewNATSEventBus(natsBroker, logger.Logger)
	}

	natsBroker := natsAdapter.NewNATSBroker(config.NATSConfig{
		URL: "nats://localhost:4222",
	}, logger.Logger)

	return natsAdapter.NewNATSEventBus(natsBroker, logger.Logger)
}
