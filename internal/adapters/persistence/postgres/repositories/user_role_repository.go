package repositories

import (
	"context"
	"time"

	"github.com/tranvuongduy2003/go-mvc/internal/adapters/persistence/postgres/models"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/repositories"
	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
	"gorm.io/gorm"
)

// userRoleRepository implements the UserRoleRepository interface
type userRoleRepository struct {
	db *gorm.DB
}

// NewUserRoleRepository creates a new UserRoleRepository instance
func NewUserRoleRepository(db *gorm.DB) repositories.UserRoleRepository {
	return &userRoleRepository{
		db: db,
	}
}

// AssignRoleToUser assigns a role to a user
func (r *userRoleRepository) AssignRoleToUser(ctx context.Context, userID, roleID string, assignedBy *string, expiresAt *time.Time) error {
	userRoleModel := &models.UserRoleModel{
		UserID:     userID,
		RoleID:     roleID,
		AssignedBy: assignedBy,
		AssignedAt: time.Now(),
		ExpiresAt:  expiresAt,
		IsActive:   true,
	}

	if err := r.db.WithContext(ctx).Create(userRoleModel).Error; err != nil {
		return err
	}
	return nil
}

// RevokeRoleFromUser revokes a role from a user
func (r *userRoleRepository) RevokeRoleFromUser(ctx context.Context, userID, roleID string) error {
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&models.UserRoleModel{}).Error; err != nil {
		return err
	}
	return nil
}

// GetUserRole retrieves a specific user-role assignment
func (r *userRoleRepository) GetUserRole(ctx context.Context, userID, roleID string) (*repositories.UserRole, error) {
	var userRoleModel models.UserRoleModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		First(&userRoleModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.modelToDomain(&userRoleModel), nil
}

// GetUserRoleByID retrieves a user-role assignment by ID
func (r *userRoleRepository) GetUserRoleByID(ctx context.Context, id string) (*repositories.UserRole, error) {
	var userRoleModel models.UserRoleModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&userRoleModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.modelToDomain(&userRoleModel), nil
}

// GetUserRoles retrieves all role assignments for a user
func (r *userRoleRepository) GetUserRoles(ctx context.Context, userID string) ([]*repositories.UserRole, error) {
	var userRoleModels []models.UserRoleModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&userRoleModels).Error; err != nil {
		return nil, err
	}

	userRoles := make([]*repositories.UserRole, 0, len(userRoleModels))
	for _, model := range userRoleModels {
		userRoles = append(userRoles, r.modelToDomain(&model))
	}

	return userRoles, nil
}

// GetActiveUserRoles retrieves all active role assignments for a user
func (r *userRoleRepository) GetActiveUserRoles(ctx context.Context, userID string) ([]*repositories.UserRole, error) {
	var userRoleModels []models.UserRoleModel

	query := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Where("expires_at IS NULL OR expires_at > ?", time.Now())

	if err := query.Find(&userRoleModels).Error; err != nil {
		return nil, err
	}

	userRoles := make([]*repositories.UserRole, 0, len(userRoleModels))
	for _, model := range userRoleModels {
		userRoles = append(userRoles, r.modelToDomain(&model))
	}

	return userRoles, nil
}

// GetRoleUsers retrieves all user assignments for a role
func (r *userRoleRepository) GetRoleUsers(ctx context.Context, roleID string) ([]*repositories.UserRole, error) {
	var userRoleModels []models.UserRoleModel
	if err := r.db.WithContext(ctx).Where("role_id = ?", roleID).Find(&userRoleModels).Error; err != nil {
		return nil, err
	}

	userRoles := make([]*repositories.UserRole, 0, len(userRoleModels))
	for _, model := range userRoleModels {
		userRoles = append(userRoles, r.modelToDomain(&model))
	}

	return userRoles, nil
}

// GetActiveRoleUsers retrieves all active user assignments for a role
func (r *userRoleRepository) GetActiveRoleUsers(ctx context.Context, roleID string) ([]*repositories.UserRole, error) {
	var userRoleModels []models.UserRoleModel

	query := r.db.WithContext(ctx).
		Where("role_id = ? AND is_active = ?", roleID, true).
		Where("expires_at IS NULL OR expires_at > ?", time.Now())

	if err := query.Find(&userRoleModels).Error; err != nil {
		return nil, err
	}

	userRoles := make([]*repositories.UserRole, 0, len(userRoleModels))
	for _, model := range userRoleModels {
		userRoles = append(userRoles, r.modelToDomain(&model))
	}

	return userRoles, nil
}

// List retrieves a paginated list of user-role assignments
func (r *userRoleRepository) List(ctx context.Context, params repositories.ListUserRolesParams) ([]*repositories.UserRole, *pagination.Pagination, error) {
	var userRoleModels []models.UserRoleModel
	var total int64

	query := r.db.WithContext(ctx).Model(&models.UserRoleModel{})

	// Apply filters
	if params.UserID != "" {
		query = query.Where("user_id = ?", params.UserID)
	}

	if params.RoleID != "" {
		query = query.Where("role_id = ?", params.RoleID)
	}

	if params.AssignedBy != "" {
		query = query.Where("assigned_by = ?", params.AssignedBy)
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if params.IsExpired != nil {
		if *params.IsExpired {
			query = query.Where("expires_at IS NOT NULL AND expires_at <= ?", time.Now())
		} else {
			query = query.Where("expires_at IS NULL OR expires_at > ?", time.Now())
		}
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Create pagination object
	paginationObj := pagination.NewPagination(params.Page, params.Limit)
	paginationObj.SetTotal(total)

	// Apply sorting
	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortDir := params.SortDir
	if sortDir != "ASC" && sortDir != "DESC" {
		sortDir = "DESC"
	}
	query = query.Order(sortBy + " " + sortDir)

	// Apply pagination
	query = query.Offset(paginationObj.Offset()).Limit(paginationObj.PageSize)

	// Execute query
	if err := query.Find(&userRoleModels).Error; err != nil {
		return nil, nil, err
	}

	// Convert models to domain entities
	userRoles := make([]*repositories.UserRole, 0, len(userRoleModels))
	for _, model := range userRoleModels {
		userRoles = append(userRoles, r.modelToDomain(&model))
	}

	return userRoles, paginationObj, nil
}

// UpdateUserRole updates a user-role assignment
func (r *userRoleRepository) UpdateUserRole(ctx context.Context, userRole *repositories.UserRole) error {
	userRoleModel := r.domainToModel(userRole)
	if err := r.db.WithContext(ctx).Model(&userRoleModel).Where("id = ?", userRoleModel.ID).Updates(userRoleModel).Error; err != nil {
		return err
	}
	return nil
}

// ActivateUserRole activates a user-role assignment
func (r *userRoleRepository) ActivateUserRole(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Model(&models.UserRoleModel{}).Where("id = ?", id).Update("is_active", true).Error; err != nil {
		return err
	}
	return nil
}

// DeactivateUserRole deactivates a user-role assignment
func (r *userRoleRepository) DeactivateUserRole(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Model(&models.UserRoleModel{}).Where("id = ?", id).Update("is_active", false).Error; err != nil {
		return err
	}
	return nil
}

// SetExpiration sets or updates expiration for a user-role assignment
func (r *userRoleRepository) SetExpiration(ctx context.Context, userID, roleID string, expiresAt *time.Time) error {
	if err := r.db.WithContext(ctx).Model(&models.UserRoleModel{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Update("expires_at", expiresAt).Error; err != nil {
		return err
	}
	return nil
}

// IsUserRoleExpired checks if a user-role assignment is expired
func (r *userRoleRepository) IsUserRoleExpired(ctx context.Context, userID, roleID string) (bool, error) {
	var userRoleModel models.UserRoleModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		First(&userRoleModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil // No assignment means not expired
		}
		return false, err
	}

	if userRoleModel.ExpiresAt == nil {
		return false, nil // No expiration set
	}

	return userRoleModel.ExpiresAt.Before(time.Now()), nil
}

// GetExpiredUserRoles retrieves all expired user-role assignments
func (r *userRoleRepository) GetExpiredUserRoles(ctx context.Context) ([]*repositories.UserRole, error) {
	var userRoleModels []models.UserRoleModel

	query := r.db.WithContext(ctx).
		Where("expires_at IS NOT NULL AND expires_at <= ? AND is_active = ?", time.Now(), true)

	if err := query.Find(&userRoleModels).Error; err != nil {
		return nil, err
	}

	userRoles := make([]*repositories.UserRole, 0, len(userRoleModels))
	for _, model := range userRoleModels {
		userRoles = append(userRoles, r.modelToDomain(&model))
	}

	return userRoles, nil
}

// CleanupExpiredRoles deactivates all expired user-role assignments
func (r *userRoleRepository) CleanupExpiredRoles(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).Model(&models.UserRoleModel{}).
		Where("expires_at IS NOT NULL AND expires_at <= ? AND is_active = ?", time.Now(), true).
		Update("is_active", false)

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

// UserHasRole checks if a user currently has a specific role (active and not expired)
func (r *userRoleRepository) UserHasRole(ctx context.Context, userID, roleID string) (bool, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&models.UserRoleModel{}).
		Where("user_id = ? AND role_id = ? AND is_active = ?", userID, roleID, true).
		Where("expires_at IS NULL OR expires_at > ?", time.Now())

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// UserHasRoleName checks if a user currently has a specific role by name
func (r *userRoleRepository) UserHasRoleName(ctx context.Context, userID, roleName string) (bool, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&models.UserRoleModel{}).
		Joins("INNER JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.name = ? AND user_roles.is_active = ? AND roles.is_active = ?",
			userID, roleName, true, true).
		Where("user_roles.expires_at IS NULL OR user_roles.expires_at > ?", time.Now())

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// CountUsersByRole counts users assigned to a specific role
func (r *userRoleRepository) CountUsersByRole(ctx context.Context, roleID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.UserRoleModel{}).
		Where("role_id = ? AND is_active = ?", roleID, true).
		Where("expires_at IS NULL OR expires_at > ?", time.Now()).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountRolesByUser counts roles assigned to a specific user
func (r *userRoleRepository) CountRolesByUser(ctx context.Context, userID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.UserRoleModel{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Where("expires_at IS NULL OR expires_at > ?", time.Now()).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Exists checks if a user-role assignment exists
func (r *userRoleRepository) Exists(ctx context.Context, userID, roleID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.UserRoleModel{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// domainToModel converts domain entity to GORM model
func (r *userRoleRepository) domainToModel(userRole *repositories.UserRole) *models.UserRoleModel {
	return &models.UserRoleModel{
		ID:         userRole.ID,
		UserID:     userRole.UserID,
		RoleID:     userRole.RoleID,
		AssignedBy: userRole.AssignedBy,
		AssignedAt: userRole.AssignedAt,
		ExpiresAt:  userRole.ExpiresAt,
		IsActive:   userRole.IsActive,
		CreatedAt:  userRole.CreatedAt,
		UpdatedAt:  userRole.UpdatedAt,
		Version:    userRole.Version,
	}
}

// modelToDomain converts GORM model to domain entity
func (r *userRoleRepository) modelToDomain(userRoleModel *models.UserRoleModel) *repositories.UserRole {
	return &repositories.UserRole{
		ID:         userRoleModel.ID,
		UserID:     userRoleModel.UserID,
		RoleID:     userRoleModel.RoleID,
		AssignedBy: userRoleModel.AssignedBy,
		AssignedAt: userRoleModel.AssignedAt,
		ExpiresAt:  userRoleModel.ExpiresAt,
		IsActive:   userRoleModel.IsActive,
		CreatedAt:  userRoleModel.CreatedAt,
		UpdatedAt:  userRoleModel.UpdatedAt,
		Version:    userRoleModel.Version,
	}
}
