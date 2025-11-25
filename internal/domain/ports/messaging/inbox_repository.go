package messaging

import (
	"context"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/messaging"
)

// InboxRepository defines the interface for inbox message persistence
type InboxRepository interface {
	// Create stores a new inbox message
	Create(ctx context.Context, message *messaging.InboxMessage) error

	// CreateWithTx stores a new inbox message within a transaction
	CreateWithTx(ctx context.Context, tx interface{}, message *messaging.InboxMessage) error

	// GetByMessageID retrieves an inbox message by message ID and consumer ID
	GetByMessageID(ctx context.Context, messageID uuid.UUID, consumerID string) (*messaging.InboxMessage, error)

	// Exists checks if a message ID has already been processed by a consumer
	Exists(ctx context.Context, messageID uuid.UUID, consumerID string) (bool, error)

	// Update updates an existing inbox message
	Update(ctx context.Context, message *messaging.InboxMessage) error

	// MarkAsProcessed marks a message as processed
	MarkAsProcessed(ctx context.Context, messageID uuid.UUID, consumerID string) error

	// Delete removes an inbox message
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteOldMessages removes old inbox messages for cleanup
	DeleteOldMessages(ctx context.Context, olderThan int64) error
}
