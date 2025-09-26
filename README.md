# Go MVC - Enterprise Grade Go Web Application

[![Go Version](https://img.shields.io/badge/go-1.24.5-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/docker-ready-blue.svg)](docker-compose.yml)

A modern, scalable Go web application built with **Clean Architecture**, **Domain-Driven Design (DDD)**, and **Enterprise patterns**. Features comprehensive observability stack with Prometheus, Grafana, and Jaeger tracing.

## 🚀 Features

### Core Features
- **Clean Architecture**: Separated layers (Domain, Application, Infrastructure, Handlers)
- **Domain-Driven Design**: Rich domain models with business logic encapsulation
- **Dependency Injection**: Uber FX for modular dependency management
- **CQRS Pattern**: Command Query Responsibility Segregation
- **Repository Pattern**: Abstract data access layer

### Infrastructure
- **HTTP Framework**: Gin with custom middleware stack
- **Database**: PostgreSQL with GORM ORM
- **Cache**: Redis for high-performance caching
- **Message Queue**: RabbitMQ for async processing
- **Authentication**: JWT-based with security middleware
- **Email Service**: SMTP integration with MailCatcher for development testing

### Observability Stack
- **Metrics**: Prometheus + Grafana dashboards
- **Distributed Tracing**: Jaeger with OpenTelemetry
- **Logging**: Structured logging with Zap
- **Health Checks**: Comprehensive health monitoring
- **Performance Profiling**: pprof integration

### Development Tools
- **Hot Reload**: Air for development
- **Code Quality**: golangci-lint, gosec, gofumpt
- **Testing**: Unit, integration, and benchmark tests
- **Documentation**: Swagger/OpenAPI generation
- **Database Migrations**: golang-migrate/migrate v4.19.0+ with timestamp-based versioning

## 📁 Project Structure

```
go-mvc/
├── cmd/                    # Application entry points
│   ├── main.go            # Main server application
│   ├── cli/               # Command-line interface
│   ├── worker/            # Background worker
│   └── migrate/           # Database migration tool
├── internal/              # Private application code
│   ├── adapters/          # External adapters (Infrastructure layer)
│   │   ├── cache/         # Cache implementations
│   │   ├── external/      # External service clients
│   │   ├── messaging/     # Message queue adapters
│   │   ├── persistence/   # Database adapters
│   │   └── repositories/  # Repository implementations
│   ├── application/       # Application layer (Use cases)
│   │   ├── commands/      # Command handlers
│   │   ├── queries/       # Query handlers
│   │   ├── dto/           # Data transfer objects
│   │   ├── services/      # Application services
│   │   └── validators/    # Input validation
│   ├── core/              # Core domain layer
│   │   ├── domain/        # Domain entities and business logic
│   │   └── ports/         # Interface definitions
│   ├── fx_modules/        # Dependency injection modules
│   ├── handlers/          # HTTP handlers and middleware
│   └── shared/            # Shared utilities and infrastructure
├── pkg/                   # Public packages (can be imported)
├── configs/               # Configuration files
├── deployments/           # Deployment configurations
├── docs/                  # Documentation
├── scripts/              # Build and deployment scripts
└── api/                  # API specifications
```

## 🛠️ Quick Start

### Prerequisites
- **Go 1.24.5+**
- **Docker & Docker Compose**
- **Make** (for build automation)

### 1. Clone and Setup
```bash
git clone https://github.com/tranvuongduy2003/go-mvc.git
cd go-mvc

# Setup development environment
make setup
```

### 2. Start Development Services
```bash
# Start database and cache
make docker-up-db

# Start monitoring stack
make monitoring
```

### 3. Run Application
```bash
# Development with hot reload
make dev

# Or run directly
make run
```

### 4. Access Services
- **Application**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **Metrics**: http://localhost:8080/metrics
- **MailCatcher**: http://localhost:1080 (email testing interface)
- **Prometheus**: http://localhost:9091
- **Grafana**: http://localhost:3000 (admin/admin)
- **Jaeger**: http://localhost:16686

## 🏗️ Architecture Overview

This application follows **Clean Architecture** principles with clear separation of concerns:

### 1. **Domain Layer** (`internal/core/domain/`)
- **Entities**: Core business objects
- **Value Objects**: Immutable data structures
- **Domain Events**: Business event definitions
- **Specifications**: Business rule definitions

### 2. **Application Layer** (`internal/application/`)
- **Commands**: Write operations (CQRS)
- **Queries**: Read operations (CQRS)
- **Services**: Application business logic
- **DTOs**: Data transfer objects
- **Validators**: Input validation logic

### 3. **Infrastructure Layer** (`internal/adapters/`)
- **Persistence**: Database repositories
- **Cache**: Caching implementations
- **External**: Third-party service clients
- **Messaging**: Event bus implementations

### 4. **Presentation Layer** (`internal/handlers/`)
- **HTTP**: REST API handlers
- **Middleware**: Cross-cutting concerns
- **Validators**: Request validation
- **Responses**: Response formatting

## 🔧 Development Commands

### Build & Run
```bash
make build          # Build server binary
make build-all      # Build all binaries
make run           # Run server
make dev           # Run with hot reload
```

### Testing
```bash
make test                # Run all tests
make test-unit          # Run unit tests
make test-integration   # Run integration tests
make test-coverage      # Run with coverage report
make benchmark          # Run benchmarks
```

### Code Quality
```bash
make lint          # Run linter
make format        # Format code
make vet           # Run go vet
make security      # Security scan
```

### Database
```bash
# Migration management (golang-migrate/migrate)
make migrate-up      # Apply all pending migrations
make migrate-down    # Rollback migrations
make migrate-create  # Create new migration
make migrate-status  # Show migration status
make migrate-version # Show current version

# See docs/DEVELOPMENT.md#database-migrations for detailed guide
# Or docs/MIGRATIONS.md for comprehensive migration documentation
```

### Docker
```bash
make docker-up           # Start all services
make docker-up-db        # Start database only
make docker-up-monitoring # Start monitoring stack
make docker-down         # Stop services
```

### Monitoring
```bash
make monitoring      # Start monitoring with URLs
make metrics        # View metrics
make health         # Health check
make trace-test     # Generate test traces
```

## 📊 Monitoring & Observability

### Metrics (Prometheus + Grafana)
- **HTTP Metrics**: Request duration, status codes, throughput
- **System Metrics**: CPU, memory, goroutines
- **Database Metrics**: Connection pool, query performance
- **Cache Metrics**: Hit/miss ratios, operation latency

### Distributed Tracing (Jaeger)
- **Request Tracing**: End-to-end request flow
- **Database Queries**: Query execution traces
- **External Calls**: Third-party service calls
- **Error Tracking**: Exception and error traces

### Logging (Zap)
- **Structured Logging**: JSON formatted logs
- **Log Levels**: Debug, Info, Warn, Error
- **Request Logging**: HTTP request/response logs
- **Error Logging**: Detailed error information

## 🚀 Deployment

### Docker Deployment
```bash
# Build Docker image
make docker-build

# Start production stack
docker-compose up -d
```

### Environment Configuration
- **Development**: `configs/development.yaml`
- **Production**: `configs/production.yaml`

## 🤝 Contributing

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Development Guidelines
- Follow **Clean Architecture** principles
- Write **comprehensive tests**
- Add **proper documentation**
- Use **conventional commits**
- Ensure **code quality** with linters

## 📚 Documentation

### Core Documentation
- [**📁 Project Structure**](docs/PROJECT_STRUCTURE.md) - Complete directory structure with detailed explanations
- [**🏛️ Architecture Guide**](docs/ARCHITECTURE.md) - Clean Architecture implementation and design patterns
- [**🛠️ Development Guide**](docs/DEVELOPMENT.md) - Development setup, testing, and best practices
- [**🚀 Deployment Guide**](docs/DEPLOYMENT.md) - Production deployment with Docker and Kubernetes

### API & Technical Guides
- [**📋 API Documentation**](docs/API.md) - REST API endpoints, examples, and usage
- [**� Email Service Guide**](docs/EMAIL_SERVICE.md) - Email service implementation and MailCatcher testing
- [**�🔌 Dependency Injection**](docs/DEPENDENCY_INJECTION.md) - Uber FX usage patterns and modules
- [**🔍 Tracing Guide**](docs/TRACING.md) - OpenTelemetry and Jaeger setup
- [**🛡️ RBAC Usage**](docs/RBAC_USAGE.md) - Role-based access control implementation

## 🛡️ Security

- **JWT Authentication**: Secure token-based auth
- **CORS Protection**: Cross-origin resource sharing
- **Rate Limiting**: Request throttling
- **Security Headers**: XSS, CSRF protection
- **Input Validation**: Comprehensive request validation
- **SQL Injection Prevention**: Parameterized queries

## 📈 Performance

- **Connection Pooling**: Database connection optimization
- **Caching Strategy**: Redis-based caching
- **Goroutine Management**: Efficient concurrent processing
- **Memory Management**: Optimized memory usage
- **Profiling**: Built-in performance profiling

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

## 🙋‍♂️ Support

For questions, issues, or contributions:

- **GitHub Issues**: [Create an issue](https://github.com/tranvuongduy2003/go-mvc/issues)
- **Email**: tranvuongduy2003@gmail.com
- **Documentation**: Check the [docs/](docs/) directory

---

⭐ **Star this repository** if you find it helpful!