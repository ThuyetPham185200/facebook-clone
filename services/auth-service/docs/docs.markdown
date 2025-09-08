# curl to test register/
curl -X POST http://localhost:9000/register   -H "Content-Type: application/json"   -d '{
    "username": "thuyetpq",
    "email": "thuyetpq@example.com",
    "password": "123456a@"
  }'

curl -X POST http://localhost:9000/login   -H "Content-Type: application/json"   -d '{
    "login": "thuyetpq",
    "password": "123456a@"
}'