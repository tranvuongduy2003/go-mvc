package auth

import "time"

type UserRole struct {
	ID         string
	UserID     string
	RoleID     string
	AssignedBy *string
	AssignedAt time.Time
	ExpiresAt  *time.Time
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Version    int64
}
