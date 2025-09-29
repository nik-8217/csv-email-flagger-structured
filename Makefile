# CSV Email Flagger Makefile
# This Makefile provides convenient commands for development, testing, and deployment

# Variables
BINARY_NAME=csv-email-flagger
BUILD_DIR=build
DOCKER_IMAGE=csv-email-flagger
DOCKER_TAG=latest
PORT=8080
PROCESS_MODE=sequential

# Go related variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: all build clean test test-unit test-functional test-coverage run run-parallel setup deps fmt lint docker-build docker-run docker-stop help

# Default target
all: clean deps fmt lint test build

# Help target
help: ## Show this help message
	@echo "$(BLUE)CSV Email Flagger - Available Commands:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""

# Setup and Dependencies
setup: ## Initial project setup
	@echo "$(BLUE)Setting up CSV Email Flagger project...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@mkdir -p storage
	@echo "$(GREEN)✓ Project directories created$(NC)"

deps: ## Install dependencies
	@echo "$(BLUE)Installing dependencies...$(NC)"
	$(GOMOD) tidy
	$(GOMOD) download
	@echo "$(GREEN)✓ Dependencies installed$(NC)"

# Code Quality
fmt: ## Format Go code
	@echo "$(BLUE)Formatting code...$(NC)"
	$(GOFMT) ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

lint: ## Run linter (if available)
	@echo "$(BLUE)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
		echo "$(GREEN)✓ Linting completed$(NC)"; \
	else \
		echo "$(YELLOW)⚠ golangci-lint not found, skipping linting$(NC)"; \
		echo "$(YELLOW)  Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
	fi

# Testing
test: test-unit test-functional ## Run all tests

test-unit: ## Run unit tests
	@echo "$(BLUE)Running unit tests...$(NC)"
	$(GOTEST) ./tests/unit/... -v
	@echo "$(GREEN)✓ Unit tests completed$(NC)"

test-functional: ## Run functional tests
	@echo "$(BLUE)Running functional tests...$(NC)"
	$(GOTEST) ./tests/functional/... -v
	@echo "$(GREEN)✓ Functional tests completed$(NC)"

test-coverage: ## Run tests with coverage report
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	$(GOTEST) ./... -coverprofile=coverage.out
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(NC)"

test-verbose: ## Run all tests with verbose output
	@echo "$(BLUE)Running all tests with verbose output...$(NC)"
	$(GOTEST) ./... -v -race
	@echo "$(GREEN)✓ All tests completed$(NC)"

# Building
build: ## Build the application
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server
	@echo "$(GREEN)✓ Build completed: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

build-linux: ## Build for Linux
	@echo "$(BLUE)Building for Linux...$(NC)"
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux ./cmd/server
	@echo "$(GREEN)✓ Linux build completed: $(BUILD_DIR)/$(BINARY_NAME)-linux$(NC)"

build-windows: ## Build for Windows
	@echo "$(BLUE)Building for Windows...$(NC)"
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows.exe ./cmd/server
	@echo "$(GREEN)✓ Windows build completed: $(BUILD_DIR)/$(BINARY_NAME)-windows.exe$(NC)"

build-all: build build-linux build-windows ## Build for all platforms

# Running
run: ## Run the application (sequential mode)
	@echo "$(BLUE)Starting CSV Email Flagger (sequential mode)...$(NC)"
	@echo "$(YELLOW)Server will be available at: http://localhost:$(PORT)$(NC)"
	@echo "$(YELLOW)Press Ctrl+C to stop$(NC)"
	PROCESS_MODE=$(PROCESS_MODE) PORT=$(PORT) $(GOCMD) run ./cmd/server

run-parallel: ## Run the application (parallel mode)
	@echo "$(BLUE)Starting CSV Email Flagger (parallel mode)...$(NC)"
	@echo "$(YELLOW)Server will be available at: http://localhost:$(PORT)$(NC)"
	@echo "$(YELLOW)Press Ctrl+C to stop$(NC)"
	PROCESS_MODE=parallel PORT=$(PORT) $(GOCMD) run ./cmd/server

run-build: build ## Build and run the application
	@echo "$(BLUE)Starting CSV Email Flagger from build...$(NC)"
	@echo "$(YELLOW)Server will be available at: http://localhost:$(PORT)$(NC)"
	@echo "$(YELLOW)Press Ctrl+C to stop$(NC)"
	PROCESS_MODE=$(PROCESS_MODE) PORT=$(PORT) ./$(BUILD_DIR)/$(BINARY_NAME)

# Docker
docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(NC)"
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "$(GREEN)✓ Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)$(NC)"

docker-run: ## Run Docker container
	@echo "$(BLUE)Starting Docker container...$(NC)"
	@echo "$(YELLOW)Server will be available at: http://localhost:$(PORT)$(NC)"
	@echo "$(YELLOW)Press Ctrl+C to stop$(NC)"
	docker run -p $(PORT):$(PORT) -e PORT=$(PORT) -e PROCESS_MODE=$(PROCESS_MODE) $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-run-parallel: ## Run Docker container in parallel mode
	@echo "$(BLUE)Starting Docker container (parallel mode)...$(NC)"
	@echo "$(YELLOW)Server will be available at: http://localhost:$(PORT)$(NC)"
	@echo "$(YELLOW)Press Ctrl+C to stop$(NC)"
	docker run -p $(PORT):$(PORT) -e PORT=$(PORT) -e PROCESS_MODE=parallel $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-stop: ## Stop all running containers
	@echo "$(BLUE)Stopping Docker containers...$(NC)"
	@docker ps -q --filter "ancestor=$(DOCKER_IMAGE):$(DOCKER_TAG)" | xargs -r docker stop
	@echo "$(GREEN)✓ Docker containers stopped$(NC)"

# Development
dev: ## Start development server with auto-reload (requires air)
	@echo "$(BLUE)Starting development server with auto-reload...$(NC)"
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "$(YELLOW)⚠ air not found, falling back to regular run$(NC)"; \
		echo "$(YELLOW)  Install with: go install github.com/cosmtrek/air@latest$(NC)"; \
		$(MAKE) run; \
	fi

# API Testing
test-api: ## Test API endpoints (requires server to be running)
	@echo "$(BLUE)Testing API endpoints...$(NC)"
	@echo "$(YELLOW)Make sure the server is running on port $(PORT)$(NC)"
	@echo ""
	@echo "$(BLUE)1. Health Check:$(NC)"
	@curl -s http://localhost:$(PORT)/healthz || echo "$(RED)✗ Health check failed$(NC)"
	@echo ""
	@echo "$(BLUE)2. Swagger JSON:$(NC)"
	@curl -s http://localhost:$(PORT)/swagger.json | head -c 100 || echo "$(RED)✗ Swagger endpoint failed$(NC)"
	@echo "..."
	@echo ""
	@echo "$(GREEN)✓ API tests completed$(NC)"

# File Management
clean: ## Clean build artifacts and temporary files
	@echo "$(BLUE)Cleaning up...$(NC)"
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	rm -rf storage/*.upload storage/*.csv
	@echo "$(GREEN)✓ Cleanup completed$(NC)"

clean-storage: ## Clean only storage files
	@echo "$(BLUE)Cleaning storage files...$(NC)"
	rm -rf storage/*.upload storage/*.csv
	@echo "$(GREEN)✓ Storage cleaned$(NC)"

# Installation
install: build ## Install the binary to GOPATH/bin
	@echo "$(BLUE)Installing $(BINARY_NAME)...$(NC)"
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/
	@echo "$(GREEN)✓ Installed to $(GOPATH)/bin/$(BINARY_NAME)$(NC)"

# Release
release: clean test build-all ## Create a release build
	@echo "$(BLUE)Creating release...$(NC)"
	@mkdir -p release
	@cp $(BUILD_DIR)/$(BINARY_NAME) release/
	@cp $(BUILD_DIR)/$(BINARY_NAME)-linux release/
	@cp $(BUILD_DIR)/$(BINARY_NAME)-windows.exe release/
	@cp README.md release/
	@cp Dockerfile release/
	@echo "$(GREEN)✓ Release created in release/ directory$(NC)"

# Monitoring
logs: ## Show application logs (if running in background)
	@echo "$(BLUE)Showing application logs...$(NC)"
	@if [ -f app.log ]; then \
		tail -f app.log; \
	else \
		echo "$(YELLOW)No log file found. Run with: make run > app.log 2>&1 &$(NC)"; \
	fi

# Quick Commands
quick-test: ## Quick test run
	@echo "$(BLUE)Running quick tests...$(NC)"
	$(GOTEST) ./... -short
	@echo "$(GREEN)✓ Quick tests completed$(NC)"

quick-run: ## Quick run without building
	@echo "$(BLUE)Quick run...$(NC)"
	$(GOCMD) run ./cmd/server

# Status
status: ## Show project status
	@echo "$(BLUE)CSV Email Flagger - Project Status$(NC)"
	@echo ""
	@echo "$(BLUE)Go Version:$(NC) $$(go version)"
	@echo "$(BLUE)Build Directory:$(NC) $(BUILD_DIR)/"
	@echo "$(BLUE)Storage Directory:$(NC) storage/"
	@echo "$(BLUE)Port:$(NC) $(PORT)"
	@echo "$(BLUE)Process Mode:$(NC) $(PROCESS_MODE)"
	@echo ""
	@if [ -d "$(BUILD_DIR)" ]; then \
		echo "$(GREEN)✓ Build directory exists$(NC)"; \
	else \
		echo "$(RED)✗ Build directory missing$(NC)"; \
	fi
	@if [ -d "storage" ]; then \
		echo "$(GREEN)✓ Storage directory exists$(NC)"; \
	else \
		echo "$(RED)✗ Storage directory missing$(NC)"; \
	fi
	@if [ -f "$(BUILD_DIR)/$(BINARY_NAME)" ]; then \
		echo "$(GREEN)✓ Binary built$(NC)"; \
	else \
		echo "$(YELLOW)⚠ Binary not built$(NC)"; \
	fi

# Default target when no argument is provided
.DEFAULT_GOAL := help
