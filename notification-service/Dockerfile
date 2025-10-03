# --- Stage 1: Build ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source
COPY . .

# Build the Go binary
RUN go build -o app ./cmd

# --- Stage 2: Minimal Runtime ---
FROM alpine:latest

# Install ca-certificates for HTTPS calls (needed for ADC)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from build stage
COPY --from=builder /app/app .

# Expose port (optional, mostly for documentation)
EXPOSE 8087

CMD ["./app"]
