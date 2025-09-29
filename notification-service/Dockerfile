# --- Stage 1: Build ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source
COPY . .

# Build the Go binary
RUN go build -o notification-service ./cmd

# --- Stage 2: Minimal Runtime ---
FROM alpine:latest

WORKDIR /root/

# Copy binary from build stage
COPY --from=builder /app/notification-service .

# Expose port (optional, mostly for documentation)  
EXPOSE 8087

CMD ["./notification-service"]
