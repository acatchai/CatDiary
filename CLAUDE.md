# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Development
- `make dev` — Start Docker services (MySQL + Redis) then run the Go API server
- `make up` — `docker compose up -d` (start dependencies only)
- `make run` — `go run backend/cmd/api/main.go` (run API server standalone)
- `make down` — Stop and remove Docker containers
- `make logs` — Tail Docker compose logs
- `make restart` — Restart Docker services

Without Makefile:
- `go run ./backend/cmd/api/main.go` — Run API server from project root
- `cd backend && go run ./cmd/api/main.go` — Run from backend directory

### Build
No build scripts configured — standard `go build ./cmd/api/main.go` in `backend/`.

### Dependencies
- `cd backend && go mod tidy` — Sync Go module dependencies
- `cd backend && go mod download` — Download all dependencies

### Testing
No test files or test infrastructure currently exist in the project.

## Project Overview

**CatDiary (猫咪日记)** — a full-stack diary management system with immersive Markdown writing support, multi-level data persistence, and LLM extensibility.

### Tech Stack
- **Backend**: Go 1.26, CloudWeGo Hertz (HTTP framework), GORM (ORM), go-redis, JWT auth
- **Infra**: Docker Compose (MySQL 8.0 + Redis 7.4)
- **Frontend**: Not yet implemented (planned: React 18 / Vue 3, Vite, Tailwind CSS)

### Backend Architecture (Monorepo)

The backend follows a standard 3-layer architecture inside `backend/internal/`:

```
backend/
├── cmd/api/main.go           — Entry point: loads config, inits DB, registers routes, starts server
├── internal/
│   ├── config/               — Env loading (godotenv)
│   ├── model/                — GORM domain models (User, Diary)
│   ├── repository/           — Data access layer (MySQL via GORM)
│   ├── service/              — Business logic layer
│   ├── handler/              — HTTP request/response layer (Hertz handlers)
│   ├── middleware/           — JWT auth middleware
│   └── router/               — Route registration (maps paths to handlers)
├── pkg/
│   └── utils/                — JWT token generate/parse
└── go.mod
```

### Layer Data Flow
`Handler (req/res) → Service (business logic) → Repository (data access) → GORM → MySQL`

- **Handlers** parse requests, validate input (via go-playground/validator), and call service layer
- **Services** contain business logic, return domain errors (`ErrUserNotFound`, `ErrDiaryNotFound`, etc.)
- **Repositories** perform CRUD via GORM, scoped by `user_id` for multi-tenancy safety
- **Models** define GORM entities with explicit `TableName()` for singular table names

### Key Design Decisions

1. **JWT auth**: Stateless token auth. `middleware.RequireAuth()` parses Bearer token and injects `user_id` into request context. Token expiry: 7 days.
2. **Validation**: Hertz server is configured with `go-playground/validator/v10` for struct tag validation (e.g., `validate:"required,min=3,max=50"`).
3. **Diary CRUD**: All diary operations are scoped to the authenticated user. PUT is full replacement, PATCH is partial update.
4. **Draft & Upload handlers exist** as stubs but use Redis/JSONL — check current implementation before modifying.
5. **Auto-migration**: GORM AutoMigrate runs on startup — schema changes happen automatically.

### Environment
- `.env` file at project root (gitignored) — requires `CATDIARY_MYSQL_DSN` (MySQL DSN)
- Optional `CATDIARY_ENV_FILE` env var to specify a custom env file path
- Default `.env` loads from cwd and parent directory

### API Routes (all under `/api/v1`)
| Group | Auth | Routes |
|---|---|---|
| `/healthz`, `/readyz` | No | GET (liveness/readiness) |
| `/auth` | Partial | POST register/login (public), POST logout + GET me (auth required) |
| `/users` | Yes | GET/PATCH me, PATCH me/password |
| `/diaries` | Yes | POST create, GET list, GET/:id, PUT/:id, PATCH/:id, DELETE/:id |
| `/drafts` | Yes | PUT/GET/DELETE diary (Redis-backed) |
| `/uploads` | Yes | POST create (file upload) |

### Commit Convention
Follow `.trae/rules/git-commit-message.md`: Chinese commit messages ending with "喵~", imperative tone, no emoji, title under 72 chars.
