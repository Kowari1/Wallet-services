.PHONY: all build test clean run help docker-up docker-down swagger migrate

BIN_DIR=bin
APP_NAME=wallet-app
DOCKER_COMPOSE=docker-compose.yml

all: build

build:
	@echo "Building all services..."
	go build -o ../$(BIN_DIR)/gw-currency-wallet ./cmd/main.go
	@echo "Build complete!"

test:
	@echo "Running tests..."
	go test ./internal/... -v

docker-up:
	docker compose -f $(DOCKER_COMPOSE) up -d

docker-down:
	docker compose -f $(DOCKER_COMPOSE) down

swagger:
	swag init -g cmd/main.go -o internal/docs

migrate-up:
	migrate -path internal/migrations -database "postgres://user:pass@localhost:5432/wallet?sslmode=disable" up

migrate-down:
	migrate -path internal/migrations -database "postgres://user:pass@localhost:5432/wallet?sslmode=disable" down

run-wallet:
	go run cmd/main.go

clean:
	@echo "Cleaning binaries..."
	rm -rf $(BIN_DIR)

help:
	@echo "Available commands:"
	@echo "  make build       - Build all services"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Clean binaries"
	@echo "  make docker-up   - Start Docker containers"
	@echo "  make docker-down - Stop Docker containers"
	@echo "  make swagger     - Generate Swagger docs"
	@echo "  make migrate     - Run database migrations"
	@echo "  make run-wallet  - Run wallet service"
