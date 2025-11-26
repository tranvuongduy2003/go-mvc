package messaging

import (
	"context"

	"github.com/google/uuid"
)

// OutboxRepository defines the interface for outbox message persistence
type OutboxRepository interface {
	// Create stores a new outbox message
	Create(ctx context.Context, message *OutboxMessage) error

	// CreateWithTx stores a new outbox message within a transaction
	CreateWithTx(ctx context.Context, tx interface{}, message *OutboxMessage) error

	// GetPendingMessages retrieves pending outbox messages for processing
	GetPendingMessages(ctx context.Context, limit int) ([]*OutboxMessage, error)

	// GetByID retrieves an outbox message by ID
	GetByID(ctx context.Context, id uuid.UUID) (*OutboxMessage, error)

	// Update updates an existing outbox message
	Update(ctx context.Context, message *OutboxMessage) error

	// UpdateStatus updates the status of an outbox message
	UpdateStatus(ctx context.Context, id uuid.UUID, status OutboxMessageStatus) error

	// Delete removes an outbox message (typically after successful processing)
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteOldProcessedMessages removes old processed messages for cleanup
	DeleteOldProcessedMessages(ctx context.Context, olderThan int64) error
}
