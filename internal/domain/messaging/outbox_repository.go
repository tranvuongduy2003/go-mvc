package messaging

import (
	"context"

	"github.com/google/uuid"
)

type OutboxRepository interface {
	Create(ctx context.Context, message *OutboxMessage) error

	CreateWithTx(ctx context.Context, tx interface{}, message *OutboxMessage) error

	GetPendingMessages(ctx context.Context, limit int) ([]*OutboxMessage, error)

	GetByID(ctx context.Context, id uuid.UUID) (*OutboxMessage, error)

	Update(ctx context.Context, message *OutboxMessage) error

	UpdateStatus(ctx context.Context, id uuid.UUID, status OutboxMessageStatus) error

	Delete(ctx context.Context, id uuid.UUID) error

	DeleteOldProcessedMessages(ctx context.Context, olderThan int64) error
}
