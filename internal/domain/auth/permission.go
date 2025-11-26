package auth

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/shared/events"
)

// Permission represents the permission aggregate root for RBAC
type Permission struct {
	id          PermissionID
	name        PermissionName
	resource    Resource
	action      Action
	description string
	isActive    bool
	createdAt   time.Time
	updatedAt   time.Time
	version     int64
	events      []events.DomainEvent
}

// PermissionID value object
type PermissionID struct {
	value string
}

// NewPermissionID creates a new permission ID
func NewPermissionID() PermissionID {
	return PermissionID{value: uuid.New().String()}
}

// NewPermissionIDFromString creates a permission ID from string
func NewPermissionIDFromString(id string) (PermissionID, error) {
	if id == "" {
		return PermissionID{}, errors.New("permission ID cannot be empty")
	}
	if _, err := uuid.Parse(id); err != nil {
		return PermissionID{}, errors.New("invalid permission ID format")
	}
	return PermissionID{value: id}, nil
}

// String returns the string representation of permission ID
func (id PermissionID) String() string {
	return id.value
}

// Equals checks if two permission IDs are equal
func (id PermissionID) Equals(other PermissionID) bool {
	return id.value == other.value
}

// PermissionName value object
type PermissionName struct {
	value string
}

// NewPermissionName creates a new permission name
func NewPermissionName(name string) (PermissionName, error) {
	if err := validatePermissionName(name); err != nil {
		return PermissionName{}, err
	}
	return PermissionName{value: strings.TrimSpace(name)}, nil
}

// String returns the string representation of permission name
func (p PermissionName) String() string {
	return p.value
}

// Equals checks if two permission names are equal
func (p PermissionName) Equals(other PermissionName) bool {
	return p.value == other.value
}

// Resource value object represents what the permission applies to
type Resource struct {
	value string
}

// NewResource creates a new resource
func NewResource(resource string) (Resource, error) {
	if err := validateResource(resource); err != nil {
		return Resource{}, err
	}
	return Resource{value: strings.ToLower(strings.TrimSpace(resource))}, nil
}

// String returns the string representation of resource
func (r Resource) String() string {
	return r.value
}

// Equals checks if two resources are equal
func (r Resource) Equals(other Resource) bool {
	return r.value == other.value
}

// Action value object represents what action can be performed
type Action struct {
	value string
}

// NewAction creates a new action
func NewAction(action string) (Action, error) {
	if err := validateAction(action); err != nil {
		return Action{}, err
	}
	return Action{value: strings.ToLower(strings.TrimSpace(action))}, nil
}

// String returns the string representation of action
func (a Action) String() string {
	return a.value
}

// Equals checks if two actions are equal
func (a Action) Equals(other Action) bool {
	return a.value == other.value
}

// Validation functions
func validatePermissionName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return errors.New("permission name cannot be empty")
	}

	if len(name) < 3 {
		return errors.New("permission name must be at least 3 characters long")
	}

	if len(name) > 100 {
		return errors.New("permission name cannot be longer than 100 characters")
	}

	// Permission names should follow format: resource:action (e.g., users:read, posts:create)
	validFormat := regexp.MustCompile(`^[a-z][a-z0-9_]*:[a-z][a-z0-9_]*$`)
	if !validFormat.MatchString(name) {
		return errors.New("permission name must follow format 'resource:action' with lowercase letters, numbers, and underscores")
	}

	return nil
}

func validateResource(resource string) error {
	resource = strings.TrimSpace(resource)

	if resource == "" {
		return errors.New("resource cannot be empty")
	}

	if len(resource) < 2 {
		return errors.New("resource must be at least 2 characters long")
	}

	if len(resource) > 50 {
		return errors.New("resource cannot be longer than 50 characters")
	}

	// Resources should be lowercase, alphanumeric with underscores
	validFormat := regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
	if !validFormat.MatchString(resource) {
		return errors.New("resource must start with lowercase letter and contain only lowercase letters, numbers, and underscores")
	}

	return nil
}

func validateAction(action string) error {
	action = strings.TrimSpace(action)

	if action == "" {
		return errors.New("action cannot be empty")
	}

	if len(action) < 2 {
		return errors.New("action must be at least 2 characters long")
	}

	if len(action) > 50 {
		return errors.New("action cannot be longer than 50 characters")
	}

	// Actions should be lowercase, alphanumeric with underscores
	validFormat := regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
	if !validFormat.MatchString(action) {
		return errors.New("action must start with lowercase letter and contain only lowercase letters, numbers, and underscores")
	}

	// Check for standard CRUD actions
	standardActions := []string{"create", "read", "update", "delete", "list", "manage", "execute", "view", "edit", "publish", "approve"}
	isStandard := false
	for _, std := range standardActions {
		if action == std || strings.HasPrefix(action, std+"_") || strings.HasSuffix(action, "_"+std) {
			isStandard = true
			break
		}
	}

	if !isStandard {
		return errors.New("action should be a standard CRUD operation or contain standard action words")
	}

	return nil
}

// NewPermission creates a new permission with validation
func NewPermission(name, resource, action, description string) (*Permission, error) {
	permissionName, err := NewPermissionName(name)
	if err != nil {
		return nil, err
	}

	resourceObj, err := NewResource(resource)
	if err != nil {
		return nil, err
	}

	actionObj, err := NewAction(action)
	if err != nil {
		return nil, err
	}

	// Ensure permission name matches resource:action pattern
	expectedName := resourceObj.String() + ":" + actionObj.String()
	if permissionName.String() != expectedName {
		return nil, errors.New("permission name must match 'resource:action' pattern")
	}

	if len(description) > 255 {
		return nil, errors.New("permission description cannot be longer than 255 characters")
	}

	now := time.Now()
	permission := &Permission{
		id:          NewPermissionID(),
		name:        permissionName,
		resource:    resourceObj,
		action:      actionObj,
		description: strings.TrimSpace(description),
		isActive:    true,
		createdAt:   now,
		updatedAt:   now,
		version:     1,
		events:      make([]events.DomainEvent, 0),
	}

	// Add domain event
	permissionUUID, _ := uuid.Parse(permission.id.String())
	event := events.NewBaseDomainEvent("permission.created", permissionUUID, "permission", map[string]interface{}{
		"permission_id": permission.id.String(),
		"name":          permission.name.String(),
		"resource":      permission.resource.String(),
		"action":        permission.action.String(),
		"description":   permission.description,
		"created_at":    permission.createdAt,
	})
	permission.events = append(permission.events, event)

	return permission, nil
}

// Getters
func (p *Permission) ID() PermissionID {
	return p.id
}

func (p *Permission) Name() PermissionName {
	return p.name
}

func (p *Permission) Resource() Resource {
	return p.resource
}

func (p *Permission) Action() Action {
	return p.action
}

func (p *Permission) Description() string {
	return p.description
}

func (p *Permission) IsActive() bool {
	return p.isActive
}

func (p *Permission) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Permission) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p *Permission) Version() int64 {
	return p.version
}

func (p *Permission) DomainEvents() []events.DomainEvent {
	return p.events
}

// Business methods
func (p *Permission) UpdateDescription(newDescription string) error {
	if len(newDescription) > 255 {
		return errors.New("permission description cannot be longer than 255 characters")
	}

	oldDescription := p.description
	p.description = strings.TrimSpace(newDescription)
	p.updatedAt = time.Now()
	p.version++

	// Add domain event
	permissionUUID, _ := uuid.Parse(p.id.String())
	event := events.NewBaseDomainEvent("permission.description_updated", permissionUUID, "permission", map[string]interface{}{
		"permission_id":   p.id.String(),
		"old_description": oldDescription,
		"new_description": p.description,
		"updated_at":      p.updatedAt,
	})
	p.events = append(p.events, event)

	return nil
}

func (p *Permission) Activate() {
	if !p.isActive {
		p.isActive = true
		p.updatedAt = time.Now()
		p.version++

		// Add domain event
		permissionUUID, _ := uuid.Parse(p.id.String())
		event := events.NewBaseDomainEvent("permission.activated", permissionUUID, "permission", map[string]interface{}{
			"permission_id": p.id.String(),
			"name":          p.name.String(),
			"updated_at":    p.updatedAt,
		})
		p.events = append(p.events, event)
	}
}

func (p *Permission) Deactivate() {
	if p.isActive {
		p.isActive = false
		p.updatedAt = time.Now()
		p.version++

		// Add domain event
		permissionUUID, _ := uuid.Parse(p.id.String())
		event := events.NewBaseDomainEvent("permission.deactivated", permissionUUID, "permission", map[string]interface{}{
			"permission_id": p.id.String(),
			"name":          p.name.String(),
			"updated_at":    p.updatedAt,
		})
		p.events = append(p.events, event)
	}
}

// ClearEvents clears all domain events
func (p *Permission) ClearEvents() {
	p.events = make([]events.DomainEvent, 0)
}

// Equals checks if two permissions are equal
func (p *Permission) Equals(other *Permission) bool {
	if other == nil {
		return false
	}
	return p.id.Equals(other.id)
}

// String returns string representation of the permission
func (p *Permission) String() string {
	return p.name.String()
}

// AppliesTo checks if this permission applies to a specific resource and action
func (p *Permission) AppliesTo(resource, action string) bool {
	return p.isActive &&
		p.resource.String() == strings.ToLower(resource) &&
		p.action.String() == strings.ToLower(action)
}

// Factory function for reconstruction from persistence
func ReconstructPermission(id, name, resource, action, description string, isActive bool, createdAt, updatedAt time.Time, version int64) (*Permission, error) {
	permissionID, err := NewPermissionIDFromString(id)
	if err != nil {
		return nil, err
	}

	permissionName, err := NewPermissionName(name)
	if err != nil {
		return nil, err
	}

	resourceObj, err := NewResource(resource)
	if err != nil {
		return nil, err
	}

	actionObj, err := NewAction(action)
	if err != nil {
		return nil, err
	}

	return &Permission{
		id:          permissionID,
		name:        permissionName,
		resource:    resourceObj,
		action:      actionObj,
		description: description,
		isActive:    isActive,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		version:     version,
		events:      make([]events.DomainEvent, 0),
	}, nil
}
