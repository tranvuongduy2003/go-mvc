package models

import (
	"time"
)

// PermissionModel represents the GORM model for Permission entity
type PermissionModel struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null;size:100" json:"name"`
	Resource    string    `gorm:"not null;size:50" json:"resource"`
	Action      string    `gorm:"not null;size:50" json:"action"`
	Description string    `gorm:"size:255" json:"description"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Version     int64     `gorm:"default:1" json:"version"`
}

// TableName returns the table name for the PermissionModel
func (PermissionModel) TableName() string {
	return "permissions"
}
