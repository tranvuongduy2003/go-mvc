# AI API Generation Rules

> ⚠️ **IMPORTANT**: All generated code must follow **[AI Coding Standards](../appendix/AI_CODING_STANDARDS.md)**
> - Use self-documenting code with clear naming
> - NO COMMENTS except for complex algorithms, security, or compliance
> - Code quality over brevity

## 📋 Table of Contents
- [Overview](#overview)
- [Coding Standards](#coding-standards)
- [User Story Template](#user-story-template)
- [API Generation Process](#api-generation-process)
- [Layer-by-Layer Guidelines](#layer-by-layer-guidelines)
- [Code Conventions](#code-conventions)
- [File Structure Templates](#file-structure-templates)
- [Validation Rules](#validation-rules)
- [Error Handling](#error-handling)
- [Testing Guidelines](#testing-guidelines)
- [Examples](#examples)

## 🎯 Overview

This document defines detailed rules for AI to automatically generate a complete API set from User Stories following the Go MVC project's Clean Architecture.

### Goals
- **Automation**: AI only needs User Story to generate complete code
- **Consistency**: Ensure generated code follows current patterns and conventions
- **Quality**: Generated code is production-ready with complete validation, error handling, and tests
- **Architecture**: Strictly follow Clean Architecture with 4 layers

### Automatic Process
1. **Analyze User Story**: Extract required information
2. **Generate Domain Layer**: Create entities, value objects, repository interfaces
3. **Generate Application Layer**: Create commands/queries, DTOs, services, validators
4. **Generate Infrastructure Layer**: Create repository implementations, migrations
5. **Generate Presentation Layer**: Create HTTP handlers, routes, middleware
6. **Update Dependency Injection**: Link all components

## 🤖 Coding Standards

### CRITICAL: All AI-generated code MUST follow these principles

1. **NO COMMENTS** except for:
   - Complex algorithms that cannot be expressed in code
   - Security-sensitive code
   - Compliance/regulatory requirements
   - Workarounds for external library bugs

2. **Self-Documenting Code**:
   - Use descriptive function names
   - Use clear variable names
   - Break complex functions into smaller ones
   - Use guard clauses
   - Use custom types for clarity

3. **See Complete Standards**: [AI Coding Standards](../appendix/AI_CODING_STANDARDS.md)

### Example: Good vs Bad Code

**❌ BAD - Too Many Comments**
```go
// CreateUser creates a new user
func CreateUser(email, password string) (*User, error) {
    // Validate email
    if !isValid(email) {
        return nil, errors.New("invalid email")
    }
    // Hash password
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
    // Create user
    user := &User{Email: email, Password: string(hash)}
    // Save to database
    return user, db.Save(user)
}
```

**✅ GOOD - Self-Documenting**
```go
func CreateUser(email, password string) (*User, error) {
    if err := ValidateEmail(email); err != nil {
        return nil, err
    }
    
    hashedPassword, err := HashPassword(password)
    if err != nil {
        return nil, err
    }
    
    user := &User{
        Email:        email,
        PasswordHash: hashedPassword,
    }
    
    return user, SaveUser(user)
}
```

## 📝 User Story Template

### Required Format

```markdown
## User Story: [Feature Name]

### Business Description
- **Actor**: [Who will use - User, Admin, System, etc.]
- **Action**: [What action - Create, Update, Delete, Get, List, etc.]
- **Object**: [What object - Product, Order, User, etc.]
- **Purpose**: [Purpose/benefit]

### Functional Requirements
- **Inputs**: 
  - field1: type (validation rules)
  - field2: type (validation rules)
  - ...
- **Outputs**:
  - Định nghĩa response structure
- **Business Rules**:
  - Rule 1: [Mô tả business logic]
  - Rule 2: [Mô tả constraints]
  - ...

### Technical Specifications  
- **HTTP Method**: GET/POST/PUT/DELETE
- **Endpoint**: /api/v1/[resource]
- **Authentication**: Required/Optional
- **Authorization**: [Role-based rules]
- **Pagination**: Yes/No (if applicable)

### Database Impact
- **Tables**: [Affected tables]
- **Relationships**: [Foreign keys, joins]
- **Indexes**: [Required indexes]
- **Migration**: [Schema changes needed]

### Validation Rules
- **Required fields**: [List]
- **Format validation**: [Email, phone, etc.]
- **Business validation**: [Unique constraints, etc.]
- **Size limits**: [String lengths, file sizes, etc.]

### Error Scenarios
- **Client Errors (4xx)**:
  - 400: [Specific validation errors]
  - 401: [Authentication scenarios]
  - 403: [Authorization scenarios]
  - 404: [Not found scenarios]
  - 409: [Conflict scenarios]
- **Server Errors (5xx)**:
  - 500: [Internal server errors]

### Performance Requirements
- **Response Time**: [Expected latency]
- **Throughput**: [Expected RPS]
- **Caching**: [Cache strategy if applicable]

### Integration Requirements
- **External APIs**: [Third-party services]
- **Message Queue**: [Async processing]
- **Email/SMS**: [Notification requirements]
```

## 🔄 API Generation Process

### Step 1: Domain Analysis
AI phải phân tích User Story để xác định:

1. **Domain Entity**: 
   - Tên entity chính (vd: Product, Order, User)
   - Value objects cần thiết (vd: Email, Money, ProductName)
   - Business methods và invariants

2. **Repository Interface**:
   - CRUD methods cần thiết
   - Custom query methods
   - Specifications cho complex queries

3. **Domain Events**:
   - Events cần fire (vd: ProductCreated, OrderUpdated)
   - Event payload structure

### Step 2: Application Layer Design
1. **Commands vs Queries**:
   - Write operations → Commands
   - Read operations → Queries

2. **DTOs**:
   - Request DTOs cho input
   - Response DTOs cho output
   - Internal DTOs cho layer communication

3. **Services**:
   - Application services cho orchestration
   - External service interfaces

4. **Validators**:
   - Input validation rules
   - Business validation rules

### Step 3: Infrastructure Implementation
1. **Database Models**:
   - GORM models với correct tags
   - Relationships definition
   - Indexes và constraints

2. **Repository Implementation**:
   - Implement domain repository interfaces
   - Query optimizations
   - Error handling

3. **Migrations**:
   - Database schema changes
   - Data migrations if needed

### Step 4: Presentation Layer
1. **HTTP Handlers**:
   - REST endpoints implementation
   - Request/response mapping
   - Error handling

2. **Routes**:
   - Route definitions
   - Middleware assignments

3. **API Documentation**:
   - Swagger annotations
   - Request/response examples

## 🏗️ Layer-by-Layer Guidelines

### 1. Domain Layer (`internal/domain/`)

#### Entity Creation Rules
```go
// Template for domain entity
package [entity_name]

import (
    "time"
    "errors"
    "github.com/google/uuid"
    "github.com/tranvuongduy2003/go-mvc/internal/domain/shared/events"
)

// [Entity] represents the [entity] aggregate root
type [Entity] struct {
    id        [Entity]ID
    // Other fields as value objects
    createdAt time.Time
    updatedAt time.Time
    version   int64
    events    []events.DomainEvent
}

// Constructor
func New[Entity](...) (*[Entity], error) {
    // Validation logic
    // Business rules enforcement
    // Return entity with events
}

// Business methods
func (e *[Entity]) [BusinessMethod](...) error {
    // Business logic
    // Add domain events
    return nil
}

// Getters (no setters to maintain encapsulation)
func (e *[Entity]) ID() [Entity]ID { return e.id }
func (e *[Entity]) Events() []events.DomainEvent { return e.events }
func (e *[Entity]) ClearEvents() { e.events = []events.DomainEvent{} }
```

#### Value Object Rules
```go
// Template for value objects
type [ValueObject] struct {
    value [type]
}

func New[ValueObject](value [type]) ([ValueObject], error) {
    // Validation logic
    if /* validation fails */ {
        return [ValueObject]{}, errors.New("validation message")
    }
    return [ValueObject]{value: value}, nil
}

func (vo [ValueObject]) String() string {
    return vo.value
}

// Add comparison methods if needed
func (vo [ValueObject]) Equals(other [ValueObject]) bool {
    return vo.value == other.value
}
```

#### Repository Interface Rules
```go
// Template for repository interface
package repositories

import (
    "context"
    "[entity_package]"
)

type [Entity]Repository interface {
    // Standard CRUD
    Create(ctx context.Context, entity *[entity_package].[Entity]) error
    GetByID(ctx context.Context, id [entity_package].[Entity]ID) (*[entity_package].[Entity], error)
    Update(ctx context.Context, entity *[entity_package].[Entity]) error
    Delete(ctx context.Context, id [entity_package].[Entity]ID) error
    
    // Listing with pagination
    List(ctx context.Context, offset, limit int) ([]*[entity_package].[Entity], int64, error)
    
    // Custom queries based on User Story
    // Example: FindByStatus, FindByUserID, etc.
}
```

### 2. Application Layer (`internal/application/`)

#### Command/Query Rules
```go
// Template for Command
package commands

import (
    "context"
    "[domain_package]"
    "github.com/tranvuongduy2003/go-mvc/internal/domain/repositories"
)

// [Action][Entity]Command represents a command to [action] [entity]
type [Action][Entity]Command struct {
    // Command fields
}

// [Action][Entity]CommandHandler handles [action] [entity] command
type [Action][Entity]CommandHandler struct {
    [entity]Repo repositories.[Entity]Repository
    // Other dependencies
}

func New[Action][Entity]CommandHandler(
    [entity]Repo repositories.[Entity]Repository,
    // Other dependencies
) *[Action][Entity]CommandHandler {
    return &[Action][Entity]CommandHandler{
        [entity]Repo: [entity]Repo,
    }
}

func (h *[Action][Entity]CommandHandler) Handle(ctx context.Context, cmd [Action][Entity]Command) ([Response], error) {
    // 1. Validate command (if needed)
    // 2. Load existing entities (for updates)
    // 3. Execute business logic
    // 4. Save changes
    // 5. Publish events
    // 6. Return response
}
```

```go
// Template for Query
package queries

import (
    "context"
    "[dto_package]"
    "github.com/tranvuongduy2003/go-mvc/internal/domain/repositories"
)

// [Action][Entity]Query represents a query to [action] [entity]
type [Action][Entity]Query struct {
    // Query parameters
}

// [Action][Entity]QueryHandler handles [action] [entity] query
type [Action][Entity]QueryHandler struct {
    [entity]Repo repositories.[Entity]Repository
}

func New[Action][Entity]QueryHandler([entity]Repo repositories.[Entity]Repository) *[Action][Entity]QueryHandler {
    return &[Action][Entity]QueryHandler{[entity]Repo: [entity]Repo}
}

func (h *[Action][Entity]QueryHandler) Handle(ctx context.Context, query [Action][Entity]Query) ([Response], error) {
    // 1. Execute query
    // 2. Map to DTOs
    // 3. Return response
}
```

#### DTO Rules
```go
// Request DTO Template
package dto

import (
    "time"
)

type [Action][Entity]Request struct {
    // Fields matching API request
    // Use json tags for serialization
    // Use validate tags for validation
    Field1 string `json:"field1" validate:"required,min=2,max=100"`
    Field2 int    `json:"field2" validate:"required,gt=0"`
}

// Response DTO Template
type [Entity]Response struct {
    ID        string    `json:"id"`
    // Fields for API response
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Mapping methods
func To[Entity]Response(entity *[domain_package].[Entity]) [Entity]Response {
    return [Entity]Response{
        ID: entity.ID().String(),
        // Map other fields
        CreatedAt: entity.CreatedAt(),
        UpdatedAt: entity.UpdatedAt(),
    }
}
```

#### Validator Rules
```go
// Template for validator
package validators

import (
    "[dto_package]"
    apperrors "github.com/tranvuongduy2003/go-mvc/pkg/errors"
)

type I[Entity]Validator interface {
    Validate[Action][Entity]Request(req [dto_package].[Action][Entity]Request) []apperrors.ValidationError
}

type [Entity]Validator struct{}

func New[Entity]Validator() I[Entity]Validator {
    return &[Entity]Validator{}
}

func (v *[Entity]Validator) Validate[Action][Entity]Request(req [dto_package].[Action][Entity]Request) []apperrors.ValidationError {
    var errors []apperrors.ValidationError
    
    // Field validation
    if req.Field1 == "" {
        errors = append(errors, apperrors.ValidationError{
            Field:   "field1",
            Message: "Field1 is required",
        })
    }
    
    // Business validation
    // Check database constraints
    // Validate business rules
    
    return errors
}
```

### 3. Infrastructure Layer (`internal/infrastructure/`)

#### Repository Implementation Rules
```go
// Template for repository implementation
package repositories

import (
    "context"
    "gorm.io/gorm"
    "[domain_package]"
    "[models_package]"
    "github.com/tranvuongduy2003/go-mvc/internal/domain/repositories"
)

type [Entity]Repository struct {
    db *gorm.DB
}

func New[Entity]Repository(db *gorm.DB) repositories.[Entity]Repository {
    return &[Entity]Repository{db: db}
}

func (r *[Entity]Repository) Create(ctx context.Context, entity *[domain_package].[Entity]) error {
    model := [models_package].From[Entity](entity)
    if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
        return err
    }
    return nil
}

func (r *[Entity]Repository) GetByID(ctx context.Context, id [domain_package].[Entity]ID) (*[domain_package].[Entity], error) {
    var model [models_package].[Entity]
    err := r.db.WithContext(ctx).First(&model, "id = ?", id.String()).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, err
    }
    return model.To[Entity](), nil
}

// Implement other methods...
```

#### Database Model Rules
```go
// Template for database model
package models

import (
    "time"
    "[domain_package]"
    "gorm.io/gorm"
)

type [Entity] struct {
    ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
    // Other fields mapped to database columns
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Table name
func ([Entity]) TableName() string {
    return "[table_name]"
}

// Convert from domain entity to model
func From[Entity](entity *[domain_package].[Entity]) *[Entity] {
    return &[Entity]{
        ID: entity.ID().String(),
        // Map other fields
    }
}

// Convert from model to domain entity
func (m *[Entity]) To[Entity]() *[domain_package].[Entity] {
    // Reconstruct domain entity with proper validation
    // Handle potential errors appropriately
}
```

#### Migration Rules
```sql
-- Template for migration
-- migrate create -ext sql -dir internal/infrastructure/persistence/postgres/migrations -seq create_[table_name]_table

-- Create table
CREATE TABLE IF NOT EXISTS [table_name] (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Define columns based on domain model
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_[table_name]_deleted_at ON [table_name](deleted_at);
-- Add other indexes based on query patterns

-- Add comments
COMMENT ON TABLE [table_name] IS '[Description of the table]';
COMMENT ON COLUMN [table_name].id IS '[Description of ID column]';
```

### 4. Presentation Layer (`internal/presentation/http/handlers/v1/`)

#### HTTP Handler Rules
```go
// Template for HTTP handler
package v1

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "[dto_package]"
    "[service_package]"
    "[validator_package]"
    "github.com/tranvuongduy2003/go-mvc/pkg/response"
)

// [Entity]Handler handles HTTP requests for [entity] operations
type [Entity]Handler struct {
    [entity]Service   *[service_package].[Entity]Service
    [entity]Validator [validator_package].I[Entity]Validator
}

func New[Entity]Handler(
    [entity]Service *[service_package].[Entity]Service,
    [entity]Validator [validator_package].I[Entity]Validator,
) *[Entity]Handler {
    return &[Entity]Handler{
        [entity]Service:   [entity]Service,
        [entity]Validator: [entity]Validator,
    }
}

// [Action][Entity] [description]
// @Summary [Summary]
// @Description [Description]
// @Tags [entity]
// @Accept json
// @Produce json
// @Param [param] body dto.[Action][Entity]Request true "[Description]"
// @Success [code] {object} response.APIResponse{data=dto.[Entity]Response}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/[resource] [method]
func (h *[Entity]Handler) [Action][Entity](c *gin.Context) {
    var req [dto_package].[Action][Entity]Request
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, err)
        return
    }

    // Validate request
    if validationErrors := h.[entity]Validator.Validate[Action][Entity]Request(req); len(validationErrors) > 0 {
        response.ValidationError(c, validationErrors)
        return
    }

    // Call service
    result, err := h.[entity]Service.[Action][Entity](c.Request.Context(), req)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, result, "Operation successful")
}
```

## 📏 Code Conventions

### 1. Naming Conventions
- **Files**: snake_case (user_handler.go, create_user_command.go)
- **Packages**: lowercase single word (user, product, order)
- **Types**: PascalCase (UserService, CreateUserCommand)
- **Functions/Methods**: PascalCase for public, camelCase for private
- **Constants**: SCREAMING_SNAKE_CASE
- **Variables**: camelCase

### 2. Import Organization
```go
import (
    // Standard library
    "context"
    "time"
    
    // Third-party packages
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    
    // Internal packages
    "github.com/tranvuongduy2003/go-mvc/internal/domain/user"
    "github.com/tranvuongduy2003/go-mvc/pkg/errors"
)
```

### 3. Error Handling
```go
// Use custom error types
if err != nil {
    return apperrors.NewInternalError("Failed to create user", err)
}

// Handle specific error cases
if user == nil {
    return apperrors.NewNotFoundError("User not found")
}

// Validation errors
if validationErrors := validator.Validate(req); len(validationErrors) > 0 {
    return apperrors.NewValidationError("Validation failed", validationErrors)
}
```

### 4. Logging
```go
import "github.com/tranvuongduy2003/go-mvc/internal/shared/logger"

// Log with context
logger.Info(ctx, "User created successfully", "user_id", userID)
logger.Error(ctx, "Failed to create user", "error", err)
```

## 📁 File Structure Templates

### Complete file structure for một entity mới:

```
internal/
├── domain/
│   ├── [entity]/
│   │   ├── [entity].go              # Domain entity
│   │   ├── [entity]_id.go          # Entity ID value object
│   │   ├── value_objects.go         # Value objects
│   │   └── events.go               # Domain events
│   ├── repositories/
│   │   └── [entity]_repository.go   # Repository interface
│   └── shared/
│       └── events/                  # Shared event interfaces
├── application/
│   ├── commands/
│   │   └── [entity]/
│   │       ├── create_[entity]_command.go
│   │       ├── update_[entity]_command.go
│   │       └── delete_[entity]_command.go
│   ├── queries/
│   │   └── [entity]/
│   │       ├── get_[entity]_query.go
│   │       └── list_[entity]_query.go
│   ├── dto/
│   │   └── [entity]/
│   │       └── [entity]_dto.go          # Request/Response DTOs
│   ├── services/
│   │   └── [entity]_service.go          # Application service
│   └── validators/
│       └── [entity]/
│           └── [entity]_validator.go    # Validators
├── infrastructure/
│   └── persistence/
│       └── postgres/
│           ├── models/
│           │   └── [entity].go          # Database model
│           ├── repositories/
│           │   └── [entity]_repository.go # Repository implementation
│           └── migrations/
│               ├── [timestamp]_create_[entity]_table.up.sql
│               └── [timestamp]_create_[entity]_table.down.sql
├── presentation/
│   └── http/
│       ├── handlers/
│       │   └── v1/
│       │       └── [entity]_handler.go   # HTTP handlers
│       └── middleware/                   # HTTP middleware
└── modules/
    └── [entity].go                       # Feature-based DI module
```

## ✅ Validation Rules

### 1. Input Validation
- **Required fields**: Sử dụng `validate:"required"`
- **String length**: `validate:"min=2,max=100"`
- **Email format**: `validate:"email"`
- **Custom validation**: Implement trong validator

### 2. Business Validation
- **Uniqueness**: Check trong database
- **Relationships**: Verify foreign keys exist
- **Business rules**: Implement trong domain layer

### 3. Authorization
- **Role-based**: Check user roles
- **Resource-based**: Check ownership
- **Method-based**: Different permissions for CRUD

## 🚨 Error Handling

### 1. Error Types
```go
// Client errors (4xx)
apperrors.NewBadRequestError("Invalid request")
apperrors.NewNotFoundError("Resource not found")
apperrors.NewConflictError("Resource already exists")
apperrors.NewValidationError("Validation failed", validationErrors)

// Server errors (5xx)
apperrors.NewInternalError("Internal server error", err)
```

### 2. Error Response Format
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "email",
        "message": "Email is required"
      }
    ]
  }
}
```

## 🧪 Testing Guidelines

### 1. Unit Tests
- **Domain Layer**: Test business logic và invariants
- **Application Layer**: Test use case workflows
- **Infrastructure Layer**: Test database operations
- **Presentation Layer**: Test HTTP handlers

### 2. Test Structure
```go
func Test[Entity][Action](t *testing.T) {
    // Arrange
    // Act
    // Assert
}

func Test[Entity][Action]_ShouldReturn[Expected]_When[Condition](t *testing.T) {
    // Specific test case
}
```

## 📚 Examples

### Example User Story
```markdown
## User Story: Create Product

### Business Description
- **Actor**: Authenticated User (Admin role)
- **Action**: Create  
- **Object**: Product
- **Purpose**: Allow admin to add new products to catalog

### Functional Requirements
- **Inputs**:
  - name: string (required, 2-100 chars)
  - description: string (optional, max 1000 chars)
  - price: decimal (required, > 0)
  - category_id: uuid (required, must exist)
- **Outputs**:
  - Created product with ID and timestamps
- **Business Rules**:
  - Product name must be unique within category
  - Price must be positive
  - Category must exist and be active

### Technical Specifications
- **HTTP Method**: POST
- **Endpoint**: /api/v1/products
- **Authentication**: Required (JWT)
- **Authorization**: Admin role required
- **Pagination**: N/A

### Database Impact
- **Tables**: products, categories (relationship)
- **Relationships**: products.category_id -> categories.id
- **Indexes**: idx_products_name_category, idx_products_category_id
- **Migration**: Create products table

### Validation Rules
- **Required fields**: name, price, category_id
- **Format validation**: price must be decimal
- **Business validation**: name unique per category, category exists
- **Size limits**: name 2-100 chars, description max 1000

### Error Scenarios
- **400**: Invalid JSON, validation errors
- **401**: No authentication token
- **403**: Not admin role
- **409**: Product name already exists in category
- **500**: Database errors, external service failures

### Performance Requirements
- **Response Time**: < 500ms
- **Throughput**: 100 RPS
- **Caching**: Cache category lookup
```

Từ User Story này, AI sẽ sinh ra:
1. Domain: Product entity, ProductName/Price value objects, ProductRepository interface
2. Application: CreateProductCommand, ProductDTO, ProductValidator, ProductService
3. Infrastructure: Product model, repository implementation, migration
4. Presentation: ProductHandler với CreateProduct endpoint
5. DI: Binding tất cả dependencies

## 🔄 Dependency Injection Updates

Sau khi sinh code, AI phải cập nhật các DI modules:

```go
// internal/modules/[entity].go
package modules

import (
    "go.uber.org/fx"
    // Import all necessary packages
)

var [Entity]Module = fx.Module("[entity]",
    // Repository
    fx.Provide(repositories.New[Entity]Repository),
    
    // Commands
    fx.Provide(commands.New[Action][Entity]CommandHandler),
    
    // Queries  
    fx.Provide(queries.New[Action][Entity]QueryHandler),
    
    // Services
    fx.Provide(services.New[Entity]Service),
    
    // Validators
    fx.Provide(validators.New[Entity]Validator),
    
    // Handlers
    fx.Provide(handlers.New[Entity]Handler),
    
    // Routes
    fx.Invoke(routes.Setup[Entity]Routes),
)
```

## 📋 Checklist cho AI

Khi sinh API từ User Story, AI phải đảm bảo:

### ✅ Domain Layer
- [ ] Entity với proper value objects và business methods
- [ ] Repository interface với tất cả methods cần thiết  
- [ ] Domain events nếu cần
- [ ] Business validation logic

### ✅ Application Layer
- [ ] Commands cho write operations
- [ ] Queries cho read operations
- [ ] DTOs cho request/response
- [ ] Validators với đầy đủ validation rules
- [ ] Services cho orchestration

### ✅ Infrastructure Layer
- [ ] Database model với correct GORM tags
- [ ] Repository implementation
- [ ] Database migration files
- [ ] Indexes và constraints

### ✅ Presentation Layer
- [ ] HTTP handlers với Swagger documentation
- [ ] Routes setup
- [ ] Middleware assignments
- [ ] Error handling

### ✅ Integration
- [ ] DI module updates
- [ ] All dependencies properly injected
- [ ] Routes registered
- [ ] Migrations executable

### ✅ Quality
- [ ] Proper error handling throughout
- [ ] Logging at appropriate levels
- [ ] Input validation và sanitization
- [ ] Security considerations (auth, authorization)
- [ ] Performance optimizations (indexes, caching)

Với bộ rules này, AI có thể tự động sinh ra một API hoàn chỉnh chỉ từ User Story được format đúng!