# Makefile for g0 load tester

# Binary name
BINARY_NAME=g0

# Build directory
BUILD_DIR=.

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

.PHONY: all build clean run test help

# Default target
all: clean build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "Build complete!"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@if [ -f $(BUILD_DIR)/$(BINARY_NAME) ]; then \
		rm -f $(BUILD_DIR)/$(BINARY_NAME); \
		echo "Removed $(BINARY_NAME)"; \
	else \
		echo "No binary to clean"; \
	fi

# Run the load test
run:
	@echo "Running load test..."
	./$(BINARY_NAME) run --url https://httpbin.org/get --concurrency 50 --duration 10s

# Clean, build, and run test
test: clean build run
	@echo "Test complete!"

# Help target
help:
	@echo "Available targets:"
	@echo "  make build  - Build the application"
	@echo "  make clean  - Remove build artifacts"
	@echo "  make run    - Run the load test"
	@echo "  make test   - Clean, build, and run test"
	@echo "  make all    - Clean and build"
	@echo "  make help   - Show this help message"

