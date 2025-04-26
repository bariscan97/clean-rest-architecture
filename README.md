# Clean REST Architecture – Go

> **Modern, layered Golang REST API boilerplate** featuring Clean Architecture principles, Chi router, PostgreSQL (pgx) and Docker‑first workflow.

---

## ✨ Features

| Layer | Key Points |
|-------|------------|
| **Domain** (`internal/domains`) | Pure Go structs (entities) – no external imports |
| **Repository** (`internal/repository`) | Interface‑first, pgx implementation, easy to swap DB |
| **Handler / Delivery** (`internal/handler`) | Chi handlers, DTO<->Entity mapping, auth middleware |
| **App** (`app/routes`, `app/middleware`) | Central router & middlewares, graceful‑shutdown |
| **Pkg** (`pkg/*`) | Config (Viper), DB factory, JWT token maker (Zap logging) |

Additional goodies  
* Pagination & filtering helpers  
* JWT auth with 32‑byte secret key check  
* SQL migrations folder  
* Fully containerised (multi‑stage Dockerfile + Compose)  

---

## 🔧 Tech Stack

* **Go 1.22** (modules)  
* **Chi** HTTP router  
* **pgx/v5** PostgreSQL driver  
* **Viper** config + env overrides  
* **Zap** structured logger  
* **Docker + Compose** for local/dev  

---

## 📂 Project Layout

```text
clean-rest-architecture/
├─ cmd/               # entrypoint – main.go
├─ app/
│   └─ routes/        # http router definitions
├─ internal/
│   ├─ domains/       # business entities (no deps)
│   ├─ repository/    # repo interfaces + pg impl
│   └─ handler/       # http handlers (DTO layer)
├─ pkg/
│   ├─ config/        # Viper wrapper
│   ├─ database/      # pgx pool factory
│   └─ token/         # JWT maker/claims
├─ configs/           # config.yaml (default values)
├─ migrations/        # SQL schema & seed
├─ Dockerfile         # multi‑stage build
└─ docker-compose.yml # dev stack (api + postgres)
```

---

## 🚀 Quick Start (Docker)

```bash
# 1. Clone
$ git clone https://github.com/bariscan97/clean-rest-architecture.git && cd clean-rest-architecture

# 2. Build & run stack (API + Postgres)
$ docker compose up --build
```

*API lives on* **`http://localhost:3000`**.  
PostgreSQL exposed on `localhost:5432` (user/pass `postgres`).

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SECRET_KEY` | – | **Required** – ≥ 32‑byte key for JWT HMAC |
| `DATABASE_HOST` | `postgres` | Compose service name |
| `DATABASE_PORT` | `5432` | – |
| `DATABASE_USER` | `postgres` | – |
| `DATABASE_PASSWORD` | `postgres` | – |
| `DATABASE_DBNAME` | `my_db` | – |
| `DATABASE_SSLMODE` | `disable` | Set to `require` in prod |

> For local dev you can create a `.env` file; compose will pick it up automatically.

---

## 🖥️ Running Without Docker

```bash
# prerequisites: Go 1.22, PostgreSQL running locally
$ cp configs/config.yaml.example configs/config.yaml   # customise if needed
$ export SECRET_KEY="super‑secret‑32‑bytes‑key"
$ go run ./cmd            # runs on :3000 by default
```

---

## 🗄️ Database Migrations

```bash
# Example with golang‑migrate CLI
$ migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/my_db?sslmode=disable" up
```

> *Hint*: add a `migrate` service to `docker-compose.yml` in CI/CD.

---

## 📑 API Reference (v1)

```
GET    /health                 – liveness probe

# Posts
POST   /api/v1/posts           – create post        (auth required)
GET    /api/v1/posts           – list posts         (?user_id=&page=&limit=)
GET    /api/v1/posts/{id}/comments – nested comments
PATCH  /api/v1/posts/{id}      – update own post    (auth)
DELETE /api/v1/posts/{id}      – delete own post    (auth)
```

> Full swagger file coming soon – contributions welcome!

---

## 🧪 Testing

```bash
$ go test ./...       # unit + repository tests
```

Mocking is done via interfaces; repository layer can be swapped with a fake in tests.

---

## 🤝 Contributing

1. Fork & create feature branch (`git checkout -b feature/my-thing`)
2. Run `make vet fmt test` & ensure `go test ./...` passes
3. Submit PR 🚀

Please open an issue first for large changes.

---

## 📜 License

Distributed under the MIT License. See `LICENSE` for more information.

---

### Author

**Barış Karakuş** – [github.com/bariscan97](https://github.com/bariscan97)

