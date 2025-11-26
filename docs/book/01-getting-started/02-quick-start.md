# Chapter 2: Quick Start

## Get Started in 15 Minutes

This chapter will get you from zero to running application in just 15 minutes. We'll cover installation, running your first API, and understanding the basic structure.

## Prerequisites Checklist

Before starting, ensure you have:

- ✅ **Go 1.24.5+** installed ([Download](https://golang.org/dl/))
- ✅ **Docker & Docker Compose** installed ([Download](https://www.docker.com/products/docker-desktop))
- ✅ **Make** installed (usually pre-installed on macOS/Linux)
- ✅ **Git** installed for cloning the repository
- ✅ **8GB RAM** available for running services

### Verify Installation

```bash
# Check Go version
go version
# Should output: go version go1.24.5 or higher

# Check Docker
docker --version
# Should output: Docker version 20.10.x or higher

# Check Docker Compose
docker-compose --version
# Should output: docker-compose version 1.29.x or higher

# Check Make
make --version
# Should output: GNU Make 3.81 or higher
```

## Step 1: Clone and Setup (3 minutes)

### Clone the Repository

```bash
# Clone the repository
git clone https://github.com/tranvuongduy2003/go-mvc.git

# Navigate to project directory
cd go-mvc

# Check the structure
ls -la
```

You should see:
```
.
├── cmd/                    # Application entry points
├── internal/               # Private application code
├── pkg/                    # Public packages
├── configs/                # Configuration files
├── docs/                   # Documentation
├── scripts/                # Build and deployment scripts
├── docker-compose.yml      # Docker services definition
├── Makefile               # Build automation
├── go.mod                 # Go module definition
└── README.md              # Project overview
```

### Install Dependencies

```bash
# Download Go dependencies
go mod download

# Verify dependencies
go mod verify
```

## Step 2: Start Development Services (5 minutes)

### Start Infrastructure Services

```bash
# Start PostgreSQL, Redis, and other services
docker-compose up -d postgres redis

# Wait for services to be ready (about 30 seconds)
docker-compose ps
```

You should see services running:
```
NAME                COMMAND                  STATUS
go-mvc-postgres     "docker-entrypoint.s…"   Up
go-mvc-redis        "docker-entrypoint.s…"   Up
```

### Configure Environment

```bash
# Copy example configuration
cp configs/development.yaml.example configs/development.yaml

# The default configuration works out of the box!
# No changes needed for local development
```

### Run Database Migrations

```bash
# Create database schema
make migrate-up

# Verify migration status
make migrate-status
```

Expected output:
```
Version   Dirty  Timestamp
1         false  2024-01-01 00:00:00 +0000 UTC
2         false  2024-01-02 00:00:00 +0000 UTC
```

## Step 3: Start the Application (2 minutes)

### Build and Run

```bash
# Build the application
make build

# Run the application
make run
```

Or use hot reload for development:

```bash
# Install Air (if not already installed)
go install github.com/cosmtrek/air@latest

# Run with hot reload
make dev
```

### Verify Application is Running

```bash
# Check health endpoint
curl http://localhost:8080/health

# Expected response:
{
  "status": "ok",
  "version": "1.0.0",
  "services": {
    "database": "healthy",
    "cache": "healthy"
  }
}
```

## Step 4: Make Your First API Call (5 minutes)

### Create a User Account

```bash
# Register a new user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "Admin123!",
    "name": "Admin User"
  }'
```

Response:
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "admin@example.com",
    "name": "Admin User",
    "created_at": "2024-01-01T00:00:00Z"
  },
  "message": "User registered successfully"
}
```

### Login

```bash
# Login to get JWT token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "Admin123!"
  }'
```

Response:
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

### Make Authenticated Request

```bash
# Save token for convenience
export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Get current user profile
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN"
```

Response:
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "admin@example.com",
    "name": "Admin User",
    "roles": ["user"],
    "permissions": ["read:users", "write:users"],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

## Step 5: Explore the API (Bonus)

### Access API Documentation

Open your browser and navigate to:

```
http://localhost:8080/swagger/index.html
```

You'll see the complete API documentation with:
- All available endpoints
- Request/response schemas
- Try-it-out functionality
- Authentication setup

### Available Endpoints

**Authentication**
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh token
- `POST /api/v1/auth/logout` - Logout

**User Management**
- `GET /api/v1/users` - List users (admin only)
- `GET /api/v1/users/:id` - Get user by ID
- `GET /api/v1/users/me` - Get current user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user (admin only)

**Role & Permission**
- `GET /api/v1/roles` - List roles
- `POST /api/v1/roles` - Create role (admin only)
- `GET /api/v1/permissions` - List permissions

## Development Workflow

### Hot Reload Development

```bash
# Start with Air (watches for file changes)
make dev

# Make changes to any .go file
# Application automatically rebuilds and restarts
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test ./internal/domain/user/...

# Run with verbose output
make test-verbose
```

### Database Management

```bash
# Create new migration
make migrate-create name=add_users_table

# Run migrations
make migrate-up

# Rollback last migration
make migrate-down

# Rollback all migrations
make migrate-reset

# Check migration status
make migrate-status
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run security checks
make security-check

# Run all quality checks
make quality
```

## Common Commands Reference

### Application Management

```bash
make build          # Build application
make run            # Run application
make dev            # Run with hot reload
make clean          # Clean build artifacts
```

### Database

```bash
make migrate-up     # Run migrations
make migrate-down   # Rollback migration
make migrate-create # Create new migration
make migrate-status # Check migration status
make migrate-reset  # Reset database
```

### Testing

```bash
make test           # Run tests
make test-coverage  # Test with coverage
make test-verbose   # Verbose test output
make test-race      # Test with race detector
```

### Code Quality

```bash
make fmt            # Format code
make lint           # Run linter
make vet            # Run go vet
make security-check # Security audit
```

### Docker

```bash
make docker-build   # Build Docker image
make docker-run     # Run in Docker
make docker-push    # Push to registry
```

## Troubleshooting

### Port Already in Use

**Error**: `bind: address already in use`

**Solution**:
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or change port in configs/development.yaml
server:
  port: 8081
```

### Database Connection Failed

**Error**: `connection refused`

**Solution**:
```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# Restart PostgreSQL
docker-compose restart postgres

# Check logs
docker-compose logs postgres
```

### Migration Failed

**Error**: `migration version X failed`

**Solution**:
```bash
# Check migration status
make migrate-status

# Force version (careful!)
migrate -path internal/infrastructure/persistence/postgres/migrations \
  -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" \
  force <version>

# Then try again
make migrate-up
```

### Dependencies Not Found

**Error**: `package not found`

**Solution**:
```bash
# Clean module cache
go clean -modcache

# Download dependencies again
go mod download

# Tidy dependencies
go mod tidy
```

## Next Steps

Congratulations! 🎉 You now have a running Go MVC application.

### Continue Learning

1. **[Chapter 3: Quick Reference](03-quick-reference.md)** - Essential commands and configurations
2. **[Chapter 4: Architecture Overview](../02-architecture/01-architecture-overview.md)** - Understand the architecture
3. **[Chapter 12: Development Workflow](../03-development-guide/01-development-workflow.md)** - Deep dive into development

### Build Your First Feature

Ready to build something? Try these guides:

- **[Chapter 23: AI Quick Start](../05-ai-development/01-ai-quick-start.md)** - Generate your first API with AI
- **[Chapter 14: API Development](../03-development-guide/04-api-development.md)** - Manual API development
- **[Chapter 16: Authentication & Authorization](../04-features/01-authentication.md)** - Add auth to your features

### Explore Features

- **[Chapter 17: Background Jobs](../04-features/02-background-jobs.md)** - Async processing
- **[Chapter 18: Email Service](../04-features/03-email-service.md)** - Send emails
- **[Chapter 22: Distributed Tracing](../04-features/07-tracing.md)** - Monitor your application

---

**You're all set! Happy coding!** 🚀

Need help? Check the [troubleshooting guide](03-quick-reference.md#troubleshooting) or [open an issue](https://github.com/tranvuongduy2003/go-mvc/issues).
