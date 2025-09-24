package dto

import (
	"time"

	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/user"
)

// Authentication Request DTOs

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Phone    string `json:"phone" validate:"omitempty"`
	Password string `json:"password" validate:"required,min=8"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// ResetPasswordRequest represents a password reset initiation request
type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ConfirmPasswordResetRequest represents a password reset confirmation request
type ConfirmPasswordResetRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// VerifyEmailRequest represents an email verification request
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

// ResendVerificationRequest represents a resend verification email request
type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Authentication Response DTOs

// TokensDTO represents authentication tokens
type TokensDTO struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	TokenType             string    `json:"token_type"`
}

// AuthUserDTO represents an authenticated user's basic info
type AuthUserDTO struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginResponse represents a successful login response
type LoginResponse struct {
	User   AuthUserDTO `json:"user"`
	Tokens TokensDTO   `json:"tokens"`
}

// RegisterResponse represents a successful registration response
type RegisterResponse struct {
	User   AuthUserDTO `json:"user"`
	Tokens TokensDTO   `json:"tokens"`
}

// RefreshTokenResponse represents a successful token refresh response
type RefreshTokenResponse struct {
	Tokens TokensDTO `json:"tokens"`
}

// UserProfileResponse represents the current user's profile
type UserProfileResponse struct {
	User        AuthUserDTO         `json:"user"`
	Roles       []string            `json:"roles"`
	Permissions []PermissionInfoDTO `json:"permissions"`
}

// PermissionInfoDTO represents permission information in responses
type PermissionInfoDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
	GrantedBy   string `json:"granted_by,omitempty"`
}

// RoleInfoDTO represents role information in responses
type RoleInfoDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Authorization DTOs

// CheckPermissionRequest represents a permission check request
type CheckPermissionRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	Resource string `json:"resource" validate:"required"`
	Action   string `json:"action" validate:"required"`
}

// CheckPermissionByNameRequest represents a permission check by name request
type CheckPermissionByNameRequest struct {
	UserID         string `json:"user_id" validate:"required"`
	PermissionName string `json:"permission_name" validate:"required"`
}

// CheckRoleRequest represents a role check request
type CheckRoleRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	RoleName string `json:"role_name" validate:"required"`
}

// CheckMultiplePermissionsRequest represents a multiple permissions check request
type CheckMultiplePermissionsRequest struct {
	UserID      string   `json:"user_id" validate:"required"`
	Permissions []string `json:"permissions" validate:"required,min=1"`
}

// AssignRoleRequest represents a role assignment request
type AssignRoleRequest struct {
	UserID     string     `json:"user_id" validate:"required"`
	RoleID     string     `json:"role_id" validate:"required"`
	AssignedBy string     `json:"assigned_by,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}

// RevokeRoleRequest represents a role revocation request
type RevokeRoleRequest struct {
	UserID string `json:"user_id" validate:"required"`
	RoleID string `json:"role_id" validate:"required"`
}

// Authorization Response DTOs

// CheckPermissionResponse represents a permission check response
type CheckPermissionResponse struct {
	HasPermission bool `json:"has_permission"`
}

// CheckMultiplePermissionsResponse represents a multiple permissions check response
type CheckMultiplePermissionsResponse struct {
	Permissions map[string]bool `json:"permissions"`
}

// GetUserRolesResponse represents user roles response
type GetUserRolesResponse struct {
	Roles []string `json:"roles"`
}

// GetUserPermissionsResponse represents user permissions response
type GetUserPermissionsResponse struct {
	Permissions []string `json:"permissions"`
}

// GetEffectivePermissionsResponse represents effective permissions response
type GetEffectivePermissionsResponse struct {
	Permissions []PermissionInfoDTO `json:"permissions"`
}

// General Response DTOs

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}

// StatusResponse represents a simple status response
type StatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// Conversion Functions

// ToAuthUserDTO converts a domain user to AuthUserDTO
func ToAuthUserDTO(domainUser *user.User) AuthUserDTO {
	return AuthUserDTO{
		ID:        domainUser.ID(),
		Email:     domainUser.Email(),
		Name:      domainUser.Name(),
		Phone:     domainUser.Phone(),
		IsActive:  domainUser.IsActive(),
		CreatedAt: domainUser.CreatedAt(),
		UpdatedAt: domainUser.UpdatedAt(),
	}
}

// ToPermissionInfoDTO converts service PermissionInfo to PermissionInfoDTO
func ToPermissionInfoDTO(permissionInfo interface{}) PermissionInfoDTO {
	// This would need to be implemented based on the actual PermissionInfo type
	// For now, return empty struct - would need actual implementation
	return PermissionInfoDTO{}
}

// ToPermissionInfoDTOs converts slice of service PermissionInfo to DTOs
func ToPermissionInfoDTOs(permissionInfos []interface{}) []PermissionInfoDTO {
	dtos := make([]PermissionInfoDTO, len(permissionInfos))
	for i, info := range permissionInfos {
		dtos[i] = ToPermissionInfoDTO(info)
	}
	return dtos
}
