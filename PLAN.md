# 📆 40-Day Project Plan: GoSocial App

A realistic and structured 40-day development plan to build a Facebook-like social app using **Golang**, PostgreSQL, and Redis.

---

## 🧠 Phase 1: System Design (Day 1–7)

| Day | Task                                                             |
| --- | ---------------------------------------------------------------- |
| 1   | Clarify **functional** and **non-functional** requirements       |
| 2   | Define **user stories** and workflows (e.g., follow, post, feed) |
| 3   | Design **ERD** (Entities: User, Post, Follow)                    |
| 4   | Define **API contract** (RESTful routes with inputs/outputs)     |
| 5   | Draft **high-level system architecture** (services, DB, cache)   |
| 6   | Design **DB schema** (PostgreSQL tables, indexes)                |
| 7   | Finalize and **review architecture** decisions                   |

---

## 🔧 Phase 2: Project Setup (Day 8–12)

| Day | Task                                                            |
| --- | --------------------------------------------------------------- |
| 8   | Initialize Go project: folder layout, Go Modules                |
| 9   | Setup REST router (Gin or Gorilla Mux)                          |
| 10  | Connect to PostgreSQL + Redis                                   |
| 11  | Add env config, error handling, and logging                     |
| 12  | Implement **mock login** & user registration (JWT/session stub) |

---

## 🧍️‍️ Phase 3: Core Features (Day 13–25)

| Day   | Task                                             |
| ----- | ------------------------------------------------ |
| 13–14 | `POST /posts` – Create a post                    |
| 15–16 | `GET /posts` – Get a user's own posts            |
| 17–18 | `POST /follow`, `POST /unfollow` – Follow system |
| 19–20 | Design feed logic (read-time/write-time fan-out) |
| 21–22 | Implement **Redis-based feed caching**           |
| 23–24 | `GET /feed?page=1` – Paginated feed              |
| 25    | Write **unit tests** for all core features       |

---

## ⚙️ Phase 4: Background & Optimization (Day 26–32)

| Day   | Task                                                    |
| ----- | ------------------------------------------------------- |
| 26    | Setup **background job queue** (e.g., `asynq`)          |
| 27–28 | Implement **fan-out on write** (push post to followers) |
| 29–30 | Handle **cache invalidation** & avoid race conditions   |
| 31    | DB & cache **performance tuning**                       |
| 32    | Add rate limiting, input validation                     |

---

## ✅ Phase 5: Testing & Deployment (Day 33–38)

| Day   | Task                                                   |
| ----- | ------------------------------------------------------ |
| 33    | Integration & end-to-end tests                         |
| 34    | Setup **CI pipeline** (GitHub Actions)                 |
| 35–36 | Dockerize: `Dockerfile`, `docker-compose.yml`          |
| 37    | Deploy to **VM**, **fly.io**, or **Render.com**        |
| 38    | Add **monitoring/logging** (Grafana, Prometheus, Logs) |

---

## 🧼 Phase 6: Polish & Wrap-up (Day 39–40)

| Day | Task                                  |
| --- | ------------------------------------- |
| 39  | UX feedback, bug fixing, polish       |
| 40  | Final **demo**, README, documentation |

---
