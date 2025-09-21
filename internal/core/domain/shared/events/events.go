package events

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent represents a domain event
type DomainEvent interface {
	GetID() uuid.UUID
	GetType() string
	GetAggregateID() uuid.UUID
	GetAggregateType() string
	GetVersion() int
	GetOccurredAt() time.Time
	GetData() interface{}
	SetVersion(version int)
}

// BaseDomainEvent provides common fields for all domain events
type BaseDomainEvent struct {
	ID            uuid.UUID   `json:"id"`
	Type          string      `json:"type"`
	AggregateID   uuid.UUID   `json:"aggregate_id"`
	AggregateType string      `json:"aggregate_type"`
	Version       int         `json:"version"`
	OccurredAt    time.Time   `json:"occurred_at"`
	Data          interface{} `json:"data"`
}

// GetID returns event ID
func (e *BaseDomainEvent) GetID() uuid.UUID {
	return e.ID
}

// GetType returns event type
func (e *BaseDomainEvent) GetType() string {
	return e.Type
}

// GetAggregateID returns aggregate ID
func (e *BaseDomainEvent) GetAggregateID() uuid.UUID {
	return e.AggregateID
}

// GetAggregateType returns aggregate type
func (e *BaseDomainEvent) GetAggregateType() string {
	return e.AggregateType
}

// GetVersion returns event version
func (e *BaseDomainEvent) GetVersion() int {
	return e.Version
}

// GetOccurredAt returns when event occurred
func (e *BaseDomainEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetData returns event data
func (e *BaseDomainEvent) GetData() interface{} {
	return e.Data
}

// SetVersion sets event version
func (e *BaseDomainEvent) SetVersion(version int) {
	e.Version = version
}

// NewBaseDomainEvent creates a new base domain event
func NewBaseDomainEvent(eventType string, aggregateID uuid.UUID, aggregateType string, data interface{}) *BaseDomainEvent {
	return &BaseDomainEvent{
		ID:            uuid.New(),
		Type:          eventType,
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		Version:       1,
		OccurredAt:    time.Now().UTC(),
		Data:          data,
	}
}
