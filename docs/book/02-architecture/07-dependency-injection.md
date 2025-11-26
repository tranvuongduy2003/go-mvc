# Dependency Injection với Uber Fx

## Tổng quan

Ứng dụng đã được tích hợp dependency injection sử dụng Uber Fx, một framework mạnh mẽ để quản lý dependencies và lifecycle của ứng dụng Go.

> **Note**: Dependency injection modules được đặt trong folder `internal/di/` (viết tắt của "Dependency Injection") thay vì tên cũ `fx_modules` để tuân theo naming conventions phổ biến trong industry.

## Tính năng Fx Dependency Injection

### 1. **Modular Architecture**
- **InfrastructureModule**: Cung cấp các dependencies cơ bản (config, logger, database, tracing)
- **RepositoryModule**: Quản lý repository layer với interface binding
- **DomainModule**: Cung cấp domain services
- **ApplicationModule**: Quản lý application services và validators
- **HandlerModule**: Cung cấp HTTP handlers
- **ServerModule**: Quản lý HTTP server và routing

### 2. **Automatic Dependency Resolution**
- Tự động resolve dependencies thông qua constructor injection
- Type-safe dependency injection với compile-time checking
- Interface binding cho loose coupling

### 3. **Lifecycle Management**
- Graceful startup và shutdown
- Lifecycle hooks cho từng component
- Resource cleanup tự động

## Cấu trúc Modules

### Infrastructure Module (`internal/di/infrastructure.go`)
```go
var InfrastructureModule = fx.Module("infrastructure",
    fx.Provide(
        NewConfig,           // Application configuration
        NewLogger,           // Structured logging
        NewDatabaseManager,  // Database connection manager
        NewDatabase,         // Primary database connection
        NewPasswordHasher,   // Password hashing service
        NewTracingService,   // Distributed tracing
    ),
)
```

### Repository Module (`internal/di/repository.go`)
```go
var RepositoryModule = fx.Module("repository",
    fx.Provide(
        NewUserRepository,
        fx.Annotate(
            NewUserRepository,
            fx.As(new(user.Repository)), // Interface binding
        ),
    ),
)
```

### Domain Module (`internal/di/domain.go`)
```go
var DomainModule = fx.Module("domain",
    fx.Provide(
        NewUserDomainService, // Domain business logic
    ),
)
```

### Application Module (`internal/di/application.go`)
```go
var ApplicationModule = fx.Module("application",
    fx.Provide(
        NewJWTService,            // JWT token service
        NewUserApplicationService, // Application service layer
        NewUserValidator,         // Request validation
    ),
)
```

### Handler Module (`internal/di/handler.go`)
```go
var HandlerModule = fx.Module("handler",
    fx.Provide(
        NewUserHandler, // HTTP user endpoints
        NewAuthHandler, // HTTP auth endpoints
    ),
)
```

### Server Module (`internal/di/server.go`)
```go
var ServerModule = fx.Module("server",
    fx.Provide(
        NewHTTPServer, // HTTP server instance
        NewGinRouter,  // Gin router with middleware
    ),
    fx.Invoke(RegisterRoutes), // Route registration
)
```

## Constructor Functions

### Infrastructure Constructors
```go
// NewConfig provides application configuration
func NewConfig() (*config.AppConfig, error) {
    return config.LoadConfig("development")
}

// NewLogger provides application logger
func NewLogger(cfg *config.AppConfig) (*logger.Logger, error) {
    return logger.NewLogger(cfg.Logger)
}

// NewDatabase provides primary database connection
func NewDatabase(params DatabaseParams) *gorm.DB {
    return params.Manager.Primary()
}
```

### Parameter Structs
```go
// DatabaseParams holds parameters for database provider
type DatabaseParams struct {
    fx.In
    Manager *database.Manager
}

// ApplicationParams holds parameters for application service
type ApplicationParams struct {
    fx.In
    UserService *user.Service
    Repository  user.Repository
    JWTService  *jwt.Service
    Logger      *logger.Logger
    Tracing     *tracing.TracingService
}
```

## Main Application (`cmd/fx_api/main.go`)

```go
func main() {
    fx.New(
        // Infrastructure modules
        di.InfrastructureModule,
        
        // Repository layer
        di.RepositoryModule,
        
        // Domain layer
        di.DomainModule,
        
        // Application layer
        di.ApplicationModule,
        
        // Handler layer
        di.HandlerModule,
        
        // HTTP Server
        di.ServerModule,
        
        // Lifecycle hooks
        fx.Invoke(di.InfrastructureLifecycle),
        fx.Invoke(di.HTTPServerLifecycle),
        
        // Logger configuration
        fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
            return &fxevent.ZapLogger{Logger: logger}
        }),
    ).Run()
}
```

## Lifecycle Management

### Infrastructure Lifecycle
```go
func InfrastructureLifecycle(
    lc fx.Lifecycle,
    manager *database.Manager,
    tracingService *tracing.TracingService,
    zapLogger *zap.Logger,
) {
    lc.Append(fx.Hook{
        OnStart: func(context.Context) error {
            zapLogger.Info("Infrastructure started successfully")
            return nil
        },
        OnStop: func(ctx context.Context) error {
            zapLogger.Info("Shutting down infrastructure...")
            
            // Shutdown tracing
            if err := tracingService.Shutdown(ctx); err != nil {
                zapLogger.Error("Failed to shutdown tracing", zap.Error(err))
            }
            
            // Close database connections
            if err := manager.Close(); err != nil {
                zapLogger.Error("Failed to close database connections", zap.Error(err))
            }
            
            return nil
        },
    })
}
```

### HTTP Server Lifecycle
```go
func HTTPServerLifecycle(
    lc fx.Lifecycle,
    server *http.Server,
    config *config.AppConfig,
    zapLogger *zap.Logger,
) {
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            zapLogger.Info("Starting HTTP server", 
                zap.String("addr", server.Addr))
            
            go func() {
                if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
                    zapLogger.Fatal("Failed to start HTTP server", zap.Error(err))
                }
            }()
            
            return nil
        },
        OnStop: func(ctx context.Context) error {
            zapLogger.Info("Shutting down HTTP server...")
            return server.Shutdown(ctx)
        },
    })
}
```

## Advantages của Fx

### 1. **Type Safety**
- Compile-time dependency checking
- No runtime injection errors
- Clear dependency graph

### 2. **Modularity**
- Clear separation of concerns
- Reusable modules
- Easy testing with module substitution

### 3. **Lifecycle Management**
- Automatic startup ordering
- Graceful shutdown
- Resource cleanup

### 4. **Debugging**
- Clear dependency visualization
- Startup/shutdown logging
- Error propagation

## Running the Application

### Demo Application
```bash
# Build and run Fx demo
go run cmd/fx_demo/main.go

# Test endpoints
curl http://localhost:8080/health
curl http://localhost:8080/fx-demo
```

### Full Application
```bash
# Build and run full Fx application
go run cmd/fx_api/main.go

# Test API endpoints
curl http://localhost:8080/api/v1/users
curl http://localhost:8080/api/v1/auth/login
```

## Advantages over Manual DI

### 1. **Automatic Resolution**
- No need to manually wire dependencies
- Fx resolves dependency graph automatically
- Prevents circular dependencies

### 2. **Lifecycle Integration**
- Built-in startup/shutdown hooks
- Graceful resource management
- Signal handling

### 3. **Error Handling**
- Clear error messages for missing dependencies
- Startup failure prevention
- Resource leak prevention

### 4. **Testing Support**
- Easy mocking with fx.Replace
- Module isolation for testing
- Integration test support

## Best Practices

### 1. **Module Organization**
- One module per architectural layer
- Clear module boundaries
- Minimal cross-module dependencies

### 2. **Constructor Design**
- Use parameter structs for complex dependencies
- Return interfaces when possible
- Handle errors appropriately

### 3. **Lifecycle Management**
- Register cleanup in OnStop hooks
- Use context for cancellation
- Log startup/shutdown events

### 4. **Testing**
- Create test modules for mocking
- Use fx.Replace for test overrides
- Test module composition separately

## Performance Impact

- **Startup Time**: ~10-50ms overhead for dependency resolution
- **Memory Usage**: ~1-2MB for dependency graph
- **Runtime**: Zero overhead after startup
- **Shutdown**: ~100-500ms for graceful cleanup

Fx dependency injection cung cấp một kiến trúc mạnh mẽ, type-safe và dễ maintain cho ứng dụng enterprise Go! 🚀