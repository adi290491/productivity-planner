# psql -p 5433 -h localhost -U adityasawant -d productivity_planner << 'EOF'
# -- Create test user
# INSERT INTO users (id, email, name, password_hash, created_at)
# VALUES (
#   '11111111-1111-1111-1111-111111111111',
#   'focususer@example.com',
#   'Focus User',
#   'hashed_password_1',
#   NOW()
# );
# -- Verify it exists
# SELECT id, email, name FROM users;
# EOF

# hey -n 10000 -c 50 \
#   -m POST \
#   -H "X-USER-ID: 11111111-1111-1111-1111-111111111111" \
#   -H "Content-Type: application/json" \
#   -d '{"session_type":"focus"}' \
#   http://localhost:8085/sessions/v1/start-session

# # hey -n 10000 -c 50 \
# #   -m PATCH \
# #   -H "X-USER-ID: 11111111-1111-1111-1111-111111111111" \
# #   -H "Content-Type: application/json" \
# #   -d '{"session_type":"focus"}' \
# #   http://localhost:8085/sessions/v1/stop-session

# PGPASSWORD='S10dulkar' psql -p 5433 -h localhost -U adityasawant -d productivity_planner << 'EOF'
# -- Create test user
# DELETE FROM sessions WHERE user_id = '11111111-1111-1111-1111-111111111111';
# -- Verify it exists
# SELECT * from sessions where user_id = '11111111-1111-1111-1111-111111111111';
# EOF


#!/bin/bash

BASE_URL="http://localhost:8085"
USER_ID="11111111-1111-1111-1111-111111111111"

# Clean up first
PGPASSWORD='S10dulkar' psql -p 5433 -h localhost -U adityasawant -d productivity_planner << EOF
DELETE FROM sessions WHERE user_id = '$USER_ID';
EOF

echo "Testing START endpoint..."
hey -n 5000 -c 1 \
  -m POST \
  -H "X-USER-ID: $USER_ID" \
  -H "Content-Type: application/json" \
  -d '{"session_type":"focus"}' \
  "$BASE_URL/sessions/v1/start-session" \
  > start_results.txt &

START_PID=$!

# Small delay to ensure first start happens
sleep 0.5

echo "Testing STOP endpoint..."
hey -n 5000 -c 1 \
  -m PATCH \
  -H "X-USER-ID: $USER_ID" \
  -H "Content-Type: application/json" \
  -d '{"session_type":"focus"}' \
  "$BASE_URL/sessions/v1/stop-session" \
  > stop_results.txt &

STOP_PID=$!

# Wait for both to complete
wait $START_PID
wait $STOP_PID

echo "Results:"
echo "=========="
echo "START requests:"
grep "Status code distribution:" -A 5 start_results.txt
echo ""
echo "STOP requests:"
grep "Status code distribution:" -A 5 stop_results.txt