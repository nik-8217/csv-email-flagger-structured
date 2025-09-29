# ---------- Build Stage ----------
FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o csv-email-flagger ./cmd/server

# ---------- Runtime Stage ----------
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/csv-email-flagger .
RUN mkdir -p /app/storage
EXPOSE 8080
ENV PROCESS_MODE=sequential
ENV PORT=8080
CMD ["./csv-email-flagger"]
