# Code Generation Guidelines by Layer

## 📋 Table of Contents
- [Overview](#overview)
- [Domain Layer Generation](#domain-layer-generation)
- [Application Layer Generation](#application-layer-generation)
- [Infrastructure Layer Generation](#infrastructure-layer-generation)
- [Presentation Layer Generation](#presentation-layer-generation)
- [Integration & DI](#integration--di)
- [Testing Generation](#testing-generation)
- [Code Quality Checklist](#code-quality-checklist)

## 🎯 Overview

Tài liệu này cung cấp hướng dẫn chi tiết cho AI về cách sinh code cho từng layer của Clean Architecture. Mỗi layer có những patterns và conventions cụ thể cần tuân thủ.

### Code Generation Flow
```
User Story → Analysis → Domain → Application → Infrastructure → Presentation → Integration → Testing
```

## 🏛️ Domain Layer Generation

### 1. Entity Generation (`internal/core/domain/[entity]/`)

#### Entity Structure Template
```go
package [entity]

import (
    "errors"
    "time"
    "github.com/google/uuid"
    "github.com/tranvuongduy2003/go-mvc/internal/core/domain/shared/events"
)

// [Entity] represents the [entity] aggregate root
type [Entity] struct {
    id        [Entity]ID
    // Value objects và other fields
    createdAt time.Time
    updatedAt time.Time
    version   int64        // For optimistic locking
    events    []events.DomainEvent
}

// Constructor with validation
func New[Entity](
    // Parameters for creation
) (*[Entity], error) {
    // 1. Validate input parameters
    if /* validation condition */ {
        return nil, errors.New("validation error message")
    }
    
    // 2. Create value objects
    // 3. Apply business rules
    // 4. Create entity
    entity := &[Entity]{
        id:        New[Entity]ID(),
        // Set other fields
        createdAt: time.Now(),
        updatedAt: time.Now(),
        version:   1,
        events:    []events.DomainEvent{},
    }
    
    // 5. Add domain events
    entity.addEvent([Entity]CreatedEvent{
        EntityID: entity.id,
        // Event payload
    })
    
    return entity, nil
}

// Business methods
func (e *[Entity]) [BusinessMethod](...) error {
    // 1. Validate preconditions
    // 2. Apply business logic
    // 3. Update entity state
    // 4. Add domain events if needed
    // 5. Update timestamp và version
    e.updatedAt = time.Now()
    e.version++
    
    return nil
}

// Getters (no direct setters)
func (e *[Entity]) ID() [Entity]ID { return e.id }
func (e *[Entity]) CreatedAt() time.Time { return e.createdAt }
func (e *[Entity]) UpdatedAt() time.Time { return e.updatedAt }
func (e *[Entity]) Version() int64 { return e.version }

// Event management
func (e *[Entity]) Events() []events.DomainEvent { return e.events }
func (e *[Entity]) ClearEvents() { e.events = []events.DomainEvent{} }
func (e *[Entity]) addEvent(event events.DomainEvent) {
    e.events = append(e.events, event)
}
```

#### Value Object Generation Rules
```go
// Simple value object
type [ValueObject] struct {
    value [type]
}

func New[ValueObject](value [type]) ([ValueObject], error) {
    // Validation logic based on business rules
    if /* validation condition */ {
        return [ValueObject]{}, errors.New("specific validation message")
    }
    
    // Additional processing if needed
    processedValue := /* processing logic */
    
    return [ValueObject]{value: processedValue}, nil
}

func (vo [ValueObject]) String() string {
    return vo.value // or formatted version
}

func (vo [ValueObject]) Equals(other [ValueObject]) bool {
    return vo.value == other.value
}

// Add specific methods based on business needs
func (vo [ValueObject]) [BusinessMethod]() [returnType] {
    // Business logic specific to this value object
}
```

#### Domain Event Generation
```go
// Domain event structure
type [Entity][Action]Event struct {
    EntityID  [Entity]ID
    // Event payload fields
    Timestamp time.Time
}

func (e [Entity][Action]Event) EventType() string {
    return "[entity].[action]"
}

func (e [Entity][Action]Event) AggregateID() string {
    return e.EntityID.String()
}

func (e [Entity][Action]Event) OccurredAt() time.Time {
    return e.Timestamp
}
```

### 2. Repository Interface Generation (`internal/core/ports/repositories/`)

```go
package repositories

import (
    "context"
    "[entity_package]"
    "github.com/tranvuongduy2003/go-mvc/pkg/pagination"
)

type [Entity]Repository interface {
    // Standard CRUD operations
    Create(ctx context.Context, entity *[entity_package].[Entity]) error
    GetByID(ctx context.Context, id [entity_package].[Entity]ID) (*[entity_package].[Entity], error)
    Update(ctx context.Context, entity *[entity_package].[Entity]) error
    Delete(ctx context.Context, id [entity_package].[Entity]ID) error
    
    // Listing with pagination
    List(ctx context.Context, pagination pagination.Pagination) ([]*[entity_package].[Entity], int64, error)
    
    // Exists check
    Exists(ctx context.Context, id [entity_package].[Entity]ID) (bool, error)
    
    // Custom queries based on User Story requirements
    // Generate these based on business needs identified in User Story
    FindBy[Criteria](ctx context.Context, criteria [CriteriaType]) ([]*[entity_package].[Entity], error)
    
    // For uniqueness checks
    ExistsByUnique[Field](ctx context.Context, field [FieldType]) (bool, error)
    
    // Specifications pattern for complex queries
    FindBySpecification(ctx context.Context, spec Specification) ([]*[entity_package].[Entity], error)
}
```

## 🎯 Application Layer Generation

### 1. Command Generation (`internal/application/commands/[entity]/`)

#### Command Structure
```go
package commands

import (
    "context"
    "[domain_package]"
    "[dto_package]"
    "github.com/tranvuongduy2003/go-mvc/internal/core/ports/repositories"
    "github.com/tranvuongduy2003/go-mvc/internal/core/ports/messaging"
    apperrors "github.com/tranvuongduy2003/go-mvc/pkg/errors"
)

// [Action][Entity]Command represents a command to [action] [entity]
type [Action][Entity]Command struct {
    // Command fields from User Story inputs
    Field1 [type] `json:"field1" validate:"required"`
    Field2 [type] `json:"field2" validate:"omitempty,min=1"`
    // Include user context for authorization
    UserID string `json:"user_id"`
}

// [Action][Entity]CommandHandler handles [action] [entity] command
type [Action][Entity]CommandHandler struct {
    [entity]Repo   repositories.[Entity]Repository
    eventBus       messaging.EventBus
    // Other dependencies based on User Story
}

func New[Action][Entity]CommandHandler(
    [entity]Repo repositories.[Entity]Repository,
    eventBus messaging.EventBus,
    // Other dependencies
) *[Action][Entity]CommandHandler {
    return &[Action][Entity]CommandHandler{
        [entity]Repo: [entity]Repo,
        eventBus:     eventBus,
    }
}

func (h *[Action][Entity]CommandHandler) Handle(ctx context.Context, cmd [Action][Entity]Command) ([dto_package].[Entity]Response, error) {
    // 1. Load existing entities if needed (for updates)
    if /* update operation */ {
        existing, err := h.[entity]Repo.GetByID(ctx, /* id */)
        if err != nil {
            return [dto_package].[Entity]Response{}, apperrors.NewInternalError("Failed to load entity", err)
        }
        if existing == nil {
            return [dto_package].[Entity]Response{}, apperrors.NewNotFoundError("[Entity] not found")
        }
    }
    
    // 2. Create or update domain entity
    entity, err := /* domain operation based on command type */
    if err != nil {
        return [dto_package].[Entity]Response{}, apperrors.NewBadRequestError(err.Error())
    }
    
    // 3. Persist changes
    if err := h.[entity]Repo.[Operation](ctx, entity); err != nil {
        return [dto_package].[Entity]Response{}, apperrors.NewInternalError("Failed to persist entity", err)
    }
    
    // 4. Publish domain events
    for _, event := range entity.Events() {
        if err := h.eventBus.Publish(ctx, event); err != nil {
            // Log error but don't fail the operation
            // logger.Error(ctx, "Failed to publish event", "error", err)
        }
    }
    entity.ClearEvents()
    
    // 5. Return response DTO
    return [dto_package].To[Entity]Response(entity), nil
}
```

### 2. Query Generation (`internal/application/queries/[entity]/`)

```go
package queries

import (
    "context"
    "[dto_package]"
    "github.com/tranvuongduy2003/go-mvc/internal/core/ports/repositories"
    "github.com/tranvuongduy2003/go-mvc/pkg/pagination"
    apperrors "github.com/tranvuongduy2003/go-mvc/pkg/errors"
)

// [Action][Entity]Query represents a query to [action] [entity]
type [Action][Entity]Query struct {
    // Query parameters
    ID string `json:"id" validate:"omitempty,uuid"`
    // Filtering parameters
    Filter[Field] [type] `json:"filter_field,omitempty"`
    // Pagination
    Pagination pagination.Pagination `json:"pagination"`
    // Sorting
    SortBy    string `json:"sort_by,omitempty"`
    SortOrder string `json:"sort_order,omitempty" validate:"omitempty,oneof=asc desc"`
}

// [Action][Entity]QueryHandler handles [action] [entity] query
type [Action][Entity]QueryHandler struct {
    [entity]Repo repositories.[Entity]Repository
}

func New[Action][Entity]QueryHandler([entity]Repo repositories.[Entity]Repository) *[Action][Entity]QueryHandler {
    return &[Action][Entity]QueryHandler{[entity]Repo: [entity]Repo}
}

func (h *[Action][Entity]QueryHandler) Handle(ctx context.Context, query [Action][Entity]Query) ([dto_package].[Entity]ListResponse, error) {
    // 1. Validate query parameters
    if query.Pagination.Limit <= 0 {
        query.Pagination.Limit = 10 // Default page size
    }
    if query.Pagination.Limit > 100 {
        query.Pagination.Limit = 100 // Max page size
    }
    
    // 2. Execute repository query
    entities, total, err := h.[entity]Repo.[QueryMethod](ctx, query.Pagination)
    if err != nil {
        return [dto_package].[Entity]ListResponse{}, apperrors.NewInternalError("Failed to query entities", err)
    }
    
    // 3. Convert to response DTOs
    var responses []dto.[Entity]Response
    for _, entity := range entities {
        responses = append(responses, [dto_package].To[Entity]Response(entity))
    }
    
    // 4. Build paginated response
    return [dto_package].[Entity]ListResponse{
        Data: responses,
        Pagination: [dto_package].PaginationResponse{
            Page:       query.Pagination.Offset/query.Pagination.Limit + 1,
            Limit:      query.Pagination.Limit,
            Total:      total,
            TotalPages: (total + int64(query.Pagination.Limit) - 1) / int64(query.Pagination.Limit),
            HasNext:    query.Pagination.Offset+query.Pagination.Limit < int(total),
            HasPrev:    query.Pagination.Offset > 0,
        },
    }, nil
}
```

### 3. DTO Generation (`internal/application/dto/[entity]/`)

```go
package dto

import (
    "time"
    "[domain_package]"
)

// Request DTOs
type Create[Entity]Request struct {
    // Fields based on User Story inputs với proper validation tags
    Field1 string  `json:"field1" validate:"required,min=2,max=100"`
    Field2 int     `json:"field2" validate:"required,gt=0"`
    Field3 *string `json:"field3,omitempty" validate:"omitempty,max=1000"`
}

type Update[Entity]Request struct {
    // Fields that can be updated
    Field1 *string `json:"field1,omitempty" validate:"omitempty,min=2,max=100"`
    Field2 *int    `json:"field2,omitempty" validate:"omitempty,gt=0"`
}

// Response DTOs
type [Entity]Response struct {
    ID        string    `json:"id"`
    Field1    string    `json:"field1"`
    Field2    int       `json:"field2"`
    Field3    string    `json:"field3"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type [Entity]ListResponse struct {
    Data       []Entity]Response   `json:"data"`
    Pagination PaginationResponse `json:"pagination"`
}

type PaginationResponse struct {
    Page       int  `json:"page"`
    Limit      int  `json:"limit"`
    Total      int64 `json:"total"`
    TotalPages int64 `json:"total_pages"`
    HasNext    bool `json:"has_next"`
    HasPrev    bool `json:"has_prev"`
}

// Mapping functions
func To[Entity]Response(entity *[domain_package].[Entity]) [Entity]Response {
    return [Entity]Response{
        ID:        entity.ID().String(),
        // Map other fields với proper conversion
        Field1:    entity.Field1().String(),
        Field2:    entity.Field2(),
        CreatedAt: entity.CreatedAt(),
        UpdatedAt: entity.UpdatedAt(),
    }
}

func From[Action][Entity]Request(req [Action][Entity]Request) *[domain_package].[Entity]CreationData {
    return &[domain_package].[Entity]CreationData{
        // Map request fields to domain creation data
        Field1: req.Field1,
        Field2: req.Field2,
    }
}
```

### 4. Validator Generation (`internal/application/validators/[entity]/`)

```go
package validators

import (
    "[dto_package]"
    "github.com/tranvuongduy2003/go-mvc/internal/core/ports/repositories"
    apperrors "github.com/tranvuongduy2003/go-mvc/pkg/errors"
    "github.com/tranvuongduy2003/go-mvc/pkg/validator"
)

type I[Entity]Validator interface {
    ValidateCreate[Entity]Request(req [dto_package].Create[Entity]Request) []apperrors.ValidationError
    ValidateUpdate[Entity]Request(req [dto_package].Update[Entity]Request) []apperrors.ValidationError
}

type [Entity]Validator struct {
    [entity]Repo repositories.[Entity]Repository
    validator    *validator.Validator
}

func New[Entity]Validator(
    [entity]Repo repositories.[Entity]Repository,
    validator *validator.Validator,
) I[Entity]Validator {
    return &[Entity]Validator{
        [entity]Repo: [entity]Repo,
        validator:    validator,
    }
}

func (v *[Entity]Validator) ValidateCreate[Entity]Request(req [dto_package].Create[Entity]Request) []apperrors.ValidationError {
    var errors []apperrors.ValidationError
    
    // 1. Structural validation using validator tags
    if validationErrors := v.validator.Validate(req); len(validationErrors) > 0 {
        errors = append(errors, validationErrors...)
    }
    
    // 2. Business validation rules from User Story
    
    // Example: Uniqueness check
    if req.Field1 != "" {
        exists, err := v.[entity]Repo.ExistsByUniqueField1(context.Background(), req.Field1)
        if err != nil {
            errors = append(errors, apperrors.ValidationError{
                Field:   "field1",
                Message: "Failed to check field uniqueness",
            })
        } else if exists {
            errors = append(errors, apperrors.ValidationError{
                Field:   "field1",
                Message: "Field1 already exists",
            })
        }
    }
    
    // Example: Reference validation
    if req.ReferenceID != "" {
        exists, err := v.referenceRepo.Exists(context.Background(), req.ReferenceID)
        if err != nil || !exists {
            errors = append(errors, apperrors.ValidationError{
                Field:   "reference_id",
                Message: "Referenced entity does not exist",
            })
        }
    }
    
    // 3. Complex business rules validation
    if /* complex business rule */ {
        errors = append(errors, apperrors.ValidationError{
            Field:   "business_rule",
            Message: "Business rule violation message",
        })
    }
    
    return errors
}
```

### 5. Service Generation (`internal/application/services/`)

```go
package services

import (
    "context"
    "[command_package]"
    "[query_package]"
    "[dto_package]"
)

type [Entity]Service struct {
    // Command handlers
    create[Entity]Handler *[command_package].Create[Entity]CommandHandler
    update[Entity]Handler *[command_package].Update[Entity]CommandHandler
    delete[Entity]Handler *[command_package].Delete[Entity]CommandHandler
    
    // Query handlers
    get[Entity]Handler  *[query_package].Get[Entity]QueryHandler
    list[Entity]Handler *[query_package].List[Entity]QueryHandler
}

func New[Entity]Service(
    create[Entity]Handler *[command_package].Create[Entity]CommandHandler,
    update[Entity]Handler *[command_package].Update[Entity]CommandHandler,
    delete[Entity]Handler *[command_package].Delete[Entity]CommandHandler,
    get[Entity]Handler  *[query_package].Get[Entity]QueryHandler,
    list[Entity]Handler *[query_package].List[Entity]QueryHandler,
) *[Entity]Service {
    return &[Entity]Service{
        create[Entity]Handler: create[Entity]Handler,
        update[Entity]Handler: update[Entity]Handler,
        delete[Entity]Handler: delete[Entity]Handler,
        get[Entity]Handler:    get[Entity]Handler,
        list[Entity]Handler:   list[Entity]Handler,
    }
}

// Public service methods
func (s *[Entity]Service) Create[Entity](ctx context.Context, req [dto_package].Create[Entity]Request) ([dto_package].[Entity]Response, error) {
    cmd := [command_package].Create[Entity]Command{
        // Map request to command
        Field1: req.Field1,
        Field2: req.Field2,
    }
    
    return s.create[Entity]Handler.Handle(ctx, cmd)
}

func (s *[Entity]Service) Get[Entity](ctx context.Context, id string) ([dto_package].[Entity]Response, error) {
    query := [query_package].Get[Entity]Query{ID: id}
    return s.get[Entity]Handler.Handle(ctx, query)
}

func (s *[Entity]Service) List[Entity](ctx context.Context, req [dto_package].List[Entity]Request) ([dto_package].[Entity]ListResponse, error) {
    query := [query_package].List[Entity]Query{
        Pagination: req.Pagination,
        // Map other filters
    }
    return s.list[Entity]Handler.Handle(ctx, query)
}

// Continue for other operations...
```

## ⚙️ Infrastructure Layer Generation

### 1. Database Model Generation (`internal/adapters/persistence/postgres/models/`)

```go
package models

import (
    "time"
    "[domain_package]"
    "gorm.io/gorm"
)

type [Entity] struct {
    ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
    // Map all domain fields to database columns
    Field1    string    `gorm:"column:field1;type:varchar(100);not null" json:"field1"`
    Field2    int       `gorm:"column:field2;type:integer;not null" json:"field2"`
    Field3    *string   `gorm:"column:field3;type:text" json:"field3"`
    
    // Foreign keys based on relationships
    CategoryID *string  `gorm:"column:category_id;type:uuid" json:"category_id"`
    
    // Standard timestamps
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
    Version   int64     `gorm:"column:version;default:1" json:"version"`
    
    // Define relationships
    Category *Category `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`
    // One-to-many relationships
    Children []ChildModel `gorm:"foreignKey:ParentID;references:ID" json:"children,omitempty"`
}

// Table name
func ([Entity]) TableName() string {
    return "[table_name]"
}

// Hooks for business logic
func (m *[Entity]) BeforeCreate(tx *gorm.DB) error {
    // Pre-creation logic
    return nil
}

func (m *[Entity]) BeforeUpdate(tx *gorm.DB) error {
    // Pre-update logic
    m.Version++
    return nil
}

// Convert from domain entity to model
func From[Entity](entity *[domain_package].[Entity]) *[Entity] {
    return &[Entity]{
        ID:     entity.ID().String(),
        Field1: entity.Field1().String(),
        Field2: entity.Field2(),
        // Handle optional fields
        Field3: func() *string {
            if val := entity.Field3(); val != "" {
                return &val
            }
            return nil
        }(),
        // Map other fields
        Version: entity.Version(),
    }
}

// Convert from model to domain entity
func (m *[Entity]) To[Entity]() (*[domain_package].[Entity], error) {
    // Reconstruct value objects với validation
    field1, err := [domain_package].NewField1(m.Field1)
    if err != nil {
        return nil, err
    }
    
    // Reconstruct entity
    entity := &[domain_package].[Entity]{}
    // Use reflection or manual mapping to set private fields
    // This might require domain factory methods
    
    return entity, nil
}
```

### 2. Repository Implementation (`internal/adapters/persistence/postgres/repositories/`)

```go
package repositories

import (
    "context"
    "gorm.io/gorm"
    "[domain_package]"
    "[models_package]"
    "github.com/tranvuongduy2003/go-mvc/internal/core/ports/repositories"
    "github.com/tranvuongduy2003/go-mvc/pkg/pagination"
    apperrors "github.com/tranvuongduy2003/go-mvc/pkg/errors"
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
        // Handle specific database errors
        if /* unique constraint violation */ {
            return apperrors.NewConflictError("Entity already exists")
        }
        return apperrors.NewInternalError("Failed to create entity", err)
    }
    
    return nil
}

func (r *[Entity]Repository) GetByID(ctx context.Context, id [domain_package].[Entity]ID) (*[domain_package].[Entity], error) {
    var model [models_package].[Entity]
    
    err := r.db.WithContext(ctx).
        Preload("RelatedEntities"). // Load related entities if needed
        First(&model, "id = ?", id.String()).Error
        
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil // Not found, return nil without error
        }
        return nil, apperrors.NewInternalError("Failed to get entity", err)
    }
    
    entity, err := model.To[Entity]()
    if err != nil {
        return nil, apperrors.NewInternalError("Failed to convert model to entity", err)
    }
    
    return entity, nil
}

func (r *[Entity]Repository) Update(ctx context.Context, entity *[domain_package].[Entity]) error {
    model := [models_package].From[Entity](entity)
    
    // Optimistic locking check
    result := r.db.WithContext(ctx).
        Where("id = ? AND version = ?", model.ID, model.Version-1).
        Updates(model)
        
    if result.Error != nil {
        return apperrors.NewInternalError("Failed to update entity", result.Error)
    }
    
    if result.RowsAffected == 0 {
        return apperrors.NewConflictError("Entity was modified by another process")
    }
    
    return nil
}

func (r *[Entity]Repository) Delete(ctx context.Context, id [domain_package].[Entity]ID) error {
    result := r.db.WithContext(ctx).
        Delete(&[models_package].[Entity]{}, "id = ?", id.String())
        
    if result.Error != nil {
        return apperrors.NewInternalError("Failed to delete entity", result.Error)
    }
    
    if result.RowsAffected == 0 {
        return apperrors.NewNotFoundError("Entity not found")
    }
    
    return nil
}

func (r *[Entity]Repository) List(ctx context.Context, pag pagination.Pagination) ([]*[domain_package].[Entity], int64, error) {
    var models []models.[Entity]
    var total int64
    
    // Count total records
    if err := r.db.WithContext(ctx).
        Model(&[models_package].[Entity]{}).
        Count(&total).Error; err != nil {
        return nil, 0, apperrors.NewInternalError("Failed to count entities", err)
    }
    
    // Get paginated records
    if err := r.db.WithContext(ctx).
        Preload("RelatedEntities").
        Offset(pag.Offset).
        Limit(pag.Limit).
        Find(&models).Error; err != nil {
        return nil, 0, apperrors.NewInternalError("Failed to list entities", err)
    }
    
    // Convert models to entities
    var entities []*[domain_package].[Entity]
    for _, model := range models {
        entity, err := model.To[Entity]()
        if err != nil {
            return nil, 0, apperrors.NewInternalError("Failed to convert model to entity", err)
        }
        entities = append(entities, entity)
    }
    
    return entities, total, nil
}

// Custom queries based on User Story
func (r *[Entity]Repository) FindBy[Criteria](ctx context.Context, criteria [CriteriaType]) ([]*[domain_package].[Entity], error) {
    var models []models.[Entity]
    
    query := r.db.WithContext(ctx)
    
    // Build query based on criteria
    if /* criteria condition */ {
        query = query.Where("field = ?", criteria.Field)
    }
    
    if err := query.Find(&models).Error; err != nil {
        return nil, apperrors.NewInternalError("Failed to find entities", err)
    }
    
    // Convert to entities
    var entities []*[domain_package].[Entity]
    for _, model := range models {
        entity, err := model.To[Entity]()
        if err != nil {
            return nil, apperrors.NewInternalError("Failed to convert model", err)
        }
        entities = append(entities, entity)
    }
    
    return entities, nil
}

func (r *[Entity]Repository) ExistsByUnique[Field](ctx context.Context, field [FieldType]) (bool, error) {
    var count int64
    
    err := r.db.WithContext(ctx).
        Model(&[models_package].[Entity]{}).
        Where("[field_column] = ?", field).
        Count(&count).Error
        
    if err != nil {
        return false, apperrors.NewInternalError("Failed to check existence", err)
    }
    
    return count > 0, nil
}
```

### 3. Migration Generation (`internal/adapters/persistence/postgres/migrations/`)

```sql
-- [timestamp]_create_[entity]_table.up.sql

-- Create table với proper constraints
CREATE TABLE IF NOT EXISTS [table_name] (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Map domain fields to columns
    field1 VARCHAR(100) NOT NULL,
    field2 INTEGER NOT NULL CHECK (field2 > 0),
    field3 TEXT,
    
    -- Foreign keys based on relationships
    category_id UUID REFERENCES categories(id) ON DELETE RESTRICT,
    
    -- Version for optimistic locking
    version BIGINT NOT NULL DEFAULT 1,
    
    -- Standard timestamps
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes based on User Story requirements
CREATE INDEX IF NOT EXISTS idx_[table_name]_deleted_at ON [table_name](deleted_at);
CREATE INDEX IF NOT EXISTS idx_[table_name]_category_id ON [table_name](category_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_[table_name]_field1 ON [table_name](field1) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_[table_name]_created_at ON [table_name](created_at DESC) WHERE deleted_at IS NULL;

-- Unique constraints from business rules
CREATE UNIQUE INDEX IF NOT EXISTS idx_[table_name]_field1_category_unique 
    ON [table_name](field1, category_id) WHERE deleted_at IS NULL;

-- Triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_[table_name]_updated_at 
    BEFORE UPDATE ON [table_name] 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Comments for documentation
COMMENT ON TABLE [table_name] IS 'Table for storing [entity] entities';
COMMENT ON COLUMN [table_name].id IS 'Unique identifier for [entity]';
COMMENT ON COLUMN [table_name].field1 IS 'Business description of field1';
COMMENT ON COLUMN [table_name].version IS 'Version for optimistic locking';
```

```sql
-- [timestamp]_create_[entity]_table.down.sql

-- Drop triggers first
DROP TRIGGER IF EXISTS update_[table_name]_updated_at ON [table_name];

-- Drop indexes
DROP INDEX IF EXISTS idx_[table_name]_field1_category_unique;
DROP INDEX IF EXISTS idx_[table_name]_created_at;
DROP INDEX IF EXISTS idx_[table_name]_field1;
DROP INDEX IF EXISTS idx_[table_name]_category_id;
DROP INDEX IF EXISTS idx_[table_name]_deleted_at;

-- Drop table
DROP TABLE IF EXISTS [table_name];

-- Drop function if not used by other tables
DROP FUNCTION IF EXISTS update_updated_at_column();
```

## 🌐 Presentation Layer Generation

### 1. HTTP Handler Generation (`internal/handlers/http/rest/v1/`)

```go
package v1

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "[dto_package]"
    "[service_package]"
    "[validator_package]"
    "github.com/tranvuongduy2003/go-mvc/pkg/response"
    "github.com/tranvuongduy2003/go-mvc/pkg/pagination"
    apperrors "github.com/tranvuongduy2003/go-mvc/pkg/errors"
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

// Create[Entity] creates a new [entity]
// @Summary Create a new [entity]
// @Description Create a new [entity] with the provided data
// @Tags [entity]
// @Accept json
// @Produce json
// @Param [entity] body dto.Create[Entity]Request true "[Entity] creation data"
// @Success 201 {object} response.APIResponse{data=dto.[Entity]Response}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/[entities] [post]
// @Security Bearer
func (h *[Entity]Handler) Create[Entity](c *gin.Context) {
    var req [dto_package].Create[Entity]Request
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, apperrors.NewBadRequestError("Invalid JSON format"))
        return
    }

    // Extract user context for authorization
    userID := c.GetString("user_id") // From JWT middleware
    userRole := c.GetString("user_role")
    
    // Authorization check based on User Story
    if /* authorization rule from User Story */ {
        response.Error(c, apperrors.NewForbiddenError("Access denied"))
        return
    }

    // Validate request
    if validationErrors := h.[entity]Validator.ValidateCreate[Entity]Request(req); len(validationErrors) > 0 {
        response.ValidationError(c, validationErrors)
        return
    }

    // Add user context to request
    req.UserID = userID

    // Call service
    result, err := h.[entity]Service.Create[Entity](c.Request.Context(), req)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Created(c, result, "[Entity] created successfully")
}

// Get[Entity] retrieves a [entity] by ID
// @Summary Get [entity] by ID
// @Description Get a specific [entity] by its ID
// @Tags [entity]
// @Accept json
// @Produce json
// @Param id path string true "[Entity] ID" Format(uuid)
// @Success 200 {object} response.APIResponse{data=dto.[Entity]Response}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/[entities]/{id} [get]
// @Security Bearer
func (h *[Entity]Handler) Get[Entity](c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        response.Error(c, apperrors.NewBadRequestError("ID is required"))
        return
    }

    // Authorization check if needed
    userID := c.GetString("user_id")
    if /* resource-based authorization */ {
        response.Error(c, apperrors.NewForbiddenError("Access denied"))
        return
    }

    result, err := h.[entity]Service.Get[Entity](c.Request.Context(), id)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, result, "[Entity] retrieved successfully")
}

// List[Entity] retrieves a paginated list of [entities]
// @Summary List [entities] với pagination
// @Description Get a paginated list of [entities] với optional filtering
// @Tags [entity]
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param sort_by query string false "Sort field" Enums([entity]_field1, [entity]_field2, created_at, updated_at)
// @Param sort_order query string false "Sort order" Enums(asc, desc) default(desc)
// @Param filter_field query string false "Filter by field value"
// @Success 200 {object} response.APIResponse{data=dto.[Entity]ListResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/[entities] [get]
// @Security Bearer
func (h *[Entity]Handler) List[Entity](c *gin.Context) {
    // Parse pagination parameters
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    
    // Validate pagination parameters
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 10
    }

    req := [dto_package].List[Entity]Request{
        Pagination: pagination.Pagination{
            Offset: (page - 1) * limit,
            Limit:  limit,
        },
        SortBy:    c.Query("sort_by"),
        SortOrder: c.DefaultQuery("sort_order", "desc"),
        // Map query parameters to filters
        FilterField: c.Query("filter_field"),
    }

    // Authorization - filter based on user permissions
    userID := c.GetString("user_id")
    userRole := c.GetString("user_role")
    if /* apply user-based filtering */ {
        req.UserID = userID
    }

    result, err := h.[entity]Service.List[Entity](c.Request.Context(), req)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, result, "[Entity] list retrieved successfully")
}

// Update[Entity] updates an existing [entity]
// @Summary Update [entity]
// @Description Update an existing [entity] với the provided data
// @Tags [entity]
// @Accept json
// @Produce json
// @Param id path string true "[Entity] ID" Format(uuid)
// @Param [entity] body dto.Update[Entity]Request true "[Entity] update data"
// @Success 200 {object} response.APIResponse{data=dto.[Entity]Response}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/[entities]/{id} [put]
// @Security Bearer
func (h *[Entity]Handler) Update[Entity](c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        response.Error(c, apperrors.NewBadRequestError("ID is required"))
        return
    }

    var req [dto_package].Update[Entity]Request
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, apperrors.NewBadRequestError("Invalid JSON format"))
        return
    }

    // Authorization check
    userID := c.GetString("user_id")
    if /* resource ownership or role check */ {
        response.Error(c, apperrors.NewForbiddenError("Access denied"))
        return
    }

    // Validate request
    if validationErrors := h.[entity]Validator.ValidateUpdate[Entity]Request(req); len(validationErrors) > 0 {
        response.ValidationError(c, validationErrors)
        return
    }

    result, err := h.[entity]Service.Update[Entity](c.Request.Context(), id, req)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, result, "[Entity] updated successfully")
}

// Delete[Entity] deletes a [entity]
// @Summary Delete [entity]
// @Description Delete a specific [entity] by its ID
// @Tags [entity]
// @Accept json
// @Produce json
// @Param id path string true "[Entity] ID" Format(uuid)
// @Success 204 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/[entities]/{id} [delete]
// @Security Bearer
func (h *[Entity]Handler) Delete[Entity](c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        response.Error(c, apperrors.NewBadRequestError("ID is required"))
        return
    }

    // Authorization check
    userID := c.GetString("user_id")
    userRole := c.GetString("user_role")
    if /* deletion authorization rules */ {
        response.Error(c, apperrors.NewForbiddenError("Access denied"))
        return
    }

    err := h.[entity]Service.Delete[Entity](c.Request.Context(), id)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.NoContent(c, "[Entity] deleted successfully")
}
```

### 2. Route Setup Generation

```go
// internal/handlers/http/rest/v1/routes.go

package v1

import (
    "github.com/gin-gonic/gin"
    "github.com/tranvuongduy2003/go-mvc/internal/handlers/http/middleware"
)

func Setup[Entity]Routes(
    router *gin.RouterGroup,
    [entity]Handler *[Entity]Handler,
    authMiddleware *middleware.AuthMiddleware,
    rbacMiddleware *middleware.RBACMiddleware,
) {
    [entities] := router.Group("/[entities]")
    
    // Apply authentication to all routes
    [entities].Use(authMiddleware.RequireAuth())
    
    // Public routes (if any) - based on User Story
    
    // Protected routes
    [entities].POST("", 
        rbacMiddleware.RequireRole("admin"), // Based on User Story authorization
        [entity]Handler.Create[Entity])
        
    [entities].GET("", [entity]Handler.List[Entity]) // May need role-based filtering
    
    [entities].GET("/:id", [entity]Handler.Get[Entity])
    
    [entities].PUT("/:id", 
        rbacMiddleware.RequireRoleOrOwnership("admin", "owner"),
        [entity]Handler.Update[Entity])
        
    [entities].DELETE("/:id", 
        rbacMiddleware.RequireRole("admin"),
        [entity]Handler.Delete[Entity])
}
```

## 🔧 Integration & DI

### DI Module Generation (`internal/di/modules/`)

```go
package modules

import (
    "go.uber.org/fx"
    
    // Domain
    "github.com/tranvuongduy2003/go-mvc/internal/core/ports/repositories"
    
    // Application
    "[command_package]"
    "[query_package]"
    "[service_package]"
    "[validator_package]"
    
    // Infrastructure
    "[repo_impl_package]"
    
    // Presentation
    "[handler_package]"
)

var [Entity]Module = fx.Module("[entity]",
    // Repository
    fx.Provide(fx.Annotate(
        [repo_impl_package].New[Entity]Repository,
        fx.As(new(repositories.[Entity]Repository)),
    )),
    
    // Command handlers
    fx.Provide([command_package].NewCreate[Entity]CommandHandler),
    fx.Provide([command_package].NewUpdate[Entity]CommandHandler),
    fx.Provide([command_package].NewDelete[Entity]CommandHandler),
    
    // Query handlers
    fx.Provide([query_package].NewGet[Entity]QueryHandler),
    fx.Provide([query_package].NewList[Entity]QueryHandler),
    
    // Application services
    fx.Provide([service_package].New[Entity]Service),
    
    // Validators
    fx.Provide(fx.Annotate(
        [validator_package].New[Entity]Validator,
        fx.As(new([validator_package].I[Entity]Validator)),
    )),
    
    // HTTP handlers
    fx.Provide([handler_package].New[Entity]Handler),
    
    // Route setup
    fx.Invoke([handler_package].Setup[Entity]Routes),
)
```

### Main DI Update (`internal/di/application.go`)

```go
// Add to existing ApplicationModule
var ApplicationModule = fx.Module("application",
    // Existing modules...
    
    // Add new entity module
    [Entity]Module,
)
```

## 🧪 Testing Generation

### 1. Unit Test Templates

```go
// Domain entity test
package [entity]_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "[domain_package]"
)

func Test[Entity]_New[Entity]_ShouldCreateValidEntity_WhenValidDataProvided(t *testing.T) {
    // Arrange
    // Valid test data based on User Story
    
    // Act
    entity, err := [domain_package].New[Entity](/* valid params */)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, entity)
    assert.NotEmpty(t, entity.ID())
    // Assert business rules
}

func Test[Entity]_New[Entity]_ShouldReturnError_WhenInvalidDataProvided(t *testing.T) {
    // Test various invalid scenarios from User Story validation rules
}

func Test[Entity]_BusinessMethod_ShouldApplyBusinessLogic_WhenValidConditions(t *testing.T) {
    // Test business methods và domain events
}
```

```go
// Application service test
package [service]_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func Test[Entity]Service_Create[Entity]_ShouldReturnEntity_WhenValidRequest(t *testing.T) {
    // Arrange
    mockRepo := &MockEntityRepository{}
    service := NewEntityService(mockRepo, ...)
    
    req := dto.CreateEntityRequest{
        // Valid request data
    }
    
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Entity")).Return(nil)
    
    // Act
    result, err := service.CreateEntity(context.Background(), req)
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, result.ID)
    mockRepo.AssertExpectations(t)
}
```

```go
// HTTP handler test
package v1_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestEntityHandler_CreateEntity_ShouldReturn201_WhenValidRequest(t *testing.T) {
    // Arrange
    gin.SetMode(gin.TestMode)
    router := gin.New()
    
    mockService := &MockEntityService{}
    handler := NewEntityHandler(mockService, mockValidator)
    
    router.POST("/entities", handler.CreateEntity)
    
    req := dto.CreateEntityRequest{
        // Valid request data
    }
    reqBody, _ := json.Marshal(req)
    
    mockService.On("CreateEntity", mock.Anything, req).Return(expectedResponse, nil)
    
    // Act
    w := httptest.NewRecorder()
    httpReq, _ := http.NewRequest("POST", "/entities", bytes.NewBuffer(reqBody))
    httpReq.Header.Set("Content-Type", "application/json")
    
    router.ServeHTTP(w, httpReq)
    
    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)
    mockService.AssertExpectations(t)
}
```

## ✅ Code Quality Checklist

### AI Generated Code Must Include:

#### ✅ Domain Layer
- [ ] Entities với proper encapsulation (private fields, public methods)
- [ ] Value objects với validation logic
- [ ] Business methods với domain logic
- [ ] Domain events for important business actions
- [ ] Repository interfaces với all necessary methods
- [ ] Proper error handling với domain-specific errors

#### ✅ Application Layer  
- [ ] Commands cho write operations với proper validation
- [ ] Queries cho read operations với pagination support
- [ ] DTOs với JSON tags và validation tags
- [ ] Request/Response mapping functions
- [ ] Validators với both structural và business validation
- [ ] Application services với proper error handling

#### ✅ Infrastructure Layer
- [ ] Database models với proper GORM tags
- [ ] Repository implementations với error handling
- [ ] Database migrations với proper constraints và indexes
- [ ] Foreign key relationships correctly defined
- [ ] Optimistic locking implementation
- [ ] Proper database connection handling

#### ✅ Presentation Layer
- [ ] HTTP handlers với Swagger documentation
- [ ] Request binding và validation
- [ ] Authorization checks based on User Story
- [ ] Proper HTTP status codes
- [ ] Response formatting với consistent structure
- [ ] Route definitions với middleware

#### ✅ Integration
- [ ] DI modules với all dependencies
- [ ] Route registration
- [ ] Middleware assignments
- [ ] Error handling chain
- [ ] Logging integration
- [ ] Metrics collection

#### ✅ Quality Attributes
- [ ] Consistent naming conventions
- [ ] Proper import organization
- [ ] Error messages are user-friendly
- [ ] Business rules are enforced
- [ ] Security considerations addressed
- [ ] Performance optimizations (indexes, caching)
- [ ] Input sanitization và validation
- [ ] Audit logging for sensitive operations

### Generated Files Checklist:

```
✅ Domain Layer:
   ├── internal/core/domain/[entity]/
   │   ├── [entity].go                    # Main entity
   │   ├── value_objects.go               # Value objects
   │   └── events.go                      # Domain events
   └── internal/core/ports/repositories/
       └── [entity]_repository.go         # Repository interface

✅ Application Layer:
   ├── internal/application/commands/[entity]/
   │   ├── create_[entity]_command.go     # Create command
   │   ├── update_[entity]_command.go     # Update command
   │   └── delete_[entity]_command.go     # Delete command
   ├── internal/application/queries/[entity]/
   │   ├── get_[entity]_query.go          # Get query
   │   └── list_[entity]_query.go         # List query
   ├── internal/application/dto/[entity]/
   │   └── [entity]_dto.go                # DTOs và mappers
   ├── internal/application/services/
   │   └── [entity]_service.go            # Application service
   └── internal/application/validators/[entity]/
       └── [entity]_validator.go          # Validators

✅ Infrastructure Layer:
   ├── internal/adapters/persistence/postgres/models/
   │   └── [entity].go                    # Database model
   ├── internal/adapters/persistence/postgres/repositories/
   │   └── [entity]_repository.go         # Repository implementation
   └── internal/adapters/persistence/postgres/migrations/
       ├── [timestamp]_create_[entity]_table.up.sql
       └── [timestamp]_create_[entity]_table.down.sql

✅ Presentation Layer:
   ├── internal/handlers/http/rest/v1/
   │   └── [entity]_handler.go            # HTTP handlers
   └── internal/handlers/http/rest/v1/
       └── routes.go                      # Route setup

✅ Integration:
   ├── internal/di/modules/
   │   └── [entity].go                    # DI module
   └── Updated main DI files
```

Với bộ guidelines này, AI có thể sinh ra code hoàn chỉnh, chất lượng cao, và sẵn sàng production cho bất kỳ User Story nào được format đúng theo template!