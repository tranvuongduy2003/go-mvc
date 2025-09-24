package repositories

import (
	"context"
	"strings"

	"github.com/tranvuongduy2003/go-mvc/internal/adapters/persistence/postgres/models"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/permission"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/repositories"
	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
	"gorm.io/gorm"
)

// permissionRepository implements the PermissionRepository interface
type permissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository creates a new PermissionRepository instance
func NewPermissionRepository(db *gorm.DB) repositories.PermissionRepository {
	return &permissionRepository{
		db: db,
	}
}

// Create saves a new permission to the database
func (r *permissionRepository) Create(ctx context.Context, permEntity *permission.Permission) error {
	permModel := r.domainToModel(permEntity)
	if err := r.db.WithContext(ctx).Create(permModel).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a permission by ID
func (r *permissionRepository) GetByID(ctx context.Context, id string) (*permission.Permission, error) {
	var permModel models.PermissionModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&permModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.modelToDomain(&permModel)
}

// GetByName retrieves a permission by name
func (r *permissionRepository) GetByName(ctx context.Context, name string) (*permission.Permission, error) {
	var permModel models.PermissionModel
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&permModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.modelToDomain(&permModel)
}

// GetByResourceAndAction retrieves a permission by resource and action
func (r *permissionRepository) GetByResourceAndAction(ctx context.Context, resource, action string) (*permission.Permission, error) {
	var permModel models.PermissionModel
	if err := r.db.WithContext(ctx).Where("resource = ? AND action = ?", resource, action).First(&permModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.modelToDomain(&permModel)
}

// Update updates an existing permission
func (r *permissionRepository) Update(ctx context.Context, permEntity *permission.Permission) error {
	permModel := r.domainToModel(permEntity)
	if err := r.db.WithContext(ctx).Model(&permModel).Where("id = ?", permModel.ID).Updates(permModel).Error; err != nil {
		return err
	}
	return nil
}

// Delete soft deletes a permission
func (r *permissionRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.PermissionModel{}).Error; err != nil {
		return err
	}
	return nil
}

// List retrieves permissions with pagination
func (r *permissionRepository) List(ctx context.Context, params repositories.ListPermissionsParams) ([]*permission.Permission, *pagination.Pagination, error) {
	var permModels []models.PermissionModel
	var total int64

	query := r.db.WithContext(ctx).Model(&models.PermissionModel{})

	// Apply search filter
	if params.Search != "" {
		searchPattern := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(resource) LIKE ? OR LOWER(action) LIKE ? OR LOWER(description) LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Apply resource filter
	if params.Resource != "" {
		query = query.Where("resource = ?", params.Resource)
	}

	// Apply action filter
	if params.Action != "" {
		query = query.Where("action = ?", params.Action)
	}

	// Apply isActive filter
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
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
	sortDir := strings.ToUpper(params.SortDir)
	if sortDir != "ASC" && sortDir != "DESC" {
		sortDir = "DESC"
	}
	query = query.Order(sortBy + " " + sortDir)

	// Apply pagination
	query = query.Offset(paginationObj.Offset()).Limit(paginationObj.PageSize)

	// Execute query
	if err := query.Find(&permModels).Error; err != nil {
		return nil, nil, err
	}

	// Convert models to domain entities
	permissions := make([]*permission.Permission, 0, len(permModels))
	for _, permModel := range permModels {
		permEntity, err := r.modelToDomain(&permModel)
		if err != nil {
			return nil, nil, err
		}
		permissions = append(permissions, permEntity)
	}

	return permissions, paginationObj, nil
}

// GetActivePermissions retrieves all active permissions
func (r *permissionRepository) GetActivePermissions(ctx context.Context) ([]*permission.Permission, error) {
	var permModels []models.PermissionModel
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&permModels).Error; err != nil {
		return nil, err
	}

	permissions := make([]*permission.Permission, 0, len(permModels))
	for _, permModel := range permModels {
		permEntity, err := r.modelToDomain(&permModel)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permEntity)
	}

	return permissions, nil
}

// GetPermissionsByResource retrieves all permissions for a specific resource
func (r *permissionRepository) GetPermissionsByResource(ctx context.Context, resource string) ([]*permission.Permission, error) {
	var permModels []models.PermissionModel
	if err := r.db.WithContext(ctx).Where("resource = ? AND is_active = ?", resource, true).Find(&permModels).Error; err != nil {
		return nil, err
	}

	permissions := make([]*permission.Permission, 0, len(permModels))
	for _, permModel := range permModels {
		permEntity, err := r.modelToDomain(&permModel)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permEntity)
	}

	return permissions, nil
}

// GetPermissionsByAction retrieves all permissions for a specific action
func (r *permissionRepository) GetPermissionsByAction(ctx context.Context, action string) ([]*permission.Permission, error) {
	var permModels []models.PermissionModel
	if err := r.db.WithContext(ctx).Where("action = ? AND is_active = ?", action, true).Find(&permModels).Error; err != nil {
		return nil, err
	}

	permissions := make([]*permission.Permission, 0, len(permModels))
	for _, permModel := range permModels {
		permEntity, err := r.modelToDomain(&permModel)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permEntity)
	}

	return permissions, nil
}

// Exists checks if a permission exists by ID
func (r *permissionRepository) Exists(ctx context.Context, id string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.PermissionModel{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByName checks if a permission exists by name
func (r *permissionRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.PermissionModel{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByResourceAndAction checks if a permission exists by resource and action
func (r *permissionRepository) ExistsByResourceAndAction(ctx context.Context, resource, action string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.PermissionModel{}).Where("resource = ? AND action = ?", resource, action).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Count returns the total number of permissions
func (r *permissionRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.PermissionModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Activate activates a permission
func (r *permissionRepository) Activate(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Model(&models.PermissionModel{}).Where("id = ?", id).Update("is_active", true).Error; err != nil {
		return err
	}
	return nil
}

// Deactivate deactivates a permission
func (r *permissionRepository) Deactivate(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Model(&models.PermissionModel{}).Where("id = ?", id).Update("is_active", false).Error; err != nil {
		return err
	}
	return nil
}

// GetPermissionsByUserID retrieves all permissions for a user (through roles)
func (r *permissionRepository) GetPermissionsByUserID(ctx context.Context, userID string) ([]*permission.Permission, error) {
	var permModels []models.PermissionModel

	query := r.db.WithContext(ctx).
		Model(&models.PermissionModel{}).
		Joins("INNER JOIN role_permissions rp ON permissions.id = rp.permission_id").
		Joins("INNER JOIN user_roles ur ON rp.role_id = ur.role_id").
		Where("ur.user_id = ?", userID)

	if err := query.Find(&permModels).Error; err != nil {
		return nil, err
	}

	permissions := make([]*permission.Permission, 0, len(permModels))
	for _, permModel := range permModels {
		permEntity, err := r.modelToDomain(&permModel)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permEntity)
	}

	return permissions, nil
}

// GetActivePermissionsByUserID retrieves all active permissions for a user (through roles)
func (r *permissionRepository) GetActivePermissionsByUserID(ctx context.Context, userID string) ([]*permission.Permission, error) {
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

	permissions := make([]*permission.Permission, 0, len(permModels))
	for _, permModel := range permModels {
		permEntity, err := r.modelToDomain(&permModel)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permEntity)
	}

	return permissions, nil
}

// GetPermissionsByRoleID retrieves all permissions assigned to a role
func (r *permissionRepository) GetPermissionsByRoleID(ctx context.Context, roleID string) ([]*permission.Permission, error) {
	var permModels []models.PermissionModel

	query := r.db.WithContext(ctx).
		Model(&models.PermissionModel{}).
		Joins("INNER JOIN role_permissions rp ON permissions.id = rp.permission_id").
		Where("rp.role_id = ?", roleID)

	if err := query.Find(&permModels).Error; err != nil {
		return nil, err
	}

	permissions := make([]*permission.Permission, 0, len(permModels))
	for _, permModel := range permModels {
		permEntity, err := r.modelToDomain(&permModel)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permEntity)
	}

	return permissions, nil
}

// GetActivePermissionsByRoleID retrieves all active permissions assigned to a role
func (r *permissionRepository) GetActivePermissionsByRoleID(ctx context.Context, roleID string) ([]*permission.Permission, error) {
	var permModels []models.PermissionModel

	query := r.db.WithContext(ctx).
		Model(&models.PermissionModel{}).
		Joins("INNER JOIN role_permissions rp ON permissions.id = rp.permission_id").
		Where("rp.role_id = ? AND rp.is_active = ? AND permissions.is_active = ?", roleID, true, true)

	if err := query.Find(&permModels).Error; err != nil {
		return nil, err
	}

	permissions := make([]*permission.Permission, 0, len(permModels))
	for _, permModel := range permModels {
		permEntity, err := r.modelToDomain(&permModel)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permEntity)
	}

	return permissions, nil
}

// UserHasPermission checks if a user has a specific permission
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

// UserHasPermissionByName checks if a user has a specific permission by name
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

// domainToModel converts domain entity to GORM model
func (r *permissionRepository) domainToModel(permEntity *permission.Permission) *models.PermissionModel {
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

// modelToDomain converts GORM model to domain entity
func (r *permissionRepository) modelToDomain(permModel *models.PermissionModel) (*permission.Permission, error) {
	return permission.ReconstructPermission(
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
