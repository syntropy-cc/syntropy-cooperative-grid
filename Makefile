# Syntropy Cooperative Grid - Makefile

.PHONY: help dev-setup dev-start dev-stop test lint clean build deploy

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
dev-setup: ## Setup development environment
	@echo "Setting up development environment..."
	@./scripts/development/local-dev-setup.sh

dev-start: ## Start local development services
	@echo "Starting development services..."
	@docker-compose -f docker-compose.dev.yml up -d

dev-stop: ## Stop local development services
	@echo "Stopping development services..."
	@docker-compose -f docker-compose.dev.yml down

# Testing
test: ## Run all tests
	@echo "Running tests..."
	@./scripts/development/run-tests.sh

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	@go test ./... -v

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@./scripts/testing/integration-tests/run-integration-tests.sh

# Code quality
lint: ## Run linters
	@echo "Running linters..."
	@golangci-lint run
	@terraform fmt -check -recursive
	@ansible-lint infrastructure/ansible/

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@terraform fmt -recursive
	@black --line-length 88 .

# Build
build: ## Build all components
	@echo "Building all components..."
	@./scripts/development/build-images.sh

build-genesis: ## Build genesis node image
	@echo "Building genesis node..."
	@docker build -t syntropy/genesis-node:latest -f build/docker/genesis-node/Dockerfile .

# Infrastructure
terraform-init: ## Initialize Terraform
	@cd infrastructure/terraform/environments/genesis && terraform init

terraform-plan: ## Plan Terraform changes
	@cd infrastructure/terraform/environments/genesis && terraform plan

terraform-apply: ## Apply Terraform changes
	@cd infrastructure/terraform/environments/genesis && terraform apply

# Deployment
deploy-genesis: ## Deploy genesis node
	@echo "Deploying genesis node..."
	@./scripts/deployment/deploy-genesis.sh

deploy-worker: ## Deploy worker node
	@echo "Deploying worker node..."
	@./scripts/deployment/add-worker-node.sh

# Monitoring
logs: ## View logs from all services
	@echo "Viewing logs..."
	@docker-compose -f docker-compose.dev.yml logs -f

status: ## Check status of all services
	@echo "Checking service status..."
	@kubectl get pods -A

# Cleanup
clean: ## Clean build artifacts and temporary files
	@echo "Cleaning up..."
	@rm -rf build/temp/*
	@rm -rf .terraform/
	@docker system prune -f

clean-all: ## Clean everything including volumes
	@echo "Cleaning everything..."
	@docker-compose -f docker-compose.dev.yml down -v
	@docker system prune -af --volumes

# Security
security-scan: ## Run security scans
	@echo "Running security scans..."
	@./scripts/security/security-scan.sh

vulnerability-check: ## Check for vulnerabilities
	@echo "Checking vulnerabilities..."
	@./scripts/security/vulnerability-check.sh

# Documentation
docs-serve: ## Serve documentation locally
	@echo "Serving documentation..."
	@cd web/docs-site && hugo server

docs-build: ## Build documentation
	@echo "Building documentation..."
	@cd web/docs-site && hugo

# Version management
version: ## Show current version
	@cat VERSION

version-bump: ## Bump version (usage: make version-bump VERSION=0.2.0)
	@echo $(VERSION) > VERSION
	@git add VERSION
	@git commit -m "chore: bump version to $(VERSION)"
	@git tag v$(VERSION)
