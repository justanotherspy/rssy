.PHONY: help backend-run frontend-run backend-build frontend-build backend-test frontend-test clean install-deps

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install-deps: ## Install all dependencies
	@echo "Installing Go dependencies..."
	cd backend && go mod download
	@echo "Installing Node dependencies..."
	cd frontend && npm install

backend-build: ## Build Go backend
	@echo "Building backend..."
	cd backend && go build -o bin/api ./cmd/api

backend-run: ## Run Go backend in development mode
	@echo "Starting backend server..."
	cd backend && go run ./cmd/api

backend-test: ## Run backend tests
	cd backend && go test -v ./...

backend-vet: ## Run Go vet
	cd backend && go vet ./...

frontend-build: ## Build frontend for production
	@echo "Building frontend..."
	cd frontend && npm run build

frontend-run: ## Run frontend in development mode
	@echo "Starting frontend dev server..."
	cd frontend && npm run dev

frontend-test: ## Run frontend tests
	cd frontend && npm run test

frontend-check: ## Run frontend type checking
	cd frontend && npm run check

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf backend/bin
	rm -rf frontend/build
	rm -rf frontend/.svelte-kit

dev: ## Run both backend and frontend in development mode (requires two terminals)
	@echo "Run 'make backend-run' in one terminal and 'make frontend-run' in another"

format-go: ## Format Go code
	cd backend && go fmt ./...

format-frontend: ## Format frontend code
	cd frontend && npm run format

lint: backend-vet frontend-check ## Run all linters
