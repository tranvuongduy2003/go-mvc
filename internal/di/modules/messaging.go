package modules

import (
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"

	messagingServices "github.com/tranvuongduy2003/go-mvc/internal/application/services/messaging"
	messagingPorts "github.com/tranvuongduy2003/go-mvc/internal/domain/ports/messaging"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/config"
	jobHandlers "github.com/tranvuongduy2003/go-mvc/internal/infrastructure/jobs/handlers"
	natsAdapter "github.com/tranvuongduy2003/go-mvc/internal/infrastructure/messaging/nats"
	postgresMessaging "github.com/tranvuongduy2003/go-mvc/internal/infrastructure/persistence/postgres/messaging"
	"github.com/tranvuongduy2003/go-mvc/internal/interfaces/http/middleware"
)

// MessagingModule provides messaging-related dependencies
var MessagingModule = fx.Module("messaging",
	fx.Provide(
		// Repositories
		NewOutboxRepository,
		NewInboxRepository,
		NewMessageDeduplicationRepository,

		// Services
		NewOutboxService,
		NewInboxService,

		// Enhanced NATS Adapter
		NewDeduplicatedNATSBroker,

		// Background Jobs
		NewOutboxProcessorJob,

		// Middleware
		NewIdempotencyMiddleware,
	),
)

// NewOutboxRepository creates outbox repository
func NewOutboxRepository(db *gorm.DB) messagingPorts.OutboxRepository {
	return postgresMessaging.NewOutboxRepository(db)
}

// NewInboxRepository creates inbox repository
func NewInboxRepository(db *gorm.DB) messagingPorts.InboxRepository {
	return postgresMessaging.NewInboxRepository(db)
}

// NewMessageDeduplicationRepository creates message deduplication repository
func NewMessageDeduplicationRepository(db *gorm.DB) messagingPorts.MessageDeduplicationRepository {
	return postgresMessaging.NewMessageDeduplicationRepository(db)
}

// NewOutboxService creates outbox service
func NewOutboxService(outboxRepo messagingPorts.OutboxRepository) *messagingServices.OutboxService {
	return messagingServices.NewOutboxService(outboxRepo)
}

// NewInboxService creates inbox service
func NewInboxService(
	inboxRepo messagingPorts.InboxRepository,
	dedupRepo messagingPorts.MessageDeduplicationRepository,
) *messagingServices.InboxService {
	return messagingServices.NewInboxService(inboxRepo, dedupRepo)
}

// NewDeduplicatedNATSBroker creates enhanced NATS broker with deduplication
func NewDeduplicatedNATSBroker(
	natsBroker *natsAdapter.NATSBroker,
	inboxService *messagingServices.InboxService,
	cfg *config.AppConfig,
) *natsAdapter.DeduplicatedNATSBroker {
	consumerID := "go-mvc-service" // Or get from config
	return natsAdapter.NewDeduplicatedNATSBroker(natsBroker, inboxService, consumerID)
}

// NewOutboxProcessorJob creates outbox processor background job
func NewOutboxProcessorJob(
	outboxService *messagingServices.OutboxService,
	publisher messagingPorts.Publisher,
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

// NewIdempotencyMiddleware creates idempotency middleware
func NewIdempotencyMiddleware(
	inboxService *messagingServices.InboxService,
	logger *zap.Logger,
	cfg *config.AppConfig,
) *middleware.IdempotencyMiddleware {
	ttl := 24 * time.Hour // Could be configurable
	return middleware.NewIdempotencyMiddleware(inboxService, logger, ttl)
}
