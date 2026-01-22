# --- Stage 1: Build ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o session-service ./cmd/server

# --- Stage 2: Minimal Runtime ---
FROM alpine:latest

WORKDIR /root/

# Copy binary from build stage
COPY --from=builder /app/session-service .

# Expose port (optional, mostly for documentation)
EXPOSE 8080

CMD ["./session-service"]
