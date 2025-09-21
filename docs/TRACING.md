# Distributed Tracing với OpenTelemetry

## Tổng quan

Ứng dụng đã được tích hợp distributed tracing sử dụng OpenTelemetry để theo dõi và debug các request qua các tầng của ứng dụng.

## Tính năng Tracing

### 1. Tracing Infrastructure
- **OpenTelemetry SDK**: Core tracing functionality
- **OTLP HTTP Exporter**: Gửi traces tới Jaeger/OTLP endpoint
- **Gin Instrumentation**: Tự động trace HTTP requests
- **Custom Spans**: Manual tracing cho business logic

### 2. Tracing Levels
- **HTTP Layer**: Automatic request/response tracing
- **Application Service Layer**: Business logic tracing
- **Repository Layer**: Database operation tracing
- **External Service Layer**: Third-party service calls

### 3. Traced Components
- HTTP requests và responses
- Database operations (CRUD)
- Application service calls
- Domain service operations
- External service integrations

## Cấu hình

### Development Config (`configs/development.yaml`)
```yaml
tracing:
  enabled: true
  service_name: "go-mvc"
  endpoint: "http://localhost:4318"  # Jaeger OTLP HTTP endpoint
  sample_rate: 1.0  # 100% sampling for development
```

### Production Config (`configs/production.yaml`)
```yaml
tracing:
  enabled: true
  service_name: "go-mvc"
  endpoint: "http://jaeger-collector:4318"
  sample_rate: 0.1  # 10% sampling for production
```

## Setup Jaeger (Local Development)

### Sử dụng Docker
```bash
# Run Jaeger All-in-One
docker run -d --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 9411:9411 \
  jaegertracing/all-in-one:latest
```

### Hoặc sử dụng Docker Compose
```yaml
version: '3'
services:
  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686"  # Jaeger UI
      - "4317:4317"    # OTLP gRPC
      - "4318:4318"    # OTLP HTTP
    environment:
      - COLLECTOR_OTLP_ENABLED=true
```

## Chạy ứng dụng

```bash
# Build và chạy
go run cmd/api/main.go

# Hoặc build binary
go build -o bin/server cmd/api/main.go
./bin/server
```

## Test Tracing

### 1. Health Check với Tracing
```bash
curl http://localhost:8080/health
```

### 2. Example Endpoint với Custom Tracing
```bash
curl http://localhost:8080/example
```

### 3. Error Endpoint để test Error Tracing
```bash
curl http://localhost:8080/error
```

## Xem Traces

1. Mở Jaeger UI: http://localhost:16686
2. Chọn service "go-mvc"
3. Click "Find Traces" để xem danh sách traces
4. Click vào một trace để xem chi tiết

## Trace Information

Mỗi trace sẽ bao gồm:

### HTTP Spans
- Method, Path, Status Code
- Request/Response headers
- Client IP, User Agent
- Processing time

### Application Spans
- Service operation names
- Input parameters
- Business logic metrics
- Error information

### Database Spans
- SQL operation type (SELECT, INSERT, UPDATE, DELETE)
- Table names
- Query execution time
- Database connection info

## Custom Tracing trong Code

### 1. Service Layer Tracing
```go
func (s *UserApplicationService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserDTO, error) {
    ctx, span := s.tracing.StartServiceSpan(ctx, "UserApplicationService", "CreateUser")
    defer span.End()
    
    span.SetAttributes(
        attribute.String("user.email", req.Email),
        attribute.String("user.username", req.Username),
    )
    
    // Business logic...
    
    if err != nil {
        s.tracing.RecordError(span, err)
        return nil, err
    }
    
    return result, nil
}
```

### 2. Repository Layer Tracing
```go
func (r *UserRepository) Create(ctx context.Context, user *userDomain.User) error {
    ctx, span := r.tracing.StartDatabaseSpan(ctx, "INSERT", "users")
    defer span.End()
    
    span.SetAttributes(
        attribute.String("user.id", user.ID().String()),
        attribute.String("user.email", user.Email().String()),
    )
    
    if err := r.db.WithContext(ctx).Create(userModel).Error; err != nil {
        r.tracing.RecordError(span, err)
        return err
    }
    
    return nil
}
```

### 3. HTTP Handler Tracing
```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    ctx := middleware.TraceContext(c)
    ctx, span := h.tracing.StartHTTPSpan(ctx, c.Request.Method, c.FullPath())
    defer span.End()
    
    // Handler logic...
    
    span.SetAttributes(
        attribute.String("user.id", user.ID.String()),
        attribute.String("response.status", "success"),
    )
}
```

## Performance Impact

- **Development**: 100% sampling - suitable for debugging
- **Production**: 10% sampling - minimal performance impact
- **Resource Usage**: ~1-2% CPU overhead, ~5MB memory per 1000 spans

## Troubleshooting

### 1. Tracing không hoạt động
- Kiểm tra Jaeger đang chạy: `docker ps | grep jaeger`
- Kiểm tra endpoint trong config
- Kiểm tra logs: có thông báo lỗi tracing không

### 2. Không thấy traces trong Jaeger UI
- Đợi vài giây để traces được flush
- Kiểm tra time range trong Jaeger UI
- Kiểm tra service name

### 3. Performance issues
- Giảm sample_rate trong config
- Kiểm tra batch size và flush interval

## Monitoring Metrics

Các metrics quan trọng để monitor:
- Trace completion rate
- Average span duration
- Error rate per service
- Database query performance
- External service latency