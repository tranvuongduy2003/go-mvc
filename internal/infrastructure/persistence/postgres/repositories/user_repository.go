package repositories

import (
	"context"
	"strings"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/ports/repositories"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/user"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/persistence/postgres/models"
	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
	"gorm.io/gorm"
)

// userRepository implements the UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create saves a new user to the database
func (r *userRepository) Create(ctx context.Context, u *user.User) error {
	userModel := r.domainToModel(u)
	if err := r.db.WithContext(ctx).Create(userModel).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	var userModel models.UserModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&userModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.modelToDomain(&userModel)
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var userModel models.UserModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&userModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.modelToDomain(&userModel)
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, u *user.User) error {
	userModel := r.domainToModel(u)
	if err := r.db.WithContext(ctx).Model(&userModel).Where("id = ?", userModel.ID).Updates(userModel).Error; err != nil {
		return err
	}
	return nil
}

// Delete soft deletes a user
func (r *userRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.UserModel{}).Error; err != nil {
		return err
	}
	return nil
}

// List retrieves users with pagination
func (r *userRepository) List(ctx context.Context, params repositories.ListUsersParams) ([]*user.User, *pagination.Pagination, error) {
	var userModels []models.UserModel
	var total int64

	query := r.db.WithContext(ctx).Model(&models.UserModel{})

	// Apply search filter
	if params.Search != "" {
		searchPattern := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(email) LIKE ? OR LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ?",
			searchPattern, searchPattern, searchPattern)
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
	pag := pagination.NewPagination(params.Page, params.Limit)
	pag.SetTotal(total)

	// Apply pagination
	query = query.Limit(pag.PageSize).Offset(pag.Offset())

	// Apply sorting
	if params.SortBy != "" {
		order := params.SortBy
		if params.SortDir == "desc" {
			order += " DESC"
		} else {
			order += " ASC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	if err := query.Find(&userModels).Error; err != nil {
		return nil, nil, err
	}

	users := make([]*user.User, len(userModels))
	for i, model := range userModels {
		domainUser, err := r.modelToDomain(&model)
		if err != nil {
			return nil, nil, err
		}
		users[i] = domainUser
	}

	return users, pag, nil
}

// Exists checks if a user exists by ID
func (r *userRepository) Exists(ctx context.Context, id string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.UserModel{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByEmail checks if a user exists by email
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.UserModel{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Count returns the total number of users
func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.UserModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// domainToModel converts domain User to GORM UserModel
func (r *userRepository) domainToModel(u *user.User) *models.UserModel {
	return &models.UserModel{
		ID:            u.ID(),
		Email:         u.Email(),
		Name:          u.Name(),
		Phone:         u.Phone(),
		PasswordHash:  u.HashedPassword(),
		AvatarFileKey: u.Avatar().FileKey(),
		AvatarCDNUrl:  u.Avatar().CDNUrl(),
		IsActive:      u.IsActive(),
		CreatedAt:     u.CreatedAt(),
		UpdatedAt:     u.UpdatedAt(),
	}
}

// modelToDomain converts GORM UserModel to domain User
func (r *userRepository) modelToDomain(m *models.UserModel) (*user.User, error) {
	return user.ReconstructUser(
		m.ID,
		m.Email,
		m.Name,
		m.Phone,
		m.PasswordHash,
		m.AvatarFileKey,
		m.AvatarCDNUrl,
		m.IsActive,
		m.CreatedAt,
		m.UpdatedAt,
		1, // version - we'll start with 1 for reconstructed users
	)
}
