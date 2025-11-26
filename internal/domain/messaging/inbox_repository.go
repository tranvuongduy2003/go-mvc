package messaging

import (
	"context"

	"github.com/google/uuid"
)

type InboxRepository interface {
	Create(ctx context.Context, message *InboxMessage) error

	CreateWithTx(ctx context.Context, tx interface{}, message *InboxMessage) error

	GetByMessageID(ctx context.Context, messageID uuid.UUID, consumerID string) (*InboxMessage, error)

	Exists(ctx context.Context, messageID uuid.UUID, consumerID string) (bool, error)

	Update(ctx context.Context, message *InboxMessage) error

	MarkAsProcessed(ctx context.Context, messageID uuid.UUID, consumerID string) error

	Delete(ctx context.Context, id uuid.UUID) error

	DeleteOldMessages(ctx context.Context, olderThan int64) error
}
