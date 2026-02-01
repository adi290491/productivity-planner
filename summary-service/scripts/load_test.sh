#!/bin/bash

set -e

SERVICE_URL="${SERVICE_URL:-http://localhost:8082}"
USER_ID="${USER_ID:-11111111-1111-1111-1111-111111111111}"
DATE=$(date +%Y-%m-%d)
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BENCHMARK_DIR="benchmark"
RESULTS_FILE="${BENCHMARK_DIR}/load-test_${TIMESTAMP}.txt"

REQUESTS="${REQUESTS:-10000}"
CONCURRENCY="${CONCURRENCY:-1000}"

mkdir -p "${BENCHMARK_DIR}"

# Check if hey is installed
if ! command -v hey &> /dev/null; then
    echo "Error: 'hey' is not installed"
    echo "Install with: go install github.com/rakyll/hey@latest"
    exit 1
fi

# Check service health
if ! curl -s -f "${SERVICE_URL}/health" > /dev/null; then
    echo "Error: Service is not responding at ${SERVICE_URL}"
    exit 1
fi

# Initialize results file
cat > "${RESULTS_FILE}" << EOF
=== Summary Service Load Test Results ===
Date: $(date)
Test: load_test.sh
Service: ${SERVICE_URL}
Requests: ${REQUESTS}
Concurrency: ${CONCURRENCY}

EOF

echo "Running load tests..."
echo "Service: ${SERVICE_URL}"
echo "Requests: ${REQUESTS}, Concurrency: ${CONCURRENCY}"
echo ""

# Test 1: Daily Summary
echo "Test 1: GET /summary/daily"
cat >> "${RESULTS_FILE}" << EOF
========================================
Test 1: GET /summary/daily
========================================

EOF

hey -n "${REQUESTS}" \
    -c "${CONCURRENCY}" \
    -H "X-USER-ID: ${USER_ID}" \
    "${SERVICE_URL}/summary/daily?date=${DATE}" \
    >> "${RESULTS_FILE}" 2>&1

echo "" >> "${RESULTS_FILE}"

# Test 2: Weekly Summary
echo "Test 2: GET /summary/weekly"
cat >> "${RESULTS_FILE}" << EOF
========================================
Test 2: GET /summary/weekly
========================================

EOF

hey -n "${REQUESTS}" \
    -c "${CONCURRENCY}" \
    -H "X-USER-ID: ${USER_ID}" \
    "${SERVICE_URL}/summary/weekly?start_date=${DATE}" \
    >> "${RESULTS_FILE}" 2>&1

echo "" >> "${RESULTS_FILE}"

# Test 3: Health endpoint (baseline)
echo "Test 3: GET /health (baseline)"
cat >> "${RESULTS_FILE}" << EOF
========================================
Test 3: GET /health (Baseline)
========================================

EOF

hey -n "${REQUESTS}" \
    -c "${CONCURRENCY}" \
    "${SERVICE_URL}/health" \
    >> "${RESULTS_FILE}" 2>&1

# Add completion timestamp
cat >> "${RESULTS_FILE}" << EOF

========================================
Test completed at: $(date)
Results saved to: ${RESULTS_FILE}
========================================
EOF

echo ""
echo "=== Test Complete ==="
echo "Results saved to: ${RESULTS_FILE}"
echo ""

# Show summary
echo "Summary:"
grep -A 3 "Summary:" "${RESULTS_FILE}" | head -n 15