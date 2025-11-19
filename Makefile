.PHONY: build test proto lint clean help

# Variables
BINARY_NAME := majhong_cli
BINARY_PATH := bin/$(BINARY_NAME)
MAIN_PACKAGE := ./cmd/majhong_cli
GO := go
PROTOC := protoc
GOLANGCI_LINT := golangci-lint

# Default target
.DEFAULT_GOAL := help

# Build the CLI binary
build:
	@echo "Building $(BINARY_NAME)..."
	mkdir -p bin/
	$(GO) build -o $(BINARY_PATH) $(MAIN_PACKAGE)
	@echo "Build complete: $(BINARY_PATH)"

# Run all tests with verbose output and race detector
test:
	@echo "Running tests..."
	$(GO) test -v -race ./...

# Generate Go code from protobuf schemas
proto:
	@echo "Generating protobuf code..."
	@if [ ! -d "proto" ]; then \
		echo "Warning: proto/ directory not found"; \
	fi
	@mkdir -p pkg/mahjong/pb
	$(PROTOC) --go_out=pkg/mahjong/pb proto/*.proto 2>/dev/null || echo "No proto files found or protoc not available"

# Run golangci-lint
lint:
	@echo "Running linter..."
	$(GOLANGCI_LINT) run

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	@echo "Clean complete"

# Display available targets
help:
	@echo "CLI Mahjong Engine - Makefile Targets"
	@echo "======================================"
	@echo ""
	@echo "Available targets:"
	@echo "  build   - Build the CLI binary (output: bin/majhong_cli)"
	@echo "  test    - Run all tests with race detector"
	@echo "  proto   - Generate Go code from protobuf schemas"
	@echo "  lint    - Run golangci-lint"
	@echo "  clean   - Clean build artifacts"
	@echo "  help    - Display this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build      # Build the binary"
	@echo "  make test       # Run tests"
	@echo "  make proto      # Generate protobuf code"
	@echo "  make clean      # Clean artifacts"
	@echo ""
