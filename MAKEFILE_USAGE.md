# Makefile Usage Guide

This document provides a comprehensive guide on how to use the Makefile for the CSV Email Flagger project.

## Quick Start

```bash
# Show all available commands
make help

# Initial setup
make setup

# Install dependencies
make deps

# Run tests
make test

# Build the application
make build

# Run the application
make run
```

## Command Categories

### 🚀 **Setup & Dependencies**

| Command | Description |
|---------|-------------|
| `make setup` | Initial project setup (creates directories) |
| `make deps` | Install Go dependencies |
| `make help` | Show all available commands |

### 🧪 **Testing**

| Command | Description |
|---------|-------------|
| `make test` | Run all tests (unit + functional) |
| `make test-unit` | Run only unit tests |
| `make test-functional` | Run only functional tests |
| `make test-coverage` | Run tests with coverage report |
| `make test-verbose` | Run all tests with verbose output and race detection |
| `make quick-test` | Quick test run (short tests only) |

### 🔨 **Building**

| Command | Description |
|---------|-------------|
| `make build` | Build the application |
| `make build-linux` | Build for Linux |
| `make build-windows` | Build for Windows |
| `make build-all` | Build for all platforms |
| `make release` | Create a release build with all platforms |

### 🏃 **Running**

| Command | Description |
|---------|-------------|
| `make run` | Run in sequential mode (default) |
| `make run-parallel` | Run in parallel mode |
| `make run-build` | Build and run the application |
| `make quick-run` | Quick run without building |

### 🐳 **Docker**

| Command | Description |
|---------|-------------|
| `make docker-build` | Build Docker image |
| `make docker-run` | Run Docker container (sequential mode) |
| `make docker-run-parallel` | Run Docker container (parallel mode) |
| `make docker-stop` | Stop all running containers |

### 🛠️ **Development**

| Command | Description |
|---------|-------------|
| `make dev` | Start development server with auto-reload (requires air) |
| `make fmt` | Format Go code |
| `make lint` | Run linter (if available) |
| `make test-api` | Test API endpoints (requires server running) |

### 🧹 **Cleanup**

| Command | Description |
|---------|-------------|
| `make clean` | Clean build artifacts and temporary files |
| `make clean-storage` | Clean only storage files |

### 📊 **Monitoring & Status**

| Command | Description |
|---------|-------------|
| `make status` | Show project status |
| `make logs` | Show application logs (if running in background) |

### 📦 **Installation**

| Command | Description |
|---------|-------------|
| `make install` | Install the binary to GOPATH/bin |

## Common Workflows

### Development Workflow

```bash
# 1. Initial setup
make setup deps

# 2. Run tests to ensure everything works
make test

# 3. Start development server
make dev
# OR
make run
```

### Testing Workflow

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Quick test during development
make quick-test
```

### Build & Deploy Workflow

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Create release package
make release

# Docker deployment
make docker-build
make docker-run
```

### CI/CD Pipeline

```bash
# Setup
make setup deps

# Code quality
make fmt lint

# Testing
make test-coverage

# Building
make build-all

# Docker
make docker-build
```

## Configuration

### Environment Variables

You can override default values by setting environment variables:

```bash
# Change port
PORT=3000 make run

# Change process mode
PROCESS_MODE=parallel make run

# Change Docker tag
DOCKER_TAG=v1.0.0 make docker-build
```

### Custom Configuration

Edit the Makefile variables at the top:

```makefile
BINARY_NAME=csv-email-flagger
BUILD_DIR=build
DOCKER_IMAGE=csv-email-flagger
DOCKER_TAG=latest
PORT=8080
PROCESS_MODE=sequential
```

## Troubleshooting

### Common Issues

1. **Go version mismatch**
   ```bash
   # Check Go version
   go version
   
   # Update go.mod if needed
   # Edit go.mod to match your Go version
   ```

2. **Dependencies not found**
   ```bash
   make deps
   ```

3. **Build fails**
   ```bash
   make clean
   make deps
   make build
   ```

4. **Tests fail**
   ```bash
   make clean
   make test-verbose
   ```

5. **Docker build fails**
   ```bash
   # Check Docker is running
   docker --version
   
   # Clean Docker cache
   docker system prune
   
   # Rebuild
   make docker-build
   ```

### Getting Help

```bash
# Show all commands
make help

# Check project status
make status

# Run verbose tests for debugging
make test-verbose
```

## Examples

### Complete Development Setup

```bash
# Clone and setup
git clone <repository>
cd csv-email-flagger-structured
make setup deps

# Verify everything works
make test

# Start development
make run
```

### Production Deployment

```bash
# Build and test
make clean test build

# Docker deployment
make docker-build
make docker-run-parallel

# Or install locally
make install
```

### Release Process

```bash
# Full release build
make clean test build-all release

# Check release contents
ls -la release/
```

This Makefile provides a comprehensive set of commands for all aspects of development, testing, building, and deployment of the CSV Email Flagger project.
