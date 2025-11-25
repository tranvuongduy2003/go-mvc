package messaging

import (
	"time"

	"github.com/google/uuid"
)

// MessageDeduplication represents a record for message deduplication
// This provides a lightweight alternative to the full inbox pattern for simple deduplication
type MessageDeduplication struct {
	ID          uuid.UUID `json:"id" db:"id"`
	MessageID   uuid.UUID `json:"message_id" db:"message_id"`
	ConsumerID  string    `json:"consumer_id" db:"consumer_id"`
	EventType   string    `json:"event_type" db:"event_type"`
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`
}

// NewMessageDeduplication creates a new message deduplication record
func NewMessageDeduplication(messageID uuid.UUID, consumerID, eventType string, ttl time.Duration) *MessageDeduplication {
	now := time.Now()
	return &MessageDeduplication{
		ID:          uuid.New(),
		MessageID:   messageID,
		ConsumerID:  consumerID,
		EventType:   eventType,
		ProcessedAt: now,
		CreatedAt:   now,
		ExpiresAt:   now.Add(ttl),
	}
}

// IsExpired checks if the deduplication record has expired
func (m *MessageDeduplication) IsExpired() bool {
	return time.Now().After(m.ExpiresAt)
}
