docker-compose restart gateway 

for i in {1..12}; do
  response=$(curl -s -w "%{http_code}" \
    -X POST http://localhost:8000/users/login \
    -H "Origin: http://localhost:5173" \
    -H "Content-Type: application/json" \
    -H "X-Forwarded-For: 192.168.1.100" \
    -d '{"email":"test@example.com","password":"test123"}' \
    -o /dev/null)
  echo "Request $i: $response"
done