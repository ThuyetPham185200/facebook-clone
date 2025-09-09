# curl to test register/
# Register
curl -X POST http://localhost:9000/register   -H "Content-Type: application/json"   -d '{
    "username": "thuyetpq",
    "email": "thuyetpq@example.com",
    "password": "123456a@"
  }'

# Login
curl -X POST http://localhost:9000/login   -H "Content-Type: application/json"   -d '{
    "login": "thuyetpq",
    "password": "123456a@"
}'

# Change password
curl -X PUT http://localhost:9000/me/password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTc0MTA5MDMsInVzZXJfaWQiOiJjZmU1YmI2Ni03OTUyLTRhYWYtYThjOS1hNzJlN2VhZDIzZTIifQ.Sfo-kRMBdBtYkco7D0jUAoymP8judwNjuzn5WL5xhHw" \
  -d '{
    "old_password": "123456a@",
    "new_password": "1234567a@"
  }'
