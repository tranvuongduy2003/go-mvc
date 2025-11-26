package dto

import (
	"time"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/user"
)

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

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone,omitempty"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListUsersRequest struct {
	Page     int    `json:"page" validate:"min=1"`
	Limit    int    `json:"limit" validate:"min=1,max=100"`
	Search   string `json:"search"`
	SortBy   string `json:"sort_by"`
	SortDir  string `json:"sort_dir"`
	IsActive *bool  `json:"is_active"`
}

type ListUsersResponse struct {
	Users      []UserResponse `json:"users"`
	Pagination PaginationDTO  `json:"pagination"`
}

type PaginationDTO struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
	Pages    int   `json:"pages"`
}

func UserResponseFromDomain(u *user.User) UserResponse {
	return UserResponse{
		ID:        u.ID(),
		Email:     u.Email(),
		Name:      u.Name(),
		Phone:     u.Phone(),
		AvatarURL: u.Avatar().CDNUrl(),
		IsActive:  u.IsActive(),
		CreatedAt: u.CreatedAt(),
		UpdatedAt: u.UpdatedAt(),
	}
}

func UserResponseListFromDomain(users []*user.User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, u := range users {
		responses[i] = UserResponseFromDomain(u)
	}
	return responses
}
