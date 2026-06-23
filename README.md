<p align="center">
  <img src="letsgo.png" alt="LetsGO Logo" width="220" />
</p>

<!-- <h1 align="center">LetsGO</h1> -->

<p align="center">
  <strong>Go API Boilerplate</strong> · Modular, production-ready, and built with Fiber, PostgreSQL, Redis, SQLC, JWT, and Swagger.
</p>

<p align="center">
  <a href="#highlights">Highlights</a> ·
  <a href="#project-structure">Structure</a> ·
  <a href="#makefile-commands-and-usage">Commands</a>
</p>

---

LetsGO is a Go API boilerplate built with Fiber, PostgreSQL, Redis, SQLC, JWT authentication, Swagger documentation, and a simple cron scheduler. It is designed to be a practical starting point for REST APIs that need a clean module structure, generated database access, and a ready-to-use development workflow.

The project follows a module-based architecture under `internal/modules`, where each feature owns its own handler, service, repository, routes, and module wiring. Shared concerns such as configuration, JWT, Redis, response helpers, and middleware live in the `shared` and `internal` packages.

<h2 id="highlights">Highlights ✨</h2>

- Fiber-based HTTP API with versioned routing under `/api/v1`
- SQLC-generated database access layer for strongly typed PostgreSQL queries
- JWT-based authentication with protected routes and bearer token middleware
- Redis service wrapper for caching and distributed coordination
- Swagger UI and OpenAPI docs generated from annotations
- Built-in health, auth, user, and permission modules with route scaffolding
- Database migrations and seed commands ready for local development
- Cron scheduler bootstrap for recurring jobs

## Tech Stack 🧰

- Go 1.25.3
- Fiber v2
- PostgreSQL with `pgx` and SQLC
- Redis
- JWT
- Swagger / Swaggo
- `robfig/cron` for background jobs

<h2 id="project-structure">Project Structure 📁</h2>

```text
cmd/
	api/         Main HTTP server entrypoint
	seed/        Database seed entrypoint
db/
	migrations/  SQL migrations
	sql/         Raw SQL used by SQLC
	sqlc/        Generated Go code from SQLC
internal/
	app/         Application bootstrap and shared dependencies
	config/      Environment loading and configuration
	middleware/  Cross-cutting HTTP middleware
	modules/     Feature modules (auth, health, permission, user)
	routes/      Global route registration and Swagger metadata
	scheduler/   Cron scheduler bootstrap and sample task
shared/
	env/         Environment helpers
	jwt/         JWT provider and claims types
	redis/       Redis client and service wrapper
	response/    API response helpers
	json/        JSON utilities
	format/      Formatting helpers
```

## Architecture 🏗️

Each feature module is organized around the same structure:

- `handler.go` handles HTTP request and response logic.
- `service.go` contains business rules.
- `repository.go` talks to SQLC queries or persistence details.
- `routes.go` binds endpoints to Fiber routers.
- `module.go` wires the dependencies together.

The application is assembled in `cmd/api/main.go` by creating the app, registering routes, starting the scheduler, and then running the Fiber server.

## API Documentation 📚

Swagger documentation is mounted under `/api/docs/*` after the application starts.

To regenerate the documentation assets, run:

```bash
make docs
```

## Configuration ⚙️

Copy `.env.example` to `.env` and adjust the values for your local environment.

| Variable | Description | Default |
| --- | --- | --- |
| `APP_PORT` | HTTP server port | `3000` |
| `APP_HOST` | Public host used in Swagger metadata | `localhost:3000` |
| `APP_TIMEZONE` | Time zone used by scheduled jobs | `UTC` |
| `DB_HOST` | PostgreSQL host | `127.0.0.1` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USERNAME` | PostgreSQL username | `postgres` |
| `DB_PASSWORD` | PostgreSQL password | `password` |
| `DB_DATABASE` | PostgreSQL database name | `letsgo` |
| `DB_SSLMODE` | PostgreSQL SSL mode | `disable` |
| `REDIS_HOST` | Redis host | `127.0.0.1` |
| `REDIS_PORT` | Redis port | `6379` |
| `JWT_SECRET` | Signing secret for JWT tokens | `secret` |
| `JWT_ISSUER` | JWT issuer value | `letsgo` |
| `JWT_EXPIRY` | Access token expiry in seconds | `3600` |

## Getting Started 🚀

### Prerequisites

- Go 1.25.3 or newer
- PostgreSQL
- Redis
- `make`
- Optional tooling for local workflows: `air`, `sqlc`, `migrate`, and `swag`

### Local Setup

1. Clone the repository.
2. Create your environment file.

```bash
cp .env.example .env
```

3. Start PostgreSQL and Redis using your preferred local setup.
4. Run database migrations.

```bash
make migrate
```

5. Seed the initial data set.

```bash
make seed
```

6. Start the API in development mode.

```bash
make dev
```

<h2 id="makefile-commands-and-usage">Makefile Commands and Usage 🛠️</h2>

The repository exposes the following `make` commands as the primary workflow for development and maintenance.

- `make dev`
  - Starts the development server with hot reload using the `.air.toml` configuration.
  - Use this command during active development to see code changes without restarting the server manually.

- `make build`
  - Builds the production binary at `./bin/letsgo`.
  - Use this when packaging or deploying the application.

- `make migration NAME=my_migration`
  - Generates a new migration file pair under `db/migrations`.
  - Replace `my_migration` with a descriptive name such as `add_users_table`.

- `make migrate`
  - Applies pending database migrations using values from `.env`.
  - Ensure `.env` is present and `DB_*` variables are configured before running.

- `make migrate-down`
  - Reverts the last migration batch.
  - Use this to roll back changes during development or testing.

- `make migrate-fresh`
  - Drops the database schema and re-applies all migrations from scratch.
  - Use cautiously; this resets database state.

- `make seed`
  - Runs the seed entrypoint at `cmd/seed` to insert sample data.
  - Useful after migrations to populate required baseline records.

- `make sqlc`
  - Regenerates SQLC code from `db/sql` and migration schema.
  - Run this whenever SQL query files or schema definitions change.

- `make docs`
  - Generates Swagger documentation from annotations in `cmd/api/main.go`.
  - Use this to refresh the API docs served under `/api/docs/*`.

### Usage Guidelines

- Always create or refresh `.env` from `.env.example` before running environment-dependent commands.
- Use `make dev` for local API development and `make build` for production packaging.
- Apply migrations with `make migrate` before running the server, then seed data with `make seed` if needed.
- For database schema changes, create a migration, then run `make migrate` and `make sqlc` if the query layer is affected.
- If you need to reset local state, `make migrate-fresh` is the safest choice for development, but do not use it in production.

## Database Layer 🗄️

The database layer is based on SQLC and PostgreSQL migrations:

- SQL migration files live in `db/migrations`.
- Query definitions live in `db/sql`.
- Generated SQLC code is written to `db/sqlc`.

This keeps query access strongly typed while keeping the SQL itself easy to review.

## Authentication 🔐

Authentication is based on JWT access and refresh tokens.

- Tokens are issued from the auth module.
- Protected routes require an `Authorization` header.
- The middleware accepts both raw tokens and `Bearer <token>` format, but bearer format is recommended.

## Scheduler ⏱️

The scheduler starts automatically with the API server. It currently runs a sample job every minute and logs the current users count using the configured application time zone. This is a placeholder for real recurring jobs and can be extended inside `internal/scheduler`.

## Extending the Boilerplate 🧩

To add a new feature module:

1. Create a new folder under `internal/modules/<feature>`.
2. Add `handler.go`, `service.go`, `repository.go`, `routes.go`, and `module.go`.
3. Register the module in `internal/routes/routes.go`.
4. Add or update SQL queries in `db/sql` if the module needs database access.
5. Run `make sqlc` if the query layer changes.

## Production Build 🚢

Build the application binary with:

```bash
make build
```

The resulting executable is written to `./bin/letsgo`.

## Notes 📝

- The application expects PostgreSQL and Redis to be reachable from the configured host and port values.
- Swagger metadata is populated from the runtime configuration, so keep `APP_HOST` aligned with the address you want to expose.
- `APP_TIMEZONE` affects scheduler timestamps only and does not change database storage behavior.
