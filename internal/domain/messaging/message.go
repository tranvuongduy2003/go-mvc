package messaging

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

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

func (m *Message) WithCorrelationID(correlationID string) *Message {
	m.CorrelationID = &correlationID
	return m
}

func (m *Message) WithCausationID(causationID string) *Message {
	m.CausationID = &causationID
	return m
}

func (m *Message) WithMetadata(key string, value interface{}) *Message {
	if m.Metadata == nil {
		m.Metadata = make(map[string]interface{})
	}
	m.Metadata[key] = value
	return m
}

func (m *Message) GetMetadata(key string) (interface{}, bool) {
	if m.Metadata == nil {
		return nil, false
	}
	value, exists := m.Metadata[key]
	return value, exists
}
