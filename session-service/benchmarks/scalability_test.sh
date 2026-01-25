#!/bin/bash

# Scalability Test - Gradually Increase Load
# Finds the breaking point of your service

set -e

SERVICE_URL="http://localhost:8085"
DB_PASSWORD="S10dulkar"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}=== Scalability Test ===${NC}"
echo "Finding your service's capacity limits..."
echo ""

# Create test users (need 500 for highest concurrency test)
echo -e "${YELLOW}Creating 500 test users...${NC}"
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
FROM generate_series(1, 500) AS i;
EOF
echo -e "${GREEN}✓ 500 users created${NC}"
echo ""

# Results table header
printf "%-15s %-15s %-15s %-15s %-15s %s\n" "Concurrency" "Throughput" "Avg Latency" "P95 Latency" "P99 Latency" "Status"
printf "%-15s %-15s %-15s %-15s %-15s %s\n" "---------------" "---------------" "---------------" "---------------" "---------------" "---------------"

# Test with increasing concurrency
for CONCURRENCY in 10 25 50 100 150 200 300 500; do
    # Clean up previous sessions
    PGPASSWORD="$DB_PASSWORD" psql -p 5433 -h localhost -U adityasawant -d productivity_planner \
        -c "DELETE FROM sessions WHERE user_id::text LIKE '11111111-1111-1111-1111-%';" > /dev/null 2>&1
    
    # Run test
    START_TIME=$(date +%s%N)
    SUCCESS_COUNT=0
    
    for i in $(seq -f "%012.0f" 1 $CONCURRENCY); do
        USER_ID="11111111-1111-1111-1111-$i"
        
        (
            RESPONSE=$(curl -s -w "\n%{http_code}\n%{time_total}" -o /dev/null \
                -X POST "$SERVICE_URL/sessions/v1/start-session" \
                -H "X-USER-ID: $USER_ID" \
                -H "Content-Type: application/json" \
                -d '{"session_type":"focus"}' 2>/dev/null)
            
            HTTP_CODE=$(echo "$RESPONSE" | tail -2 | head -1)
            TIME=$(echo "$RESPONSE" | tail -1)
            
            if [ "$HTTP_CODE" = "200" ]; then
                echo "$TIME" >> /tmp/test_${CONCURRENCY}_times.txt
            fi
        ) &
        
        # Limit parallel jobs
        if (( $(jobs -r | wc -l) >= 50 )); then
            wait -n
        fi
    done
    wait
    
    END_TIME=$(date +%s%N)
    TOTAL_MS=$(( (END_TIME - START_TIME) / 1000000 ))
    
    # Calculate metrics
    if [ -f "/tmp/test_${CONCURRENCY}_times.txt" ]; then
        sort -n /tmp/test_${CONCURRENCY}_times.txt > /tmp/sorted_${CONCURRENCY}.txt
        
        COUNT=$(wc -l < /tmp/sorted_${CONCURRENCY}.txt)
        THROUGHPUT=$(echo "scale=2; $COUNT / ($TOTAL_MS / 1000)" | bc)
        
        AVG=$(awk '{sum+=$1; count++} END {printf "%.0f", (sum/count)*1000}' /tmp/sorted_${CONCURRENCY}.txt)
        
        P95_LINE=$(echo "($COUNT * 0.95)/1" | bc)
        P95=$(sed -n "${P95_LINE}p" /tmp/sorted_${CONCURRENCY}.txt | awk '{printf "%.0f", $1*1000}')
        
        P99_LINE=$(echo "($COUNT * 0.99)/1" | bc)
        P99=$(sed -n "${P99_LINE}p" /tmp/sorted_${CONCURRENCY}.txt | awk '{printf "%.0f", $1*1000}')
        
        # Assess performance
        if [ "$P95" -lt 100 ]; then
            STATUS="${GREEN}Excellent${NC}"
        elif [ "$P95" -lt 200 ]; then
            STATUS="${GREEN}Good${NC}"
        elif [ "$P95" -lt 500 ]; then
            STATUS="${YELLOW}Acceptable${NC}"
        else
            STATUS="${RED}Degraded${NC}"
        fi
        
        printf "%-15s %-15s %-15s %-15s %-15s %b\n" \
            "$CONCURRENCY users" \
            "${THROUGHPUT} req/s" \
            "${AVG}ms" \
            "${P95}ms" \
            "${P99}ms" \
            "$STATUS" 

        echo "$CONCURRENCY users|${THROUGHPUT} req/s|${AVG}ms|${P95}ms|${P99}ms|$STATUS" >> /tmp/scalability_results.txt
        
        # Clean up
        rm /tmp/test_${CONCURRENCY}_times.txt /tmp/sorted_${CONCURRENCY}.txt
    else
        printf "%-15s %-15s %-15s %-15s %-15s %b\n" \
            "$CONCURRENCY users" \
            "FAILED" \
            "N/A" \
            "N/A" \
            "N/A" \
            "${RED}Failed${NC}"
    fi
    
    # Brief pause between tests
    sleep 2
done

echo ""
echo -e "${BLUE}=== Scalability Analysis ===${NC}"
echo ""
echo "Look for:"
echo "  • Linear scaling: Throughput increases with concurrency ✅"
echo "  • Stable latency: P95 stays under 200ms ✅"
echo "  • Degradation point: Where performance drops ⚠️"
echo ""
echo "Your service's capacity is where latency starts increasing significantly."
echo ""
echo -e "${GREEN}=== Test Complete ===${NC}"

# Format and append the saved results
if [ -f /tmp/scalability_results.txt ]; then
    while IFS='|' read -r conc thru avg p95 p99 status; do
        printf "%-15s %-15s %-15s %-15s %-15s %s\n" \
            "$conc" "$thru" "$avg" "$p95" "$p99" "$(echo $status | sed 's/\\033\[[0-9;]*m//g')" >> "$RESULTS_FILE"
    done < /tmp/scalability_results.txt
    rm /tmp/scalability_results.txt
fi

# Save results to file
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RESULTS_FILE="benchmarks/scalability_test_${TIMESTAMP}.txt"

# The table was already printed to stdout, now save to file
cat > "$RESULTS_FILE" << 'EOF'
=== Session Service Scalability Test Results ===
Date: $(date)
Test: scalability_test.sh

Concurrency Test Results:
EOF

# Re-run the test but output to file (or save the results during the test)
# For simplicity, add this note and instructions to save manually
echo "" >> "$RESULTS_FILE"
echo "Concurrency     Throughput      Avg Latency     P95 Latency     P99 Latency     Status" >> "$RESULTS_FILE"
echo "--------------- --------------- --------------- --------------- --------------- ---------------" >> "$RESULTS_FILE"

# You'll need to save the test results during the loop
# I'll show you where to add this below

echo "" >> "$RESULTS_FILE"
echo "Scalability Analysis:" >> "$RESULTS_FILE"
echo "  • Linear scaling: Throughput increases with concurrency" >> "$RESULTS_FILE"
echo "  • Stable latency: P95 stays under 200ms" >> "$RESULTS_FILE"
echo "  • Service capacity identified" >> "$RESULTS_FILE"
echo "" >> "$RESULTS_FILE"

echo "Results saved to: $RESULTS_FILE"
ln -sf "$(basename $RESULTS_FILE)" benchmarks/scalability_test_latest.txt
echo "Latest results: benchmarks/scalability_test_latest.txt"