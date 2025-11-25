package models

import (
	"time"
)

// UserModel represents the GORM model for User entity
type UserModel struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Email         string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Name          string    `gorm:"not null;size:100" json:"name"`
	PasswordHash  string    `gorm:"column:password_hash;not null;size:255" json:"-"`
	Phone         string    `gorm:"size:20" json:"phone"`
	AvatarFileKey string    `gorm:"column:avatar_file_key;size:500" json:"avatar_file_key"`
	AvatarCDNUrl  string    `gorm:"column:avatar_cdn_url;size:1000" json:"avatar_cdn_url"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName returns the table name for the UserModel
func (UserModel) TableName() string {
	return "users"
}
