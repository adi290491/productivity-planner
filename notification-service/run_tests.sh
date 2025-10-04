#!/bin/bash

echo "Running notification service tests with coverage..."

# Run all tests with basic coverage display
echo -e "\n=== Running Tests with Basic Coverage ==="
go test ./... -cover

# Run all tests with detailed coverage profile
echo -e "\n=== Generating Detailed Coverage Profile ==="
go test ./... -coverprofile=cover.out

# Display detailed coverage by function and package
echo -e "\n=== Coverage by Function and Package ==="
go tool cover -func=cover.out

# Generate HTML coverage report
echo -e "\n=== Generating HTML Coverage Report ==="
go tool cover -html=cover.out -o coverage.html
echo "HTML coverage report generated: coverage.html"

echo -e "\n=== Running Integration Tests (if enabled) ==="
echo "Set RUN_INTEGRATION_TESTS=true to run integration tests"
echo "Set TEST_DATABASE_URL for custom test database connection"

# Run integration tests if enabled
if [ "$RUN_INTEGRATION_TESTS" = "true" ]; then
    echo "Running integration tests..."
    go test -v -tags=integration ./notification/integration_test.go
else
    echo "Integration tests skipped (set RUN_INTEGRATION_TESTS=true to enable)"
fi

echo -e "\n=== Test Summary ==="
echo "✅ Unit tests completed"
echo "✅ Models tests completed"
echo "✅ Handler tests completed"
echo "📊 Coverage report available in coverage.html"
