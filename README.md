# CSV Email Flagger

A high-performance Go service that processes CSV files and adds a boolean `hasEmail` flag to each row based on whether the row contains a valid email address. The service supports both sequential and parallel processing modes for optimal performance.

## Quick Start

```bash
# Clone and setup
git clone <repository-url>
cd csv-email-flagger-structured

# Quick setup and run
make setup deps test run

# Or run in parallel mode
make run-parallel
```

The service will be available at `http://localhost:8080`

## Features

- **Email Detection**: Uses regex pattern matching to identify valid email addresses in CSV data
- **Flexible Processing**: Supports both sequential and parallel processing modes
- **Robust Error Handling**: Handles blank rows, malformed CSV, and various edge cases
- **Header Management**: Automatically adds `hasEmail` column header if not present
- **File Cleanup**: Built-in cleanup mechanisms for temporary files
- **RESTful API**: Clean HTTP API for file upload, status checking, and download
- **Comprehensive Testing**: Extensive unit and functional test coverage
- **Makefile Support**: Complete build automation and development workflow

## Architecture

### Design Overview

The service follows a clean architecture pattern with clear separation of concerns:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP API      │    │   Job Manager   │    │   Transform     │
│   (handlers)    │───▶│   (processor)   │───▶│   (parallel/    │
│                 │    │                 │    │    sequential)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Router        │    │   Job Store     │    │   Storage       │
│   (mux)         │    │   (in-memory)   │    │   (file system) │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Components

1. **API Layer** (`internal/api/`): HTTP handlers and routing
2. **Job Management** (`internal/jobs/`): Job lifecycle and status tracking
3. **Transform Engine** (`internal/transform/`): CSV processing logic
4. **Storage** (`internal/storage/`): File management and cleanup
5. **Logging** (`pkg/logger/`): Structured logging with logrus

### Processing Modes

- **Sequential**: Processes CSV row by row (default)
- **Parallel**: Uses worker goroutines for concurrent processing (configurable worker count)

## Installation & Setup

### Prerequisites

- Go go 1.22 or later
- Docker (optional, for containerized deployment)
- Make (optional, for using Makefile commands)

### Local Development

#### Using Makefile (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd csv-email-flagger-structured
   ```

2. **Quick setup and run**
   ```bash
   # Show all available commands
   make help
   
   # Initial setup
   make setup deps
   
   # Run tests
   make test
   
   # Run the service
   make run
   ```

3. **Development workflow**
   ```bash
   # Run in parallel mode
   make run-parallel
   
   # Run with auto-reload (requires air)
   make dev
   
   # Build and run
   make build run-build
   ```

#### Manual Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd csv-email-flagger-structured
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Run the service**
   ```bash
   go run ./cmd/server
   ```

4. **Run with specific processing mode**
   ```bash
   PROCESS_MODE=parallel go run ./cmd/server
   ```

### Docker Deployment

#### Using Makefile (Recommended)

1. **Build and run with Docker**
   ```bash
   # Build Docker image
   make docker-build
   
   # Run container (sequential mode)
   make docker-run
   
   # Run container (parallel mode)
   make docker-run-parallel
   
   # Stop containers
   make docker-stop
   ```

#### Manual Docker Commands

1. **Build the image**
   ```bash
   docker build -t csv-email-flagger .
   ```

2. **Run the container**
   ```bash
   # Sequential mode (default)
   docker run -p 8080:8080 csv-email-flagger
   
   # Parallel mode
   docker run -p 8080:8080 -e PROCESS_MODE=parallel csv-email-flagger
   ```

3. **Run with custom port**
   ```bash
   docker run -p 3000:3000 -e PORT=3000 csv-email-flagger
   ```

## API Reference

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/upload` | Upload and process a CSV file |
| `GET` | `/api/status/{id}` | Get job status by ID |
| `GET` | `/api/download/{id}` | Download processed CSV file |
| `POST` | `/api/cleanup` | Clean up old temporary files |
| `GET` | `/healthz` | Health check endpoint |
| `GET` | `/swagger.json` | OpenAPI specification |

### Request/Response Examples

#### Upload CSV File
```bash
curl -X POST -F "file=@data.csv" http://localhost:8080/api/upload
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "mode": "parallel"
}
```

#### Check Job Status
```bash
curl http://localhost:8080/api/status/550e8400-e29b-41d4-a716-446655440000
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "DONE",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:05Z",
  "mode": "parallel"
}
```

#### Download Processed File
```bash
curl -O http://localhost:8080/api/download/550e8400-e29b-41d4-a716-446655440000
```

#### Cleanup Old Files
```bash
curl -X POST http://localhost:8080/api/cleanup
```

**Response:**
```json
{
  "message": "cleanup completed"
}
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `PROCESS_MODE` | `sequential` | Processing mode (`sequential` or `parallel`) |

### Processing Modes

- **Sequential Mode**: Processes CSV files row by row, suitable for smaller files
- **Parallel Mode**: Uses multiple worker goroutines for concurrent processing, ideal for large files

## Testing

### Using Makefile (Recommended)

1. **Run all tests**
   ```bash
   make test
   ```

2. **Run specific test types**
   ```bash
   # Unit tests only
   make test-unit
   
   # Functional tests only
   make test-functional
   
   # Tests with coverage
   make test-coverage
   
   # Verbose tests with race detection
   make test-verbose
   ```

3. **Quick testing during development**
   ```bash
   make quick-test
   ```

### Manual Testing

1. **Run all tests**
   ```bash
   go test ./... -v
   ```

2. **Run specific test packages**
   ```bash
   # Unit tests only
   go test ./tests/unit/... -v
   
   # Functional tests only
   go test ./tests/functional/... -v
   ```

3. **Run tests with coverage**
   ```bash
   go test ./... -v -cover
   ```

4. **Generate coverage report**
   ```bash
   go test ./... -coverprofile=coverage.out
   go tool cover -html=coverage.out -o coverage.html
   ```

### Test Categories

#### Unit Tests (`tests/unit/`)
- Transform function validation
- Email regex pattern testing
- Edge case handling (blank rows, malformed CSV)
- Header management logic

#### Functional Tests (`tests/functional/`)
- End-to-end API testing
- File upload and processing workflows
- Error handling scenarios
- Job status management

### Test Data

The test suite includes various CSV samples:
- Standard CSV with valid emails
- CSV with blank/incomplete rows
- CSV with existing `hasEmail` header
- Malformed CSV data
- Empty files
- Large files for performance testing

## Error Handling

The service handles various error scenarios gracefully:

- **Empty CSV files**: Returns appropriate error message
- **Malformed CSV**: Logs error and fails job processing
- **Missing files**: Returns 400 Bad Request
- **Invalid job IDs**: Returns 400 Bad Request
- **File I/O errors**: Logs error and updates job status
- **Processing errors**: Cleans up temporary files and reports failure

## File Management

### Storage Structure
```
storage/
├── {job-id}.upload    # Original uploaded file
└── {job-id}.csv       # Processed output file
```

### Cleanup Mechanisms

1. **Automatic Cleanup**: Files older than 24 hours are automatically cleaned up
2. **Manual Cleanup**: Use `/api/cleanup` endpoint for immediate cleanup
3. **Error Cleanup**: Failed jobs automatically clean up their temporary files

## Performance Considerations

### Sequential vs Parallel Processing

- **Sequential**: Lower memory usage, suitable for small to medium files
- **Parallel**: Higher throughput for large files, configurable worker count

### Memory Management

- Streaming CSV processing to handle large files
- Proper file handle management with defer statements
- Automatic cleanup of temporary files

### Scalability

- In-memory job store (consider external storage for production)
- Stateless processing (can be horizontally scaled)
- Configurable worker pools for parallel processing

## Development

### Project Structure
```
csv-email-flagger-structured/
├── cmd/server/           # Application entry point
├── internal/
│   ├── api/             # HTTP handlers and routing
│   ├── jobs/            # Job management
│   ├── storage/         # File storage utilities
│   └── transform/       # CSV processing logic
├── pkg/logger/          # Logging utilities
├── tests/
│   ├── unit/            # Unit tests
│   └── functional/      # Integration tests
├── storage/             # File storage directory
├── build/               # Build artifacts (created by make)
├── Makefile             # Build automation and development commands
├── MAKEFILE_USAGE.md    # Detailed Makefile usage guide
├── Dockerfile           # Docker container configuration
├── go.mod               # Go module dependencies
├── go.sum               # Go module checksums
└── README.md            # This file
```

### Makefile Commands

The project includes a comprehensive Makefile for development, testing, and deployment:

#### Quick Commands
```bash
make help          # Show all available commands
make setup         # Initial project setup
make deps          # Install dependencies
make test          # Run all tests
make build         # Build the application
make run           # Run the service
make clean         # Clean build artifacts
```

#### Development Commands
```bash
make dev           # Development server with auto-reload
make fmt           # Format Go code
make lint          # Run linter
make test-coverage # Run tests with coverage report
make quick-test    # Quick test run
```

#### Docker Commands
```bash
make docker-build        # Build Docker image
make docker-run          # Run Docker container
make docker-run-parallel # Run in parallel mode
make docker-stop         # Stop containers
```

#### Build Commands
```bash
make build-all     # Build for all platforms
make release       # Create release package
make install       # Install binary to GOPATH/bin
```

For a complete list of commands and detailed usage, see [MAKEFILE_USAGE.md](MAKEFILE_USAGE.md).

### Adding New Features

1. **New API endpoints**: Add handlers in `internal/api/handlers.go`
2. **New processing modes**: Extend `internal/transform/` package
3. **New storage backends**: Implement interfaces in `internal/storage/`
4. **New test cases**: Add to appropriate test package

### Code Style

- Follow Go standard formatting (`gofmt`)
- Use meaningful variable and function names
- Add comprehensive error handling
- Include unit tests for new functionality
- Document public APIs

## Troubleshooting

### Common Issues

1. **Port already in use**
   ```bash
   # Using Makefile
   PORT=3000 make run
   
   # Manual
   PORT=3000 go run ./cmd/server
   ```

2. **Permission denied on storage directory**
   ```bash
   # Using Makefile
   make setup
   
   # Manual
   chmod 755 storage/
   ```

3. **Memory issues with large files**
   ```bash
   # Using Makefile
   make run-parallel
   
   # Manual
   PROCESS_MODE=parallel go run ./cmd/server
   ```

4. **Go version mismatch**
   ```bash
   # Check Go version
   go version
   
   # Update go.mod if needed (currently set to go 1.22)
   ```

5. **Makefile not found or not working**
   ```bash
   # Check if make is installed
   make --version
   
   # Use manual commands instead
   go mod tidy
   go test ./...
   go run ./cmd/server
   ```

6. **Docker build fails**
   ```bash
   # Check Docker is running
   docker --version
   
   # Clean Docker cache
   docker system prune
   
   # Rebuild
   make docker-build
   ```

### Logging

The service uses structured logging with different levels:
- `INFO`: General operational messages
- `WARN`: Non-critical issues
- `ERROR`: Processing errors
- `FATAL`: Critical errors that stop the service

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

[Add your license information here]
