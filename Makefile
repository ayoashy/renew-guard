.PHONY: help build run test clean docker-build docker-up docker-down migrate-up migrate-down

# Default target
help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application locally"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down   - Stop Docker containers"
	@echo "  make deps          - Download dependencies"
	@echo "  make tidy          - Tidy go modules"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/renew-guard cmd/app/main.go

# Run the application
run:
	@echo "Running application..."
	go run cmd/app/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download

# Tidy go modules
tidy:
	@echo "Tidying go modules..."
	go mod tidy

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker-compose build

# Start Docker containers
docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d

# Stop Docker containers
docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

# View Docker logs
docker-logs:
	docker-compose logs -f app

# Restart application container
docker-restart:
	docker-compose restart app
