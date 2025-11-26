# 🤖 AGENT.md - AI Assistant Guide for Go MVC Project

> **Comprehensive guide for AI IDEs, Code Assistants, and Agents working on this Go MVC project**  
> **Last Updated**: November 26, 2025  
> **Project Version**: 1.0.0  
> **Go Version**: 1.24.5

---

## 📋 Table of Contents
- [Quick Overview](#quick-overview)
- [Architecture & Design Patterns](#architecture--design-patterns)
- [Project Structure](#project-structure)
- [Core Technologies](#core-technologies)
- [Development Workflow](#development-workflow)
- [Code Generation Rules](#code-generation-rules)
- [Best Practices & Conventions](#best-practices--conventions)
- [Common Tasks](#common-tasks)
- [Testing Strategy](#testing-strategy)
- [Deployment & Operations](#deployment--operations)

---

## 🎯 Quick Overview

### What is this project?
A **production-ready Go web application** built with:
- ✅ **Clean Architecture** (4 layers: Domain, Application, Infrastructure, Presentation)
- ✅ **Domain-Driven Design (DDD)** with rich domain models
- ✅ **CQRS Pattern** (Command Query Responsibility Segregation)
- ✅ **Dependency Injection** via Uber FX
- ✅ **Full Observability** (Metrics, Tracing, Logging)

### Key Capabilities
1. **RESTful API** with Gin framework
2. **JWT Authentication** with RBAC authorization
3. **PostgreSQL** with GORM ORM
4. **Redis** caching layer
5. **NATS** messaging with inbox/outbox pattern
6. **Background Jobs** with Redis queue
7. **Email Service** with SMTP
8. **File Storage** with MinIO
9. **Distributed Tracing** with Jaeger
10. **Metrics** with Prometheus & Grafana

### Current Features
- **User Management**: CRUD operations, profile, avatar upload
- **Authentication**: Register, login, JWT tokens, refresh tokens
- **Authorization**: Role-based access control (RBAC) with permissions
- **Password Management**: Change, reset, forgot password flow
- **Email Verification**: Email verification with resend capability
- **Message Deduplication**: Idempotent message processing
- **Background Jobs**: Async task processing
- **Observability**: Full metrics, tracing, structured logging

---

## 🏛️ Architecture & Design Patterns

### Clean Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                    PRESENTATION LAYER                       │
│  internal/presentation/                                     │
│  - HTTP handlers (Gin controllers)                          │
│  - Middleware (auth, logging, metrics, etc.)                │
│  - Routes registration                                       │
└─────────────────────────────────────────────────────────────┘
                            ↓ depends on
┌─────────────────────────────────────────────────────────────┐
│                   APPLICATION LAYER                         │
│  internal/application/                                      │
│  - Commands (write operations - CQRS)                       │
│  - Queries (read operations - CQRS)                         │
│  - Services (orchestration logic)                           │
│  - DTOs (data transfer objects)                             │
│  - Validators (input validation)                            │
└─────────────────────────────────────────────────────────────┘
                            ↓ depends on
┌─────────────────────────────────────────────────────────────┐
│                      DOMAIN LAYER                           │
│  internal/domain/                                           │
│  - Entities (business objects with ID)                      │
│  - Value Objects (immutable data)                           │
│  - Repository Interfaces (persistence contracts)            │
│  - Service Interfaces (contracts/ports)                     │
│  - Domain Events                                             │
│  - Business Rules & Validation                              │
└─────────────────────────────────────────────────────────────┘
                            ↑ implemented by
┌─────────────────────────────────────────────────────────────┐
│                  INFRASTRUCTURE LAYER                       │
│  internal/infrastructure/                                   │
│  - Database (GORM implementation)                           │
│  - Cache (Redis implementation)                             │
│  - Messaging (NATS implementation)                          │
│  - External Services (MinIO, SMTP, etc.)                    │
│  - Persistence (repository implementations)                 │
│  - Security (password hashing, token generation)            │
└─────────────────────────────────────────────────────────────┘
```

### Dependency Rules (CRITICAL!)
1. **Domain Layer**: NO dependencies on other layers (pure business logic)
2. **Application Layer**: Depends ONLY on Domain layer
3. **Infrastructure Layer**: Implements Domain & Application interfaces
4. **Presentation Layer**: Depends on Application & Domain (uses interfaces)

**❌ NEVER**: Application → Infrastructure (always use interfaces!)  
**✅ ALWAYS**: Define interfaces in Domain/Application, implement in Infrastructure

### Design Patterns Used

#### 1. Repository Pattern
```go
// Domain defines interface
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id string) (*User, error)
    // ...
}

// Infrastructure implements
type userRepository struct {
    db *gorm.DB
}
```

#### 2. CQRS (Command Query Responsibility Segregation)
```
Commands (Write):
- internal/application/commands/auth/login_command.go
- internal/application/commands/user/create_user_command.go

Queries (Read):
- internal/application/queries/auth/get_user_profile_query.go
- internal/application/queries/user/list_users_query.go
```

#### 3. Dependency Injection (Uber FX)
```go
// Module-based DI in internal/modules/
var UserModule = fx.Module("user",
    fx.Provide(
        NewUserRepository,
        NewUserService,
        NewCreateUserCommandHandler,
        // ...
    ),
)
```

#### 4. Middleware Chain Pattern
```go
// internal/presentation/http/middleware/
- Auth → Authorization → Rate Limiting → Logging → Handler
```

#### 5. Inbox/Outbox Pattern
```go
// For reliable message processing
- internal/domain/messaging/outbox_message.go
- internal/domain/messaging/inbox_message.go
```

---

## 📁 Project Structure

### High-Level Overview
```
go-mvc/
├── cmd/                        # Application entry points
│   ├── main.go                # HTTP server (Gin + Fx)
│   ├── cli/main.go            # CLI commands (Cobra)
│   ├── worker/main.go         # Background worker
│   └── migrate/main.go        # Database migrations
│
├── internal/                   # Private application code
│   ├── domain/                # Domain layer (entities, interfaces)
│   ├── application/           # Application layer (use cases)
│   ├── infrastructure/        # Infrastructure layer (implementations)
│   ├── presentation/          # Presentation layer (HTTP)
│   └── modules/               # Dependency injection modules
│
├── pkg/                       # Public reusable packages
│   ├── errors/                # Error types
│   ├── jwt/                   # JWT utilities
│   ├── pagination/            # Pagination helpers
│   ├── response/              # HTTP response helpers
│   └── validator/             # Validation utilities
│
├── configs/                   # Configuration files
│   ├── development.yaml       # Dev environment config
│   └── production.yaml        # Prod environment config
│
├── docs/                      # Documentation
│   ├── AGENT.md              # This file (AI assistant guide)
│   ├── ARCHITECTURE.md       # Architecture deep dive
│   ├── AI_API_GENERATION_RULES.md  # API generation rules
│   ├── DEVELOPMENT.md        # Development guide
│   └── ...
│
├── scripts/                   # Build and utility scripts
├── docker-compose.yml         # Docker services
├── Dockerfile                 # Production image
├── Dockerfile.dev             # Development image
├── Makefile                   # Build automation
└── go.mod                     # Go dependencies
```

### Domain Layer (`internal/domain/`)
```
domain/
├── domain.go                  # Domain module (Fx)
├── user/
│   ├── user.go               # User entity
│   └── user_repository.go    # Repository interface
├── auth/
│   ├── role.go               # Role entity
│   ├── permission.go         # Permission entity
│   ├── user_role.go          # User-Role relationship
│   ├── role_repository.go    # Repository interfaces
│   └── permission_repository.go
├── messaging/
│   ├── message.go            # Message entities
│   ├── outbox_message.go     # Outbox pattern
│   ├── inbox_message.go      # Inbox pattern
│   └── *_repository.go       # Repository interfaces
├── job/
│   ├── job.go                # Job entity
│   └── job_types.go          # Job type constants
├── contracts/                 # Service interfaces (ports)
│   ├── auth_service.go       # Auth service interface
│   ├── user_service.go       # User service interface
│   └── file_storage_service.go
└── shared/
    ├── events/               # Domain events
    └── valueobject/          # Value objects
```

### Application Layer (`internal/application/`)
```
application/
├── application.go            # Application module (Fx)
├── commands/                 # Write operations (CQRS)
│   ├── command.go           # Base command interface
│   ├── auth/
│   │   ├── login_command.go
│   │   ├── register_command.go
│   │   ├── refresh_token_command.go
│   │   └── ... (10+ commands)
│   └── user/
│       ├── create_user_command.go
│       ├── update_user_command.go
│       ├── delete_user_command.go
│       └── upload_avatar_command.go
│
├── queries/                  # Read operations (CQRS)
│   ├── query.go             # Base query interface
│   ├── auth/
│   │   ├── get_user_profile_query.go
│   │   └── get_user_permissions_query.go
│   └── user/
│       ├── get_user_query.go
│       └── list_users_query.go
│
├── services/                 # Application services
│   ├── auth_service.go      # Auth orchestration
│   ├── authorization_service.go  # RBAC service
│   ├── user_service.go      # User orchestration
│   └── messaging/
│       ├── outbox_service.go
│       └── inbox_service.go
│
├── dto/                      # Data transfer objects
│   ├── auth/
│   │   ├── login_dto.go
│   │   ├── register_dto.go
│   │   └── token_dto.go
│   └── user/
│       ├── user_dto.go
│       └── update_user_dto.go
│
├── validators/               # Input validation
│   └── user/
│       └── user_validator.go
│
└── event_handlers/           # Domain event handlers
    └── user_event_handler.go
```

### Infrastructure Layer (`internal/infrastructure/`)
```
infrastructure/
├── infrastructure.go         # Infrastructure module (Fx)
├── config/
│   └── config.go            # Configuration loading
├── database/
│   └── database.go          # Database connection
├── cache/
│   └── cache.go             # Redis cache implementation
├── messaging/
│   └── nats/
│       ├── nats.go          # NATS broker
│       └── deduplicated_nats.go
├── external/                 # External service clients
│   ├── email_service.go     # Email sending
│   ├── smtp_service.go      # SMTP client
│   ├── file_storage_service.go  # MinIO client
│   ├── sms_service.go       # SMS service
│   └── push_notification_service.go
├── persistence/
│   └── postgres/
│       ├── models/          # GORM models (DB schema)
│       │   ├── user.go
│       │   ├── role.go
│       │   └── permission.go
│       ├── repositories/    # Repository implementations
│       │   ├── user_repository.go
│       │   ├── role_repository.go
│       │   └── permission_repository.go
│       └── messaging/
│           ├── outbox_repository.go
│           └── inbox_repository.go
├── security/
│   └── security.go          # Password hashing, token gen
├── logger/
│   └── logger.go            # Structured logging (Zap)
├── tracing/
│   └── tracing.go           # OpenTelemetry tracing
├── metrics/
│   └── metrics.go           # Prometheus metrics
├── jobs/                     # Background job system
│   ├── scheduler/
│   ├── worker/
│   ├── redis/               # Redis queue
│   └── handlers/            # Job handlers
└── utils/
    └── utils.go             # Utility functions
```

### Presentation Layer (`internal/presentation/`)
```
presentation/
├── presentation.go           # Presentation module (Fx)
└── http/
    ├── handlers/
    │   ├── handler.go       # Handler module (Fx)
    │   └── v1/             # API v1 handlers
    │       ├── auth_handler.go
    │       └── user_handler.go
    └── middleware/
        ├── manager.go       # Middleware manager
        ├── auth.go         # JWT authentication
        ├── authorization.go # RBAC authorization
        ├── cors.go         # CORS handling
        ├── logger.go       # Request/response logging
        ├── metrics.go      # Prometheus metrics
        ├── ratelimit.go    # Rate limiting
        ├── recovery.go     # Panic recovery
        ├── security.go     # Security headers
        ├── tracing.go      # Distributed tracing
        └── idempotency.go  # Idempotency key handling
```

### Dependency Injection (`internal/modules/`)
```
modules/
├── user.go                   # User module dependencies
├── auth.go                   # Auth module dependencies
├── job.go                    # Job module dependencies
└── messaging.go              # Messaging module dependencies
```

---

## ⚙️ Core Technologies

### Backend Stack
| Technology | Version | Purpose |
|------------|---------|---------|
| **Go** | 1.24.5 | Primary language |
| **Gin** | 1.10.1 | HTTP framework |
| **GORM** | 1.30.0 | ORM for PostgreSQL |
| **Uber Fx** | 1.24.0 | Dependency injection |
| **Zap** | 1.26.0 | Structured logging |
| **Viper** | 1.17.0 | Configuration management |

### Infrastructure
| Technology | Version | Purpose |
|------------|---------|---------|
| **PostgreSQL** | 15+ | Primary database |
| **Redis** | 7+ | Caching & job queue |
| **NATS** | 2.10+ | Message broker |
| **MinIO** | Latest | Object storage (S3-compatible) |
| **MailCatcher** | Latest | Email testing (dev) |

### Observability
| Technology | Version | Purpose |
|------------|---------|---------|
| **Prometheus** | Latest | Metrics collection |
| **Grafana** | Latest | Metrics visualization |
| **Jaeger** | Latest | Distributed tracing |
| **OpenTelemetry** | 1.38.0 | Tracing instrumentation |

### Security
| Technology | Purpose |
|------------|---------|
| **JWT** (golang-jwt/jwt/v5) | Token-based authentication |
| **bcrypt** | Password hashing |
| **validator/v10** | Input validation |

### Development Tools
| Tool | Purpose |
|------|---------|
| **Air** | Hot reload |
| **golangci-lint** | Code linting |
| **gosec** | Security scanning |
| **golang-migrate** | Database migrations |

---

## 🔧 Development Workflow

### Setup Development Environment

```bash
# 1. Clone repository
git clone https://github.com/tranvuongduy2003/go-mvc.git
cd go-mvc

# 2. Install dependencies
go mod download

# 3. Start infrastructure services
make docker-up-db      # PostgreSQL, Redis, NATS, MinIO, MailCatcher

# 4. Run database migrations
make migrate-up

# 5. Start application (with hot reload)
make dev

# Or run without hot reload
make run
```

### Common Make Commands

```bash
# Development
make dev                # Run with hot reload (Air)
make run                # Run without hot reload
make build              # Build binary

# Database
make docker-up-db       # Start database services
make docker-down        # Stop all services
make migrate-up         # Run migrations
make migrate-down       # Rollback migrations
make migrate-create NAME=create_xyz_table  # Create new migration

# Testing
make test               # Run all tests
make test-coverage      # Run tests with coverage report
make test-integration   # Run integration tests

# Code Quality
make lint               # Run linter
make fmt                # Format code
make vet                # Run go vet

# Monitoring
make monitoring         # Start Prometheus + Grafana + Jaeger

# Docker
make docker-build       # Build production image
make docker-run         # Run in Docker

# Cleanup
make clean              # Remove binaries and temp files
```

### File Naming Conventions

#### Entity Files
```
<entity_name>.go            # Single entity
user.go, role.go, permission.go
```

#### Repository Files
```
<entity_name>_repository.go # Repository interface (domain)
user_repository.go, role_repository.go
```

#### Command Files (Application Layer)
```
<action>_<entity>_command.go
create_user_command.go
update_user_command.go
login_command.go
```

#### Query Files (Application Layer)
```
<action>_<entity>_query.go
get_user_query.go
list_users_query.go
```

#### Handler Files (Presentation Layer)
```
<entity>_handler.go
user_handler.go, auth_handler.go
```

#### Service Files
```
<domain>_service.go
auth_service.go, user_service.go
```

#### DTO Files
```
<entity>_dto.go
user_dto.go, login_dto.go
```

### Code Structure Pattern

#### Entity (Domain Layer)
```go
package user

import (
    "github.com/google/uuid"
    "time"
)

// User represents a user entity (Aggregate Root)
type User struct {
    ID                string
    Email             string
    Username          string
    PasswordHash      string
    FirstName         string
    LastName          string
    IsEmailVerified   bool
    EmailVerifiedAt   *time.Time
    CreatedAt         time.Time
    UpdatedAt         time.Time
}

// NewUser creates a new user (Factory method)
func NewUser(email, username, passwordHash string) *User {
    return &User{
        ID:           uuid.New().String(),
        Email:        email,
        Username:     username,
        PasswordHash: passwordHash,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
}

// Business methods
func (u *User) VerifyEmail() {
    u.IsEmailVerified = true
    now := time.Now()
    u.EmailVerifiedAt = &now
    u.UpdatedAt = now
}

func (u *User) UpdateProfile(firstName, lastName string) {
    u.FirstName = firstName
    u.LastName = lastName
    u.UpdatedAt = time.Now()
}
```

#### Repository Interface (Domain Layer)
```go
package user

import "context"

type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id string) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindByUsername(ctx context.Context, username string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, offset, limit int) ([]*User, int64, error)
}
```

#### Command Handler (Application Layer)
```go
package commands

import (
    "context"
    "github.com/tranvuongduy2003/go-mvc/internal/domain/user"
    "github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
)

type CreateUserCommand struct {
    Email     string
    Username  string
    Password  string
    FirstName string
    LastName  string
}

type CreateUserCommandHandler struct {
    userRepo       user.UserRepository
    passwordHasher contracts.PasswordHasher
}

func NewCreateUserCommandHandler(
    userRepo user.UserRepository,
    passwordHasher contracts.PasswordHasher,
) *CreateUserCommandHandler {
    return &CreateUserCommandHandler{
        userRepo:       userRepo,
        passwordHasher: passwordHasher,
    }
}

func (h *CreateUserCommandHandler) Handle(
    ctx context.Context,
    cmd CreateUserCommand,
) (*user.User, error) {
    // 1. Hash password
    hashedPassword, err := h.passwordHasher.Hash(cmd.Password)
    if err != nil {
        return nil, err
    }

    // 2. Create user entity
    newUser := user.NewUser(cmd.Email, cmd.Username, hashedPassword)
    newUser.FirstName = cmd.FirstName
    newUser.LastName = cmd.LastName

    // 3. Persist to database
    if err := h.userRepo.Create(ctx, newUser); err != nil {
        return nil, err
    }

    return newUser, nil
}
```

#### HTTP Handler (Presentation Layer)
```go
package v1

import (
    "github.com/gin-gonic/gin"
    "github.com/tranvuongduy2003/go-mvc/internal/application/commands/user"
    "github.com/tranvuongduy2003/go-mvc/pkg/response"
    "net/http"
)

type UserHandler struct {
    createUserHandler *commands.CreateUserCommandHandler
}

func NewUserHandler(
    createUserHandler *commands.CreateUserCommandHandler,
) *UserHandler {
    return &UserHandler{
        createUserHandler: createUserHandler,
    }
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User creation data"
// @Success 201 {object} response.Response{data=UserResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request", err)
        return
    }

    cmd := commands.CreateUserCommand{
        Email:     req.Email,
        Username:  req.Username,
        Password:  req.Password,
        FirstName: req.FirstName,
        LastName:  req.LastName,
    }

    user, err := h.createUserHandler.Handle(c.Request.Context(), cmd)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to create user", err)
        return
    }

    response.Success(c, http.StatusCreated, "User created successfully", toUserResponse(user))
}
```

---

## 🎨 Code Generation Rules

### When to Generate New API

Use AI to generate complete API when you have:
1. **User Story** or **Feature Request**
2. **Database Schema** requirements
3. **API Endpoints** specification
4. **Business Rules** definition

### Generation Steps (Automated by AI)

#### Step 1: Domain Layer
```
1. Create entity in internal/domain/<entity>/
   - <entity>.go (entity with business logic)
   - <entity>_repository.go (repository interface)

2. Define service interface in internal/domain/contracts/
   - <entity>_service.go (if needed)
```

#### Step 2: Application Layer
```
1. Create commands in internal/application/commands/<entity>/
   - create_<entity>_command.go
   - update_<entity>_command.go
   - delete_<entity>_command.go

2. Create queries in internal/application/queries/<entity>/
   - get_<entity>_query.go
   - list_<entity>s_query.go

3. Create DTOs in internal/application/dto/<entity>/
   - <entity>_dto.go
   - create_<entity>_dto.go
   - update_<entity>_dto.go

4. Create validators in internal/application/validators/<entity>/
   - <entity>_validator.go

5. Create service (if needed) in internal/application/services/
   - <entity>_service.go
```

#### Step 3: Infrastructure Layer
```
1. Create GORM model in internal/infrastructure/persistence/postgres/models/
   - <entity>.go

2. Create repository implementation in internal/infrastructure/persistence/postgres/repositories/
   - <entity>_repository.go

3. Create migration file
   - scripts/migrations/YYYYMMDDHHMMSS_create_<entity>_table.up.sql
   - scripts/migrations/YYYYMMDDHHMMSS_create_<entity>_table.down.sql
```

#### Step 4: Presentation Layer
```
1. Create handler in internal/presentation/http/handlers/v1/
   - <entity>_handler.go

2. Add routes in internal/presentation/presentation.go
   - RegisterRoutes() function
```

#### Step 5: Dependency Injection
```
1. Create or update module in internal/modules/
   - <entity>.go

2. Register module in cmd/main.go
```

### Code Generation Template

See `docs/AI_API_GENERATION_RULES.md` for detailed templates and examples.

### Validation Rules
```go
// Always use validator tags
type CreateUserRequest struct {
    Email     string `json:"email" binding:"required,email"`
    Username  string `json:"username" binding:"required,min=3,max=50"`
    Password  string `json:"password" binding:"required,min=8"`
    FirstName string `json:"first_name" binding:"required"`
    LastName  string `json:"last_name" binding:"required"`
}
```

### Error Handling
```go
// Use custom error types from pkg/errors
import "github.com/tranvuongduy2003/go-mvc/pkg/errors"

// Domain errors
return nil, errors.NewNotFoundError("User not found")
return nil, errors.NewValidationError("Invalid email format")
return nil, errors.NewConflictError("Username already exists")

// HTTP error responses
response.Error(c, http.StatusBadRequest, "Invalid request", err)
response.Error(c, http.StatusNotFound, "User not found", err)
response.Error(c, http.StatusConflict, "Username already exists", err)
```

---

## 📚 Best Practices & Conventions

### 1. Naming Conventions

#### Variables
```go
// Use camelCase for local variables
userID := "123"
firstName := "John"

// Use PascalCase for exported variables
const MaxRetries = 3
var DefaultTimeout = 30 * time.Second
```

#### Functions & Methods
```go
// Exported functions (PascalCase)
func NewUserService() *UserService { }
func (s *UserService) CreateUser() error { }

// Private functions (camelCase)
func validateEmail(email string) bool { }
func hashPassword(password string) string { }
```

#### Interfaces
```go
// End with "er" for behavior interfaces
type Reader interface { }
type Writer interface { }
type Closer interface { }

// Or descriptive names for complex interfaces
type UserRepository interface { }
type AuthService interface { }
```

#### Structs
```go
// PascalCase for exported
type User struct { }
type CreateUserCommand struct { }

// camelCase for private
type userRepository struct { }
```

### 2. Error Handling

```go
// ✅ Always handle errors
result, err := someFunction()
if err != nil {
    return nil, fmt.Errorf("failed to do something: %w", err)
}

// ✅ Use custom error types
if err != nil {
    return nil, errors.NewNotFoundError("user not found")
}

// ❌ Never ignore errors
result, _ := someFunction() // Bad!

// ✅ Return early
if err != nil {
    return err
}
// continue with happy path
```

### 3. Context Usage

```go
// ✅ Always pass context as first parameter
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    // Use ctx for timeouts, cancellation, tracing
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return user, nil
}

// ✅ Propagate context through call chain
func (h *UserHandler) GetUser(c *gin.Context) {
    user, err := h.service.GetUser(c.Request.Context(), id)
    // ...
}
```

### 4. Dependency Injection

```go
// ✅ Use constructor injection
type UserService struct {
    repo   user.UserRepository
    logger *logger.Logger
}

func NewUserService(
    repo user.UserRepository,
    logger *logger.Logger,
) *UserService {
    return &UserService{
        repo:   repo,
        logger: logger,
    }
}

// ✅ Define in Fx module
fx.Provide(
    NewUserService,
)
```

### 5. Logging

```go
import "github.com/tranvuongduy2003/go-mvc/internal/infrastructure/logger"

// ✅ Use structured logging
logger.Info("User created", 
    zap.String("user_id", user.ID),
    zap.String("email", user.Email),
)

logger.Error("Failed to create user",
    zap.Error(err),
    zap.String("email", email),
)

// ❌ Avoid unstructured logging
logger.Info(fmt.Sprintf("User %s created", user.ID)) // Bad!
```

### 6. Testing

```go
// Test file naming: <file>_test.go
// user_service_test.go

func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    repo := &mockUserRepository{}
    service := NewUserService(repo)
    
    // Act
    user, err := service.CreateUser(context.Background(), CreateUserCommand{
        Email: "test@example.com",
    })
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "test@example.com", user.Email)
}
```

### 7. Comments & Documentation

```go
// ✅ Document exported types, functions, methods
// User represents a user in the system.
// It contains authentication credentials and profile information.
type User struct {
    ID    string
    Email string
}

// CreateUser creates a new user account.
// It validates the input, hashes the password, and persists the user.
func (s *UserService) CreateUser(ctx context.Context, cmd CreateUserCommand) (*User, error) {
    // Implementation
}

// ✅ Use Swagger comments for HTTP handlers
// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user account with validation
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User data"
// @Success 201 {object} response.Response{data=UserResponse}
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    // Implementation
}
```

### 8. Configuration

```go
// ✅ Use Viper for configuration
// configs/development.yaml
app:
  name: "go-mvc"
  environment: "development"
  port: 8080

database:
  host: "localhost"
  port: 5432
  name: "go_mvc"

// Access in code
cfg, err := config.LoadConfig("configs/development.yaml")
port := cfg.App.Port
```

### 9. Middleware Usage

```go
// Apply middleware in order
router.Use(
    middleware.RequestIDMiddleware(),      // 1. Request ID
    middleware.SecureHeaders(),            // 2. Security headers
    middleware.CORSMiddleware(),           // 3. CORS
    middleware.RecoveryMiddleware(),       // 4. Panic recovery
    middleware.LoggerMiddleware(),         // 5. Logging
    middleware.MetricsMiddleware(),        // 6. Metrics
    middleware.TracingMiddleware(),        // 7. Tracing
)

// Protected routes
protected := router.Group("/api/v1")
protected.Use(
    authMiddleware.RequireAuth(),          // Authentication
    rbacMiddleware.RequireRole("admin"),   // Authorization
)
```

### 10. Database Transactions

```go
// ✅ Use transactions for multiple operations
func (s *UserService) CreateUserWithRoles(ctx context.Context, cmd CreateUserCommand) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Create user
        if err := tx.Create(&user).Error; err != nil {
            return err
        }
        
        // Assign roles
        for _, roleID := range cmd.RoleIDs {
            userRole := &UserRole{UserID: user.ID, RoleID: roleID}
            if err := tx.Create(&userRole).Error; err != nil {
                return err
            }
        }
        
        return nil
    })
}
```

---

## 🧪 Testing Strategy

### Test Types

#### 1. Unit Tests
```bash
# Run all unit tests
make test

# Run specific package
go test ./internal/application/services/...

# With coverage
make test-coverage
```

#### 2. Integration Tests
```bash
# Run integration tests (requires running services)
make test-integration
```

#### 3. End-to-End Tests
```bash
# Run E2E tests
make test-e2e
```

### Test Structure

```go
// internal/application/services/user_service_test.go
package services_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type mockUserRepository struct {
    mock.Mock
}

func (m *mockUserRepository) Create(ctx context.Context, user *user.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    repo := new(mockUserRepository)
    repo.On("Create", mock.Anything, mock.Anything).Return(nil)
    
    service := NewUserService(repo)
    
    // Act
    user, err := service.CreateUser(context.Background(), CreateUserCommand{
        Email: "test@example.com",
    })
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    repo.AssertExpectations(t)
}
```

### Test Coverage Goals
- **Domain Layer**: 90%+ coverage
- **Application Layer**: 85%+ coverage
- **Infrastructure Layer**: 70%+ coverage
- **Presentation Layer**: 75%+ coverage

---

## 🚀 Common Tasks

### Task 1: Add a New Entity (Complete Flow)

```bash
# Example: Adding a "Product" entity

# 1. Domain Layer
touch internal/domain/product/product.go
touch internal/domain/product/product_repository.go

# 2. Application Layer
mkdir -p internal/application/commands/product
touch internal/application/commands/product/create_product_command.go
touch internal/application/commands/product/update_product_command.go

mkdir -p internal/application/queries/product
touch internal/application/queries/product/get_product_query.go
touch internal/application/queries/product/list_products_query.go

mkdir -p internal/application/dto/product
touch internal/application/dto/product/product_dto.go

# 3. Infrastructure Layer
touch internal/infrastructure/persistence/postgres/models/product.go
touch internal/infrastructure/persistence/postgres/repositories/product_repository.go

# Create migration
make migrate-create NAME=create_products_table

# 4. Presentation Layer
touch internal/presentation/http/handlers/v1/product_handler.go

# 5. Dependency Injection
touch internal/modules/product.go

# 6. Update cmd/main.go to include ProductModule
```

### Task 2: Add a New API Endpoint

```go
// 1. Define handler method in internal/presentation/http/handlers/v1/
func (h *ProductHandler) GetProduct(c *gin.Context) {
    // Implementation
}

// 2. Register route in internal/presentation/presentation.go
func RegisterRoutes(params RouteParams) {
    v1API := params.Router.Group("/api/v1")
    
    products := v1API.Group("/products")
    products.Use(authMiddleware.RequireAuth())
    {
        products.GET("/:id", productHandler.GetProduct)
    }
}
```

### Task 3: Add Database Migration

```bash
# Create migration
make migrate-create NAME=add_status_to_users

# Edit migration files
# scripts/migrations/YYYYMMDDHHMMSS_add_status_to_users.up.sql
ALTER TABLE users ADD COLUMN status VARCHAR(20) DEFAULT 'active';

# scripts/migrations/YYYYMMDDHHMMSS_add_status_to_users.down.sql
ALTER TABLE users DROP COLUMN status;

# Run migration
make migrate-up

# Rollback if needed
make migrate-down
```

### Task 4: Add Background Job

```go
// 1. Define job handler in internal/infrastructure/jobs/handlers/
type SendWelcomeEmailHandler struct {
    emailService *external.EmailService
}

func (h *SendWelcomeEmailHandler) Handle(ctx context.Context, payload []byte) error {
    // Implementation
}

// 2. Register job handler
jobScheduler.RegisterHandler("send_welcome_email", welcomeEmailHandler)

// 3. Enqueue job
jobScheduler.Enqueue(ctx, job.Job{
    Type:    "send_welcome_email",
    Payload: payload,
})
```

### Task 5: Add Middleware

```go
// 1. Create middleware in internal/presentation/http/middleware/
func CustomMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Before request
        
        c.Next()
        
        // After request
    }
}

// 2. Apply middleware
router.Use(CustomMiddleware())

// Or to specific routes
protected := router.Group("/api/v1")
protected.Use(CustomMiddleware())
```

---

## 🔒 Security Best Practices

### 1. Authentication
```go
// JWT token validation
authMiddleware.RequireAuth()

// Public routes (no auth)
router.POST("/auth/login", authHandler.Login)
router.POST("/auth/register", authHandler.Register)
```

### 2. Authorization (RBAC)
```go
// Role-based access
rbacMiddleware.RequireRole("admin")
rbacMiddleware.RequireAnyRole("admin", "editor")

// Permission-based access
rbacMiddleware.RequirePermission("users:create")
rbacMiddleware.RequireAllPermissions("users:read", "users:update")
```

### 3. Input Validation
```go
// Always validate input
type CreateUserRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

// Bind and validate
if err := c.ShouldBindJSON(&req); err != nil {
    response.Error(c, http.StatusBadRequest, "Invalid request", err)
    return
}
```

### 4. Password Security
```go
// Use bcrypt for password hashing
passwordHasher := security.NewPasswordHasher(12) // Cost factor 12
hashedPassword, err := passwordHasher.Hash(password)

// Verify password
isValid := passwordHasher.Compare(hashedPassword, password)
```

### 5. SQL Injection Prevention
```go
// ✅ Use GORM parameterized queries
db.Where("email = ?", email).First(&user)

// ❌ Never use string concatenation
db.Where("email = " + email).First(&user) // Vulnerable!
```

### 6. Rate Limiting
```go
// Global rate limiting
router.Use(middleware.GlobalRateLimitMiddleware(100, 200))

// Per-route rate limiting
router.GET("/api/v1/expensive", 
    middleware.RateLimitMiddleware(10, 20),
    handler.ExpensiveOperation,
)
```

### 7. Security Headers
```go
// Applied via SecureHeaders middleware
- X-Content-Type-Options: nosniff
- X-Frame-Options: DENY
- X-XSS-Protection: 1; mode=block
- Strict-Transport-Security: max-age=31536000
- Content-Security-Policy: default-src 'self'
```

---

## 📊 Monitoring & Observability

### Metrics (Prometheus)
```go
// Access metrics at http://localhost:8080/metrics

// Custom business metrics
import "github.com/tranvuongduy2003/go-mvc/internal/infrastructure/metrics"

metrics.RecordUserCreated()
metrics.RecordLoginAttempt(success)
metrics.RecordAPICall(endpoint, method, statusCode, duration)
```

### Distributed Tracing (Jaeger)
```go
// Traces automatically collected via OpenTelemetry middleware
// View traces at http://localhost:16686

// Add custom spans
import "go.opentelemetry.io/otel"

tracer := otel.Tracer("user-service")
ctx, span := tracer.Start(ctx, "CreateUser")
defer span.End()

// Add attributes
span.SetAttributes(
    attribute.String("user.email", email),
    attribute.String("user.id", userID),
)
```

### Logging
```go
// Use structured logging
logger.Info("User created",
    zap.String("user_id", user.ID),
    zap.String("email", user.Email),
)

// Log levels
logger.Debug("Debug message")
logger.Info("Info message")
logger.Warn("Warning message")
logger.Error("Error message", zap.Error(err))
```

### Health Checks
```bash
# Application health
curl http://localhost:8080/health

# Response
{
    "status": "healthy",
    "timestamp": "2025-11-26T10:00:00Z",
    "services": {
        "database": "healthy",
        "redis": "healthy",
        "nats": "healthy"
    }
}
```

---

## 🚀 Deployment & Operations

### Environment Configuration

#### Development
```yaml
# configs/development.yaml
app:
  environment: development
  debug: true
  port: 8080

database:
  host: localhost
  port: 5432
  name: go_mvc_dev

redis:
  host: localhost
  port: 6379
```

#### Production
```yaml
# configs/production.yaml
app:
  environment: production
  debug: false
  port: 8080

database:
  host: ${DB_HOST}
  port: ${DB_PORT}
  name: ${DB_NAME}
  ssl_mode: require

redis:
  host: ${REDIS_HOST}
  port: ${REDIS_PORT}
  password: ${REDIS_PASSWORD}
```

### Docker Deployment

```bash
# Build production image
make docker-build

# Run with docker-compose
docker-compose up -d

# Check logs
docker-compose logs -f app

# Stop services
docker-compose down
```

### Database Migrations

```bash
# Run all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Rollback to specific version
migrate -path scripts/migrations -database "postgres://..." goto <version>

# Check migration status
migrate -path scripts/migrations -database "postgres://..." version
```

### Monitoring Setup

```bash
# Start full monitoring stack
make monitoring

# Access dashboards
- Grafana: http://localhost:3000 (admin/admin)
- Prometheus: http://localhost:9091
- Jaeger: http://localhost:16686
```

---

## 📚 Additional Resources

### Documentation Files
- **ARCHITECTURE.md**: Deep dive into architecture patterns
- **AI_API_GENERATION_RULES.md**: Complete API generation guide
- **DEVELOPMENT.md**: Development setup and workflows
- **DEPLOYMENT.md**: Production deployment guide
- **API.md**: API documentation and examples
- **MIGRATIONS.md**: Database migration guide
- **BACKGROUND_JOBS.md**: Background job system
- **MESSAGE_DEDUPLICATION.md**: Idempotent message processing
- **NATS_MESSAGING.md**: NATS messaging patterns
- **EMAIL_SERVICE.md**: Email service integration
- **FILE_UPLOAD.md**: File upload and storage
- **RBAC_USAGE.md**: Role-based access control
- **TRACING.md**: Distributed tracing setup

### Quick Reference
- **README.md**: Project overview and quick start
- **Makefile**: All available commands
- **docker-compose.yml**: Infrastructure services

---

## 🎯 Key Principles for AI Agents

### 1. **Always Follow Clean Architecture**
- Respect layer boundaries
- Use interfaces for dependencies
- Domain layer has NO external dependencies
- Application layer depends ONLY on domain
- Infrastructure implements interfaces from domain/application

### 2. **Use CQRS Pattern**
- Write operations → Commands
- Read operations → Queries
- Separate models if needed

### 3. **Apply DDD Principles**
- Rich domain models with business logic
- Value objects for immutable data
- Domain events for business events
- Aggregate roots for consistency boundaries

### 4. **Follow Conventions**
- Naming: entity_action_layer.go
- Structure: domain → application → infrastructure → presentation
- Tests: <file>_test.go
- Migrations: YYYYMMDDHHMMSS_description.up/down.sql

### 5. **Generate Complete Features**
When generating new features, include:
- ✅ Domain entities and interfaces
- ✅ Commands and queries
- ✅ DTOs and validators
- ✅ Repository implementations
- ✅ HTTP handlers and routes
- ✅ Database migrations
- ✅ Dependency injection wiring
- ✅ Unit tests
- ✅ Swagger documentation

### 6. **Error Handling**
- Always handle errors
- Use custom error types
- Return wrapped errors with context
- Log errors with structured data

### 7. **Security First**
- Validate all inputs
- Use prepared statements (GORM handles this)
- Hash passwords with bcrypt
- Implement rate limiting
- Use HTTPS in production
- Apply security headers

### 8. **Performance**
- Use context for cancellation
- Implement caching where appropriate
- Use database indexes
- Paginate large result sets
- Monitor with metrics

### 9. **Testing**
- Write tests for business logic
- Mock external dependencies
- Use table-driven tests
- Aim for high coverage in domain/application layers

### 10. **Documentation**
- Comment exported types and functions
- Use Swagger annotations for APIs
- Keep docs updated with code changes
- Include examples in documentation
- **DO NOT create summary/completion markdown files after implementing features**
- Only create new documentation files if explicitly requested
- Update existing docs (README.md, AGENT.md, CHANGELOG.md) instead
- See [AI Coding Standards](docs/AI_CODING_STANDARDS.md) for code documentation principles

---

## 🚫 AI Work Standards

### IMPORTANT: Do NOT Create These Files

After completing any implementation task, **DO NOT** automatically create:
- ❌ IMPLEMENTATION_COMPLETE.md
- ❌ TASK_SUMMARY.md
- ❌ FEATURE_CHECKLIST.md
- ❌ ARCHITECTURE_VISUALIZATION.md
- ❌ QUICK_REFERENCE.md
- ❌ Any other summary/status markdown files

### What to Do Instead

✅ **DO:**
- Update existing documentation files (README.md, CHANGELOG.md)
- Report completion status directly to user
- Update relevant docs in `docs/` directory if needed
- Only create new .md files when explicitly requested

✅ **Example:**
```
After implementing feature X:
1. Update README.md with new feature (if significant)
2. Add entry to CHANGELOG.md
3. Update API documentation in docs/ (if needed)
4. Tell user "Feature X completed successfully"
```

❌ **DON'T:**
```
After implementing feature X:
1. Create FEATURE_X_COMPLETE.md
2. Create FEATURE_X_SUMMARY.md
3. Create ARCHITECTURE_UPDATE.md
```

**Rule**: Focus on code quality and essential documentation only. Avoid creating verbose summary files that clutter the repository.

---

## 🤝 AI Assistant Checklist

When working on this project, ensure:

- [ ] Understand the current architecture
- [ ] Follow Clean Architecture principles
- [ ] Use CQRS pattern (commands vs queries)
- [ ] Implement dependency injection via Fx
- [ ] Create migrations for database changes
- [ ] Add proper error handling
- [ ] Write structured logs
- [ ] Include input validation
- [ ] Add authentication/authorization as needed
- [ ] Write unit tests
- [ ] Update Swagger documentation
- [ ] Follow naming conventions
- [ ] Add metrics and tracing spans
- [ ] Update dependency injection modules
- [ ] Register routes properly
- [ ] Test locally before committing
- [ ] **DO NOT create summary markdown files after completion**

---

## 📞 Contact & Support

For questions or issues:
- **GitHub**: https://github.com/tranvuongduy2003/go-mvc
- **Email**: tranvuongduy2003@gmail.com
- **Documentation**: See `docs/` folder

---

**Last Updated**: November 26, 2025  
**Maintained By**: Tran Vuong Duy  
**Version**: 1.0.0  
**Go Version**: 1.24.5

---

*This document is designed to help AI assistants, code agents, and developers understand and work with this Go MVC project efficiently. Keep it updated as the project evolves.*
