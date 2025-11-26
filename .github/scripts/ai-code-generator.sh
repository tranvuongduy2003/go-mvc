#!/bin/bash

# AI-Assisted Code Generation Script
# This script helps generate boilerplate code based on templates

set -e

ENTITY_NAME="$1"
OPERATION="${2:-full}" # full, model, service, handler, repository

if [ -z "$ENTITY_NAME" ]; then
    echo "Usage: ./ai-code-generator.sh <EntityName> [operation]"
    echo "Example: ./ai-code-generator.sh Product full"
    echo ""
    echo "Operations:"
    echo "  full        - Generate all files (model, repository, service, handler)"
    echo "  model       - Generate only model"
    echo "  repository  - Generate only repository"
    echo "  service     - Generate only service"
    echo "  handler     - Generate only handler"
    exit 1
fi

ENTITY_LOWER=$(echo "$ENTITY_NAME" | tr '[:upper:]' '[:lower:]')
ENTITY_SNAKE=$(echo "$ENTITY_NAME" | sed 's/\([A-Z]\)/_\1/g' | tr '[:upper:]' '[:lower:]' | sed 's/^_//')

echo "🤖 AI Code Generator"
echo "==================="
echo "Entity: $ENTITY_NAME"
echo "Operation: $OPERATION"
echo ""

# Create directories
mkdir -p "internal/domain/${ENTITY_LOWER}"
mkdir -p "internal/application/dto/${ENTITY_LOWER}"
mkdir -p "internal/application/commands/${ENTITY_LOWER}"
mkdir -p "internal/application/queries/${ENTITY_LOWER}"
mkdir -p "internal/application/services"
mkdir -p "internal/infrastructure/persistence/repositories"
mkdir -p "internal/presentation/http/handlers"

# Generate model
generate_model() {
    MODEL_FILE="internal/domain/${ENTITY_LOWER}/${ENTITY_LOWER}.go"
    
    if [ -f "$MODEL_FILE" ]; then
        echo "⚠️  Model already exists: $MODEL_FILE"
        return
    fi
    
    cat > "$MODEL_FILE" << EOF
package ${ENTITY_LOWER}

import (
	"time"
	"github.com/google/uuid"
)

// ${ENTITY_NAME} represents a ${ENTITY_LOWER} entity
type ${ENTITY_NAME} struct {
	ID        uuid.UUID  \`gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"\`
	Name      string     \`gorm:"type:varchar(255);not null" json:"name" validate:"required"\`
	CreatedAt time.Time  \`gorm:"not null;default:now()" json:"created_at"\`
	UpdatedAt time.Time  \`gorm:"not null;default:now()" json:"updated_at"\`
	DeletedAt *time.Time \`gorm:"index" json:"deleted_at,omitempty"\`
}

// TableName specifies the table name for ${ENTITY_NAME}
func (${ENTITY_NAME}) TableName() string {
	return "${ENTITY_SNAKE}s"
}
EOF
    
    echo "✅ Generated: $MODEL_FILE"
}

# Generate repository interface
generate_repository_interface() {
    REPO_FILE="internal/domain/repositories/${ENTITY_LOWER}_repository.go"
    
    if [ -f "$REPO_FILE" ]; then
        echo "⚠️  Repository interface already exists: $REPO_FILE"
        return
    fi
    
    cat > "$REPO_FILE" << EOF
package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/${ENTITY_LOWER}"
	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
)

// ${ENTITY_NAME}Repository defines the interface for ${ENTITY_LOWER} data operations
type ${ENTITY_NAME}Repository interface {
	Create(ctx context.Context, entity *${ENTITY_LOWER}.${ENTITY_NAME}) error
	GetByID(ctx context.Context, id uuid.UUID) (*${ENTITY_LOWER}.${ENTITY_NAME}, error)
	List(ctx context.Context, params pagination.Params) ([]${ENTITY_LOWER}.${ENTITY_NAME}, *pagination.Metadata, error)
	Update(ctx context.Context, entity *${ENTITY_LOWER}.${ENTITY_NAME}) error
	Delete(ctx context.Context, id uuid.UUID) error
}
EOF
    
    echo "✅ Generated: $REPO_FILE"
}

# Generate repository implementation
generate_repository_impl() {
    IMPL_FILE="internal/infrastructure/persistence/repositories/${ENTITY_LOWER}_repository.go"
    
    if [ -f "$IMPL_FILE" ]; then
        echo "⚠️  Repository implementation already exists: $IMPL_FILE"
        return
    fi
    
    cat > "$IMPL_FILE" << EOF
package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/${ENTITY_LOWER}"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/repositories"
	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
	"gorm.io/gorm"
)

type ${ENTITY_LOWER}Repository struct {
	db *gorm.DB
}

// New${ENTITY_NAME}Repository creates a new instance of ${ENTITY_NAME}Repository
func New${ENTITY_NAME}Repository(db *gorm.DB) repositories.${ENTITY_NAME}Repository {
	return &${ENTITY_LOWER}Repository{db: db}
}

func (r *${ENTITY_LOWER}Repository) Create(ctx context.Context, entity *${ENTITY_LOWER}.${ENTITY_NAME}) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *${ENTITY_LOWER}Repository) GetByID(ctx context.Context, id uuid.UUID) (*${ENTITY_LOWER}.${ENTITY_NAME}, error) {
	var entity ${ENTITY_LOWER}.${ENTITY_NAME}
	if err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *${ENTITY_LOWER}Repository) List(ctx context.Context, params pagination.Params) ([]${ENTITY_LOWER}.${ENTITY_NAME}, *pagination.Metadata, error) {
	var entities []${ENTITY_LOWER}.${ENTITY_NAME}
	var total int64

	query := r.db.WithContext(ctx).Model(&${ENTITY_LOWER}.${ENTITY_NAME}{})

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	offset := (params.Page - 1) * params.Limit
	if err := query.Offset(offset).Limit(params.Limit).Find(&entities).Error; err != nil {
		return nil, nil, err
	}

	metadata := pagination.NewMetadata(total, params.Page, params.Limit)
	return entities, metadata, nil
}

func (r *${ENTITY_LOWER}Repository) Update(ctx context.Context, entity *${ENTITY_LOWER}.${ENTITY_NAME}) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *${ENTITY_LOWER}Repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&${ENTITY_LOWER}.${ENTITY_NAME}{}, "id = ?", id).Error
}
EOF
    
    echo "✅ Generated: $IMPL_FILE"
}

# Generate service
generate_service() {
    SERVICE_FILE="internal/application/services/${ENTITY_LOWER}_service.go"
    
    if [ -f "$SERVICE_FILE" ]; then
        echo "⚠️  Service already exists: $SERVICE_FILE"
        return
    fi
    
    cat > "$SERVICE_FILE" << EOF
package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/${ENTITY_LOWER}"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/repositories"
	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
)

type ${ENTITY_NAME}Service struct {
	repo repositories.${ENTITY_NAME}Repository
}

func New${ENTITY_NAME}Service(repo repositories.${ENTITY_NAME}Repository) *${ENTITY_NAME}Service {
	return &${ENTITY_NAME}Service{repo: repo}
}

func (s *${ENTITY_NAME}Service) Create(ctx context.Context, entity *${ENTITY_LOWER}.${ENTITY_NAME}) error {
	return s.repo.Create(ctx, entity)
}

func (s *${ENTITY_NAME}Service) GetByID(ctx context.Context, id uuid.UUID) (*${ENTITY_LOWER}.${ENTITY_NAME}, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *${ENTITY_NAME}Service) List(ctx context.Context, params pagination.Params) ([]${ENTITY_LOWER}.${ENTITY_NAME}, *pagination.Metadata, error) {
	return s.repo.List(ctx, params)
}

func (s *${ENTITY_NAME}Service) Update(ctx context.Context, entity *${ENTITY_LOWER}.${ENTITY_NAME}) error {
	return s.repo.Update(ctx, entity)
}

func (s *${ENTITY_NAME}Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
EOF
    
    echo "✅ Generated: $SERVICE_FILE"
}

# Generate handler
generate_handler() {
    HANDLER_FILE="internal/presentation/http/handlers/${ENTITY_LOWER}_handler.go"
    
    if [ -f "$HANDLER_FILE" ]; then
        echo "⚠️  Handler already exists: $HANDLER_FILE"
        return
    fi
    
    cat > "$HANDLER_FILE" << EOF
package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/application/services"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/${ENTITY_LOWER}"
	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
	"github.com/tranvuongduy2003/go-mvc/pkg/response"
)

type ${ENTITY_NAME}Handler struct {
	service *services.${ENTITY_NAME}Service
}

func New${ENTITY_NAME}Handler(service *services.${ENTITY_NAME}Service) *${ENTITY_NAME}Handler {
	return &${ENTITY_NAME}Handler{service: service}
}

// Create handles POST /${ENTITY_LOWER}s
func (h *${ENTITY_NAME}Handler) Create(c *gin.Context) {
	var req ${ENTITY_LOWER}.${ENTITY_NAME}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Create(c.Request.Context(), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, req)
}

// GetByID handles GET /${ENTITY_LOWER}s/:id
func (h *${ENTITY_NAME}Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	entity, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Not found")
		return
	}

	response.Success(c, http.StatusOK, entity)
}

// List handles GET /${ENTITY_LOWER}s
func (h *${ENTITY_NAME}Handler) List(c *gin.Context) {
	params := pagination.GetParams(c)

	entities, metadata, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithPagination(c, http.StatusOK, entities, metadata)
}

// Update handles PUT /${ENTITY_LOWER}s/:id
func (h *${ENTITY_NAME}Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req ${ENTITY_LOWER}.${ENTITY_NAME}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	req.ID = id
	if err := h.service.Update(c.Request.Context(), &req); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, req)
}

// Delete handles DELETE /${ENTITY_LOWER}s/:id
func (h *${ENTITY_NAME}Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusNoContent, nil)
}
EOF
    
    echo "✅ Generated: $HANDLER_FILE"
}

# Execute based on operation
case $OPERATION in
    full)
        generate_model
        generate_repository_interface
        generate_repository_impl
        generate_service
        generate_handler
        ;;
    model)
        generate_model
        ;;
    repository)
        generate_repository_interface
        generate_repository_impl
        ;;
    service)
        generate_service
        ;;
    handler)
        generate_handler
        ;;
    *)
        echo "❌ Unknown operation: $OPERATION"
        exit 1
        ;;
esac

echo ""
echo "✅ Code generation complete!"
echo ""
echo "📝 Next steps:"
echo "1. Review the generated code"
echo "2. Add validation rules to the model"
echo "3. Implement business logic in the service"
echo "4. Add routes in the router"
echo "5. Create migration file"
echo "6. Write tests"
echo ""
echo "💡 Suggested commands:"
echo "  make lint        # Check code quality"
echo "  make test        # Run tests"
echo "  make build       # Build the application"
