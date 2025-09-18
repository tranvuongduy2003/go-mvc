# Enterprise Clean Architecture + DDD Structure

```
my-go-project/
├── cmd/
│   ├── api/                        # REST API server
│   │   └── main.go
│   ├── worker/                     # Background worker
│   │   └── main.go
│   ├── cli/                        # CLI commands
│   │   └── main.go
│   └── migrate/                    # Database migration tool
│       └── main.go
├── internal/
│   ├── core/                       # Core business logic (Domain Layer)
│   │   ├── domain/                 # Domain entities & aggregates
│   │   │   ├── user/
│   │   │   │   ├── entity.go       # User entity
│   │   │   │   ├── repository.go   # User repository interface
│   │   │   │   ├── service.go      # User domain service
│   │   │   │   └── events.go       # Domain events
│   │   │   └── shared/             # Shared domain objects
│   │   │       ├── valueobject/
│   │   │       ├── specification/
│   │   │       └── events/
│   │   └── ports/                  # Ports (interfaces) for external dependencies
│   │       ├── repositories/       # Repository interfaces
│   │       ├── services/           # External service interfaces
│   │       ├── cache/              # Cache interfaces
│   │       └── messaging/          # Message broker interfaces
│   ├── application/                # Application Layer (Use Cases)
│   │   ├── commands/               # Command handlers (CQRS)
│   │   │   ├── user/
│   │   │   └── shared/
│   │   ├── queries/                # Query handlers (CQRS)
│   │   │   ├── user/
│   │   │   └── shared/
│   │   ├── services/               # Application services
│   │   ├── dto/                    # Data Transfer Objects
│   │   ├── validators/             # Business logic validators
│   │   └── events/                 # Event handlers
│   ├── adapters/                   # Adapters (Infrastructure Layer)
│   │   ├── persistence/            # Database implementations
│   │   │   ├── postgres/
│   │   │   │   ├── repositories/
│   │   │   │   ├── models/
│   │   │   │   └── migrations/
│   │   │   └── redis/
│   │   ├── messaging/              # Message broker implementations
│   │   │   ├── nats/
│   │   │   └── kafka/
│   │   ├── external/               # External API clients
│   │   │   ├── payment/
│   │   │   └── notification/
│   │   ├── cache/                  # Cache implementations
│   │   └── monitoring/             # Monitoring implementations
│   ├── handlers/                   # Interface Layer (Presentation)
│   │   ├── http/                   # HTTP handlers
│   │   │   ├── rest/               # REST API handlers
│   │   │   │   ├── v1/
│   │   │   │   └── v2/
│   │   │   ├── middleware/
│   │   │   ├── responses/
│   │   │   └── validators/
│   │   ├── grpc/                   # gRPC handlers
│   │   ├── graphql/                # GraphQL handlers
│   │   └── cli/                    # CLI handlers
│   └── shared/                     # Shared infrastructure
│       ├── config/                 # Configuration management
│       ├── logger/                 # Logging infrastructure
│       ├── database/               # Database connections
│       ├── metrics/                # Metrics collection
│       ├── tracing/                # Distributed tracing
│       ├── security/               # Security utilities
│       ├── middleware/             # Common middleware
│       └── utils/                  # Utility functions
├── pkg/                           # Public reusable packages
│   ├── errors/                     # Error handling
│   ├── validator/                  # Validation utilities
│   ├── pagination/                 # Pagination utilities
│   ├── response/                   # Standard API responses
│   ├── jwt/                        # JWT utilities
│   ├── crypto/                     # Cryptography utilities
│   ├── client/                     # HTTP client utilities
│   └── testing/                    # Testing utilities
├── api/                           # API specifications
│   ├── openapi/                    # OpenAPI/Swagger specs
│   ├── proto/                      # Protocol Buffer definitions
│   └── graphql/                    # GraphQL schemas
├── web/                           # Web assets (if needed)
│   ├── static/
│   └── templates/
├── configs/                       # Configuration files
│   ├── local.yaml
│   ├── development.yaml
│   ├── staging.yaml
│   ├── production.yaml
│   └── docker.yaml
├── deployments/                   # Deployment configurations
│   ├── docker/
│   ├── k8s/
│   └── helm/
├── scripts/                       # Build and deployment scripts
│   ├── build.sh
│   ├── test.sh
│   ├── deploy.sh
│   └── db/
├── docs/                          # Documentation
│   ├── api/                        # API documentation
│   ├── architecture/               # Architecture documentation
│   └── deployment/                 # Deployment guides
├── tests/                         # Test files
│   ├── unit/
│   ├── integration/
│   ├── e2e/
│   ├── fixtures/
│   └── mocks/
├── tools/                         # Development tools
│   └── migrate/
├── .github/                       # GitHub workflows
│   └── workflows/
├── go.mod
├── go.sum
├── go.work                        # Go workspace (if needed)
├── Makefile
├── docker-compose.yml
├── docker-compose.prod.yml
├── Dockerfile
├── Dockerfile.prod
├── .golangci.yml
├── .gitignore
├── .env.example
└── README.md
```

## Giải thích từng layer:

### 1. Domain Layer (`internal/domain/`)
- **Entities**: Các đối tượng nghiệp vụ core, có identity duy nhất
- **Value Objects**: Các đối tượng không có identity, chỉ có value
- **Repository Interfaces**: Định nghĩa contract để truy cập data
- **Domain Services**: Logic nghiệp vụ không thuộc về entity nào cụ thể

### 2. Use Case Layer (`internal/usecase/`)
- Chứa application logic, orchestrate các domain objects
- Implement business workflows
- Không phụ thuộc vào infrastructure

### 3. Infrastructure Layer (`internal/infrastructure/`)
- Database implementations
- External APIs
- File system
- Logging, configuration

### 4. Interface Layer (`internal/interface/`)
- HTTP controllers
- GraphQL resolvers
- CLI commands
- Presenters để format output

## Tham khảo các dependencies:
# Core Dependencies
go get github.com/gin-gonic/gin@v1.9.1                 # HTTP framework
go get github.com/spf13/viper@v1.17.0                  # Configuration
go get go.uber.org/zap@v1.26.0                         # Structured logging
go get go.uber.org/fx@v1.20.0                          # Dependency injection framework

# Database & ORM
go get gorm.io/gorm@v1.25.5                           # ORM
go get gorm.io/driver/postgres@v1.5.4                 # PostgreSQL driver
go get github.com/golang-migrate/migrate/v4@v4.16.2   # Database migrations
go get go.uber.org/multierr@v1.11.0                   # Multiple error handling

# Validation & Serialization
go get github.com/go-playground/validator/v10@v10.16.0 # Validation
go get github.com/go-playground/locales@v0.14.1       # Localization for validation
go get github.com/go-playground/universal-translator@v0.18.1 # Translation
go get github.com/shopspring/decimal@v1.3.1           # Decimal arithmetic

# Security & Authentication
go get github.com/golang-jwt/jwt/v5@v5.1.0             # JWT
go get golang.org/x/crypto@v0.17.0                     # Cryptography
go get github.com/google/uuid@v1.4.0                   # UUID generation

# Caching & Redis
go get github.com/redis/go-redis/v9@v9.3.0             # Redis client
go get github.com/patrickmn/go-cache@v2.1.0            # In-memory cache

# HTTP Client & External APIs
go get github.com/go-resty/resty/v2@v2.10.0            # HTTP client
go get github.com/hashicorp/go-retryablehttp@v0.7.5    # Retryable HTTP client

# Monitoring & Observability
go get github.com/prometheus/client_golang@v1.17.0     # Prometheus metrics
go get go.opentelemetry.io/otel@v1.21.0               # OpenTelemetry
go get go.opentelemetry.io/otel/trace@v1.21.0         # Tracing
go get go.opentelemetry.io/otel/metric@v1.21.0        # Metrics
go get github.com/uptrace/opentelemetry-go-extra/otelgorm@v0.2.3 # GORM tracing

# Message Queue & Event Streaming
go get github.com/nats-io/nats.go@v1.31.0              # NATS messaging
go get github.com/hibiken/asynq@v0.24.1                # Background job processing

# Testing
go get github.com/stretchr/testify@v1.8.4              # Testing toolkit
go get github.com/testcontainers/testcontainers-go@v0.26.0 # Integration testing
go get github.com/golang/mock@v1.6.0                   # Mock generation
go get github.com/DATA-DOG/go-sqlmock@v1.5.0          # SQL mocking

# Documentation
go get github.com/swaggo/swag@v1.16.2                  # Swagger generation
go get github.com/swaggo/gin-swagger@v1.6.0            # Gin Swagger middleware
go get github.com/swaggo/files@v1.0.1                  # Swagger files

# Utilities
go get github.com/joho/godotenv@v1.4.0                 # Environment variables
go get github.com/sony/gobreaker@v0.5.0                # Circuit breaker
go get github.com/golang/groupcache@v0.0.0-20210331224755-41bb18bfe9da # Group cache
go get github.com/robfig/cron/v3@v3.0.1                # Cron jobs

# Development tools
go install github.com/cosmtrek/air@v1.49.0             # Hot reload
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2 # Linter
go install github.com/swaggo/swag/cmd/swag@v1.16.2     # Swagger CLI
go install golang.org/x/tools/cmd/goimports@latest     # Import formatter
go install github.com/golang/mock/mockgen@v1.6.0       # Mock generator


// internal/shared/config/config.go
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// AppConfig contains all application configuration
type AppConfig struct {
	App        App        `mapstructure:"app"`
	Server     Server     `mapstructure:"server"`
	Database   Database   `mapstructure:"database"`
	Redis      Redis      `mapstructure:"redis"`
	Logger     Logger     `mapstructure:"logger"`
	JWT        JWT        `mapstructure:"jwt"`
	Metrics    Metrics    `mapstructure:"metrics"`
	Tracing    Tracing    `mapstructure:"tracing"`
	RateLimit  RateLimit  `mapstructure:"rate_limit"`
	Messaging  Messaging  `mapstructure:"messaging"`
	External   External   `mapstructure:"external"`
	Feature    Feature    `mapstructure:"feature"`
}

type App struct {
	Name         string `mapstructure:"name"`
	Version      string `mapstructure:"version"`
	Environment  string `mapstructure:"environment"`
	Debug        bool   `mapstructure:"debug"`
	Timezone     string `mapstructure:"timezone"`
	GracefulStop time.Duration `mapstructure:"graceful_stop"`
}

type Server struct {
	HTTP ServerHTTP `mapstructure:"http"`
	GRPC ServerGRPC `mapstructure:"grpc"`
}

type ServerHTTP struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	MaxHeaderBytes  int           `mapstructure:"max_header_bytes"`
	TLS             TLS           `mapstructure:"tls"`
	CORS            CORS          `mapstructure:"cors"`
}

type ServerGRPC struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type TLS struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

type CORS struct {
	Enabled        bool     `mapstructure:"enabled"`
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
	ExposedHeaders []string `mapstructure:"exposed_headers"`
	MaxAge         int      `mapstructure:"max_age"`
}

type Database struct {
	Primary DatabaseConnection `mapstructure:"primary"`
	Replica DatabaseConnection `mapstructure:"replica"`
}

type DatabaseConnection struct {
	Driver          string        `mapstructure:"driver"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	LogLevel        string        `mapstructure:"log_level"`
}

type Redis struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type Logger struct {
	Level       string `mapstructure:"level"`
	Encoding    string `mapstructure:"encoding"`
	OutputPaths []string `mapstructure:"output_paths"`
	ErrorPaths  []string `mapstructure:"error_paths"`
	Development bool   `mapstructure:"development"`
}

type JWT struct {
	Secret           string        `mapstructure:"secret"`
	AccessExpiry     time.Duration `mapstructure:"access_expiry"`
	RefreshExpiry    time.Duration `mapstructure:"refresh_expiry"`
	Issuer           string        `mapstructure:"issuer"`
	Audience         string        `mapstructure:"audience"`
}

type Metrics struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
	Port    int    `mapstructure:"port"`
}

type Tracing struct {
	Enabled     bool    `mapstructure:"enabled"`
	ServiceName string  `mapstructure:"service_name"`
	Endpoint    string  `mapstructure:"endpoint"`
	SampleRate  float64 `mapstructure:"sample_rate"`
}

type RateLimit struct {
	Enabled bool `mapstructure:"enabled"`
	RPS     int  `mapstructure:"rps"`
	Burst   int  `mapstructure:"burst"`
}

type Messaging struct {
	NATS NATSConfig `mapstructure:"nats"`
}

type NATSConfig struct {
	URL             string        `mapstructure:"url"`
	MaxReconnects   int           `mapstructure:"max_reconnects"`
	ReconnectWait   time.Duration `mapstructure:"reconnect_wait"`
	Timeout         time.Duration `mapstructure:"timeout"`
	DrainTimeout    time.Duration `mapstructure:"drain_timeout"`
}

type External struct {
	PaymentService PaymentServiceConfig `mapstructure:"payment_service"`
	EmailService   EmailServiceConfig   `mapstructure:"email_service"`
}

type PaymentServiceConfig struct {
	BaseURL    string        `mapstructure:"base_url"`
	APIKey     string        `mapstructure:"api_key"`
	Timeout    time.Duration `mapstructure:"timeout"`
	MaxRetries int           `mapstructure:"max_retries"`
}

type EmailServiceConfig struct {
	Provider string `mapstructure:"provider"`
	APIKey   string `mapstructure:"api_key"`
	From     string `mapstructure:"from"`
}

type Feature struct {
	EnableSwagger     bool `mapstructure:"enable_swagger"`
	EnablePprof       bool `mapstructure:"enable_pprof"`
	EnableHealthCheck bool `mapstructure:"enable_health_check"`
	EnableMetrics     bool `mapstructure:"enable_metrics"`
}

// GetDSN returns database connection string
func (db *DatabaseConnection) GetDSN() string {
	switch db.Driver {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
			db.Host, db.Port, db.Username, db.Password, db.Database, db.SSLMode)
	default:
		return ""
	}
}

// LoadConfig loads configuration from various sources
func LoadConfig(configPath string) (*AppConfig, error) {
	v := viper.New()
	
	// Set config file
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./configs")
		v.AddConfigPath("../configs")
	}
	
	// Set defaults
	setDefaults(v)
	
	// Read environment variables
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	
	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}
	
	var config AppConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}
	
	// Validate config
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "my-enterprise-app")
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.debug", false)
	v.SetDefault("app.timezone", "UTC")
	v.SetDefault("app.graceful_stop", "30s")
	
	// Server defaults
	v.SetDefault("server.http.host", "0.0.0.0")
	v.SetDefault("server.http.port", 8080)
	v.SetDefault("server.http.read_timeout", "30s")
	v.SetDefault("server.http.write_timeout", "30s")
	v.SetDefault("server.http.idle_timeout", "120s")
	v.SetDefault("server.http.max_header_bytes", 1048576)
	
	// CORS defaults
	v.SetDefault("server.http.cors.enabled", true)
	v.SetDefault("server.http.cors.allowed_origins", []string{"*"})
	v.SetDefault("server.http.cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("server.http.cors.allowed_headers", []string{"*"})
	v.SetDefault("server.http.cors.max_age", 86400)
	
	// Database defaults
	v.SetDefault("database.primary.driver", "postgres")
	v.SetDefault("database.primary.host", "localhost")
	v.SetDefault("database.primary.port", 5432)
	v.SetDefault("database.primary.username", "postgres")
	v.SetDefault("database.primary.password", "postgres")
	v.SetDefault("database.primary.database", "myapp")
	v.SetDefault("database.primary.ssl_mode", "disable")
	v.SetDefault("database.primary.max_open_conns", 25)
	v.SetDefault("database.primary.max_idle_conns", 25)
	v.SetDefault("database.primary.conn_max_lifetime", "5m")
	v.SetDefault("database.primary.conn_max_idle_time", "5m")
	v.SetDefault("database.primary.log_level", "warn")
	
	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.pool_size", 10)
	v.SetDefault("redis.min_idle_conns", 5)
	v.SetDefault("redis.dial_timeout", "5s")
	v.SetDefault("redis.read_timeout", "3s")
	v.SetDefault("redis.write_timeout", "3s")
	v.SetDefault("redis.idle_timeout", "5m")
	
	// Logger defaults
	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.encoding", "json")
	v.SetDefault("logger.output_paths", []string{"stdout"})
	v.SetDefault("logger.error_paths", []string{"stderr"})
	v.SetDefault("logger.development", false)
	
	// JWT defaults
	v.SetDefault("jwt.access_expiry", "15m")
	v.SetDefault("jwt.refresh_expiry", "7d")
	v.SetDefault("jwt.issuer", "my-enterprise-app")
	v.SetDefault("jwt.audience", "my-enterprise-app")
	
	// Metrics defaults
	v.SetDefault("metrics.enabled", true)
	v.SetDefault("metrics.path", "/metrics")
	v.SetDefault("metrics.port", 9090)
	
	// Tracing defaults
	v.SetDefault("tracing.enabled", false)
	v.SetDefault("tracing.service_name", "my-enterprise-app")
	v.SetDefault("tracing.sample_rate", 0.1)
	
	// Rate limiting defaults
	v.SetDefault("rate_limit.enabled", true)
	v.SetDefault("rate_limit.rps", 100)
	v.SetDefault("rate_limit.burst", 200)
	
	// Feature flags defaults
	v.SetDefault("feature.enable_swagger", true)
	v.SetDefault("feature.enable_pprof", false)
	v.SetDefault("feature.enable_health_check", true)
	v.SetDefault("feature.enable_metrics", true)
}

// validateConfig validates the configuration
func validateConfig(config *AppConfig) error {
	if config.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}
	
	if config.Server.HTTP.Port <= 0 || config.Server.HTTP.Port > 65535 {
		return fmt.Errorf("server.http.port must be between 1 and 65535")
	}
	
	if config.JWT.Secret == "" && config.App.Environment == "production" {
		return fmt.Errorf("jwt.secret is required in production")
	}
	
	return nil
}

