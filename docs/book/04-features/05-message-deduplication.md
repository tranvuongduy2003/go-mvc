# Message Deduplication Implementation

Hệ thống này đã được triển khai với các tính năng khắc phục trùng lặp message như yêu cầu:

## 🔹 1. Message ID duy nhất

### Outbox Pattern
- Mỗi message được gán một UUID duy nhất (`message_id`)
- Producer lưu message vào `outbox_messages` table trong cùng transaction với dữ liệu chính
- Background job đọc từ outbox và publish message với unique ID

### Cách sử dụng:
```go
// Trong business logic, sử dụng outbox service
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Tạo user
        user := &domain.User{...}
        if err := tx.Create(user).Error; err != nil {
            return err
        }
        
        // Lưu event vào outbox (cùng transaction)
        return s.outboxService.StoreMessage(ctx, tx, "user.created", user.ID.String(), user)
    })
}
```

## 🔹 2. Idempotent Consumer

### Inbox Pattern
Consumer sử dụng `inbox_messages` table để tracking message đã xử lý:

```go
// Trong message handler
func (h *UserEventHandler) HandleUserCreated(msg messaging.Message) error {
    ctx := context.Background()
    
    return h.inboxService.ProcessWithInboxPattern(
        ctx, nil, msg.ID, msg.EventType, "user-service",
        func(ctx context.Context, tx interface{}) error {
            // Business logic chỉ chạy nếu message chưa được xử lý
            return h.processUserCreatedEvent(ctx, msg.Payload)
        },
    )
}
```

### Message Deduplication (Lightweight)
Cho các trường hợp đơn giản hơn:

```go
func (h *Handler) ProcessMessage(msg messaging.Message) error {
    ctx := context.Background()
    ttl := 24 * time.Hour
    
    return h.inboxService.ProcessWithIdempotency(
        ctx, msg.ID, msg.EventType, "consumer-id", ttl,
        func(ctx context.Context) error {
            // Business logic
            return h.doWork(msg)
        },
    )
}
```

## 🔹 3. HTTP Idempotency Middleware

### Sử dụng middleware:
```go
// Trong router setup
func SetupRoutes(r *gin.Engine, deps *Dependencies) {
    // Apply idempotency middleware
    idempotencyMiddleware := deps.IdempotencyMiddleware
    
    api := r.Group("/api/v1")
    api.Use(idempotencyMiddleware.Handler())
    
    // Hoặc với custom options
    api.Use(idempotencyMiddleware.WithOptions(middleware.IdempotencyOptions{
        TTL:            time.Hour,
        RequireKey:     true,
        IgnoredPaths:   []string{"/health", "/metrics"},
        IgnoredMethods: []string{"GET", "HEAD"},
    }))
    
    api.POST("/users", userHandler.CreateUser)
}
```

### Client usage:
```bash
# Client gửi request với Idempotency-Key header
curl -X POST /api/v1/users \
  -H "Idempotency-Key: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
  
# Request duplicate sẽ return 409 Conflict
```

## 🔹 4. NATS Enhanced Adapter

### Subscription với deduplication:
```go
// Lightweight deduplication
broker.SubscribeWithDeduplication("user.events", handler, 24*time.Hour)

// Full inbox pattern
broker.SubscribeWithInbox("order.events", handler)

// Idempotent subscription
broker.CreateIdempotentSubscription("payment.events", handler, 
    nats.IdempotencyOptions{
        UseInboxPattern:  true,
        DeduplicationTTL: time.Hour,
    })
```

## 🔹 5. Background Jobs

### Outbox Processor Job
```go
// Job chạy định kỳ để xử lý outbox messages
job := NewOutboxProcessorJob(outboxService, publisher, batchSize, retryDelay)

// Execute job
if err := job.ExecuteWithRetry(ctx); err != nil {
    log.Printf("Outbox processing failed: %v", err)
}

// Cleanup old messages
if err := job.CleanupOldMessages(ctx, 7); err != nil { // Remove messages older than 7 days
    log.Printf("Cleanup failed: %v", err)
}
```

## 🗄️ Database Schema

### Tables created:
1. **outbox_messages** - Outbox pattern implementation
2. **inbox_messages** - Full inbox pattern for consumers
3. **message_deduplication** - Lightweight deduplication with TTL

### Indexes tối ưu:
- Unique constraints cho message_id + consumer_id
- Indexes cho performance queries
- Partial indexes cho expired records cleanup

## 🔧 Configuration

### Dependency Injection:
```go
// Trong main.go hoặc DI setup
app := fx.New(
    // ... other modules
    modules.MessagingModule, // Thêm messaging module
    fx.Invoke(func(
        outboxJob *handlers.OutboxProcessorJob,
        scheduler *scheduler.Scheduler,
    ) {
        // Schedule outbox processing job
        scheduler.AddJob("outbox-processor", outboxJob, 30*time.Second)
    }),
)
```

## 🧪 Testing

### Unit tests:
- Test outbox service với mock repositories
- Test inbox service deduplication logic
- Test middleware idempotency behavior

### Integration tests:
- Test end-to-end message flow
- Test duplicate message handling
- Test failure scenarios và retry logic

## 📊 Monitoring

### Metrics to track:
- Outbox processing latency
- Message deduplication hit rate
- Failed message retry counts
- HTTP idempotency conflicts

### Cleanup jobs:
```go
// Schedule cleanup jobs
scheduler.AddJob("cleanup-outbox", func() {
    outboxService.CleanupOldMessages(ctx, 7) // 7 days
}, 24*time.Hour)

scheduler.AddJob("cleanup-deduplication", func() {
    inboxService.CleanupExpiredDeduplicationRecords(ctx)
}, 4*time.Hour)
```

## 🚀 Benefits

✅ **Exactly-once processing**: Message chỉ được xử lý đúng 1 lần
✅ **Reliable publishing**: Outbox pattern đảm bảo message không bị mất
✅ **HTTP idempotency**: API calls an toàn với retry
✅ **Performance optimized**: Indexes và cleanup jobs tự động
✅ **Flexible patterns**: Hỗ trợ cả inbox pattern và lightweight deduplication
✅ **Monitoring ready**: Built-in metrics và logging

Hệ thống này cung cấp giải pháp toàn diện cho vấn đề trùng lặp message trong distributed systems.