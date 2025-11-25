package models

import (
	"time"
)

// UserRoleModel represents the GORM model for UserRole junction entity
type UserRoleModel struct {
	ID         string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID     string     `gorm:"type:uuid;not null;index" json:"user_id"`
	RoleID     string     `gorm:"type:uuid;not null;index" json:"role_id"`
	AssignedBy *string    `gorm:"type:uuid;index" json:"assigned_by,omitempty"`
	AssignedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"assigned_at"`
	ExpiresAt  *time.Time `gorm:"index" json:"expires_at,omitempty"`
	IsActive   bool       `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	Version    int64      `gorm:"default:1" json:"version"`

	// Relations
	User           UserModel  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Role           RoleModel  `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"role,omitempty"`
	AssignedByUser *UserModel `gorm:"foreignKey:AssignedBy;constraint:OnDelete:SET NULL" json:"assigned_by_user,omitempty"`
}

// TableName returns the table name for the UserRoleModel
func (UserRoleModel) TableName() string {
	return "user_roles"
}
