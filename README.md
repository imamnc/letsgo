# LetsGO

LetsGO is a Go API boilerplate built with Fiber, PostgreSQL, Redis, SQLC, JWT authentication, Swagger documentation, and a simple cron scheduler. It is designed to be a practical starting point for REST APIs that need a clean module structure, generated database access, and a ready-to-use development workflow.

The project follows a module-based architecture under `internal/modules`, where each feature owns its own handler, service, repository, routes, and module wiring. Shared concerns such as configuration, JWT, Redis, response helpers, and middleware live in the `shared` and `internal` packages.

## Highlights

- Fiber-based HTTP API with versioned routing under `/api/v1`
- SQLC-generated database access layer for strongly typed PostgreSQL queries
- JWT-based authentication with protected routes and bearer token middleware
- Redis service wrapper for caching and distributed coordination
- Swagger UI and OpenAPI docs generated from annotations
- Database migrations and seed commands ready for local development
- Cron scheduler bootstrap for recurring jobs

## Tech Stack

- Go 1.25.3
- Fiber v2
- PostgreSQL with `pgx` and SQLC
- Redis
- JWT
- Swagger / Swaggo
- `robfig/cron` for background jobs

## Project Structure

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

## Architecture

Each feature module is organized around the same structure:

- `handler.go` handles HTTP request and response logic.
- `service.go` contains business rules.
- `repository.go` talks to SQLC queries or persistence details.
- `routes.go` binds endpoints to Fiber routers.
- `module.go` wires the dependencies together.

The application is assembled in `cmd/api/main.go` by creating the app, registering routes, starting the scheduler, and then running the Fiber server.

## Available Modules

### Health

Provides a lightweight service check endpoint.

- `GET /api/v1/health`

### Auth

Handles token issuance, refresh flows, and current-user lookup.

- `POST /api/v1/auth/access-token`
- `POST /api/v1/auth/refresh-token`
- `GET /api/v1/auth/user`

The `GET /api/v1/auth/user` endpoint requires an `Authorization` header with a bearer access token.

### Users

Provides user CRUD operations and permission assignment endpoints.

- `GET /api/v1/users`
- `POST /api/v1/users`
- `GET /api/v1/users/:id`
- `PUT /api/v1/users/:id`
- `DELETE /api/v1/users/:id`
- `GET /api/v1/users/:id/permissions`
- `POST /api/v1/users/:id/permissions`
- `DELETE /api/v1/users/:id/permissions`
- `PUT /api/v1/users/:id/permissions`

### Permissions

Provides permission CRUD operations.

- `GET /api/v1/permissions`
- `POST /api/v1/permissions`
- `GET /api/v1/permissions/:id`
- `PUT /api/v1/permissions/:id`
- `DELETE /api/v1/permissions/:id`

## Configuration

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

## Getting Started

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

## Common Commands

The repository includes a `Makefile` with the following targets:

- `make dev` starts the development server with hot reload.
- `make build` compiles the production binary to `./bin/letsgo`.
- `make migration NAME=my_migration` creates a new SQL migration pair.
- `make migrate` applies pending migrations.
- `make migrate-down` rolls back the last migration batch.
- `make migrate-fresh` drops the database schema and re-applies migrations.
- `make seed` runs the database seed entrypoint.
- `make sqlc` regenerates SQLC output from the query files and migrations.
- `make docs` regenerates Swagger documentation.

## API Documentation

Swagger documentation is mounted under `/api/docs/*` after the application starts.

To regenerate the documentation assets, run:

```bash
make docs
```

## Database Layer

The database layer is based on SQLC and PostgreSQL migrations:

- SQL migration files live in `db/migrations`.
- Query definitions live in `db/sql`.
- Generated SQLC code is written to `db/sqlc`.

This keeps query access strongly typed while keeping the SQL itself easy to review.

## Authentication

Authentication is based on JWT access and refresh tokens.

- Tokens are issued from the auth module.
- Protected routes require an `Authorization` header.
- The middleware accepts both raw tokens and `Bearer <token>` format, but bearer format is recommended.

## Scheduler

The scheduler starts automatically with the API server. It currently runs a sample job every minute and logs the current users count using the configured application time zone. This is a placeholder for real recurring jobs and can be extended inside `internal/scheduler`.

## Extending the Boilerplate

To add a new feature module:

1. Create a new folder under `internal/modules/<feature>`.
2. Add `handler.go`, `service.go`, `repository.go`, `routes.go`, and `module.go`.
3. Register the module in `internal/routes/routes.go`.
4. Add or update SQL queries in `db/sql` if the module needs database access.
5. Run `make sqlc` if the query layer changes.

## Production Build

Build the application binary with:

```bash
make build
```

The resulting executable is written to `./bin/letsgo`.

## Notes

- The application expects PostgreSQL and Redis to be reachable from the configured host and port values.
- Swagger metadata is populated from the runtime configuration, so keep `APP_HOST` aligned with the address you want to expose.
- `APP_TIMEZONE` affects scheduler timestamps only and does not change database storage behavior.
