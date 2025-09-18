package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/application/dto"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/shared/valueobject"
	userDomain "github.com/tranvuongduy2003/go-mvc/internal/core/domain/user"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
)

// CreateUserCommand represents a command to create a user
type CreateUserCommand struct {
	Email     string
	Username  string
	Password  string
	FirstName string
	LastName  string
}

// CreateUserCommandHandler handles user creation commands
type CreateUserCommandHandler struct {
	userService *userDomain.Service
	logger      *logger.Logger
}

// NewCreateUserCommandHandler creates a new create user command handler
func NewCreateUserCommandHandler(userService *userDomain.Service, logger *logger.Logger) *CreateUserCommandHandler {
	return &CreateUserCommandHandler{
		userService: userService,
		logger:      logger,
	}
}

// Handle handles the create user command
func (h *CreateUserCommandHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*dto.UserDTO, error) {
	h.logger.Infof("Handling create user command for email: %s", cmd.Email)

	// Create email value object
	email, err := valueobject.NewEmail(cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	// Create user through domain service
	user, err := h.userService.CreateUser(ctx, email, cmd.Username, cmd.Password, cmd.FirstName, cmd.LastName)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Convert to DTO
	userDTO := h.toUserDTO(user)

	h.logger.Infof("User created successfully: %s", user.ID())
	return userDTO, nil
}

// UpdateUserCommand represents a command to update a user
type UpdateUserCommand struct {
	UserID    uuid.UUID
	FirstName *string
	LastName  *string
	Email     *string
}

// UpdateUserCommandHandler handles user update commands
type UpdateUserCommandHandler struct {
	userService *userDomain.Service
	repository  userDomain.Repository
	logger      *logger.Logger
}

// NewUpdateUserCommandHandler creates a new update user command handler
func NewUpdateUserCommandHandler(userService *userDomain.Service, repository userDomain.Repository, logger *logger.Logger) *UpdateUserCommandHandler {
	return &UpdateUserCommandHandler{
		userService: userService,
		repository:  repository,
		logger:      logger,
	}
}

// Handle handles the update user command
func (h *UpdateUserCommandHandler) Handle(ctx context.Context, cmd UpdateUserCommand) (*dto.UserDTO, error) {
	h.logger.Infof("Handling update user command for user: %s", cmd.UserID)

	// Get user
	user, err := h.repository.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update profile if provided
	if cmd.FirstName != nil || cmd.LastName != nil {
		firstName := user.FirstName()
		lastName := user.LastName()

		if cmd.FirstName != nil {
			firstName = *cmd.FirstName
		}
		if cmd.LastName != nil {
			lastName = *cmd.LastName
		}

		if err := h.userService.UpdateUserProfile(ctx, cmd.UserID, firstName, lastName); err != nil {
			return nil, fmt.Errorf("failed to update user profile: %w", err)
		}
	}

	// Update email if provided
	if cmd.Email != nil {
		email, err := valueobject.NewEmail(*cmd.Email)
		if err != nil {
			return nil, fmt.Errorf("invalid email: %w", err)
		}

		if err := user.UpdateEmail(email); err != nil {
			return nil, fmt.Errorf("failed to update email: %w", err)
		}

		if err := h.repository.Update(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to save user: %w", err)
		}
	}

	// Get updated user
	updatedUser, err := h.repository.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated user: %w", err)
	}

	// Convert to DTO
	userDTO := h.toUserDTO(updatedUser)

	h.logger.Infof("User updated successfully: %s", cmd.UserID)
	return userDTO, nil
}

// ChangePasswordCommand represents a command to change user password
type ChangePasswordCommand struct {
	UserID          uuid.UUID
	CurrentPassword string
	NewPassword     string
}

// ChangePasswordCommandHandler handles password change commands
type ChangePasswordCommandHandler struct {
	userService *userDomain.Service
	logger      *logger.Logger
}

// NewChangePasswordCommandHandler creates a new change password command handler
func NewChangePasswordCommandHandler(userService *userDomain.Service, logger *logger.Logger) *ChangePasswordCommandHandler {
	return &ChangePasswordCommandHandler{
		userService: userService,
		logger:      logger,
	}
}

// Handle handles the change password command
func (h *ChangePasswordCommandHandler) Handle(ctx context.Context, cmd ChangePasswordCommand) error {
	h.logger.Infof("Handling change password command for user: %s", cmd.UserID)

	if err := h.userService.ChangePassword(ctx, cmd.UserID, cmd.CurrentPassword, cmd.NewPassword); err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}

	h.logger.Infof("Password changed successfully for user: %s", cmd.UserID)
	return nil
}

// DeleteUserCommand represents a command to delete a user
type DeleteUserCommand struct {
	UserID uuid.UUID
}

// DeleteUserCommandHandler handles user deletion commands
type DeleteUserCommandHandler struct {
	userService *userDomain.Service
	logger      *logger.Logger
}

// NewDeleteUserCommandHandler creates a new delete user command handler
func NewDeleteUserCommandHandler(userService *userDomain.Service, logger *logger.Logger) *DeleteUserCommandHandler {
	return &DeleteUserCommandHandler{
		userService: userService,
		logger:      logger,
	}
}

// Handle handles the delete user command
func (h *DeleteUserCommandHandler) Handle(ctx context.Context, cmd DeleteUserCommand) error {
	h.logger.Infof("Handling delete user command for user: %s", cmd.UserID)

	if err := h.userService.DeleteUser(ctx, cmd.UserID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	h.logger.Infof("User deleted successfully: %s", cmd.UserID)
	return nil
}

// ActivateUserCommand represents a command to activate a user
type ActivateUserCommand struct {
	UserID uuid.UUID
}

// ActivateUserCommandHandler handles user activation commands
type ActivateUserCommandHandler struct {
	userService *userDomain.Service
	logger      *logger.Logger
}

// NewActivateUserCommandHandler creates a new activate user command handler
func NewActivateUserCommandHandler(userService *userDomain.Service, logger *logger.Logger) *ActivateUserCommandHandler {
	return &ActivateUserCommandHandler{
		userService: userService,
		logger:      logger,
	}
}

// Handle handles the activate user command
func (h *ActivateUserCommandHandler) Handle(ctx context.Context, cmd ActivateUserCommand) error {
	h.logger.Infof("Handling activate user command for user: %s", cmd.UserID)

	if err := h.userService.ActivateUser(ctx, cmd.UserID); err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}

	h.logger.Infof("User activated successfully: %s", cmd.UserID)
	return nil
}

// DeactivateUserCommand represents a command to deactivate a user
type DeactivateUserCommand struct {
	UserID uuid.UUID
}

// DeactivateUserCommandHandler handles user deactivation commands
type DeactivateUserCommandHandler struct {
	userService *userDomain.Service
	logger      *logger.Logger
}

// NewDeactivateUserCommandHandler creates a new deactivate user command handler
func NewDeactivateUserCommandHandler(userService *userDomain.Service, logger *logger.Logger) *DeactivateUserCommandHandler {
	return &DeactivateUserCommandHandler{
		userService: userService,
		logger:      logger,
	}
}

// Handle handles the deactivate user command
func (h *DeactivateUserCommandHandler) Handle(ctx context.Context, cmd DeactivateUserCommand) error {
	h.logger.Infof("Handling deactivate user command for user: %s", cmd.UserID)

	if err := h.userService.DeactivateUser(ctx, cmd.UserID); err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	h.logger.Infof("User deactivated successfully: %s", cmd.UserID)
	return nil
}

// ChangeRoleCommand represents a command to change user role
type ChangeRoleCommand struct {
	UserID uuid.UUID
	Role   userDomain.Role
}

// ChangeRoleCommandHandler handles user role change commands
type ChangeRoleCommandHandler struct {
	userService *userDomain.Service
	logger      *logger.Logger
}

// NewChangeRoleCommandHandler creates a new change role command handler
func NewChangeRoleCommandHandler(userService *userDomain.Service, logger *logger.Logger) *ChangeRoleCommandHandler {
	return &ChangeRoleCommandHandler{
		userService: userService,
		logger:      logger,
	}
}

// Handle handles the change role command
func (h *ChangeRoleCommandHandler) Handle(ctx context.Context, cmd ChangeRoleCommand) error {
	h.logger.Infof("Handling change role command for user: %s to role: %s", cmd.UserID, cmd.Role)

	switch cmd.Role {
	case userDomain.RoleAdmin:
		if err := h.userService.PromoteToAdmin(ctx, cmd.UserID); err != nil {
			return fmt.Errorf("failed to promote to admin: %w", err)
		}
	case userDomain.RoleModerator:
		if err := h.userService.PromoteToModerator(ctx, cmd.UserID); err != nil {
			return fmt.Errorf("failed to promote to moderator: %w", err)
		}
	case userDomain.RoleUser:
		if err := h.userService.DemoteToUser(ctx, cmd.UserID); err != nil {
			return fmt.Errorf("failed to demote to user: %w", err)
		}
	default:
		return fmt.Errorf("invalid role: %s", cmd.Role)
	}

	h.logger.Infof("User role changed successfully: %s to %s", cmd.UserID, cmd.Role)
	return nil
}

// Helper method to convert domain user to DTO
func (h *CreateUserCommandHandler) toUserDTO(user *userDomain.User) *dto.UserDTO {
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

// Helper method to convert domain user to DTO for update handler
func (h *UpdateUserCommandHandler) toUserDTO(user *userDomain.User) *dto.UserDTO {
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
