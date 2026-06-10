BINARY := ./bin/letsgo
CMD := ./cmd/api
NAME ?= create_users

.PHONY: up dev build sqlc migration migrate migrate-down

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

sqlc:
	@echo "Generating sqlc code..."
	sqlc generate