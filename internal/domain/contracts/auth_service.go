package contracts

import (
	"context"
	"time"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/user"
)

type AuthTokens struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	TokenType             string    `json:"token_type"`
}

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type AuthenticatedUser struct {
	User   *user.User  `json:"user"`
	Tokens *AuthTokens `json:"tokens"`
}

type AuthService interface {
	Register(ctx context.Context, req *RegisterRequest) (*AuthenticatedUser, error)

	Login(ctx context.Context, credentials *LoginCredentials) (*AuthenticatedUser, error)

	RefreshToken(ctx context.Context, refreshToken string) (*AuthTokens, error)

	ValidateToken(ctx context.Context, accessToken string) (*user.User, error)

	GetUserFromToken(ctx context.Context, token string) (*user.User, error)
}

type TokenManagementService interface {
	Logout(ctx context.Context, userID string) error

	LogoutAll(ctx context.Context, userID string) error

	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)

	BlacklistToken(ctx context.Context, token string) error
}

type PasswordManagementService interface {
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error

	ResetPassword(ctx context.Context, email string) error

	ConfirmPasswordReset(ctx context.Context, token, newPassword string) error
}

type EmailVerificationService interface {
	VerifyEmail(ctx context.Context, token string) error

	ResendVerificationEmail(ctx context.Context, email string) error
}

type AuthorizationService interface {
	UserHasPermission(ctx context.Context, userID, resource, action string) (bool, error)

	UserHasPermissionByName(ctx context.Context, userID, permissionName string) (bool, error)

	UserHasRole(ctx context.Context, userID, roleName string) (bool, error)

	UserHasAnyRole(ctx context.Context, userID string, roleNames []string) (bool, error)

	UserHasAllRoles(ctx context.Context, userID string, roleNames []string) (bool, error)

	GetUserPermissions(ctx context.Context, userID string) ([]string, error)

	GetUserRoles(ctx context.Context, userID string) ([]string, error)

	CheckMultiplePermissions(ctx context.Context, userID string, permissions []string) (map[string]bool, error)

	CanAccessResource(ctx context.Context, userID, resource, action string) error

	IsAdmin(ctx context.Context, userID string) (bool, error)

	IsModerator(ctx context.Context, userID string) (bool, error)

	GetEffectivePermissions(ctx context.Context, userID string) ([]PermissionInfo, error)
}

type PermissionInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
	GrantedBy   string `json:"granted_by,omitempty"`
}
