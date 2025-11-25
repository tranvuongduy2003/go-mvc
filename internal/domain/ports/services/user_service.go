package services

import "context"

// UserCommandService handles user write operations
// Following Single Responsibility Principle - separated from query operations
type UserCommandService interface {
	// CreateUser creates a new user account
	CreateUser(ctx context.Context, req CreateUserRequest) (UserResponse, error)

	// UpdateUser updates an existing user's profile
	UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) (UserResponse, error)

	// DeleteUser removes a user account
	DeleteUser(ctx context.Context, userID string) error

	// UploadAvatar uploads and updates user avatar
	UploadAvatar(ctx context.Context, userID string, req UploadAvatarRequest) (UserResponse, error)

	// ActivateUser activates a deactivated user
	ActivateUser(ctx context.Context, userID string) error

	// DeactivateUser deactivates an active user
	DeactivateUser(ctx context.Context, userID string) error
}

// UserQueryService handles user read operations
// Following Single Responsibility Principle - separated from command operations
type UserQueryService interface {
	// GetUserByID retrieves a user by their ID
	GetUserByID(ctx context.Context, userID string) (UserResponse, error)

	// GetUserByEmail retrieves a user by their email
	GetUserByEmail(ctx context.Context, email string) (UserResponse, error)

	// ListUsers retrieves a paginated list of users
	ListUsers(ctx context.Context, req ListUsersRequest) (ListUsersResponse, error)

	// CountUsers returns the total number of users
	CountUsers(ctx context.Context) (int64, error)

	// UserExists checks if a user exists by ID
	UserExists(ctx context.Context, userID string) (bool, error)
}

// Request/Response DTOs
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Phone    string `json:"phone" validate:"omitempty"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Phone string `json:"phone" validate:"omitempty"`
}

type UploadAvatarRequest struct {
	FileData    []byte `json:"-"`
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Phone     string `json:"phone,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ListUsersRequest struct {
	Page     int    `json:"page" validate:"min=1"`
	Limit    int    `json:"limit" validate:"min=1,max=100"`
	Search   string `json:"search"`
	SortBy   string `json:"sort_by"`
	SortDir  string `json:"sort_dir" validate:"oneof=asc desc"`
	IsActive *bool  `json:"is_active"`
}

type ListUsersResponse struct {
	Users      []UserResponse `json:"users"`
	Pagination PaginationInfo `json:"pagination"`
}

type PaginationInfo struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
	Pages    int   `json:"pages"`
}
