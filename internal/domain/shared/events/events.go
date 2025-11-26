package events

import (
	"time"

	"github.com/google/uuid"
)

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

type BaseDomainEvent struct {
	ID            uuid.UUID   `json:"id"`
	Type          string      `json:"type"`
	AggregateID   uuid.UUID   `json:"aggregate_id"`
	AggregateType string      `json:"aggregate_type"`
	Version       int         `json:"version"`
	OccurredAt    time.Time   `json:"occurred_at"`
	Data          interface{} `json:"data"`
}

func (e *BaseDomainEvent) GetID() uuid.UUID {
	return e.ID
}

func (e *BaseDomainEvent) GetType() string {
	return e.Type
}

func (e *BaseDomainEvent) GetAggregateID() uuid.UUID {
	return e.AggregateID
}

func (e *BaseDomainEvent) GetAggregateType() string {
	return e.AggregateType
}

func (e *BaseDomainEvent) GetVersion() int {
	return e.Version
}

func (e *BaseDomainEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

func (e *BaseDomainEvent) GetData() interface{} {
	return e.Data
}

func (e *BaseDomainEvent) SetVersion(version int) {
	e.Version = version
}

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
