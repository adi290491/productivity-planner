#!/bin/bash

# Configuration
BASE_URL="http://localhost:8085"
NUM_USERS=100
CYCLES_PER_USER=100
DB_HOST="localhost"
DB_PORT="5433"
DB_USER="adityasawant"
DB_NAME="productivity_planner"
export PGPASSWORD='S10dulkar'

echo "=========================================="
echo "Session Service Load Test"
echo "=========================================="
echo "Configuration:"
echo "  - Base URL: $BASE_URL"
echo "  - Number of users: $NUM_USERS"
echo "  - Cycles per user: $CYCLES_PER_USER"
echo "  - Total requests: $((NUM_USERS * CYCLES_PER_USER * 2))"
echo "=========================================="

# Clean up and create test users
echo ""
echo "Setting up test users..."
psql -p $DB_PORT -h $DB_HOST -U $DB_USER -d $DB_NAME << 'EOF'
-- Clean up existing test data
DELETE FROM sessions WHERE user_id LIKE '00000000-0000-0000-0000-%';
DELETE FROM users WHERE id LIKE '00000000-0000-0000-0000-%';

-- Create test users
DO $$
BEGIN
  FOR i IN 1..100 LOOP
    INSERT INTO users (id, email, name, password_hash, created_at)
    VALUES (
      '00000000-0000-0000-0000-' || lpad(i::text, 12, '0'),
      'loadtest' || i || '@example.com',
      'Load Test User ' || i,
      'hashed_password',
      NOW()
    );
  END LOOP;
END $$;

SELECT COUNT(*) as "Test users created" FROM users WHERE id LIKE '00000000-0000-0000-0000-%';
EOF

echo "Test users created successfully!"

# Function to run start-stop cycles for one user
run_user_load_test() {
  local user_num=$1
  local user_id=$(printf "00000000-0000-0000-0000-%012d" $user_num)
  local success_count=0
  local start_errors=0
  local stop_errors=0
  
  for cycle in $(seq 1 $CYCLES_PER_USER); do
    # Start session
    start_response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/sessions/v1/start-session" \
      -H "X-USER-ID: $user_id" \
      -H "Content-Type: application/json" \
      -d '{"session_type":"focus"}' 2>/dev/null)
    
    start_code=$(echo "$start_response" | tail -n1)
    
    if [ "$start_code" != "200" ]; then
      ((start_errors++))
      continue
    fi
    
    # Simulate session duration (10-100ms)
    sleep 0.0$(( RANDOM % 90 + 10 ))
    
    # Stop session
    stop_response=$(curl -s -w "\n%{http_code}" -X PATCH "$BASE_URL/sessions/v1/stop-session" \
      -H "X-USER-ID: $user_id" \
      -H "Content-Type: application/json" \
      -d '{"session_type":"focus"}' 2>/dev/null)
    
    stop_code=$(echo "$stop_response" | tail -n1)
    
    if [ "$stop_code" != "200" ]; then
      ((stop_errors++))
      continue
    fi
    
    ((success_count++))
  done
  
  echo "$user_num,$success_count,$start_errors,$stop_errors"
}

# Export function and variables for parallel execution
export -f run_user_load_test
export BASE_URL CYCLES_PER_USER

# Record start time
start_time=$(date +%s)

echo ""
echo "Starting load test..."
echo "Running $NUM_USERS users in parallel..."

# Create results directory
mkdir -p /tmp/loadtest_results
rm -f /tmp/loadtest_results/*

# Run all users in parallel and collect results
seq 1 $NUM_USERS | xargs -P $NUM_USERS -I {} bash -c 'run_user_load_test {} > /tmp/loadtest_results/user_{}.txt'

# Record end time
end_time=$(date +%s)
duration=$((end_time - start_time))

echo "Load test completed in ${duration}s!"

# Aggregate results
echo ""
echo "=========================================="
echo "Load Test Results"
echo "=========================================="

total_success=0
total_start_errors=0
total_stop_errors=0

for file in /tmp/loadtest_results/user_*.txt; do
  if [ -f "$file" ]; then
    IFS=',' read -r user_num success start_err stop_err < "$file"
    ((total_success += success))
    ((total_start_errors += start_err))
    ((total_stop_errors += stop_err))
  fi
done

total_cycles=$((NUM_USERS * CYCLES_PER_USER))
total_requests=$((total_cycles * 2))
successful_requests=$((total_success * 2))
failed_requests=$((total_start_errors + total_stop_errors))

echo "Duration: ${duration}s"
echo "Throughput: $((successful_requests / duration)) successful requests/sec"
echo ""
echo "Request Summary:"
echo "  Total cycles attempted: $total_cycles"
echo "  Successful cycles: $total_success ($(( total_success * 100 / total_cycles ))%)"
echo "  Total requests: $total_requests"
echo "  Successful requests: $successful_requests"
echo "  Failed requests: $failed_requests"
echo ""
echo "Error Breakdown:"
echo "  START errors: $total_start_errors"
echo "  STOP errors: $total_stop_errors"

# Database verification
echo ""
echo "=========================================="
echo "Database Verification"
echo "=========================================="

psql -p $DB_PORT -h $DB_HOST -U $DB_USER -d $DB_NAME << 'EOF'
SELECT 
  COUNT(*) as total_sessions,
  COUNT(CASE WHEN end_time IS NOT NULL THEN 1 END) as completed_sessions,
  COUNT(CASE WHEN end_time IS NULL THEN 1 END) as active_sessions,
  ROUND(AVG(EXTRACT(EPOCH FROM (end_time - start_time)))::numeric, 2) as avg_duration_seconds,
  ROUND(MIN(EXTRACT(EPOCH FROM (end_time - start_time)))::numeric, 2) as min_duration_seconds,
  ROUND(MAX(EXTRACT(EPOCH FROM (end_time - start_time)))::numeric, 2) as max_duration_seconds
FROM sessions 
WHERE user_id LIKE '00000000-0000-0000-0000-%';

-- Show distribution by user (sample 10 users)
SELECT 
  user_id,
  COUNT(*) as session_count,
  COUNT(CASE WHEN end_time IS NOT NULL THEN 1 END) as completed
FROM sessions 
WHERE user_id LIKE '00000000-0000-0000-0000-%'
GROUP BY user_id
ORDER BY user_id
LIMIT 10;
EOF

# Cleanup
echo ""
echo "=========================================="
echo "Cleaning up test data..."
psql -p $DB_PORT -h $DB_HOST -U $DB_USER -d $DB_NAME << 'EOF'
DELETE FROM sessions WHERE user_id LIKE '00000000-0000-0000-0000-%';
DELETE FROM users WHERE id LIKE '00000000-0000-0000-0000-%';
SELECT 'Test data cleaned up' as status;
EOF

rm -rf /tmp/loadtest_results

echo "Load test complete!"
echo "=========================================="