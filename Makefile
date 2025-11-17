.PHONY: proto test build lint docker clean help

# Build variables
BINARY_NAME=majhong_cli
PROTO_DIR=proto
PKG_PROTO_DIR=pkg/proto
GO_FILES=$(shell find . -name '*.go' -type f)

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

proto: ## Generate Go code from protobuf files
	@echo "Generating protobuf code..."
	@mkdir -p $(PKG_PROTO_DIR)
	@protoc --go_out=$(PKG_PROTO_DIR) --go_opt=paths=source_relative \
		--proto_path=$(PROTO_DIR) $(PROTO_DIR)/*.proto
	@echo "Protobuf code generated successfully"

test: ## Run all tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | tail -1

test-contract: ## Run contract tests only
	@echo "Running contract tests..."
	@go test -v ./test/contract/...

test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	@go test -v ./test/integration/...

build: ## Build the CLI binary
	@echo "Building $(BINARY_NAME)..."
	@go build -o bin/$(BINARY_NAME) cmd/majhong_cli/main.go
	@echo "Binary built: bin/$(BINARY_NAME)"

lint: ## Run golangci-lint
	@echo "Running linter..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@gofmt -s -w $(GO_FILES)

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

tidy: ## Tidy go modules
	@echo "Tidying modules..."
	@go mod tidy

docker: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t majhong_cli:latest .
	@echo "Docker image built: majhong_cli:latest"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out
	@echo "Clean complete"

install-deps: ## Install development dependencies
	@echo "Installing dependencies..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go get -u github.com/alecthomas/kong
	@go get -u github.com/gdamore/tcell/v2
	@go get -u google.golang.org/protobuf
	@go mod tidy
	@echo "Dependencies installed"

all: lint test build ## Run lint, test, and build
