package repositories

import (
	"context"
	"time"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/persistence/postgres/models"
	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
	"gorm.io/gorm"
)

type userRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) auth.UserRoleRepository {
	return &userRoleRepository{
		db: db,
	}
}

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

func (r *userRoleRepository) RevokeRoleFromUser(ctx context.Context, userID, roleID string) error {
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&models.UserRoleModel{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRoleRepository) GetUserRole(ctx context.Context, userID, roleID string) (*auth.UserRole, error) {
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

func (r *userRoleRepository) GetUserRoleByID(ctx context.Context, id string) (*auth.UserRole, error) {
	var userRoleModel models.UserRoleModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&userRoleModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.modelToDomain(&userRoleModel), nil
}

func (r *userRoleRepository) GetUserRoles(ctx context.Context, userID string) ([]*auth.UserRole, error) {
	var userRoleModels []models.UserRoleModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&userRoleModels).Error; err != nil {
		return nil, err
	}

	userRoles := make([]*auth.UserRole, 0, len(userRoleModels))
	for _, model := range userRoleModels {
		userRoles = append(userRoles, r.modelToDomain(&model))
	}

	return userRoles, nil
}

func (r *userRoleRepository) GetActiveUserRoles(ctx context.Context, userID string) ([]*auth.UserRole, error) {
	var userRoleModels []models.UserRoleModel

	query := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Where("expires_at IS NULL OR expires_at > ?", time.Now())

	if err := query.Find(&userRoleModels).Error; err != nil {
		return nil, err
	}

	userRoles := make([]*auth.UserRole, 0, len(userRoleModels))
	for _, model := range userRoleModels {
		userRoles = append(userRoles, r.modelToDomain(&model))
	}

	return userRoles, nil
}

func (r *userRoleRepository) GetRoleUsers(ctx context.Context, roleID string) ([]*auth.UserRole, error) {
	var userRoleModels []models.UserRoleModel
	if err := r.db.WithContext(ctx).Where("role_id = ?", roleID).Find(&userRoleModels).Error; err != nil {
		return nil, err
	}

	userRoles := make([]*auth.UserRole, 0, len(userRoleModels))
	for _, model := range userRoleModels {
		userRoles = append(userRoles, r.modelToDomain(&model))
	}

	return userRoles, nil
}

func (r *userRoleRepository) GetActiveRoleUsers(ctx context.Context, roleID string) ([]*auth.UserRole, error) {
	var userRoleModels []models.UserRoleModel

	query := r.db.WithContext(ctx).
		Where("role_id = ? AND is_active = ?", roleID, true).
		Where("expires_at IS NULL OR expires_at > ?", time.Now())

	if err := query.Find(&userRoleModels).Error; err != nil {
		return nil, err
	}

	userRoles := make([]*auth.UserRole, 0, len(userRoleModels))
	for _, model := range userRoleModels {
		userRoles = append(userRoles, r.modelToDomain(&model))
	}

	return userRoles, nil
}

func (r *userRoleRepository) List(ctx context.Context, params auth.ListUserRolesParams) ([]*auth.UserRole, *pagination.Pagination, error) {
	var userRoleModels []models.UserRoleModel
	var total int64

	query := r.db.WithContext(ctx).Model(&models.UserRoleModel{})

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

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	paginationObj := pagination.NewPagination(params.Page, params.Limit)
	paginationObj.SetTotal(total)

	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortDir := params.SortDir
	if sortDir != "ASC" && sortDir != "DESC" {
		sortDir = "DESC"
	}
	query = query.Order(sortBy + " " + sortDir)

	query = query.Offset(paginationObj.Offset()).Limit(paginationObj.PageSize)

	if err := query.Find(&userRoleModels).Error; err != nil {
		return nil, nil, err
	}

	userRoles := make([]*auth.UserRole, 0, len(userRoleModels))
	for _, model := range userRoleModels {
		userRoles = append(userRoles, r.modelToDomain(&model))
	}

	return userRoles, paginationObj, nil
}

func (r *userRoleRepository) UpdateUserRole(ctx context.Context, userRole *auth.UserRole) error {
	userRoleModel := r.domainToModel(userRole)
	if err := r.db.WithContext(ctx).Model(&userRoleModel).Where("id = ?", userRoleModel.ID).Updates(userRoleModel).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRoleRepository) ActivateUserRole(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Model(&models.UserRoleModel{}).Where("id = ?", id).Update("is_active", true).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRoleRepository) DeactivateUserRole(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Model(&models.UserRoleModel{}).Where("id = ?", id).Update("is_active", false).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRoleRepository) SetExpiration(ctx context.Context, userID, roleID string, expiresAt *time.Time) error {
	if err := r.db.WithContext(ctx).Model(&models.UserRoleModel{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Update("expires_at", expiresAt).Error; err != nil {
		return err
	}
	return nil
}

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

func (r *userRoleRepository) GetExpiredUserRoles(ctx context.Context) ([]*auth.UserRole, error) {
	var userRoleModels []models.UserRoleModel

	query := r.db.WithContext(ctx).
		Where("expires_at IS NOT NULL AND expires_at <= ? AND is_active = ?", time.Now(), true)

	if err := query.Find(&userRoleModels).Error; err != nil {
		return nil, err
	}

	userRoles := make([]*auth.UserRole, 0, len(userRoleModels))
	for _, model := range userRoleModels {
		userRoles = append(userRoles, r.modelToDomain(&model))
	}

	return userRoles, nil
}

func (r *userRoleRepository) CleanupExpiredRoles(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).Model(&models.UserRoleModel{}).
		Where("expires_at IS NOT NULL AND expires_at <= ? AND is_active = ?", time.Now(), true).
		Update("is_active", false)

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

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

func (r *userRoleRepository) Exists(ctx context.Context, userID, roleID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.UserRoleModel{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRoleRepository) domainToModel(userRole *auth.UserRole) *models.UserRoleModel {
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

func (r *userRoleRepository) modelToDomain(userRoleModel *models.UserRoleModel) *auth.UserRole {
	return &auth.UserRole{
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
