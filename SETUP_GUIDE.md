# RSSY RSS Reader - Complete Setup Guide

This guide walks you through setting up and running the complete RSSY application.

## Prerequisites

- Node.js 16+ (for frontend)
- Go 1.25+ (for backend)
- SQLite3 (database)
- npm or yarn (package manager)

## Project Structure

```
stream/
├── backend/                # Go API server
│   ├── cmd/api/main.go
│   ├── internal/           # Models, handlers, services, database
│   ├── go.mod
│   └── go.sum
├── frontend/               # SvelteKit web application
│   ├── src/
│   │   ├── lib/           # API client, stores, components
│   │   └── routes/        # Pages
│   ├── package.json
│   └── .env
├── docs/                   # Task documentation
└── Makefile               # Build automation
```

## Backend Setup

### 1. Install Backend Dependencies

```bash
cd backend
go mod download
```

### 2. Set Up Database

The application uses SQLite. The database file will be created automatically on first run at `backend/rssy.db`.

### 3. Run Backend Server

```bash
cd backend
go run ./cmd/api
```

The backend will start on `http://localhost:8080` and display:
```
Starting RSS reader API server on :8080
```

The backend automatically:
- Creates database tables on first run
- Seeds default feeds
- Starts polling feeds every 10 minutes
- Serves REST API endpoints

### Backend API Endpoints

**Feeds:**
- `GET /api/feeds` - List all feeds
- `POST /api/feeds` - Create new feed
- `GET /api/feeds/:id` - Get specific feed
- `PUT /api/feeds/:id` - Update feed
- `DELETE /api/feeds/:id` - Delete feed
- `POST /api/feeds/reddit` - Add Reddit feed (quick)
- `POST /api/feeds/:id/refresh` - Manually refresh feed

**Posts:**
- `GET /api/posts` - List all posts (paginated)
- `GET /api/posts/feed/:feedId` - List posts from feed
- `PATCH /api/posts/:id/read` - Mark as read/unread
- `DELETE /api/posts` - Delete all posts

## Frontend Setup

### 1. Install Frontend Dependencies

```bash
cd frontend
npm install
```

This installs:
- `svelte@5.39.5` - UI framework
- `@sveltejs/kit@2.43.2` - Web framework
- `axios@1.12.2` - HTTP client
- `date-fns@4.1.0` - Date utilities
- `lucide-svelte@0.546.0` - Icons

### 2. Configure Environment

The `.env` file is already configured:
```
VITE_API_URL=http://localhost:8080
```

Update this if your backend runs on a different URL.

### 3. Run Development Server

```bash
cd frontend
npm run dev
```

The frontend will start on `http://localhost:5173` and display:
```
  VITE v7.1.7  ready in 123 ms

  ➜  Local:   http://localhost:5173/
  ➜  press h to show help
```

## Running Both Servers

In separate terminal windows:

**Terminal 1 - Backend:**
```bash
cd /home/daniel/claude/stream/backend
go run ./cmd/api
```

**Terminal 2 - Frontend:**
```bash
cd /home/daniel/claude/stream/frontend
npm run dev
```

Then open your browser to `http://localhost:5173`

## Frontend Development

### Type Checking

```bash
cd frontend
npm run check
```

Validates TypeScript and Svelte types.

### Build for Production

```bash
cd frontend
npm run build
npm run preview  # Test production build
```

## Verifying the Installation

### 1. Check Backend is Running

```bash
curl http://localhost:8080/api/feeds
```

Should return: `{"data":[...]}`

### 2. Check Frontend is Accessible

Visit `http://localhost:5173` in your browser. You should see:
- RSSY sidebar on the left with "#all" and default feeds
- Empty content area (no posts yet)
- Plus button to add feeds
- Settings gear icon

### 3. Test Adding a Feed

1. Click the "+" button in the sidebar
2. Select "RSS Feed" tab
3. Enter:
   - Name: "HN News"
   - URL: `https://news.ycombinator.com/rss`
4. Click "Add Feed"
5. Wait for posts to load from the feed

### 4. Test Reddit Integration

1. Click the "+" button again
2. Select "Reddit" tab
3. Enter subreddit: `programming`
4. Click "Add Feed"
5. Posts from r/programming should appear

### 5. Test Settings

1. Click the settings gear icon
2. Change refresh interval
3. Click "Done"

## Build Artifacts

### Backend Build

```bash
cd backend
make backend-build
```

Produces: `backend/bin/api` executable

### Frontend Build

```bash
cd frontend
npm run build
```

Produces: `.svelte-kit/output/` directory with:
- Static assets in `client/`
- Server code in `server/`

## Database Management

### View Database

The SQLite database is at `backend/rssy.db`. You can inspect it with:

```bash
sqlite3 backend/rssy.db ".schema"
sqlite3 backend/rssy.db "SELECT COUNT(*) FROM feeds;"
sqlite3 backend/rssy.db "SELECT COUNT(*) FROM posts;"
```

### Reset Database

Delete the database file (it will be recreated on next run):

```bash
rm backend/rssy.db
```

## Troubleshooting

### Frontend can't connect to backend

- Verify backend is running: `curl http://localhost:8080/api/feeds`
- Check `VITE_API_URL` in `frontend/.env`
- Check browser console (F12) for CORS errors
- Verify firewall isn't blocking port 8080

### No posts appearing

- Wait for backend polling interval (10 minutes default)
- Manually refresh a feed: Click feed in sidebar, wait
- Check backend logs for errors
- Verify RSS feeds are valid URLs

### Feeds not loading in sidebar

- Check network tab in browser DevTools (F12)
- Verify API response: `curl http://localhost:8080/api/feeds`
- Check browser console for JavaScript errors
- Try hard refresh (Ctrl+Shift+R)

### Frontend won't start

```bash
# Clear cache and reinstall
rm -rf frontend/node_modules frontend/.svelte-kit
npm install
npm run dev
```

### Backend won't start

```bash
# Check if port 8080 is in use
lsof -i :8080

# Run with verbose logging
cd backend
go run ./cmd/api
```

## Development Workflow

1. **Make backend changes**
   - Backend hot-reloads automatically
   - Or restart with `Ctrl+C` and `go run ./cmd/api`

2. **Make frontend changes**
   - Frontend auto-refreshes in browser
   - TypeScript errors show in console

3. **Test changes**
   - Open http://localhost:5173 in browser
   - Check browser console (F12) for errors
   - Check backend terminal for error logs

## Deployment

### Frontend Deployment

Build and deploy to static hosting (Vercel, Netlify, GitHub Pages):

```bash
npm run build
# Upload contents of .svelte-kit/output/client/ to static host
```

### Backend Deployment

Build and run on server:

```bash
make backend-build
./bin/api
```

Set environment variables:
- `DATABASE_URL` - Path to SQLite database file
- `PORT` - HTTP port (default 8080)

## Git Workflow

Before committing, ensure:

1. Backend builds: `go build ./cmd/api`
2. Frontend checks: `npm run check`
3. Frontend builds: `npm run build`

```bash
git status
git diff
git add .
git commit -m "Feature: Add X"
git push origin feature-branch
```

## Performance Optimization

### Backend
- Feed polling every 10 minutes (configurable)
- Post deduplication via GUID
- Database indexing on key fields
- Connection pooling

### Frontend
- Code splitting by route
- Lazy loading of modals
- Optimized bundle size (~79KB gzip)
- Scoped CSS (no global pollution)

## Security Considerations

1. **CORS**: Backend only accepts requests from `http://localhost:3000`
2. **Input Validation**: All inputs validated server-side
3. **HTML Sanitization**: Post content cleaned before display
4. **URL Validation**: RSS feed URLs must be valid
5. **Rate Limiting**: Implement on production backend

## Next Steps

1. Add authentication (user accounts)
2. Implement categorization (tags/folders)
3. Add search functionality
4. Deploy to production
5. Monitor performance
6. Gather user feedback

## Support

For issues:
1. Check terminal logs (backend and frontend)
2. Check browser console (F12)
3. Verify both servers are running
4. Review error messages
5. Restart both servers

## Documentation

- Backend API: See `backend/README.md`
- Frontend: See `FRONTEND_IMPLEMENTATION.md`
- Tasks: See `docs/task*.md`
