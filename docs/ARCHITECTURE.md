# Architecture Documentation

## 📋 Table of Contents
- [Overview](#overview)
- [Clean Architecture Layers](#clean-architecture-layers)
- [Detailed Directory Structure](#detailed-directory-structure)
- [Design Patterns](#design-patterns)
- [Dependency Flow](#dependency-flow)
- [Data Flow](#data-flow)

## 🏛️ Overview

This Go MVC application follows **Clean Architecture** principles, ensuring:
- **Independence**: Business logic is independent of frameworks, UI, and databases
- **Testability**: Business logic can be tested without external dependencies
- **Maintainability**: Code is organized in layers with clear responsibilities
- **Scalability**: Architecture supports horizontal and vertical scaling

## 🔄 Clean Architecture Layers

```
┌─────────────────────────────────────────┐
│              Frameworks & Drivers        │  ← External
│  ┌─────────────────────────────────────┐ │
│  │         Interface Adapters          │ │  ← Adapters
│  │  ┌─────────────────────────────────┐│ │
│  │  │        Application Business      ││ │  ← Use Cases
│  │  │  ┌─────────────────────────────┐││ │
│  │  │  │      Enterprise Business    │││ │  ← Domain
│  │  │  │         Rules (Domain)      │││ │
│  │  │  └─────────────────────────────┘││ │
│  │  └─────────────────────────────────┘│ │
│  └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

### 1. **Domain Layer** (innermost)
- Pure business logic
- No external dependencies
- Domain entities and value objects

### 2. **Application Layer**
- Use cases and business workflows
- Orchestrates domain objects
- Defines interfaces for external layers

### 3. **Interface Adapters**
- Controllers, presenters, gateways
- Converts data between layers
- Implements application interfaces

### 4. **Frameworks & Drivers** (outermost)
- Web frameworks, databases, external APIs
- Implementation details
- Can be easily replaced

## 📁 Detailed Directory Structure

### Root Level Files

```
├── go.mod                    # Go module definition and dependencies
├── go.sum                    # Dependency checksums for security
├── Makefile                  # Build automation and development commands
├── README.md                 # Project documentation and quick start
├── docker-compose.yml        # Multi-container Docker application definition
├── Dockerfile               # Production Docker image definition
├── Dockerfile.dev           # Development Docker image with debugging tools
├── .env                     # Environment variables for local development
├── .env.example             # Template for environment variables
├── .gitignore               # Git ignore patterns
└── guide.md                 # Additional project guidance
```

### `/cmd` - Application Entry Points

```bash
cmd/
├── main.go                   # 🎯 Main HTTP server entry point
├── cli/
│   └── main.go              # 🔧 Command-line interface application
├── migrate/
│   └── main.go              # 🗃️ Database migration runner
└── worker/
    └── main.go              # ⚡ Background job processor
```

**Purpose**: Contains the main applications for this project. Each subdirectory is an executable.

**Details**:
- **`main.go`**: Primary HTTP server with Gin router, middleware, and FX dependency injection
- **`cli/main.go`**: Administrative commands using Cobra CLI framework
- **`migrate/main.go`**: Database schema migration using golang-migrate/migrate v4.19.0+ with timestamp-based versioning
- **`worker/main.go`**: Background job processor for async tasks (email, notifications, etc.)

### `/internal` - Private Application Code

The `/internal` directory contains private application code that cannot be imported by other applications.

#### `/internal/core` - Domain Layer (Clean Architecture Core)

```bash
internal/core/
├── domain/                   # 🏛️ Domain entities and business logic
│   ├── shared/              # Common domain constructs
│   │   ├── events/
│   │   │   └── events.go    # Domain event definitions
│   │   ├── specification/   # Business rule specifications
│   │   └── valueobject/
│   │       └── valueobject.go # Immutable value objects
│   └── user/                # User domain aggregate
└── ports/                   # 🔌 Interface definitions (Dependency Inversion)
    ├── cache/               # Cache interface contracts
    ├── messaging/           # Message bus interface contracts
    ├── repositories/        # Data access interface contracts
    └── services/            # External service interface contracts
```

**Purpose**: Contains the core business logic and domain model.

**Key Files**:
- **`events.go`**: Domain events (UserCreated, UserUpdated, etc.)
- **`valueobject.go`**: Immutable objects (Email, UserID, Money, etc.)
- **`ports/`**: Interfaces that define contracts for external dependencies

#### `/internal/application` - Application Layer (Use Cases)

```bash
internal/application/
├── commands/                 # 🎯 Write operations (CQRS Commands)
│   └── shared/              # Common command structures
├── queries/                 # 📊 Read operations (CQRS Queries)
│   └── shared/              # Common query structures
├── dto/                     # 📦 Data Transfer Objects
├── services/                # 🔧 Application services (orchestration)
├── events/                  # 📨 Application event handlers
└── validators/              # ✅ Input validation logic
```

**Purpose**: Implements use cases and business workflows. Orchestrates domain objects.

**Key Concepts**:
- **Commands**: Handle write operations (CreateUser, UpdateProfile, etc.)
- **Queries**: Handle read operations (GetUser, ListUsers, etc.)
- **DTOs**: Data structures for transferring data between layers
- **Services**: Coordinate multiple domain objects and external services

#### `/internal/adapters` - Infrastructure Layer (External Concerns)

```bash
internal/adapters/
├── cache/                    # 💾 Cache implementations
│   ├── cache.go             # Redis cache adapter
│   └── errors.go            # Cache-specific error types
├── external/                # 🌐 External service clients
│   ├── services.go          # Service registry and clients
│   ├── notification/        # Email, SMS, push notification services
│   └── payment/             # Payment gateway integrations
├── messaging/               # 📬 Message queue implementations
│   └── rabbitmq/            # RabbitMQ adapter for async messaging
├── monitoring/              # 📊 Observability implementations
├── persistence/             # 🗄️ Data storage implementations
│   ├── postgres/            # PostgreSQL specific implementations
│   │   ├── migrations/      # Database schema migrations
│   │   ├── models/          # GORM model definitions
│   │   └── repositories/    # Repository implementations
│   └── redis/               # Redis specific implementations
└── repositories/            # 📚 Repository interface implementations
```

**Purpose**: Implements interfaces defined in the core layer. Contains all external dependencies.

**Key Components**:
- **Cache**: Redis-based caching with connection pooling
- **External**: HTTP clients for third-party APIs
- **Messaging**: RabbitMQ for async communication
- **Persistence**: Database repositories with GORM ORM

#### `/internal/fx_modules` - Dependency Injection Modules

```bash
internal/fx_modules/
├── application.go            # 🎯 Application layer DI bindings
├── domain.go                # 🏛️ Domain layer DI bindings
├── handler.go               # 🌐 HTTP handler DI bindings
├── infrastructure.go        # ⚙️ Infrastructure layer DI bindings
└── server.go                # 🚀 Server configuration and startup
```

**Purpose**: Uber FX dependency injection modules for clean separation of concerns.

**Key Files**:
- **`infrastructure.go`**: Database, cache, external service connections
- **`application.go`**: Use case and service bindings
- **`handler.go`**: HTTP routes and middleware setup
- **`server.go`**: Server startup, graceful shutdown, and lifecycle management

#### `/internal/handlers` - Presentation Layer

```bash
internal/handlers/
└── http/                     # 🌐 HTTP-specific handlers
    ├── middleware/           # 🛡️ HTTP middleware components
    │   ├── auth.go          # JWT authentication middleware
    │   ├── authorization.go # RBAC authorization middleware
    │   ├── cors.go          # Cross-Origin Resource Sharing
    │   ├── logger.go        # HTTP request/response logging
    │   ├── manager.go       # Middleware manager and chaining
    │   ├── metrics.go       # Prometheus metrics collection
    │   ├── ratelimit.go     # Rate limiting and throttling
    │   ├── recovery.go      # Panic recovery and error handling
    │   ├── security.go      # Security headers and protection
    │   └── tracing.go       # Distributed tracing with Jaeger
    ├── responses/           # 📄 Response formatting utilities
    ├── rest/                # 🔗 REST API endpoints
    │   └── v1/              # API version 1 endpoints
    └── validators/          # ✅ Request validation logic
```

**Purpose**: HTTP transport layer. Handles HTTP requests and responses.

**Middleware Stack**:
1. **Recovery**: Panic recovery and error handling
2. **Logger**: Structured request/response logging
3. **Tracing**: OpenTelemetry distributed tracing
4. **Metrics**: Prometheus metrics collection
5. **CORS**: Cross-origin resource sharing
6. **Security**: Security headers (CSP, XSS protection)
7. **Rate Limiting**: Request throttling
8. **Authentication**: JWT token validation
9. **Authorization**: RBAC permission checking

#### `/internal/shared` - Shared Infrastructure

```bash
internal/shared/
├── config/
│   └── config.go            # ⚙️ Configuration management (Viper)
├── database/
│   └── database.go          # 🗄️ Database connection and management
├── logger/
│   └── logger.go            # 📝 Structured logging (Zap)
├── metrics/
│   └── metrics.go           # 📊 Prometheus metrics setup
├── security/
│   └── security.go          # 🔒 Security utilities and helpers
├── tracing/
│   └── tracing.go           # 🔍 OpenTelemetry tracing setup
└── utils/
    └── utils.go             # 🛠️ Common utility functions
```

**Purpose**: Shared utilities and infrastructure components used across the application.

**Key Components**:
- **Config**: Environment-based configuration with Viper
- **Database**: GORM connection with connection pooling
- **Logger**: Zap logger with structured JSON output
- **Metrics**: Prometheus metrics registry and collectors
- **Tracing**: OpenTelemetry setup with Jaeger exporter

### `/pkg` - Public Packages

```bash
pkg/
├── crypto/                   # 🔐 Cryptographic utilities
├── errors/
│   └── errors.go            # 🚫 Custom error types and handling
├── jwt/
│   └── jwt.go               # 🎫 JWT token generation and validation
├── pagination/
│   └── pagination.go        # 📄 Pagination utilities
├── response/
│   └── response.go          # 📤 Standardized API response format
└── validator/
    └── validator.go         # ✅ Input validation utilities
```

**Purpose**: Public packages that can be imported by other applications or services.

**Key Packages**:
- **Errors**: Custom error types with error codes and messages
- **JWT**: Token generation, validation, and claims management
- **Pagination**: Offset/limit and cursor-based pagination
- **Response**: Standardized API response format with metadata
- **Validator**: Input validation with custom rules

### `/configs` - Configuration Files

```bash
configs/
├── development.yaml          # 🔧 Development environment config
├── production.yaml           # 🚀 Production environment config
├── prometheus.yml            # 📊 Prometheus scraping configuration
├── redis.conf               # 💾 Redis server configuration
└── grafana/                 # 📈 Grafana monitoring setup
    ├── dashboards/
    │   └── go-mvc-dashboard.json # Pre-built monitoring dashboard
    └── provisioning/
        ├── dashboards/
        │   └── dashboard.yml  # Dashboard auto-provisioning
        └── datasources/
            └── prometheus.yml # Prometheus datasource config
```

**Purpose**: Environment-specific configuration files and monitoring setup.

**Key Files**:
- **`development.yaml`**: Local development settings (debug logging, local URLs)
- **`production.yaml`**: Production settings (optimized logging, production URLs)
- **`prometheus.yml`**: Metrics collection configuration
- **`go-mvc-dashboard.json`**: Grafana dashboard with application metrics

### `/deployments` - Deployment Configurations

```bash
deployments/
└── docker/                   # 🐳 Docker deployment files
```

**Purpose**: Deployment configurations for different environments (Docker, Kubernetes, etc.).

### `/docs` - Documentation

```bash
docs/
├── api/                     # 📚 API documentation
├── DEPENDENCY_INJECTION.md  # 🔌 FX dependency injection guide
├── RBAC_USAGE.md            # 🔐 Role-based access control guide
└── TRACING.md               # 🔍 Distributed tracing setup
```

**Purpose**: Project documentation, API specs, and architectural guides.

### `/scripts` - Utility Scripts

```bash
scripts/
├── db/                      # 🗄️ Database utility scripts
└── init-db.sql             # 🎯 Database initialization script
```

**Purpose**: Build scripts, database initialization, and deployment automation.

### `/api` - API Specifications

```bash
api/
└── openapi/                 # 📋 OpenAPI/Swagger specifications
```

**Purpose**: API documentation and contract definitions.

## 🎯 Design Patterns Used

### 1. **Repository Pattern**
- **Location**: `internal/core/ports/repositories/`, `internal/adapters/persistence/`
- **Purpose**: Abstract data access layer
- **Implementation**: Interface in ports, concrete implementation in adapters

### 2. **Dependency Injection**
- **Framework**: Uber FX
- **Location**: `internal/fx_modules/`
- **Purpose**: Inversion of control and testability

### 3. **CQRS (Command Query Responsibility Segregation)**
- **Location**: `internal/application/commands/`, `internal/application/queries/`
- **Purpose**: Separate read and write operations

### 4. **Domain Events**
- **Location**: `internal/core/domain/shared/events/`
- **Purpose**: Decouple domain logic and side effects

### 5. **Middleware Pattern**
- **Location**: `internal/handlers/http/middleware/`
- **Purpose**: Cross-cutting concerns (auth, logging, metrics)

### 6. **Strategy Pattern**
- **Location**: External service implementations
- **Purpose**: Pluggable algorithms and services

## 🔄 Dependency Flow

```
Handlers → Application → Domain ← Adapters
    ↓           ↓          ↑         ↑
   HTTP    → Use Cases → Entities → Database
              Events     Rules      Cache
                                   External APIs
```

**Rules**:
- **Domain** has NO dependencies on outer layers
- **Application** depends only on Domain
- **Adapters** depend on Application and Domain
- **Handlers** depend on Application layer only

## 📊 Data Flow

### 1. **Request Flow (Inbound)**
```
HTTP Request → Middleware → Handler → Application Service → Domain Entity → Repository → Database
```

### 2. **Response Flow (Outbound)**
```
Database → Repository → Domain Entity → Application Service → Handler → Middleware → HTTP Response
```

### 3. **Event Flow (Async)**
```
Domain Event → Event Handler → External Service → Message Queue → Worker → Side Effects
```

## 🔧 Configuration Management

### Environment Variables
- **Development**: Loaded from `.env` file
- **Production**: Loaded from environment or config files
- **Priority**: ENV > Config File > Defaults

### Configuration Structure
```yaml
server:
  host: "localhost"
  port: 8080
  timeout: 30s

database:
  host: "localhost"
  port: 5432
  name: "go_mvc_dev"
  user: "postgres"
  password: "postgres"

redis:
  host: "localhost"
  port: 6379
  db: 0

logging:
  level: "info"
  format: "json"

tracing:
  enabled: true
  endpoint: "localhost:4318"
  service_name: "go-mvc"
```

## 🛡️ Security Architecture

### Authentication & Authorization
- **JWT Tokens**: Stateless authentication with refresh tokens
- **Authentication Middleware**: `auth.go` - Token validation and user context
- **Authorization Middleware**: `authorization.go` - RBAC permission and role checking
- **RBAC System**: Complete role-based access control implementation

### Security Measures
- **CORS**: Cross-origin resource sharing protection
- **Rate Limiting**: Request throttling per IP/user
- **Security Headers**: XSS, CSP, HSTS protection
- **Input Validation**: Comprehensive request validation
- **SQL Injection**: Parameterized queries with GORM

## 📈 Performance Considerations

### Database
- **Connection Pooling**: Optimized database connections
- **Query Optimization**: Efficient queries with proper indexes
- **Caching**: Redis for frequently accessed data

### HTTP
- **Keep-Alive**: Persistent connections
- **Compression**: Gzip response compression
- **Static Assets**: Efficient static file serving

### Monitoring
- **Metrics**: Prometheus for system metrics
- **Tracing**: Jaeger for request tracing
- **Profiling**: pprof for performance analysis
- **Health Checks**: Comprehensive health monitoring

This architecture ensures scalability, maintainability, and testability while following industry best practices for enterprise Go applications.