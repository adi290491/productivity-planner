docker-compose restart gateway 

# for i in {1..21}; do
#   curl -s -X POST http://localhost:8000/users/login \
#     -H "Origin: http://localhost:5173" \
#     -H "Content-Type: application/json" \
#     -H "X-Forwarded-For: 192.168.1.100" \
#     -d '{"email":"test@example.com","password":"test123"}' \
#     -o /dev/null &
# done
# wait# Simulate attacker sending rapid requests
# for i in {1..5}; do
#   echo "=== Request $i ==="
#   curl -i -X POST http://localhost:8000/users/login \
#     -H "Origin: http://localhost:5173" \
#     -H "Content-Type: application/json" \
#     -H "X-Forwarded-For: 192.168.1.100" \
#     -d '{"email":"test@example.com","password":"test123"}' 2>/dev/null \
#     | grep -E "HTTP|X-RateLimit"
#   echo ""
# done
seq 21 | parallel -j 21 "curl -s -w '%{http_code}\n' \
  -X POST http://localhost:8000/users/login \
  -H 'Origin: http://localhost:5173' \
  -H 'Content-Type: application/json' \
  -H 'X-Forwarded-For: 192.168.1.100' \
  -d '{\"email\":\"test@example.com\",\"password\":\"test123\"}' \
  -o /dev/null" | sort | uniq -c