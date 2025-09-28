package messaging

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/messaging"
	messagingPorts "github.com/tranvuongduy2003/go-mvc/internal/core/ports/messaging"
	"gorm.io/gorm"
)

type gormInboxRepository struct {
	db *gorm.DB
}

// NewInboxRepository creates a new inbox repository using GORM
func NewInboxRepository(db *gorm.DB) messagingPorts.InboxRepository {
	return &gormInboxRepository{
		db: db,
	}
}

func (r *gormInboxRepository) Create(ctx context.Context, message *messaging.InboxMessage) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *gormInboxRepository) CreateWithTx(ctx context.Context, tx interface{}, message *messaging.InboxMessage) error {
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return r.Create(ctx, message)
	}
	return gormTx.WithContext(ctx).Create(message).Error
}

func (r *gormInboxRepository) GetByMessageID(ctx context.Context, messageID uuid.UUID, consumerID string) (*messaging.InboxMessage, error) {
	var message messaging.InboxMessage
	err := r.db.WithContext(ctx).
		Where("message_id = ? AND consumer_id = ?", messageID, consumerID).
		First(&message).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &message, nil
}

func (r *gormInboxRepository) Exists(ctx context.Context, messageID uuid.UUID, consumerID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&messaging.InboxMessage{}).
		Where("message_id = ? AND consumer_id = ?", messageID, consumerID).
		Count(&count).Error

	return count > 0, err
}

func (r *gormInboxRepository) Update(ctx context.Context, message *messaging.InboxMessage) error {
	return r.db.WithContext(ctx).Save(message).Error
}

func (r *gormInboxRepository) MarkAsProcessed(ctx context.Context, messageID uuid.UUID, consumerID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&messaging.InboxMessage{}).
		Where("message_id = ? AND consumer_id = ?", messageID, consumerID).
		Updates(map[string]interface{}{
			"status":       messaging.InboxMessageStatusProcessed,
			"processed_at": &now,
			"updated_at":   now,
		}).Error
}

func (r *gormInboxRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&messaging.InboxMessage{}, "id = ?", id).Error
}

func (r *gormInboxRepository) DeleteOldMessages(ctx context.Context, olderThan int64) error {
	return r.db.WithContext(ctx).
		Delete(&messaging.InboxMessage{}, "created_at < ?", time.Unix(olderThan, 0)).Error
}
