# Notification Service Testing Guide

This document describes the comprehensive testing setup for the notification service, including unit tests, integration tests, and coverage analysis.

## Test Structure

### Unit Tests
- **Handler Tests** (`cmd/handler_test.go`): Tests for HTTP request handlers with authentication and validation
- **Service Tests** (`notification/notification_service_test.go`): Tests for business logic and database operations  
- **Models Tests** (`models/models_test.go`): Tests for data models and JSON marshaling/unmarshaling

### Integration Tests
- **Subscription Integration** (`notification/integration_test.go`): Tests for PubSub subscription functionality
- **Database Integration**: Tests for complete database workflows with PostgreSQL

## Running Tests

### Basic Unit Tests
```bash
# Run all unit tests
go test ./cmd ./notification ./models

# Run with verbose output
go test -v ./cmd ./notification ./models

# Run specific test file
go test -v ./cmd/handler_test.go
```

### Coverage Analysis
```bash
# Run tests with coverage
go test -coverprofile=coverage.out ./cmd ./notification ./models

# View coverage summary
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
```

### Integration Tests
```bash
# Set environment variables for integration tests
export RUN_INTEGRATION_TESTS=true
export TEST_DATABASE_URL="postgres://postgres:password@localhost:5432/test_notifications?sslmode=disable"

# Run integration tests
go test -v ./notification/integration_test.go
```

### Using Test Script
```bash
# Make script executable
chmod +x run_tests.sh

# Run all tests with coverage
./run_tests.sh

# Run with integration tests enabled
RUN_INTEGRATION_TESTS=true ./run_tests.sh
```

## Test Coverage Goals

The test suite aims for **>70% code coverage** across:

### Handler Layer (cmd/handler_test.go)
✅ **Covered Scenarios:**
- Authentication enforcement (401 errors)
- Authorization checks (403 errors) 
- Input validation (400 errors)
- Success responses (200 OK)
- Service error handling (500 errors)
- Case-insensitive notification type validation
- UUID parsing and validation

### Service Layer (notification/notification_service_test.go)
✅ **Covered Scenarios:**
- Database CRUD operations
- User notification retrieval (exists/not found)
- Notification marking as read (daily/weekly)
- Invalid notification type handling
- Database UPSERT operations
- Success/failure event processing
- Statistics tracking

### Models Layer (models/models_test.go)
✅ **Covered Scenarios:**
- JSON marshaling/unmarshaling for all models
- Field assignments and data integrity
- camelCase JSON field naming
- Error event structures
- Success event structures with user arrays

### Integration Layer (notification/integration_test.go)
✅ **Covered Scenarios:**
- PubSub subscription setup
- End-to-end message processing
- Database persistence verification
- Error handling workflows
- Multi-user notification scenarios

## Test Database Setup

### PostgreSQL Test Database
```sql
-- Create test database
CREATE DATABASE test_notifications;

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE test_notifications TO postgres;
```

### Environment Variables
```bash
# Test database connection
export TEST_DATABASE_URL="postgres://postgres:password@localhost:5432/test_notifications?sslmode=disable"

# Integration test flag  
export RUN_INTEGRATION_TESTS=true

# PubSub emulator (if testing with emulator)
export PUBSUB_EMULATOR_HOST=localhost:8085
```

## Mock Strategy

The tests use **interface-based mocking** for clean unit testing:

### NotificationServiceInterface
```go
type NotificationServiceInterface interface {
    ProcessDailyTrendNotifications(ctx context.Context) error
    ProcessWeeklyTrendNotifications(ctx context.Context) error
    GetUserNotification(userID uuid.UUID) (*models.UserNotificationResponse, error)
    MarkNotificationAsRead(userID uuid.UUID, notificationType string) error
    GetStats() notification.ProcessingStats
}
```

### MockNotificationService
Implements the interface for handler testing without requiring actual database or PubSub connections.

## Key Test Features

### Authentication Testing
- JWT token extraction from context
- User ID validation and parsing
- Authorization checks (users can only access their own data)
- Multiple authentication error scenarios

### Input Validation Testing
- UUID format validation
- Notification type normalization (case-insensitive)
- Required field validation
- Boundary condition testing

### Database Integration Testing  
- PostgreSQL connection with test database
- GORM model auto-migration
- Transaction testing and rollback
- Multi-user scenarios
- UPSERT operation verification

### Error Handling Testing
- Invalid JSON message processing
- Missing required fields
- Database connection errors
- PubSub subscription errors
- Graceful error recovery

## Coverage Metrics

Expected coverage percentages:
- **Handler Layer**: >90% (comprehensive HTTP endpoint testing)
- **Service Layer**: >80% (core business logic coverage)
- **Models Layer**: >95% (simple data structure testing)
- **Overall**: >70% (meeting the coverage requirement)

## Running Specific Test Categories

```bash
# Test only authentication/authorization
go test -v -run "Test.*Auth" ./cmd

# Test only database operations  
go test -v -run "Test.*Notification.*" ./notification

# Test only JSON marshaling
go test -v -run "Test.*JSON" ./models

# Test error handling scenarios
go test -v -run "Test.*Error" ./...
```

## Continuous Integration

The test suite is designed to run in CI environments:

```yaml
# Example CI configuration
test:
  script:
    - go test -coverprofile=coverage.out ./cmd ./notification ./models
    - go tool cover -func=coverage.out
    - if [ "$CI_COMMIT_BRANCH" = "main" ]; then RUN_INTEGRATION_TESTS=true go test ./notification/integration_test.go; fi
```

## Test Data Management

### Test Database Cleanup
All tests use isolated test data and clean up after execution:
- Integration tests truncate tables before/after
- Unit tests use in-memory SQLite or isolated PostgreSQL transactions
- No test data pollution between test runs

### UUID Generation
Tests use `uuid.New()` for generating unique test data, ensuring no conflicts between parallel test runs.

This comprehensive testing approach ensures the notification service is thoroughly validated across all layers with high code coverage and reliable CI/CD integration.
