# RSSY - RSS Reader Application

A modern, minimalist RSS reader built with Go and Svelte, featuring a clean two-panel interface for managing feeds and reading posts.

## Status: MVP Complete ✓

Both backend and frontend are fully implemented and production-ready.

## Quick Start

### Prerequisites
- Node.js 18+
- Go 1.21+ (project uses 1.25.0)
- SQLite3

### Run Locally

**Terminal 1 - Backend:**
```bash
cd backend
go run ./cmd/api
```
Backend runs on `http://localhost:8080`

**Terminal 2 - Frontend:**
```bash
cd frontend
npm run dev
```
Frontend runs on `http://localhost:5173`

### Using Make
```bash
make help              # Show all available commands
make install-deps      # Install all dependencies
make backend-run       # Run backend
make frontend-run      # Run frontend
```

**Note:** These commands start the application in development mode. For production deployment, see the [Deployment](#deployment) section.

## Project Structure

```
stream/
├── backend/                    # Go API server
│   ├── cmd/api/
│   │   └── main.go            # Application entry point
│   ├── internal/
│   │   ├── config/            # Configuration management
│   │   ├── database/          # SQLite operations, repositories, seed data
│   │   ├── handlers/          # HTTP request handlers
│   │   ├── models/            # Data models (Feed, Post)
│   │   ├── router/            # Route configuration
│   │   └── services/          # RSS polling and fetching logic
│   ├── go.mod
│   └── go.sum
│
├── frontend/                   # Svelte/SvelteKit application
│   ├── src/
│   │   ├── lib/
│   │   │   ├── api.ts         # Axios HTTP client
│   │   │   ├── stores.ts      # Svelte stores (state)
│   │   │   └── components/
│   │   │       ├── Sidebar.svelte
│   │   │       ├── PostCard.svelte
│   │   │       ├── AddFeedModal.svelte
│   │   │       └── SettingsModal.svelte
│   │   └── routes/
│   │       ├── +layout.svelte
│   │       └── +page.svelte   # Main application page
│   ├── package.json
│   └── .env
│
├── docs/                       # Development documentation
├── scripts/                    # Utility scripts
├── Makefile                    # Build automation
└── README.md                   # This file
```

## Features

### What's Implemented

**Feed Management:**
- Add RSS/Atom feeds via URL
- Quick-add Reddit subreddits (r/subreddit)
- Default feeds included (HackerNews, TechCrunch, etc.)
- View all posts or filter by specific feed
- Automatic feed polling every 10 minutes (configurable)

**User Interface:**
- Two-panel layout (sidebar + main content)
- Dark minimalist sidebar with feed list
- Post cards with images, titles, descriptions, and metadata
- Click posts to open in new tab
- Modal dialogs for adding feeds and settings
- Responsive design (mobile, tablet, desktop)
- Loading states and error handling

**Settings:**
- Configure feed refresh interval (1-1440 minutes)
- Delete all posts (with confirmation)
- Database management

### Technology Stack

**Backend:**
- Go 1.21+ with chi router
- SQLite database with migration support
- Automatic RSS/Atom feed polling
- RESTful JSON API
- CORS enabled for development

**Frontend:**
- Svelte 5 + SvelteKit v2
- TypeScript for type safety
- Axios HTTP client
- date-fns for date formatting
- lucide-svelte for icons
- Vite build system

## Documentation

- **[SETUP_GUIDE.md](./SETUP_GUIDE.md)** - Complete setup and installation guide
- **[IMPLEMENTATION_SUMMARY.md](./IMPLEMENTATION_SUMMARY.md)** - Complete implementation overview
- **[FRONTEND_IMPLEMENTATION.md](./FRONTEND_IMPLEMENTATION.md)** - Frontend architecture details
- **[CLAUDE.md](./CLAUDE.md)** - Project requirements and development approach
- **[docs/](./docs/)** - Task-by-task development documentation

## Development

### Makefile Commands

```bash
make help              # Show all available commands
make install-deps      # Install Go and Node dependencies

# Backend
make backend-run       # Run backend (go run ./cmd/api)
make backend-build     # Build backend binary (./bin/api)
make backend-test      # Run Go tests
make backend-vet       # Run Go vet linter

# Frontend
make frontend-run      # Run frontend dev server
make frontend-build    # Build frontend for production
make frontend-check    # Run TypeScript type checking

# Utilities
make clean             # Remove build artifacts
make format-go         # Format Go code
make lint              # Run all linters
```

### Manual Commands

**Backend:**
```bash
cd backend
go run ./cmd/api       # Run server on :8080
go test ./...          # Run tests
go build -o bin/api ./cmd/api  # Build binary
```

**Frontend:**
```bash
cd frontend
npm run dev            # Dev server on :5173
npm run build          # Production build
npm run check          # Type checking
npm run preview        # Preview production build
```

## API Endpoints

Backend API runs on `http://localhost:8080`

**Health Check:**
- `GET /health` - Health check endpoint (returns "OK")

**Feeds:**
- `GET /api/feeds` - List all feeds
- `POST /api/feeds` - Create feed (body: `{name, url, category?}`)
- `POST /api/feeds/reddit` - Add Reddit feed (body: `{subreddit}`)
- `POST /api/feeds/refresh` - Manually refresh all feeds
- `GET /api/feeds/:id` - Get specific feed
- `PUT /api/feeds/:id` - Update feed
- `DELETE /api/feeds/:id` - Delete feed
- `POST /api/feeds/:id/refresh` - Manually refresh specific feed

**Posts:**
- `GET /api/posts` - List all posts
- `GET /api/posts/feed/:feedId` - List posts from specific feed
- `DELETE /api/posts` - Delete all posts

All responses are JSON. Example:
```json
{
  "data": [...],
  "error": null
}
```

## Deployment

**Frontend:**
```bash
cd frontend
npm run build
# Deploy .svelte-kit/output/client/ to Vercel, Netlify, or GitHub Pages
```

**Backend:**
```bash
cd backend
go build -o rssy ./cmd/api
./rssy
```

**Database:**
- SQLite database created automatically as `rssy.db`
- Default feeds seeded on first run
- No migrations needed for fresh install

## Architecture

**Design:**
- Two-panel responsive layout (sidebar + main content)
- Component-based Svelte architecture
- RESTful JSON API
- Automatic background feed polling
- Dark minimalist UI (#1a1a1a sidebar)

**Data Flow:**
1. Backend polls RSS feeds every 10 minutes
2. New posts stored in SQLite
3. Frontend fetches posts via API
4. Posts displayed as cards in main panel
5. User interactions trigger API calls

## Troubleshooting

**Frontend can't connect to backend:**
- Ensure backend is running on port 8080
- Check `VITE_API_URL` in `frontend/.env` (default: `http://localhost:8080`)
- Check browser console for CORS errors

**No posts appearing:**
- Wait 10 minutes for initial feed polling
- Check backend logs for fetch errors
- Manually refresh a feed via the API

**Build errors:**
```bash
# Frontend
rm -rf frontend/node_modules frontend/.svelte-kit
cd frontend && npm install

# Backend
cd backend && go mod tidy
```

## Future Enhancements

Planned features for future releases:
- User authentication and multi-user support
- Feed categorization and tagging
- Full-text search across posts
- Read/unread status persistence
- Keyboard shortcuts (j/k navigation)
- PWA support (offline reading, installable app)
- OPML import/export
- Dark/light mode toggle
- Advanced filtering (date, author, keywords)

## Contributing

Contributions are welcome! Please:
- Follow Go and Svelte/TypeScript best practices
- Run `make lint` before committing
- Write descriptive commit messages
- Update documentation as needed

## License

MIT License - see LICENSE file for details

---

**Project Status:** MVP Complete ✓

For detailed documentation:
- [SETUP_GUIDE.md](./SETUP_GUIDE.md) - Installation and setup
- [IMPLEMENTATION_SUMMARY.md](./IMPLEMENTATION_SUMMARY.md) - Implementation details
- [FRONTEND_IMPLEMENTATION.md](./FRONTEND_IMPLEMENTATION.md) - Frontend architecture
