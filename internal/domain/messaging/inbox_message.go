package messaging

import (
	"time"

	"github.com/google/uuid"
)

type InboxMessageStatus string

const (
	InboxMessageStatusReceived  InboxMessageStatus = "received"
	InboxMessageStatusProcessed InboxMessageStatus = "processed"
	InboxMessageStatusIgnored   InboxMessageStatus = "ignored"
)

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

func (m *InboxMessage) MarkAsProcessed() {
	now := time.Now()
	m.Status = InboxMessageStatusProcessed
	m.ProcessedAt = &now
	m.UpdatedAt = now
}

func (m *InboxMessage) MarkAsIgnored() {
	m.Status = InboxMessageStatusIgnored
	m.UpdatedAt = time.Now()
}

func (m *InboxMessage) IsProcessed() bool {
	return m.Status == InboxMessageStatusProcessed
}

func (m *InboxMessage) IsDuplicate() bool {
	return m.Status != InboxMessageStatusReceived
}
