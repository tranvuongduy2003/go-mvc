package messaging

import (
	"context"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/messaging"
)

// OutboxRepository defines the interface for outbox message persistence
type OutboxRepository interface {
	// Create stores a new outbox message
	Create(ctx context.Context, message *messaging.OutboxMessage) error

	// CreateWithTx stores a new outbox message within a transaction
	CreateWithTx(ctx context.Context, tx interface{}, message *messaging.OutboxMessage) error

	// GetPendingMessages retrieves pending outbox messages for processing
	GetPendingMessages(ctx context.Context, limit int) ([]*messaging.OutboxMessage, error)

	// GetByID retrieves an outbox message by ID
	GetByID(ctx context.Context, id uuid.UUID) (*messaging.OutboxMessage, error)

	// Update updates an existing outbox message
	Update(ctx context.Context, message *messaging.OutboxMessage) error

	// UpdateStatus updates the status of an outbox message
	UpdateStatus(ctx context.Context, id uuid.UUID, status messaging.OutboxMessageStatus) error

	// Delete removes an outbox message (typically after successful processing)
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteOldProcessedMessages removes old processed messages for cleanup
	DeleteOldProcessedMessages(ctx context.Context, olderThan int64) error
}
