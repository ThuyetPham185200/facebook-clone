gateway-api/
 ├── cmd/
 │   └── gateway-api/
 │       └── main.go      # entrypoint
 └── internal/
     ├── app/             # orchestration (App struct quản lý lifecycle)
     ├── config/          # Quản lý toàn bộ cấu hình của service (port, DB URL, JWT secret, Redis host, rate-limit…).
     ├── middleware/
     ├── server/
     └── api/
