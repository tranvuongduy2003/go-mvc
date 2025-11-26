package auth

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/shared/events"
)

// Role represents the role aggregate root for RBAC
type Role struct {
	id          RoleID
	name        RoleName
	description string
	isActive    bool
	createdAt   time.Time
	updatedAt   time.Time
	version     int64
	events      []events.DomainEvent
}

// RoleID value object
type RoleID struct {
	value string
}

// NewRoleID creates a new role ID
func NewRoleID() RoleID {
	return RoleID{value: uuid.New().String()}
}

// NewRoleIDFromString creates a role ID from string
func NewRoleIDFromString(id string) (RoleID, error) {
	if id == "" {
		return RoleID{}, errors.New("role ID cannot be empty")
	}
	if _, err := uuid.Parse(id); err != nil {
		return RoleID{}, errors.New("invalid role ID format")
	}
	return RoleID{value: id}, nil
}

// String returns the string representation of role ID
func (id RoleID) String() string {
	return id.value
}

// Equals checks if two role IDs are equal
func (id RoleID) Equals(other RoleID) bool {
	return id.value == other.value
}

// RoleName value object
type RoleName struct {
	value string
}

// NewRoleName creates a new role name
func NewRoleName(name string) (RoleName, error) {
	if err := validateRoleName(name); err != nil {
		return RoleName{}, err
	}
	return RoleName{value: strings.TrimSpace(name)}, nil
}

// String returns the string representation of role name
func (r RoleName) String() string {
	return r.value
}

// Equals checks if two role names are equal
func (r RoleName) Equals(other RoleName) bool {
	return r.value == other.value
}

// validateRoleName validates role name format and business rules
func validateRoleName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return errors.New("role name cannot be empty")
	}

	if len(name) < 2 {
		return errors.New("role name must be at least 2 characters long")
	}

	if len(name) > 50 {
		return errors.New("role name cannot be longer than 50 characters")
	}

	// Role names should be uppercase, alphanumeric with underscores
	validFormat := regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	if !validFormat.MatchString(name) {
		return errors.New("role name must start with uppercase letter and contain only uppercase letters, numbers, and underscores")
	}

	return nil
}

// NewRole creates a new role with basic validation
func NewRole(name string, description string) (*Role, error) {
	roleName, err := NewRoleName(name)
	if err != nil {
		return nil, err
	}

	if len(description) > 255 {
		return nil, errors.New("role description cannot be longer than 255 characters")
	}

	now := time.Now()
	role := &Role{
		id:          NewRoleID(),
		name:        roleName,
		description: strings.TrimSpace(description),
		isActive:    true,
		createdAt:   now,
		updatedAt:   now,
		version:     1,
		events:      make([]events.DomainEvent, 0),
	}

	// Add domain event
	roleUUID, _ := uuid.Parse(role.id.String())
	event := events.NewBaseDomainEvent("role.created", roleUUID, "role", map[string]interface{}{
		"role_id":     role.id.String(),
		"role_name":   role.name.String(),
		"description": role.description,
		"created_at":  role.createdAt,
	})
	role.events = append(role.events, event)

	return role, nil
}

// Getters
func (r *Role) ID() RoleID {
	return r.id
}

func (r *Role) Name() RoleName {
	return r.name
}

func (r *Role) Description() string {
	return r.description
}

func (r *Role) IsActive() bool {
	return r.isActive
}

func (r *Role) CreatedAt() time.Time {
	return r.createdAt
}

func (r *Role) UpdatedAt() time.Time {
	return r.updatedAt
}

func (r *Role) Version() int64 {
	return r.version
}

func (r *Role) DomainEvents() []events.DomainEvent {
	return r.events
}

// Business methods
func (r *Role) UpdateDescription(newDescription string) error {
	if len(newDescription) > 255 {
		return errors.New("role description cannot be longer than 255 characters")
	}

	oldDescription := r.description
	r.description = strings.TrimSpace(newDescription)
	r.updatedAt = time.Now()
	r.version++

	// Add domain event
	roleUUID, _ := uuid.Parse(r.id.String())
	event := events.NewBaseDomainEvent("role.description_updated", roleUUID, "role", map[string]interface{}{
		"role_id":         r.id.String(),
		"old_description": oldDescription,
		"new_description": r.description,
		"updated_at":      r.updatedAt,
	})
	r.events = append(r.events, event)

	return nil
}

func (r *Role) Activate() {
	if !r.isActive {
		r.isActive = true
		r.updatedAt = time.Now()
		r.version++

		// Add domain event
		roleUUID, _ := uuid.Parse(r.id.String())
		event := events.NewBaseDomainEvent("role.activated", roleUUID, "role", map[string]interface{}{
			"role_id":    r.id.String(),
			"role_name":  r.name.String(),
			"updated_at": r.updatedAt,
		})
		r.events = append(r.events, event)
	}
}

func (r *Role) Deactivate() {
	if r.isActive {
		r.isActive = false
		r.updatedAt = time.Now()
		r.version++

		// Add domain event
		roleUUID, _ := uuid.Parse(r.id.String())
		event := events.NewBaseDomainEvent("role.deactivated", roleUUID, "role", map[string]interface{}{
			"role_id":    r.id.String(),
			"role_name":  r.name.String(),
			"updated_at": r.updatedAt,
		})
		r.events = append(r.events, event)
	}
}

// ClearEvents clears all domain events
func (r *Role) ClearEvents() {
	r.events = make([]events.DomainEvent, 0)
}

// Equals checks if two roles are equal
func (r *Role) Equals(other *Role) bool {
	if other == nil {
		return false
	}
	return r.id.Equals(other.id)
}

// String returns string representation of the role
func (r *Role) String() string {
	return r.name.String()
}

// Factory function for reconstruction from persistence
func ReconstructRole(id, name, description string, isActive bool, createdAt, updatedAt time.Time, version int64) (*Role, error) {
	roleID, err := NewRoleIDFromString(id)
	if err != nil {
		return nil, err
	}

	roleName, err := NewRoleName(name)
	if err != nil {
		return nil, err
	}

	return &Role{
		id:          roleID,
		name:        roleName,
		description: description,
		isActive:    isActive,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		version:     version,
		events:      make([]events.DomainEvent, 0),
	}, nil
}
