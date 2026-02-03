.PHONY: run build test clean dev help install-deps

# Default target
help:
    @echo "Available commands:"
    @echo "  make run          - Run the application"
    @echo "  make build        - Build the application"
    @echo "  make test         - Run tests"
    @echo "  make clean        - Remove build artifacts"
    @echo "  make dev          - Run with hot reload (requires air)"
    @echo "  make install-deps - Install Go dependencies"

# Run the application
run:
    go run main.go

# Build the application
build:
    @echo "Building..."
    go build -o kasir-api.exe main.go
    @echo "Build complete: kasir-api.exe"

# Run tests
test:
    go test -v ./...

# Clean build artifacts
clean:
    @echo "Cleaning..."
    @if exist kasir-api.exe del kasir-api.exe
    @echo "Clean complete"

# Install dependencies
install-deps:
    go mod download
    go mod tidy

# Development mode with hot reload
dev:
    @echo "Starting development server..."
    air