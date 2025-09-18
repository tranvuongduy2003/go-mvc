package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/shared/valueobject"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/security"
)

// Service provides domain services for user operations
type Service struct {
	repository     Repository
	passwordHasher *security.PasswordHasher
	logger         *logger.Logger
}

// NewService creates a new user domain service
func NewService(repository Repository, passwordHasher *security.PasswordHasher, logger *logger.Logger) *Service {
	return &Service{
		repository:     repository,
		passwordHasher: passwordHasher,
		logger:         logger,
	}
}

// CreateUser creates a new user with validation
func (s *Service) CreateUser(ctx context.Context, email valueobject.Email, username, password, firstName, lastName string) (*User, error) {
	// Check if user already exists
	if exists, err := s.repository.ExistsByEmail(ctx, email); err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	} else if exists {
		return nil, ErrEmailAlreadyExists
	}

	if exists, err := s.repository.ExistsByUsername(ctx, username); err != nil {
		return nil, fmt.Errorf("failed to check username existence: %w", err)
	} else if exists {
		return nil, ErrUsernameAlreadyExists
	}

	// Validate password strength
	if err := security.ValidatePassword(password); err != nil {
		return nil, fmt.Errorf("password validation failed: %w", err)
	}

	// Hash password
	hashedPassword, err := s.passwordHasher.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user entity
	user, err := NewUser(email, username, hashedPassword, firstName, lastName)
	if err != nil {
		return nil, fmt.Errorf("failed to create user entity: %w", err)
	}

	// Create user profile
	profile := NewProfile(user.ID())
	user.SetProfile(profile)

	// Save to repository
	if err := s.repository.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// Create profile
	if err := s.repository.CreateProfile(ctx, profile); err != nil {
		s.logger.Errorf("Failed to create user profile: %v", err)
		// Don't fail the user creation if profile creation fails
	}

	s.logger.Infof("User created successfully: %s", user.ID())
	return user, nil
}

// ValidateCredentials validates user credentials
func (s *Service) ValidateCredentials(ctx context.Context, email valueobject.Email, password string) (*User, error) {
	user, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive() {
		return nil, ErrUserInactive
	}

	if user.IsDeleted() {
		return nil, ErrUserDeleted
	}

	if !s.passwordHasher.Verify(password, user.Password()) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

// ChangePassword changes user password with validation
func (s *Service) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify current password
	if !s.passwordHasher.Verify(currentPassword, user.Password()) {
		return ErrPasswordMismatch
	}

	// Validate new password
	if err := security.ValidatePassword(newPassword); err != nil {
		return fmt.Errorf("new password validation failed: %w", err)
	}

	// Hash new password
	hashedPassword, err := s.passwordHasher.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	if err := user.UpdatePassword(hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Save changes
	if err := s.repository.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to save password change: %w", err)
	}

	s.logger.Infof("Password changed for user: %s", userID)
	return nil
}

// ResetPassword resets user password (admin operation)
func (s *Service) ResetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Validate new password
	if err := security.ValidatePassword(newPassword); err != nil {
		return fmt.Errorf("password validation failed: %w", err)
	}

	// Hash new password
	hashedPassword, err := s.passwordHasher.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	if err := user.UpdatePassword(hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Save changes
	if err := s.repository.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to save password reset: %w", err)
	}

	s.logger.Infof("Password reset for user: %s", userID)
	return nil
}

// DeactivateUser deactivates a user
func (s *Service) DeactivateUser(ctx context.Context, userID uuid.UUID) error {
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if err := user.Deactivate(); err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	if err := s.repository.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to save user deactivation: %w", err)
	}

	s.logger.Infof("User deactivated: %s", userID)
	return nil
}

// ActivateUser activates a user
func (s *Service) ActivateUser(ctx context.Context, userID uuid.UUID) error {
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if err := user.Activate(); err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}

	if err := s.repository.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to save user activation: %w", err)
	}

	s.logger.Infof("User activated: %s", userID)
	return nil
}

// DeleteUser soft deletes a user
func (s *Service) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if err := user.Delete(); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if err := s.repository.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to save user deletion: %w", err)
	}

	s.logger.Infof("User deleted: %s", userID)
	return nil
}

// PromoteToAdmin promotes a user to admin role
func (s *Service) PromoteToAdmin(ctx context.Context, userID uuid.UUID) error {
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if err := user.ChangeRole(RoleAdmin); err != nil {
		return fmt.Errorf("failed to promote user to admin: %w", err)
	}

	if err := s.repository.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to save user promotion: %w", err)
	}

	s.logger.Infof("User promoted to admin: %s", userID)
	return nil
}

// PromoteToModerator promotes a user to moderator role
func (s *Service) PromoteToModerator(ctx context.Context, userID uuid.UUID) error {
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if err := user.ChangeRole(RoleModerator); err != nil {
		return fmt.Errorf("failed to promote user to moderator: %w", err)
	}

	if err := s.repository.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to save user promotion: %w", err)
	}

	s.logger.Infof("User promoted to moderator: %s", userID)
	return nil
}

// DemoteToUser demotes a user to regular user role
func (s *Service) DemoteToUser(ctx context.Context, userID uuid.UUID) error {
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if err := user.ChangeRole(RoleUser); err != nil {
		return fmt.Errorf("failed to demote user: %w", err)
	}

	if err := s.repository.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to save user demotion: %w", err)
	}

	s.logger.Infof("User demoted to user: %s", userID)
	return nil
}

// UpdateUserProfile updates user profile information
func (s *Service) UpdateUserProfile(ctx context.Context, userID uuid.UUID, firstName, lastName string) error {
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if err := user.UpdateProfile(firstName, lastName); err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	if err := s.repository.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to save user profile update: %w", err)
	}

	s.logger.Infof("User profile updated: %s", userID)
	return nil
}

// IsEmailAvailable checks if email is available for use
func (s *Service) IsEmailAvailable(ctx context.Context, email valueobject.Email) (bool, error) {
	exists, err := s.repository.ExistsByEmail(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to check email availability: %w", err)
	}
	return !exists, nil
}

// IsUsernameAvailable checks if username is available for use
func (s *Service) IsUsernameAvailable(ctx context.Context, username string) (bool, error) {
	exists, err := s.repository.ExistsByUsername(ctx, username)
	if err != nil {
		return false, fmt.Errorf("failed to check username availability: %w", err)
	}
	return !exists, nil
}
