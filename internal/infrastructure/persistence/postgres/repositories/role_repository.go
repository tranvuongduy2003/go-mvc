package repositories

import (
	"context"
	"strings"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/ports/repositories"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/role"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/persistence/postgres/models"
	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
	"gorm.io/gorm"
)

// roleRepository implements the RoleRepository interface
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new RoleRepository instance
func NewRoleRepository(db *gorm.DB) repositories.RoleRepository {
	return &roleRepository{
		db: db,
	}
}

// Create saves a new role to the database
func (r *roleRepository) Create(ctx context.Context, roleEntity *role.Role) error {
	roleModel := r.domainToModel(roleEntity)
	if err := r.db.WithContext(ctx).Create(roleModel).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a role by ID
func (r *roleRepository) GetByID(ctx context.Context, id string) (*role.Role, error) {
	var roleModel models.RoleModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&roleModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.modelToDomain(&roleModel)
}

// GetByName retrieves a role by name
func (r *roleRepository) GetByName(ctx context.Context, name string) (*role.Role, error) {
	var roleModel models.RoleModel
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&roleModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.modelToDomain(&roleModel)
}

// Update updates an existing role
func (r *roleRepository) Update(ctx context.Context, roleEntity *role.Role) error {
	roleModel := r.domainToModel(roleEntity)
	if err := r.db.WithContext(ctx).Model(&roleModel).Where("id = ?", roleModel.ID).Updates(roleModel).Error; err != nil {
		return err
	}
	return nil
}

// Delete soft deletes a role
func (r *roleRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.RoleModel{}).Error; err != nil {
		return err
	}
	return nil
}

// List retrieves roles with pagination
func (r *roleRepository) List(ctx context.Context, params repositories.ListRolesParams) ([]*role.Role, *pagination.Pagination, error) {
	var roleModels []models.RoleModel
	var total int64

	query := r.db.WithContext(ctx).Model(&models.RoleModel{})

	// Apply search filter
	if params.Search != "" {
		searchPattern := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?",
			searchPattern, searchPattern)
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
	if err := query.Find(&roleModels).Error; err != nil {
		return nil, nil, err
	}

	// Convert models to domain entities
	roles := make([]*role.Role, 0, len(roleModels))
	for _, roleModel := range roleModels {
		roleEntity, err := r.modelToDomain(&roleModel)
		if err != nil {
			return nil, nil, err
		}
		roles = append(roles, roleEntity)
	}

	return roles, paginationObj, nil
}

// GetActiveRoles retrieves all active roles
func (r *roleRepository) GetActiveRoles(ctx context.Context) ([]*role.Role, error) {
	var roleModels []models.RoleModel
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&roleModels).Error; err != nil {
		return nil, err
	}

	roles := make([]*role.Role, 0, len(roleModels))
	for _, roleModel := range roleModels {
		roleEntity, err := r.modelToDomain(&roleModel)
		if err != nil {
			return nil, err
		}
		roles = append(roles, roleEntity)
	}

	return roles, nil
}

// Exists checks if a role exists by ID
func (r *roleRepository) Exists(ctx context.Context, id string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.RoleModel{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByName checks if a role exists by name
func (r *roleRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.RoleModel{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Count returns the total number of roles
func (r *roleRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.RoleModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Activate activates a role
func (r *roleRepository) Activate(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Model(&models.RoleModel{}).Where("id = ?", id).Update("is_active", true).Error; err != nil {
		return err
	}
	return nil
}

// Deactivate deactivates a role
func (r *roleRepository) Deactivate(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Model(&models.RoleModel{}).Where("id = ?", id).Update("is_active", false).Error; err != nil {
		return err
	}
	return nil
}

// GetRolesByUserID retrieves all roles assigned to a user
func (r *roleRepository) GetRolesByUserID(ctx context.Context, userID string) ([]*role.Role, error) {
	var roleModels []models.RoleModel

	query := r.db.WithContext(ctx).
		Model(&models.RoleModel{}).
		Joins("INNER JOIN user_roles ur ON roles.id = ur.role_id").
		Where("ur.user_id = ?", userID)

	if err := query.Find(&roleModels).Error; err != nil {
		return nil, err
	}

	roles := make([]*role.Role, 0, len(roleModels))
	for _, roleModel := range roleModels {
		roleEntity, err := r.modelToDomain(&roleModel)
		if err != nil {
			return nil, err
		}
		roles = append(roles, roleEntity)
	}

	return roles, nil
}

// GetActiveRolesByUserID retrieves all active roles assigned to a user
func (r *roleRepository) GetActiveRolesByUserID(ctx context.Context, userID string) ([]*role.Role, error) {
	var roleModels []models.RoleModel

	query := r.db.WithContext(ctx).
		Model(&models.RoleModel{}).
		Joins("INNER JOIN user_roles ur ON roles.id = ur.role_id").
		Where("ur.user_id = ? AND ur.is_active = ? AND roles.is_active = ?", userID, true, true).
		Where("ur.expires_at IS NULL OR ur.expires_at > CURRENT_TIMESTAMP")

	if err := query.Find(&roleModels).Error; err != nil {
		return nil, err
	}

	roles := make([]*role.Role, 0, len(roleModels))
	for _, roleModel := range roleModels {
		roleEntity, err := r.modelToDomain(&roleModel)
		if err != nil {
			return nil, err
		}
		roles = append(roles, roleEntity)
	}

	return roles, nil
}

// domainToModel converts domain entity to GORM model
func (r *roleRepository) domainToModel(roleEntity *role.Role) *models.RoleModel {
	return &models.RoleModel{
		ID:          roleEntity.ID().String(),
		Name:        roleEntity.Name().String(),
		Description: roleEntity.Description(),
		IsActive:    roleEntity.IsActive(),
		CreatedAt:   roleEntity.CreatedAt(),
		UpdatedAt:   roleEntity.UpdatedAt(),
		Version:     roleEntity.Version(),
	}
}

// modelToDomain converts GORM model to domain entity
func (r *roleRepository) modelToDomain(roleModel *models.RoleModel) (*role.Role, error) {
	return role.ReconstructRole(
		roleModel.ID,
		roleModel.Name,
		roleModel.Description,
		roleModel.IsActive,
		roleModel.CreatedAt,
		roleModel.UpdatedAt,
		roleModel.Version,
	)
}
