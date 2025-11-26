# Go MVC - Enterprise Grade Go Web Application

[![Go Version](https://img.shields.io/badge/go-1.24.5-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/docker-ready-blue.svg)](docker-compose.yml)
[![Clean Code](https://img.shields.io/badge/code-self--documenting-brightgreen.svg)](docs/AI_CODING_STANDARDS.md)
[![CI/CD](https://img.shields.io/badge/CI%2FCD-GitHub%20Actions-2088FF.svg)](.github/workflows/)
[![AI Powered](https://img.shields.io/badge/AI-powered-ff69b4.svg)](docs/AI_DEVOPS.md)

A modern, production-ready Go web application built with **Clean Architecture**, **Domain-Driven Design (DDD)**, and **Enterprise patterns**. Features comprehensive observability, AI-powered development, automated CI/CD, and self-documenting code standards.

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
- **MCP Agents**: AI-powered API testing, database management, and GitHub automation agents

### 🤖 MCP Agents (Model Context Protocol)
- **API Testing Agent**: 6 tools for REST API testing (GET, POST, PUT, PATCH, DELETE, TEST)
- **Database Agent**: 6 tools for PostgreSQL management (connect, query, schema, migrate, analyze, generate SQL)
- **GitHub Agent**: 10 tools for GitHub automation (repo info, issues, PRs, workflows, code search, branches)
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

## 🚀 DevOps & CI/CD

### Automated CI/CD Pipeline

Complete DevOps automation với GitHub Actions và AI Assistant:

**6 Workflows Tự động**:
- ✅ **CI Pipeline**: Build, test, lint với PostgreSQL/Redis/NATS
- ✅ **Security**: Gosec, Trivy, CodeQL, secret scanning
- ✅ **Docker**: Multi-platform builds, GHCR publishing
- ✅ **Release**: Automated releases, binaries, deployment
- ✅ **Dependencies**: Auto-update Go modules, Actions, Docker
- ✅ **AI Assistant**: AI code review, docs, coverage analysis

**AI Development Scripts**:
```bash
# AI code review
./.github/scripts/ai-code-review.sh

# Generate CRUD code
./.github/scripts/ai-code-generator.sh Product full

# Analyze workflow
./.github/scripts/ai-workflow-optimizer.sh
```

**AI Bot Commands** (trong PR comments):
```bash
/ai review      # Trigger AI code review
/ai optimize    # Analyze workflow  
/ai docs        # Generate documentation
/ai help        # Show available commands
```

**Quick Start**:
```bash
# 1. Push code triggers CI/CD automatically
git push origin master

# 2. Create PR với AI assistance
gh pr create
# Comment: /ai review

# 3. Release với tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
# Auto: build binaries, create release, deploy
```

> 📖 **Complete DevOps Guide**:
> - [CI/CD Documentation](docs/CICD.md) - Workflows and automation
> - [AI DevOps](docs/AI_DEVOPS.md) - AI-powered development features
> - [MCP Agents](mcp/README.md) - AI testing and automation tools

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

> 📖 **See [AGENT.md](AGENT.md)** for comprehensive architecture guide and AI assistant instructions

## 🤖 AI-Powered API Generation

This project includes comprehensive AI assistance for development. The AGENT.md file contains complete instructions for AI assistants working on this codebase.

### Quick AI Generation

**AI will auto-generate complete, production-ready code following Clean Architecture:**

- ✅ **Domain Layer**: Entities, value objects, repository interfaces
- ✅ **Application Layer**: Commands/queries, DTOs, validators, services  
- ✅ **Infrastructure Layer**: Repository implementations, migrations
- ✅ **Presentation Layer**: HTTP handlers, routes, middleware
- ✅ **Integration**: Complete dependency injection setup
- ✅ **Self-Documenting Code**: Clear names, minimal comments

### AI Documentation & Standards
- **[AGENT.md](AGENT.md)** - Complete AI assistant guide
- **[AI Coding Standards](docs/AI_CODING_STANDARDS.md)** - Self-documenting code principles
- **[AI DevOps](docs/AI_DEVOPS.md)** - AI-powered development workflow

### Code Quality Principles

**Self-Documenting Code** - Code should be clear through naming, not comments:
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

// ✅ GOOD - Clear naming, no comments needed
func CreateUser(email, password string) (*User, error) {
    if err := ValidateEmail(email); err != nil {
        return nil, err
    }
    ...
}
```

> 📖 **See [AI Coding Standards](docs/AI_CODING_STANDARDS.md)** for complete self-documenting code guide

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

> 📖 **MCP Agent Documentation**:
> - [MCP README](mcp/README.md) - Complete user guide with all 3 agents
> - [GitHub Agent Setup](mcp/GITHUB_AGENT_SETUP.md) - GitHub automation guide
> - [MCP Architecture](mcp/ARCHITECTURE.md) - Technical architecture

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
- Use **self-documenting code** (see [AI Coding Standards](docs/AI_CODING_STANDARDS.md))
- Write **comprehensive tests**
- Add **proper documentation** to AGENT.md or README.md
- Use **conventional commits**
- Ensure **code quality** with linters
- **NO unnecessary comments** - code should be clear through naming
- **DO NOT create summary markdown files** after implementing features

## 📚 Documentation

### Core Documentation
- **[AGENT.md](AGENT.md)** - Complete AI assistant guide and architecture
- **[README.md](README.md)** - This file, project overview and quick start
- **[CI/CD Guide](docs/CICD.md)** - Complete CI/CD automation documentation
- **[AI DevOps](docs/AI_DEVOPS.md)** - AI-powered development workflow
- **[AI Coding Standards](docs/AI_CODING_STANDARDS.md)** - Self-documenting code principles

### MCP Agents
- **[MCP README](mcp/README.md)** - API Testing, Database, and GitHub agents
- **[GitHub Agent Setup](mcp/GITHUB_AGENT_SETUP.md)** - GitHub automation setup
- **[MCP Architecture](mcp/ARCHITECTURE.md)** - Technical architecture details
- **[MCP Changelog](mcp/CHANGELOG.md)** - Version history and changes

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
- **Documentation**: See [AGENT.md](AGENT.md) for comprehensive guide
- **Email**: tranvuongduy2003@gmail.com

### Quick Documentation Links
- 🤖 [AGENT.md](AGENT.md) - AI assistant complete guide
- 📖 [README.md](README.md) - Project overview (this file)
- 🚀 [CI/CD Guide](docs/CICD.md) - DevOps automation
- 🤖 [AI DevOps](docs/AI_DEVOPS.md) - AI development workflow
- 📝 [AI Coding Standards](docs/AI_CODING_STANDARDS.md) - Code quality principles
- 🔧 [MCP Agents](mcp/README.md) - Testing and automation tools

---

⭐ **Star this repository** if you find it helpful!

**Built with ❤️ using Clean Architecture, DDD, and Self-Documenting Code principles**