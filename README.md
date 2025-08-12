# 🔵 GoSocial - A Facebook-like Social Feed System

A scalable, high-performance backend service written in **Golang** that supports:

* Creating posts
* Following users
* Viewing a feed in chronological order

> ⚙️ Technologies: Go, PostgreSQL, Redis, Docker

---

## 📦 Features

* 📝 Create & publish posts
* ➕ Follow/Unfollow other users
* 📜 Paginated newsfeed (sorted by time)
* ⚡ Fast feed read with Redis caching
* 🧱 Clean architecture (Handlers, Services, Repositories)
* 🐳 Containerized with Docker & Docker Compose

---

## 🚀 Quick Start

### 1. Clone the repo

```bash
git clone https://github.com/your-username/gosocial.git
cd gosocial
```

### 2. Start using Docker

```bash
docker-compose up --build
```

### 3. Access

* API Server[: ](http://localhost:8080)[http://localhost:8080](http://localhost:8080)
* PostgreSQL: `localhost:5432`
* Redis: `localhost:6379`

---

## 🗂️ Project Structure

```
facebook-clone/
│
├── services/
│   ├── gateway-api/                # API Gateway (entry point)
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   ├── handler/
│   │   │   ├── service/
│   │   │   ├── repository/
│   │   │   └── ...
│   │   ├── pkg/
│   │   └── go.mod
│   │
│   ├── user-service/               # User Service
│   ├── follow-service/             # Follow Service
│   ├── post-service/               # Post Service
│   ├── feed-service/               # Feed Service
│   ├── like-service/               # Like Service
│   ├── comment-service/            # Comment Service
│   ├── notification-service/       # Notification Service
│   ├── content-delivery-service/   # Media/CDN handling
│   └── auth-service/               # Authentication & token
│
├── shared/                         # Share code dùng chung cho tất cả service
│   ├── configs/
│   ├── logger/
│   ├── middleware/
│   ├── models/
│   ├── utils/
│   ├── proto/                      # Nếu dùng gRPC
│   └── events/                     # Event schemas cho Kafka/NATS
│
├── deployments/                    # Docker-compose, k8s manifests
│   ├── docker-compose.yml
│   ├── k8s/
│   └── ...
│
├── docs/
│   ├── architecture-diagram.png
│   ├── api-design.md
│   └── ...
│
├── go.work                         # Nếu dùng Go workspace
└── README.md
```

---

## ⚙️ API Endpoints

### User

* `POST /follow` – Follow a user
* `POST /unfollow` – Unfollow a user

### Post

* `POST /posts` – Create a post

### Feed

* `GET /feed?page=1&size=10` – Get paginated feed

---

## 🧪 Testing

```bash
go test ./...
```

---

## 📌 Environment Variables

Create a `.env` file:

```env
PORT=8080
DB_URL=postgres://user:password@db:5432/social_db?sslmode=disable
REDIS_ADDR=redis:6379
```

---

## 🏗️ System Design

### 🧱 High-Level Architecture

```
Users
  ↓
Web Server (API Gateway)
  ↓
───────────────────────────────────────────────
| Microservices:                             |
|                                            |
|  • User Service (CRUD, follow logic)       |
|  • Post Service (create & store posts)     |
|  • Feed Service (generate & fetch feeds)   |
|  • Auth Service (JWT/session validation)   |
───────────────────────────────────────────────
  ↓
Cache (Redis) ←→ DB (PostgreSQL)
```

### ⚙️ Components

| Component    | Description                              |
| ------------ | ---------------------------------------- |
| Web Server   | API layer using Go (e.g., Gorilla Mux)   |
| User Service | Handles following logic & user info      |
| Post Service | Stores posts in DB                       |
| Feed Service | Fetches posts from followed users        |
| Redis Cache  | Caches user feeds for fast reads         |
| PostgreSQL   | Main storage for users, posts, relations |

### 🔁 Feed Generation Strategy

* **Fan-out on write** (for active users): When a user posts, push it to followers’ Redis feeds
* **Fan-out on read** (for inactive users): Read from DB on-demand
* Use **cache invalidation** to ensure eventual consistency

---

## 🐳 Docker Setup

The project uses Docker and Docker Compose for ease of setup and deployment.

**Dockerfile:** Builds the Go application binary.

**docker-compose.yml:**

* `app`: Golang application container
* `db`: PostgreSQL instance
* `redis`: Redis instance

---

## 📄 License

MIT © 2025 Your Name
