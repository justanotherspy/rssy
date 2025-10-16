# Task 1: Project Setup & Structure

## Overview
Set up the initial project structure, initialize both Go backend and SvelteKit frontend projects, and create the foundational build automation tools.

## Goals
- Create organized directory structure
- Initialize Go module and project scaffolding
- Initialize SvelteKit project with TypeScript
- Set up Makefile for common development tasks
- Verify all tools and dependencies are installed
- Create initial git repository structure

## Prerequisites
Ensure the following are installed:
- Go 1.25 or higher
- Node.js 18+ and npm
- SQLite3
- Git
- Make

## Detailed Steps

### 1. Directory Structure Creation

Create the following directory structure:

```
stream/
├── backend/
│   ├── cmd/
│   │   └── api/
│   ├── internal/
│   │   ├── models/
│   │   ├── handlers/
│   │   ├── services/
│   │   └── database/
│   └── pkg/
├── frontend/
├── docs/
└── scripts/
```

**Commands:**
```bash
cd /home/daniel/claude/stream
mkdir -p backend/cmd/api
mkdir -p backend/internal/{models,handlers,services,database}
mkdir -p backend/pkg
mkdir -p frontend
mkdir -p scripts
```

### 2. Initialize Go Backend

**Step 2.1: Initialize Go module**
```bash
cd backend
go mod init github.com/justanotherspy/rssy
```

**Step 2.2: Add initial dependencies**
```bash
# HTTP router (using gorilla/mux or chi)
go get github.com/go-chi/chi/v5

# SQLite driver
go get github.com/mattn/go-sqlite3

# RSS/Atom parser
go get github.com/mmcdole/gofeed

# Environment variable handling
go get github.com/joho/godotenv

# CORS middleware
go get github.com/go-chi/cors
```

**Step 2.3: Create initial main.go**

Create [backend/cmd/api/main.go](backend/cmd/api/main.go) with basic structure:
```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    fmt.Println("RSSY API Server")
    log.Println("Starting server on :8080")

    // Basic health check endpoint
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
```

**Step 2.4: Verify Go setup**
```bash
cd backend
go mod tidy
go build ./cmd/api
./api  # Should start server on port 8080
```

### 3. Initialize SvelteKit Frontend

**Step 3.1: Create SvelteKit project**
```bash
cd /home/daniel/claude/stream/frontend
npm create svelte@latest .
```

When prompted, select:
- Skeleton project
- Yes, using TypeScript syntax
- Add ESLint for code linting
- Add Prettier for code formatting
- Add Playwright for browser testing (optional)
- Add Vitest for unit testing (optional)

**Step 3.2: Install dependencies**
```bash
npm install
```

**Step 3.3: Add additional frontend dependencies**
```bash
# For API calls
npm install axios

# For date formatting
npm install date-fns

# For icons (optional but recommended)
npm install lucide-svelte
```

**Step 3.4: Configure TypeScript**

Verify [frontend/tsconfig.json](frontend/tsconfig.json) is properly configured for SvelteKit.

**Step 3.5: Test frontend**
```bash
npm run dev
```

Should start dev server on http://localhost:5173

### 4. Create Makefile

Create [Makefile](Makefile) in the project root:

```makefile
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
```

### 5. Create Initial Configuration Files

**Step 5.1: Backend .env template**

Create [backend/.env.example](backend/.env.example):
```env
# Server Configuration
PORT=8080
HOST=localhost

# Database
DATABASE_PATH=./rssy.db

# RSS Polling
FEED_REFRESH_INTERVAL=10m

# CORS
ALLOWED_ORIGINS=http://localhost:5173
```

**Step 5.2: Create .gitignore**

Create [.gitignore](.gitignore) in project root:
```gitignore
# Go
backend/bin/
backend/*.db
backend/.env

# Node
frontend/node_modules/
frontend/build/
frontend/.svelte-kit/
frontend/.env

# OS
.DS_Store
*.swp
*.swo
*~

# IDE
.vscode/
.idea/
*.iml
```

### 6. Set Up Git Repository

**Step 6.1: Initialize repository**
```bash
cd /home/daniel/claude/stream
git init
git add .
git commit -m "Initial project setup

- Go backend scaffolding with chi router
- SvelteKit frontend with TypeScript
- Makefile for build automation
- Basic project structure and configuration"
```

**Step 6.2: Create initial branch**
```bash
git checkout -b main
```

### 7. Verification Checklist

Run through this checklist to ensure everything is set up correctly:

- [ ] Directory structure created
- [ ] Go module initialized and dependencies installed
- [ ] Go backend compiles without errors
- [ ] Go backend runs and responds to health check endpoint
- [ ] SvelteKit project created with TypeScript
- [ ] Frontend dependencies installed
- [ ] Frontend dev server starts successfully
- [ ] Makefile targets work (`make help`, `make backend-build`, etc.)
- [ ] Git repository initialized
- [ ] .gitignore properly excludes build artifacts and dependencies

### 8. Testing the Setup

**Test Backend:**
```bash
# Terminal 1
cd /home/daniel/claude/stream
make backend-run

# Terminal 2
curl http://localhost:8080/health
# Should return: OK
```

**Test Frontend:**
```bash
# Terminal 1
cd /home/daniel/claude/stream
make frontend-run

# Browser
# Visit http://localhost:5173
# Should see SvelteKit welcome page
```

## Success Criteria

- Both backend and frontend projects are initialized
- All dependencies are installed and working
- Backend server starts and responds to requests
- Frontend dev server starts and displays properly
- Makefile commands execute successfully
- Git repository is initialized with proper structure
- All verification checklist items are complete

## Next Steps

After completing Task 1, proceed to Task 2 (Data Model Design) where we will:
- Design the database schema
- Create SQLite tables for feeds and posts
- Implement Go models and database layer
- Set up migrations or initialization scripts

## Notes

- Keep both servers running during development
- Backend API will run on port 8080
- Frontend dev server will run on port 5173
- Frontend will proxy API requests to backend during development
- Use `make help` to see all available commands

## Agent & Hook Recommendations for This Task

### Recommended Agent Usage
- **Direct tool usage** for file creation and directory setup
- **Bash commands** for package installations and verification
- Keep operations simple and sequential

### Hook Configuration
Since this is the setup phase, hooks are not critical yet. However, consider adding:
- Post-file-save hook for auto-formatting Go and TypeScript files
- Pre-commit hook will be valuable once we have code to test (Task 2+)

## Troubleshooting

**Go module issues:**
- Ensure GOPATH is set correctly
- Run `go clean -modcache` if you encounter module errors

**Frontend build issues:**
- Clear npm cache: `npm cache clean --force`
- Delete node_modules and reinstall: `rm -rf node_modules && npm install`

**Port conflicts:**
- If port 8080 is in use, modify PORT in backend/.env
- If port 5173 is in use, modify vite.config.ts in frontend

**SQLite issues:**
- Ensure SQLite3 development libraries are installed
- On Ubuntu/Debian: `sudo apt-get install libsqlite3-dev`
- On macOS: SQLite should be pre-installed
