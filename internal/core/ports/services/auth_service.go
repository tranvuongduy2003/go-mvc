package services

import (
	"context"
	"time"

	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/user"
)

// AuthTokens represents authentication tokens
type AuthTokens struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	TokenType             string    `json:"token_type"`
}

// LoginCredentials represents login credentials
type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest represents user registration data
type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// AuthenticatedUser represents an authenticated user with tokens
type AuthenticatedUser struct {
	User   *user.User  `json:"user"`
	Tokens *AuthTokens `json:"tokens"`
}

// AuthService defines the interface for authentication operations
type AuthService interface {
	// Register creates a new user account
	Register(ctx context.Context, req *RegisterRequest) (*AuthenticatedUser, error)

	// Login authenticates a user with email and password
	Login(ctx context.Context, credentials *LoginCredentials) (*AuthenticatedUser, error)

	// RefreshToken generates new access token using refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*AuthTokens, error)

	// Logout invalidates user tokens
	Logout(ctx context.Context, userID string) error

	// LogoutAll invalidates all tokens for a user across all devices
	LogoutAll(ctx context.Context, userID string) error

	// ValidateToken validates an access token and returns user info
	ValidateToken(ctx context.Context, accessToken string) (*user.User, error)

	// ChangePassword changes user password
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error

	// ResetPassword initiates password reset process
	ResetPassword(ctx context.Context, email string) error

	// ConfirmPasswordReset completes password reset with token
	ConfirmPasswordReset(ctx context.Context, token, newPassword string) error

	// VerifyEmail verifies user email with verification token
	VerifyEmail(ctx context.Context, token string) error

	// ResendVerificationEmail resends email verification
	ResendVerificationEmail(ctx context.Context, email string) error

	// GetUserFromToken extracts user information from a valid token
	GetUserFromToken(ctx context.Context, token string) (*user.User, error)

	// IsTokenBlacklisted checks if a token is blacklisted
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)

	// BlacklistToken adds a token to blacklist
	BlacklistToken(ctx context.Context, token string) error
}

// AuthorizationService defines the interface for authorization operations
type AuthorizationService interface {
	// UserHasPermission checks if user has specific permission
	UserHasPermission(ctx context.Context, userID, resource, action string) (bool, error)

	// UserHasPermissionByName checks if user has specific permission by name
	UserHasPermissionByName(ctx context.Context, userID, permissionName string) (bool, error)

	// UserHasRole checks if user has specific role
	UserHasRole(ctx context.Context, userID, roleName string) (bool, error)

	// UserHasAnyRole checks if user has any of the specified roles
	UserHasAnyRole(ctx context.Context, userID string, roleNames []string) (bool, error)

	// UserHasAllRoles checks if user has all specified roles
	UserHasAllRoles(ctx context.Context, userID string, roleNames []string) (bool, error)

	// GetUserPermissions retrieves all permissions for a user
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)

	// GetUserRoles retrieves all roles for a user
	GetUserRoles(ctx context.Context, userID string) ([]string, error)

	// CheckMultiplePermissions checks multiple permissions at once
	CheckMultiplePermissions(ctx context.Context, userID string, permissions []string) (map[string]bool, error)

	// CanAccessResource checks if user can access a resource with specific action
	CanAccessResource(ctx context.Context, userID, resource, action string) error

	// IsAdmin checks if user has admin privileges
	IsAdmin(ctx context.Context, userID string) (bool, error)

	// IsModerator checks if user has moderator privileges
	IsModerator(ctx context.Context, userID string) (bool, error)

	// GetEffectivePermissions gets all effective permissions (through all roles)
	GetEffectivePermissions(ctx context.Context, userID string) ([]PermissionInfo, error)
}

// PermissionInfo represents permission information
type PermissionInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
	GrantedBy   string `json:"granted_by,omitempty"`
}
