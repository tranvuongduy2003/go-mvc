# NATS Messaging Implementation

This document explains the NATS messaging implementation for handling domain events and asynchronous communication in the Go MVC application.

## Overview

The application uses NATS as a message broker for:
- Domain event publishing and handling
- Asynchronous communication between services
- Decoupling business logic from side effects
- Event-driven architecture patterns

## Architecture

### Core Components

1. **Messaging Ports** (`internal/core/ports/messaging/`)
   - Defines interfaces for messaging operations
   - MessageBroker, EventBus, Publisher, Subscriber interfaces
   - Event and Message abstractions

2. **NATS Adapter** (`internal/adapters/messaging/nats/`)
   - Implements messaging interfaces using NATS
   - Connection management with auto-reconnection
   - Publisher/Subscriber functionality
   - Event bus implementation

3. **Domain Events** (`internal/core/domain/shared/events/`)
   - Defines domain events (UserCreated, UserUpdated, etc.)
   - Event serialization and metadata

4. **Event Handlers** (`internal/application/events/handlers/`)
   - Handles domain events asynchronously
   - Business logic for event processing

## Configuration

### NATS Configuration (configs/development.yaml)

```yaml
messaging:
  nats:
    url: "nats://localhost:4222"
    max_reconnects: 10
    reconnect_wait: "2s"
    timeout: "5s"
    drain_timeout: "30s"
```

### Docker Compose NATS Service

```yaml
nats:
  image: nats:2.10-alpine
  container_name: dev-nats
  ports:
    - "4222:4222"  # Client connections
    - "8222:8222"  # HTTP monitoring
  command: ["-js", "-m", "8222"]
```

## Usage Examples

### Publishing Events

```go
// In command handlers
func (h *UploadAvatarCommandHandler) Handle(ctx context.Context, cmd UploadAvatarCommand) error {
    // Business logic...
    
    // Publish event
    event := events.NewUserAvatarUploadedEvent(userID, avatarURL, fileKey, version)
    if err := h.eventBus.PublishEvent(ctx, event); err != nil {
        // Handle error (log, retry, etc.)
    }
    
    return nil
}
```

### Subscribing to Events

```go
// In event handlers
func (h *UserEventHandler) HandleUserCreated(ctx context.Context, event messaging.Event) error {
    // Parse event data
    var userCreated events.UserCreatedEvent
    data, _ := event.EventData()
    json.Unmarshal(data, &userCreated)
    
    // Handle the event
    // - Send welcome email
    // - Create user profile
    // - Update analytics
    
    return nil
}

// Setup subscriptions
func (h *UserEventHandler) SetupEventSubscriptions(eventBus messaging.EventBus) error {
    _, err := eventBus.SubscribeToEvent("user.created", h.HandleUserCreated)
    return err
}
```

### Direct Messaging (Pub/Sub)

```go
// Publishing messages
func (s *Service) PublishMessage(ctx context.Context, subject string, data []byte) error {
    return s.messageBroker.Publish(ctx, subject, data)
}

// Subscribing to messages
func (s *Service) SubscribeToMessages(subject string, handler messaging.MessageHandler) error {
    subscription, err := s.messageBroker.Subscribe(subject, handler)
    if err != nil {
        return err
    }
    
    // Store subscription for cleanup
    s.subscriptions = append(s.subscriptions, subscription)
    return nil
}
```

## Event Types

### User Events

1. **user.created**
   - Published when a new user is registered
   - Contains user ID, email, full name
   - Triggers welcome email, profile creation

2. **user.updated**
   - Published when user information is updated
   - Contains updated user details
   - Triggers cache invalidation, search index updates

3. **user.avatar.uploaded**
   - Published when user uploads avatar
   - Contains user ID, avatar URL, file key
   - Triggers thumbnail generation, cache updates

## Event Flow

```
Command Handler -> Domain Event -> Event Bus -> NATS -> Event Handler -> Side Effects
     ↓               ↓              ↓          ↓         ↓              ↓
Update User    -> UserUpdated -> PublishEvent -> NATS -> HandleEvent -> Send Email
                                                                     -> Update Cache
                                                                     -> Update Search
```

## Error Handling

### Event Publishing Failures
- Events are published after successful business operations
- Publishing failures are logged but don't fail the operation
- Consider implementing event store for guaranteed delivery

### Event Processing Failures
- Failed event handlers are logged with context
- Implement retry mechanisms for transient failures
- Consider dead letter queues for persistent failures

### Connection Failures
- NATS client automatically reconnects
- Configurable reconnection attempts and delays
- Connection health checks and monitoring

## Monitoring

### NATS Monitoring
- HTTP monitoring interface at http://localhost:8222
- Connection statistics and message metrics
- Queue depths and consumer information

### Application Metrics
- Event publishing success/failure rates
- Event processing latencies
- Connection status and health

### Logging
- Structured logging for all messaging operations
- Event publishing and consumption logs
- Error tracking with context

## Best Practices

### Event Design
- Make events immutable and self-contained
- Include all necessary data in event payload
- Use semantic versioning for event schemas
- Keep events small and focused

### Error Handling
- Always log event processing failures
- Implement idempotent event handlers
- Use dead letter queues for poison messages
- Monitor event processing metrics

### Performance
- Use queue subscriptions for load balancing
- Batch related events when possible
- Configure appropriate timeouts
- Monitor memory usage and connection pools

### Security
- Use NATS authentication in production
- Encrypt sensitive event payloads
- Implement access controls for topics
- Monitor for suspicious activity

## Production Considerations

### Scalability
- Use NATS clustering for high availability
- Implement proper queue groups for load distribution
- Monitor resource usage and scaling metrics
- Consider NATS JetStream for persistence

### Reliability
- Implement event sourcing for critical events
- Use persistent queues for important messages
- Set up proper monitoring and alerting
- Implement circuit breakers for external dependencies

### Deployment
- Use external NATS service (NATS Cloud, etc.)
- Configure proper security and authentication
- Set up monitoring and log aggregation
- Implement proper backup and recovery procedures

## Testing

### Unit Testing
- Mock the messaging interfaces for command/query tests
- Test event serialization/deserialization
- Test event handler logic in isolation

### Integration Testing
- Use embedded NATS server for integration tests
- Test complete event flow from publishing to handling
- Test error scenarios and recovery

### Load Testing
- Test message throughput and latency
- Test system behavior under high event volumes
- Test failover and recovery scenarios

This NATS messaging implementation provides a robust foundation for event-driven architecture in the Go MVC application, enabling scalable and maintainable asynchronous communication patterns.