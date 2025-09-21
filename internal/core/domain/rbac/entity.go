package rbac

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/shared/valueobject"
)

// RBAC domain errors
var (
	ErrRoleNotFound       = errors.New("role not found")
	ErrPermissionNotFound = errors.New("permission not found")
	ErrUserRoleNotFound   = errors.New("user role not found")
	ErrRoleAlreadyExists  = errors.New("role already exists")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrInvalidRole        = errors.New("invalid role")
	ErrInvalidPermission  = errors.New("invalid permission")
)

// Role represents a role in the RBAC system
type Role struct {
	ID          uuid.UUID            `json:"id"`
	Name        string               `json:"name"`
	DisplayName string               `json:"display_name"`
	Description string               `json:"description"`
	Permissions []Permission         `json:"permissions"`
	UserIDs     []uuid.UUID          `json:"user_ids,omitempty"`
	IsActive    bool                 `json:"is_active"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	AuditLog    valueobject.AuditLog `json:"audit_log"`
}

// Permission represents a permission in the RBAC system
type Permission struct {
	ID          uuid.UUID            `json:"id"`
	Name        string               `json:"name"`
	Resource    string               `json:"resource"`
	Action      string               `json:"action"`
	Description string               `json:"description"`
	IsActive    bool                 `json:"is_active"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	AuditLog    valueobject.AuditLog `json:"audit_log"`
}

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	UserID     uuid.UUID            `json:"user_id"`
	RoleID     uuid.UUID            `json:"role_id"`
	AssignedBy uuid.UUID            `json:"assigned_by"`
	AssignedAt time.Time            `json:"assigned_at"`
	ExpiresAt  *time.Time           `json:"expires_at,omitempty"`
	IsActive   bool                 `json:"is_active"`
	AuditLog   valueobject.AuditLog `json:"audit_log"`
}

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	RoleID       uuid.UUID            `json:"role_id"`
	PermissionID uuid.UUID            `json:"permission_id"`
	GrantedBy    uuid.UUID            `json:"granted_by"`
	GrantedAt    time.Time            `json:"granted_at"`
	IsActive     bool                 `json:"is_active"`
	AuditLog     valueobject.AuditLog `json:"audit_log"`
}

// NewRole creates a new role
func NewRole(name, description string, createdBy uuid.UUID) *Role {
	return &Role{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Permissions: []Permission{},
		UserIDs:     []uuid.UUID{},
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		AuditLog: valueobject.AuditLog{
			CreatedBy: createdBy,
			UpdatedBy: createdBy,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// NewPermission creates a new permission
func NewPermission(name, resource, action, description string, createdBy uuid.UUID) *Permission {
	return &Permission{
		ID:          uuid.New(),
		Name:        name,
		Resource:    resource,
		Action:      action,
		Description: description,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		AuditLog: valueobject.AuditLog{
			CreatedBy: createdBy,
			UpdatedBy: createdBy,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// NewUserRole creates a new user role assignment
func NewUserRole(userID, roleID, assignedBy uuid.UUID, expiresAt *time.Time) *UserRole {
	return &UserRole{
		UserID:     userID,
		RoleID:     roleID,
		AssignedBy: assignedBy,
		AssignedAt: time.Now(),
		ExpiresAt:  expiresAt,
		IsActive:   true,
		AuditLog: valueobject.AuditLog{
			CreatedBy: assignedBy,
			UpdatedBy: assignedBy,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// NewRolePermission creates a new role permission assignment
func NewRolePermission(roleID, permissionID, grantedBy uuid.UUID) *RolePermission {
	return &RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
		GrantedBy:    grantedBy,
		GrantedAt:    time.Now(),
		IsActive:     true,
		AuditLog: valueobject.AuditLog{
			CreatedBy: grantedBy,
			UpdatedBy: grantedBy,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// HasPermission checks if the role has a specific permission
func (r *Role) HasPermission(resource, action string) bool {
	for _, permission := range r.Permissions {
		if permission.IsActive &&
			permission.Resource == resource &&
			permission.Action == action {
			return true
		}
	}
	return false
}

// AddPermission adds a permission to the role
func (r *Role) AddPermission(permission Permission, grantedBy uuid.UUID) {
	// Check if permission already exists
	for _, p := range r.Permissions {
		if p.ID == permission.ID {
			return
		}
	}

	r.Permissions = append(r.Permissions, permission)
	r.UpdatedAt = time.Now()
	r.AuditLog.UpdatedBy = grantedBy
	r.AuditLog.UpdatedAt = time.Now()
}

// RemovePermission removes a permission from the role
func (r *Role) RemovePermission(permissionID uuid.UUID, removedBy uuid.UUID) {
	for i, permission := range r.Permissions {
		if permission.ID == permissionID {
			r.Permissions = append(r.Permissions[:i], r.Permissions[i+1:]...)
			r.UpdatedAt = time.Now()
			r.AuditLog.UpdatedBy = removedBy
			r.AuditLog.UpdatedAt = time.Now()
			break
		}
	}
}

// IsExpired checks if a user role assignment is expired
func (ur *UserRole) IsExpired() bool {
	if ur.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*ur.ExpiresAt)
}

// IsValid checks if a user role assignment is valid
func (ur *UserRole) IsValid() bool {
	return ur.IsActive && !ur.IsExpired()
}

// Common permission constants
const (
	// Resources
	ResourceUser   = "user"
	ResourceRole   = "role"
	ResourceSystem = "system"
	ResourceReport = "report"
	ResourceAudit  = "audit"

	// Actions
	ActionCreate = "create"
	ActionRead   = "read"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionList   = "list"
	ActionExport = "export"
	ActionAdmin  = "admin"
)

// Common role names
const (
	RoleAdmin     = "admin"
	RoleUser      = "user"
	RoleModerator = "moderator"
	RoleGuest     = "guest"
)
