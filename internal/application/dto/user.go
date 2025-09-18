package dto

import (
	"time"

	"github.com/google/uuid"
)

// UserDTO represents a user data transfer object
type UserDTO struct {
	ID        uuid.UUID   `json:"id"`
	Email     string      `json:"email"`
	Username  string      `json:"username"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	FullName  string      `json:"full_name"`
	Role      string      `json:"role"`
	IsActive  bool        `json:"is_active"`
	Profile   *ProfileDTO `json:"profile,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// ProfileDTO represents a user profile data transfer object
type ProfileDTO struct {
	ID          uuid.UUID         `json:"id"`
	UserID      uuid.UUID         `json:"user_id"`
	Avatar      string            `json:"avatar,omitempty"`
	Bio         string            `json:"bio,omitempty"`
	DateOfBirth *time.Time        `json:"date_of_birth,omitempty"`
	Phone       string            `json:"phone,omitempty"`
	Address     string            `json:"address,omitempty"`
	City        string            `json:"city,omitempty"`
	Country     string            `json:"country,omitempty"`
	Website     string            `json:"website,omitempty"`
	SocialLinks map[string]string `json:"social_links,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// CreateUserRequest represents request to create a user
type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string `json:"last_name" validate:"max=100"`
}

// UpdateUserRequest represents request to update a user
type UpdateUserRequest struct {
	FirstName *string `json:"first_name" validate:"omitempty,min=1,max=100"`
	LastName  *string `json:"last_name" validate:"omitempty,max=100"`
	Email     *string `json:"email" validate:"omitempty,email"`
}

// ChangePasswordRequest represents request to change password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// ResetPasswordRequest represents request to reset password with token
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// UpdateProfileRequest represents request to update user profile
type UpdateProfileRequest struct {
	Avatar      *string           `json:"avatar" validate:"omitempty,url"`
	Bio         *string           `json:"bio" validate:"omitempty,max=500"`
	DateOfBirth *time.Time        `json:"date_of_birth"`
	Phone       *string           `json:"phone" validate:"omitempty"`
	Address     *string           `json:"address" validate:"omitempty,max=200"`
	City        *string           `json:"city" validate:"omitempty,max=100"`
	Country     *string           `json:"country" validate:"omitempty,max=100"`
	Website     *string           `json:"website" validate:"omitempty,url"`
	SocialLinks map[string]string `json:"social_links"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	User         *UserDTO `json:"user"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresAt    int64    `json:"expires_at"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ForgotPasswordRequest represents request to initiate password reset
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// VerifyEmailRequest represents request to verify email
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

// ResendVerificationRequest represents request to resend verification email
type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// UserListRequest represents request to list users
type UserListRequest struct {
	Page      int     `json:"page" validate:"min=1"`
	Limit     int     `json:"limit" validate:"min=1,max=100"`
	Sort      string  `json:"sort" validate:"omitempty,oneof=created_at updated_at email username first_name last_name"`
	Order     string  `json:"order" validate:"omitempty,oneof=asc desc"`
	Search    string  `json:"search"`
	Role      *string `json:"role" validate:"omitempty,oneof=admin user moderator"`
	IsActive  *bool   `json:"is_active"`
	IsDeleted *bool   `json:"is_deleted"`
}

// UserListResponse represents response for listing users
type UserListResponse struct {
	Users      []*UserDTO  `json:"users"`
	Pagination *Pagination `json:"pagination"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// CheckAvailabilityRequest represents request to check username/email availability
type CheckAvailabilityRequest struct {
	Email    *string `json:"email" validate:"omitempty,email"`
	Username *string `json:"username" validate:"omitempty,min=3,max=50"`
}

// CheckAvailabilityResponse represents response for availability check
type CheckAvailabilityResponse struct {
	EmailAvailable    *bool `json:"email_available,omitempty"`
	UsernameAvailable *bool `json:"username_available,omitempty"`
}

// ChangeRoleRequest represents request to change user role
type ChangeRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin user moderator"`
}

// UserStatsResponse represents user statistics
type UserStatsResponse struct {
	TotalUsers     int64 `json:"total_users"`
	ActiveUsers    int64 `json:"active_users"`
	InactiveUsers  int64 `json:"inactive_users"`
	AdminUsers     int64 `json:"admin_users"`
	ModeratorUsers int64 `json:"moderator_users"`
	RegularUsers   int64 `json:"regular_users"`
	UsersThisMonth int64 `json:"users_this_month"`
	UsersToday     int64 `json:"users_today"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Code    string            `json:"code,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// HealthCheckResponse represents health check response
type HealthCheckResponse struct {
	Status      string                 `json:"status"`
	Version     string                 `json:"version"`
	Environment string                 `json:"environment"`
	Timestamp   time.Time              `json:"timestamp"`
	Uptime      string                 `json:"uptime"`
	Services    map[string]interface{} `json:"services"`
}
