package dto

import (
	"time"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/user"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Phone    string `json:"phone" validate:"omitempty"`
	Password string `json:"password" validate:"required,min=8"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ConfirmPasswordResetRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type TokensDTO struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	TokenType             string    `json:"token_type"`
}

type AuthUserDTO struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginResponse struct {
	User   AuthUserDTO `json:"user"`
	Tokens TokensDTO   `json:"tokens"`
}

type RegisterResponse struct {
	User   AuthUserDTO `json:"user"`
	Tokens TokensDTO   `json:"tokens"`
}

type RefreshTokenResponse struct {
	Tokens TokensDTO `json:"tokens"`
}

type UserProfileResponse struct {
	User        AuthUserDTO         `json:"user"`
	Roles       []string            `json:"roles"`
	Permissions []PermissionInfoDTO `json:"permissions"`
}

type PermissionInfoDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
	GrantedBy   string `json:"granted_by,omitempty"`
}

type RoleInfoDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CheckPermissionRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	Resource string `json:"resource" validate:"required"`
	Action   string `json:"action" validate:"required"`
}

type CheckPermissionByNameRequest struct {
	UserID         string `json:"user_id" validate:"required"`
	PermissionName string `json:"permission_name" validate:"required"`
}

type CheckRoleRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	RoleName string `json:"role_name" validate:"required"`
}

type CheckMultiplePermissionsRequest struct {
	UserID      string   `json:"user_id" validate:"required"`
	Permissions []string `json:"permissions" validate:"required,min=1"`
}

type AssignRoleRequest struct {
	UserID     string     `json:"user_id" validate:"required"`
	RoleID     string     `json:"role_id" validate:"required"`
	AssignedBy string     `json:"assigned_by,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}

type RevokeRoleRequest struct {
	UserID string `json:"user_id" validate:"required"`
	RoleID string `json:"role_id" validate:"required"`
}

type CheckPermissionResponse struct {
	HasPermission bool `json:"has_permission"`
}

type CheckMultiplePermissionsResponse struct {
	Permissions map[string]bool `json:"permissions"`
}

type GetUserRolesResponse struct {
	Roles []string `json:"roles"`
}

type GetUserPermissionsResponse struct {
	Permissions []string `json:"permissions"`
}

type GetEffectivePermissionsResponse struct {
	Permissions []PermissionInfoDTO `json:"permissions"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type StatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

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

func ToPermissionInfoDTO(permissionInfo interface{}) PermissionInfoDTO {
	return PermissionInfoDTO{}
}

func ToPermissionInfoDTOs(permissionInfos []interface{}) []PermissionInfoDTO {
	dtos := make([]PermissionInfoDTO, len(permissionInfos))
	for i, info := range permissionInfos {
		dtos[i] = ToPermissionInfoDTO(info)
	}
	return dtos
}
