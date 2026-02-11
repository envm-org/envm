.PHONY: build run test clean migration migrate-up migrate-down migrate-status lint format install-hooks build-cli build-docs build-all run-docs help

# Load environment variables from .env
ifneq (,$(wildcard .env))
    include .env
    export
endif

BINARY_NAME=envm-server
MIGRATION_DIR=./internal/adapters/postgresql/schema/migrations
DB_DRIVER=postgres
DB_STRING=$(DATABASE_URI)

build:
	go build -o bin/$(BINARY_NAME) ./cmd

run: build
	./bin/$(BINARY_NAME)

test:
	go test -v ./...

clean:
	go clean
	rm -f bin/$(BINARY_NAME)

migration:
	@if [ -z "$(name)" ]; then echo "Error: name argument is required (e.g., make migration name=create_users)"; exit 1; fi
	go run github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATION_DIR) create $(name) sql

migrate-up:
	go run github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATION_DIR) $(DB_DRIVER) "$(DB_STRING)" up

migrate-down:
	go run github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATION_DIR) $(DB_DRIVER) "$(DB_STRING)" down

migrate-status:
	go run github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATION_DIR) $(DB_DRIVER) "$(DB_STRING)" status

dev:
	go run github.com/air-verse/air

security:
	govulncheck ./...

lint:
	golangci-lint run ./...

format:
	gofmt -w .

install-hooks:
	git config core.hooksPath githooks

build-cli:
	cd cli && go build -o ../bin/envm .

build-docs:
	cd docs && npm install && npm run build

build-all: build build-cli build-docs

run-docs:
	cd docs && npm run start

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build           Build the server binary"
	@echo "  run             Run the server locally"
	@echo "  test            Run tests"
	@echo "  clean           Clean build artifacts"
	@echo "  migration       Create a new migration (usage: make migration name=...)"
	@echo "  migrate-up      Run migrations up"
	@echo "  migrate-down    Run migrations down"
	@echo "  migrate-status  Check migration status"
	@echo "  dev             Run server with air (hot reload)"
	@echo "  security        Run vulnerability checks"
	@echo "  lint            Run linters"
	@echo "  format          Format code"
	@echo "  install-hooks   Install git hooks"
	@echo "  build-cli       Build the CLI binary"
	@echo "  build-docs      Build the documentation"
	@echo "  build-all       Build everything (server, cli, docs)"
	@echo "  run-docs        Run the documentation server"
