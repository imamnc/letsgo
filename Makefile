BINARY := ./bin/letsgo
CMD := ./cmd/api
# Migration name for creating new migrations
NAME ?= create_users

# ANSI escape codes for colored output
ESC := \033
BOLD := $(ESC)[1m
CYAN := $(ESC)[36m
GREEN := $(ESC)[32m
YELLOW := $(ESC)[33m
MAGENTA := $(ESC)[35m
RESET := $(ESC)[0m

# Define phony targets to prevent conflicts with files of the same name
.PHONY: dev build sqlc migration migrate migrate-down migrate-fresh seed

# =============================================================
# Default target when running `make` without arguments
# =============================================================

dev:
	@printf "$(CYAN)$(BOLD)╔══════════════════════════════════════════════╗$(RESET)\n"
	@printf "$(CYAN)$(BOLD)║                  LetsGO                      ║$(RESET)\n"
	@printf "$(CYAN)$(BOLD)║           Just clone and LetsGo              ║$(RESET)\n"
	@printf "$(CYAN)$(BOLD)╚══════════════════════════════════════════════╝$(RESET)\n"
	@printf "$(CYAN)$(BOLD)🚀 Starting development server with hot reload...$(RESET)\n"
	air -c .air.toml

build:
	@printf "$(GREEN)📦 Building binary $(BINARY)...$(RESET)\n"
	@mkdir -p $(dir $(BINARY))
	go build -o $(BINARY) $(CMD)

migration:
	@printf "$(MAGENTA)🧱 Creating migration $(NAME)...$(RESET)\n"
	migrate create -ext sql -dir db/migrations $(NAME)

migrate:
	@printf "$(GREEN)⬆️  Running database migrations...$(RESET)\n"
	@set -a; . ./.env 2>/dev/null || true; set +a; \
	migrate -path db/migrations -database "postgres://$$DB_USERNAME:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_DATABASE?sslmode=$$DB_SSLMODE" up

migrate-down:
	@printf "$(YELLOW)⬇️  Reverting database migrations...$(RESET)\n"
	@set -a; . ./.env 2>/dev/null || true; set +a; \
	migrate -path db/migrations -database "postgres://$$DB_USERNAME:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_DATABASE?sslmode=$$DB_SSLMODE" down

migrate-fresh:
	@printf "$(YELLOW)♻️  Resetting database and running migrations...$(RESET)\n"
	@set -a; . ./.env 2>/dev/null || true; set +a; \
	migrate -path db/migrations -database "postgres://$$DB_USERNAME:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_DATABASE?sslmode=$$DB_SSLMODE" drop -f && \
	migrate -path db/migrations -database "postgres://$$DB_USERNAME:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_DATABASE?sslmode=$$DB_SSLMODE" up

seed:
	@printf "$(GREEN)🌱 Seeding database...$(RESET)\n"
	@set -a; . ./.env 2>/dev/null || true; set +a; \
	go run ./cmd/seed

sqlc:
	@printf "$(CYAN)🛠️ Generating sqlc code...$(RESET)\n"
	sqlc generate