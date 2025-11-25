package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// UserCreatedEvent represents a user creation domain event
type UserCreatedEvent struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	AggregateID_ string    `json:"aggregate_id"`
	Version_     int       `json:"version"`
	OccurredAt   time.Time `json:"occurred_at"`
}

// NewUserCreatedEvent creates a new user created event
func NewUserCreatedEvent(userID, email, fullName string) *UserCreatedEvent {
	return &UserCreatedEvent{
		ID:           uuid.New().String(),
		UserID:       userID,
		Email:        email,
		FullName:     fullName,
		AggregateID_: userID,
		Version_:     1,
		OccurredAt:   time.Now(),
	}
}

// EventType returns the event type
func (e *UserCreatedEvent) EventType() string {
	return "user.created"
}

// EventData returns the event data as JSON bytes
func (e *UserCreatedEvent) EventData() ([]byte, error) {
	return json.Marshal(e)
}

// AggregateID returns the aggregate ID
func (e *UserCreatedEvent) AggregateID() string {
	return e.AggregateID_
}

// Version returns the event version
func (e *UserCreatedEvent) Version() int {
	return e.Version_
}

// Timestamp returns the event timestamp
func (e *UserCreatedEvent) Timestamp() int64 {
	return e.OccurredAt.Unix()
}

// UserUpdatedEvent represents a user update domain event
type UserUpdatedEvent struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	AggregateID_ string    `json:"aggregate_id"`
	Version_     int       `json:"version"`
	OccurredAt   time.Time `json:"occurred_at"`
}

// NewUserUpdatedEvent creates a new user updated event
func NewUserUpdatedEvent(userID, email, fullName string, version int) *UserUpdatedEvent {
	return &UserUpdatedEvent{
		ID:           uuid.New().String(),
		UserID:       userID,
		Email:        email,
		FullName:     fullName,
		AggregateID_: userID,
		Version_:     version,
		OccurredAt:   time.Now(),
	}
}

// EventType returns the event type
func (e *UserUpdatedEvent) EventType() string {
	return "user.updated"
}

// EventData returns the event data as JSON bytes
func (e *UserUpdatedEvent) EventData() ([]byte, error) {
	return json.Marshal(e)
}

// AggregateID returns the aggregate ID
func (e *UserUpdatedEvent) AggregateID() string {
	return e.AggregateID_
}

// Version returns the event version
func (e *UserUpdatedEvent) Version() int {
	return e.Version_
}

// Timestamp returns the event timestamp
func (e *UserUpdatedEvent) Timestamp() int64 {
	return e.OccurredAt.Unix()
}

// UserAvatarUploadedEvent represents a user avatar upload domain event
type UserAvatarUploadedEvent struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	AvatarURL    string    `json:"avatar_url"`
	FileKey      string    `json:"file_key"`
	AggregateID_ string    `json:"aggregate_id"`
	Version_     int       `json:"version"`
	OccurredAt   time.Time `json:"occurred_at"`
}

// NewUserAvatarUploadedEvent creates a new user avatar uploaded event
func NewUserAvatarUploadedEvent(userID, avatarURL, fileKey string, version int) *UserAvatarUploadedEvent {
	return &UserAvatarUploadedEvent{
		ID:           uuid.New().String(),
		UserID:       userID,
		AvatarURL:    avatarURL,
		FileKey:      fileKey,
		AggregateID_: userID,
		Version_:     version,
		OccurredAt:   time.Now(),
	}
}

// EventType returns the event type
func (e *UserAvatarUploadedEvent) EventType() string {
	return "user.avatar.uploaded"
}

// EventData returns the event data as JSON bytes
func (e *UserAvatarUploadedEvent) EventData() ([]byte, error) {
	return json.Marshal(e)
}

// AggregateID returns the aggregate ID
func (e *UserAvatarUploadedEvent) AggregateID() string {
	return e.AggregateID_
}

// Version returns the event version
func (e *UserAvatarUploadedEvent) Version() int {
	return e.Version_
}

// Timestamp returns the event timestamp
func (e *UserAvatarUploadedEvent) Timestamp() int64 {
	return e.OccurredAt.Unix()
}
