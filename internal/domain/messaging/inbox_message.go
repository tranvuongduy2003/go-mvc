package messaging

import (
	"time"

	"github.com/google/uuid"
)

// InboxMessageStatus represents the status of an inbox message
type InboxMessageStatus string

const (
	InboxMessageStatusReceived  InboxMessageStatus = "received"
	InboxMessageStatusProcessed InboxMessageStatus = "processed"
	InboxMessageStatusIgnored   InboxMessageStatus = "ignored"
)

// InboxMessage represents a message stored in the inbox pattern
// This is used for deduplication on the consumer side
type InboxMessage struct {
	ID          uuid.UUID          `json:"id" db:"id"`
	MessageID   uuid.UUID          `json:"message_id" db:"message_id"`
	EventType   string             `json:"event_type" db:"event_type"`
	ConsumerID  string             `json:"consumer_id" db:"consumer_id"`
	Status      InboxMessageStatus `json:"status" db:"status"`
	ReceivedAt  time.Time          `json:"received_at" db:"received_at"`
	ProcessedAt *time.Time         `json:"processed_at,omitempty" db:"processed_at"`
	CreatedAt   time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" db:"updated_at"`
}

// NewInboxMessage creates a new inbox message
func NewInboxMessage(messageID uuid.UUID, eventType, consumerID string) *InboxMessage {
	now := time.Now()
	return &InboxMessage{
		ID:         uuid.New(),
		MessageID:  messageID,
		EventType:  eventType,
		ConsumerID: consumerID,
		Status:     InboxMessageStatusReceived,
		ReceivedAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// MarkAsProcessed marks the inbox message as successfully processed
func (m *InboxMessage) MarkAsProcessed() {
	now := time.Now()
	m.Status = InboxMessageStatusProcessed
	m.ProcessedAt = &now
	m.UpdatedAt = now
}

// MarkAsIgnored marks the inbox message as ignored (duplicate)
func (m *InboxMessage) MarkAsIgnored() {
	m.Status = InboxMessageStatusIgnored
	m.UpdatedAt = time.Now()
}

// IsProcessed checks if the message has been processed
func (m *InboxMessage) IsProcessed() bool {
	return m.Status == InboxMessageStatusProcessed
}

// IsDuplicate checks if this message is a duplicate (already received)
func (m *InboxMessage) IsDuplicate() bool {
	return m.Status != InboxMessageStatusReceived
}
