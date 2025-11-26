package messaging

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/messaging"
	"gorm.io/gorm"
)

type gormOutboxRepository struct {
	db *gorm.DB
}

func NewOutboxRepository(db *gorm.DB) messaging.OutboxRepository {
	return &gormOutboxRepository{
		db: db,
	}
}

func (r *gormOutboxRepository) Create(ctx context.Context, message *messaging.OutboxMessage) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *gormOutboxRepository) CreateWithTx(ctx context.Context, tx interface{}, message *messaging.OutboxMessage) error {
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return r.Create(ctx, message)
	}
	return gormTx.WithContext(ctx).Create(message).Error
}

func (r *gormOutboxRepository) GetPendingMessages(ctx context.Context, limit int) ([]*messaging.OutboxMessage, error) {
	var messages []*messaging.OutboxMessage
	err := r.db.WithContext(ctx).
		Where("status IN ?", []messaging.OutboxMessageStatus{
			messaging.OutboxMessageStatusPending,
			messaging.OutboxMessageStatusFailed,
		}).
		Order("created_at ASC").
		Limit(limit).
		Find(&messages).Error

	return messages, err
}

func (r *gormOutboxRepository) GetByID(ctx context.Context, id uuid.UUID) (*messaging.OutboxMessage, error) {
	var message messaging.OutboxMessage
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&message).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &message, nil
}

func (r *gormOutboxRepository) Update(ctx context.Context, message *messaging.OutboxMessage) error {
	return r.db.WithContext(ctx).Save(message).Error
}

func (r *gormOutboxRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status messaging.OutboxMessageStatus) error {
	return r.db.WithContext(ctx).
		Model(&messaging.OutboxMessage{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

func (r *gormOutboxRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&messaging.OutboxMessage{}, "id = ?", id).Error
}

func (r *gormOutboxRepository) DeleteOldProcessedMessages(ctx context.Context, olderThan int64) error {
	return r.db.WithContext(ctx).
		Delete(&messaging.OutboxMessage{},
			"status = ? AND processed_at < ?",
			messaging.OutboxMessageStatusProcessed,
			time.Unix(olderThan, 0)).Error
}
