#!/bin/bash

# Load Testing Script for User Service
# Uses 'hey' tool: https://github.com/rakyll/hey

set -e

# Configuration
SERVICE_URL="${SERVICE_URL:-http://localhost:8081}"
OUTPUT_DIR="./load_test_results"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}User Service Load Testing${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Check if hey is installed
if ! command -v hey &> /dev/null; then
    echo -e "${RED}Error: 'hey' is not installed${NC}"
    echo "Install with: go install github.com/rakyll/hey@latest"
    exit 1
fi

# Check if service is running
if ! curl -s "${SERVICE_URL}/health" > /dev/null; then
    echo -e "${RED}Error: Service is not running at ${SERVICE_URL}${NC}"
    echo "Start the service first: go run cmd/server/main.go"
    exit 1
fi

echo -e "${GREEN}✓ Service is running at ${SERVICE_URL}${NC}"
echo ""

# Create output directory
mkdir -p "${OUTPUT_DIR}"

# Test 1: Health endpoint
echo -e "${YELLOW}Test 1: Health Endpoint${NC}"
hey -n 1000 -c 50 -m GET \
    "${SERVICE_URL}/health" \
    > "${OUTPUT_DIR}/${TIMESTAMP}_health.txt"
echo -e "${GREEN}✓ Results saved to ${OUTPUT_DIR}/${TIMESTAMP}_health.txt${NC}"
echo ""

# Test 2: Signup endpoint
echo -e "${YELLOW}Test 2: Signup Endpoint${NC}"
SIGNUP_PAYLOAD='{"email":"loadtest@example.com","password":"password123","name":"Load Test User"}'
hey -n 100 -c 10 -m POST \
    -H "Content-Type: application/json" \
    -d "${SIGNUP_PAYLOAD}" \
    "${SERVICE_URL}/users/signup" \
    > "${OUTPUT_DIR}/${TIMESTAMP}_signup.txt"
echo -e "${GREEN}✓ Results saved to ${OUTPUT_DIR}/${TIMESTAMP}_signup.txt${NC}"
echo ""

# Test 3: Login endpoint (requires user to exist)
echo -e "${YELLOW}Test 3: Login Endpoint${NC}"
LOGIN_PAYLOAD='{"email":"loadtest@example.com","password":"password123"}'
hey -n 500 -c 25 -m POST \
    -H "Content-Type: application/json" \
    -d "${LOGIN_PAYLOAD}" \
    "${SERVICE_URL}/users/login" \
    > "${OUTPUT_DIR}/${TIMESTAMP}_login.txt"
echo -e "${GREEN}✓ Results saved to ${OUTPUT_DIR}/${TIMESTAMP}_login.txt${NC}"
echo ""

# Test 4: Batch endpoint
echo -e "${YELLOW}Test 4: Batch Users Endpoint${NC}"
BATCH_PAYLOAD='{"user_ids":[]}'
hey -n 500 -c 25 -m POST \
    -H "Content-Type: application/json" \
    -d "${BATCH_PAYLOAD}" \
    "${SERVICE_URL}/users/batch" \
    > "${OUTPUT_DIR}/${TIMESTAMP}_batch.txt"
echo -e "${GREEN}✓ Results saved to ${OUTPUT_DIR}/${TIMESTAMP}_batch.txt${NC}"
echo ""

# Test 5: Stress test - High concurrency
echo -e "${YELLOW}Test 5: Stress Test (High Concurrency)${NC}"
hey -n 2000 -c 100 -m GET \
    "${SERVICE_URL}/health" \
    > "${OUTPUT_DIR}/${TIMESTAMP}_stress.txt"
echo -e "${GREEN}✓ Results saved to ${OUTPUT_DIR}/${TIMESTAMP}_stress.txt${NC}"
echo ""

# Test 6: Sustained load test
echo -e "${YELLOW}Test 6: Sustained Load Test (30 seconds)${NC}"
hey -z 30s -c 50 -m GET \
    "${SERVICE_URL}/health" \
    > "${OUTPUT_DIR}/${TIMESTAMP}_sustained.txt"
echo -e "${GREEN}✓ Results saved to ${OUTPUT_DIR}/${TIMESTAMP}_sustained.txt${NC}"
echo ""

# Generate summary report
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Generating Summary Report${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

SUMMARY_FILE="${OUTPUT_DIR}/${TIMESTAMP}_summary.txt"

{
    echo "Load Test Summary Report"
    echo "Generated: $(date)"
    echo "Service URL: ${SERVICE_URL}"
    echo ""
    echo "========================================="
    echo ""
    
    for file in "${OUTPUT_DIR}/${TIMESTAMP}"_*.txt; do
        if [[ "$file" != *"summary"* ]]; then
            filename=$(basename "$file")
            echo "--- ${filename} ---"
            # Extract key metrics
            grep -E "Requests/sec|Average|Fastest|Slowest|Status code" "$file" || true
            echo ""
        fi
    done
} > "${SUMMARY_FILE}"

echo -e "${GREEN}✓ Summary report saved to ${SUMMARY_FILE}${NC}"
echo ""

# Display summary
cat "${SUMMARY_FILE}"

echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}Load testing complete!${NC}"
echo -e "${BLUE}All results saved in: ${OUTPUT_DIR}/${NC}"
echo -e "${BLUE}========================================${NC}"