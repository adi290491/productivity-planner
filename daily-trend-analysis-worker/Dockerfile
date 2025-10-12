FROM golang:1.24-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

# Copy all source
COPY . .
    
# Build the Go binary
RUN go build -o daily-trend ./cmd

FROM alpine:latest

WORKDIR /root/
    
# Copy binary from build stage
COPY --from=builder /app/daily-trend .

CMD ["./daily-trend"]