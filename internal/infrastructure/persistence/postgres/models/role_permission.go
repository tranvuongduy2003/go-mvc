package models

import (
	"time"
)

type RolePermissionModel struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	RoleID       string    `gorm:"type:uuid;not null;index" json:"role_id"`
	PermissionID string    `gorm:"type:uuid;not null;index" json:"permission_id"`
	GrantedBy    *string   `gorm:"type:uuid;index" json:"granted_by,omitempty"`
	GrantedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"granted_at"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Version      int64     `gorm:"default:1" json:"version"`

	Role          RoleModel       `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"role,omitempty"`
	Permission    PermissionModel `gorm:"foreignKey:PermissionID;constraint:OnDelete:CASCADE" json:"permission,omitempty"`
	GrantedByUser *UserModel      `gorm:"foreignKey:GrantedBy;constraint:OnDelete:SET NULL" json:"granted_by_user,omitempty"`
}

func (RolePermissionModel) TableName() string {
	return "role_permissions"
}
