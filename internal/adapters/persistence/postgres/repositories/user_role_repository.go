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

// UserRoleRepositoryImpl implements the RBAC UserRoleRepository interface
type UserRoleRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRoleRepository creates a new user role repository
func NewUserRoleRepository(db *gorm.DB) rbac.UserRoleRepository {
	return &UserRoleRepositoryImpl{
		db: db,
	}
}

// AssignRoleToUser assigns a role to a user
func (r *UserRoleRepositoryImpl) AssignRoleToUser(ctx context.Context, userRole *rbac.UserRole) error {
	model := &models.UserRole{
		UserID:    userRole.UserID,
		RoleID:    userRole.RoleID,
		CreatedAt: userRole.AssignedAt,
		UpdatedAt: time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	return nil
}

// RemoveRoleFromUser removes a role from a user
func (r *UserRoleRepositoryImpl) RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID, removedBy uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&models.UserRole{}).Error; err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}

	return nil
}

// GetUserRoles retrieves all roles for a user
func (r *UserRoleRepositoryImpl) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*rbac.UserRole, error) {
	var userRoleModels []models.UserRole
	if err := r.db.WithContext(ctx).
		Preload("Role").
		Where("user_id = ?", userID).
		Find(&userRoleModels).Error; err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	userRoles := make([]*rbac.UserRole, len(userRoleModels))
	for i, model := range userRoleModels {
		userRoles[i] = r.toDomain(&model)
	}

	return userRoles, nil
}

// GetActiveUserRoles retrieves active roles for a user
func (r *UserRoleRepositoryImpl) GetActiveUserRoles(ctx context.Context, userID uuid.UUID) ([]*rbac.UserRole, error) {
	var userRoleModels []models.UserRole
	if err := r.db.WithContext(ctx).
		Preload("Role").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.is_active = ?", userID, true).
		Find(&userRoleModels).Error; err != nil {
		return nil, fmt.Errorf("failed to get active user roles: %w", err)
	}

	userRoles := make([]*rbac.UserRole, len(userRoleModels))
	for i, model := range userRoleModels {
		userRoles[i] = r.toDomain(&model)
	}

	return userRoles, nil
}

// GetRoleUsers retrieves all users for a role
func (r *UserRoleRepositoryImpl) GetRoleUsers(ctx context.Context, roleID uuid.UUID) ([]*rbac.UserRole, error) {
	var userRoleModels []models.UserRole
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("role_id = ?", roleID).
		Find(&userRoleModels).Error; err != nil {
		return nil, fmt.Errorf("failed to get role users: %w", err)
	}

	userRoles := make([]*rbac.UserRole, len(userRoleModels))
	for i, model := range userRoleModels {
		userRoles[i] = r.toDomain(&model)
	}

	return userRoles, nil
}

// UserHasRole checks if a user has a specific role
func (r *UserRoleRepositoryImpl) UserHasRole(ctx context.Context, userID, roleID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.UserRole{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check if user has role: %w", err)
	}

	return count > 0, nil
}

// HasRole checks if a user has a specific role (alias for UserHasRole)
func (r *UserRoleRepositoryImpl) HasRole(ctx context.Context, userID, roleID uuid.UUID) (bool, error) {
	return r.UserHasRole(ctx, userID, roleID)
}

// IsRoleExpired checks if a user's role is expired
func (r *UserRoleRepositoryImpl) IsRoleExpired(ctx context.Context, userID, roleID uuid.UUID) (bool, error) {
	// Since our model doesn't have expiry, always return false
	// You can implement this if you add expiry functionality
	return false, nil
}

// GetExpiredRoles retrieves expired user roles
func (r *UserRoleRepositoryImpl) GetExpiredRoles(ctx context.Context) ([]*rbac.UserRole, error) {
	// Since our model doesn't have expiry, return empty slice
	// You can implement this if you add expiry functionality
	return []*rbac.UserRole{}, nil
}

// AssignMultipleRolesToUser assigns multiple roles to a user
func (r *UserRoleRepositoryImpl) AssignMultipleRolesToUser(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID, assignedBy uuid.UUID) error {
	return r.BulkAssignRolesToUser(ctx, userID, roleIDs, assignedBy)
}

// RemoveAllRolesFromUser removes all roles from a user
func (r *UserRoleRepositoryImpl) RemoveAllRolesFromUser(ctx context.Context, userID uuid.UUID, removedBy uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.UserRole{}).Error; err != nil {
		return fmt.Errorf("failed to remove all roles from user: %w", err)
	}

	return nil
}

// GetUsersWithRole retrieves user IDs for a specific role
func (r *UserRoleRepositoryImpl) GetUsersWithRole(ctx context.Context, roleID uuid.UUID) ([]uuid.UUID, error) {
	return r.GetUsersByRole(ctx, roleID)
}

// List retrieves user roles with pagination
func (r *UserRoleRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*rbac.UserRole, error) {
	var userRoleModels []models.UserRole
	query := r.db.WithContext(ctx).Preload("Role").Preload("User")

	if offset > 0 {
		query = query.Offset(offset)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&userRoleModels).Error; err != nil {
		return nil, fmt.Errorf("failed to list user roles: %w", err)
	}

	userRoles := make([]*rbac.UserRole, len(userRoleModels))
	for i, model := range userRoleModels {
		userRoles[i] = r.toDomain(&model)
	}

	return userRoles, nil
}

// Count returns the total number of user roles
func (r *UserRoleRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.UserRole{}).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count user roles: %w", err)
	}

	return count, nil
}

// GetUserRoleHistory retrieves user role history
func (r *UserRoleRepositoryImpl) GetUserRoleHistory(ctx context.Context, userID uuid.UUID) ([]*rbac.UserRole, error) {
	// For now, just return current roles as we don't have history tracking
	// You can implement history tracking by adding a separate table
	return r.GetUserRoles(ctx, userID)
}

// UserHasAnyRole checks if a user has any of the specified roles
func (r *UserRoleRepositoryImpl) UserHasAnyRole(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.UserRole{}).
		Where("user_id = ? AND role_id IN ?", userID, roleIDs).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check if user has any role: %w", err)
	}

	return count > 0, nil
}

// GetUsersByRole retrieves user IDs for a specific role
func (r *UserRoleRepositoryImpl) GetUsersByRole(ctx context.Context, roleID uuid.UUID) ([]uuid.UUID, error) {
	var userIDs []uuid.UUID
	if err := r.db.WithContext(ctx).
		Model(&models.UserRole{}).
		Where("role_id = ?", roleID).
		Pluck("user_id", &userIDs).Error; err != nil {
		return nil, fmt.Errorf("failed to get users by role: %w", err)
	}

	return userIDs, nil
}

// BulkAssignRolesToUser assigns multiple roles to a user
func (r *UserRoleRepositoryImpl) BulkAssignRolesToUser(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID, assignedBy uuid.UUID) error {
	userRoles := make([]models.UserRole, len(roleIDs))
	now := time.Now()

	for i, roleID := range roleIDs {
		userRoles[i] = models.UserRole{
			UserID:    userID,
			RoleID:    roleID,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	if err := r.db.WithContext(ctx).Create(&userRoles).Error; err != nil {
		return fmt.Errorf("failed to bulk assign roles to user: %w", err)
	}

	return nil
}

// BulkRemoveRolesFromUser removes multiple roles from a user
func (r *UserRoleRepositoryImpl) BulkRemoveRolesFromUser(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID, removedBy uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND role_id IN ?", userID, roleIDs).
		Delete(&models.UserRole{}).Error; err != nil {
		return fmt.Errorf("failed to bulk remove roles from user: %w", err)
	}

	return nil
}

// toDomain converts a user role model to domain entity
func (r *UserRoleRepositoryImpl) toDomain(model *models.UserRole) *rbac.UserRole {
	return &rbac.UserRole{
		UserID:     model.UserID,
		RoleID:     model.RoleID,
		AssignedAt: model.CreatedAt,
	}
}
