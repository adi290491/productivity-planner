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

hey -n 10000 -c 50 \
  -m POST \
  -H "X-USER-ID: 11111111-1111-1111-1111-111111111111" \
  -H "Content-Type: application/json" \
  -d '{"session_type":"focus"}' \
  http://localhost:8085/sessions/v1/start-session

# hey -n 10000 -c 50 \
#   -m PATCH \
#   -H "X-USER-ID: 11111111-1111-1111-1111-111111111111" \
#   -H "Content-Type: application/json" \
#   -d '{"session_type":"focus"}' \
#   http://localhost:8085/sessions/v1/stop-session

PGPASSWORD='S10dulkar' psql -p 5433 -h localhost -U adityasawant -d productivity_planner << 'EOF'
-- Create test user
DELETE FROM sessions WHERE user_id = '11111111-1111-1111-1111-111111111111';
-- Verify it exists
SELECT * from sessions where user_id = '11111111-1111-1111-1111-111111111111';
EOF


