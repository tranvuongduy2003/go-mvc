package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	userCommands "github.com/tranvuongduy2003/go-mvc/internal/application/commands/user"
	"github.com/tranvuongduy2003/go-mvc/internal/application/dto"
	userQueries "github.com/tranvuongduy2003/go-mvc/internal/application/queries/user"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/shared/valueobject"
	userDomain "github.com/tranvuongduy2003/go-mvc/internal/core/domain/user"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
	"github.com/tranvuongduy2003/go-mvc/pkg/jwt"
)

// UserApplicationService coordinates user-related operations
type UserApplicationService struct {
	// Command handlers
	createUserHandler     *userCommands.CreateUserCommandHandler
	updateUserHandler     *userCommands.UpdateUserCommandHandler
	changePasswordHandler *userCommands.ChangePasswordCommandHandler
	deleteUserHandler     *userCommands.DeleteUserCommandHandler
	activateUserHandler   *userCommands.ActivateUserCommandHandler
	deactivateUserHandler *userCommands.DeactivateUserCommandHandler
	changeRoleHandler     *userCommands.ChangeRoleCommandHandler

	// Query handlers
	getUserByIDHandler       *userQueries.GetUserByIDQueryHandler
	getUserByEmailHandler    *userQueries.GetUserByEmailQueryHandler
	getUserByUsernameHandler *userQueries.GetUserByUsernameQueryHandler
	listUsersHandler         *userQueries.ListUsersQueryHandler
	checkAvailabilityHandler *userQueries.CheckAvailabilityQueryHandler
	searchUsersHandler       *userQueries.SearchUsersQueryHandler

	// Services
	userService *userDomain.Service
	jwtService  *jwt.Service
	logger      *logger.Logger
}

// NewUserApplicationService creates a new user application service
func NewUserApplicationService(
	userService *userDomain.Service,
	repository userDomain.Repository,
	jwtService *jwt.Service,
	logger *logger.Logger,
) *UserApplicationService {
	return &UserApplicationService{
		// Initialize command handlers
		createUserHandler:     userCommands.NewCreateUserCommandHandler(userService, logger),
		updateUserHandler:     userCommands.NewUpdateUserCommandHandler(userService, repository, logger),
		changePasswordHandler: userCommands.NewChangePasswordCommandHandler(userService, logger),
		deleteUserHandler:     userCommands.NewDeleteUserCommandHandler(userService, logger),
		activateUserHandler:   userCommands.NewActivateUserCommandHandler(userService, logger),
		deactivateUserHandler: userCommands.NewDeactivateUserCommandHandler(userService, logger),
		changeRoleHandler:     userCommands.NewChangeRoleCommandHandler(userService, logger),

		// Initialize query handlers
		getUserByIDHandler:       userQueries.NewGetUserByIDQueryHandler(repository, logger),
		getUserByEmailHandler:    userQueries.NewGetUserByEmailQueryHandler(repository, logger),
		getUserByUsernameHandler: userQueries.NewGetUserByUsernameQueryHandler(repository, logger),
		listUsersHandler:         userQueries.NewListUsersQueryHandler(repository, logger),
		checkAvailabilityHandler: userQueries.NewCheckAvailabilityQueryHandler(repository, logger),
		searchUsersHandler:       userQueries.NewSearchUsersQueryHandler(repository, logger),

		// Services
		userService: userService,
		jwtService:  jwtService,
		logger:      logger,
	}
}

// CreateUser creates a new user
func (s *UserApplicationService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserDTO, error) {
	cmd := userCommands.CreateUserCommand{
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	return s.createUserHandler.Handle(ctx, cmd)
}

// UpdateUser updates an existing user
func (s *UserApplicationService) UpdateUser(ctx context.Context, userID uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserDTO, error) {
	cmd := userCommands.UpdateUserCommand{
		UserID:    userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	return s.updateUserHandler.Handle(ctx, cmd)
}

// ChangePassword changes user password
func (s *UserApplicationService) ChangePassword(ctx context.Context, userID uuid.UUID, req *dto.ChangePasswordRequest) error {
	cmd := userCommands.ChangePasswordCommand{
		UserID:          userID,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}

	return s.changePasswordHandler.Handle(ctx, cmd)
}

// DeleteUser deletes a user
func (s *UserApplicationService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	cmd := userCommands.DeleteUserCommand{
		UserID: userID,
	}

	return s.deleteUserHandler.Handle(ctx, cmd)
}

// ActivateUser activates a user
func (s *UserApplicationService) ActivateUser(ctx context.Context, userID uuid.UUID) error {
	cmd := userCommands.ActivateUserCommand{
		UserID: userID,
	}

	return s.activateUserHandler.Handle(ctx, cmd)
}

// DeactivateUser deactivates a user
func (s *UserApplicationService) DeactivateUser(ctx context.Context, userID uuid.UUID) error {
	cmd := userCommands.DeactivateUserCommand{
		UserID: userID,
	}

	return s.deactivateUserHandler.Handle(ctx, cmd)
}

// ChangeRole changes user role
func (s *UserApplicationService) ChangeRole(ctx context.Context, userID uuid.UUID, req *dto.ChangeRoleRequest) error {
	role := userDomain.Role(req.Role)
	cmd := userCommands.ChangeRoleCommand{
		UserID: userID,
		Role:   role,
	}

	return s.changeRoleHandler.Handle(ctx, cmd)
}

// GetUserByID gets a user by ID
func (s *UserApplicationService) GetUserByID(ctx context.Context, userID uuid.UUID) (*dto.UserDTO, error) {
	query := userQueries.GetUserByIDQuery{
		UserID: userID,
	}

	return s.getUserByIDHandler.Handle(ctx, query)
}

// GetUserByEmail gets a user by email
func (s *UserApplicationService) GetUserByEmail(ctx context.Context, email string) (*dto.UserDTO, error) {
	query := userQueries.GetUserByEmailQuery{
		Email: email,
	}

	return s.getUserByEmailHandler.Handle(ctx, query)
}

// GetUserByUsername gets a user by username
func (s *UserApplicationService) GetUserByUsername(ctx context.Context, username string) (*dto.UserDTO, error) {
	query := userQueries.GetUserByUsernameQuery{
		Username: username,
	}

	return s.getUserByUsernameHandler.Handle(ctx, query)
}

// ListUsers lists users with pagination and filters
func (s *UserApplicationService) ListUsers(ctx context.Context, req *dto.UserListRequest) (*dto.UserListResponse, error) {
	query := userQueries.ListUsersQuery{
		Page:      req.Page,
		Limit:     req.Limit,
		Sort:      req.Sort,
		Order:     req.Order,
		Search:    req.Search,
		Role:      req.Role,
		IsActive:  req.IsActive,
		IsDeleted: req.IsDeleted,
	}

	return s.listUsersHandler.Handle(ctx, query)
}

// CheckAvailability checks username/email availability
func (s *UserApplicationService) CheckAvailability(ctx context.Context, req *dto.CheckAvailabilityRequest) (*dto.CheckAvailabilityResponse, error) {
	query := userQueries.CheckAvailabilityQuery{
		Email:    req.Email,
		Username: req.Username,
	}

	return s.checkAvailabilityHandler.Handle(ctx, query)
}

// SearchUsers searches users
func (s *UserApplicationService) SearchUsers(ctx context.Context, searchQuery string, searchType string, page, limit int) (*dto.UserListResponse, error) {
	query := userQueries.SearchUsersQuery{
		Query: searchQuery,
		Type:  searchType,
		Page:  page,
		Limit: limit,
	}

	return s.searchUsersHandler.Handle(ctx, query)
}

// Login authenticates a user and returns tokens
func (s *UserApplicationService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	s.logger.Infof("Handling login request for email: %s", req.Email)

	// Create email value object
	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	// Validate credentials
	user, err := s.userService.ValidateCredentials(ctx, email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID().String(), string(user.Role()))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID().String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Get expiration time
	expiresAt := s.jwtService.GetAccessTokenExpirationTime()

	// Convert user to DTO
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

	response := &dto.LoginResponse{
		User:         userDTO,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}

	s.logger.Infof("User logged in successfully: %s", user.ID())
	return response, nil
}

// RefreshToken refreshes access token using refresh token
func (s *UserApplicationService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.LoginResponse, error) {
	s.logger.Infof("Handling refresh token request")

	// Validate refresh token
	claims, err := s.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get user ID from claims
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in token: %w", err)
	}

	// Get user
	userDTO, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is active
	if !userDTO.IsActive {
		return nil, fmt.Errorf("user is not active")
	}

	// Generate new tokens
	accessToken, err := s.jwtService.GenerateAccessToken(userDTO.ID.String(), userDTO.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(userDTO.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Get expiration time
	expiresAt := s.jwtService.GetAccessTokenExpirationTime()

	response := &dto.LoginResponse{
		User:         userDTO,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}

	s.logger.Infof("Token refreshed successfully for user: %s", userID)
	return response, nil
}

// ResetPassword resets user password (admin operation)
func (s *UserApplicationService) ResetPassword(ctx context.Context, userID uuid.UUID, req *dto.ResetPasswordRequest) error {
	return s.userService.ResetPassword(ctx, userID, req.NewPassword)
}

// GetUserStats returns user statistics
func (s *UserApplicationService) GetUserStats(ctx context.Context) (*dto.UserStatsResponse, error) {
	s.logger.Infof("Handling get user stats request")

	// This would require additional methods in the repository to get statistics
	// For now, return a basic implementation
	stats := &dto.UserStatsResponse{
		TotalUsers:     0,
		ActiveUsers:    0,
		InactiveUsers:  0,
		AdminUsers:     0,
		ModeratorUsers: 0,
		RegularUsers:   0,
		UsersThisMonth: 0,
		UsersToday:     0,
	}

	// TODO: Implement actual statistics gathering
	// This would involve adding methods to the repository to count users by various criteria

	return stats, nil
}
