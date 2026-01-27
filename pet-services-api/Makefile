docs:
	@echo "Gerando documentação Swagger..."
	swag init --dir ./cmd/api,./internal/transport/http --parseDependency --parseInternal --output ./internal/transport/http/docs
.PHONY: help build run test migrate migrate-down clean doc
doc:
	@echo "Gerando documentação Swagger..."
	swag init -g cmd/api/main.go

help:
	@echo "Available targets:"
	@echo "  make build        - Build the API binary"
	@echo "  make run          - Run the API server"
	@echo "  make test         - Run all tests"
	@echo "  make migrate      - Run database migrations"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make lint         - Run linters (gofmt, go vet)"
	@echo "  make fmt          - Format all Go files"

build:
	@echo "Building API..."
	@go build -o bin/api ./cmd/api

run: build
	@echo "Starting API server..."
	@./bin/api

test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

test-coverage: test
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

migrate:
	@echo "Running migrations..."
	@go run ./cmd/migrate

migrate-down:
	@echo "⚠️  Warning: Migration rollback not implemented."

lint:
	@echo "Linting..."
	@gofmt -l -w ./...
	@go vet ./...

fmt:
	@echo "Formatting Go files..."
	@gofmt -s -w ./cmd ./internal

docker-up:
	@echo "Starting PostgreSQL container..."
	@docker-compose up -d postgres

docker-down:
	@echo "Stopping PostgreSQL container..."
	@docker-compose down

docker-logs:
	@docker-compose logs -f postgres

clean:
	@echo "Cleaning..."
	@rm -f bin/api coverage.out coverage.html
	@go clean

install-deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

setup: install-deps docker-up migrate
	@echo "✅ Setup complete! Run 'make run' to start the API."
