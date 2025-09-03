# curl to test register/
curl -X POST http://localhost:9000/register   -H "Content-Type: application/json"   -d '{
    "username": "thuyetpq",
    "email": "thuyetpq@example.com",
    "password": "123456a@"
  }'
