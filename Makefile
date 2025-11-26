# Go MVC Makefile
.PHONY: help build run test clean docker-up docker-down migrate lint format setup dev monitoring mcp-build mcp-test mcp-clean

# Variables
APP_NAME=go-mvc
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT ?= $(shell git rev-parse HEAD)
BUILD_DIR=bin
DOCKER_COMPOSE=docker-compose
CONFIG_FILE ?= configs/development.yaml
MIGRATION_PATH=internal/adapters/persistence/postgres/migrations
DATABASE_URL ?= postgresql://postgres:postgres@localhost:5432/go_mvc_dev?sslmode=disable
MIGRATE_CMD=$(shell if command -v migrate >/dev/null 2>&1; then echo "migrate"; else echo "~/go/bin/migrate"; fi)

# MCP Variables
MCP_DIR=mcp
MCP_DIST=$(MCP_DIR)/dist

# Go build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"
CGO_ENABLED ?= 0
GOOS ?= darwin
GOARCH ?= amd64

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
NC=\033[0m # No Color

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ''
	@echo '$(YELLOW)Quick Start:$(NC)'
	@echo '  make setup           # Setup Go development environment'
	@echo '  make mcp-all         # Setup and build MCP agents'
	@echo '  make run             # Run the main application'
	@echo '  make dev             # Run with hot reload'
	@echo ''
	@echo '$(YELLOW)MCP Agents:$(NC)'
	@echo '  make mcp-status      # Check MCP agents status'
	@echo '  make mcp-test        # Test MCP agents'
	@echo '  make mcp-docs        # View documentation'

# ==========================================
# Build Commands
# ==========================================

build: ## Build main server binary
	@echo "$(YELLOW)Building $(APP_NAME) server...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(BUILD_DIR)/server cmd/main.go
	@echo "$(GREEN)Server built successfully$(NC)"

build-cli: ## Build CLI binary
	@echo "$(YELLOW)Building CLI...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(BUILD_DIR)/cli cmd/cli/main.go
	@echo "$(GREEN)CLI built successfully$(NC)"

build-worker: ## Build worker binary
	@echo "$(YELLOW)Building worker...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(BUILD_DIR)/worker cmd/worker/main.go
	@echo "$(GREEN)Worker built successfully$(NC)"

build-migrate: ## Build migration tool
	@echo "$(YELLOW)Building migration tool...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(BUILD_DIR)/migrate cmd/migrate/main.go
	@echo "$(GREEN)Migration tool built successfully$(NC)"

build-all: build build-cli build-worker build-migrate ## Build all binaries

# ==========================================
# Run Commands
# ==========================================

run: ## Run the main server
	@echo "$(YELLOW)Starting $(APP_NAME) server...$(NC)"
	@go run cmd/main.go

run-cli: ## Run the CLI
	@echo "$(YELLOW)Starting CLI...$(NC)"
	@go run cmd/cli/main.go

run-worker: ## Run the background worker
	@echo "$(YELLOW)Starting worker...$(NC)"
	@go run cmd/worker/main.go

dev: ## Run with hot reload using air
	@echo "$(YELLOW)Starting $(APP_NAME) with hot reload...$(NC)"
	@air -c .air.toml

# ==========================================
# Testing Commands
# ==========================================

test: ## Run all tests
	@echo "$(YELLOW)Running tests...$(NC)"
	@go test -v -race ./...

test-unit: ## Run unit tests only
	@echo "$(YELLOW)Running unit tests...$(NC)"
	@go test -v -race -short ./...

test-integration: ## Run integration tests only
	@echo "$(YELLOW)Running integration tests...$(NC)"
	@go test -v -race -tags=integration ./...

test-coverage: ## Run tests with coverage
	@echo "$(YELLOW)Running tests with coverage...$(NC)"
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

test-coverage-func: ## Show test coverage by function
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

benchmark: ## Run benchmarks
	@echo "$(YELLOW)Running benchmarks...$(NC)"
	@go test -bench=. -benchmem ./...

# ==========================================
# Code Quality Commands
# ==========================================

lint: ## Run golangci-lint
	@echo "$(YELLOW)Running linter...$(NC)"
	@golangci-lint run ./...

lint-fix: ## Run golangci-lint with auto-fix
	@echo "$(YELLOW)Running linter with auto-fix...$(NC)"
	@golangci-lint run ./... --fix

format: ## Format code
	@echo "$(YELLOW)Formatting code...$(NC)"
	@go fmt ./...
	@goimports -w .

vet: ## Run go vet
	@echo "$(YELLOW)Running go vet...$(NC)"
	@go vet ./...

security: ## Run gosec security scanner
	@echo "$(YELLOW)Running security scan...$(NC)"
	@gosec ./...

# ==========================================
# Dependency Management
# ==========================================

deps: ## Download and tidy dependencies
	@echo "$(YELLOW)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@go mod verify

deps-update: ## Update all dependencies
	@echo "$(YELLOW)Updating dependencies...$(NC)"
	@go get -u ./...
	@go mod tidy

deps-audit: ## Audit dependencies for vulnerabilities
	@echo "$(YELLOW)Auditing dependencies...$(NC)"
	@go list -json -deps ./... | nancy sleuth

# ==========================================
# Code Generation
# ==========================================

generate: ## Run go generate
	@echo "$(YELLOW)Generating code...$(NC)"
	@go generate ./...

mocks: ## Generate mocks (if mockgen is available)
	@echo "$(YELLOW)Generating mocks...$(NC)"
	@find . -name "*_mock.go" -delete
	@go generate ./...

swagger: ## Generate Swagger documentation
	@echo "$(YELLOW)Generating Swagger documentation...$(NC)"
	@swag init -g cmd/main.go -o api/openapi

# ==========================================
# Database Commands
# ==========================================

migrate-up: ## Run database migrations up
	@echo "$(YELLOW)Running migrations up...$(NC)"
	@$(MIGRATE_CMD) -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" up

migrate-down: ## Run database migrations down
	@echo "$(YELLOW)Running migrations down...$(NC)"
	@$(MIGRATE_CMD) -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" down

migrate-down-1: ## Run database migrations down by 1 step
	@echo "$(YELLOW)Running migration down by 1...$(NC)"
	@$(MIGRATE_CMD) -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" down 1

migrate-drop: ## Drop all migrations (DANGER!)
	@echo "$(RED)Dropping all migrations...$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		$(MIGRATE_CMD) -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" drop; \
	fi

migrate-force: ## Force migration version
	@read -p "Enter migration version: " version; \
	$(MIGRATE_CMD) -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" force $$version

migrate-create: ## Create new migration file
	@if [ -z "$(name)" ]; then \
		read -p "Enter migration name: " migration_name; \
	else \
		migration_name="$(name)"; \
	fi; \
	echo "$(YELLOW)Creating migration: $$migration_name$(NC)"; \
	$(MIGRATE_CMD) create -ext sql -dir $(MIGRATION_PATH) $$migration_name

migrate-version: ## Show current migration version
	@$(MIGRATE_CMD) -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" version

migrate-status: ## Show migration status
	@echo "$(YELLOW)Migration Status:$(NC)"
	@echo "Migration Path: $(MIGRATION_PATH)"
	@echo "Database URL: $(DATABASE_URL)"
	@echo -n "Current Version: "
	@$(MIGRATE_CMD) -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" version 2>/dev/null || echo "No migrations applied"
	@echo "Available Migrations:"
	@ls -la $(MIGRATION_PATH)/ 2>/dev/null | grep -E '\.(up|down)\.sql$$' | wc -l | xargs -I {} echo "  {} migration files found"

# ==========================================
# Docker Commands
# ==========================================

docker-build: ## Build Docker image
	@echo "$(YELLOW)Building Docker image...$(NC)"
	@docker build -t $(APP_NAME):$(VERSION) .
	@docker build -t $(APP_NAME):latest .

docker-up: ## Start all Docker services
	@echo "$(YELLOW)Starting Docker services...$(NC)"
	@$(DOCKER_COMPOSE) up -d

docker-up-db: ## Start only database services
	@echo "$(YELLOW)Starting database services...$(NC)"
	@$(DOCKER_COMPOSE) up -d postgres redis

docker-up-monitoring: ## Start monitoring stack
	@echo "$(YELLOW)Starting monitoring stack...$(NC)"
	@$(DOCKER_COMPOSE) up -d prometheus grafana jaeger

docker-down: ## Stop Docker services
	@echo "$(YELLOW)Stopping Docker services...$(NC)"
	@$(DOCKER_COMPOSE) down

docker-down-volumes: ## Stop Docker services and remove volumes
	@echo "$(RED)Stopping Docker services and removing volumes...$(NC)"
	@$(DOCKER_COMPOSE) down -v

docker-logs: ## View Docker logs
	@$(DOCKER_COMPOSE) logs -f

docker-ps: ## Show Docker container status
	@$(DOCKER_COMPOSE) ps

docker-restart: ## Restart Docker services
	@echo "$(YELLOW)Restarting Docker services...$(NC)"
	@$(DOCKER_COMPOSE) restart

docker-clean: ## Clean Docker system
	@echo "$(YELLOW)Cleaning Docker system...$(NC)"
	@docker system prune -f
	@docker volume prune -f

# ==========================================
# Monitoring & Observability
# ==========================================

monitoring: docker-up-monitoring ## Start monitoring stack and show URLs
	@echo "$(GREEN)Monitoring stack started!$(NC)"
	@echo "$(YELLOW)Prometheus:$(NC) http://localhost:9091"
	@echo "$(YELLOW)Grafana:$(NC) http://localhost:3000 (admin/admin)"
	@echo "$(YELLOW)Jaeger:$(NC) http://localhost:16686"

metrics: ## View application metrics
	@echo "$(YELLOW)Opening metrics endpoint...$(NC)"
	@curl -s http://localhost:8080/metrics | head -20

health: ## Check application health
	@echo "$(YELLOW)Checking application health...$(NC)"
	@curl -f http://localhost:8080/health || echo "$(RED)Health check failed$(NC)"

trace-test: ## Generate test traces
	@echo "$(YELLOW)Generating test traces...$(NC)"
	@for i in {1..5}; do \
		curl -s http://localhost:8080/api/v1/trace-test > /dev/null; \
		echo "Trace $$i sent"; \
		sleep 1; \
	done
	@echo "$(GREEN)Test traces generated$(NC)"

# ==========================================
# Development Setup
# ==========================================

setup: ## Setup development environment
	@echo "$(YELLOW)Setting up development environment...$(NC)"
	@go mod download
	@$(DOCKER_COMPOSE) up -d postgres redis
	@echo "$(YELLOW)Waiting for services to be ready...$(NC)"
	@sleep 10
	@make migrate-up
	@echo "$(GREEN)Development environment ready!$(NC)"
	@echo "$(GREEN)Run 'make run' or 'make dev' to start the application$(NC)"

setup-db: ## Setup database only
	@echo "$(YELLOW)Setting up database...$(NC)"
	@$(DOCKER_COMPOSE) up -d postgres
	@echo "$(YELLOW)Waiting for PostgreSQL to be ready...$(NC)"
	@sleep 10
	@make migrate-up
	@echo "$(GREEN)Database setup completed!$(NC)"

reset-db: ## Reset database (drop and recreate)
	@echo "$(RED)Resetting database...$(NC)"
	@read -p "This will delete all data. Are you sure? [y/N] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo ""; \
		make migrate-drop; \
		make migrate-up; \
		echo "$(GREEN)Database reset completed!$(NC)"; \
	else \
		echo ""; \
		echo "$(YELLOW)Database reset cancelled$(NC)"; \
	fi

setup-tools: ## Install development tools
	@echo "$(YELLOW)Installing development tools...$(NC)"
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "$(GREEN)Development tools installed successfully$(NC)"

# ==========================================
# Cleanup
# ==========================================

clean: ## Clean build artifacts and cache
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@go clean -cache -testcache
	@echo "$(GREEN)Cleanup completed$(NC)"

clean-all: clean docker-clean ## Clean everything including Docker
	@echo "$(GREEN)Full cleanup completed$(NC)"

# ==========================================
# Utility Commands
# ==========================================

version: ## Show version information
	@echo "App Name: $(APP_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Go Version: $(shell go version)"

status: ## Show project status
	@echo "$(YELLOW)Project Status:$(NC)"
	@echo "Git Branch: $(shell git branch --show-current 2>/dev/null || echo 'unknown')"
	@echo "Git Status: $(shell git status --porcelain 2>/dev/null | wc -l | xargs) uncommitted changes"
	@echo "Go Version: $(shell go version)"
	@echo "Docker Status:"
	@docker ps --format 'table {{.Names}}\t{{.Status}}' | grep -E "(postgres|redis|prometheus|grafana|jaeger)" || echo "  No containers running"

# ==========================================
# Performance & Profiling
# ==========================================

pprof-cpu: ## Profile CPU usage
	@echo "$(YELLOW)Profiling CPU usage...$(NC)"
	@go tool pprof http://localhost:8080/debug/pprof/profile

pprof-mem: ## Profile memory usage
	@echo "$(YELLOW)Profiling memory usage...$(NC)"
	@go tool pprof http://localhost:8080/debug/pprof/heap

pprof-goroutine: ## Profile goroutines
	@echo "$(YELLOW)Profiling goroutines...$(NC)"
	@go tool pprof http://localhost:8080/debug/pprof/goroutine

# ==========================================
# MCP Agents Commands
# ==========================================

mcp-setup: ## Setup MCP agents (install dependencies)
	@echo "$(YELLOW)Setting up MCP agents...$(NC)"
	@cd $(MCP_DIR) && npm install
	@echo "$(GREEN)MCP agents dependencies installed$(NC)"

mcp-build: ## Build MCP agents
	@echo "$(YELLOW)Building MCP agents...$(NC)"
	@cd $(MCP_DIR) && npm run build
	@echo "$(GREEN)MCP agents built successfully$(NC)"

mcp-test: ## Run MCP agents integration tests
	@echo "$(YELLOW)Testing MCP agents...$(NC)"
	@cd $(MCP_DIR) && node tests/integration/test-agents.js
	@echo "$(GREEN)MCP agents tests completed$(NC)"

mcp-dev: ## Watch and rebuild MCP agents on changes
	@echo "$(YELLOW)Starting MCP agents development mode...$(NC)"
	@cd $(MCP_DIR) && npm run dev

mcp-clean: ## Clean MCP agents build artifacts
	@echo "$(YELLOW)Cleaning MCP agents...$(NC)"
	@rm -rf $(MCP_DIST)
	@cd $(MCP_DIR) && rm -rf node_modules package-lock.json
	@echo "$(GREEN)MCP agents cleaned$(NC)"

mcp-status: ## Show MCP agents status
	@echo "$(YELLOW)MCP Agents Status:$(NC)"
	@echo "MCP Directory: $(MCP_DIR)"
	@if [ -d "$(MCP_DIST)" ]; then \
		echo "Build Status: $(GREEN)Built$(NC)"; \
		echo "API Agent: $(MCP_DIST)/agents/api/api-agent.server.js"; \
		echo "Database Agent: $(MCP_DIST)/agents/database/database-agent.server.js"; \
		ls -lh $(MCP_DIST)/agents/api/api-agent.server.js $(MCP_DIST)/agents/database/database-agent.server.js 2>/dev/null || echo "$(RED)Server files not found$(NC)"; \
	else \
		echo "Build Status: $(RED)Not built$(NC)"; \
		echo "Run 'make mcp-build' to build agents"; \
	fi
	@if [ -d "$(MCP_DIR)/node_modules" ]; then \
		echo "Dependencies: $(GREEN)Installed$(NC)"; \
	else \
		echo "Dependencies: $(RED)Not installed$(NC)"; \
		echo "Run 'make mcp-setup' to install dependencies"; \
	fi

mcp-docs: ## Show MCP documentation files
	@echo "$(YELLOW)MCP Documentation:$(NC)"
	@ls -lh $(MCP_DIR)/*.md 2>/dev/null | awk '{printf "  %-20s %s\n", $$9, $$5}' | sed 's|$(MCP_DIR)/||'
	@echo ""
	@echo "$(GREEN)Available documentation:$(NC)"
	@echo "  README.md           - User guide"
	@echo "  ARCHITECTURE.md     - Technical architecture"
	@echo "  QUICK_REFERENCE.md  - Developer quick start"
	@echo "  CHANGELOG.md        - Version history"
	@echo "  REFACTORING_COMPLETE.md - Refactoring summary"

mcp-config-show: ## Show MCP configuration example
	@echo "$(YELLOW)MCP Configuration Example:$(NC)"
	@echo "Add this to your MCP client config file:"
	@echo ""
	@echo "macOS: ~/Library/Application Support/Claude/claude_desktop_config.json"
	@echo "Windows: %APPDATA%\\Claude\\claude_desktop_config.json"
	@echo "Linux: ~/.config/Claude/claude_desktop_config.json"
	@echo ""
	@echo "{"
	@echo "  \"mcpServers\": {"
	@echo "    \"api-tester\": {"
	@echo "      \"command\": \"node\","
	@echo "      \"args\": [\"$(shell pwd)/$(MCP_DIST)/agents/api/api-agent.server.js\"]"
	@echo "    },"
	@echo "    \"database-agent\": {"
	@echo "      \"command\": \"node\","
	@echo "      \"args\": [\"$(shell pwd)/$(MCP_DIST)/agents/database/database-agent.server.js\"]"
	@echo "    }"
	@echo "  }"
	@echo "}"

mcp-all: mcp-setup mcp-build mcp-test ## Setup, build and test MCP agents

mcp-rebuild: mcp-clean mcp-setup mcp-build ## Clean rebuild MCP agents

# Default target
.DEFAULT_GOAL := help