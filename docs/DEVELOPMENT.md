# Development Guide

## 📋 Table of Contents
- [Development Environment Setup](#development-environment-setup)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Code Quality](#code-quality)
- [Debugging](#debugging)
- [Contributing](#contributing)
- [Best Practices](#best-practices)

## 🛠️ Development Environment Setup

### Prerequisites

Ensure you have the following installed:

```bash
# Check Go version (required: 1.24.5+)
go version

# Check Docker and Docker Compose
docker --version
docker-compose --version

# Check Make
make --version

# Check Git
git --version
```

### Required Tools

| Tool | Version | Purpose |
|------|---------|---------|
| **Go** | 1.24.5+ | Primary language |
| **Docker** | 20.0+ | Containerization |
| **Docker Compose** | 2.0+ | Multi-container orchestration |
| **Make** | 4.0+ | Build automation |
| **Git** | 2.30+ | Version control |

### Optional Development Tools

```bash
# Install development tools
make setup-tools
```

This installs:
- **Air**: Hot reload for Go applications
- **golangci-lint**: Comprehensive Go linter
- **Swagger**: API documentation generator
- **goimports**: Import formatter
- **gofumpt**: Stricter gofmt
- **gosec**: Security analyzer

### Initial Setup

1. **Clone the repository:**
```bash
git clone https://github.com/tranvuongduy2003/go-mvc.git
cd go-mvc
```

2. **Setup development environment:**
```bash
# This will:
# - Download Go dependencies
# - Start database and Redis containers
# - Run database migrations
make setup
```

3. **Create environment file:**
```bash
cp .env.example .env
# Edit .env with your local settings
```

4. **Verify setup:**
```bash
# Check if all services are running
make status
```

## 📁 Project Structure

### Key Development Directories

```bash
# Source code
internal/           # Private application code
├── core/          # Domain layer (business logic)
├── application/   # Use cases and application services
├── adapters/      # Infrastructure implementations
├── handlers/      # HTTP handlers and middleware
└── shared/        # Shared utilities

# Public packages
pkg/               # Reusable packages

# Configuration
configs/           # Environment configurations

# Entry points
cmd/               # Application executables
├── main.go       # Main HTTP server
├── cli/          # Command-line tools
├── worker/       # Background workers
└── migrate/      # Database migration tool
```

### Development Files

```bash
# Build and automation
Makefile          # Development commands
docker-compose.yml # Local services

# Environment
.env.example      # Environment template
.env              # Local environment (create from example)

# Documentation
docs/             # Project documentation
README.md         # Project overview
```

## 🔄 Development Workflow

### Daily Development Cycle

1. **Start Development Environment:**
```bash
# Start database and cache services
make docker-up-db

# Start monitoring (optional)
make monitoring
```

2. **Run Application with Hot Reload:**
```bash
# Hot reload with Air
make dev

# Or run directly
make run
```

3. **Make Changes and Test:**
```bash
# Run tests
make test

# Run linter
make lint

# Format code
make format
```

4. **Commit Changes:**
```bash
git add .
git commit -m "feat: add new feature"
git push
```

### Common Development Commands

```bash
# Development
make dev              # Run with hot reload
make run              # Run without hot reload
make build            # Build binary
make clean            # Clean build artifacts

# Testing
make test             # Run all tests
make test-unit        # Unit tests only
make test-integration # Integration tests only
make test-coverage    # Test with coverage report

# Code Quality
make lint             # Run linter
make format           # Format code
make vet              # Run go vet
make security         # Security scan

# Database
make migrate-up       # Apply migrations
make migrate-down     # Rollback migrations
make migrate-create   # Create new migration

# Docker
make docker-up        # Start all services
make docker-down      # Stop all services
make docker-logs      # View container logs

# Monitoring
make monitoring       # Start monitoring stack
make metrics          # View application metrics
make health           # Health check
```

## 🧪 Testing

### Test Structure

```bash
# Test files should be located alongside source code
internal/
├── application/
│   ├── services/
│   │   ├── user_service.go
│   │   └── user_service_test.go      # Unit tests
├── adapters/
│   ├── persistence/
│   │   ├── user_repository.go
│   │   └── user_repository_test.go   # Integration tests

# Integration tests (optional separate directory)
tests/
├── integration/      # Integration tests
├── e2e/             # End-to-end tests
└── mocks/           # Generated mocks
```

### Test Types

#### 1. Unit Tests
Test individual functions and methods in isolation.

```go
// internal/application/services/user_service_test.go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := mocks.NewUserRepository(t)
    service := NewUserService(mockRepo)
    
    user := &domain.User{
        Email: "test@example.com",
        Name:  "Test User",
    }
    
    mockRepo.On("Create", mock.Anything).Return(nil)
    
    // Act
    err := service.CreateUser(context.Background(), user)
    
    // Assert
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

#### 2. Integration Tests
Test component interactions with real dependencies.

```go
// internal/adapters/persistence/user_repository_test.go
// +build integration

func TestUserRepository_Integration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    repo := NewUserRepository(db)
    
    // Test repository operations
    user := &models.User{
        Email: "integration@example.com",
        Name:  "Integration Test",
    }
    
    err := repo.Create(context.Background(), user)
    assert.NoError(t, err)
    assert.NotEmpty(t, user.ID)
}
```

#### 3. End-to-End Tests
Test complete workflows through HTTP API.

```go
// tests/e2e/user_test.go
func TestUserAPI_E2E(t *testing.T) {
    // Setup test server
    app := setupTestApp(t)
    defer app.Cleanup()
    
    // Test user registration
    body := `{"email":"e2e@example.com","password":"test123","name":"E2E Test"}`
    resp := app.POST("/api/v1/auth/register").
        WithBytes([]byte(body)).
        Expect().
        Status(201)
    
    // Verify response
    resp.JSON().Object().
        Value("success").Boolean().True()
}
```

### Running Tests

```bash
# All tests
make test

# Unit tests only
make test-unit

# Integration tests (requires database)
make test-integration

# E2E tests
make test-e2e

# Test with coverage
make test-coverage

# Benchmarks
make benchmark
```

### Test Database Setup

Integration tests require a test database:

```bash
# Start test database
docker run -d --name test-postgres \
  -e POSTGRES_DB=go_mvc_test \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5433:5432 \
  postgres:15-alpine

# Run integration tests
go test -tags=integration ./...
```

## 🔍 Code Quality

### Linting

We use `golangci-lint` with comprehensive checks:

```bash
# Run linter
make lint

# Run with auto-fix
make lint-fix

# Configuration in .golangci.yml
```

### Code Formatting

```bash
# Format code
make format

# This runs:
# - go fmt (basic formatting)
# - goimports (import organization)
# - gofumpt (stricter formatting)
```

### Security Scanning

```bash
# Run security scanner
make security

# This runs gosec to find security issues
```

### Pre-commit Hooks

Set up pre-commit hooks to ensure code quality:

```bash
# Create .git/hooks/pre-commit
#!/bin/bash
set -e

echo "Running pre-commit checks..."

# Format code
make format

# Run linter
make lint

# Run tests
make test-unit

echo "Pre-commit checks passed!"
```

## 🐛 Debugging

### Development Debugging

#### 1. Using Delve Debugger

```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug main application
dlv debug cmd/main.go

# Debug with arguments
dlv debug cmd/main.go -- --config=configs/development.yaml

# Debug tests
dlv test ./internal/application/services
```

#### 2. VS Code Debug Configuration

Create `.vscode/launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug Application",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "./cmd/main.go",
            "env": {
                "ENV": "development"
            },
            "args": []
        },
        {
            "name": "Debug Tests",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "./internal/application/services"
        }
    ]
}
```

#### 3. Logging for Debug

```go
// Use structured logging
logger.Debug("User operation started",
    zap.String("user_id", userID),
    zap.String("operation", "create"),
)

// Add request ID for tracing
logger.Info("Processing request",
    zap.String("request_id", requestID),
    zap.String("method", "POST"),
    zap.String("path", "/api/v1/users"),
)
```

### Production Debugging

#### 1. Health Checks

```bash
# Application health
curl http://localhost:8080/health

# Check with Make
make health
```

#### 2. Metrics and Monitoring

```bash
# View metrics
make metrics

# Grafana dashboard
open http://localhost:3000

# Prometheus
open http://localhost:9091
```

#### 3. Distributed Tracing

```bash
# Jaeger UI
open http://localhost:16686

# Generate test traces
make trace-test
```

#### 4. Performance Profiling

```bash
# CPU profiling
make pprof-cpu

# Memory profiling
make pprof-mem

# Goroutine profiling
make pprof-goroutine
```

## 🤝 Contributing

### Git Workflow

1. **Create Feature Branch:**
```bash
git checkout -b feature/user-authentication
```

2. **Make Changes:**
```bash
# Edit files
# Run tests
make test

# Format and lint
make format
make lint
```

3. **Commit Changes:**
```bash
git add .
git commit -m "feat(auth): implement JWT authentication"
```

4. **Push and Create PR:**
```bash
git push origin feature/user-authentication
# Create Pull Request on GitHub
```

### Commit Message Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```bash
# Format: type(scope): description

feat(auth): add JWT authentication
fix(db): resolve connection pool issues
docs(api): update endpoint documentation
test(user): add unit tests for user service
refactor(handler): simplify error handling
chore(deps): update dependencies
```

### Pull Request Guidelines

1. **Clear Description:** Explain what and why
2. **Tests:** Include relevant tests
3. **Documentation:** Update docs if needed
4. **Small Changes:** Keep PRs focused and small
5. **Code Review:** Request review from team members

### Code Review Checklist

- [ ] Code follows project conventions
- [ ] Tests are included and passing
- [ ] Documentation is updated
- [ ] No security vulnerabilities
- [ ] Performance considerations addressed
- [ ] Error handling is appropriate

## 📚 Best Practices

### Go Best Practices

#### 1. Package Organization
```go
// Good: Clear package responsibility
package userservice

// Bad: Generic package names
package utils
```

#### 2. Interface Design
```go
// Good: Small, focused interfaces
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id string) (*User, error)
}

// Bad: Large interfaces
type UserService interface {
    // 20+ methods...
}
```

#### 3. Error Handling
```go
// Good: Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// Good: Use custom error types
var ErrUserNotFound = errors.New("user not found")
```

#### 4. Context Usage
```go
// Good: Pass context as first parameter
func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    // Use ctx for cancellation, timeouts, values
}
```

### Architecture Best Practices

#### 1. Dependency Direction
```go
// Good: Dependencies point inward
// Handler -> Application -> Domain
// Infrastructure implements interfaces

// Bad: Domain depends on infrastructure
```

#### 2. Domain Logic
```go
// Good: Rich domain models
func (u *User) ChangeEmail(newEmail string) error {
    if !u.IsEmailValid(newEmail) {
        return ErrInvalidEmail
    }
    u.Email = newEmail
    return nil
}

// Bad: Anemic domain models
type User struct {
    ID    string
    Email string
}
```

#### 3. Use Cases
```go
// Good: Clear use case responsibility
type CreateUserUseCase struct {
    userRepo UserRepository
    emailService EmailService
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, req CreateUserRequest) error {
    // Validation
    // Business logic
    // Persistence
    // Side effects
}
```

### Testing Best Practices

#### 1. Test Naming
```go
// Good: Descriptive test names
func TestUserService_CreateUser_WhenEmailExists_ReturnsError(t *testing.T) {}

// Bad: Generic test names
func TestCreateUser(t *testing.T) {}
```

#### 2. Test Structure (AAA Pattern)
```go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := mocks.NewUserRepository(t)
    service := NewUserService(mockRepo)
    user := &User{Email: "test@example.com"}
    
    // Act
    err := service.CreateUser(context.Background(), user)
    
    // Assert
    assert.NoError(t, err)
}
```

#### 3. Table-Driven Tests
```go
func TestUserValidation(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "user@example.com", false},
        {"invalid email", "invalid", true},
        {"empty email", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Performance Best Practices

#### 1. Database Queries
```go
// Good: Use select specific fields
SELECT id, name, email FROM users WHERE active = true

// Bad: Select all fields
SELECT * FROM users
```

#### 2. Caching Strategy
```go
// Good: Cache expensive operations
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    // Check cache first
    if user := s.cache.Get(id); user != nil {
        return user, nil
    }
    
    // Fallback to database
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Cache result
    s.cache.Set(id, user)
    return user, nil
}
```

#### 3. Connection Pooling
```yaml
# Database configuration
database:
  max_open_conns: 25
  max_idle_conns: 10
  conn_max_lifetime: 5m
```

### Security Best Practices

#### 1. Input Validation
```go
// Good: Validate all inputs
func CreateUser(req CreateUserRequest) error {
    if err := req.Validate(); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    // Process request
}
```

#### 2. SQL Injection Prevention
```go
// Good: Use parameterized queries
query := "SELECT * FROM users WHERE email = $1"
err := db.QueryRow(query, email).Scan(&user)

// Bad: String concatenation
query := "SELECT * FROM users WHERE email = '" + email + "'"
```

#### 3. Secrets Management
```bash
# Good: Use environment variables
export JWT_SECRET="your-secret-key"

# Bad: Hardcode secrets
const jwtSecret = "hardcoded-secret"
```

## 🚀 Production Readiness

### Checklist

- [ ] All tests passing
- [ ] Security scan clean
- [ ] Performance tested
- [ ] Monitoring configured
- [ ] Error handling comprehensive
- [ ] Documentation complete
- [ ] Deployment scripts ready
- [ ] Environment configs set

### Deployment Preparation

```bash
# Build production binary
GOOS=linux GOARCH=amd64 make build

# Run security scan
make security

# Test production config
CONFIG_FILE=configs/production.yaml make run

# Build Docker image
make docker-build
```

For deployment instructions, see [Deployment Guide](DEPLOYMENT.md).