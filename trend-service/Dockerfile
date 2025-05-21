# --- Stage 1: Build ---
    FROM golang:1.23-alpine AS builder

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
    
    WORKDIR /root/
    
    # Copy binary from build stage
    COPY --from=builder /app/app .
    
    # Expose port (optional, mostly for documentation)
    # EXPOSE 8080 (customize per service)
    
    CMD ["./app"]
    