package messaging

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MessageDeduplicationRepository interface {
	Create(ctx context.Context, dedup *MessageDeduplication) error

	CreateWithTx(ctx context.Context, tx interface{}, dedup *MessageDeduplication) error

	Exists(ctx context.Context, messageID uuid.UUID, consumerID string) (bool, error)

	GetByMessageID(ctx context.Context, messageID uuid.UUID, consumerID string) (*MessageDeduplication, error)

	Delete(ctx context.Context, messageID uuid.UUID, consumerID string) error

	DeleteExpired(ctx context.Context) error

	CreateIfNotExists(ctx context.Context, messageID uuid.UUID, consumerID, eventType string, ttl time.Duration) (bool, error)
}
