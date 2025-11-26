package repositories

import (
	"context"
	"time"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/persistence/postgres/models"
	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
	"gorm.io/gorm"
)

type rolePermissionRepository struct {
	db *gorm.DB
}

func NewRolePermissionRepository(db *gorm.DB) auth.RolePermissionRepository {
	return &rolePermissionRepository{
		db: db,
	}
}

func (r *rolePermissionRepository) GrantPermissionToRole(ctx context.Context, roleID, permissionID string, grantedBy *string) error {
	rolePermModel := &models.RolePermissionModel{
		RoleID:       roleID,
		PermissionID: permissionID,
		GrantedBy:    grantedBy,
		GrantedAt:    time.Now(),
		IsActive:     true,
	}

	if err := r.db.WithContext(ctx).Create(rolePermModel).Error; err != nil {
		return err
	}
	return nil
}

func (r *rolePermissionRepository) RevokePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	if err := r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&models.RolePermissionModel{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *rolePermissionRepository) GetRolePermission(ctx context.Context, roleID, permissionID string) (*auth.RolePermission, error) {
	var rolePermModel models.RolePermissionModel
	if err := r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		First(&rolePermModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.modelToDomain(&rolePermModel), nil
}

func (r *rolePermissionRepository) GetRolePermissionByID(ctx context.Context, id string) (*auth.RolePermission, error) {
	var rolePermModel models.RolePermissionModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&rolePermModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.modelToDomain(&rolePermModel), nil
}

func (r *rolePermissionRepository) GetRolePermissions(ctx context.Context, roleID string) ([]*auth.RolePermission, error) {
	var rolePermModels []models.RolePermissionModel
	if err := r.db.WithContext(ctx).Where("role_id = ?", roleID).Find(&rolePermModels).Error; err != nil {
		return nil, err
	}

	rolePermissions := make([]*auth.RolePermission, 0, len(rolePermModels))
	for _, model := range rolePermModels {
		rolePermissions = append(rolePermissions, r.modelToDomain(&model))
	}

	return rolePermissions, nil
}

func (r *rolePermissionRepository) GetActiveRolePermissions(ctx context.Context, roleID string) ([]*auth.RolePermission, error) {
	var rolePermModels []models.RolePermissionModel

	query := r.db.WithContext(ctx).Where("role_id = ? AND is_active = ?", roleID, true)

	if err := query.Find(&rolePermModels).Error; err != nil {
		return nil, err
	}

	rolePermissions := make([]*auth.RolePermission, 0, len(rolePermModels))
	for _, model := range rolePermModels {
		rolePermissions = append(rolePermissions, r.modelToDomain(&model))
	}

	return rolePermissions, nil
}

func (r *rolePermissionRepository) GetPermissionRoles(ctx context.Context, permissionID string) ([]*auth.RolePermission, error) {
	var rolePermModels []models.RolePermissionModel
	if err := r.db.WithContext(ctx).Where("permission_id = ?", permissionID).Find(&rolePermModels).Error; err != nil {
		return nil, err
	}

	rolePermissions := make([]*auth.RolePermission, 0, len(rolePermModels))
	for _, model := range rolePermModels {
		rolePermissions = append(rolePermissions, r.modelToDomain(&model))
	}

	return rolePermissions, nil
}

func (r *rolePermissionRepository) GetActivePermissionRoles(ctx context.Context, permissionID string) ([]*auth.RolePermission, error) {
	var rolePermModels []models.RolePermissionModel

	query := r.db.WithContext(ctx).Where("permission_id = ? AND is_active = ?", permissionID, true)

	if err := query.Find(&rolePermModels).Error; err != nil {
		return nil, err
	}

	rolePermissions := make([]*auth.RolePermission, 0, len(rolePermModels))
	for _, model := range rolePermModels {
		rolePermissions = append(rolePermissions, r.modelToDomain(&model))
	}

	return rolePermissions, nil
}

func (r *rolePermissionRepository) List(ctx context.Context, params auth.ListRolePermissionsParams) ([]*auth.RolePermission, *pagination.Pagination, error) {
	var rolePermModels []models.RolePermissionModel
	var total int64

	query := r.db.WithContext(ctx).Model(&models.RolePermissionModel{})

	if params.RoleID != "" {
		query = query.Where("role_id = ?", params.RoleID)
	}

	if params.PermissionID != "" {
		query = query.Where("permission_id = ?", params.PermissionID)
	}

	if params.GrantedBy != "" {
		query = query.Where("granted_by = ?", params.GrantedBy)
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if params.Resource != "" || params.Action != "" {
		query = query.Joins("INNER JOIN permissions ON role_permissions.permission_id = permissions.id")

		if params.Resource != "" {
			query = query.Where("permissions.resource = ?", params.Resource)
		}

		if params.Action != "" {
			query = query.Where("permissions.action = ?", params.Action)
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

	if err := query.Find(&rolePermModels).Error; err != nil {
		return nil, nil, err
	}

	rolePermissions := make([]*auth.RolePermission, 0, len(rolePermModels))
	for _, model := range rolePermModels {
		rolePermissions = append(rolePermissions, r.modelToDomain(&model))
	}

	return rolePermissions, paginationObj, nil
}

func (r *rolePermissionRepository) UpdateRolePermission(ctx context.Context, rolePermission *auth.RolePermission) error {
	rolePermModel := r.domainToModel(rolePermission)
	if err := r.db.WithContext(ctx).Model(&rolePermModel).Where("id = ?", rolePermModel.ID).Updates(rolePermModel).Error; err != nil {
		return err
	}
	return nil
}

func (r *rolePermissionRepository) ActivateRolePermission(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Model(&models.RolePermissionModel{}).Where("id = ?", id).Update("is_active", true).Error; err != nil {
		return err
	}
	return nil
}

func (r *rolePermissionRepository) DeactivateRolePermission(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Model(&models.RolePermissionModel{}).Where("id = ?", id).Update("is_active", false).Error; err != nil {
		return err
	}
	return nil
}

func (r *rolePermissionRepository) RoleHasPermission(ctx context.Context, roleID, permissionID string) (bool, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&models.RolePermissionModel{}).
		Where("role_id = ? AND permission_id = ? AND is_active = ?", roleID, permissionID, true)

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *rolePermissionRepository) RoleHasPermissionByName(ctx context.Context, roleID, permissionName string) (bool, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&models.RolePermissionModel{}).
		Joins("INNER JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ? AND permissions.name = ? AND role_permissions.is_active = ? AND permissions.is_active = ?",
			roleID, permissionName, true, true)

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *rolePermissionRepository) RoleHasResourceAction(ctx context.Context, roleID, resource, action string) (bool, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&models.RolePermissionModel{}).
		Joins("INNER JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ? AND permissions.resource = ? AND permissions.action = ?", roleID, resource, action).
		Where("role_permissions.is_active = ? AND permissions.is_active = ?", true, true)

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *rolePermissionRepository) CountPermissionsByRole(ctx context.Context, roleID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.RolePermissionModel{}).
		Where("role_id = ? AND is_active = ?", roleID, true).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *rolePermissionRepository) CountRolesByPermission(ctx context.Context, permissionID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.RolePermissionModel{}).
		Where("permission_id = ? AND is_active = ?", permissionID, true).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *rolePermissionRepository) Exists(ctx context.Context, roleID, permissionID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.RolePermissionModel{}).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *rolePermissionRepository) BulkGrantPermissions(ctx context.Context, roleID string, permissionIDs []string, grantedBy *string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		for _, permissionID := range permissionIDs {
			var count int64
			if err := tx.Model(&models.RolePermissionModel{}).
				Where("role_id = ? AND permission_id = ?", roleID, permissionID).
				Count(&count).Error; err != nil {
				return err
			}

			if count > 0 {
				continue
			}

			rolePermModel := &models.RolePermissionModel{
				RoleID:       roleID,
				PermissionID: permissionID,
				GrantedBy:    grantedBy,
				GrantedAt:    now,
				IsActive:     true,
			}

			if err := tx.Create(rolePermModel).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *rolePermissionRepository) BulkRevokePermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, permissionID := range permissionIDs {
			if err := tx.Where("role_id = ? AND permission_id = ?", roleID, permissionID).
				Delete(&models.RolePermissionModel{}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *rolePermissionRepository) SyncRolePermissions(ctx context.Context, roleID string, permissionIDs []string, grantedBy *string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var currentRolePerms []models.RolePermissionModel
		if err := tx.Where("role_id = ?", roleID).Find(&currentRolePerms).Error; err != nil {
			return err
		}

		currentPermissionIDs := make(map[string]bool)
		for _, rp := range currentRolePerms {
			currentPermissionIDs[rp.PermissionID] = true
		}

		desiredPermissionIDs := make(map[string]bool)
		for _, pid := range permissionIDs {
			desiredPermissionIDs[pid] = true
		}

		now := time.Now()

		for _, permissionID := range permissionIDs {
			if !currentPermissionIDs[permissionID] {
				rolePermModel := &models.RolePermissionModel{
					RoleID:       roleID,
					PermissionID: permissionID,
					GrantedBy:    grantedBy,
					GrantedAt:    now,
					IsActive:     true,
				}

				if err := tx.Create(rolePermModel).Error; err != nil {
					return err
				}
			}
		}

		for _, rp := range currentRolePerms {
			if !desiredPermissionIDs[rp.PermissionID] {
				if err := tx.Where("id = ?", rp.ID).Delete(&models.RolePermissionModel{}).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *rolePermissionRepository) GetRolePermissionsByResource(ctx context.Context, roleID, resource string) ([]*auth.RolePermission, error) {
	var rolePermModels []models.RolePermissionModel

	query := r.db.WithContext(ctx).Model(&models.RolePermissionModel{}).
		Joins("INNER JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ? AND permissions.resource = ? AND role_permissions.is_active = ?",
			roleID, resource, true)

	if err := query.Find(&rolePermModels).Error; err != nil {
		return nil, err
	}

	rolePermissions := make([]*auth.RolePermission, 0, len(rolePermModels))
	for _, model := range rolePermModels {
		rolePermissions = append(rolePermissions, r.modelToDomain(&model))
	}

	return rolePermissions, nil
}

func (r *rolePermissionRepository) GetRolePermissionsByAction(ctx context.Context, roleID, action string) ([]*auth.RolePermission, error) {
	var rolePermModels []models.RolePermissionModel

	query := r.db.WithContext(ctx).Model(&models.RolePermissionModel{}).
		Joins("INNER JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ? AND permissions.action = ? AND role_permissions.is_active = ?",
			roleID, action, true)

	if err := query.Find(&rolePermModels).Error; err != nil {
		return nil, err
	}

	rolePermissions := make([]*auth.RolePermission, 0, len(rolePermModels))
	for _, model := range rolePermModels {
		rolePermissions = append(rolePermissions, r.modelToDomain(&model))
	}

	return rolePermissions, nil
}

func (r *rolePermissionRepository) domainToModel(rolePerm *auth.RolePermission) *models.RolePermissionModel {
	return &models.RolePermissionModel{
		ID:           rolePerm.ID,
		RoleID:       rolePerm.RoleID,
		PermissionID: rolePerm.PermissionID,
		GrantedBy:    rolePerm.GrantedBy,
		GrantedAt:    rolePerm.GrantedAt,
		IsActive:     rolePerm.IsActive,
		CreatedAt:    rolePerm.CreatedAt,
		UpdatedAt:    rolePerm.UpdatedAt,
		Version:      rolePerm.Version,
	}
}

func (r *rolePermissionRepository) modelToDomain(rolePermModel *models.RolePermissionModel) *auth.RolePermission {
	return &auth.RolePermission{
		ID:           rolePermModel.ID,
		RoleID:       rolePermModel.RoleID,
		PermissionID: rolePermModel.PermissionID,
		GrantedBy:    rolePermModel.GrantedBy,
		GrantedAt:    rolePermModel.GrantedAt,
		IsActive:     rolePermModel.IsActive,
		CreatedAt:    rolePermModel.CreatedAt,
		UpdatedAt:    rolePermModel.UpdatedAt,
		Version:      rolePermModel.Version,
	}
}
