# Hướng dẫn Tổ chức Source Code - Best Practices

## 🎯 Tổng Quan Cấu Trúc Tối Ưu

### Nguyên tắc chính:
1. **Clean Architecture**: Dependencies chỉ point inward
2. **DDD**: Domain là trung tâm, isolate hoàn toàn
3. **CQRS**: Tách biệt hoàn toàn Commands và Queries
4. **SOLID**: Tuân thủ đầy đủ 5 nguyên tắc

---

## 📁 Cấu Trúc Directory Được Đề Xuất

```
go-mvc/
├── cmd/                                    # Application Entry Points
│   ├── api/                                # HTTP API Server
│   │   └── main.go
│   ├── worker/                             # Background Worker
│   │   └── main.go
│   ├── cli/                                # CLI Tools
│   │   └── main.go
│   └── migrator/                           # Database Migrator
│       └── main.go
│
├── internal/
│   ├── domain/                             # DOMAIN LAYER (Core Business Logic)
│   │   ├── user/                           # User Bounded Context
│   │   │   ├── entity.go                   # User aggregate root
│   │   │   ├── value_objects.go            # Email, Name, Phone, Password
│   │   │   ├── events.go                   # UserCreated, UserUpdated, UserDeleted
│   │   │   ├── repository.go               # Repository interface (port)
│   │   │   ├── specifications.go           # Business rules specifications
│   │   │   └── errors.go                   # Domain-specific errors
│   │   │
│   │   ├── auth/                           # Authentication Bounded Context
│   │   │   ├── entity.go
│   │   │   ├── value_objects.go
│   │   │   ├── events.go
│   │   │   ├── repository.go
│   │   │   └── service.go                  # Domain service (if needed)
│   │   │
│   │   ├── authorization/                  # Authorization Bounded Context
│   │   │   ├── role.go
│   │   │   ├── permission.go
│   │   │   ├── policy.go
│   │   │   └── repository.go
│   │   │
│   │   └── shared/                         # Shared Domain
│   │       ├── value_objects/              # Common value objects
│   │       ├── events/                     # Base domain events
│   │       └── specifications/             # Common specifications
│   │
│   ├── application/                        # APPLICATION LAYER (Use Cases)
│   │   ├── user/                           # User Use Cases
│   │   │   ├── commands/                   # Write Operations (CQRS)
│   │   │   │   ├── create/
│   │   │   │   │   ├── command.go          # CreateUserCommand
│   │   │   │   │   ├── handler.go          # CreateUserCommandHandler
│   │   │   │   │   └── validator.go        # Command validation
│   │   │   │   ├── update/
│   │   │   │   │   ├── command.go
│   │   │   │   │   ├── handler.go
│   │   │   │   │   └── validator.go
│   │   │   │   ├── delete/
│   │   │   │   └── upload_avatar/
│   │   │   │
│   │   │   ├── queries/                    # Read Operations (CQRS)
│   │   │   │   ├── get_by_id/
│   │   │   │   │   ├── query.go            # GetUserByIdQuery
│   │   │   │   │   ├── handler.go          # GetUserByIdQueryHandler
│   │   │   │   │   └── dto.go              # Query response DTO
│   │   │   │   ├── list/
│   │   │   │   │   ├── query.go
│   │   │   │   │   ├── handler.go
│   │   │   │   │   └── dto.go
│   │   │   │   └── search/
│   │   │   │
│   │   │   ├── events/                     # Event Handlers
│   │   │   │   ├── user_created_handler.go
│   │   │   │   ├── user_updated_handler.go
│   │   │   │   └── user_deleted_handler.go
│   │   │   │
│   │   │   └── dto/                        # Shared DTOs for this context
│   │   │       ├── user_response.go
│   │   │       └── user_request.go
│   │   │
│   │   ├── auth/                           # Auth Use Cases
│   │   │   ├── commands/
│   │   │   │   ├── login/
│   │   │   │   ├── register/
│   │   │   │   ├── logout/
│   │   │   │   ├── refresh_token/
│   │   │   │   ├── change_password/
│   │   │   │   ├── reset_password/
│   │   │   │   └── verify_email/
│   │   │   ├── queries/
│   │   │   │   ├── get_profile/
│   │   │   │   └── validate_token/
│   │   │   ├── events/
│   │   │   └── dto/
│   │   │
│   │   ├── authorization/                  # Authorization Use Cases
│   │   │   ├── commands/
│   │   │   ├── queries/
│   │   │   └── dto/
│   │   │
│   │   └── common/                         # Common Application Layer
│   │       ├── interfaces/                 # Application interfaces
│   │       │   ├── command_bus.go
│   │       │   ├── query_bus.go
│   │       │   ├── event_bus.go
│   │       │   └── unit_of_work.go
│   │       ├── behaviors/                  # Cross-cutting behaviors
│   │       │   ├── logging_behavior.go
│   │       │   ├── validation_behavior.go
│   │       │   ├── transaction_behavior.go
│   │       │   └── retry_behavior.go
│   │       └── errors/                     # Application errors
│   │           └── application_errors.go
│   │
│   ├── infrastructure/                     # INFRASTRUCTURE LAYER
│   │   ├── persistence/                    # Data Persistence
│   │   │   ├── postgres/
│   │   │   │   ├── user/                   # User repository implementation
│   │   │   │   │   ├── repository.go       # Implements domain.user.Repository
│   │   │   │   │   ├── mapper.go           # Domain <-> DB model mapper
│   │   │   │   │   └── queries.go          # SQL queries
│   │   │   │   ├── auth/
│   │   │   │   ├── authorization/
│   │   │   │   ├── models/                 # GORM/SQL models
│   │   │   │   │   ├── user_model.go
│   │   │   │   │   ├── role_model.go
│   │   │   │   │   └── permission_model.go
│   │   │   │   ├── migrations/             # Database migrations
│   │   │   │   └── seeds/                  # Database seeds
│   │   │   │
│   │   │   └── redis/                      # Redis implementations
│   │   │       ├── cache_repository.go
│   │   │       └── session_repository.go
│   │   │
│   │   ├── messaging/                      # Message Queue
│   │   │   ├── nats/
│   │   │   │   ├── publisher.go            # Implements EventBus
│   │   │   │   ├── subscriber.go
│   │   │   │   └── connection.go
│   │   │   ├── rabbitmq/
│   │   │   └── kafka/
│   │   │
│   │   ├── external/                       # External Services
│   │   │   ├── storage/                    # File Storage
│   │   │   │   ├── s3/
│   │   │   │   │   └── s3_storage.go       # Implements FileStorageService
│   │   │   │   ├── minio/
│   │   │   │   │   └── minio_storage.go
│   │   │   │   └── local/
│   │   │   │       └── local_storage.go
│   │   │   ├── email/                      # Email Services
│   │   │   │   ├── smtp/
│   │   │   │   │   └── smtp_service.go     # Implements EmailService
│   │   │   │   ├── sendgrid/
│   │   │   │   └── ses/
│   │   │   ├── sms/                        # SMS Services
│   │   │   │   ├── twilio/
│   │   │   │   └── nexmo/
│   │   │   └── notification/               # Push Notifications
│   │   │       ├── fcm/
│   │   │       └── apns/
│   │   │
│   │   ├── jobs/                           # Background Jobs
│   │   │   ├── handlers/                   # Job handlers
│   │   │   ├── scheduler/                  # Job scheduler
│   │   │   └── worker/                     # Worker implementation
│   │   │
│   │   ├── cache/                          # Cache Implementation
│   │   │   ├── redis_cache.go              # Implements CacheService
│   │   │   └── memory_cache.go
│   │   │
│   │   ├── logging/                        # Logging Implementation
│   │   │   ├── zap_logger.go               # Implements Logger interface
│   │   │   └── structured_logger.go
│   │   │
│   │   ├── tracing/                        # Distributed Tracing
│   │   │   ├── jaeger_tracer.go
│   │   │   └── opentelemetry.go
│   │   │
│   │   └── monitoring/                     # Metrics & Monitoring
│   │       ├── prometheus_metrics.go
│   │       └── health_check.go
│   │
│   ├── interfaces/                         # INTERFACE ADAPTERS (Presentation)
│   │   ├── http/                           # HTTP Interface
│   │   │   ├── rest/                       # REST API
│   │   │   │   ├── v1/                     # API version 1
│   │   │   │   │   ├── user/
│   │   │   │   │   │   ├── handler.go      # User HTTP handler
│   │   │   │   │   │   ├── routes.go       # User routes
│   │   │   │   │   │   └── dto.go          # HTTP-specific DTOs
│   │   │   │   │   ├── auth/
│   │   │   │   │   │   ├── handler.go
│   │   │   │   │   │   ├── routes.go
│   │   │   │   │   │   └── dto.go
│   │   │   │   │   └── router.go           # Main v1 router
│   │   │   │   └── v2/                     # API version 2
│   │   │   │
│   │   │   ├── graphql/                    # GraphQL API (if needed)
│   │   │   │   ├── resolvers/
│   │   │   │   ├── schema/
│   │   │   │   └── server.go
│   │   │   │
│   │   │   ├── middleware/                 # HTTP Middlewares
│   │   │   │   ├── auth.go
│   │   │   │   ├── authorization.go
│   │   │   │   ├── cors.go
│   │   │   │   ├── logger.go
│   │   │   │   ├── rate_limit.go
│   │   │   │   ├── recovery.go
│   │   │   │   ├── metrics.go
│   │   │   │   ├── tracing.go
│   │   │   │   └── validation.go
│   │   │   │
│   │   │   ├── responses/                  # Standardized responses
│   │   │   │   ├── success.go
│   │   │   │   ├── error.go
│   │   │   │   └── pagination.go
│   │   │   │
│   │   │   └── server.go                   # HTTP server setup
│   │   │
│   │   ├── grpc/                           # gRPC Interface (if needed)
│   │   │   ├── proto/                      # Protocol buffers
│   │   │   ├── services/                   # gRPC service implementations
│   │   │   └── server.go
│   │   │
│   │   └── cli/                            # CLI Interface
│   │       ├── commands/                   # CLI commands
│   │       └── app.go                      # CLI app setup
│   │
│   └── di/                                 # DEPENDENCY INJECTION
│       ├── container.go                    # Main DI container
│       ├── wire.go                         # Wire code generation (optional)
│       └── modules/                        # DI modules by layer
│           ├── domain_module.go
│           ├── application_module.go
│           ├── infrastructure_module.go
│           ├── interface_module.go
│           └── server_module.go
│
├── pkg/                                    # PUBLIC LIBRARIES (Reusable)
│   ├── errors/                             # Error handling utilities
│   │   ├── errors.go
│   │   └── error_codes.go
│   ├── validator/                          # Input validation
│   │   └── validator.go
│   ├── jwt/                                # JWT utilities
│   │   └── jwt.go
│   ├── pagination/                         # Pagination helpers
│   │   └── pagination.go
│   ├── response/                           # Response helpers
│   │   └── response.go
│   ├── crypto/                             # Encryption utilities
│   │   └── crypto.go
│   └── converter/                          # Type converters
│       └── converter.go
│
├── configs/                                # Configuration Files
│   ├── development.yaml
│   ├── production.yaml
│   ├── testing.yaml
│   └── config_schema.json
│
├── docs/                                   # Documentation
│   ├── API.md
│   ├── ARCHITECTURE.md
│   ├── DEPLOYMENT.md
│   ├── DEVELOPMENT.md
│   ├── REFACTORING_REPORT.md
│   ├── REFACTORING_SUMMARY.md
│   ├── ARCHITECTURE_EXAMPLES.md
│   └── diagrams/                           # Architecture diagrams
│       ├── clean_architecture.png
│       ├── ddd_bounded_contexts.png
│       └── cqrs_flow.png
│
├── scripts/                                # Build & Dev Scripts
│   ├── build.sh
│   ├── test.sh
│   ├── migrate.sh
│   ├── seed.sh
│   └── docker-build.sh
│
├── tests/                                  # Tests
│   ├── unit/                               # Unit tests (by layer)
│   │   ├── domain/
│   │   ├── application/
│   │   └── infrastructure/
│   ├── integration/                        # Integration tests
│   │   ├── api/
│   │   ├── database/
│   │   └── messaging/
│   ├── e2e/                                # End-to-end tests
│   │   └── scenarios/
│   ├── fixtures/                           # Test data
│   └── mocks/                              # Generated mocks
│
├── migrations/                             # Database Migrations
│   ├── postgres/
│   │   ├── 001_create_users_table.up.sql
│   │   ├── 001_create_users_table.down.sql
│   │   └── ...
│   └── redis/
│
├── deployments/                            # Deployment Configs
│   ├── docker/
│   │   ├── Dockerfile
│   │   ├── Dockerfile.dev
│   │   └── docker-compose.yml
│   ├── kubernetes/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── ingress.yaml
│   └── terraform/
│
├── .github/                                # GitHub specific
│   ├── workflows/                          # CI/CD pipelines
│   │   ├── ci.yml
│   │   ├── cd.yml
│   │   └── release.yml
│   └── CODEOWNERS
│
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── .gitignore
├── .env.example
└── LICENSE
```

---

## 🎯 Các Nguyên Tắc Tổ Chức Quan Trọng

### 1. Domain Layer (Core)

#### ✅ Nên làm:
- **Bounded Contexts**: Mỗi domain context (user, auth, authorization) nên độc lập
- **Value Objects**: Immutable, self-validating
- **Entities**: Rich domain models với business logic
- **Domain Events**: Publish sau khi aggregate state changes
- **Repository Interfaces**: Định nghĩa trong domain, implement ở infrastructure
- **Domain Services**: Chỉ khi logic không thuộc về entity nào

#### ❌ Không nên:
- Import bất kỳ thứ gì từ application, infrastructure, hoặc interfaces layer
- Có dependencies đến database, HTTP, external services
- Chứa DTO hoặc data mapping logic
- Chứa framework-specific code

**Ví dụ cấu trúc Domain tốt:**
```go
// internal/domain/user/entity.go
package user

type User struct {
    id        UserID
    email     Email
    name      Name
    password  Password
    events    []DomainEvent
}

func NewUser(email, name, password string) (*User, error) {
    // Business validation
    // Create value objects
    // Raise UserCreated event
}

func (u *User) UpdateProfile(name, phone string) error {
    // Business logic
    // Raise UserUpdated event
}

// internal/domain/user/repository.go
type Repository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id UserID) (*User, error)
    FindByEmail(ctx context.Context, email Email) (*User, error)
}
```

### 2. Application Layer (Use Cases)

#### ✅ Nên làm:
- **Một handler cho mỗi use case**: CreateUserCommandHandler, GetUserByIdQueryHandler
- **Grouping**: commands/, queries/, events/ trong mỗi bounded context
- **DTOs**: Riêng biệt cho commands, queries, và responses
- **Validation**: Validate input ở command/query level
- **Transaction Management**: Use Unit of Work pattern
- **Event Publishing**: Publish events after successful operations

#### ❌ Không nên:
- Chứa domain logic (đó là việc của domain layer)
- Direct database access (dùng repositories)
- HTTP-specific logic (thuộc interfaces layer)

**Ví dụ cấu trúc Application tốt:**
```go
// internal/application/user/commands/create/command.go
package create

type Command struct {
    Email    string
    Name     string
    Password string
}

func (c Command) Validate() error {
    // Input validation
}

// internal/application/user/commands/create/handler.go
type Handler struct {
    userRepo     domain.UserRepository
    eventBus     EventBus
    unitOfWork   UnitOfWork
}

func (h *Handler) Handle(ctx context.Context, cmd Command) (*DTO, error) {
    // 1. Validate command
    // 2. Create domain entity
    // 3. Save via repository
    // 4. Publish events
    // 5. Return DTO
}

// internal/application/user/commands/create/dto.go
type DTO struct {
    ID    string
    Email string
    Name  string
}
```

### 3. Infrastructure Layer

#### ✅ Nên làm:
- **Implement ports**: Repository implementations, service implementations
- **One implementation per file**: Dễ maintain và test
- **Mapper pattern**: Domain <-> DB model conversion
- **Configuration**: Load từ environment hoặc config files
- **Error wrapping**: Wrap infrastructure errors thành domain errors

#### ❌ Không nên:
- Leak infrastructure concerns vào domain
- Return infrastructure-specific types (GORM models, etc.)

**Ví dụ cấu trúc Infrastructure tốt:**
```go
// internal/infrastructure/persistence/postgres/user/repository.go
package user

type postgresRepository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.UserRepository {
    return &postgresRepository{db: db}
}

func (r *postgresRepository) Save(ctx context.Context, user *domain.User) error {
    model := r.toModel(user)  // Domain -> DB model
    return r.db.WithContext(ctx).Save(model).Error
}

func (r *postgresRepository) FindByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
    var model UserModel
    err := r.db.WithContext(ctx).First(&model, "id = ?", id.String()).Error
    if err != nil {
        return nil, err
    }
    return r.toDomain(&model)  // DB model -> Domain
}

// internal/infrastructure/persistence/postgres/user/mapper.go
func (r *postgresRepository) toDomain(model *UserModel) (*domain.User, error) {
    return domain.ReconstructUser(
        model.ID,
        model.Email,
        model.Name,
        model.HashedPassword,
        // ...
    )
}
```

### 4. Interfaces Layer (Presentation)

#### ✅ Nên làm:
- **Versioning**: /api/v1, /api/v2
- **Handler per route**: UserHandler, AuthHandler
- **HTTP-specific DTOs**: Request/Response structs
- **Middleware**: Cross-cutting concerns
- **Error handling**: Convert domain errors to HTTP responses

#### ❌ Không nên:
- Chứa business logic
- Direct repository access

**Ví dụ cấu trúc Interfaces tốt:**
```go
// internal/interfaces/http/rest/v1/user/handler.go
package user

type Handler struct {
    createUserHandler *create.Handler
    getUserHandler    *getbyid.Handler
}

func (h *Handler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        responses.BadRequest(c, err)
        return
    }
    
    cmd := create.Command{
        Email:    req.Email,
        Name:     req.Name,
        Password: req.Password,
    }
    
    dto, err := h.createUserHandler.Handle(c.Request.Context(), cmd)
    if err != nil {
        responses.Error(c, err)
        return
    }
    
    responses.Created(c, toResponse(dto))
}
```

---

## 🔄 CQRS Best Practices

### Command Structure
```
commands/
  └── create_user/
      ├── command.go       # Command struct
      ├── handler.go       # Command handler
      ├── validator.go     # Business validation
      └── dto.go          # Response DTO
```

### Query Structure
```
queries/
  └── get_user_by_id/
      ├── query.go        # Query struct
      ├── handler.go      # Query handler
      └── dto.go         # Result DTO (read model)
```

### Nguyên tắc:
- ✅ Commands modify state, returns void or ID
- ✅ Queries read state, never modify
- ✅ Separate read/write models nếu cần optimize
- ✅ Commands use domain models
- ✅ Queries có thể bypass domain, query DB directly cho performance

---

## 🎯 Dependency Injection Best Practices

### Module Structure
```go
// internal/di/modules/user_module.go
package modules

func NewUserModule() fx.Option {
    return fx.Module("user",
        // Domain
        fx.Provide(user.NewRepository),
        
        // Application - Commands
        fx.Provide(createuser.NewHandler),
        fx.Provide(updateuser.NewHandler),
        
        // Application - Queries
        fx.Provide(getuserbyid.NewHandler),
        fx.Provide(listusers.NewHandler),
        
        // Infrastructure
        fx.Provide(postgres.NewUserRepository),
        
        // Interfaces
        fx.Provide(userhttp.NewHandler),
    )
}
```

---

## 📝 Naming Conventions

### Files
- `entity.go` - Domain entities
- `value_objects.go` - Value objects
- `repository.go` - Repository interface
- `command.go` - Command struct
- `handler.go` - Command/Query handler
- `dto.go` - Data Transfer Objects
- `mapper.go` - Domain <-> DB mapping

### Packages
- Lowercase, no underscores: `user`, `auth`, `createuser`
- Descriptive: `commands`, `queries`, `events`
- Context-based: `user/commands/create`, `user/queries/getbyid`

### Types
- Entities: `User`, `Product`, `Order`
- Value Objects: `Email`, `Money`, `Address`
- Commands: `CreateUserCommand`, `UpdateUserCommand`
- Queries: `GetUserByIDQuery`, `ListUsersQuery`
- Handlers: `CreateUserCommandHandler`, `GetUserByIDQueryHandler`
- DTOs: `UserResponse`, `CreateUserRequest`

---

## 🧪 Testing Structure

```
tests/
├── unit/
│   ├── domain/
│   │   └── user/
│   │       ├── entity_test.go
│   │       └── value_objects_test.go
│   ├── application/
│   │   └── user/
│   │       ├── commands/
│   │       └── queries/
│   └── infrastructure/
│       └── persistence/
│
├── integration/
│   ├── api/
│   ├── database/
│   └── messaging/
│
└── e2e/
    └── scenarios/
```

---

## 🚀 Migration Guide

### Bước 1: Reorganize Domain
1. Di chuyển entities từ `internal/core/domain/` 
2. Tạo bounded contexts: `domain/user/`, `domain/auth/`
3. Tách value objects vào file riêng
4. Tạo repository interfaces trong domain

### Bước 2: Restructure Application
1. Group commands by use case: `create/`, `update/`, `delete/`
2. Group queries by use case: `getbyid/`, `list/`, `search/`
3. Mỗi use case trong folder riêng
4. Tạo DTOs riêng cho từng use case

### Bước 3: Clean Infrastructure
1. Group implementations by technology: `postgres/`, `redis/`, `nats/`
2. Tạo mappers giữa domain và DB models
3. Implement repository interfaces

### Bước 4: Refactor Interfaces
1. Version APIs: `v1/`, `v2/`
2. Group handlers by bounded context
3. Tạo HTTP-specific DTOs

### Bước 5: Update DI
1. Create modules per bounded context
2. Wire dependencies properly
3. Ensure unidirectional dependencies

---

## ✅ Checklist

### Domain Layer
- [ ] No external dependencies
- [ ] Rich domain models
- [ ] Value objects are immutable
- [ ] Repository interfaces in domain
- [ ] Domain events implemented
- [ ] Business rules in domain

### Application Layer
- [ ] One handler per use case
- [ ] Commands separated from queries
- [ ] Input validation in commands/queries
- [ ] DTOs for data transfer
- [ ] Event handlers implemented
- [ ] Transaction management

### Infrastructure Layer
- [ ] Implements domain interfaces
- [ ] No domain logic
- [ ] Proper error handling
- [ ] Configuration management
- [ ] Mappers between layers

### Interface Layer
- [ ] API versioning
- [ ] HTTP-specific DTOs
- [ ] Middleware implemented
- [ ] Error responses standardized
- [ ] Documentation (Swagger/OpenAPI)

### General
- [ ] Tests per layer
- [ ] Documentation updated
- [ ] Dependencies unidirectional
- [ ] SOLID principles followed
- [ ] Code compiles and tests pass

---

## 📚 Tài Liệu Tham Khảo

1. **Clean Architecture** - Robert C. Martin
2. **Domain-Driven Design** - Eric Evans
3. **Implementing Domain-Driven Design** - Vaughn Vernon
4. **CQRS Journey** - Microsoft patterns & practices
5. **Go Best Practices** - Effective Go

---

**Lưu ý**: Đây là best practices được đề xuất. Tùy vào quy mô và yêu cầu project, bạn có thể điều chỉnh cho phù hợp. Quan trọng là giữ nguyên tắc Clean Architecture và separation of concerns.
