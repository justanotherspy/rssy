# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project: RSSY RSS Reader

RSSY (pronounced "rizzy") is a full-stack RSS reader with a Go backend, Svelte frontend, and SQLite database.

**Status:** MVP Complete - Both backend and frontend are production-ready.

## Technology Stack

- **Backend:** Go 1.21+ (project uses 1.25.0) with chi router, SQLite3
- **Frontend:** Svelte 5, SvelteKit v2, TypeScript, Axios
- **Database:** SQLite with automatic schema initialization
- **Build:** Makefile for automation

## Common Commands

### Development

**Note:** These commands start the application in development mode.

```bash
# Quick start (requires 2 terminals)
make backend-run      # Terminal 1: Start API on :8080
make frontend-run     # Terminal 2: Start UI on :5173

# Install dependencies
make install-deps     # Install both Go and Node dependencies

# Testing and linting
make backend-test     # Run Go tests
make backend-vet      # Run Go vet
make frontend-check   # TypeScript type checking
make lint             # Run all linters

# Building
make backend-build    # Builds to ./backend/bin/api
make frontend-build   # Builds to ./frontend/.svelte-kit/output

# Cleanup
make clean            # Remove build artifacts

# Formatting
make format-go        # Format Go code with gofmt
cd frontend && npm run format  # Format frontend with prettier
```

### Backend-only commands
```bash
cd backend
go run ./cmd/api              # Run server
go test ./...                 # Run tests
go test ./internal/database   # Test specific package
go build -o bin/api ./cmd/api # Build binary
```

### Frontend-only commands
```bash
cd frontend
npm run dev       # Dev server with HMR
npm run build     # Production build
npm run preview   # Preview production build
npm run check     # Type checking (runs svelte-check)
```

## Architecture Overview

### Backend Structure

The Go backend follows a clean layered architecture:

1. **Entry Point** (`cmd/api/main.go`):
   - Loads configuration from env vars (or defaults)
   - Initializes SQLite database and schema
   - Seeds default feeds (HackerNews, TechCrunch, etc.)
   - Starts background feed poller (runs every 10 minutes by default)
   - Creates HTTP server with graceful shutdown

2. **Configuration** (`internal/config`):
   - Environment-based config with sensible defaults
   - Supports `.env` file or environment variables
   - Key settings: PORT, DATABASE_PATH, FEED_REFRESH_INTERVAL, ALLOWED_ORIGINS

3. **Database Layer** (`internal/database`):
   - `database.go`: Main DB struct and connection management
   - `schema.go`: Table creation and migrations
   - `seed.go`: Default feed data
   - `feed_repository.go`: Feed CRUD operations
   - `post_repository.go`: Post CRUD operations
   - Pattern: Repository pattern separating data access from business logic

4. **Models** (`internal/models`):
   - `feed.go`: Feed struct with validation tags
   - `post.go`: Post struct with RSS/Atom field mapping
   - These map directly to database tables

5. **Handlers** (`internal/handlers`):
   - HTTP request handlers using chi router
   - `handlers.go`: Handler struct holding DB dependency
   - `feed_handlers.go`: Feed CRUD endpoints
   - `post_handlers.go`: Post CRUD endpoints
   - Pattern: Each handler validates input → calls repository → returns JSON

6. **Router** (`internal/router/router.go`):
   - Uses chi router with middleware (logger, recoverer, CORS)
   - RESTful route structure under `/api`
   - CORS configured for local frontend development

7. **Services** (`internal/services`):
   - `poller.go`: Background goroutine that polls feeds at intervals
   - `fetcher.go`: RSS/Atom feed fetching and parsing logic
   - Uses `gofeed` library for RSS parsing
   - Automatically deduplicates posts by GUID

**Key Backend Flow:**
- Application starts → DB initialized → Default feeds seeded → Poller starts
- Poller runs in background, fetching all active feeds every N minutes
- HTTP handlers serve API requests from frontend
- Posts stored in SQLite, deduplicated by GUID to prevent duplicates

### Frontend Structure

The Svelte frontend uses a component-based architecture:

1. **Entry Point** (`src/routes/+page.svelte`):
   - Main application page with two-panel layout
   - Loads feeds and posts on mount
   - Manages modal visibility and data synchronization

2. **API Client** (`src/lib/api.ts`):
   - Axios-based HTTP client with TypeScript interfaces
   - Exports `feedsApi` and `postsApi` with typed methods
   - Base URL from `VITE_API_URL` environment variable

3. **State Management** (`src/lib/stores.ts`):
   - Svelte writable stores for global reactive state
   - Stores: feeds, posts, selectedFeedId, loading, error, modals
   - Pattern: Components subscribe to stores, update via set()

4. **Components** (`src/lib/components/`):
   - `Sidebar.svelte`: Dark sidebar with feed list, action buttons
   - `PostCard.svelte`: Individual post display with image, metadata, content
   - `AddFeedModal.svelte`: Tabbed modal (RSS vs Reddit) for adding feeds
   - `SettingsModal.svelte`: Settings for refresh interval, delete all posts

**Key Frontend Flow:**
- App loads → Fetch feeds and posts from API → Display in two panels
- User clicks feed → Posts filtered by feed_id → PostCards rendered
- User adds feed → Modal submits to API → Feeds reloaded
- Background poller (backend) adds new posts → User refreshes or polls API

### Data Flow

```
Backend Poller (every 10min) → Fetches RSS feeds → Parses posts → Stores in SQLite
                                                                           ↓
Frontend → Calls /api/posts → Backend queries SQLite → Returns JSON → Displays posts
```

### API Endpoints

Backend serves RESTful JSON API on port 8080:

**Health Check:**
- `GET /health` - Health check endpoint (returns "OK")

**Feeds:**
- `GET /api/feeds` - List all feeds
- `POST /api/feeds` - Create feed (body: {name, url, category?})
- `POST /api/feeds/reddit` - Quick-add Reddit feed (body: {subreddit})
- `POST /api/feeds/refresh` - Manually refresh all feeds
- `GET /api/feeds/:id` - Get feed by ID
- `PUT /api/feeds/:id` - Update feed
- `DELETE /api/feeds/:id` - Delete feed
- `POST /api/feeds/:id/refresh` - Manually trigger specific feed refresh

**Posts:**
- `GET /api/posts` - List all posts (with limit/offset)
- `GET /api/posts/feed/:feedId` - List posts from specific feed
- `PATCH /api/posts/:id/read` - Mark post read/unread
- `DELETE /api/posts` - Delete all posts

All responses follow pattern: `{"data": [...], "error": null}`

## Database Schema

SQLite database (`rssy.db`) with two main tables:

**feeds table:**
- Primary key: id
- Unique constraint on url
- Tracks: name, url, category, last_fetched_at, error_count, is_active

**posts table:**
- Primary key: id
- Foreign key: feed_id → feeds(id) ON DELETE CASCADE
- Unique constraint on guid (prevents duplicate posts)
- Stores: title, link, description, content, author, published_at, image_url

## Environment Configuration

**Backend** (optional `.env` in backend/):
```
PORT=8080
DATABASE_PATH=./rssy.db
FEED_REFRESH_INTERVAL=10m
ALLOWED_ORIGINS=http://localhost:5173
```

**Frontend** (`.env` in frontend/):
```
VITE_API_URL=http://localhost:8080
```

## Design Principles

**API:**
- RESTful with clear resource naming
- JSON request/response format
- Proper HTTP status codes (200, 201, 400, 404, 500)
- Structured error responses

**Frontend:**
- Component-based with clear separation of concerns
- Minimalist flat design with dark sidebar (#1a1a1a)
- Responsive layout (works on mobile, tablet, desktop)
- Loading states and error handling throughout

**Backend:**
- Clean architecture with layered separation
- Repository pattern for data access
- Graceful shutdown with context cancellation
- Background services managed with goroutines and contexts

## Module Path

Backend uses Go module: `github.com/justanotherspy/rssy`

When adding new packages or imports, use this module path.

## Testing

Backend has test files for critical components:
- Run all tests: `go test ./...`
- Run specific package: `go test ./internal/database`
- Verbose output: `go test -v ./...`

Frontend uses svelte-check for type validation:
- Run: `npm run check` or `make frontend-check`

## Git Workflow

Follow standard feature branch workflow:
1. Create branch from main: `git checkout -b feature/xyz`
2. Make changes and test locally
3. Check status: `git status`
4. Review diffs: `git diff`
5. Stage and commit with clear messages
6. Push branch: `git push -u origin feature/xyz`
7. Create pull request

## Common Tasks

**Add a new API endpoint:**
1. Add handler method to `internal/handlers/feed_handlers.go` or `post_handlers.go`
2. Register route in `internal/router/router.go`
3. Update frontend API client in `src/lib/api.ts`
4. Use endpoint in appropriate component

**Add a new Svelte component:**
1. Create `.svelte` file in `src/lib/components/`
2. Import and use in `src/routes/+page.svelte` or other components
3. Update stores if component needs global state

**Modify database schema:**
1. Update structs in `internal/models/`
2. Update schema in `internal/database/schema.go`
3. Update repository queries in `internal/database/*_repository.go`
4. Delete `rssy.db` to test fresh schema (or write migration)

**Change feed polling interval:**
- Set `FEED_REFRESH_INTERVAL` environment variable (e.g., `5m`, `30m`, `1h`)
- Or modify default in `internal/config/config.go`

## Dependencies

**Backend:**
- `github.com/go-chi/chi/v5` - HTTP router
- `github.com/go-chi/cors` - CORS middleware
- `github.com/mattn/go-sqlite3` - SQLite driver
- `github.com/mmcdole/gofeed` - RSS/Atom parser
- `github.com/joho/godotenv` - .env file loading

**Frontend:**
- `axios` - HTTP client
- `date-fns` - Date formatting
- `lucide-svelte` - Icon components
- Svelte 5 and SvelteKit v2 for framework
