package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/tranvuongduy2003/go-mvc/internal/adapters/persistence/postgres/models"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/rbac"
)

// RoleRepositoryImpl implements the RBAC RoleRepository interface
type RoleRepositoryImpl struct {
	db *gorm.DB
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *gorm.DB) rbac.RoleRepository {
	return &RoleRepositoryImpl{
		db: db,
	}
}

// Create creates a new role
func (r *RoleRepositoryImpl) Create(ctx context.Context, role *rbac.Role) error {
	model := &models.Role{
		ID:          role.ID,
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Description: role.Description,
		IsActive:    role.IsActive,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

// GetByID retrieves a role by ID
func (r *RoleRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*rbac.Role, error) {
	var model models.Role
	if err := r.db.WithContext(ctx).
		Preload("Permissions").
		First(&model, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, rbac.ErrRoleNotFound
		}
		return nil, fmt.Errorf("failed to get role by ID: %w", err)
	}

	return r.toDomain(&model), nil
}

// GetByName retrieves a role by name
func (r *RoleRepositoryImpl) GetByName(ctx context.Context, name string) (*rbac.Role, error) {
	var model models.Role
	if err := r.db.WithContext(ctx).
		Preload("Permissions").
		First(&model, "name = ?", name).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, rbac.ErrRoleNotFound
		}
		return nil, fmt.Errorf("failed to get role by name: %w", err)
	}

	return r.toDomain(&model), nil
}

// Update updates a role
func (r *RoleRepositoryImpl) Update(ctx context.Context, role *rbac.Role) error {
	model := &models.Role{
		ID:          role.ID,
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Description: role.Description,
		IsActive:    role.IsActive,
		UpdatedAt:   role.UpdatedAt,
	}

	if err := r.db.WithContext(ctx).
		Model(model).
		Where("id = ?", role.ID).
		Updates(model).Error; err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	return nil
}

// Delete deletes a role
func (r *RoleRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.Role{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}

// List retrieves roles with pagination
func (r *RoleRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*rbac.Role, error) {
	var models []models.Role
	query := r.db.WithContext(ctx).Preload("Permissions")

	if offset > 0 {
		query = query.Offset(offset)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	roles := make([]*rbac.Role, len(models))
	for i, model := range models {
		roles[i] = r.toDomain(&model)
	}

	return roles, nil
}

// GetActiveRoles retrieves all active roles
func (r *RoleRepositoryImpl) GetActiveRoles(ctx context.Context) ([]*rbac.Role, error) {
	var models []models.Role
	if err := r.db.WithContext(ctx).
		Preload("Permissions").
		Where("is_active = ?", true).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get active roles: %w", err)
	}

	roles := make([]*rbac.Role, len(models))
	for i, model := range models {
		roles[i] = r.toDomain(&model)
	}

	return roles, nil
}

// ExistsByName checks if a role exists by name
func (r *RoleRepositoryImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Role{}).
		Where("name = ?", name).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check if role exists by name: %w", err)
	}

	return count > 0, nil
}

// Count returns the total number of roles
func (r *RoleRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Role{}).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count roles: %w", err)
	}

	return count, nil
}

// GetRolesWithPermissions retrieves roles with their permissions
func (r *RoleRepositoryImpl) GetRolesWithPermissions(ctx context.Context, roleIDs []uuid.UUID) ([]*rbac.Role, error) {
	var models []models.Role
	if err := r.db.WithContext(ctx).
		Preload("Permissions").
		Where("id IN ?", roleIDs).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get roles with permissions: %w", err)
	}

	roles := make([]*rbac.Role, len(models))
	for i, model := range models {
		roles[i] = r.toDomain(&model)
	}

	return roles, nil
}

// GetUserRoles retrieves roles for a specific user
func (r *RoleRepositoryImpl) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*rbac.Role, error) {
	var models []models.Role
	if err := r.db.WithContext(ctx).
		Preload("Permissions").
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roles := make([]*rbac.Role, len(models))
	for i, model := range models {
		roles[i] = r.toDomain(&model)
	}

	return roles, nil
}

// SearchByName searches roles by name with pagination
func (r *RoleRepositoryImpl) SearchByName(ctx context.Context, name string, offset, limit int) ([]*rbac.Role, error) {
	var models []models.Role
	query := r.db.WithContext(ctx).
		Preload("Permissions").
		Where("name ILIKE ?", "%"+name+"%")

	if offset > 0 {
		query = query.Offset(offset)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to search roles by name: %w", err)
	}

	roles := make([]*rbac.Role, len(models))
	for i, model := range models {
		roles[i] = r.toDomain(&model)
	}

	return roles, nil
}

// AddPermissionToRole adds a permission to a role
func (r *RoleRepositoryImpl) AddPermissionToRole(ctx context.Context, roleID, permissionID, grantedBy uuid.UUID) error {
	rolePermission := &models.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(rolePermission).Error; err != nil {
		return fmt.Errorf("failed to add permission to role: %w", err)
	}

	return nil
}

// RemovePermissionFromRole removes a permission from a role
func (r *RoleRepositoryImpl) RemovePermissionFromRole(ctx context.Context, roleID, permissionID, removedBy uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&models.RolePermission{}).Error; err != nil {
		return fmt.Errorf("failed to remove permission from role: %w", err)
	}

	return nil
}

// GetRolePermissions retrieves permissions for a specific role
func (r *RoleRepositoryImpl) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*rbac.Permission, error) {
	var models []models.Permission
	if err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}

	permissions := make([]*rbac.Permission, len(models))
	for i, model := range models {
		permissions[i] = &rbac.Permission{
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

	return permissions, nil
}

// toDomain converts a role model to domain entity
func (r *RoleRepositoryImpl) toDomain(model *models.Role) *rbac.Role {
	role := &rbac.Role{
		ID:          model.ID,
		Name:        model.Name,
		DisplayName: model.DisplayName,
		Description: model.Description,
		IsActive:    model.IsActive,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}

	// Convert permissions
	if len(model.Permissions) > 0 {
		permissions := make([]rbac.Permission, len(model.Permissions))
		for i, perm := range model.Permissions {
			permissions[i] = rbac.Permission{
				ID:          perm.ID,
				Name:        perm.DisplayName,
				Resource:    perm.Resource,
				Action:      perm.Action,
				Description: perm.Description,
				IsActive:    perm.IsActive,
				CreatedAt:   perm.CreatedAt,
				UpdatedAt:   perm.UpdatedAt,
			}
		}
		role.Permissions = permissions
	}

	return role
}
