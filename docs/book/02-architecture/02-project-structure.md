# Project Structure Overview

## 📁 Complete Directory Structure

```
go-mvc/
├── 📄 README.md                     # Project ov├── ⚙️ configs/                    # Configuration Files
│   ├── 📄 development.yaml        # Development environment config
│   ├── 📄 production.yaml         # Production environment config
│   ├── 📄 prometheus.yml          # Prometheus scraping configuration
│   ├── 📄 redis.conf              # Redis server configuration
│   ├── 🌐 traefik/                # Traefik reverse proxy configuration
│   │   └── 📄 dynamic.yml         # Dynamic routing configuration
│   └── 📈 grafana/                # Grafana monitoring setup and quick start guide
├── 📄 Makefile                      # Build automation and development commands
├── 📄 go.mod                        # Go module definition and dependencies
├── 📄 go.sum                        # Dependency checksums for security
├── 📄 docker-compose.yml            # Multi-container Docker application
├── 📄 Dockerfile                    # Production Docker image
├── 📄 Dockerfile.dev               # Development Docker image
├── 📄 .env.example                 # Environment variables template
├── 📄 .gitignore                   # Git ignore patterns
├── 📄 guide.md                     # Additional project guidance
│
├── 🎯 cmd/                         # Application Entry Points
│   ├── 📄 main.go                  # Main HTTP server (FX + Gin + Middleware)
│   ├── 🔧 cli/
│   │   └── 📄 main.go              # CLI tool (Cobra commands)
│   ├── 🗃️ migrate/
│   │   └── 📄 main.go              # Database migration runner
│   └── ⚡ worker/
│       └── 📄 main.go              # Background job processor
│
├── 🏗️ internal/                    # Private Application Code
│   ├── 🏛️ domain/                  # Domain Layer (Clean Architecture Core)
│   │   ├── 📄 domain.go            # Domain module (Fx)
│   │   ├── user/                   # User domain aggregate
│   │   │   ├── 📄 user.go          # User entity with business logic
│   │   │   └── 📄 user_repository.go # User repository interface
│   │   ├── auth/                   # Authentication domain
│   │   │   ├── 📄 role.go          # Role entity
│   │   │   ├── 📄 permission.go    # Permission entity
│   │   │   ├── 📄 user_role.go     # User-Role relationship
│   │   │   └── 📄 *_repository.go  # Repository interfaces
│   │   ├── messaging/              # Messaging domain (inbox/outbox pattern)
│   │   │   ├── 📄 message.go       # Message entities
│   │   │   ├── 📄 outbox_message.go # Outbox pattern
│   │   │   ├── 📄 inbox_message.go # Inbox pattern
│   │   │   └── 📄 *_repository.go  # Repository interfaces
│   │   ├── job/                    # Background job domain
│   │   │   ├── 📄 job.go           # Job entity
│   │   │   └── 📄 job_types.go     # Job type constants
│   │   ├── contracts/              # Service interfaces (ports)
│   │   │   ├── 📄 auth_service.go  # Auth service interface
│   │   │   ├── 📄 user_service.go  # User service interface
│   │   │   └── 📄 file_storage_service.go
│   │   ├── repositories/           # Common repository interfaces
│   │   └── shared/                 # Shared domain constructs
│   │       ├── events/             # Domain events
│   │       │   ├── 📄 events.go    # Event definitions
│   │       │   └── 📄 user_events.go
│   │       └── valueobject/        # Value objects
│   │           └── 📄 valueobject.go
│   │
│   ├── 🎯 application/             # Application Layer (Use Cases)
│   │   ├── 📄 application.go       # Application module (Fx)
│   │   ├── commands/               # Write operations (CQRS Commands)
│   │   │   ├── 📄 command.go       # Base command interface
│   │   │   ├── auth/               # Auth commands
│   │   │   │   ├── 📄 login_command.go
│   │   │   │   ├── 📄 register_command.go
│   │   │   │   └── ... (10+ auth commands)
│   │   │   └── user/               # User commands
│   │   │       ├── 📄 create_user_command.go
│   │   │       ├── 📄 update_user_command.go
│   │   │       └── 📄 upload_avatar_command.go
│   │   ├── queries/                # Read operations (CQRS Queries)
│   │   │   ├── 📄 query.go         # Base query interface
│   │   │   ├── auth/               # Auth queries
│   │   │   └── user/               # User queries
│   │   ├── dto/                    # Data Transfer Objects
│   │   │   ├── auth/               # Auth DTOs
│   │   │   └── user/               # User DTOs
│   │   ├── services/               # Application services (orchestration)
│   │   │   ├── 📄 auth_service.go  # Auth service implementation
│   │   │   ├── 📄 authorization_service.go # RBAC service
│   │   │   ├── 📄 user_service.go  # User service implementation
│   │   │   └── messaging/          # Messaging services
│   │   │       ├── 📄 outbox_service.go
│   │   │       └── 📄 inbox_service.go
│   │   ├── event_handlers/         # Application event handlers
│   │   │   └── 📄 user_event_handler.go
│   │   └── validators/             # Input validation logic
│   │       └── user/
│   │           └── 📄 user_validator.go
│   │
│   ├── ⚙️ infrastructure/          # Infrastructure Layer (External Concerns)
│   │   ├── 📄 infrastructure.go    # Infrastructure module (Fx)
│   │   ├── config/                 # Configuration management
│   │   │   └── 📄 config.go        # Viper-based config loader
│   │   ├── database/               # Database connection
│   │   │   └── 📄 database.go      # GORM database manager
│   │   ├── 💾 cache/               # Cache implementations
│   │   │   ├── 📄 cache.go         # Redis cache service
│   │   │   └── 📄 errors.go        # Cache-specific errors
│   │   ├── 🌐 external/            # External service clients
│   │   │   ├── 📄 email_service.go # Email service client
│   │   │   ├── 📄 smtp_service.go  # SMTP email implementation
│   │   │   ├── 📄 file_storage_service.go # MinIO file storage
│   │   │   ├── 📄 push_notification_service.go
│   │   │   └── 📄 sms_service.go
│   │   ├── 📬 messaging/           # Message broker implementations
│   │   │   └── nats/               # NATS message broker
│   │   │       ├── 📄 nats.go      # NATS adapter
│   │   │       └── 📄 deduplicated_nats.go # With deduplication
│   │   ├── 🗄️ persistence/         # Data storage implementations
│   │   │   └── postgres/           # PostgreSQL implementations
│   │   │       ├── models/         # GORM models
│   │   │       │   ├── 📄 user.go
│   │   │       │   ├── 📄 role.go
│   │   │       │   └── 📄 permission.go
│   │   │       ├── repositories/   # Repository implementations
│   │   │       │   ├── 📄 user_repository.go
│   │   │       │   ├── 📄 role_repository.go
│   │   │       │   └── 📄 permission_repository.go
│   │   │       └── messaging/      # Messaging repositories
│   │   │           ├── 📄 outbox_repository.go
│   │   │           └── 📄 inbox_repository.go
│   │   ├── security/               # Security utilities
│   │   │   └── 📄 security.go      # Password hashing, token generation
│   │   ├── logger/                 # Logging
│   │   │   └── 📄 logger.go        # Zap logger wrapper
│   │   ├── tracing/                # Distributed tracing
│   │   │   └── 📄 tracing.go       # OpenTelemetry setup
│   │   ├── metrics/                # Metrics collection
│   │   │   └── 📄 metrics.go       # Prometheus metrics
│   │   ├── jobs/                   # Background job system
│   │   │   ├── scheduler/          # Job scheduler
│   │   │   ├── worker/             # Job worker
│   │   │   ├── redis/              # Redis queue
│   │   │   ├── handlers/           # Job handlers
│   │   │   └── metrics/            # Job metrics
│   │   └── utils/                  # Infrastructure utilities
│   │
│   ├── 🔌 modules/                 # Dependency Injection Modules (Uber FX)
│   │   ├── 📄 user.go              # User module DI
│   │   ├── 📄 auth.go              # Auth module DI
│   │   ├── 📄 job.go               # Job module DI
│   │   └── 📄 messaging.go         # Messaging module DI
│   │
│   └── 🌐 presentation/            # Presentation Layer
│       ├── 📄 presentation.go      # Presentation module (Fx)
│       └── http/                   # HTTP-specific handlers
│           ├── handlers/           # HTTP handlers
│           │   ├── 📄 handler.go   # Handler module (Fx)
│           │   └── v1/             # API v1
│           │       ├── 📄 auth_handler.go
│           │       └── 📄 user_handler.go
│           └── middleware/         # HTTP middleware components
│               ├── 📄 manager.go   # Middleware manager
│               ├── 📄 auth.go      # JWT authentication
│               ├── 📄 authorization.go # RBAC authorization
│               ├── 📄 cors.go      # CORS handling
│               ├── 📄 logger.go    # Request/response logging
│               ├── 📄 metrics.go   # Prometheus metrics
│               ├── 📄 ratelimit.go # Rate limiting
│               ├── 📄 recovery.go  # Panic recovery
│               ├── 📄 security.go  # Security headers
│               ├── 📄 tracing.go   # Distributed tracing
│               └── 📄 idempotency.go # Idempotency key
│   │       ├── 📄 responses/       # Response formatting utilities
│   │       ├── 🔗 rest/            # REST API endpoints
│   │       │   └── v1/             # API version 1 endpoints
│   │       └── ✅ validators/      # Request validation logic
│   │
│   └── 🔧 shared/                  # Shared Infrastructure
│       ├── config/
│       │   └── 📄 config.go        # Configuration management (Viper)
│       ├── database/
│       │   └── 📄 database.go      # Database connection (GORM)
│       ├── logger/
│       │   └── 📄 logger.go        # Structured logging (Zap)
│       ├── metrics/
│       │   └── 📄 metrics.go       # Prometheus metrics setup
│       ├── security/
│       │   └── 📄 security.go      # Security utilities
│       ├── tracing/
│       │   └── 📄 tracing.go       # OpenTelemetry tracing setup
│       └── utils/
│           └── 📄 utils.go         # Common utility functions
│
├── 📦 pkg/                        # Public Packages (Importable)
│   ├── 🔐 crypto/                 # Cryptographic utilities
│   ├── 🚫 errors/
│   │   └── 📄 errors.go           # Custom error types and handling
│   ├── 🎫 jwt/
│   │   └── 📄 jwt.go              # JWT token generation/validation
│   ├── 📄 pagination/
│   │   └── 📄 pagination.go       # Pagination utilities
│   ├── 📤 response/
│   │   └── 📄 response.go         # Standardized API response format
│   └── ✅ validator/
│       └── 📄 validator.go        # Input validation utilities
│
├── ⚙️ configs/                    # Configuration Files
│   ├── 📄 development.yaml        # Development environment config
│   ├── 📄 production.yaml         # Production environment config
│   ├── 📄 prometheus.yml          # Prometheus scraping configuration
│   ├── 📄 redis.conf              # Redis server configuration
│   └── 📈 grafana/               # Grafana monitoring setup
│       ├── dashboards/
│       │   └── 📄 go-mvc-dashboard.json # Pre-built dashboard
│       └── provisioning/
│           ├── dashboards/
│           │   └── 📄 dashboard.yml     # Dashboard auto-provisioning
│           └── datasources/
│               └── 📄 prometheus.yml    # Prometheus datasource
│
├── 🚀 deployments/               # Deployment Configurations
│   └── docker/                   # Docker deployment files
│
├── 📚 docs/                      # Documentation
│   ├── 📄 ARCHITECTURE.md        # Detailed architecture guide
│   ├── 📄 API.md                 # API documentation and examples
│   ├── 📄 DEVELOPMENT.md         # Development setup and guidelines
│   ├── 📄 DEPLOYMENT.md          # Production deployment guide
│   ├── 📄 DEPENDENCY_INJECTION.md # FX dependency injection guide
│   ├── 📄 RBAC_USAGE.md          # Role-based access control
│   ├── 📄 TRACING.md             # Distributed tracing setup
│   └── api/                      # API specifications
│
├── 🔧 scripts/                   # Utility Scripts
│   ├── db/                       # Database utility scripts
│   └── 📄 init-db.sql            # Database initialization
│
└── 📋 api/                       # API Specifications
    └── openapi/                  # OpenAPI/Swagger specifications
```

## 🏛️ Architecture Layers Explained

### 1. **Core Layer** (`internal/core/`)
- **Purpose**: Contains pure business logic and domain models
- **Dependencies**: None (innermost layer)
- **Key Components**:
  - **Domain Entities**: Core business objects with behavior
  - **Value Objects**: Immutable data structures
  - **Domain Events**: Business event definitions
  - **Ports**: Interface definitions for external dependencies

### 2. **Application Layer** (`internal/application/`)
- **Purpose**: Orchestrates business workflows and use cases
- **Dependencies**: Only Core layer
- **Key Components**:
  - **Commands**: Handle write operations (Create, Update, Delete)
  - **Queries**: Handle read operations (Get, List, Search)
  - **Services**: Coordinate domain objects and external services
  - **DTOs**: Data transfer between layers

### 3. **Infrastructure Layer** (`internal/adapters/`)
- **Purpose**: Implements external concerns and technical details
- **Dependencies**: Application and Core layers
- **Key Components**:
  - **Persistence**: Database repositories and models
  - **Cache**: Redis caching implementations
  - **External**: Third-party service clients
  - **Messaging**: Event bus and message queue adapters

### 4. **Presentation Layer** (`internal/handlers/`)
- **Purpose**: Handles HTTP requests and responses
- **Dependencies**: Application layer
- **Key Components**:
  - **REST Handlers**: HTTP endpoint implementations
  - **Middleware**: Cross-cutting concerns (auth, logging, metrics)
  - **Validators**: Request validation
  - **Responses**: Response formatting

## 🔄 Data Flow

```
HTTP Request → Middleware → Handler → Application Service → Domain Entity → Repository → Database
              ↓
         Logging, Metrics,
         Tracing, Security
```

## 🎯 Key Design Patterns

### 1. **Repository Pattern**
- **Interfaces**: `internal/core/ports/repositories/`
- **Implementations**: `internal/adapters/persistence/`

### 2. **Dependency Injection**
- **Framework**: Uber FX
- **Modules**: `internal/di/`

### 3. **CQRS (Command Query Responsibility Segregation)**
- **Commands**: `internal/application/commands/`
- **Queries**: `internal/application/queries/`

### 4. **Middleware Pattern**
- **Location**: `internal/handlers/http/middleware/`
- **Order**: Recovery → Logger → Tracing → Metrics → CORS → Security → Rate Limit → Authentication → Authorization

### 5. **Domain Events**
- **Events**: `internal/core/domain/shared/events/`
- **Handlers**: `internal/application/events/`

## 📊 Observability Stack

### Metrics (Prometheus + Grafana)
- **Collection**: `internal/handlers/http/middleware/metrics.go`
- **Exposition**: `/metrics` endpoint
- **Visualization**: Grafana dashboard at `http://localhost:3000`

### Distributed Tracing (Jaeger)
- **Implementation**: `internal/shared/tracing/tracing.go`
- **Middleware**: `internal/handlers/http/middleware/tracing.go`
- **UI**: Jaeger at `http://localhost:16686`

### Logging (Zap)
- **Configuration**: `internal/shared/logger/logger.go`
- **Structured**: JSON format with contextual fields
- **Integration**: Request/response logging middleware

## 🛡️ Security Features

### Authentication & Authorization
- **JWT Tokens**: `pkg/jwt/jwt.go`
- **Authentication Middleware**: `internal/handlers/http/middleware/auth.go` - Token validation
- **Authorization Middleware**: `internal/handlers/http/middleware/authorization.go` - RBAC permissions
- **Security Headers**: XSS, CSRF, HSTS protection

### Input Validation
- **Request Validation**: `internal/handlers/http/validators/`
- **Domain Validation**: `pkg/validator/validator.go`
- **SQL Injection Prevention**: Parameterized queries with GORM

### Rate Limiting
- **Implementation**: `internal/handlers/http/middleware/ratelimit.go`
- **Strategy**: Token bucket per IP address

## 🔧 Development Tools

### Build Automation
- **Makefile**: Comprehensive development commands
- **Hot Reload**: Air for development
- **Code Quality**: golangci-lint, gosec, gofumpt

### Testing
- **Unit Tests**: Alongside source code
- **Integration Tests**: With test database
- **Mocks**: Generated with mockgen
- **Coverage**: HTML coverage reports

### Database
- **Migrations**: golang-migrate/migrate v4.19.0+ with timestamp-based versioning and up/down SQL files
- **ORM**: GORM with PostgreSQL
- **Connection Pooling**: Optimized configuration

## 🚀 Deployment Options

### Docker
- **Development**: `docker-compose.yml`
- **Production**: Multi-stage Dockerfile with optimizations

### Kubernetes
- **Manifests**: Namespace, ConfigMap, Secrets, Deployments
- **Ingress**: SSL termination and load balancing
- **Scaling**: Horizontal Pod Autoscaler

### Monitoring
- **Prometheus**: Metrics collection and alerting
- **Grafana**: Dashboards and visualization
- **Jaeger**: Distributed tracing analysis

This structure ensures:
- **Maintainability**: Clear separation of concerns
- **Testability**: Dependency injection and mocking
- **Scalability**: Stateless design and caching
- **Observability**: Comprehensive monitoring and tracing
- **Security**: Multiple layers of protection
- **Performance**: Optimized database and caching