.PHONY: test test-unit test-integration test-coverage bench load-test perf-test

# Run all unit tests
test-unit:
	go test -v -short ./...

# Run integration tests only
test-integration:
	go test -v -tags=integration ./internal/repository/

# Run all tests
tests:
	go test -v ./...
	go test -v -tags=integration ./internal/repository

# Run tests with coverage
test-coverage:
	go test -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
bench:
	@chmod +x scripts/benchmark.sh
	./scripts/benchmark.sh

# Run load tests
load-test:
	@chmod +x scripts/load-test.sh
	./scripts/load-test.sh

# Run complete performance test suite
perf-test:
	@chmod +x scripts/performance_test.sh
	./scripts/performance_test.sh

# Build the service
build:
	go build -o bin/server cmd/server/main.go

# Run the service
run:
	go run cmd/server/main.go

# Clean build artifacts and test results
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -rf benchmark_results/
	rm -rf load_test_results/
	rm -f test_results.txt

# Install dependencies
deps:
	go mod download
	go mod tidy
	go install github.com/rakyll/hey@latest