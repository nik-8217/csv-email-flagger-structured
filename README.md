# CSV Email Flagger

A Go service to process CSV files and add a boolean flag if a row contains an email.

## Run locally
```bash
go mod tidy
go run ./cmd/server
```

## Run with Docker
```bash
docker build -t csv-email-flagger .
docker run -p 8080:8080 -e PROCESS_MODE=parallel csv-email-flagger
```

## API Endpoints
- POST /api/upload
- GET /api/status/{id}
- GET /api/download/{id}
- GET /swagger.json
- GET /healthz

## Tests
```bash
go test ./... -v
```
