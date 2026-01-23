#!/bin/bash

# Single-user sequential benchmark
# Measures maximum throughput for one user doing start-stop cycles

BASE_URL="http://localhost:8085"
USER_ID="11111111-1111-1111-1111-111111111111"
NUM_CYCLES=1000
DB_HOST="localhost"
DB_PORT="5433"
DB_USER="adityasawant"
DB_NAME="productivity_planner"
export PGPASSWORD='S10dulkar'

echo "=========================================="
echo "Single-User Sequential Benchmark"
echo "=========================================="
echo "Configuration:"
echo "  - User ID: $USER_ID"
echo "  - Cycles: $NUM_CYCLES"
echo "  - Total requests: $((NUM_CYCLES * 2))"
echo "=========================================="

# Ensure user exists and clean up sessions
psql -p $DB_PORT -h $DB_HOST -U $DB_USER -d $DB_NAME << EOF
-- Ensure test user exists
INSERT INTO users (id, email, name, password_hash, created_at)
VALUES (
  '$USER_ID',
  'benchmark@example.com',
  'Benchmark User',
  'hashed_password',
  NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Clean up existing sessions
DELETE FROM sessions WHERE user_id = '$USER_ID';
SELECT 'Database prepared' as status;
EOF

echo ""
echo "Starting benchmark..."

start_time=$(date +%s.%N)
success_count=0
start_errors=0
stop_errors=0

for i in $(seq 1 $NUM_CYCLES); do
  # Start session
  start_code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/sessions/v1/start-session" \
    -H "X-USER-ID: $USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"session_type":"focus"}')
  
  if [ "$start_code" != "200" ]; then
    ((start_errors++))
    echo "Cycle $i: START failed with code $start_code"
    # Try to recover by stopping any active session
    curl -s -o /dev/null -X PATCH "$BASE_URL/sessions/v1/stop-session" \
      -H "X-USER-ID: $USER_ID" \
      -H "Content-Type: application/json" \
      -d '{"session_type":"focus"}'
    continue
  fi
  
  # Stop session
  stop_code=$(curl -s -o /dev/null -w "%{http_code}" -X PATCH "$BASE_URL/sessions/v1/stop-session" \
    -H "X-USER-ID: $USER_ID" \
    -H "Content-Type: application/json" \
    -d '{"session_type":"focus"}')
  
  if [ "$stop_code" != "200" ]; then
    ((stop_errors++))
    echo "Cycle $i: STOP failed with code $stop_code"
    continue
  fi
  
  ((success_count++))
  
  # Progress indicator every 100 cycles
  if [ $((i % 100)) -eq 0 ]; then
    echo "Completed $i/$NUM_CYCLES cycles..."
  fi
done

end_time=$(date +%s.%N)

# Calculate duration
duration=$(echo "$end_time - $start_time" | bc)
requests_per_sec=$(echo "scale=2; ($success_count * 2) / $duration" | bc)
avg_cycle_time=$(echo "scale=4; $duration / $success_count" | bc)

echo ""
echo "=========================================="
echo "Benchmark Results"
echo "=========================================="
echo "Duration: ${duration}s"
echo ""
echo "Cycle Summary:"
echo "  Attempted: $NUM_CYCLES"
echo "  Successful: $success_count ($(( success_count * 100 / NUM_CYCLES ))%)"
echo "  Failed: $((NUM_CYCLES - success_count))"
echo ""
echo "Error Breakdown:"
echo "  START errors: $start_errors"
echo "  STOP errors: $stop_errors"
echo ""
echo "Performance Metrics:"
echo "  Total requests: $((success_count * 2))"
echo "  Requests/sec: $requests_per_sec"
echo "  Avg cycle time: ${avg_cycle_time}s"
echo "  Avg request time: $(echo "scale=4; $avg_cycle_time / 2" | bc)s"

# Database verification
echo ""
echo "=========================================="
echo "Database Verification"
echo "=========================================="

psql -p $DB_PORT -h $DB_HOST -U $DB_USER -d $DB_NAME << EOF
SELECT 
  COUNT(*) as total_sessions,
  COUNT(CASE WHEN end_time IS NOT NULL THEN 1 END) as completed_sessions,
  COUNT(CASE WHEN end_time IS NULL THEN 1 END) as active_sessions,
  ROUND(AVG(EXTRACT(EPOCH FROM (end_time - start_time)))::numeric, 4) as avg_duration_seconds
FROM sessions 
WHERE user_id = '$USER_ID';
EOF

echo ""
echo "Benchmark complete!"
echo "=========================================="