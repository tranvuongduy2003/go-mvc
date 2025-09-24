package models

import (
	"time"
)

// RoleModel represents the GORM model for Role entity
type RoleModel struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null;size:50" json:"name"`
	Description string    `gorm:"size:255" json:"description"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Version     int64     `gorm:"default:1" json:"version"`
}

// TableName returns the table name for the RoleModel
func (RoleModel) TableName() string {
	return "roles"
}
