# Task 3: Go Backend API Development

## Overview
Build the complete REST API backend with HTTP handlers, implement RSS feed fetching service, create background polling mechanism, and add middleware for CORS and error handling.

## Goals
- Create RESTful API endpoints for feeds and posts
- Implement RSS/Atom feed parser integration
- Build background service for periodic feed polling
- Add proper error handling and logging
- Implement CORS middleware for frontend integration
- Create configuration management
- Add graceful shutdown handling

## Prerequisites
- Task 1 completed (project structure set up)
- Task 2 completed (database layer implemented)
- gofeed library installed
- chi router installed

## API Design

### Endpoints Overview

#### Feed Endpoints
- `GET /api/feeds` - List all feeds
- `GET /api/feeds/:id` - Get specific feed
- `POST /api/feeds` - Create new feed
- `PUT /api/feeds/:id` - Update feed
- `DELETE /api/feeds/:id` - Delete feed
- `POST /api/feeds/reddit` - Quick add Reddit subreddit

#### Post Endpoints
- `GET /api/posts` - List all posts (paginated)
- `GET /api/posts/feed/:feedId` - Get posts for specific feed
- `PATCH /api/posts/:id/read` - Mark post as read/unread
- `DELETE /api/posts` - Delete all posts (reset)

#### Settings Endpoints
- `GET /api/settings` - Get current settings
- `PUT /api/settings` - Update settings (e.g., refresh interval)

#### System Endpoints
- `GET /health` - Health check
- `POST /api/feeds/refresh` - Manually trigger feed refresh
- `POST /api/feeds/:id/refresh` - Refresh specific feed

## Detailed Steps

### 1. Create Configuration Package

**Step 1.1: Create config.go**

Create [backend/internal/config/config.go](backend/internal/config/config.go):

```go
package config

import (
    "log"
    "os"
    "strconv"
    "time"

    "github.com/joho/godotenv"
)

type Config struct {
    Port                 string
    Host                 string
    DatabasePath         string
    FeedRefreshInterval  time.Duration
    AllowedOrigins       []string
}

func Load() *Config {
    // Load .env file if it exists
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    port := getEnv("PORT", "8080")
    host := getEnv("HOST", "localhost")
    dbPath := getEnv("DATABASE_PATH", "./rssy.db")

    refreshInterval := getEnvAsDuration("FEED_REFRESH_INTERVAL", "10m")
    allowedOrigins := getEnvAsSlice("ALLOWED_ORIGINS", []string{"http://localhost:5173"})

    return &Config{
        Port:                port,
        Host:                host,
        DatabasePath:        dbPath,
        FeedRefreshInterval: refreshInterval,
        AllowedOrigins:      allowedOrigins,
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsDuration(key, defaultValue string) time.Duration {
    valueStr := getEnv(key, defaultValue)
    duration, err := time.ParseDuration(valueStr)
    if err != nil {
        log.Printf("Invalid duration for %s, using default: %s", key, defaultValue)
        duration, _ = time.ParseDuration(defaultValue)
    }
    return duration
}

func getEnvAsSlice(key string, defaultValue []string) []string {
    if value := os.Getenv(key); value != "" {
        return []string{value}
    }
    return defaultValue
}
```

### 2. Create HTTP Handlers

**Step 2.1: Create handlers package structure**

Create [backend/internal/handlers/handlers.go](backend/internal/handlers/handlers.go):

```go
package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/justanotherspy/rssy/internal/database"
)

type Handler struct {
    db *database.DB
}

func New(db *database.DB) *Handler {
    return &Handler{db: db}
}

// Response helpers
type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(Response{
        Success: status < 400,
        Data:    data,
    })
}

func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(Response{
        Success: false,
        Error:   message,
    })
}
```

**Step 2.2: Create feed handlers**

Create [backend/internal/handlers/feed_handlers.go](backend/internal/handlers/feed_handlers.go):

```go
package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "github.com/justanotherspy/rssy/internal/models"
)

// GetAllFeeds handles GET /api/feeds
func (h *Handler) GetAllFeeds(w http.ResponseWriter, r *http.Request) {
    feeds, err := h.db.GetAllFeeds()
    if err != nil {
        h.respondError(w, http.StatusInternalServerError, "Failed to retrieve feeds")
        return
    }

    h.respondJSON(w, http.StatusOK, feeds)
}

// GetFeedByID handles GET /api/feeds/:id
func (h *Handler) GetFeedByID(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        h.respondError(w, http.StatusBadRequest, "Invalid feed ID")
        return
    }

    feed, err := h.db.GetFeedByID(id)
    if err != nil {
        h.respondError(w, http.StatusNotFound, "Feed not found")
        return
    }

    h.respondJSON(w, http.StatusOK, feed)
}

// CreateFeed handles POST /api/feeds
func (h *Handler) CreateFeed(w http.ResponseWriter, r *http.Request) {
    var req models.CreateFeedRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    // Validate required fields
    if req.Name == "" || req.URL == "" {
        h.respondError(w, http.StatusBadRequest, "Name and URL are required")
        return
    }

    feed, err := h.db.CreateFeed(req)
    if err != nil {
        h.respondError(w, http.StatusInternalServerError, "Failed to create feed")
        return
    }

    h.respondJSON(w, http.StatusCreated, feed)
}

// UpdateFeed handles PUT /api/feeds/:id
func (h *Handler) UpdateFeed(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        h.respondError(w, http.StatusBadRequest, "Invalid feed ID")
        return
    }

    var req models.UpdateFeedRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    feed, err := h.db.UpdateFeed(id, req)
    if err != nil {
        h.respondError(w, http.StatusInternalServerError, "Failed to update feed")
        return
    }

    h.respondJSON(w, http.StatusOK, feed)
}

// DeleteFeed handles DELETE /api/feeds/:id
func (h *Handler) DeleteFeed(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        h.respondError(w, http.StatusBadRequest, "Invalid feed ID")
        return
    }

    if err := h.db.DeleteFeed(id); err != nil {
        h.respondError(w, http.StatusInternalServerError, "Failed to delete feed")
        return
    }

    h.respondJSON(w, http.StatusOK, map[string]string{"message": "Feed deleted successfully"})
}

// CreateRedditFeed handles POST /api/feeds/reddit
func (h *Handler) CreateRedditFeed(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Subreddit string `json:"subreddit"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    if req.Subreddit == "" {
        h.respondError(w, http.StatusBadRequest, "Subreddit name is required")
        return
    }

    // Create feed from subreddit
    feedReq := models.CreateFeedRequest{
        Name:     "r/" + req.Subreddit,
        URL:      "https://www.reddit.com/r/" + req.Subreddit + "/.rss",
        Category: "Reddit",
        SiteURL:  "https://www.reddit.com/r/" + req.Subreddit,
        Description: "Reddit /r/" + req.Subreddit + " feed",
    }

    feed, err := h.db.CreateFeed(feedReq)
    if err != nil {
        h.respondError(w, http.StatusInternalServerError, "Failed to create Reddit feed")
        return
    }

    h.respondJSON(w, http.StatusCreated, feed)
}
```

**Step 2.3: Create post handlers**

Create [backend/internal/handlers/post_handlers.go](backend/internal/handlers/post_handlers.go):

```go
package handlers

import (
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
)

// GetAllPosts handles GET /api/posts
func (h *Handler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
    // Parse pagination parameters
    limitStr := r.URL.Query().Get("limit")
    offsetStr := r.URL.Query().Get("offset")

    limit := 50 // default
    if limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil {
            limit = l
        }
    }

    offset := 0 // default
    if offsetStr != "" {
        if o, err := strconv.Atoi(offsetStr); err == nil {
            offset = o
        }
    }

    posts, err := h.db.GetAllPosts(limit, offset)
    if err != nil {
        h.respondError(w, http.StatusInternalServerError, "Failed to retrieve posts")
        return
    }

    h.respondJSON(w, http.StatusOK, posts)
}

// GetPostsByFeed handles GET /api/posts/feed/:feedId
func (h *Handler) GetPostsByFeed(w http.ResponseWriter, r *http.Request) {
    feedIDStr := chi.URLParam(r, "feedId")
    feedID, err := strconv.ParseInt(feedIDStr, 10, 64)
    if err != nil {
        h.respondError(w, http.StatusBadRequest, "Invalid feed ID")
        return
    }

    // Parse pagination
    limitStr := r.URL.Query().Get("limit")
    offsetStr := r.URL.Query().Get("offset")

    limit := 50
    if limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil {
            limit = l
        }
    }

    offset := 0
    if offsetStr != "" {
        if o, err := strconv.Atoi(offsetStr); err == nil {
            offset = o
        }
    }

    posts, err := h.db.GetPostsByFeedID(feedID, limit, offset)
    if err != nil {
        h.respondError(w, http.StatusInternalServerError, "Failed to retrieve posts")
        return
    }

    h.respondJSON(w, http.StatusOK, posts)
}

// MarkPostRead handles PATCH /api/posts/:id/read
func (h *Handler) MarkPostRead(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        h.respondError(w, http.StatusBadRequest, "Invalid post ID")
        return
    }

    var req struct {
        IsRead bool `json:"is_read"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    if err := h.db.MarkPostAsRead(id, req.IsRead); err != nil {
        h.respondError(w, http.StatusInternalServerError, "Failed to update post")
        return
    }

    h.respondJSON(w, http.StatusOK, map[string]string{"message": "Post updated successfully"})
}

// DeleteAllPosts handles DELETE /api/posts
func (h *Handler) DeleteAllPosts(w http.ResponseWriter, r *http.Request) {
    if err := h.db.DeleteAllPosts(); err != nil {
        h.respondError(w, http.StatusInternalServerError, "Failed to delete posts")
        return
    }

    h.respondJSON(w, http.StatusOK, map[string]string{"message": "All posts deleted successfully"})
}
```

### 3. Create Feed Fetching Service

**Step 3.1: Create fetcher service**

Create [backend/internal/services/fetcher.go](backend/internal/services/fetcher.go):

```go
package services

import (
    "log"
    "time"

    "github.com/mmcdole/gofeed"
    "github.com/justanotherspy/rssy/internal/database"
    "github.com/justanotherspy/rssy/internal/models"
)

type FeedFetcher struct {
    db     *database.DB
    parser *gofeed.Parser
}

func NewFeedFetcher(db *database.DB) *FeedFetcher {
    return &FeedFetcher{
        db:     db,
        parser: gofeed.NewParser(),
    }
}

// FetchFeed fetches and parses a single feed
func (f *FeedFetcher) FetchFeed(feed *models.Feed) error {
    log.Printf("Fetching feed: %s (%s)", feed.Name, feed.URL)

    parsedFeed, err := f.parser.ParseURL(feed.URL)
    if err != nil {
        log.Printf("Error parsing feed %s: %v", feed.Name, err)
        return err
    }

    // Process each item in the feed
    newPostCount := 0
    for _, item := range parsedFeed.Items {
        // Check if post already exists
        existing, err := f.db.GetPostByGUID(feed.ID, item.GUID)
        if err != nil {
            log.Printf("Error checking post existence: %v", err)
            continue
        }

        if existing != nil {
            continue // Post already exists
        }

        // Create new post
        post := &models.Post{
            FeedID:      feed.ID,
            Title:       item.Title,
            Link:        item.Link,
            Description: item.Description,
            Content:     item.Content,
            Author:      getAuthor(item),
            PublishedAt: getPublishedTime(item),
            ImageURL:    getImageURL(item),
            GUID:        item.GUID,
        }

        if err := f.db.CreatePost(post); err != nil {
            log.Printf("Error creating post: %v", err)
            continue
        }

        newPostCount++
    }

    // Update feed last fetched time
    if err := f.db.UpdateFeedLastFetched(feed.ID, time.Now()); err != nil {
        log.Printf("Error updating feed last fetched time: %v", err)
    }

    log.Printf("Fetched %d new posts from %s", newPostCount, feed.Name)
    return nil
}

// FetchAllFeeds fetches all active feeds
func (f *FeedFetcher) FetchAllFeeds() error {
    feeds, err := f.db.GetAllFeeds()
    if err != nil {
        return err
    }

    for _, feed := range feeds {
        if !feed.IsActive {
            continue
        }

        if err := f.FetchFeed(&feed); err != nil {
            log.Printf("Failed to fetch feed %s: %v", feed.Name, err)
            // Continue with other feeds even if one fails
        }
    }

    return nil
}

// Helper functions
func getAuthor(item *gofeed.Item) string {
    if item.Author != nil {
        return item.Author.Name
    }
    return ""
}

func getPublishedTime(item *gofeed.Item) *time.Time {
    if item.PublishedParsed != nil {
        return item.PublishedParsed
    }
    if item.UpdatedParsed != nil {
        return item.UpdatedParsed
    }
    return nil
}

func getImageURL(item *gofeed.Item) string {
    if item.Image != nil {
        return item.Image.URL
    }

    // Try to find image in enclosures
    for _, enclosure := range item.Enclosures {
        if enclosure.Type != "" && len(enclosure.Type) > 5 && enclosure.Type[:5] == "image" {
            return enclosure.URL
        }
    }

    return ""
}
```

**Step 3.2: Create polling service**

Create [backend/internal/services/poller.go](backend/internal/services/poller.go):

```go
package services

import (
    "context"
    "log"
    "time"

    "github.com/justanotherspy/rssy/internal/database"
)

type Poller struct {
    fetcher  *FeedFetcher
    interval time.Duration
    ctx      context.Context
    cancel   context.CancelFunc
}

func NewPoller(db *database.DB, interval time.Duration) *Poller {
    ctx, cancel := context.WithCancel(context.Background())
    return &Poller{
        fetcher:  NewFeedFetcher(db),
        interval: interval,
        ctx:      ctx,
        cancel:   cancel,
    }
}

// Start begins the polling loop
func (p *Poller) Start() {
    log.Printf("Starting feed poller with interval: %v", p.interval)

    // Fetch immediately on start
    go func() {
        if err := p.fetcher.FetchAllFeeds(); err != nil {
            log.Printf("Error during initial fetch: %v", err)
        }
    }()

    // Start periodic polling
    ticker := time.NewTicker(p.interval)
    go func() {
        for {
            select {
            case <-ticker.C:
                log.Println("Polling feeds...")
                if err := p.fetcher.FetchAllFeeds(); err != nil {
                    log.Printf("Error polling feeds: %v", err)
                }
            case <-p.ctx.Done():
                ticker.Stop()
                log.Println("Feed poller stopped")
                return
            }
        }
    }()
}

// Stop stops the polling loop
func (p *Poller) Stop() {
    log.Println("Stopping feed poller...")
    p.cancel()
}
```

### 4. Create Router and Middleware

**Step 4.1: Create router**

Create [backend/internal/router/router.go](backend/internal/router/router.go):

```go
package router

import (
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/cors"
    "github.com/justanotherspy/rssy/internal/handlers"
)

func New(h *handlers.Handler, allowedOrigins []string) *chi.Mux {
    r := chi.NewRouter()

    // Middleware
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)

    // CORS
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   allowedOrigins,
        AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: false,
        MaxAge:           300,
    }))

    // Health check
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    // API routes
    r.Route("/api", func(r chi.Router) {
        // Feed routes
        r.Route("/feeds", func(r chi.Router) {
            r.Get("/", h.GetAllFeeds)
            r.Post("/", h.CreateFeed)
            r.Post("/reddit", h.CreateRedditFeed)

            r.Route("/{id}", func(r chi.Router) {
                r.Get("/", h.GetFeedByID)
                r.Put("/", h.UpdateFeed)
                r.Delete("/", h.DeleteFeed)
            })
        })

        // Post routes
        r.Route("/posts", func(r chi.Router) {
            r.Get("/", h.GetAllPosts)
            r.Delete("/", h.DeleteAllPosts)

            r.Get("/feed/{feedId}", h.GetPostsByFeed)

            r.Route("/{id}", func(r chi.Router) {
                r.Patch("/read", h.MarkPostRead)
            })
        })
    })

    return r
}
```

### 5. Update main.go

Update [backend/cmd/api/main.go](backend/cmd/api/main.go) with complete server:

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/justanotherspy/rssy/internal/config"
    "github.com/justanotherspy/rssy/internal/database"
    "github.com/justanotherspy/rssy/internal/handlers"
    "github.com/justanotherspy/rssy/internal/router"
    "github.com/justanotherspy/rssy/internal/services"
)

func main() {
    log.Println("Starting RSSY API Server...")

    // Load configuration
    cfg := config.Load()
    log.Printf("Configuration loaded: Port=%s, RefreshInterval=%v", cfg.Port, cfg.FeedRefreshInterval)

    // Initialize database
    db, err := database.New(cfg.DatabasePath)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Initialize schema
    if err := db.InitSchema(); err != nil {
        log.Fatalf("Failed to initialize schema: %v", err)
    }

    // Seed default feeds
    if err := db.SeedDefaultFeeds(); err != nil {
        log.Fatalf("Failed to seed default feeds: %v", err)
    }

    // Create handlers
    h := handlers.New(db)

    // Create router
    r := router.New(h, cfg.AllowedOrigins)

    // Start feed poller
    poller := services.NewPoller(db, cfg.FeedRefreshInterval)
    poller.Start()
    defer poller.Stop()

    // Create HTTP server
    srv := &http.Server{
        Addr:         ":" + cfg.Port,
        Handler:      r,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Start server in goroutine
    go func() {
        log.Printf("Server listening on :%s", cfg.Port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()

    // Wait for interrupt signal for graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    log.Println("Server stopped")
}
```

### 6. Testing the API

**Step 6.1: Build and run**
```bash
cd backend
go mod tidy
go run ./cmd/api
```

**Step 6.2: Test endpoints**

```bash
# Health check
curl http://localhost:8080/health

# Get all feeds
curl http://localhost:8080/api/feeds

# Get posts
curl http://localhost:8080/api/posts

# Create new feed
curl -X POST http://localhost:8080/api/feeds \
  -H "Content-Type: application/json" \
  -d '{"name":"My Feed","url":"https://example.com/feed.xml"}'

# Create Reddit feed
curl -X POST http://localhost:8080/api/feeds/reddit \
  -H "Content-Type: application/json" \
  -d '{"subreddit":"golang"}'

# Get posts for specific feed
curl http://localhost:8080/api/posts/feed/1

# Mark post as read
curl -X PATCH http://localhost:8080/api/posts/1/read \
  -H "Content-Type: application/json" \
  -d '{"is_read":true}'

# Delete all posts
curl -X DELETE http://localhost:8080/api/posts
```

### 7. Add Manual Refresh Handlers

Add to [backend/internal/handlers/feed_handlers.go](backend/internal/handlers/feed_handlers.go):

```go
// RefreshAllFeeds manually triggers feed refresh
func (h *Handler) RefreshAllFeeds(w http.ResponseWriter, r *http.Request) {
    fetcher := services.NewFeedFetcher(h.db)

    if err := fetcher.FetchAllFeeds(); err != nil {
        h.respondError(w, http.StatusInternalServerError, "Failed to refresh feeds")
        return
    }

    h.respondJSON(w, http.StatusOK, map[string]string{"message": "Feeds refreshed successfully"})
}

// RefreshFeed manually triggers refresh for specific feed
func (h *Handler) RefreshFeed(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        h.respondError(w, http.StatusBadRequest, "Invalid feed ID")
        return
    }

    feed, err := h.db.GetFeedByID(id)
    if err != nil {
        h.respondError(w, http.StatusNotFound, "Feed not found")
        return
    }

    fetcher := services.NewFeedFetcher(h.db)
    if err := fetcher.FetchFeed(feed); err != nil {
        h.respondError(w, http.StatusInternalServerError, "Failed to refresh feed")
        return
    }

    h.respondJSON(w, http.StatusOK, map[string]string{"message": "Feed refreshed successfully"})
}
```

Update router to include these endpoints.

## Success Criteria

- All API endpoints respond correctly
- Feed fetching works for various RSS/Atom feeds
- Background polling runs at configured interval
- CORS properly configured for frontend
- Graceful shutdown works properly
- Error handling is consistent
- Logging provides useful information
- Manual refresh endpoints work

## Next Steps

After completing Task 3, proceed to Task 4 (Svelte Frontend) where we will:
- Create SvelteKit page structure
- Build feed list sidebar component
- Create post card components
- Implement feed management modals
- Add settings interface
- Connect to backend API
- Style the application

## Notes

- The poller runs in a separate goroutine
- Graceful shutdown ensures the poller stops cleanly
- CORS must allow the frontend origin
- Reddit feeds use the .rss endpoint
- Some feeds may have different structures (handle gracefully)
- Consider rate limiting for production use

## Agent & Hook Recommendations for This Task

### Recommended Agent Usage
- **Direct tool usage** for creating handler files
- **General-purpose agent** if refactoring needed
- Test each endpoint manually as you build

### Pre-commit Hook
Now is a good time to set up the pre-commit hook mentioned in CLAUDE.md:
```bash
cd backend && go test ./... && go vet ./...
```

### Testing Strategy
- Use curl for manual API testing
- Verify polling works by checking logs
- Test CORS by running frontend (Task 4)
- Monitor database for new posts

## Troubleshooting

**Feed parsing fails:**
- Some feeds may have invalid XML
- gofeed handles most formats, but not all
- Log errors and continue with other feeds

**CORS issues:**
- Ensure AllowedOrigins includes frontend URL
- Check browser console for CORS errors
- Verify OPTIONS requests are handled

**Polling not working:**
- Check logs for errors
- Verify feeds are marked as active
- Ensure database is accessible

**Port already in use:**
- Change PORT in .env file
- Check for other processes on port 8080
