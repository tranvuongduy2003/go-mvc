package messaging

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type OutboxMessageStatus string

const (
	OutboxMessageStatusPending   OutboxMessageStatus = "pending"
	OutboxMessageStatusProcessed OutboxMessageStatus = "processed"
	OutboxMessageStatusFailed    OutboxMessageStatus = "failed"
)

type OutboxMessage struct {
	ID           uuid.UUID           `json:"id" db:"id"`
	MessageID    uuid.UUID           `json:"message_id" db:"message_id"`
	EventType    string              `json:"event_type" db:"event_type"`
	AggregateID  string              `json:"aggregate_id" db:"aggregate_id"`
	Payload      json.RawMessage     `json:"payload" db:"payload"`
	Status       OutboxMessageStatus `json:"status" db:"status"`
	Retries      int                 `json:"retries" db:"retries"`
	MaxRetries   int                 `json:"max_retries" db:"max_retries"`
	CreatedAt    time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at" db:"updated_at"`
	ProcessedAt  *time.Time          `json:"processed_at,omitempty" db:"processed_at"`
	FailedAt     *time.Time          `json:"failed_at,omitempty" db:"failed_at"`
	ErrorMessage *string             `json:"error_message,omitempty" db:"error_message"`
}

func NewOutboxMessage(eventType, aggregateID string, payload interface{}) (*OutboxMessage, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &OutboxMessage{
		ID:          uuid.New(),
		MessageID:   uuid.New(),
		EventType:   eventType,
		AggregateID: aggregateID,
		Payload:     payloadBytes,
		Status:      OutboxMessageStatusPending,
		Retries:     0,
		MaxRetries:  3,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (m *OutboxMessage) MarkAsProcessed() {
	now := time.Now()
	m.Status = OutboxMessageStatusProcessed
	m.ProcessedAt = &now
	m.UpdatedAt = now
}

func (m *OutboxMessage) MarkAsFailed(errorMsg string) {
	now := time.Now()
	m.Status = OutboxMessageStatusFailed
	m.ErrorMessage = &errorMsg
	m.FailedAt = &now
	m.UpdatedAt = now
}

func (m *OutboxMessage) IncrementRetry() {
	m.Retries++
	m.UpdatedAt = time.Now()
}

func (m *OutboxMessage) CanRetry() bool {
	return m.Retries < m.MaxRetries
}

func (m *OutboxMessage) ShouldRetry() bool {
	return m.Status == OutboxMessageStatusFailed && m.CanRetry()
}
