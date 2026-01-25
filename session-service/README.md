# Session Service

A microservice for managing user focus/break/meeting sessions, built with Go.

## Project Structure

```
session-service/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/                     # Private application code
│   ├── config/                   # Configuration management
│   │   ├── config.go
│   │   └── config_test.go
│   ├── model/                    # Domain models
│   │   └── session.go
│   ├── repository/               # Data access layer
│   │   ├── repository.go         # Repository interface
│   │   ├── mock.go              # Mock for testing
│   │   └── postgres/
│   │       ├── session.go       # PostgreSQL implementation
│   │       └── session_test.go  # Integration tests
│   ├── service/                  # Business logic layer
│   │   ├── interface.go
│   │   ├── session.go
│   │   └── session_test.go
│   └── handler/                  # HTTP handlers
│       ├── session.go
│       ├── routes.go
│       └── handler_test.go
├── pkg/                          # Public libraries
│   └── httperr/                  # HTTP error utilities
│       └── error.go
├── migrations/                   # Database migrations
│   ├── 001_initial_schema.up.sql
│   ├── 001_initial_schema.down.sql
│   ├── 002_add_viewed_at.up.sql
│   └── 002_add_viewed_at.down.sql
├── testdata/                     # Test fixtures
│   └── schema.sql
├── Dockerfile
├── go.mod
└── go.sum
```

## Architecture

The service follows a clean architecture pattern with clear separation of concerns:

1. **Handler Layer** (`internal/handler`): HTTP request handling, input validation
2. **Service Layer** (`internal/service`): Business logic, orchestration
3. **Repository Layer** (`internal/repository`): Data access, persistence
4. **Model Layer** (`internal/model`): Domain entities

### Data Flow

```
HTTP Request → Handler → Service → Repository → Database
                  ↓
            (validation, transformation)
                  ↓
HTTP Response ← Handler ← Service ← Repository
```

## API Endpoints

### Health Checks
- `GET /health` - Service health check
- `GET /ready` - Readiness probe

### Session Operations
- `POST /sessions/v1/start-session` - Start a new session
- `PATCH /sessions/v1/stop-session` - Stop an active session

All session endpoints require `X-USER-ID` header (set by API gateway).

## Session Types

- `focus` - Focus/work session
- `break` - Break session
- `meeting` - Meeting session

## Configuration

Environment variables:

```bash
# Database
DB_HOSTNAME=<HOSTNAME>
DB_PORT=<PORT>
DB_NAME=<DATABASE>
DB_USERNAME=<USERNAME>
DB_PASSWORD=<PASSWORD>
DB_SSLMODE=<SSL_MODE>

# Server
PORT=8080
PROFILE=development
```

## Running the Service

### Local Development

```bash
# Install dependencies
go mod download

# Run the service
go run cmd/server/main.go
```

### Docker

```bash
# Build image
docker build -t session-service .

# Run container
docker run -p 8080:8080 \
  -e DB_HOSTNAME=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_NAME=productivity_planner \
  -e DB_USERNAME=postgres \
  -e DB_PASSWORD=postgres \
  session-service
```

## Testing

### Unit Tests

```bash
# Run all unit tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Integration Tests

Integration tests use dockertest to spin up a PostgreSQL container:

```bash
# Run integration tests (requires Docker)
go test -tags=integration ./internal/repository/postgres/...
```

### Test Coverage by Layer

- **Handler**: HTTP routing, request validation, response formatting
- **Service**: Business logic, UUID validation, error handling
- **Repository**: Database operations, transaction management
- **Model**: Domain entity validation

## Dependencies

- **Gin** - HTTP web framework (temporary, will be migrated to net/http)
- **GORM** - ORM (temporary, will be migrated to database/sql)
- **PostgreSQL** - Primary database
- **UUID** - Unique identifiers
- **godotenv** - Environment variable management
- **dockertest** - Integration testing with Docker

## Migration Plan

This service is in the process of being refactored:

1. ✅ **Phase 1**: Restructure code (COMPLETED)
2. 🔄 **Phase 2**: Migrate Gin → net/http (IN PROGRESS)
3. ⏳ **Phase 3**: Migrate GORM → database/sql (PENDING)

## Gateway Integration

This service sits behind a custom API gateway that handles:
- JWT authentication
- Rate limiting
- CORS
- Request routing

The gateway sets the `X-USER-ID` header for authenticated requests.
