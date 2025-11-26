# Chapter 3: Quick Reference

Essential commands, configurations, and troubleshooting guide for daily development with Go MVC.

## 📋 Command Cheat Sheet

### Application Management

```bash
# Build application
make build

# Run application
make run

# Run with hot reload (development)
make dev

# Clean build artifacts
make clean

# Run all quality checks
make quality
```

### Development

```bash
# Format code
make fmt

# Run linter
make lint

# Run tests
make test

# Test with coverage
make test-coverage

# Run security checks
make security-check
```

### Database

```bash
# Create new migration
make migrate-create name=create_users_table

# Run all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Rollback all migrations
make migrate-reset

# Check migration status
make migrate-status

# Force migration version (use carefully!)
make migrate-force version=20240101000001
```

### Docker

```bash
# Start all services
docker-compose up -d

# Start specific service
docker-compose up -d postgres

# Stop all services
docker-compose down

# View logs
docker-compose logs -f

# Rebuild and restart
docker-compose up -d --build

# Remove all data (destructive!)
docker-compose down -v
```

### Testing

```bash
# Run all tests
make test

# Run specific package tests
go test ./internal/domain/user/...

# Run tests with coverage
make test-coverage

# Run tests with race detector
make test-race

# Verbose test output
make test-verbose

# Run specific test
go test -run TestUserCreate ./internal/domain/user/
```

## ⚙️ Configuration

### Environment Variables

```bash
# Application
APP_ENV=development|production|staging
APP_PORT=8080
APP_DEBUG=true|false

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=go_mvc_dev
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=3600

# Email
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_FROM=noreply@example.com
```

### Configuration Files

**Development**: `configs/development.yaml`
```yaml
server:
  port: 8080
  timeout: 30s
  
database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  database: go_mvc_dev
  ssl_mode: disable
  
redis:
  host: localhost
  port: 6379
  database: 0
  
jwt:
  secret: dev-secret-key
  expiration: 3600
```

**Production**: `configs/production.yaml`
```yaml
server:
  port: 8080
  timeout: 30s
  
database:
  host: ${DB_HOST}
  port: ${DB_PORT}
  user: ${DB_USER}
  password: ${DB_PASSWORD}
  database: ${DB_NAME}
  ssl_mode: require
  max_connections: 100
  
redis:
  host: ${REDIS_HOST}
  port: ${REDIS_PORT}
  password: ${REDIS_PASSWORD}
  
jwt:
  secret: ${JWT_SECRET}
  expiration: 3600
```

## 🏗️ Project Structure Quick View

```
go-mvc/
├── cmd/                          # Entry points
│   ├── main.go                  # Main server
│   ├── cli/                     # CLI tool
│   ├── worker/                  # Background worker
│   └── migrate/                 # Migration tool
├── internal/                     # Private code
│   ├── domain/                  # Domain layer
│   │   ├── user/               # User aggregate
│   │   ├── auth/               # Auth aggregate
│   │   └── repositories/       # Repository interfaces
│   ├── application/             # Application layer
│   │   ├── commands/           # Write operations
│   │   ├── queries/            # Read operations
│   │   ├── dto/                # Data transfer objects
│   │   ├── services/           # Application services
│   │   └── validators/         # Input validation
│   ├── infrastructure/          # Infrastructure layer
│   │   ├── persistence/        # Database
│   │   ├── cache/              # Caching
│   │   ├── messaging/          # Message queue
│   │   └── external/           # External services
│   ├── presentation/            # Presentation layer
│   │   └── http/
│   │       ├── handlers/       # HTTP handlers
│   │       └── middleware/     # Middleware
│   └── modules/                 # DI modules
├── pkg/                          # Public packages
│   ├── errors/                  # Error types
│   ├── jwt/                     # JWT utils
│   ├── pagination/              # Pagination
│   └── response/                # HTTP responses
└── configs/                      # Configuration files
```

## 🔧 Common Tasks

### Add New Feature

```bash
# 1. Write User Story (AI-powered)
cp docs/book/05-ai-development/user-story-template.md features/new-feature.md

# 2. Generate code with AI (see Chapter 23)

# 3. Create migration
make migrate-create name=create_new_feature_table

# 4. Run migration
make migrate-up

# 5. Test
make test

# 6. Run locally
make dev
```

### Debug Application

```bash
# Enable debug logging
export APP_DEBUG=true

# Run with dlv debugger
dlv debug cmd/main.go

# Check application health
curl http://localhost:8080/health

# View metrics
curl http://localhost:8080/metrics

# Check pprof
go tool pprof http://localhost:8080/debug/pprof/heap
```

### Database Operations

```bash
# Connect to database
psql -h localhost -U postgres -d go_mvc_dev

# Dump database
pg_dump -h localhost -U postgres go_mvc_dev > backup.sql

# Restore database
psql -h localhost -U postgres go_mvc_dev < backup.sql

# Reset database completely
make migrate-reset
make migrate-up

# Seed development data
go run cmd/cli/main.go seed
```

### Performance Testing

```bash
# Benchmark tests
go test -bench=. -benchmem ./...

# Load testing with hey
hey -n 10000 -c 100 http://localhost:8080/api/v1/health

# Profile CPU
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30

# Profile memory
go tool pprof http://localhost:8080/debug/pprof/heap
```

## 🐛 Troubleshooting

### Application Won't Start

**Symptom**: Application crashes on startup

**Check**:
```bash
# Verify configuration
cat configs/development.yaml

# Check if port is available
lsof -i :8080

# Check database connectivity
psql -h localhost -U postgres -c "SELECT 1"

# Check Redis connectivity  
redis-cli ping

# View application logs
make run 2>&1 | tee app.log
```

### Database Issues

**Symptom**: Database connection errors

**Solutions**:
```bash
# Check PostgreSQL is running
docker-compose ps postgres

# Restart PostgreSQL
docker-compose restart postgres

# Check PostgreSQL logs
docker-compose logs postgres

# Verify credentials
psql -h localhost -U postgres -d go_mvc_dev -c "SELECT version()"

# Check migrations
make migrate-status
```

### Migration Failures

**Symptom**: Migration won't run or fails

**Solutions**:
```bash
# Check migration status
make migrate-status

# View migration errors
make migrate-up 2>&1 | tee migration.log

# Rollback and retry
make migrate-down
make migrate-up

# Force to specific version (last resort)
make migrate-force version=20240101000001

# Recreate database
docker-compose down postgres
docker-compose up -d postgres
make migrate-up
```

### Performance Issues

**Symptom**: Slow API responses

**Debug**:
```bash
# Check metrics
curl http://localhost:8080/metrics | grep http_request_duration

# Enable tracing
# Visit Jaeger UI: http://localhost:16686

# Profile application
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30

# Check database queries
# Enable query logging in configs/development.yaml:
database:
  log_level: debug

# Check slow queries in PostgreSQL
psql -h localhost -U postgres -c "
  SELECT query, calls, total_time, mean_time
  FROM pg_stat_statements
  ORDER BY total_time DESC
  LIMIT 10
"
```

### Memory Leaks

**Symptom**: Memory usage grows over time

**Debug**:
```bash
# Capture heap profile
curl http://localhost:8080/debug/pprof/heap > heap.prof

# Analyze heap
go tool pprof heap.prof

# Compare two heap snapshots
curl http://localhost:8080/debug/pprof/heap > heap1.prof
# Wait some time
curl http://localhost:8080/debug/pprof/heap > heap2.prof
go tool pprof -base=heap1.prof heap2.prof
```

### Testing Issues

**Symptom**: Tests fail unexpectedly

**Solutions**:
```bash
# Run tests with verbose output
go test -v ./...

# Run specific test
go test -v -run TestUserCreate ./internal/domain/user/

# Clear test cache
go clean -testcache

# Run tests with race detector
go test -race ./...

# Check for concurrent issues
go test -count=100 ./internal/domain/user/
```

### Build Issues

**Symptom**: Build fails

**Solutions**:
```bash
# Clean everything
make clean
go clean -cache -modcache -testcache

# Update dependencies
go mod tidy
go mod download

# Verify dependencies
go mod verify

# Update specific dependency
go get -u github.com/gin-gonic/gin

# Check for breaking changes
go list -u -m all
```

## 📊 Monitoring

### Health Checks

```bash
# Application health
curl http://localhost:8080/health

# Expected response:
{
  "status": "ok",
  "version": "1.0.0",
  "services": {
    "database": "healthy",
    "cache": "healthy",
    "messaging": "healthy"
  }
}
```

### Metrics

```bash
# Prometheus metrics
curl http://localhost:8080/metrics

# Common metrics to watch:
# - http_requests_total
# - http_request_duration_seconds
# - db_connections_open
# - cache_hits_total
# - cache_misses_total
```

### Logs

```bash
# Application logs (structured JSON)
tail -f logs/app.log

# Filter by level
tail -f logs/app.log | jq 'select(.level == "error")'

# Filter by user ID
tail -f logs/app.log | jq 'select(.user_id == "some-uuid")'

# Watch for errors
tail -f logs/app.log | jq 'select(.level == "error" or .level == "fatal")'
```

## 🔐 Security

### Security Checklist

```bash
# Run security audit
make security-check

# Check for vulnerabilities
govulncheck ./...

# Scan dependencies
go list -json -deps | nancy sleuth

# Lint for security issues
gosec ./...
```

### Common Security Issues

**SQL Injection Prevention**
```go
// ❌ Don't: String concatenation
query := "SELECT * FROM users WHERE email = '" + email + "'"

// ✅ Do: Use parameterized queries
query := "SELECT * FROM users WHERE email = ?"
db.Query(query, email)
```

**Password Handling**
```go
// ❌ Don't: Store plain passwords
user.Password = password

// ✅ Do: Hash passwords
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
user.Password = string(hashedPassword)
```

**JWT Secret**
```bash
# ❌ Don't: Use weak secrets
JWT_SECRET=secret123

# ✅ Do: Generate strong secrets
JWT_SECRET=$(openssl rand -base64 32)
```

## 🎯 Best Practices

### Code Organization
- Follow Clean Architecture layers
- Keep domain logic pure
- Use dependency injection
- Write testable code

### Performance
- Use caching strategically
- Optimize database queries
- Index frequently queried columns
- Use connection pooling

### Security
- Always validate input
- Use parameterized queries
- Hash passwords with bcrypt
- Implement rate limiting

### Testing
- Write unit tests for business logic
- Integration tests for workflows
- Mock external dependencies
- Achieve >80% coverage

## 📚 Additional Resources

### Documentation
- [Architecture Guide](../02-architecture/01-architecture-overview.md)
- [API Development](../03-development-guide/04-api-development.md)
- [AI Code Generation](../05-ai-development/01-ai-quick-start.md)

### External Links
- [Go Documentation](https://golang.org/doc/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Gin Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)

---

**Need more help?** Check the [full documentation](../../BOOK.md) or [open an issue](https://github.com/tranvuongduy2003/go-mvc/issues).
