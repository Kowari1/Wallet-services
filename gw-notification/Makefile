.PHONY: all build test clean run help docker-up docker-down swagger migrate

BIN_DIR=bin
APP_NAME=notification-app

all: build

build:
	@echo "Building all services..."
	go build -o ../$(BIN_DIR)/gw-notification ./cmd/main.go
	@echo "Build complete!"

run-notification:
	go run cmd/main.go

