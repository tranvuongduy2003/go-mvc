package user

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/application/dto"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/shared/valueobject"
	userDomain "github.com/tranvuongduy2003/go-mvc/internal/core/domain/user"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
)

// GetUserByIDQuery represents a query to get user by ID
type GetUserByIDQuery struct {
	UserID uuid.UUID
}

// GetUserByIDQueryHandler handles get user by ID queries
type GetUserByIDQueryHandler struct {
	repository userDomain.Repository
	logger     *logger.Logger
}

// NewGetUserByIDQueryHandler creates a new get user by ID query handler
func NewGetUserByIDQueryHandler(repository userDomain.Repository, logger *logger.Logger) *GetUserByIDQueryHandler {
	return &GetUserByIDQueryHandler{
		repository: repository,
		logger:     logger,
	}
}

// Handle handles the get user by ID query
func (h *GetUserByIDQueryHandler) Handle(ctx context.Context, query GetUserByIDQuery) (*dto.UserDTO, error) {
	h.logger.Infof("Handling get user by ID query: %s", query.UserID)

	user, err := h.repository.GetByID(ctx, query.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get user profile
	profile, err := h.repository.GetProfileByUserID(ctx, user.ID())
	if err != nil {
		h.logger.Warnf("Failed to get user profile: %v", err)
		// Continue without profile
	} else {
		user.SetProfile(profile)
	}

	return h.toUserDTO(user), nil
}

// GetUserByEmailQuery represents a query to get user by email
type GetUserByEmailQuery struct {
	Email string
}

// GetUserByEmailQueryHandler handles get user by email queries
type GetUserByEmailQueryHandler struct {
	repository userDomain.Repository
	logger     *logger.Logger
}

// NewGetUserByEmailQueryHandler creates a new get user by email query handler
func NewGetUserByEmailQueryHandler(repository userDomain.Repository, logger *logger.Logger) *GetUserByEmailQueryHandler {
	return &GetUserByEmailQueryHandler{
		repository: repository,
		logger:     logger,
	}
}

// Handle handles the get user by email query
func (h *GetUserByEmailQueryHandler) Handle(ctx context.Context, query GetUserByEmailQuery) (*dto.UserDTO, error) {
	h.logger.Infof("Handling get user by email query: %s", query.Email)

	email, err := valueobject.NewEmail(query.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	user, err := h.repository.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get user profile
	profile, err := h.repository.GetProfileByUserID(ctx, user.ID())
	if err != nil {
		h.logger.Warnf("Failed to get user profile: %v", err)
		// Continue without profile
	} else {
		user.SetProfile(profile)
	}

	return h.toUserDTO(user), nil
}

// GetUserByUsernameQuery represents a query to get user by username
type GetUserByUsernameQuery struct {
	Username string
}

// GetUserByUsernameQueryHandler handles get user by username queries
type GetUserByUsernameQueryHandler struct {
	repository userDomain.Repository
	logger     *logger.Logger
}

// NewGetUserByUsernameQueryHandler creates a new get user by username query handler
func NewGetUserByUsernameQueryHandler(repository userDomain.Repository, logger *logger.Logger) *GetUserByUsernameQueryHandler {
	return &GetUserByUsernameQueryHandler{
		repository: repository,
		logger:     logger,
	}
}

// Handle handles the get user by username query
func (h *GetUserByUsernameQueryHandler) Handle(ctx context.Context, query GetUserByUsernameQuery) (*dto.UserDTO, error) {
	h.logger.Infof("Handling get user by username query: %s", query.Username)

	user, err := h.repository.GetByUsername(ctx, query.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get user profile
	profile, err := h.repository.GetProfileByUserID(ctx, user.ID())
	if err != nil {
		h.logger.Warnf("Failed to get user profile: %v", err)
		// Continue without profile
	} else {
		user.SetProfile(profile)
	}

	return h.toUserDTO(user), nil
}

// ListUsersQuery represents a query to list users
type ListUsersQuery struct {
	Page      int
	Limit     int
	Sort      string
	Order     string
	Search    string
	Role      *string
	IsActive  *bool
	IsDeleted *bool
}

// ListUsersQueryHandler handles list users queries
type ListUsersQueryHandler struct {
	repository userDomain.Repository
	logger     *logger.Logger
}

// NewListUsersQueryHandler creates a new list users query handler
func NewListUsersQueryHandler(repository userDomain.Repository, logger *logger.Logger) *ListUsersQueryHandler {
	return &ListUsersQueryHandler{
		repository: repository,
		logger:     logger,
	}
}

// Handle handles the list users query
func (h *ListUsersQueryHandler) Handle(ctx context.Context, query ListUsersQuery) (*dto.UserListResponse, error) {
	h.logger.Infof("Handling list users query - page: %d, limit: %d", query.Page, query.Limit)

	// Build filters
	filters := userDomain.ListFilters{
		Limit:     query.Limit,
		Offset:    (query.Page - 1) * query.Limit,
		SortBy:    query.Sort,
		SortOrder: query.Order,
		Search:    query.Search,
		IsActive:  query.IsActive,
		IsDeleted: query.IsDeleted,
	}

	// Set role filter if provided
	if query.Role != nil {
		role := userDomain.Role(*query.Role)
		filters.Role = &role
	}

	// Get users
	users, err := h.repository.List(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Get total count
	totalCount, err := h.repository.Count(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Convert to DTOs
	userDTOs := make([]*dto.UserDTO, len(users))
	for i, user := range users {
		// Get user profile
		profile, err := h.repository.GetProfileByUserID(ctx, user.ID())
		if err != nil {
			h.logger.Warnf("Failed to get user profile for user %s: %v", user.ID(), err)
			// Continue without profile
		} else {
			user.SetProfile(profile)
		}

		userDTOs[i] = h.toUserDTO(user)
	}

	// Create pagination
	totalPages := int(math.Ceil(float64(totalCount) / float64(query.Limit)))
	pagination := &dto.Pagination{
		Page:       query.Page,
		Limit:      query.Limit,
		Total:      totalCount,
		TotalPages: totalPages,
		HasNext:    query.Page < totalPages,
		HasPrev:    query.Page > 1,
	}

	return &dto.UserListResponse{
		Users:      userDTOs,
		Pagination: pagination,
	}, nil
}

// CheckAvailabilityQuery represents a query to check username/email availability
type CheckAvailabilityQuery struct {
	Email    *string
	Username *string
}

// CheckAvailabilityQueryHandler handles availability check queries
type CheckAvailabilityQueryHandler struct {
	repository userDomain.Repository
	logger     *logger.Logger
}

// NewCheckAvailabilityQueryHandler creates a new check availability query handler
func NewCheckAvailabilityQueryHandler(repository userDomain.Repository, logger *logger.Logger) *CheckAvailabilityQueryHandler {
	return &CheckAvailabilityQueryHandler{
		repository: repository,
		logger:     logger,
	}
}

// Handle handles the check availability query
func (h *CheckAvailabilityQueryHandler) Handle(ctx context.Context, query CheckAvailabilityQuery) (*dto.CheckAvailabilityResponse, error) {
	h.logger.Infof("Handling check availability query")

	response := &dto.CheckAvailabilityResponse{}

	// Check email availability
	if query.Email != nil {
		email, err := valueobject.NewEmail(*query.Email)
		if err != nil {
			return nil, fmt.Errorf("invalid email: %w", err)
		}

		exists, err := h.repository.ExistsByEmail(ctx, email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email availability: %w", err)
		}
		available := !exists
		response.EmailAvailable = &available
	}

	// Check username availability
	if query.Username != nil {
		exists, err := h.repository.ExistsByUsername(ctx, *query.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to check username availability: %w", err)
		}
		available := !exists
		response.UsernameAvailable = &available
	}

	return response, nil
}

// SearchUsersQuery represents a query to search users
type SearchUsersQuery struct {
	Query string
	Page  int
	Limit int
	Type  string // "name" or "email"
}

// SearchUsersQueryHandler handles user search queries
type SearchUsersQueryHandler struct {
	repository userDomain.Repository
	logger     *logger.Logger
}

// NewSearchUsersQueryHandler creates a new search users query handler
func NewSearchUsersQueryHandler(repository userDomain.Repository, logger *logger.Logger) *SearchUsersQueryHandler {
	return &SearchUsersQueryHandler{
		repository: repository,
		logger:     logger,
	}
}

// Handle handles the search users query
func (h *SearchUsersQueryHandler) Handle(ctx context.Context, query SearchUsersQuery) (*dto.UserListResponse, error) {
	h.logger.Infof("Handling search users query: %s", query.Query)

	offset := (query.Page - 1) * query.Limit
	var users []*userDomain.User
	var err error

	switch query.Type {
	case "email":
		users, err = h.repository.SearchByEmail(ctx, query.Query, query.Limit, offset)
	case "name":
		users, err = h.repository.SearchByName(ctx, query.Query, query.Limit, offset)
	default:
		// Default to name search
		users, err = h.repository.SearchByName(ctx, query.Query, query.Limit, offset)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	// Convert to DTOs
	userDTOs := make([]*dto.UserDTO, len(users))
	for i, user := range users {
		// Get user profile
		profile, err := h.repository.GetProfileByUserID(ctx, user.ID())
		if err != nil {
			h.logger.Warnf("Failed to get user profile for user %s: %v", user.ID(), err)
			// Continue without profile
		} else {
			user.SetProfile(profile)
		}

		userDTOs[i] = h.toUserDTO(user)
	}

	// For search, we don't calculate exact total count for performance
	// We just check if there are more results
	hasNext := len(users) == query.Limit
	pagination := &dto.Pagination{
		Page:    query.Page,
		Limit:   query.Limit,
		HasNext: hasNext,
		HasPrev: query.Page > 1,
	}

	return &dto.UserListResponse{
		Users:      userDTOs,
		Pagination: pagination,
	}, nil
}

// Helper method to convert domain user to DTO
func (h *GetUserByIDQueryHandler) toUserDTO(user *userDomain.User) *dto.UserDTO {
	userDTO := &dto.UserDTO{
		ID:        user.ID(),
		Email:     user.Email().Value(),
		Username:  user.Username(),
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		FullName:  user.FullName(),
		Role:      string(user.Role()),
		IsActive:  user.IsActive(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}

	// Add profile if exists
	if profile := user.Profile(); profile != nil {
		userDTO.Profile = &dto.ProfileDTO{
			ID:          profile.ID,
			UserID:      profile.UserID,
			Avatar:      profile.Avatar,
			Bio:         profile.Bio,
			DateOfBirth: profile.DateOfBirth,
			Phone:       profile.Phone,
			Address:     profile.Address,
			City:        profile.City,
			Country:     profile.Country,
			Website:     profile.Website,
			SocialLinks: profile.SocialLinks,
			CreatedAt:   profile.CreatedAt,
			UpdatedAt:   profile.UpdatedAt,
		}
	}

	return userDTO
}

// Helper methods for other handlers (similar implementation)
func (h *GetUserByEmailQueryHandler) toUserDTO(user *userDomain.User) *dto.UserDTO {
	return h.convertToDTO(user)
}

func (h *GetUserByUsernameQueryHandler) toUserDTO(user *userDomain.User) *dto.UserDTO {
	return h.convertToDTO(user)
}

func (h *ListUsersQueryHandler) toUserDTO(user *userDomain.User) *dto.UserDTO {
	return h.convertToDTO(user)
}

func (h *SearchUsersQueryHandler) toUserDTO(user *userDomain.User) *dto.UserDTO {
	return h.convertToDTO(user)
}

// Shared conversion method
func (h *GetUserByEmailQueryHandler) convertToDTO(user *userDomain.User) *dto.UserDTO {
	userDTO := &dto.UserDTO{
		ID:        user.ID(),
		Email:     user.Email().Value(),
		Username:  user.Username(),
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		FullName:  user.FullName(),
		Role:      string(user.Role()),
		IsActive:  user.IsActive(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}

	if profile := user.Profile(); profile != nil {
		userDTO.Profile = &dto.ProfileDTO{
			ID:          profile.ID,
			UserID:      profile.UserID,
			Avatar:      profile.Avatar,
			Bio:         profile.Bio,
			DateOfBirth: profile.DateOfBirth,
			Phone:       profile.Phone,
			Address:     profile.Address,
			City:        profile.City,
			Country:     profile.Country,
			Website:     profile.Website,
			SocialLinks: profile.SocialLinks,
			CreatedAt:   profile.CreatedAt,
			UpdatedAt:   profile.UpdatedAt,
		}
	}

	return userDTO
}

// Similar conversion methods for other handlers
func (h *GetUserByUsernameQueryHandler) convertToDTO(user *userDomain.User) *dto.UserDTO {
	userDTO := &dto.UserDTO{
		ID:        user.ID(),
		Email:     user.Email().Value(),
		Username:  user.Username(),
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		FullName:  user.FullName(),
		Role:      string(user.Role()),
		IsActive:  user.IsActive(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}

	if profile := user.Profile(); profile != nil {
		userDTO.Profile = &dto.ProfileDTO{
			ID:          profile.ID,
			UserID:      profile.UserID,
			Avatar:      profile.Avatar,
			Bio:         profile.Bio,
			DateOfBirth: profile.DateOfBirth,
			Phone:       profile.Phone,
			Address:     profile.Address,
			City:        profile.City,
			Country:     profile.Country,
			Website:     profile.Website,
			SocialLinks: profile.SocialLinks,
			CreatedAt:   profile.CreatedAt,
			UpdatedAt:   profile.UpdatedAt,
		}
	}

	return userDTO
}

func (h *ListUsersQueryHandler) convertToDTO(user *userDomain.User) *dto.UserDTO {
	userDTO := &dto.UserDTO{
		ID:        user.ID(),
		Email:     user.Email().Value(),
		Username:  user.Username(),
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		FullName:  user.FullName(),
		Role:      string(user.Role()),
		IsActive:  user.IsActive(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}

	if profile := user.Profile(); profile != nil {
		userDTO.Profile = &dto.ProfileDTO{
			ID:          profile.ID,
			UserID:      profile.UserID,
			Avatar:      profile.Avatar,
			Bio:         profile.Bio,
			DateOfBirth: profile.DateOfBirth,
			Phone:       profile.Phone,
			Address:     profile.Address,
			City:        profile.City,
			Country:     profile.Country,
			Website:     profile.Website,
			SocialLinks: profile.SocialLinks,
			CreatedAt:   profile.CreatedAt,
			UpdatedAt:   profile.UpdatedAt,
		}
	}

	return userDTO
}

func (h *SearchUsersQueryHandler) convertToDTO(user *userDomain.User) *dto.UserDTO {
	userDTO := &dto.UserDTO{
		ID:        user.ID(),
		Email:     user.Email().Value(),
		Username:  user.Username(),
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		FullName:  user.FullName(),
		Role:      string(user.Role()),
		IsActive:  user.IsActive(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}

	if profile := user.Profile(); profile != nil {
		userDTO.Profile = &dto.ProfileDTO{
			ID:          profile.ID,
			UserID:      profile.UserID,
			Avatar:      profile.Avatar,
			Bio:         profile.Bio,
			DateOfBirth: profile.DateOfBirth,
			Phone:       profile.Phone,
			Address:     profile.Address,
			City:        profile.City,
			Country:     profile.Country,
			Website:     profile.Website,
			SocialLinks: profile.SocialLinks,
			CreatedAt:   profile.CreatedAt,
			UpdatedAt:   profile.UpdatedAt,
		}
	}

	return userDTO
}
