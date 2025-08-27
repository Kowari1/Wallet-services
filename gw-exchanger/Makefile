.PHONY: all build test clean run help docker-up docker-down swagger migrate

BIN_DIR=bin
APP_NAME=exchanger-appmake

all: build

build:
	@echo "Building all services..."
	go build -o ../$(BIN_DIR)/gw-exchanger ./cmd/main.go
	@echo "Build complete!"

run-exchanger:
	go run cmd/main.go

