#!/bin/bash

# Go Benchmark Script
set -e

OUTPUT_DIR="./benchmark_results"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Running Go Benchmarks${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Create output directory
mkdir -p "${OUTPUT_DIR}"

# Run all benchmarks
echo -e "${YELLOW}Running all benchmarks...${NC}"
go test -bench=. -benchmem -benchtime=3s ./... \
    > "${OUTPUT_DIR}/${TIMESTAMP}_benchmarks.txt" 2>&1

echo -e "${GREEN}✓ Benchmark results saved to ${OUTPUT_DIR}/${TIMESTAMP}_benchmarks.txt${NC}"
echo ""

# Run benchmarks with CPU profiling
echo -e "${YELLOW}Running benchmarks with CPU profiling...${NC}"
go test -bench=. -benchmem -cpuprofile="${OUTPUT_DIR}/${TIMESTAMP}_cpu.prof" \
    ./pkg/util/ > /dev/null 2>&1

echo -e "${GREEN}✓ CPU profile saved to ${OUTPUT_DIR}/${TIMESTAMP}_cpu.prof${NC}"
echo ""

# Run benchmarks with memory profiling
echo -e "${YELLOW}Running benchmarks with memory profiling...${NC}"
go test -bench=. -benchmem -memprofile="${OUTPUT_DIR}/${TIMESTAMP}_mem.prof" \
    ./pkg/util/ > /dev/null 2>&1

echo -e "${GREEN}✓ Memory profile saved to ${OUTPUT_DIR}/${TIMESTAMP}_mem.prof${NC}"
echo ""

# Display results
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Benchmark Results${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
cat "${OUTPUT_DIR}/${TIMESTAMP}_benchmarks.txt"

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}Benchmarking complete!${NC}"
echo -e "${BLUE}Results saved in: ${OUTPUT_DIR}/${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "To analyze CPU profile: go tool pprof ${OUTPUT_DIR}/${TIMESTAMP}_cpu.prof"
echo -e "To analyze memory profile: go tool pprof ${OUTPUT_DIR}/${TIMESTAMP}_mem.prof"