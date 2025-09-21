package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role represents a role in the database
type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	DisplayName string    `gorm:"not null" json:"display_name"`
	Description string    `json:"description"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Many-to-many relationship with permissions
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`

	// Many-to-many relationship with users
	Users []User `gorm:"many2many:user_roles;" json:"users"`
}

// Permission represents a permission in the database
type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Resource    string    `gorm:"not null" json:"resource"`
	Action      string    `gorm:"not null" json:"action"`
	DisplayName string    `gorm:"not null" json:"display_name"`
	Description string    `json:"description"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Many-to-many relationship with roles
	Roles []Role `gorm:"many2many:role_permissions;" json:"roles"`
}

// UserRole represents the user-role relationship
type UserRole struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	RoleID    uuid.UUID `gorm:"type:uuid;not null;index" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Foreign key relationships
	User User `gorm:"foreignKey:UserID" json:"user"`
	Role Role `gorm:"foreignKey:RoleID" json:"role"`
}

// RolePermission represents the role-permission relationship
type RolePermission struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	RoleID       uuid.UUID `gorm:"type:uuid;not null;index" json:"role_id"`
	PermissionID uuid.UUID `gorm:"type:uuid;not null;index" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Foreign key relationships
	Role       Role       `gorm:"foreignKey:RoleID" json:"role"`
	Permission Permission `gorm:"foreignKey:PermissionID" json:"permission"`
}

// User represents a user in the database (simplified for RBAC)
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Many-to-many relationship with roles
	Roles []Role `gorm:"many2many:user_roles;" json:"roles"`
}

// BeforeCreate hook for Role
func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook for Permission
func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook for UserRole
func (ur *UserRole) BeforeCreate(tx *gorm.DB) error {
	if ur.ID == uuid.Nil {
		ur.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook for RolePermission
func (rp *RolePermission) BeforeCreate(tx *gorm.DB) error {
	if rp.ID == uuid.Nil {
		rp.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook for User
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for Role
func (Role) TableName() string {
	return "roles"
}

// TableName returns the table name for Permission
func (Permission) TableName() string {
	return "permissions"
}

// TableName returns the table name for UserRole
func (UserRole) TableName() string {
	return "user_roles"
}

// TableName returns the table name for RolePermission
func (RolePermission) TableName() string {
	return "role_permissions"
}

// TableName returns the table name for User
func (User) TableName() string {
	return "users"
}
