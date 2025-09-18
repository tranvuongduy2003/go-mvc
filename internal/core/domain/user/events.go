package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/shared/valueobject"
)

// DomainEvent represents a domain event
type DomainEvent interface {
	EventID() uuid.UUID
	EventType() string
	AggregateID() uuid.UUID
	OccurredAt() time.Time
	Version() int
}

// BaseDomainEvent provides common fields for domain events
type BaseDomainEvent struct {
	eventID     uuid.UUID
	eventType   string
	aggregateID uuid.UUID
	occurredAt  time.Time
	version     int
}

// NewBaseDomainEvent creates a new base domain event
func NewBaseDomainEvent(eventType string, aggregateID uuid.UUID, version int) BaseDomainEvent {
	return BaseDomainEvent{
		eventID:     uuid.New(),
		eventType:   eventType,
		aggregateID: aggregateID,
		occurredAt:  time.Now().UTC(),
		version:     version,
	}
}

// EventID returns the event ID
func (e BaseDomainEvent) EventID() uuid.UUID {
	return e.eventID
}

// EventType returns the event type
func (e BaseDomainEvent) EventType() string {
	return e.eventType
}

// AggregateID returns the aggregate ID
func (e BaseDomainEvent) AggregateID() uuid.UUID {
	return e.aggregateID
}

// OccurredAt returns when the event occurred
func (e BaseDomainEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// Version returns the event version
func (e BaseDomainEvent) Version() int {
	return e.version
}

// UserCreatedEvent represents a user creation event
type UserCreatedEvent struct {
	BaseDomainEvent
	Email     valueobject.Email
	Username  string
	FirstName string
	LastName  string
	Role      Role
}

// NewUserCreatedEvent creates a new user created event
func NewUserCreatedEvent(user *User) *UserCreatedEvent {
	return &UserCreatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("user.created", user.ID(), 1),
		Email:           user.Email(),
		Username:        user.Username(),
		FirstName:       user.FirstName(),
		LastName:        user.LastName(),
		Role:            user.Role(),
	}
}

// UserUpdatedEvent represents a user update event
type UserUpdatedEvent struct {
	BaseDomainEvent
	Changes map[string]interface{}
}

// NewUserUpdatedEvent creates a new user updated event
func NewUserUpdatedEvent(user *User, changes map[string]interface{}) *UserUpdatedEvent {
	return &UserUpdatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("user.updated", user.ID(), 1),
		Changes:         changes,
	}
}

// UserDeletedEvent represents a user deletion event
type UserDeletedEvent struct {
	BaseDomainEvent
	Email    valueobject.Email
	Username string
}

// NewUserDeletedEvent creates a new user deleted event
func NewUserDeletedEvent(user *User) *UserDeletedEvent {
	return &UserDeletedEvent{
		BaseDomainEvent: NewBaseDomainEvent("user.deleted", user.ID(), 1),
		Email:           user.Email(),
		Username:        user.Username(),
	}
}

// UserActivatedEvent represents a user activation event
type UserActivatedEvent struct {
	BaseDomainEvent
	Email    valueobject.Email
	Username string
}

// NewUserActivatedEvent creates a new user activated event
func NewUserActivatedEvent(user *User) *UserActivatedEvent {
	return &UserActivatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("user.activated", user.ID(), 1),
		Email:           user.Email(),
		Username:        user.Username(),
	}
}

// UserDeactivatedEvent represents a user deactivation event
type UserDeactivatedEvent struct {
	BaseDomainEvent
	Email    valueobject.Email
	Username string
}

// NewUserDeactivatedEvent creates a new user deactivated event
func NewUserDeactivatedEvent(user *User) *UserDeactivatedEvent {
	return &UserDeactivatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("user.deactivated", user.ID(), 1),
		Email:           user.Email(),
		Username:        user.Username(),
	}
}

// UserPasswordChangedEvent represents a password change event
type UserPasswordChangedEvent struct {
	BaseDomainEvent
	Email    valueobject.Email
	Username string
}

// NewUserPasswordChangedEvent creates a new password changed event
func NewUserPasswordChangedEvent(user *User) *UserPasswordChangedEvent {
	return &UserPasswordChangedEvent{
		BaseDomainEvent: NewBaseDomainEvent("user.password.changed", user.ID(), 1),
		Email:           user.Email(),
		Username:        user.Username(),
	}
}

// UserRoleChangedEvent represents a role change event
type UserRoleChangedEvent struct {
	BaseDomainEvent
	Email    valueobject.Email
	Username string
	OldRole  Role
	NewRole  Role
}

// NewUserRoleChangedEvent creates a new role changed event
func NewUserRoleChangedEvent(user *User, oldRole Role) *UserRoleChangedEvent {
	return &UserRoleChangedEvent{
		BaseDomainEvent: NewBaseDomainEvent("user.role.changed", user.ID(), 1),
		Email:           user.Email(),
		Username:        user.Username(),
		OldRole:         oldRole,
		NewRole:         user.Role(),
	}
}

// UserLoggedInEvent represents a user login event
type UserLoggedInEvent struct {
	BaseDomainEvent
	Email     valueobject.Email
	Username  string
	IPAddress string
	UserAgent string
}

// NewUserLoggedInEvent creates a new user logged in event
func NewUserLoggedInEvent(user *User, ipAddress, userAgent string) *UserLoggedInEvent {
	return &UserLoggedInEvent{
		BaseDomainEvent: NewBaseDomainEvent("user.logged.in", user.ID(), 1),
		Email:           user.Email(),
		Username:        user.Username(),
		IPAddress:       ipAddress,
		UserAgent:       userAgent,
	}
}

// UserLoggedOutEvent represents a user logout event
type UserLoggedOutEvent struct {
	BaseDomainEvent
	Email     valueobject.Email
	Username  string
	IPAddress string
}

// NewUserLoggedOutEvent creates a new user logged out event
func NewUserLoggedOutEvent(user *User, ipAddress string) *UserLoggedOutEvent {
	return &UserLoggedOutEvent{
		BaseDomainEvent: NewBaseDomainEvent("user.logged.out", user.ID(), 1),
		Email:           user.Email(),
		Username:        user.Username(),
		IPAddress:       ipAddress,
	}
}

// UserProfileUpdatedEvent represents a profile update event
type UserProfileUpdatedEvent struct {
	BaseDomainEvent
	Email    valueobject.Email
	Username string
	Changes  map[string]interface{}
}

// NewUserProfileUpdatedEvent creates a new profile updated event
func NewUserProfileUpdatedEvent(user *User, changes map[string]interface{}) *UserProfileUpdatedEvent {
	return &UserProfileUpdatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("user.profile.updated", user.ID(), 1),
		Email:           user.Email(),
		Username:        user.Username(),
		Changes:         changes,
	}
}

// UserEmailVerifiedEvent represents an email verification event
type UserEmailVerifiedEvent struct {
	BaseDomainEvent
	Email    valueobject.Email
	Username string
}

// NewUserEmailVerifiedEvent creates a new email verified event
func NewUserEmailVerifiedEvent(user *User) *UserEmailVerifiedEvent {
	return &UserEmailVerifiedEvent{
		BaseDomainEvent: NewBaseDomainEvent("user.email.verified", user.ID(), 1),
		Email:           user.Email(),
		Username:        user.Username(),
	}
}

// Event types constants
const (
	EventTypeUserCreated         = "user.created"
	EventTypeUserUpdated         = "user.updated"
	EventTypeUserDeleted         = "user.deleted"
	EventTypeUserActivated       = "user.activated"
	EventTypeUserDeactivated     = "user.deactivated"
	EventTypeUserPasswordChanged = "user.password.changed"
	EventTypeUserRoleChanged     = "user.role.changed"
	EventTypeUserLoggedIn        = "user.logged.in"
	EventTypeUserLoggedOut       = "user.logged.out"
	EventTypeUserProfileUpdated  = "user.profile.updated"
	EventTypeUserEmailVerified   = "user.email.verified"
)
