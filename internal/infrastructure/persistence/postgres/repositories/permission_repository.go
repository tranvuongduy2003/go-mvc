package repositories

import (
	"context"
	"strings"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/persistence/postgres/models"
	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) auth.PermissionRepository {
	return &permissionRepository{
		db: db,
	}
}

func (r *permissionRepository) Create(ctx context.Context, permEntity *auth.Permission) error {
	permModel := r.domainToModel(permEntity)
	if err := r.db.WithContext(ctx).Create(permModel).Error; err != nil {
		return err
	}
	return nil
}

func (r *permissionRepository) GetByID(ctx context.Context, id string) (*auth.Permission, error) {
	var permModel models.PermissionModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&permModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.modelToDomain(&permModel)
}

func (r *permissionRepository) GetByName(ctx context.Context, name string) (*auth.Permission, error) {
	var permModel models.PermissionModel
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&permModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.modelToDomain(&permModel)
}

func (r *permissionRepository) GetByResourceAndAction(ctx context.Context, resource, action string) (*auth.Permission, error) {
	var permModel models.PermissionModel
	if err := r.db.WithContext(ctx).Where("resource = ? AND action = ?", resource, action).First(&permModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.modelToDomain(&permModel)
}

func (r *permissionRepository) Update(ctx context.Context, permEntity *auth.Permission) error {
	permModel := r.domainToModel(permEntity)
	if err := r.db.WithContext(ctx).Model(&permModel).Where("id = ?", permModel.ID).Updates(permModel).Error; err != nil {
		return err
	}
	return nil
}

func (r *permissionRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.PermissionModel{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *permissionRepository) List(ctx context.Context, params auth.ListPermissionsParams) ([]*auth.Permission, *pagination.Pagination, error) {
	var permModels []models.PermissionModel
	var total int64

	query := r.db.WithContext(ctx).Model(&models.PermissionModel{})

	if params.Search != "" {
		searchPattern := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(resource) LIKE ? OR LOWER(action) LIKE ? OR LOWER(description) LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}

	if params.Resource != "" {
		query = query.Where("resource = ?", params.Resource)
	}

	if params.Action != "" {
		query = query.Where("action = ?", params.Action)
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
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
	sortDir := strings.ToUpper(params.SortDir)
	if sortDir != "ASC" && sortDir != "DESC" {
		sortDir = "DESC"
	}
	query = query.Order(sortBy + " " + sortDir)

	query = query.Offset(paginationObj.Offset()).Limit(paginationObj.PageSize)

	if err := query.Find(&permModels).Error; err != nil {
		return nil, nil, err
	}

	permissions := make([]*auth.Permission, 0, len(permModels))
	for _, permModel := range permModels {
		permEntity, err := r.modelToDomain(&permModel)
		if err != nil {
			return nil, nil, err
		}
		permissions = append(permissions, permEntity)
	}

	return permissions, paginationObj, nil
}

func (r *permissionRepository) GetActivePermissions(ctx context.Context) ([]*auth.Permission, error) {
	var permModels []models.PermissionModel
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&permModels).Error; err != nil {
		return nil, err
	}

	permissions := make([]*auth.Permission, 0, len(permModels))
	for _, permModel := range permModels {
		permEntity, err := r.modelToDomain(&permModel)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permEntity)
	}

	return permissions, nil
}

func (r *permissionRepository) GetPermissionsByResource(ctx context.Context, resource string) ([]*auth.Permission, error) {
	var permModels []models.PermissionModel
	if err := r.db.WithContext(ctx).Where("resource = ? AND is_active = ?", resource, true).Find(&permModels).Error; err != nil {
		return nil, err
	}

	permissions := make([]*auth.Permission, 0, len(permModels))
	for _, permModel := range permModels {
		permEntity, err := r.modelToDomain(&permModel)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permEntity)
	}

	return permissions, nil
}

func (r *permissionRepository) GetPermissionsByAction(ctx context.Context, action string) ([]*auth.Permission, error) {
	var permModels []models.PermissionModel
	if err := r.db.WithContext(ctx).Where("action = ? AND is_active = ?", action, true).Find(&permModels).Error; err != nil {
		return nil, err
	}

	permissions := make([]*auth.Permission, 0, len(permModels))
	for _, permModel := range permModels {
		permEntity, err := r.modelToDomain(&permModel)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permEntity)
	}

	return permissions, nil
}

func (r *permissionRepository) Exists(ctx context.Context, id string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.PermissionModel{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *permissionRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.PermissionModel{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *permissionRepository) ExistsByResourceAndAction(ctx context.Context, resource, action string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.PermissionModel{}).Where("resource = ? AND action = ?", resource, action).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *permissionRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.PermissionModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *permissionRepository) Activate(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Model(&models.PermissionModel{}).Where("id = ?", id).Update("is_active", true).Error; err != nil {
		return err
	}
	return nil
}

func (r *permissionRepository) Deactivate(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Model(&models.PermissionModel{}).Where("id = ?", id).Update("is_active", false).Error; err != nil {
		return err
	}
	return nil
}

func (r *permissionRepository) GetPermissionsByUserID(ctx context.Context, userID string) ([]*auth.Permission, error) {
	var permModels []models.PermissionModel

	query := r.db.WithContext(ctx).
		Model(&models.PermissionModel{}).
		Joins("INNER JOIN role_permissions rp ON permissions.id = rp.permission_id").
		Joins("INNER JOIN user_roles ur ON rp.role_id = ur.role_id").
		Where("ur.user_id = ?", userID)

	if err := query.Find(&permModels).Error; err != nil {
		return nil, err
	}

	permissions := make([]*auth.Permission, 0, len(permModels))
	for _, permModel := range permModels {
		permEntity, err := r.modelToDomain(&permModel)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permEntity)
	}

	return permissions, nil
}

func (r *permissionRepository) GetActivePermissionsByUserID(ctx context.Context, userID string) ([]*auth.Permission, error) {
	var permModels []models.PermissionModel

	query := r.db.WithContext(ctx).
		Model(&models.PermissionModel{}).
		Joins("INNER JOIN role_permissions rp ON permissions.id = rp.permission_id").
		Joins("INNER JOIN roles r ON rp.role_id = r.id").
		Joins("INNER JOIN user_roles ur ON r.id = ur.role_id").
		Where("ur.user_id = ? AND ur.is_active = ? AND r.is_active = ? AND rp.is_active = ? AND permissions.is_active = ?",
			userID, true, true, true, true).
		Where("ur.expires_at IS NULL OR ur.expires_at > CURRENT_TIMESTAMP")

	if err := query.Find(&permModels).Error; err != nil {
		return nil, err
	}

	permissions := make([]*auth.Permission, 0, len(permModels))
	for _, permModel := range permModels {
		permEntity, err := r.modelToDomain(&permModel)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permEntity)
	}

	return permissions, nil
}

func (r *permissionRepository) GetPermissionsByRoleID(ctx context.Context, roleID string) ([]*auth.Permission, error) {
	var permModels []models.PermissionModel

	query := r.db.WithContext(ctx).
		Model(&models.PermissionModel{}).
		Joins("INNER JOIN role_permissions rp ON permissions.id = rp.permission_id").
		Where("rp.role_id = ?", roleID)

	if err := query.Find(&permModels).Error; err != nil {
		return nil, err
	}

	permissions := make([]*auth.Permission, 0, len(permModels))
	for _, permModel := range permModels {
		permEntity, err := r.modelToDomain(&permModel)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permEntity)
	}

	return permissions, nil
}

func (r *permissionRepository) GetActivePermissionsByRoleID(ctx context.Context, roleID string) ([]*auth.Permission, error) {
	var permModels []models.PermissionModel

	query := r.db.WithContext(ctx).
		Model(&models.PermissionModel{}).
		Joins("INNER JOIN role_permissions rp ON permissions.id = rp.permission_id").
		Where("rp.role_id = ? AND rp.is_active = ? AND permissions.is_active = ?", roleID, true, true)

	if err := query.Find(&permModels).Error; err != nil {
		return nil, err
	}

	permissions := make([]*auth.Permission, 0, len(permModels))
	for _, permModel := range permModels {
		permEntity, err := r.modelToDomain(&permModel)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permEntity)
	}

	return permissions, nil
}

func (r *permissionRepository) UserHasPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	var count int64

	query := r.db.WithContext(ctx).
		Model(&models.PermissionModel{}).
		Joins("INNER JOIN role_permissions rp ON permissions.id = rp.permission_id").
		Joins("INNER JOIN roles r ON rp.role_id = r.id").
		Joins("INNER JOIN user_roles ur ON r.id = ur.role_id").
		Where("ur.user_id = ? AND permissions.resource = ? AND permissions.action = ?", userID, resource, action).
		Where("ur.is_active = ? AND r.is_active = ? AND rp.is_active = ? AND permissions.is_active = ?", true, true, true, true).
		Where("ur.expires_at IS NULL OR ur.expires_at > CURRENT_TIMESTAMP")

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *permissionRepository) UserHasPermissionByName(ctx context.Context, userID, permissionName string) (bool, error) {
	var count int64

	query := r.db.WithContext(ctx).
		Model(&models.PermissionModel{}).
		Joins("INNER JOIN role_permissions rp ON permissions.id = rp.permission_id").
		Joins("INNER JOIN roles r ON rp.role_id = r.id").
		Joins("INNER JOIN user_roles ur ON r.id = ur.role_id").
		Where("ur.user_id = ? AND permissions.name = ?", userID, permissionName).
		Where("ur.is_active = ? AND r.is_active = ? AND rp.is_active = ? AND permissions.is_active = ?", true, true, true, true).
		Where("ur.expires_at IS NULL OR ur.expires_at > CURRENT_TIMESTAMP")

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *permissionRepository) domainToModel(permEntity *auth.Permission) *models.PermissionModel {
	return &models.PermissionModel{
		ID:          permEntity.ID().String(),
		Name:        permEntity.Name().String(),
		Resource:    permEntity.Resource().String(),
		Action:      permEntity.Action().String(),
		Description: permEntity.Description(),
		IsActive:    permEntity.IsActive(),
		CreatedAt:   permEntity.CreatedAt(),
		UpdatedAt:   permEntity.UpdatedAt(),
		Version:     permEntity.Version(),
	}
}

func (r *permissionRepository) modelToDomain(permModel *models.PermissionModel) (*auth.Permission, error) {
	return auth.ReconstructPermission(
		permModel.ID,
		permModel.Name,
		permModel.Resource,
		permModel.Action,
		permModel.Description,
		permModel.IsActive,
		permModel.CreatedAt,
		permModel.UpdatedAt,
		permModel.Version,
	)
}
