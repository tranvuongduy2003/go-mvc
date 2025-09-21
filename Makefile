# Enterprise Makefile
.PHONY: help build run test clean docker-up docker-down migrate lint format setup

# Variables
APP_NAME=go-mvc
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT ?= $(shell git rev-parse HEAD)
BUILD_DIR=bin
DOCKER_COMPOSE=docker-compose
CONFIG_FILE ?= configs/development.yaml
MIGRATION_PATH=migrations
DATABASE_URL ?= postgresql://postgres:postgres@localhost:5432/enterprise_app_dev?sslmode=disable

# Go build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"
CGO_ENABLED ?= 0
GOOS ?= linux
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
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-20s$(NC) %s\n", $1, $2}' $(MAKEFILE_LIST)

# ==========================================
# Build Commands
# ==========================================

build: ## Build all binaries
	@echo "$(YELLOW)Building $(APP_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-api cmd/api/main.go
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-worker cmd/worker/main.go
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-migrate cmd/migrate/main.go
	@echo "$(GREEN)Build completed successfully$(NC)"

build-api: ## Build API server only
	@echo "$(YELLOW)Building API server...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-api cmd/api/main.go
	@echo "$(GREEN)API server built successfully$(NC)"

build-worker: ## Build worker only
	@echo "$(YELLOW)Building worker...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-worker cmd/worker/main.go
	@echo "$(GREEN)Worker built successfully$(NC)"

build-migrate: ## Build migration tool only
	@echo "$(YELLOW)Building migration tool...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-migrate cmd/migrate/main.go
	@echo "$(GREEN)Migration tool built successfully$(NC)"

# ==========================================
# Run Commands
# ==========================================

run: ## Run the API server
	@echo "$(YELLOW)Starting $(APP_NAME) API server...$(NC)"
	@go run cmd/api/main.go -config=$(CONFIG_FILE)

run-worker: ## Run the background worker
	@echo "$(YELLOW)Starting $(APP_NAME) worker...$(NC)"
	@go run cmd/worker/main.go -config=$(CONFIG_FILE)

dev: ## Run with hot reload using air
	@echo "$(YELLOW)Starting $(APP_NAME) with hot reload...$(NC)"
	@air -c .air.toml

dev-worker: ## Run worker with hot reload
	@echo "$(YELLOW)Starting worker with hot reload...$(NC)"
	@air -c .air.worker.toml

# ==========================================
# Testing Commands
# ==========================================

test: ## Run tests
	@echo "$(YELLOW)Running tests...$(NC)"
	@go test -v -race ./...

test-unit: ## Run unit tests only
	@echo "$(YELLOW)Running unit tests...$(NC)"
	@go test -v -race -short ./...

test-integration: ## Run integration tests only
	@echo "$(YELLOW)Running integration tests...$(NC)"
	@go test -v -race -run Integration ./tests/integration/...

test-e2e: ## Run end-to-end tests
	@echo "$(YELLOW)Running e2e tests...$(NC)"
	@go test -v -race ./tests/e2e/...

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
	@golangci-lint run ./... --config .golangci.yml

lint-fix: ## Run golangci-lint with auto-fix
	@echo "$(YELLOW)Running linter with auto-fix...$(NC)"
	@golangci-lint run ./... --config .golangci.yml --fix

format: ## Format code
	@echo "$(YELLOW)Formatting code...$(NC)"
	@go fmt ./...
	@goimports -w .
	@gofumpt -w .

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

mocks: ## Generate mocks
	@echo "$(YELLOW)Generating mocks...$(NC)"
	@mockgen -source=internal/core/ports/repositories/user.go -destination=tests/mocks/user_repository_mock.go
	@mockgen -source=internal/core/ports/services/user.go -destination=tests/mocks/user_service_mock.go

swagger: ## Generate Swagger documentation
	@echo "$(YELLOW)Generating Swagger documentation...$(NC)"
	@swag init -g cmd/api/main.go -o api/swagger

proto: ## Generate protobuf code
	@echo "$(YELLOW)Generating protobuf code...$(NC)"
	@protoc --go_out=. --go-grpc_out=. api/proto/*.proto

# ==========================================
# Database Commands
# ==========================================

migrate-up: ## Run database migrations up
	@echo "$(YELLOW)Running migrations up...$(NC)"
	@migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" up

migrate-down: ## Run database migrations down
	@echo "$(YELLOW)Running migrations down...$(NC)"
	@migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" down

migrate-drop: ## Drop all migrations
	@echo "$(RED)Dropping all migrations...$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	if [[ $REPLY =~ ^[Yy]$ ]]; then \
		migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" drop; \
	fi

migrate-force: ## Force migration version
	@read -p "Enter migration version: " version; \
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" force $version

migrate-create: ## Create new migration file
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATION_PATH) $name

migrate-version: ## Show current migration version
	@migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" version

# ==========================================
# Docker Commands
# ==========================================

docker-build: ## Build Docker images
	@echo "$(YELLOW)Building Docker images...$(NC)"
	@docker build -t $(APP_NAME):$(VERSION) .
	@docker build -t $(APP_NAME):latest .

docker-build-multi: ## Build multi-platform Docker images
	@echo "$(YELLOW)Building multi-platform Docker images...$(NC)"
	@docker buildx build --platform linux/amd64,linux/arm64 -t $(APP_NAME):$(VERSION) -t $(APP_NAME):latest --push .

docker-up: ## Start Docker services
	@echo "$(YELLOW)Starting Docker services...$(NC)"
	@$(DOCKER_COMPOSE) up -d

docker-up-build: ## Start Docker services with build
	@echo "$(YELLOW)Starting Docker services with build...$(NC)"
	@$(DOCKER_COMPOSE) up -d --build

docker-down: ## Stop Docker services
	@echo "$(YELLOW)Stopping Docker services...$(NC)"
	@$(DOCKER_COMPOSE) down

docker-down-volumes: ## Stop Docker services and remove volumes
	@echo "$(RED)Stopping Docker services and removing volumes...$(NC)"
	@$(DOCKER_COMPOSE) down -v

docker-logs: ## View Docker logs
	@$(DOCKER_COMPOSE) logs -f

docker-logs-api: ## View API Docker logs
	@$(DOCKER_COMPOSE) logs -f api

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
# Development Setup
# ==========================================

setup: ## Setup development environment
	@echo "$(YELLOW)Setting up development environment...$(NC)"
	@cp .env.example .env
	@go mod download
	@$(DOCKER_COMPOSE) up -d postgres redis
	@echo "$(YELLOW)Waiting for services to be ready...$(NC)"
	@sleep 10
	@make migrate-up
	@echo "$(GREEN)Development environment ready!$(NC)"
	@echo "$(GREEN)Run 'make dev' to start the application$(NC)"

setup-tools: ## Install development tools
	@echo "$(YELLOW)Installing development tools...$(NC)"
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install mvdan.cc/gofumpt@latest
	@go install github.com/securecodewarrior/go-tool-gosec/v2/cmd/gosec@latest
	@go install github.com/golang/mock/mockgen@latest
	@go install github.com/sonatard/noctx/cmd/noctx@latest
	@curl -sSfL https://raw.githubusercontent.com/nancy-cli/nancy/main/install.sh | sh -s -- -b $(go env GOPATH)/bin
	@echo "$(GREEN)Development tools installed successfully$(NC)"

# ==========================================
# Monitoring & Observability
# ==========================================

metrics: ## View Prometheus metrics
	@echo "Opening metrics endpoint..."
	@open http://localhost:9090/metrics || xdg-open http://localhost:9090/metrics

health: ## Check application health
	@echo "$(YELLOW)Checking application health...$(NC)"
	@curl -f http://localhost:8080/health || echo "$(RED)Health check failed$(NC)"

# ==========================================
# Deployment
# ==========================================

deploy-staging: ## Deploy to staging
	@echo "$(YELLOW)Deploying to staging...$(NC)"
	@./scripts/deploy.sh staging

deploy-prod: ## Deploy to production
	@echo "$(YELLOW)Deploying to production...$(NC)"
	@read -p "Are you sure you want to deploy to production? [y/N] " -n 1 -r; \
	if [[ $REPLY =~ ^[Yy]$ ]]; then \
		./scripts/deploy.sh production; \
	fi

k8s-apply: ## Apply Kubernetes manifests
	@echo "$(YELLOW)Applying Kubernetes manifests...$(NC)"
	@kubectl apply -f deployments/k8s/

k8s-delete: ## Delete Kubernetes resources
	@echo "$(YELLOW)Deleting Kubernetes resources...$(NC)"
	@kubectl delete -f deployments/k8s/

# ==========================================
# Cleanup
# ==========================================

clean: ## Clean build artifacts and cache
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@go clean -cache -modcache -testcache
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

env: ## Show environment variables
	@echo "$(YELLOW)Environment variables:$(NC)"
	@env | grep -E "(APP_|DB_|REDIS_|JWT_)" | sort

status: ## Show project status
	@echo "$(YELLOW)Project Status:$(NC)"
	@echo "Git Branch: $(shell git branch --show-current)"
	@echo "Git Status: $(shell git status --porcelain | wc -l) uncommitted changes"
	@echo "Go Version: $(shell go version)"
	@echo "Docker Status: $(shell docker ps --format 'table {{.Names}}\t{{.Status}}' | grep $(APP_NAME) || echo 'No containers running')"

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
# Load Testing
# ==========================================

load-test: ## Run load tests with hey
	@echo "$(YELLOW)Running load tests...$(NC)"
	@hey -n 1000 -c 50 http://localhost:8080/health

load-test-api: ## Run API load tests
	@echo "$(YELLOW)Running API load tests...$(NC)"
	@hey -n 1000 -c 50 -H "Content-Type: application/json" http://localhost:8080/api/v1/users

# Default target
.DEFAULT_GOAL := help# Makefile
.PHONY: help build run test clean docker-up docker-down migrate lint format

# Variables
APP_NAME=my-go-project
BUILD_DIR=bin
DOCKER_COMPOSE=docker-compose

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) cmd/api/main.go

run: ## Run the application
	@echo "Running $(APP_NAME)..."
	@go run cmd/api/main.go

dev: ## Run with hot reload using air
	@echo "Running $(APP_NAME) with hot reload..."
	@air

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

docker-up: ## Start docker services
	@echo "Starting Docker services..."
	@$(DOCKER_COMPOSE) up -d

docker-down: ## Stop docker services
	@echo "Stopping Docker services..."
	@$(DOCKER_COMPOSE) down

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):latest .

lint: ## Run golangci-lint
	@echo "Running linter..."
	@golangci-lint run ./...

format: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

generate: ## Generate code (if using go generate)
	@echo "Generating code..."
	@go generate ./...

# Database targets
migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	@migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/myapp_dev?sslmode=disable" up

migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	@migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/myapp_dev?sslmode=disable" down

migrate-create: ## Create new migration file (usage: make migrate-create name=create_users_table)
	@echo "Creating migration: $(name)"
	@migrate create -ext sql -dir migrations $(name)

# Development setup
setup: ## Setup development environment
	@echo "Setting up development environment..."
	@cp .env.example .env
	@go mod download
	@$(DOCKER_COMPOSE) up -d postgres
	@sleep 5
	@make migrate-up
	@echo "Development environment ready!"