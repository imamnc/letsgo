BINARY := ./bin/letsgo
CMD := ./cmd/api

.PHONY: up dev build

up: dev

dev:
	@echo "Starting development server with hot reload..."
	air -c .air.toml

build:
	@mkdir -p $(dir $(BINARY))
	go build -o $(BINARY) $(CMD)
