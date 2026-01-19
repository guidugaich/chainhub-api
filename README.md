# ChainHub API

Backend for the ChainHub link-in-bio platform. This project is intentionally small, explicit, and readable.

## What this backend does

- Email/password auth with JWT sessions
- Public tree page by username
- Authenticated CRUD for links
- PostgreSQL for persistence
- Docker Compose for API + Postgres

## Project structure

- `cmd/api/main.go`: entrypoint and HTTP server startup
- `internal/config/config.go`: environment config and defaults
- `internal/db/db.go`: PostgreSQL connection setup
- `internal/http/router.go`: routing and middleware
- `internal/http/handlers/`: request handlers
- `internal/http/middleware/auth.go`: JWT auth middleware
- `internal/models/`: DB models
- `internal/repo/`: SQL queries and data access
- `internal/services/jwt.go`: JWT creation
- `migrations/001_init.sql`: database schema
- `Dockerfile`: API container build
- `docker-compose.yml`: API + Postgres services

## Setup (step by step)

1. Install Go 1.22+ and Docker Desktop.
2. Download dependencies:
   - `go mod tidy`
3. Start the services:
   - `docker compose up --build`
4. Migrations run automatically on startup. To run them manually:
   - `docker compose --profile migrate run --rm migrate up`

## Reset the database (ephemeral workflow)

- `docker compose down -v`

This removes the Postgres volume so you can start from a clean database.

## Environment variables

The API reads configuration from environment variables:

- `APP_ENV` (default: `development`)
- `PORT` (default: `8080`)
- `DB_HOST` (default: `localhost`)
- `DB_PORT` (default: `5432`)
- `DB_USER` (default: `chainhub`)
- `DB_PASSWORD` (default: `chainhub`)
- `DB_NAME` (default: `chainhub`)
- `DB_SSLMODE` (default: `disable`)
- `JWT_SECRET` (required)
- `RUN_MIGRATIONS` (default: `true`)
- `MIGRATIONS_PATH` (default: `file://migrations`)

`docker-compose.yml` already supplies these for local dev.

## Migrations

Use the `migrate` CLI via Docker Compose (optional; the API runs migrations by default):

- Apply all: `docker compose --profile migrate run --rm migrate up`
- Roll back one: `docker compose run --rm migrate down 1`
- Create new migration: `docker compose run --rm migrate create -ext sql -dir /migrations -seq add_feature`

Or use the Makefile shortcuts:

- `make migrate-up`
- `make migrate-down`
- `make migrate-create name=add_feature`

## Endpoints

- `POST /signup` `{ "email": "...", "password": "..." }`
- `POST /login` `{ "email": "...", "password": "..." }`
- `GET /tree/{username}`
- `GET /links?tree_id=...` (auth required)
- `POST /links` (auth required)
- `PUT /links/{id}` (auth required)
- `DELETE /links/{id}` (auth required)

## Request flow (high level)

`main.go` loads config → connects to DB → creates router → starts HTTP server → router runs middleware → handler validates input → repo runs SQL → handler writes JSON response.
