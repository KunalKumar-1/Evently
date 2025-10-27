.PHONY: help install deps migrate up down dev run build clean test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install Go dependencies
	go mod download
	go mod tidy

deps: install ## Alias for install

migrate-up: ## Run database migrations
	@echo "Running migrations..."
	go run cmd/migrate/main.go up

migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	go run cmd/migrate/main.go down

dev: ## Start development server with Air
	@echo "Starting development server with Air..."
	air

run: ## Run the application
	@echo "Starting application..."
	go run cmd/api/main.go

build: ## Build the application
	@echo "Building application..."
	go build -o bin/api cmd/api/main.go

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf tmp/
	rm -rf bin/
	go clean

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

.DEFAULT_GOAL := help
