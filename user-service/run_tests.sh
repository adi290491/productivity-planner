#!/bin/bash

# Navigate to the user-service directory (where this script is located)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "Starting test execution..."
echo "Current directory: $(pwd)"

# First, let's check if go.mod exists
if [ ! -f go.mod ]; then
    echo "ERROR: go.mod not found in current directory"
    exit 1
fi

echo "Go module file found. Checking dependencies..."

# Download dependencies
echo "Running go mod tidy..."
go mod tidy

# Try to compile first
echo "Checking if code compiles..."
go build ./...
if [ $? -ne 0 ]; then
    echo "ERROR: Code compilation failed"
    exit 1
fi

echo "Code compiles successfully. Running tests..."

# Run tests with coverage, excluding integration tests by default
go test -v -coverprofile=coverage.out ./...
test_status=$?
if [ $test_status -ne 0 ]; then
    echo "ERROR: Unit tests failed"
    exit $test_status
fi

# Check if coverage.out was created
if [ -f coverage.out ]; then
    echo "Coverage report generated successfully."
    echo "Displaying coverage summary:"
    go tool cover -func=coverage.out | tail -1
    
    echo "Detailed coverage by package:"
    go tool cover -func=coverage.out
    
    # Generate HTML coverage report
    echo "Generating HTML coverage report..."
    go tool cover -html=coverage.out -o coverage.html
    go tool cover -html=coverage.out -o coverage.html
    echo "HTML coverage report generated: coverage.html"
else
    echo "Coverage report not generated. Tests may have failed."
fi

echo "Test execution completed."
