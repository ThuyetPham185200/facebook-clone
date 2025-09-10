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

# Refresh token
curl -X POST http://localhost:9000/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTgwMjk5NjksInVzZXJfaWQiOiI4OGFjZGU0ZS1hYzFlLTRlMGEtYTY1OS0xYTk4NjE2ZDZmMmEifQ.hkYHVymiqXR-4w0VuD46j4tPzvHcfVeM963KMMefNCQ"
  }'

# Delete account
curl -X DELETE http://localhost:9000/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTc0OTA4MzgsInVzZXJfaWQiOiI3NDM2NDNjYi1lNDEyLTRlMzMtOWI4Mi00NzlmYmYwNTA4NDQifQ.PbhxEY5uXtxhGWfdqnjXZ0PiD8G9KcDoFcIY-JT_O8g"
