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

// RolePermissionRepositoryImpl implements the RBAC RolePermissionRepository interface
type RolePermissionRepositoryImpl struct {
	db *gorm.DB
}

// NewRolePermissionRepository creates a new role permission repository
func NewRolePermissionRepository(db *gorm.DB) rbac.RolePermissionRepository {
	return &RolePermissionRepositoryImpl{
		db: db,
	}
}

// GrantPermissionToRole grants a permission to a role
func (r *RolePermissionRepositoryImpl) GrantPermissionToRole(ctx context.Context, rolePermission *rbac.RolePermission) error {
	model := &models.RolePermission{
		RoleID:       rolePermission.RoleID,
		PermissionID: rolePermission.PermissionID,
		CreatedAt:    rolePermission.GrantedAt,
		UpdatedAt:    time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to grant permission to role: %w", err)
	}

	return nil
}

// RevokePermissionFromRole revokes a permission from a role
func (r *RolePermissionRepositoryImpl) RevokePermissionFromRole(ctx context.Context, roleID, permissionID uuid.UUID, revokedBy uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&models.RolePermission{}).Error; err != nil {
		return fmt.Errorf("failed to revoke permission from role: %w", err)
	}

	return nil
}

// GetRolePermissions retrieves permissions for a role
func (r *RolePermissionRepositoryImpl) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*rbac.RolePermission, error) {
	var rolePermissionModels []models.RolePermission
	if err := r.db.WithContext(ctx).
		Preload("Permission").
		Where("role_id = ?", roleID).
		Find(&rolePermissionModels).Error; err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}

	rolePermissions := make([]*rbac.RolePermission, len(rolePermissionModels))
	for i, model := range rolePermissionModels {
		rolePermissions[i] = r.toDomain(&model)
	}

	return rolePermissions, nil
}

// GetActiveRolePermissions retrieves active permissions for a role
func (r *RolePermissionRepositoryImpl) GetActiveRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*rbac.RolePermission, error) {
	var rolePermissionModels []models.RolePermission
	if err := r.db.WithContext(ctx).
		Preload("Permission").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ? AND permissions.is_active = ?", roleID, true).
		Find(&rolePermissionModels).Error; err != nil {
		return nil, fmt.Errorf("failed to get active role permissions: %w", err)
	}

	rolePermissions := make([]*rbac.RolePermission, len(rolePermissionModels))
	for i, model := range rolePermissionModels {
		rolePermissions[i] = r.toDomain(&model)
	}

	return rolePermissions, nil
}

// GetPermissionRoles retrieves roles for a permission
func (r *RolePermissionRepositoryImpl) GetPermissionRoles(ctx context.Context, permissionID uuid.UUID) ([]*rbac.RolePermission, error) {
	var rolePermissionModels []models.RolePermission
	if err := r.db.WithContext(ctx).
		Preload("Role").
		Where("permission_id = ?", permissionID).
		Find(&rolePermissionModels).Error; err != nil {
		return nil, fmt.Errorf("failed to get permission roles: %w", err)
	}

	rolePermissions := make([]*rbac.RolePermission, len(rolePermissionModels))
	for i, model := range rolePermissionModels {
		rolePermissions[i] = r.toDomain(&model)
	}

	return rolePermissions, nil
}

// HasPermission checks if a role has a specific permission
func (r *RolePermissionRepositoryImpl) HasPermission(ctx context.Context, roleID, permissionID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.RolePermission{}).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check if role has permission: %w", err)
	}

	return count > 0, nil
}

// GrantMultiplePermissionsToRole grants multiple permissions to a role
func (r *RolePermissionRepositoryImpl) GrantMultiplePermissionsToRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID, grantedBy uuid.UUID) error {
	rolePermissions := make([]models.RolePermission, len(permissionIDs))
	now := time.Now()

	for i, permissionID := range permissionIDs {
		rolePermissions[i] = models.RolePermission{
			RoleID:       roleID,
			PermissionID: permissionID,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
	}

	if err := r.db.WithContext(ctx).Create(&rolePermissions).Error; err != nil {
		return fmt.Errorf("failed to grant multiple permissions to role: %w", err)
	}

	return nil
}

// RevokeAllPermissionsFromRole revokes all permissions from a role
func (r *RolePermissionRepositoryImpl) RevokeAllPermissionsFromRole(ctx context.Context, roleID uuid.UUID, revokedBy uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("role_id = ?", roleID).
		Delete(&models.RolePermission{}).Error; err != nil {
		return fmt.Errorf("failed to revoke all permissions from role: %w", err)
	}

	return nil
}

// GetRolesWithPermission retrieves role IDs for a specific permission
func (r *RolePermissionRepositoryImpl) GetRolesWithPermission(ctx context.Context, permissionID uuid.UUID) ([]uuid.UUID, error) {
	var roleIDs []uuid.UUID
	if err := r.db.WithContext(ctx).
		Model(&models.RolePermission{}).
		Where("permission_id = ?", permissionID).
		Pluck("role_id", &roleIDs).Error; err != nil {
		return nil, fmt.Errorf("failed to get roles with permission: %w", err)
	}

	return roleIDs, nil
}

// List retrieves role permissions with pagination
func (r *RolePermissionRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*rbac.RolePermission, error) {
	var rolePermissionModels []models.RolePermission
	query := r.db.WithContext(ctx).Preload("Role").Preload("Permission")

	if offset > 0 {
		query = query.Offset(offset)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&rolePermissionModels).Error; err != nil {
		return nil, fmt.Errorf("failed to list role permissions: %w", err)
	}

	rolePermissions := make([]*rbac.RolePermission, len(rolePermissionModels))
	for i, model := range rolePermissionModels {
		rolePermissions[i] = r.toDomain(&model)
	}

	return rolePermissions, nil
}

// Count returns the total number of role permissions
func (r *RolePermissionRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.RolePermission{}).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count role permissions: %w", err)
	}

	return count, nil
}

// GetRolePermissionHistory retrieves role permission history
func (r *RolePermissionRepositoryImpl) GetRolePermissionHistory(ctx context.Context, roleID uuid.UUID) ([]*rbac.RolePermission, error) {
	// For now, just return current permissions as we don't have history tracking
	// You can implement history tracking by adding a separate table
	return r.GetRolePermissions(ctx, roleID)
}

// toDomain converts a role permission model to domain entity
func (r *RolePermissionRepositoryImpl) toDomain(model *models.RolePermission) *rbac.RolePermission {
	return &rbac.RolePermission{
		RoleID:       model.RoleID,
		PermissionID: model.PermissionID,
		GrantedAt:    model.CreatedAt,
	}
}
