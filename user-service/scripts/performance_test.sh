#!/bin/bash

# Complete Performance Testing Suite
set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}User Service Performance Test Suite${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Check prerequisites
echo -e "${YELLOW}Checking prerequisites...${NC}"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    exit 1
fi

# Check if hey is installed
if ! command -v hey &> /dev/null; then
    echo -e "${YELLOW}Installing 'hey' load testing tool...${NC}"
    go install github.com/rakyll/hey@latest
fi

echo -e "${GREEN}✓ All prerequisites met${NC}"
echo ""

# Step 1: Run unit tests
echo -e "${BLUE}Step 1: Running unit tests${NC}"
go test -v ./... 2>&1 | tee test_results.txt
echo -e "${GREEN}✓ Unit tests complete${NC}"
echo ""

# Step 2: Run benchmarks
echo -e "${BLUE}Step 2: Running benchmarks${NC}"
./scripts/benchmark.sh
echo -e "${GREEN}✓ Benchmarks complete${NC}"
echo ""

# Step 3: Start service for load testing
echo -e "${BLUE}Step 3: Starting service for load testing${NC}"
go build -o bin/server cmd/server/main.go
./bin/server &
SERVER_PID=$!

# Wait for service to start
echo -e "${YELLOW}Waiting for service to start...${NC}"
sleep 3

# Check if service is running
if ! curl -s http://localhost:8081/health > /dev/null; then
    echo -e "${RED}Error: Service failed to start${NC}"
    kill $SERVER_PID 2>/dev/null || true
    exit 1
fi

echo -e "${GREEN}✓ Service started (PID: $SERVER_PID)${NC}"
echo ""

# Step 4: Run load tests
echo -e "${BLUE}Step 4: Running load tests${NC}"
./scripts/load_test.sh
echo -e "${GREEN}✓ Load tests complete${NC}"
echo ""

# Step 5: Stop service
echo -e "${BLUE}Step 5: Stopping service${NC}"
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true
echo -e "${GREEN}✓ Service stopped${NC}"
echo ""

# Final summary
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}Performance Testing Complete!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "Results locations:"
echo "  - Unit tests: test_results.txt"
echo "  - Benchmarks: ./benchmark_results/"
echo "  - Load tests: ./load_test_results/"
echo ""