package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type UserCreatedEvent struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	AggregateID_ string    `json:"aggregate_id"`
	Version_     int       `json:"version"`
	OccurredAt   time.Time `json:"occurred_at"`
}

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

func (e *UserCreatedEvent) EventType() string {
	return "user.created"
}

func (e *UserCreatedEvent) EventData() ([]byte, error) {
	return json.Marshal(e)
}

func (e *UserCreatedEvent) AggregateID() string {
	return e.AggregateID_
}

func (e *UserCreatedEvent) Version() int {
	return e.Version_
}

func (e *UserCreatedEvent) Timestamp() int64 {
	return e.OccurredAt.Unix()
}

type UserUpdatedEvent struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	AggregateID_ string    `json:"aggregate_id"`
	Version_     int       `json:"version"`
	OccurredAt   time.Time `json:"occurred_at"`
}

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

func (e *UserUpdatedEvent) EventType() string {
	return "user.updated"
}

func (e *UserUpdatedEvent) EventData() ([]byte, error) {
	return json.Marshal(e)
}

func (e *UserUpdatedEvent) AggregateID() string {
	return e.AggregateID_
}

func (e *UserUpdatedEvent) Version() int {
	return e.Version_
}

func (e *UserUpdatedEvent) Timestamp() int64 {
	return e.OccurredAt.Unix()
}

type UserAvatarUploadedEvent struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	AvatarURL    string    `json:"avatar_url"`
	FileKey      string    `json:"file_key"`
	AggregateID_ string    `json:"aggregate_id"`
	Version_     int       `json:"version"`
	OccurredAt   time.Time `json:"occurred_at"`
}

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

func (e *UserAvatarUploadedEvent) EventType() string {
	return "user.avatar.uploaded"
}

func (e *UserAvatarUploadedEvent) EventData() ([]byte, error) {
	return json.Marshal(e)
}

func (e *UserAvatarUploadedEvent) AggregateID() string {
	return e.AggregateID_
}

func (e *UserAvatarUploadedEvent) Version() int {
	return e.Version_
}

func (e *UserAvatarUploadedEvent) Timestamp() int64 {
	return e.OccurredAt.Unix()
}
