package modules

import (
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"

	messagingServices "github.com/tranvuongduy2003/go-mvc/internal/application/services/messaging"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/messaging"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/config"
	jobHandlers "github.com/tranvuongduy2003/go-mvc/internal/infrastructure/jobs/handlers"
	natsAdapter "github.com/tranvuongduy2003/go-mvc/internal/infrastructure/messaging/nats"
	postgresMessaging "github.com/tranvuongduy2003/go-mvc/internal/infrastructure/persistence/postgres/messaging"
	"github.com/tranvuongduy2003/go-mvc/internal/presentation/http/middleware"
)

var MessagingModule = fx.Module("messaging",
	fx.Provide(
		NewOutboxRepository,
		NewInboxRepository,
		NewMessageDeduplicationRepository,

		NewOutboxService,
		NewInboxService,

		NewDeduplicatedNATSBroker,

		NewOutboxProcessorJob,

		NewIdempotencyMiddleware,
	),
)

func NewOutboxRepository(db *gorm.DB) messaging.OutboxRepository {
	return postgresMessaging.NewOutboxRepository(db)
}

func NewInboxRepository(db *gorm.DB) messaging.InboxRepository {
	return postgresMessaging.NewInboxRepository(db)
}

func NewMessageDeduplicationRepository(db *gorm.DB) messaging.MessageDeduplicationRepository {
	return postgresMessaging.NewMessageDeduplicationRepository(db)
}

func NewOutboxService(outboxRepo messaging.OutboxRepository) *messagingServices.OutboxService {
	return messagingServices.NewOutboxService(outboxRepo)
}

func NewInboxService(
	inboxRepo messaging.InboxRepository,
	dedupRepo messaging.MessageDeduplicationRepository,
) *messagingServices.InboxService {
	return messagingServices.NewInboxService(inboxRepo, dedupRepo)
}

func NewDeduplicatedNATSBroker(
	natsBroker *natsAdapter.NATSBroker,
	inboxService *messagingServices.InboxService,
	cfg *config.AppConfig,
) *natsAdapter.DeduplicatedNATSBroker {
	consumerID := "go-mvc-service" // Or get from config
	return natsAdapter.NewDeduplicatedNATSBroker(natsBroker, inboxService, consumerID)
}

func NewOutboxProcessorJob(
	outboxService *messagingServices.OutboxService,
	publisher messaging.Publisher,
	cfg *config.AppConfig,
) *jobHandlers.OutboxProcessorJob {
	batchSize := 10               // Could be configurable
	retryDelay := 5 * time.Second // Could be configurable

	return jobHandlers.NewOutboxProcessorJob(
		outboxService,
		publisher,
		batchSize,
		retryDelay,
	)
}

func NewIdempotencyMiddleware(
	inboxService *messagingServices.InboxService,
	logger *zap.Logger,
	cfg *config.AppConfig,
) *middleware.IdempotencyMiddleware {
	ttl := 24 * time.Hour // Could be configurable
	return middleware.NewIdempotencyMiddleware(inboxService, logger, ttl)
}
