# Syntropy CLI Manager - Makefile

.PHONY: help build test clean install uninstall deps fmt lint vet cross-build

# Configuration
APP_NAME=syntropy
BUILD_DIR=build
MAIN_GO_FILE=main.go
VERSION=$(shell date +%Y%m%d-%H%M%S)
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
deps: ## Install dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies updated."

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted."

lint: ## Run linter
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

# Testing
test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -cover ./...

test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	go test -v -race ./...

# Building
build: ## Build for current platform
	@echo "Building $(APP_NAME) for $(shell go env GOOS)/$(shell go env GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_GO_FILE)
	@echo "Build complete. Executable: $(BUILD_DIR)/$(APP_NAME)"

build-linux: ## Build for Linux
	@echo "Building $(APP_NAME) for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(MAIN_GO_FILE)
	@echo "Build complete. Executable: $(BUILD_DIR)/$(APP_NAME)-linux-amd64"

build-windows: ## Build for Windows
	@echo "Building $(APP_NAME) for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(MAIN_GO_FILE)
	@echo "Build complete. Executable: $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe"

build-darwin: ## Build for macOS
	@echo "Building $(APP_NAME) for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(MAIN_GO_FILE)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 $(MAIN_GO_FILE)
	@echo "Build complete. Executables: $(BUILD_DIR)/$(APP_NAME)-darwin-*"

cross-build: build-linux build-windows build-darwin ## Build for all platforms

# Installation
install: build ## Install the binary
	@echo "Installing $(APP_NAME)..."
	sudo cp $(BUILD_DIR)/$(APP_NAME) /usr/local/bin/
	@echo "Installation complete."

uninstall: ## Uninstall the binary
	@echo "Uninstalling $(APP_NAME)..."
	sudo rm -f /usr/local/bin/$(APP_NAME)
	@echo "Uninstallation complete."

# Packaging
package: cross-build ## Create distribution packages
	@echo "Creating distribution packages..."
	@mkdir -p $(BUILD_DIR)/packages
	
	# Linux package
	@if [ -f $(BUILD_DIR)/$(APP_NAME)-linux-amd64 ]; then \
		cd $(BUILD_DIR) && tar -czf packages/$(APP_NAME)-linux-amd64-$(VERSION).tar.gz $(APP_NAME)-linux-amd64; \
		echo "Linux package: $(BUILD_DIR)/packages/$(APP_NAME)-linux-amd64-$(VERSION).tar.gz"; \
	fi
	
	# Windows package
	@if [ -f $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe ]; then \
		cd $(BUILD_DIR) && zip packages/$(APP_NAME)-windows-amd64-$(VERSION).zip $(APP_NAME)-windows-amd64.exe; \
		echo "Windows package: $(BUILD_DIR)/packages/$(APP_NAME)-windows-amd64-$(VERSION).zip"; \
	fi
	
	# macOS packages
	@if [ -f $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 ]; then \
		cd $(BUILD_DIR) && tar -czf packages/$(APP_NAME)-darwin-amd64-$(VERSION).tar.gz $(APP_NAME)-darwin-amd64; \
		echo "macOS Intel package: $(BUILD_DIR)/packages/$(APP_NAME)-darwin-amd64-$(VERSION).tar.gz"; \
	fi
	@if [ -f $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 ]; then \
		cd $(BUILD_DIR) && tar -czf packages/$(APP_NAME)-darwin-arm64-$(VERSION).tar.gz $(APP_NAME)-darwin-arm64; \
		echo "macOS Apple Silicon package: $(BUILD_DIR)/packages/$(APP_NAME)-darwin-arm64-$(VERSION).tar.gz"; \
	fi
	
	@echo "Packaging complete."

# Quality checks
quality: fmt vet lint test ## Run all quality checks

# Development workflow
dev-setup: deps ## Setup development environment
	@echo "Development environment setup complete."
	@echo "Run 'make test' to verify everything works."

dev: build test ## Build and test for development

# Cleanup
clean: ## Clean build artifacts
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)
	@echo "Cleanup complete."

# Information
version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"

info: ## Show build information
	@echo "App Name: $(APP_NAME)"
	@echo "Main File: $(MAIN_GO_FILE)"
	@echo "Build Dir: $(BUILD_DIR)"
	@echo "Go Version: $(shell go version)"
	@echo "Go OS/Arch: $(shell go env GOOS)/$(shell go env GOARCH)"

# Quick commands
quick: build test ## Quick build and test

all: deps quality cross-build package ## Do everything (deps, quality, build, package)

# Example usage commands
run: build ## Build and run with arguments
	@echo "Running $(APP_NAME)..."
	./$(BUILD_DIR)/$(APP_NAME) $(ARGS)

# Example usage: make run ARGS="setup --help"
# Example usage: make run ARGS="setup run --force"
