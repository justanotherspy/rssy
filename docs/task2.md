# Task 2: Data Model Design & Database Layer

## Overview
Design and implement the database schema, create SQLite database with proper tables, and build the Go data access layer with models and database operations.

## Goals
- Design normalized database schema for feeds and posts
- Create SQLite database with proper constraints and indexes
- Implement Go models (structs) for domain entities
- Build database layer with CRUD operations
- Implement database initialization and migration logic
- Add seed data with default RSS feeds

## Prerequisites
- Task 1 completed successfully
- Backend project structure in place
- SQLite3 driver installed (`github.com/mattn/go-sqlite3`)
- Understanding of RSS feed structure

## Database Schema Design

### Entity Relationship Diagram

```
┌─────────────────────┐       ┌─────────────────────┐
│      Feeds          │       │       Posts         │
├─────────────────────┤       ├─────────────────────┤
│ id (INTEGER PK)     │───┐   │ id (INTEGER PK)     │
│ name (TEXT)         │   └──<│ feed_id (INTEGER FK)│
│ url (TEXT UNIQUE)   │       │ title (TEXT)        │
│ category (TEXT)     │       │ link (TEXT UNIQUE)  │
│ site_url (TEXT)     │       │ description (TEXT)  │
│ is_active (BOOLEAN) │       │ content (TEXT)      │
│ last_fetched (TIME) │       │ author (TEXT)       │
│ created_at (TIME)   │       │ published_at (TIME) │
│ updated_at (TIME)   │       │ image_url (TEXT)    │
└─────────────────────┘       │ guid (TEXT)         │
                              │ created_at (TIME)   │
                              │ updated_at (TIME)   │
                              └─────────────────────┘
```

### Table Definitions

#### Feeds Table
Stores RSS/Atom feed sources that users subscribe to.

```sql
CREATE TABLE IF NOT EXISTS feeds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    category TEXT,
    site_url TEXT,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT 1,
    last_fetched_at DATETIME,
    error_count INTEGER DEFAULT 0,
    last_error TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_feeds_is_active ON feeds(is_active);
CREATE INDEX idx_feeds_last_fetched ON feeds(last_fetched_at);
```

#### Posts Table
Stores individual RSS feed items/entries.

```sql
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    feed_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    link TEXT NOT NULL,
    description TEXT,
    content TEXT,
    author TEXT,
    published_at DATETIME,
    image_url TEXT,
    guid TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
    UNIQUE(feed_id, guid)
);

CREATE INDEX idx_posts_feed_id ON posts(feed_id);
CREATE INDEX idx_posts_published_at ON posts(published_at DESC);
CREATE INDEX idx_posts_is_read ON posts(is_read);
CREATE UNIQUE INDEX idx_posts_feed_guid ON posts(feed_id, guid);
```

## Detailed Steps

### 1. Create Database Package Structure

**Step 1.1: Create database.go**

Create [backend/internal/database/database.go](backend/internal/database/database.go):

```go
package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "path/filepath"

    _ "github.com/mattn/go-sqlite3"
)

type DB struct {
    *sql.DB
}

// New creates a new database connection
func New(dbPath string) (*DB, error) {
    // Ensure directory exists
    dir := filepath.Dir(dbPath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create database directory: %w", err)
    }

    // Open database connection
    db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    // Test connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    log.Printf("Connected to database: %s", dbPath)

    return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
    return db.DB.Close()
}
```

**Step 1.2: Create schema.go**

Create [backend/internal/database/schema.go](backend/internal/database/schema.go):

```go
package database

const schema = `
CREATE TABLE IF NOT EXISTS feeds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    category TEXT,
    site_url TEXT,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT 1,
    last_fetched_at DATETIME,
    error_count INTEGER DEFAULT 0,
    last_error TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_feeds_is_active ON feeds(is_active);
CREATE INDEX IF NOT EXISTS idx_feeds_last_fetched ON feeds(last_fetched_at);

CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    feed_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    link TEXT NOT NULL,
    description TEXT,
    content TEXT,
    author TEXT,
    published_at DATETIME,
    image_url TEXT,
    guid TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
    UNIQUE(feed_id, guid)
);

CREATE INDEX IF NOT EXISTS idx_posts_feed_id ON posts(feed_id);
CREATE INDEX IF NOT EXISTS idx_posts_published_at ON posts(published_at DESC);
CREATE INDEX IF NOT EXISTS idx_posts_is_read ON posts(is_read);
CREATE UNIQUE INDEX IF NOT EXISTS idx_posts_feed_guid ON posts(feed_id, guid);
`

// InitSchema initializes the database schema
func (db *DB) InitSchema() error {
    _, err := db.Exec(schema)
    if err != nil {
        return fmt.Errorf("failed to initialize schema: %w", err)
    }
    log.Println("Database schema initialized successfully")
    return nil
}
```

**Step 1.3: Create seed.go for default feeds**

Create [backend/internal/database/seed.go](backend/internal/database/seed.go):

```go
package database

import (
    "log"
)

type SeedFeed struct {
    Name        string
    URL         string
    Category    string
    SiteURL     string
    Description string
}

var defaultFeeds = []SeedFeed{
    {
        Name:        "Hacker News",
        URL:         "https://news.ycombinator.com/rss",
        Category:    "Tech",
        SiteURL:     "https://news.ycombinator.com",
        Description: "Hacker News RSS Feed",
    },
    {
        Name:        "TechCrunch",
        URL:         "https://techcrunch.com/feed/",
        Category:    "Tech",
        SiteURL:     "https://techcrunch.com",
        Description: "TechCrunch latest articles",
    },
    {
        Name:        "Reddit - Programming",
        URL:         "https://www.reddit.com/r/programming/.rss",
        Category:    "Tech",
        SiteURL:     "https://www.reddit.com/r/programming",
        Description: "Programming subreddit feed",
    },
    {
        Name:        "Ars Technica",
        URL:         "https://feeds.arstechnica.com/arstechnica/index",
        Category:    "Tech",
        SiteURL:     "https://arstechnica.com",
        Description: "Ars Technica RSS Feed",
    },
}

// SeedDefaultFeeds inserts default feeds if database is empty
func (db *DB) SeedDefaultFeeds() error {
    // Check if any feeds exist
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM feeds").Scan(&count)
    if err != nil {
        return err
    }

    if count > 0 {
        log.Println("Feeds already exist, skipping seed")
        return nil
    }

    log.Println("Seeding default feeds...")

    stmt, err := db.Prepare(`
        INSERT INTO feeds (name, url, category, site_url, description)
        VALUES (?, ?, ?, ?, ?)
    `)
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, feed := range defaultFeeds {
        _, err := stmt.Exec(feed.Name, feed.URL, feed.Category, feed.SiteURL, feed.Description)
        if err != nil {
            log.Printf("Failed to seed feed %s: %v", feed.Name, err)
            continue
        }
        log.Printf("Seeded feed: %s", feed.Name)
    }

    log.Println("Default feeds seeded successfully")
    return nil
}
```

### 2. Create Model Structs

**Step 2.1: Create feed model**

Create [backend/internal/models/feed.go](backend/internal/models/feed.go):

```go
package models

import "time"

type Feed struct {
    ID            int64      `json:"id"`
    Name          string     `json:"name"`
    URL           string     `json:"url"`
    Category      string     `json:"category"`
    SiteURL       string     `json:"site_url"`
    Description   string     `json:"description"`
    IsActive      bool       `json:"is_active"`
    LastFetchedAt *time.Time `json:"last_fetched_at"`
    ErrorCount    int        `json:"error_count"`
    LastError     string     `json:"last_error"`
    CreatedAt     time.Time  `json:"created_at"`
    UpdatedAt     time.Time  `json:"updated_at"`
}

type CreateFeedRequest struct {
    Name        string `json:"name"`
    URL         string `json:"url"`
    Category    string `json:"category"`
    SiteURL     string `json:"site_url"`
    Description string `json:"description"`
}

type UpdateFeedRequest struct {
    Name        *string `json:"name"`
    URL         *string `json:"url"`
    Category    *string `json:"category"`
    SiteURL     *string `json:"site_url"`
    Description *string `json:"description"`
    IsActive    *bool   `json:"is_active"`
}
```

**Step 2.2: Create post model**

Create [backend/internal/models/post.go](backend/internal/models/post.go):

```go
package models

import "time"

type Post struct {
    ID          int64      `json:"id"`
    FeedID      int64      `json:"feed_id"`
    Title       string     `json:"title"`
    Link        string     `json:"link"`
    Description string     `json:"description"`
    Content     string     `json:"content"`
    Author      string     `json:"author"`
    PublishedAt *time.Time `json:"published_at"`
    ImageURL    string     `json:"image_url"`
    GUID        string     `json:"guid"`
    IsRead      bool       `json:"is_read"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

type PostWithFeed struct {
    Post
    FeedName string `json:"feed_name"`
}
```

### 3. Create Repository Layer

**Step 3.1: Create feed repository**

Create [backend/internal/database/feed_repository.go](backend/internal/database/feed_repository.go):

```go
package database

import (
    "database/sql"
    "fmt"
    "time"

    "github.com/justanotherspy/rssy/internal/models"
)

// GetAllFeeds retrieves all feeds
func (db *DB) GetAllFeeds() ([]models.Feed, error) {
    query := `
        SELECT id, name, url, category, site_url, description, is_active,
               last_fetched_at, error_count, last_error, created_at, updated_at
        FROM feeds
        ORDER BY name ASC
    `

    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    feeds := []models.Feed{}
    for rows.Next() {
        var feed models.Feed
        err := rows.Scan(
            &feed.ID, &feed.Name, &feed.URL, &feed.Category, &feed.SiteURL,
            &feed.Description, &feed.IsActive, &feed.LastFetchedAt,
            &feed.ErrorCount, &feed.LastError, &feed.CreatedAt, &feed.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        feeds = append(feeds, feed)
    }

    return feeds, nil
}

// GetFeedByID retrieves a feed by ID
func (db *DB) GetFeedByID(id int64) (*models.Feed, error) {
    query := `
        SELECT id, name, url, category, site_url, description, is_active,
               last_fetched_at, error_count, last_error, created_at, updated_at
        FROM feeds
        WHERE id = ?
    `

    var feed models.Feed
    err := db.QueryRow(query, id).Scan(
        &feed.ID, &feed.Name, &feed.URL, &feed.Category, &feed.SiteURL,
        &feed.Description, &feed.IsActive, &feed.LastFetchedAt,
        &feed.ErrorCount, &feed.LastError, &feed.CreatedAt, &feed.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("feed not found")
    }
    if err != nil {
        return nil, err
    }

    return &feed, nil
}

// CreateFeed creates a new feed
func (db *DB) CreateFeed(req models.CreateFeedRequest) (*models.Feed, error) {
    query := `
        INSERT INTO feeds (name, url, category, site_url, description)
        VALUES (?, ?, ?, ?, ?)
        RETURNING id, name, url, category, site_url, description, is_active,
                  last_fetched_at, error_count, last_error, created_at, updated_at
    `

    var feed models.Feed
    err := db.QueryRow(
        query, req.Name, req.URL, req.Category, req.SiteURL, req.Description,
    ).Scan(
        &feed.ID, &feed.Name, &feed.URL, &feed.Category, &feed.SiteURL,
        &feed.Description, &feed.IsActive, &feed.LastFetchedAt,
        &feed.ErrorCount, &feed.LastError, &feed.CreatedAt, &feed.UpdatedAt,
    )

    if err != nil {
        return nil, err
    }

    return &feed, nil
}

// UpdateFeed updates an existing feed
func (db *DB) UpdateFeed(id int64, req models.UpdateFeedRequest) (*models.Feed, error) {
    // Build dynamic update query
    query := "UPDATE feeds SET updated_at = CURRENT_TIMESTAMP"
    args := []interface{}{}

    if req.Name != nil {
        query += ", name = ?"
        args = append(args, *req.Name)
    }
    if req.URL != nil {
        query += ", url = ?"
        args = append(args, *req.URL)
    }
    if req.Category != nil {
        query += ", category = ?"
        args = append(args, *req.Category)
    }
    if req.SiteURL != nil {
        query += ", site_url = ?"
        args = append(args, *req.SiteURL)
    }
    if req.Description != nil {
        query += ", description = ?"
        args = append(args, *req.Description)
    }
    if req.IsActive != nil {
        query += ", is_active = ?"
        args = append(args, *req.IsActive)
    }

    query += " WHERE id = ?"
    args = append(args, id)

    _, err := db.Exec(query, args...)
    if err != nil {
        return nil, err
    }

    return db.GetFeedByID(id)
}

// DeleteFeed deletes a feed
func (db *DB) DeleteFeed(id int64) error {
    _, err := db.Exec("DELETE FROM feeds WHERE id = ?", id)
    return err
}

// UpdateFeedLastFetched updates the last fetched timestamp
func (db *DB) UpdateFeedLastFetched(id int64, fetchTime time.Time) error {
    _, err := db.Exec(
        "UPDATE feeds SET last_fetched_at = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
        fetchTime, id,
    )
    return err
}
```

**Step 3.2: Create post repository**

Create [backend/internal/database/post_repository.go](backend/internal/database/post_repository.go):

```go
package database

import (
    "database/sql"
    "fmt"

    "github.com/justanotherspy/rssy/internal/models"
)

// GetAllPosts retrieves all posts with pagination
func (db *DB) GetAllPosts(limit, offset int) ([]models.PostWithFeed, error) {
    query := `
        SELECT p.id, p.feed_id, p.title, p.link, p.description, p.content,
               p.author, p.published_at, p.image_url, p.guid, p.is_read,
               p.created_at, p.updated_at, f.name as feed_name
        FROM posts p
        JOIN feeds f ON p.feed_id = f.id
        ORDER BY p.published_at DESC
        LIMIT ? OFFSET ?
    `

    rows, err := db.Query(query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    posts := []models.PostWithFeed{}
    for rows.Next() {
        var post models.PostWithFeed
        err := rows.Scan(
            &post.ID, &post.FeedID, &post.Title, &post.Link, &post.Description,
            &post.Content, &post.Author, &post.PublishedAt, &post.ImageURL,
            &post.GUID, &post.IsRead, &post.CreatedAt, &post.UpdatedAt,
            &post.FeedName,
        )
        if err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }

    return posts, nil
}

// GetPostsByFeedID retrieves posts for a specific feed
func (db *DB) GetPostsByFeedID(feedID int64, limit, offset int) ([]models.Post, error) {
    query := `
        SELECT id, feed_id, title, link, description, content, author,
               published_at, image_url, guid, is_read, created_at, updated_at
        FROM posts
        WHERE feed_id = ?
        ORDER BY published_at DESC
        LIMIT ? OFFSET ?
    `

    rows, err := db.Query(query, feedID, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    posts := []models.Post{}
    for rows.Next() {
        var post models.Post
        err := rows.Scan(
            &post.ID, &post.FeedID, &post.Title, &post.Link, &post.Description,
            &post.Content, &post.Author, &post.PublishedAt, &post.ImageURL,
            &post.GUID, &post.IsRead, &post.CreatedAt, &post.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }

    return posts, nil
}

// CreatePost creates a new post (used by feed fetcher)
func (db *DB) CreatePost(post *models.Post) error {
    query := `
        INSERT INTO posts (feed_id, title, link, description, content, author,
                          published_at, image_url, guid)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

    result, err := db.Exec(
        query, post.FeedID, post.Title, post.Link, post.Description,
        post.Content, post.Author, post.PublishedAt, post.ImageURL, post.GUID,
    )
    if err != nil {
        return err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return err
    }

    post.ID = id
    return nil
}

// MarkPostAsRead marks a post as read
func (db *DB) MarkPostAsRead(id int64, isRead bool) error {
    _, err := db.Exec("UPDATE posts SET is_read = ? WHERE id = ?", isRead, id)
    return err
}

// DeleteAllPosts deletes all posts (for reset functionality)
func (db *DB) DeleteAllPosts() error {
    _, err := db.Exec("DELETE FROM posts")
    return err
}

// GetPostByGUID checks if a post exists by GUID
func (db *DB) GetPostByGUID(feedID int64, guid string) (*models.Post, error) {
    query := `
        SELECT id, feed_id, title, link, description, content, author,
               published_at, image_url, guid, is_read, created_at, updated_at
        FROM posts
        WHERE feed_id = ? AND guid = ?
    `

    var post models.Post
    err := db.QueryRow(query, feedID, guid).Scan(
        &post.ID, &post.FeedID, &post.Title, &post.Link, &post.Description,
        &post.Content, &post.Author, &post.PublishedAt, &post.ImageURL,
        &post.GUID, &post.IsRead, &post.CreatedAt, &post.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    return &post, nil
}
```

### 4. Update main.go to Initialize Database

Update [backend/cmd/api/main.go](backend/cmd/api/main.go):

```go
package main

import (
    "log"
    "os"

    "github.com/justanotherspy/rssy/internal/database"
)

func main() {
    log.Println("Starting RSSY API Server...")

    // Get database path from environment or use default
    dbPath := os.Getenv("DATABASE_PATH")
    if dbPath == "" {
        dbPath = "./rssy.db"
    }

    // Initialize database
    db, err := database.New(dbPath)
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

    log.Println("Database initialized successfully")
    log.Println("Server ready")
}
```

### 5. Testing the Data Layer

**Step 5.1: Build and run**
```bash
cd backend
go mod tidy
go run ./cmd/api
```

Should output:
```
Starting RSSY API Server...
Connected to database: ./rssy.db
Database schema initialized successfully
Seeding default feeds...
Seeded feed: Hacker News
Seeded feed: TechCrunch
Seeded feed: Reddit - Programming
Seeded feed: Ars Technica
Default feeds seeded successfully
Database initialized successfully
Server ready
```

**Step 5.2: Verify database**
```bash
sqlite3 backend/rssy.db
```

Run queries:
```sql
.tables
.schema feeds
.schema posts
SELECT * FROM feeds;
```

**Step 5.3: Create test file**

Create [backend/internal/database/database_test.go](backend/internal/database/database_test.go) for basic tests.

## Success Criteria

- Database schema created with proper tables and indexes
- Foreign key constraints properly defined
- Go models defined for Feed and Post entities
- Repository layer implements all CRUD operations
- Database initializes and seeds default feeds
- No errors when running the application
- Can query database and see seeded feeds

## Next Steps

After completing Task 2, proceed to Task 3 (Go Backend API) where we will:
- Create HTTP handlers for API endpoints
- Implement feed CRUD endpoints
- Implement post retrieval endpoints
- Build RSS feed fetching service
- Implement background polling mechanism
- Add CORS and middleware

## Notes

- SQLite RETURNING clause requires SQLite 3.35.0+
- The UNIQUE constraint on (feed_id, guid) prevents duplicate posts
- CASCADE DELETE ensures posts are deleted when their feed is deleted
- Indexes optimize common query patterns
- The repository pattern abstracts database operations

## Agent & Hook Recommendations for This Task

### Recommended Agent Usage
- **Direct tool usage** for creating files
- **General-purpose agent** if complex refactoring needed
- Test database operations manually before proceeding

### Testing Approach
- Manual testing with sqlite3 CLI
- Unit tests can be added later for repository methods
- Focus on ensuring schema correctness first

## Troubleshooting

**RETURNING clause not supported:**
- Update SQLite to 3.35.0 or later
- Or split INSERT and SELECT operations

**Foreign key not enforced:**
- Ensure `?_foreign_keys=on` in connection string

**Duplicate entry errors:**
- Check UNIQUE constraints are properly defined
- Handle errors gracefully in repository layer
