# Syntropy Cooperative Grid Management System
# Makefile for development and deployment

.PHONY: help install build test clean docker up down deploy

# Default target
help: ## Show this help message
	@echo "Syntropy Cooperative Grid Management System"
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development
install: ## Install dependencies
	@echo "Installing Core dependencies..."
	cd core && go mod download && go mod tidy
	@echo "Installing CLI dependencies..."
	cd interfaces/cli && go mod download && go mod tidy
	@echo "Installing Web Backend dependencies..."
	cd interfaces/web/backend && go mod download && go mod tidy
	@echo "Installing Web Frontend dependencies..."
	cd interfaces/web/frontend && npm install
	@echo "Installing Flutter dependencies..."
	cd interfaces/mobile/flutter && flutter pub get
	@echo "Installing Desktop dependencies..."
	cd interfaces/desktop/electron && npm install

build: ## Build all applications
	@echo "Building Core..."
	cd core && go build ./...
	@echo "Building CLI..."
	cd interfaces/cli && go build -o bin/syntropy-cli cmd/main.go
	@echo "Building Web Backend..."
	cd interfaces/web/backend && go build -o bin/api-server cmd/main.go
	@echo "Building Web Frontend..."
	cd interfaces/web/frontend && npm run build
	@echo "Building Desktop App..."
	cd interfaces/desktop/electron && npm run build

test: ## Run all tests
	@echo "Running Core tests..."
	cd core && go test ./...
	@echo "Running CLI tests..."
	cd interfaces/cli && go test ./...
	@echo "Running Web Backend tests..."
	cd interfaces/web/backend && go test ./...
	@echo "Running Web Frontend tests..."
	cd interfaces/web/frontend && npm test
	@echo "Running Flutter tests..."
	cd interfaces/mobile/flutter && flutter test

lint: ## Run linters
	@echo "Running Go linters..."
	cd core && golangci-lint run
	cd interfaces/cli && golangci-lint run
	cd interfaces/web/backend && golangci-lint run
	@echo "Running frontend linters..."
	cd interfaces/web/frontend && npm run lint
	@echo "Running Flutter linters..."
	cd interfaces/mobile/flutter && flutter analyze

# Docker
docker: ## Build Docker images
	@echo "Building Docker images..."
	docker build -t syntropy/core -f deployments/docker/core.Dockerfile .
	docker build -t syntropy/cli -f deployments/docker/cli.Dockerfile .
	docker build -t syntropy/web-backend -f deployments/docker/web-backend.Dockerfile .
	docker build -t syntropy/web-frontend -f deployments/docker/web-frontend.Dockerfile .

up: ## Start development environment
	@echo "Starting development environment..."
	docker-compose -f deployments/docker/docker-compose.dev.yml up -d

down: ## Stop development environment
	@echo "Stopping development environment..."
	docker-compose -f deployments/docker/docker-compose.dev.yml down

# Database
db-migrate: ## Run database migrations
	@echo "Running database migrations..."
	cd core && go run cmd/migrate/main.go up

db-seed: ## Seed database with initial data
	@echo "Seeding database..."
	cd core && go run cmd/seed/main.go

# Deployment
deploy-staging: ## Deploy to staging environment
	@echo "Deploying to staging..."
	kubectl apply -f deployments/kubernetes/staging/
	helm upgrade --install syntropy-staging deployments/helm/syntropy --namespace staging

deploy-production: ## Deploy to production environment
	@echo "Deploying to production..."
	kubectl apply -f deployments/kubernetes/production/
	helm upgrade --install syntropy-prod deployments/helm/syntropy --namespace production

# Utilities
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf interfaces/cli/bin/
	rm -rf interfaces/web/backend/bin/
	rm -rf interfaces/web/frontend/dist/
	rm -rf interfaces/desktop/electron/dist/
	rm -rf interfaces/mobile/flutter/build/

security-scan: ## Run security scans
	@echo "Running security scans..."
	cd core && gosec ./...
	cd interfaces/cli && gosec ./...
	cd interfaces/web/backend && gosec ./...
	npm audit --prefix interfaces/web/frontend

# Development helpers
dev-cli: ## Run CLI in development mode
	@echo "Running CLI in development mode..."
	cd interfaces/cli && go run cmd/main.go

dev-web-backend: ## Run Web Backend in development mode
	@echo "Running Web Backend in development mode..."
	cd interfaces/web/backend && go run cmd/main.go

dev-web-frontend: ## Run Web Frontend in development mode
	@echo "Running Web Frontend in development mode..."
	cd interfaces/web/frontend && npm run dev

dev-mobile: ## Run mobile app in development mode
	@echo "Running mobile app in development mode..."
	cd interfaces/mobile/flutter && flutter run

dev-desktop: ## Run desktop app in development mode
	@echo "Running desktop app in development mode..."
	cd interfaces/desktop/electron && npm run dev

# Documentation
docs: ## Generate documentation
	@echo "Generating API documentation..."
	swagger generate spec -o interfaces/api/openapi/swagger.json
	@echo "Generating Go documentation..."
	cd core && godoc -http=:6060

# Release
release: ## Create a new release
	@echo "Creating release..."
	@read -p "Enter version (e.g., v1.0.0): " version; \
	git tag -a $$version -m "Release $$version"; \
	git push origin $$version; \
	gh release create $$version --title "Release $$version" --notes "Release $$version"