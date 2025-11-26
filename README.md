# Go MVC - Enterprise Grade Go Web Application

[![Go Version](https://img.shields.io/badge/go-1.24.5-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/docker-ready-blue.svg)](docker-compose.yml)
[![Clean Code](https://img.shields.io/badge/code-self--documenting-brightgreen.svg)](docs/book/appendix/AI_CODING_STANDARDS.md)

A modern, production-ready Go web application built with **Clean Architecture**, **Domain-Driven Design (DDD)**, and **Enterprise patterns**. Features comprehensive observability, AI-powered development, and self-documenting code standards.

## 🚀 Features

### 🤖 AI-Powered Development
- **AI Code Generation**: Auto-generate complete APIs from User Stories
- **Self-Documenting Code**: Minimal comments, maximum clarity through naming
- **AI Coding Standards**: Comprehensive rules for consistent, clean code
- **Production-Ready Output**: Generated code includes validation, error handling, tests
- **Architecture Compliance**: AI follows Clean Architecture patterns automatically

### Core Features
- **Clean Architecture**: Clear separation of Domain, Application, Infrastructure, Presentation layers
- **Domain-Driven Design**: Rich domain models with business logic encapsulation
- **Dependency Injection**: Uber FX for modular dependency management
- **CQRS Pattern**: Command Query Responsibility Segregation
- **Repository Pattern**: Abstract data access layer
- **Self-Documenting Code**: Clear naming eliminates need for comments

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
- **MCP Agents**: AI-powered API testing and database management agents

### 🤖 MCP Agents (Model Context Protocol)
- **API Testing Agent**: 6 tools for REST API testing (GET, POST, PUT, PATCH, DELETE, TEST)
- **Database Agent**: 6 tools for PostgreSQL management (connect, query, schema, migrate, analyze, generate SQL)
- **Modular Architecture**: Scalable, maintainable TypeScript implementation
- **Full Documentation**: Architecture, quick reference, and user guides included
- **Integration Ready**: Works with Claude Desktop and other MCP clients

## 📁 Project Structure

```
go-mvc/
├── cmd/                        # Application entry points
│   ├── main.go                # Main server application
│   ├── cli/                   # Command-line interface
│   ├── worker/                # Background worker
│   └── migrate/               # Database migration tool
├── internal/                   # Private application code
│   ├── domain/                # Domain Layer - Business logic & entities
│   │   ├── auth/             # Authentication domain
│   │   ├── user/             # User domain
│   │   ├── job/              # Background job domain
│   │   ├── messaging/        # Messaging domain
│   │   └── shared/           # Shared domain objects
│   ├── application/           # Application Layer - Use cases
│   │   ├── commands/         # Write operations (CQRS)
│   │   ├── queries/          # Read operations (CQRS)
│   │   ├── dto/              # Data transfer objects
│   │   ├── services/         # Application services
│   │   ├── validators/       # Input validation
│   │   └── event_handlers/   # Domain event handlers
│   ├── infrastructure/        # Infrastructure Layer - Technical details
│   │   ├── cache/            # Cache implementations (Redis)
│   │   ├── database/         # Database setup
│   │   ├── persistence/      # Repository implementations
│   │   ├── messaging/        # Message queue (NATS)
│   │   ├── jobs/             # Background job system
│   │   ├── security/         # Security utilities
│   │   ├── tracing/          # Distributed tracing
│   │   └── metrics/          # Prometheus metrics
│   ├── presentation/          # Presentation Layer - HTTP handlers
│   │   └── http/             # HTTP transport
│   │       ├── handlers/     # Request handlers
│   │       └── middleware/   # HTTP middleware
│   └── modules/               # Dependency injection modules
├── pkg/                       # Public reusable packages
│   ├── errors/               # Error utilities
│   ├── jwt/                  # JWT utilities
│   ├── pagination/           # Pagination helpers
│   ├── response/             # Response formatting
│   └── validator/            # Validation utilities
├── configs/                   # Configuration files
│   ├── development.yaml      # Development config
│   ├── production.yaml       # Production config
│   └── grafana/              # Grafana dashboards
├── docs/                      # Complete documentation (ebook format)
│   ├── BOOK.md               # Main table of contents
│   ├── INDEX.md              # Quick access index
│   ├── README.md             # Documentation overview
│   └── book/                 # Organized chapters
│       ├── 01-getting-started/
│       ├── 02-architecture/
│       ├── 03-development-guide/
│       ├── 04-features/
│       ├── 05-ai-development/
│       ├── 06-operations/
│       └── appendix/
└── scripts/                   # Build and deployment scripts
```

> 📖 **See [Project Structure Guide](docs/book/02-architecture/02-project-structure.md)** for detailed explanations

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

# Or run without hot reload
make run
```

### 4. Setup MCP Agents (Optional)
```bash
# Setup and build MCP agents
make mcp-all

# Test MCP agents
make mcp-test

# View MCP configuration
make mcp-config-show

# Check status
make mcp-status
```

> 📖 **See [MCP Agents Documentation](mcp/README.md)** for detailed usage

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

### 1. **Domain Layer** (`internal/domain/`)
- **Entities**: Core business objects with rich behavior
- **Value Objects**: Immutable data structures
- **Domain Events**: Business event definitions
- **Repository Interfaces**: Data access contracts
- **Domain Services**: Complex business logic

### 2. **Application Layer** (`internal/application/`)
- **Commands**: Write operations (CQRS)
- **Queries**: Read operations (CQRS)
- **Application Services**: Use case orchestration
- **DTOs**: Data transfer objects
- **Validators**: Input validation logic
- **Event Handlers**: Domain event processing

### 3. **Infrastructure Layer** (`internal/infrastructure/`)
- **Persistence**: Database repository implementations
- **Cache**: Redis caching layer
- **Messaging**: NATS message broker
- **Jobs**: Background job processing
- **Security**: Authentication & authorization
- **Tracing**: OpenTelemetry integration
- **Metrics**: Prometheus metrics

### 4. **Presentation Layer** (`internal/presentation/`)
- **HTTP Handlers**: REST API endpoints
- **Middleware**: Cross-cutting concerns (auth, logging, tracing)
- **Response Formatting**: Consistent API responses

### 5. **Modules** (`internal/modules/`)
- **Dependency Injection**: Uber FX modules for each domain
- **Lifecycle Management**: Startup/shutdown coordination

> 📖 **See [Architecture Guide](docs/book/02-architecture/01-architecture-overview.md)** for deep dive

## 🤖 AI-Powered API Generation

This project includes a comprehensive AI system that can automatically generate complete, production-ready APIs from User Stories following our strict coding standards.

### Quick AI Generation

1. **Use the User Story Template**
   ```bash
   # See template at: docs/USER_STORY_TEMPLATE.md
   ```

2. **Give AI This Instruction**
   ```markdown
   Generate a complete API following Clean Architecture from this User Story.
   Follow the rules in:
   - docs/book/appendix/AI_CODING_STANDARDS.md (self-documenting code)
   - docs/book/05-ai-development/02-api-generation-rules.md
   - docs/book/05-ai-development/03-code-generation-guidelines.md
   
   [Your User Story Here]
   ```

**AI will generate:**
- ✅ **Domain Layer**: Entities, value objects, repository interfaces
- ✅ **Application Layer**: Commands/queries, DTOs, validators, services  
- ✅ **Infrastructure Layer**: Repository implementations, migrations
- ✅ **Presentation Layer**: HTTP handlers, routes, middleware
- ✅ **Integration**: Complete dependency injection setup
- ✅ **Self-Documenting Code**: Clear names, no unnecessary comments

### AI Documentation
- **[🤖 AI Coding Standards](docs/book/appendix/AI_CODING_STANDARDS.md)** ⭐ **MUST READ** - Self-documenting code principles
- **[⚡ AI Quick Start](docs/book/05-ai-development/01-ai-quick-start.md)** - 5-minute tutorial
- **[📋 User Story Template](docs/USER_STORY_TEMPLATE.md)** - Complete template with examples
- **[⚙️ API Generation Rules](docs/book/05-ai-development/02-api-generation-rules.md)** - Comprehensive AI rules
- **[🔧 Code Generation Guidelines](docs/book/05-ai-development/03-code-generation-guidelines.md)** - Layer-by-layer guides

### Code Quality Principles

**Self-Documenting Code**
```go
// ❌ BAD - Unnecessary comments
// CreateUser creates a new user
func CreateUser(email, password string) (*User, error) {
    // Validate email
    if !isValid(email) {
        return nil, errors.New("invalid")
    }
    ...
}

// ✅ GOOD - Clear naming, no comments
func CreateUser(email, password string) (*User, error) {
    if err := ValidateEmail(email); err != nil {
        return nil, err
    }
    ...
}
```

> 📖 **See [AI Coding Standards](docs/book/appendix/AI_CODING_STANDARDS.md)** for complete guide

## 🔧 Development Commands

### Build & Run
```bash
make build          # Build server binary
make build-all      # Build all binaries
make run           # Run server
make dev           # Run with hot reload
```

### AI-Generated Code
```bash
make validate-generated [entity]  # Validate AI-generated code
make test-generated [entity]      # Test generated components
make integrate-generated [entity] # Integrate with existing codebase
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
# Migration management (golang-migrate/migrate v4.19.0+)
make migrate-up      # Apply all pending migrations
make migrate-down    # Rollback last migration
make migrate-create  # Create new timestamped migration
make migrate-status  # Show migration status
make migrate-version # Show current version

# See docs/book/03-development-guide/02-migrations.md for detailed guide
```

### Docker
```bash
make docker-up           # Start all services
make docker-up-db        # Start database only
make docker-up-monitoring # Start monitoring stack
make docker-down         # Stop services
```

### MCP Agents
```bash
make mcp-setup       # Install MCP dependencies
make mcp-build       # Build MCP agents
make mcp-test        # Run MCP tests
make mcp-all         # Setup, build, and test
make mcp-status      # Check MCP status
make mcp-docs        # View MCP documentation
make mcp-config-show # Show configuration example
make mcp-dev         # Watch mode (rebuild on changes)
make mcp-clean       # Clean MCP artifacts
make mcp-rebuild     # Clean rebuild
```

> 📖 **See [MCP Documentation](mcp/)** for detailed usage:
> - [README.md](mcp/README.md) - User guide
> - [ARCHITECTURE.md](mcp/ARCHITECTURE.md) - Technical docs
> - [QUICK_REFERENCE.md](mcp/QUICK_REFERENCE.md) - Quick start

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
- Use **self-documenting code** (see [AI Coding Standards](docs/book/appendix/AI_CODING_STANDARDS.md))
- Write **comprehensive tests**
- Add **proper documentation**
- Use **conventional commits**
- Ensure **code quality** with linters
- **NO unnecessary comments** - code should be clear through naming

## 📚 Documentation

### 📖 Complete Developer Guide (Ebook Format)
**[Start Here: Complete Book](docs/BOOK.md)** - Comprehensive documentation organized as professional ebook

**Quick Access:**
- [Quick Start Guide](docs/book/01-getting-started/02-quick-start.md) - Running in 15 minutes
- [Quick Reference](docs/book/01-getting-started/03-quick-reference.md) - Common commands & troubleshooting
- [Architecture Overview](docs/book/02-architecture/01-architecture-overview.md) - Clean Architecture deep dive
- [Development Workflow](docs/book/03-development-guide/01-development-workflow.md) - Dev environment setup

### 🤖 AI Development (Must Read for AI Assistants)
- **[AI Coding Standards](docs/book/appendix/AI_CODING_STANDARDS.md)** ⭐ **REQUIRED** - Self-documenting code principles
- [AI Quick Start](docs/book/05-ai-development/01-ai-quick-start.md) - AI-powered development intro
- [API Generation Rules](docs/book/05-ai-development/02-api-generation-rules.md) - Complete AI generation guide
- [Code Generation Guidelines](docs/book/05-ai-development/03-code-generation-guidelines.md) - Layer-by-layer templates

### 📋 Documentation Structure
```
docs/
├── BOOK.md                 # Main table of contents
├── INDEX.md                # Quick access index
├── README.md               # Documentation overview
├── MIGRATION_GUIDE.md      # Migration from old structure
└── book/                   # Organized content
    ├── 01-getting-started/     # Quick start, introduction
    ├── 02-architecture/        # Architecture guides
    ├── 03-development-guide/   # Development workflows
    ├── 04-features/            # Feature documentation
    ├── 05-ai-development/      # AI-powered development
    ├── 06-operations/          # Deployment & monitoring
    └── appendix/               # Standards, glossary, FAQ
        └── AI_CODING_STANDARDS.md  # ⭐ Essential coding standards
```

### 🔍 Key Topics
- **Getting Started**: [Introduction](docs/book/01-getting-started/01-introduction.md), [Quick Start](docs/book/01-getting-started/02-quick-start.md)
- **Architecture**: [Overview](docs/book/02-architecture/01-architecture-overview.md), [Domain Layer](docs/book/02-architecture/03-domain-layer.md), [DI](docs/book/02-architecture/07-dependency-injection.md)
- **Development**: [Workflow](docs/book/03-development-guide/01-development-workflow.md), [Migrations](docs/book/03-development-guide/02-migrations.md), [Testing](docs/book/03-development-guide/03-testing.md)
- **Features**: [Auth](docs/book/04-features/01-authentication.md), [Background Jobs](docs/book/04-features/02-background-jobs.md), [Email](docs/book/04-features/03-email-service.md), [Tracing](docs/book/04-features/07-tracing.md)
- **Operations**: [Deployment](docs/book/06-operations/01-deployment.md), [Monitoring](docs/book/06-operations/02-monitoring.md)

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
- **Documentation**: Start with [Complete Book](docs/BOOK.md) or [Quick Start](docs/book/01-getting-started/02-quick-start.md)
- **Email**: tranvuongduy2003@gmail.com

### 📖 Documentation Quick Links
- 🚀 [Quick Start (15 min)](docs/book/01-getting-started/02-quick-start.md)
- 🤖 [AI Coding Standards](docs/book/appendix/AI_CODING_STANDARDS.md) - **Essential for contributors**
- 📖 [Complete Developer Guide](docs/BOOK.md)
- 🔍 [Quick Reference](docs/book/01-getting-started/03-quick-reference.md)

---

⭐ **Star this repository** if you find it helpful!

**Built with ❤️ using Clean Architecture, DDD, and Self-Documenting Code principles**