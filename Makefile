# Makefile for Update Manager Project

.PHONY: help install build run stop test test-repo test-service test-backend test-api test-coverage test-service-coverage test-api-coverage test-pending-updates clean docker-build docker-up docker-down db-start db-stop db-status db-setup db-indexes db-logs load-test load-test-read load-test-write load-test-spike frontend-install frontend-dev frontend-stop frontend-build frontend-test-e2e frontend-test-e2e-ui frontend-test-e2e-headed frontend-install-browsers

# Variables
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_TEST=$(GO_CMD) test
GO_MOD=$(GO_CMD) mod
GO_FMT=$(GO_CMD) fmt
GO_VET=$(GO_CMD) vet

BACKEND_DIR=src/backend
FRONTEND_DIR=src/frontend
DB_DIR=src/database
BUILD_DIR=build

# Default target
help:
	@echo "Update Manager - Available commands:"
	@echo ""
	@echo "  make install       - Install Go dependencies"
	@echo "  make build         - Build backend binary"
	@echo "  make run           - Run backend server"
	@echo "  make stop          - Stop backend server"
	@echo "  make frontend-dev  - Run frontend development server"
	@echo "  make test          - Run all tests"
	@echo "  make test-repo      - Run repository tests"
	@echo "  make test-service   - Run service layer tests"
	@echo "  make test-backend   - Run all backend tests (repo + service)"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make fmt           - Format Go code"
	@echo "  make vet           - Run go vet"
	@echo "  make clean         - Clean build artifacts"
	@echo ""
	@echo "Load Testing:"
	@echo "  make load-test     - Run mixed load test"
	@echo "  make load-test-read - Run read-heavy load test"
	@echo "  make load-test-write - Run write-heavy load test"
	@echo "  make load-test-spike - Run spike test"
	@echo ""
	@echo "Database (Docker MongoDB):"
	@echo "  make db-start      - Start MongoDB Docker container"
	@echo "  make db-stop       - Stop MongoDB Docker container"
	@echo "  make db-status     - Check MongoDB container status"
	@echo "  make db-logs       - View MongoDB container logs"
	@echo "  make db-setup      - Setup database (runs automatically on first start)"
	@echo "  make db-indexes    - Create/recreate MongoDB indexes"
	@echo "  make db-cleanup    - Clean up database (interactive script)"
	@echo ""
	@echo "Frontend:"
	@echo "  make frontend-install      - Install frontend dependencies"
	@echo "  make frontend-dev          - Run frontend development server"
	@echo "  make frontend-stop         - Stop frontend development server"
	@echo "  make frontend-build        - Build frontend for production"
	@echo "  make frontend-test-e2e     - Run frontend E2E tests"
	@echo "  make frontend-test-e2e-ui  - Run frontend E2E tests with UI"
	@echo "  make frontend-test-e2e-headed - Run frontend E2E tests in headed mode"
	@echo "  make frontend-install-browsers - Install Playwright browsers"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build  - Build Docker images"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down   - Stop Docker containers"
	@echo ""

# Install dependencies
install:
	@echo "Installing Go dependencies..."
	cd $(BACKEND_DIR) && $(GO_MOD) download
	cd $(BACKEND_DIR) && $(GO_MOD) tidy

# Build backend
build:
	@echo "Building backend..."
	@mkdir -p $(BUILD_DIR)
	cd $(BACKEND_DIR) && $(GO_BUILD) -o ../../$(BUILD_DIR)/server ./cmd/server

# Run backend
run:
	@echo "Running backend server..."
	cd $(BACKEND_DIR) && $(GO_CMD) run ./cmd/server

# Stop backend server
stop:
	@echo "Stopping backend server..."
	@PORT=$${PORT:-8080}; \
	PIDS=$$(lsof -ti:$$PORT 2>/dev/null || true); \
	if [ -z "$$PIDS" ]; then \
		PIDS=$$(pgrep -f "go run.*cmd/server" 2>/dev/null || true); \
	fi; \
	if [ -z "$$PIDS" ]; then \
		PIDS=$$(pgrep -f "build.*server" 2>/dev/null || true); \
	fi; \
	if [ -n "$$PIDS" ]; then \
		echo "Found backend process(es): $$PIDS"; \
		echo "Sending SIGTERM signal..."; \
		kill -TERM $$PIDS 2>/dev/null || true; \
		sleep 2; \
		REMAINING=$$(lsof -ti:$$PORT 2>/dev/null || true); \
		if [ -n "$$REMAINING" ]; then \
			echo "Process still running, sending SIGKILL..."; \
			kill -KILL $$REMAINING 2>/dev/null || true; \
		fi; \
		echo "Backend server stopped successfully"; \
	else \
		echo "No backend server process found on port $$PORT or matching server pattern"; \
	fi

# Frontend commands
frontend-install:
	@echo "Installing frontend dependencies..."
	@echo "Note: Requires Node.js 18+. Using nvm if available..."
	@if command -v nvm >/dev/null 2>&1 || [ -s "$$HOME/.nvm/nvm.sh" ]; then \
		. $$HOME/.nvm/nvm.sh && nvm use 24.11.1 2>/dev/null || nvm use node 2>/dev/null || true; \
	fi
	cd $(FRONTEND_DIR) && npm install

frontend-dev:
	@echo "Starting frontend development server..."
	@echo "Frontend will be available at: http://localhost:3000"
	@echo "Note: Requires Node.js 18+. Using nvm if available..."
	@bash -c 'if [ -s "$$HOME/.nvm/nvm.sh" ]; then \
		source $$HOME/.nvm/nvm.sh && nvm use 24.11.1 2>/dev/null || nvm use node 2>/dev/null || true; \
	fi && cd $(FRONTEND_DIR) && npm run dev'

# Stop frontend development server
frontend-stop:
	@echo "Stopping frontend development server..."
	@PORT=$${FRONTEND_PORT:-3000}; \
	PIDS=$$(lsof -ti:$$PORT 2>/dev/null || true); \
	if [ -z "$$PIDS" ]; then \
		PIDS=$$(pgrep -f "vite" 2>/dev/null || true); \
	fi; \
	if [ -z "$$PIDS" ]; then \
		PIDS=$$(pgrep -f "npm run dev" 2>/dev/null || true); \
	fi; \
	if [ -n "$$PIDS" ]; then \
		echo "Found frontend process(es): $$PIDS"; \
		echo "Sending SIGTERM signal..."; \
		kill -TERM $$PIDS 2>/dev/null || true; \
		sleep 2; \
		REMAINING=$$(lsof -ti:$$PORT 2>/dev/null || true); \
		if [ -n "$$REMAINING" ]; then \
			echo "Process still running, sending SIGKILL..."; \
			kill -KILL $$REMAINING 2>/dev/null || true; \
		fi; \
		echo "Frontend server stopped successfully"; \
	else \
		echo "No frontend server process found on port $$PORT or matching vite/npm pattern"; \
	fi

frontend-build:
	@echo "Building frontend for production..."
	@echo "Note: Requires Node.js 18+. Using nvm if available..."
	@if command -v nvm >/dev/null 2>&1 || [ -s "$$HOME/.nvm/nvm.sh" ]; then \
		. $$HOME/.nvm/nvm.sh && nvm use 24.11.1 2>/dev/null || nvm use node 2>/dev/null || true; \
	fi
	cd $(FRONTEND_DIR) && npm run build
	@echo "Frontend build complete! Output in: $(FRONTEND_DIR)/dist"

# Frontend E2E tests
frontend-test-e2e:
	@echo "Running frontend E2E tests..."
	@echo "Note: Requires Node.js 18+. Using nvm if available..."
	@echo "Note: Frontend dev server will start automatically if not running"
	@if command -v nvm >/dev/null 2>&1 || [ -s "$$HOME/.nvm/nvm.sh" ]; then \
		. $$HOME/.nvm/nvm.sh && nvm use 24.11.1 2>/dev/null || nvm use node 2>/dev/null || true; \
	fi
	cd $(FRONTEND_DIR) && npm run test:e2e

frontend-test-e2e-ui:
	@echo "Running frontend E2E tests with UI..."
	@echo "Note: Requires Node.js 18+. Using nvm if available..."
	@echo "Note: Frontend dev server will start automatically if not running"
	@if command -v nvm >/dev/null 2>&1 || [ -s "$$HOME/.nvm/nvm.sh" ]; then \
		. $$HOME/.nvm/nvm.sh && nvm use 24.11.1 2>/dev/null || nvm use node 2>/dev/null || true; \
	fi
	cd $(FRONTEND_DIR) && npm run test:e2e:ui

frontend-test-e2e-headed:
	@echo "Running frontend E2E tests in headed mode (browser visible)..."
	@echo "Note: Requires Node.js 18+. Using nvm if available..."
	@echo "Note: Frontend dev server will start automatically if not running"
	@if command -v nvm >/dev/null 2>&1 || [ -s "$$HOME/.nvm/nvm.sh" ]; then \
		. $$HOME/.nvm/nvm.sh && nvm use 24.11.1 2>/dev/null || nvm use node 2>/dev/null || true; \
	fi
	cd $(FRONTEND_DIR) && npm run test:e2e:headed

frontend-install-browsers:
	@echo "Installing Playwright browsers and system dependencies..."
	@echo "Note: Requires Node.js 18+. Using nvm if available..."
	@if command -v nvm >/dev/null 2>&1 || [ -s "$$HOME/.nvm/nvm.sh" ]; then \
		. $$HOME/.nvm/nvm.sh && nvm use 24.11.1 2>/dev/null || nvm use node 2>/dev/null || true; \
	fi
	cd $(FRONTEND_DIR) && npx playwright install-deps
	cd $(FRONTEND_DIR) && npx playwright install
	@echo "Playwright browsers installed successfully!"

# Run tests
test:
	@echo "Running all tests..."
	cd $(BACKEND_DIR) && $(GO_TEST) -v ./...

# Run repository tests
test-repo:
	@echo "Running repository tests..."
	cd $(BACKEND_DIR) && $(GO_TEST) -v ./internal/repository

# Run service tests
test-service:
	@echo "Running service layer tests..."
	cd $(BACKEND_DIR) && $(GO_TEST) -v ./internal/service

# Run all backend tests (repository + service)
test-backend: test-repo test-service
	@echo "All backend tests completed!"

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	cd $(BACKEND_DIR) && $(GO_TEST) -v -coverprofile=coverage.out ./...
	@echo "Coverage report generated: coverage.out"

# Run service tests with coverage
test-service-coverage:
	@echo "Running service tests with coverage..."
	cd $(BACKEND_DIR) && $(GO_TEST) -v -coverprofile=service-coverage.out ./internal/service
	@echo "Service coverage report generated: service-coverage.out"

# Run API handler tests
test-api:
	@echo "Running API handler tests..."
	cd $(BACKEND_DIR) && $(GO_TEST) -v ./internal/api/handlers

# Run API tests with coverage
test-api-coverage:
	@echo "Running API tests with coverage..."
	cd $(BACKEND_DIR) && $(GO_TEST) -v -coverprofile=api-coverage.out ./internal/api/handlers
	@echo "API coverage report generated: api-coverage.out"

# Test pending updates API endpoints (requires backend running)
test-pending-updates:
	@echo "Testing Pending Updates API endpoints..."
	@echo "Make sure backend is running (make run)"
	@echo ""
	@echo "Testing deployment pending updates..."
	@curl -s http://localhost:8080/api/v1/customers/test-customer/tenants/test-tenant/deployments/test-deployment/updates | jq '.' || echo "Failed - check if test data exists and backend is running"
	@echo ""
	@echo "Testing tenant pending updates..."
	@curl -s http://localhost:8080/api/v1/customers/test-customer/tenants/test-tenant/deployments/pending-updates | jq '.' || echo "Failed - check if test data exists and backend is running"
	@echo ""
	@echo "Testing customer pending updates..."
	@curl -s http://localhost:8080/api/v1/customers/test-customer/deployments/pending-updates | jq '.' || echo "Failed - check if test data exists and backend is running"
	@echo ""
	@echo "Testing all pending updates (admin view)..."
	@curl -s "http://localhost:8080/api/v1/updates/pending?page=1&limit=20" | jq '.' || echo "Failed - check if test data exists and backend is running"

# Format code
fmt:
	@echo "Formatting Go code..."
	cd $(BACKEND_DIR) && $(GO_FMT) ./...

# Run go vet
vet:
	@echo "Running go vet..."
	cd $(BACKEND_DIR) && $(GO_VET) ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	cd $(BACKEND_DIR) && $(GO_CMD) clean

# MongoDB Docker commands
# Detect docker and compose commands
# Check docker socket access
HAS_DOCKER_ACCESS := $(shell test -r /var/run/docker.sock 2>/dev/null && echo "yes" || echo "no")
# Find docker path
DOCKER_PATH := $(shell which docker)
DOCKER_COMPOSE_PATH := $(shell which docker-compose)
# Find docker-compose or docker compose
DOCKER_COMPOSE_CMD := $(shell which docker-compose >/dev/null 2>&1 && echo "docker-compose" || (which docker >/dev/null 2>&1 && echo "docker compose" || echo "docker-compose"))
# Set DOCKER_CMD based on access
ifeq ($(HAS_DOCKER_ACCESS),yes)
	DOCKER_CMD := docker
	COMPOSE_CMD := $(DOCKER_COMPOSE_CMD)
else
	# Use full path for sudo
	DOCKER_CMD := sudo $(DOCKER_PATH)
	# For sudo, use full path or ensure PATH includes docker locations
	ifeq ($(DOCKER_COMPOSE_PATH),)
		COMPOSE_CMD := sudo env PATH=$$PATH $(DOCKER_COMPOSE_CMD)
	else
		COMPOSE_CMD := sudo $(DOCKER_COMPOSE_PATH)
	endif
endif

db-start:
	@echo "Starting MongoDB Docker container..."
	cd $(DB_DIR) && $(COMPOSE_CMD) -f docker-compose.mongodb.yml up -d
	@echo "Waiting for MongoDB to be ready..."
	@sleep 5
	@echo "MongoDB is ready!"
	@echo "Connection string: mongodb://admin:admin123@localhost:27017/updatemanager?authSource=admin"
	@echo "Mongo Express UI: http://localhost:8081 (admin/admin123)"

db-stop:
	@echo "Stopping MongoDB Docker container..."
	cd $(DB_DIR) && $(COMPOSE_CMD) -f docker-compose.mongodb.yml down
	@echo "MongoDB stopped!"

db-status:
	@echo "MongoDB Container Status:"
	@if $(DOCKER_CMD) ps --filter "name=updatemanager-mongodb" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" 2>/dev/null | grep -q updatemanager-mongodb; then \
		$(DOCKER_CMD) ps --filter "name=updatemanager-mongodb" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" 2>/dev/null; \
	else \
		echo "Container not running"; \
	fi

db-logs:
	@echo "MongoDB Container Logs:"
	cd $(DB_DIR) && $(COMPOSE_CMD) -f docker-compose.mongodb.yml logs -f mongodb

# Database setup (runs automatically on first container start, but can be run manually)
db-setup:
	@echo "Setting up MongoDB database..."
	@echo "Checking if MongoDB container is running..."
	@$(DOCKER_CMD) ps --filter "name=updatemanager-mongodb" --format "{{.Names}}" 2>/dev/null | grep -q updatemanager-mongodb || (echo "Error: MongoDB container is not running. Run 'make db-start' first." && exit 1)
	@echo "Running database setup script..."
	$(DOCKER_CMD) exec -i updatemanager-mongodb mongosh -u admin -p admin123 --authenticationDatabase admin updatemanager < $(DB_DIR)/setup-database.js
	@echo "Database setup complete!"

# Create database indexes
db-indexes:
	@echo "Creating MongoDB indexes..."
	@echo "Checking if MongoDB container is running..."
	@$(DOCKER_CMD) ps --filter "name=updatemanager-mongodb" --format "{{.Names}}" 2>/dev/null | grep -q updatemanager-mongodb || (echo "Error: MongoDB container is not running. Run 'make db-start' first." && exit 1)
	@echo "Running index creation script..."
	$(DOCKER_CMD) exec -i updatemanager-mongodb mongosh -u admin -p admin123 --authenticationDatabase admin updatemanager < $(DB_DIR)/mongodb-indexes.js
	@echo "Indexes created successfully!"

# Clean up database
db-cleanup:
	@echo "Running database cleanup script..."
	@bash scripts/cleanup-db.sh

# Docker commands
docker-build:
	@echo "Building Docker images..."
	docker-compose build

docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

# Development workflow
dev: install build run

# Load testing with Artillery
load-test:
	@echo "Running mixed load test..."
	@which artillery > /dev/null || (echo "Error: Artillery not found. Install with: npm install -g artillery" && exit 1)
	@node --version | grep -qE "^v(1[4-9]|[2-9][0-9])" || (echo "Error: Node.js 14+ required. Current version: $$(node --version). Please upgrade Node.js." && exit 1)
	@NODE_ENV=production artillery run load-tests/artillery-config.yml

load-test-read:
	@echo "Running read-heavy load test..."
	@which artillery > /dev/null || (echo "Error: Artillery not found. Install with: npm install -g artillery" && exit 1)
	@node --version | grep -qE "^v(1[4-9]|[2-9][0-9])" || (echo "Error: Node.js 14+ required. Current version: $$(node --version). Please upgrade Node.js." && exit 1)
	@NODE_ENV=production artillery run load-tests/artillery-read-heavy.yml

load-test-write:
	@echo "Running write-heavy load test..."
	@which artillery > /dev/null || (echo "Error: Artillery not found. Install with: npm install -g artillery" && exit 1)
	@node --version | grep -qE "^v(1[4-9]|[2-9][0-9])" || (echo "Error: Node.js 14+ required. Current version: $$(node --version). Please upgrade Node.js." && exit 1)
	@NODE_ENV=production artillery run load-tests/artillery-write-heavy.yml

load-test-spike:
	@echo "Running spike test..."
	@which artillery > /dev/null || (echo "Error: Artillery not found. Install with: npm install -g artillery" && exit 1)
	@node --version | grep -qE "^v(1[4-9]|[2-9][0-9])" || (echo "Error: Node.js 14+ required. Current version: $$(node --version). Please upgrade Node.js." && exit 1)
	@NODE_ENV=production artillery run load-tests/artillery-spike-test.yml

# Full setup
setup: install db-start
	@echo "Waiting for MongoDB to initialize..."
	@sleep 10
	@echo "Setup complete!"
	@echo ""
	@echo "MongoDB is running with:"
	@echo "  - Database: updatemanager"
	@echo "  - Root user: admin/admin123"
	@echo "  - App user: updatemanager/updatemanager123"
	@echo "  - Connection: mongodb://admin:admin123@localhost:27017/updatemanager?authSource=admin"
	@echo "  - Mongo Express: http://localhost:8081 (admin/admin123)"

