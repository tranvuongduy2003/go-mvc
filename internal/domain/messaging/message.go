package messaging

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Message represents a generic message with deduplication support
type Message struct {
	ID            uuid.UUID              `json:"id"`
	EventType     string                 `json:"event_type"`
	AggregateID   string                 `json:"aggregate_id"`
	Payload       json.RawMessage        `json:"payload"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
	CorrelationID *string                `json:"correlation_id,omitempty"`
	CausationID   *string                `json:"causation_id,omitempty"`
}

// NewMessage creates a new message with a unique ID
func NewMessage(eventType, aggregateID string, payload interface{}) (*Message, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Message{
		ID:          uuid.New(),
		EventType:   eventType,
		AggregateID: aggregateID,
		Payload:     payloadBytes,
		Metadata:    make(map[string]interface{}),
		Timestamp:   time.Now(),
	}, nil
}

// WithCorrelationID sets the correlation ID for message tracing
func (m *Message) WithCorrelationID(correlationID string) *Message {
	m.CorrelationID = &correlationID
	return m
}

// WithCausationID sets the causation ID for message tracing
func (m *Message) WithCausationID(causationID string) *Message {
	m.CausationID = &causationID
	return m
}

// WithMetadata adds metadata to the message
func (m *Message) WithMetadata(key string, value interface{}) *Message {
	if m.Metadata == nil {
		m.Metadata = make(map[string]interface{})
	}
	m.Metadata[key] = value
	return m
}

// GetMetadata retrieves metadata value by key
func (m *Message) GetMetadata(key string) (interface{}, bool) {
	if m.Metadata == nil {
		return nil, false
	}
	value, exists := m.Metadata[key]
	return value, exists
}
