.PHONY: build run test clean migration migrate-up migrate-down migrate-status

# Load environment variables from .env
ifneq (,$(wildcard .env))
    include .env
    export
endif

BINARY_NAME=envm
MIGRATION_DIR=./internal/adapters/postgresql/schema/migrations
DB_DRIVER=postgres
DB_STRING=$(DATABASE_URI)

build:
	go build -o bin/$(BINARY_NAME) ./internal/cmd/main.go

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
