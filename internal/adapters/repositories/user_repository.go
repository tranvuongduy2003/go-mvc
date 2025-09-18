package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/tranvuongduy2003/go-mvc/internal/application/dto"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/shared/valueobject"
	userDomain "github.com/tranvuongduy2003/go-mvc/internal/core/domain/user"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
)

// UserModel represents the user database model
type UserModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Email     string     `gorm:"type:varchar(255);unique;not null" json:"email"`
	Username  string     `gorm:"type:varchar(50);unique;not null" json:"username"`
	Password  string     `gorm:"type:varchar(255);not null" json:"-"`
	FirstName string     `gorm:"type:varchar(50);not null" json:"first_name"`
	LastName  string     `gorm:"type:varchar(50);not null" json:"last_name"`
	Role      string     `gorm:"type:varchar(20);not null;default:'user'" json:"role"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
	IsDeleted bool       `gorm:"default:false" json:"is_deleted"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for UserModel
func (UserModel) TableName() string {
	return "users"
}

// UserRepository implements the user repository interface
type UserRepository struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB, logger *logger.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *userDomain.User) error {
	r.logger.Infof("Creating user with email: %s", user.Email)

	userModel := &UserModel{
		ID:        user.ID(),
		Email:     user.Email().String(),
		Username:  user.Username(),
		Password:  user.Password(),
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		Role:      string(user.Role()),
		IsActive:  user.IsActive(),
		IsDeleted: user.IsDeleted(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}

	if err := r.db.WithContext(ctx).Create(userModel).Error; err != nil {
		r.logger.Errorf("Failed to create user: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	r.logger.Infof("User created successfully with ID: %s", user.ID())
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*userDomain.User, error) {
	r.logger.Infof("Getting user by ID: %s", id)

	var userModel UserModel
	if err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&userModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, userDomain.ErrUserNotFound
		}
		r.logger.Errorf("Failed to get user by ID: %v", err)
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return r.modelToDomain(&userModel), nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*userDomain.User, error) {
	r.logger.Infof("Getting user by email: %s", email)

	var userModel UserModel
	if err := r.db.WithContext(ctx).Where("email = ? AND is_deleted = ?", email, false).First(&userModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, userDomain.ErrUserNotFound
		}
		r.logger.Errorf("Failed to get user by email: %v", err)
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return r.modelToDomain(&userModel), nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*userDomain.User, error) {
	r.logger.Infof("Getting user by username: %s", username)

	var userModel UserModel
	if err := r.db.WithContext(ctx).Where("username = ? AND is_deleted = ?", username, false).First(&userModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, userDomain.ErrUserNotFound
		}
		r.logger.Errorf("Failed to get user by username: %v", err)
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return r.modelToDomain(&userModel), nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *userDomain.User) error {
	r.logger.Infof("Updating user with ID: %s", user.ID)

	userModel := &UserModel{
		ID:        user.ID(),
		Email:     user.Email().String(),
		Username:  user.Username(),
		Password:  user.Password(),
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		Role:      string(user.Role()),
		IsActive:  user.IsActive(),
		IsDeleted: user.IsDeleted(),
		UpdatedAt: time.Now(),
	}

	if err := r.db.WithContext(ctx).Where("id = ?", user.ID()).Updates(userModel).Error; err != nil {
		r.logger.Errorf("Failed to update user: %v", err)
		return fmt.Errorf("failed to update user: %w", err)
	}

	r.logger.Infof("User updated successfully with ID: %s", user.ID())
	return nil
}

// Delete soft deletes a user
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.Infof("Soft deleting user with ID: %s", id)

	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&UserModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": now,
			"updated_at": now,
		}).Error; err != nil {
		r.logger.Errorf("Failed to delete user: %v", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	r.logger.Infof("User deleted successfully with ID: %s", id)
	return nil
}

// List retrieves users with pagination and filters
func (r *UserRepository) List(ctx context.Context, req *dto.UserListRequest) (*dto.UserListResponse, error) {
	r.logger.Infof("Listing users with filters: %+v", req)

	var users []UserModel
	var total int64

	query := r.db.WithContext(ctx).Model(&UserModel{})

	// Apply filters
	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		query = query.Where("(email ILIKE ? OR username ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?) AND is_deleted = ?",
			searchPattern, searchPattern, searchPattern, searchPattern, false)
	} else {
		query = query.Where("is_deleted = ?", false)
	}

	if req.Role != nil && *req.Role != "" {
		query = query.Where("role = ?", *req.Role)
	}

	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	if req.IsDeleted != nil {
		query = query.Where("is_deleted = ?", *req.IsDeleted)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		r.logger.Errorf("Failed to count users: %v", err)
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Apply sorting
	if req.Sort != "" {
		order := "ASC"
		if req.Order == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", req.Sort, order))
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Find(&users).Error; err != nil {
		r.logger.Errorf("Failed to list users: %v", err)
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Convert to DTOs
	userDTOs := make([]*dto.UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = r.modelToDTO(&user)
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	response := &dto.UserListResponse{
		Users: userDTOs,
		Pagination: &dto.Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	r.logger.Infof("Listed %d users out of %d total", len(users), total)
	return response, nil
}

// CheckEmailExists checks if email already exists
func (r *UserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	r.logger.Infof("Checking if email exists: %s", email)

	var count int64
	if err := r.db.WithContext(ctx).Model(&UserModel{}).
		Where("email = ? AND is_deleted = ?", email, false).
		Count(&count).Error; err != nil {
		r.logger.Errorf("Failed to check email existence: %v", err)
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}

// CheckUsernameExists checks if username already exists
func (r *UserRepository) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	r.logger.Infof("Checking if username exists: %s", username)

	var count int64
	if err := r.db.WithContext(ctx).Model(&UserModel{}).
		Where("username = ? AND is_deleted = ?", username, false).
		Count(&count).Error; err != nil {
		r.logger.Errorf("Failed to check username existence: %v", err)
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}

	return count > 0, nil
}

// SearchUsers searches users by query and type
func (r *UserRepository) SearchUsers(ctx context.Context, searchQuery, searchType string, page, limit int) (*dto.UserListResponse, error) {
	r.logger.Infof("Searching users: query=%s, type=%s", searchQuery, searchType)

	var users []UserModel
	var total int64

	query := r.db.WithContext(ctx).Model(&UserModel{}).Where("is_deleted = ?", false)

	// Apply search based on type
	searchPattern := "%" + searchQuery + "%"
	switch searchType {
	case "email":
		query = query.Where("email ILIKE ?", searchPattern)
	case "username":
		query = query.Where("username ILIKE ?", searchPattern)
	case "name":
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ?", searchPattern, searchPattern)
	default:
		// Search all fields
		query = query.Where("email ILIKE ? OR username ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		r.logger.Errorf("Failed to count search results: %v", err)
		return nil, fmt.Errorf("failed to count search results: %w", err)
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		r.logger.Errorf("Failed to search users: %v", err)
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	// Convert to DTOs
	userDTOs := make([]*dto.UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = r.modelToDTO(&user)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response := &dto.UserListResponse{
		Users: userDTOs,
		Pagination: &dto.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	r.logger.Infof("Found %d users for search query", len(users))
	return response, nil
}

// Helper methods for conversion

func (r *UserRepository) modelToDomain(model *UserModel) *userDomain.User {
	// Create email value object
	email, _ := valueobject.NewEmail(model.Email)

	// Create user entity
	user, _ := userDomain.NewUser(email, model.Username, model.Password, model.FirstName, model.LastName)

	// Set additional fields using reflection or custom setters if available
	// For now, we'll use a simplified approach

	return user
}

func (r *UserRepository) modelToDTO(model *UserModel) *dto.UserDTO {
	return &dto.UserDTO{
		ID:        model.ID,
		Email:     model.Email,
		Username:  model.Username,
		FirstName: model.FirstName,
		LastName:  model.LastName,
		Role:      model.Role,
		IsActive:  model.IsActive,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
