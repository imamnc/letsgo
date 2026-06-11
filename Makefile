BINARY := ./bin/letsgo
CMD := ./cmd/api
# Migration name for creating new migrations
NAME ?= create_users
# Default user details for seeding
USER_NAME ?= admin
USER_EMAIL ?= admin@example.com
USER_PASSWORD ?= password

.PHONY: up dev build sqlc migration migrate migrate-down user

up: dev

dev:
	@echo "Starting development server with hot reload..."
	air -c .air.toml

build:
	@mkdir -p $(dir $(BINARY))
	go build -o $(BINARY) $(CMD)

migration:
	@echo "Creating migration $(NAME)..."
	migrate create -ext sql -dir migrations $(NAME)

migrate:
	@echo "Running database migrations..."
	@set -a; . ./.env 2>/dev/null || true; set +a; \
	migrate -path migrations -database "postgres://$$DB_USERNAME:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_DATABASE?sslmode=$$DB_SSLMODE" up

migrate-down:
	@echo "Reverting database migrations..."
	@set -a; . ./.env 2>/dev/null || true; set +a; \
	migrate -path migrations -database "postgres://$$DB_USERNAME:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_DATABASE?sslmode=$$DB_SSLMODE" down

user:
	@echo "Seeding user $(USER_EMAIL)..."
	@set -a; . ./.env 2>/dev/null || true; set +a; \
	USER_NAME="$(USER_NAME)" USER_EMAIL="$(USER_EMAIL)" USER_PASSWORD="$(USER_PASSWORD)" go run ./cmd/seed

sqlc:
	@echo "Generating sqlc code..."
	sqlc generate