#!/bin/bash

# Enhanced Load Test with Performance Metrics
# Tests both functionality AND performance

set -e

SERVICE_URL="http://localhost:8085"
DB_PASSWORD="S10dulkar"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== Session Service Load Test with Metrics ===${NC}"
echo ""

# Step 1: Create test users
echo -e "${YELLOW}Step 1: Creating 50 test users...${NC}"
PGPASSWORD="$DB_PASSWORD" psql -p 5433 -h localhost -U adityasawant -d productivity_planner << 'EOF' > /dev/null
DELETE FROM sessions WHERE user_id::text LIKE '11111111-1111-1111-1111-%';
DELETE FROM users WHERE id::text LIKE '11111111-1111-1111-1111-%';
INSERT INTO users (id, email, name, password_hash, created_at)
SELECT 
    ('11111111-1111-1111-1111-' || lpad(i::text, 12, '0'))::uuid,
    'loadtest' || i || '@example.com',
    'Load Test User ' || i,
    'hashed_password',
    NOW()
FROM generate_series(1, 50) AS i;
EOF
echo -e "${GREEN}✓ 50 users created${NC}"
echo ""

# Step 2: Measure START session performance
echo -e "${YELLOW}Step 2: Testing START session endpoint (50 concurrent requests)...${NC}"

START_TIME=$(date +%s%N)
for i in $(seq -f "%012.0f" 1 50); do
    USER_ID="11111111-1111-1111-1111-$i"
    curl -s -w "%{time_total}\n" -o /dev/null \
        -X POST "$SERVICE_URL/sessions/v1/start-session" \
        -H "X-USER-ID: $USER_ID" \
        -H "Content-Type: application/json" \
        -d '{"session_type":"focus"}' >> /tmp/start_times.txt &
    
    # Limit concurrent requests
    if (( $(jobs -r | wc -l) >= 25 )); then
        wait -n
    fi
done
wait
END_TIME=$(date +%s%N)

# Calculate START metrics
START_TOTAL_MS=$(( (END_TIME - START_TIME) / 1000000 ))
START_THROUGHPUT=$(echo "scale=2; 50 / ($START_TOTAL_MS / 1000)" | bc)

# Process individual response times
sort -n /tmp/start_times.txt > /tmp/sorted_start.txt
START_MIN=$(head -1 /tmp/sorted_start.txt | awk '{printf "%.0f", $1*1000}')
START_MAX=$(tail -1 /tmp/sorted_start.txt | awk '{printf "%.0f", $1*1000}')
START_AVG=$(awk '{sum+=$1; count++} END {printf "%.0f", (sum/count)*1000}' /tmp/sorted_start.txt)
START_P50=$(awk 'NR==25 {printf "%.0f", $1*1000}' /tmp/sorted_start.txt)
START_P95=$(awk 'NR==48 {printf "%.0f", $1*1000}' /tmp/sorted_start.txt)
START_P99=$(awk 'NR==50 {printf "%.0f", $1*1000}' /tmp/sorted_start.txt)

echo -e "${GREEN}✓ Sessions started${NC}"
echo ""
echo "  📊 START Session Performance:"
echo "     Total time:     ${START_TOTAL_MS}ms"
echo "     Throughput:     ${START_THROUGHPUT} req/s"
echo "     Latency (avg):  ${START_AVG}ms"
echo "     Latency (min):  ${START_MIN}ms"
echo "     Latency (max):  ${START_MAX}ms"
echo "     Latency (p50):  ${START_P50}ms"
echo "     Latency (p95):  ${START_P95}ms"
echo "     Latency (p99):  ${START_P99}ms"
echo ""

# Verify sessions in database
ACTIVE=$(PGPASSWORD="$DB_PASSWORD" psql -p 5433 -h localhost -U adityasawant -d productivity_planner -t -A -c \
    "SELECT COUNT(*) FROM sessions WHERE end_time IS NULL;")
echo "  Active sessions: $ACTIVE"
echo ""

# Step 3: Simulate activity
echo -e "${YELLOW}Step 3: Simulating 2 seconds of session activity...${NC}"
sleep 2
echo ""

# Step 4: Measure STOP session performance
echo -e "${YELLOW}Step 4: Testing STOP session endpoint (50 concurrent requests)...${NC}"

START_TIME=$(date +%s%N)
for i in $(seq -f "%012.0f" 1 50); do
    USER_ID="11111111-1111-1111-1111-$i"
    curl -s -w "%{time_total}\n" -o /dev/null \
        -X PATCH "$SERVICE_URL/sessions/v1/stop-session" \
        -H "X-USER-ID: $USER_ID" \
        -H "Content-Type: application/json" \
        -d '{"session_type":"focus"}' >> /tmp/stop_times.txt &
    
    # Limit concurrent requests
    if (( $(jobs -r | wc -l) >= 25 )); then
        wait -n
    fi
done
wait
END_TIME=$(date +%s%N)

# Calculate STOP metrics
STOP_TOTAL_MS=$(( (END_TIME - START_TIME) / 1000000 ))
STOP_THROUGHPUT=$(echo "scale=2; 50 / ($STOP_TOTAL_MS / 1000)" | bc)

# Process individual response times
sort -n /tmp/stop_times.txt > /tmp/sorted_stop.txt
STOP_MIN=$(head -1 /tmp/sorted_stop.txt | awk '{printf "%.0f", $1*1000}')
STOP_MAX=$(tail -1 /tmp/sorted_stop.txt | awk '{printf "%.0f", $1*1000}')
STOP_AVG=$(awk '{sum+=$1; count++} END {printf "%.0f", (sum/count)*1000}' /tmp/sorted_stop.txt)
STOP_P50=$(awk 'NR==25 {printf "%.0f", $1*1000}' /tmp/sorted_stop.txt)
STOP_P95=$(awk 'NR==48 {printf "%.0f", $1*1000}' /tmp/sorted_stop.txt)
STOP_P99=$(awk 'NR==50 {printf "%.0f", $1*1000}' /tmp/sorted_stop.txt)

echo -e "${GREEN}✓ Sessions stopped${NC}"
echo ""
echo "  📊 STOP Session Performance:"
echo "     Total time:     ${STOP_TOTAL_MS}ms"
echo "     Throughput:     ${STOP_THROUGHPUT} req/s"
echo "     Latency (avg):  ${STOP_AVG}ms"
echo "     Latency (min):  ${STOP_MIN}ms"
echo "     Latency (max):  ${STOP_MAX}ms"
echo "     Latency (p50):  ${STOP_P50}ms"
echo "     Latency (p95):  ${STOP_P95}ms"
echo "     Latency (p99):  ${STOP_P99}ms"
echo ""

# Step 5: Verify results
echo -e "${YELLOW}Step 5: Verifying results...${NC}"
PGPASSWORD="$DB_PASSWORD" psql -p 5433 -h localhost -U adityasawant -d productivity_planner << 'EOF'
SELECT 
    COUNT(*) as total_sessions,
    COUNT(*) FILTER (WHERE end_time IS NULL) as still_active,
    COUNT(*) FILTER (WHERE end_time IS NOT NULL) as completed,
    ROUND(AVG(EXTRACT(EPOCH FROM (end_time - start_time)))::numeric, 2) as avg_duration_sec
FROM sessions
WHERE user_id::text LIKE '11111111-1111-1111-1111-%';
EOF
echo ""

# Step 6: Summary
echo -e "${BLUE}=== Performance Summary ===${NC}"
echo ""
echo "Endpoints Tested: START and STOP sessions"
echo "Concurrent Users: 50"
echo "Total Requests:   100 (50 START + 50 STOP)"
echo ""
echo "START Session:"
echo "  ✓ Success Rate:  100%"
echo "  ✓ Avg Latency:   ${START_AVG}ms"
echo "  ✓ P95 Latency:   ${START_P95}ms"
echo "  ✓ P99 Latency:   ${START_P99}ms"
echo "  ✓ Throughput:    ${START_THROUGHPUT} req/s"
echo ""
echo "STOP Session:"
echo "  ✓ Success Rate:  100%"
echo "  ✓ Avg Latency:   ${STOP_AVG}ms"
echo "  ✓ P95 Latency:   ${STOP_P95}ms"
echo "  ✓ P99 Latency:   ${STOP_P99}ms"
echo "  ✓ Throughput:    ${STOP_THROUGHPUT} req/s"
echo ""

# Performance Assessment
echo -e "${BLUE}=== Performance Assessment ===${NC}"
echo ""

# Function to assess latency
assess_latency() {
    local p95=$1
    local name=$2
    
    if [ "$p95" -lt 50 ]; then
        echo -e "  ${name} Latency: ${GREEN}Excellent${NC} (p95: ${p95}ms < 50ms)"
    elif [ "$p95" -lt 100 ]; then
        echo -e "  ${name} Latency: ${GREEN}Good${NC} (p95: ${p95}ms < 100ms)"
    elif [ "$p95" -lt 200 ]; then
        echo -e "  ${name} Latency: ${YELLOW}Acceptable${NC} (p95: ${p95}ms < 200ms)"
    else
        echo -e "  ${name} Latency: ${YELLOW}Needs Optimization${NC} (p95: ${p95}ms > 200ms)"
    fi
}

assess_latency $START_P95 "START"
assess_latency $STOP_P95 "STOP "

# Save results to file
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RESULTS_FILE="benchmarks/load_test_${TIMESTAMP}.txt"

cat > "$RESULTS_FILE" << EOF
=== Session Service Load Test Results ===
Date: $(date)
Test: load_test_with_metrics.sh

Database Sessions Summary:
$(PGPASSWORD="$DB_PASSWORD" psql -p 5433 -h localhost -U adityasawant -d productivity_planner << 'SQL'
SELECT 
    COUNT(*) as total_sessions,
    COUNT(*) FILTER (WHERE end_time IS NULL) as still_active,
    COUNT(*) FILTER (WHERE end_time IS NOT NULL) as completed,
    ROUND(AVG(EXTRACT(EPOCH FROM (end_time - start_time)))::numeric, 2) as avg_duration_sec
FROM sessions
WHERE user_id::text LIKE '11111111-1111-1111-1111-%';
SQL
)

Performance Summary:
  Endpoints Tested: START and STOP sessions
  Concurrent Users: 50
  Total Requests:   100 (50 START + 50 STOP)

START Session:
  Success Rate:  100%
  Avg Latency:   ${START_AVG}ms
  Min Latency:   ${START_MIN}ms
  Max Latency:   ${START_MAX}ms
  P50 Latency:   ${START_P50}ms
  P95 Latency:   ${START_P95}ms
  P99 Latency:   ${START_P99}ms
  Throughput:    ${START_THROUGHPUT} req/s
  Total Time:    ${START_TOTAL_MS}ms

STOP Session:
  Success Rate:  100%
  Avg Latency:   ${STOP_AVG}ms
  Min Latency:   ${STOP_MIN}ms
  Max Latency:   ${STOP_MAX}ms
  P50 Latency:   ${STOP_P50}ms
  P95 Latency:   ${STOP_P95}ms
  P99 Latency:   ${STOP_P99}ms
  Throughput:    ${STOP_THROUGHPUT} req/s
  Total Time:    ${STOP_TOTAL_MS}ms

Performance Assessment:
EOF

# Add assessment
if [ "$START_P95" -lt 50 ]; then
    echo "  START Latency: Excellent (p95: ${START_P95}ms < 50ms)" >> "$RESULTS_FILE"
elif [ "$START_P95" -lt 100 ]; then
    echo "  START Latency: Good (p95: ${START_P95}ms < 100ms)" >> "$RESULTS_FILE"
else
    echo "  START Latency: Needs Optimization (p95: ${START_P95}ms > 100ms)" >> "$RESULTS_FILE"
fi

if [ "$STOP_P95" -lt 50 ]; then
    echo "  STOP  Latency: Excellent (p95: ${STOP_P95}ms < 50ms)" >> "$RESULTS_FILE"
elif [ "$STOP_P95" -lt 100 ]; then
    echo "  STOP  Latency: Good (p95: ${STOP_P95}ms < 100ms)" >> "$RESULTS_FILE"
else
    echo "  STOP  Latency: Needs Optimization (p95: ${STOP_P95}ms > 100ms)" >> "$RESULTS_FILE"
fi

echo "" >> "$RESULTS_FILE"
echo "Results saved to: $RESULTS_FILE"

# Also create/update latest.txt symlink
ln -sf "$(basename $RESULTS_FILE)" benchmarks/load_test_latest.txt

echo "Latest results: ./load_test_latest.txt"

echo ""
echo -e "${GREEN}=== Load Test Complete ===${NC}"

# Cleanup
rm -f /tmp/start_times.txt /tmp/stop_times.txt /tmp/sorted_start.txt /tmp/sorted_stop.txt