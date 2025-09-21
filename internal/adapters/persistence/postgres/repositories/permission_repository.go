package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/tranvuongduy2003/go-mvc/internal/adapters/persistence/postgres/models"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/rbac"
)

// PermissionRepositoryImpl implements the RBAC PermissionRepository interface
type PermissionRepositoryImpl struct {
	db *gorm.DB
}

// NewPermissionRepository creates a new permission repository
func NewPermissionRepository(db *gorm.DB) rbac.PermissionRepository {
	return &PermissionRepositoryImpl{
		db: db,
	}
}

// Create creates a new permission
func (r *PermissionRepositoryImpl) Create(ctx context.Context, permission *rbac.Permission) error {
	model := &models.Permission{
		ID:          permission.ID,
		Resource:    permission.Resource,
		Action:      permission.Action,
		DisplayName: permission.Name,
		Description: permission.Description,
		IsActive:    permission.IsActive,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create permission: %w", err)
	}

	return nil
}

// GetByID retrieves a permission by ID
func (r *PermissionRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*rbac.Permission, error) {
	var model models.Permission
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, rbac.ErrPermissionNotFound
		}
		return nil, fmt.Errorf("failed to get permission by ID: %w", err)
	}

	return r.toDomain(&model), nil
}

// GetByResourceAndAction retrieves a permission by resource and action
func (r *PermissionRepositoryImpl) GetByResourceAndAction(ctx context.Context, resource, action string) (*rbac.Permission, error) {
	var model models.Permission
	if err := r.db.WithContext(ctx).
		First(&model, "resource = ? AND action = ?", resource, action).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, rbac.ErrPermissionNotFound
		}
		return nil, fmt.Errorf("failed to get permission by resource and action: %w", err)
	}

	return r.toDomain(&model), nil
}

// GetByName retrieves a permission by name
func (r *PermissionRepositoryImpl) GetByName(ctx context.Context, name string) (*rbac.Permission, error) {
	var model models.Permission
	if err := r.db.WithContext(ctx).
		First(&model, "display_name = ?", name).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, rbac.ErrPermissionNotFound
		}
		return nil, fmt.Errorf("failed to get permission by name: %w", err)
	}

	return r.toDomain(&model), nil
}

// GetByAction retrieves permissions by action
func (r *PermissionRepositoryImpl) GetByAction(ctx context.Context, action string) ([]*rbac.Permission, error) {
	var models []models.Permission
	if err := r.db.WithContext(ctx).
		Where("action = ?", action).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get permissions by action: %w", err)
	}

	permissions := make([]*rbac.Permission, len(models))
	for i, model := range models {
		permissions[i] = r.toDomain(&model)
	}

	return permissions, nil
}

// Update updates a permission
func (r *PermissionRepositoryImpl) Update(ctx context.Context, permission *rbac.Permission) error {
	model := &models.Permission{
		ID:          permission.ID,
		Resource:    permission.Resource,
		Action:      permission.Action,
		DisplayName: permission.Name,
		Description: permission.Description,
		IsActive:    permission.IsActive,
		UpdatedAt:   permission.UpdatedAt,
	}

	if err := r.db.WithContext(ctx).
		Model(model).
		Where("id = ?", permission.ID).
		Updates(model).Error; err != nil {
		return fmt.Errorf("failed to update permission: %w", err)
	}

	return nil
}

// Delete deletes a permission
func (r *PermissionRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.Permission{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}

	return nil
}

// List retrieves permissions with pagination
func (r *PermissionRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*rbac.Permission, error) {
	var models []models.Permission
	query := r.db.WithContext(ctx)

	if offset > 0 {
		query = query.Offset(offset)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}

	permissions := make([]*rbac.Permission, len(models))
	for i, model := range models {
		permissions[i] = r.toDomain(&model)
	}

	return permissions, nil
}

// Count returns the total number of permissions
func (r *PermissionRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Permission{}).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count permissions: %w", err)
	}

	return count, nil
}

// GetByResource retrieves permissions by resource
func (r *PermissionRepositoryImpl) GetByResource(ctx context.Context, resource string) ([]*rbac.Permission, error) {
	var models []models.Permission
	if err := r.db.WithContext(ctx).
		Where("resource = ?", resource).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get permissions by resource: %w", err)
	}

	permissions := make([]*rbac.Permission, len(models))
	for i, model := range models {
		permissions[i] = r.toDomain(&model)
	}

	return permissions, nil
}

// GetActivePermissions retrieves all active permissions
func (r *PermissionRepositoryImpl) GetActivePermissions(ctx context.Context) ([]*rbac.Permission, error) {
	var models []models.Permission
	if err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get active permissions: %w", err)
	}

	permissions := make([]*rbac.Permission, len(models))
	for i, model := range models {
		permissions[i] = r.toDomain(&model)
	}

	return permissions, nil
}

// SearchByName searches permissions by name with pagination
func (r *PermissionRepositoryImpl) SearchByName(ctx context.Context, name string, offset, limit int) ([]*rbac.Permission, error) {
	var models []models.Permission
	query := r.db.WithContext(ctx).
		Where("display_name ILIKE ?", "%"+name+"%")

	if offset > 0 {
		query = query.Offset(offset)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to search permissions by name: %w", err)
	}

	permissions := make([]*rbac.Permission, len(models))
	for i, model := range models {
		permissions[i] = r.toDomain(&model)
	}

	return permissions, nil
}

// GetPermissionsByIDs retrieves permissions by IDs
func (r *PermissionRepositoryImpl) GetPermissionsByIDs(ctx context.Context, ids []uuid.UUID) ([]*rbac.Permission, error) {
	var models []models.Permission
	if err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get permissions by IDs: %w", err)
	}

	permissions := make([]*rbac.Permission, len(models))
	for i, model := range models {
		permissions[i] = r.toDomain(&model)
	}

	return permissions, nil
}

// ExistsByResourceAndAction checks if a permission exists by resource and action
func (r *PermissionRepositoryImpl) ExistsByResourceAndAction(ctx context.Context, resource, action string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Permission{}).
		Where("resource = ? AND action = ?", resource, action).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check if permission exists: %w", err)
	}

	return count > 0, nil
}

// CreateBulk creates multiple permissions
func (r *PermissionRepositoryImpl) CreateBulk(ctx context.Context, permissions []*rbac.Permission) error {
	permissionModels := make([]models.Permission, len(permissions))
	for i, perm := range permissions {
		permissionModels[i] = models.Permission{
			ID:          perm.ID,
			Resource:    perm.Resource,
			Action:      perm.Action,
			DisplayName: perm.Name,
			Description: perm.Description,
			IsActive:    perm.IsActive,
			CreatedAt:   perm.CreatedAt,
			UpdatedAt:   perm.UpdatedAt,
		}
	}

	if err := r.db.WithContext(ctx).Create(&permissionModels).Error; err != nil {
		return fmt.Errorf("failed to create permissions in bulk: %w", err)
	}

	return nil
}

// toDomain converts a permission model to domain entity
func (r *PermissionRepositoryImpl) toDomain(model *models.Permission) *rbac.Permission {
	return &rbac.Permission{
		ID:          model.ID,
		Name:        model.DisplayName,
		Resource:    model.Resource,
		Action:      model.Action,
		Description: model.Description,
		IsActive:    model.IsActive,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}
