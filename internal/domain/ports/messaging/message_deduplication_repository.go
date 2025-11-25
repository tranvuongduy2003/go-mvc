package messaging

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/messaging"
)

// MessageDeduplicationRepository defines the interface for message deduplication
type MessageDeduplicationRepository interface {
	// Create stores a new message deduplication record
	Create(ctx context.Context, dedup *messaging.MessageDeduplication) error

	// CreateWithTx stores a new message deduplication record within a transaction
	CreateWithTx(ctx context.Context, tx interface{}, dedup *messaging.MessageDeduplication) error

	// Exists checks if a message has already been processed
	Exists(ctx context.Context, messageID uuid.UUID, consumerID string) (bool, error)

	// GetByMessageID retrieves a deduplication record by message ID and consumer ID
	GetByMessageID(ctx context.Context, messageID uuid.UUID, consumerID string) (*messaging.MessageDeduplication, error)

	// Delete removes a deduplication record
	Delete(ctx context.Context, messageID uuid.UUID, consumerID string) error

	// DeleteExpired removes expired deduplication records
	DeleteExpired(ctx context.Context) error

	// CreateIfNotExists creates a deduplication record if it doesn't exist
	// Returns true if created, false if already exists
	CreateIfNotExists(ctx context.Context, messageID uuid.UUID, consumerID, eventType string, ttl time.Duration) (bool, error)
}
