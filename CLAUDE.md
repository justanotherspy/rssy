# Let us build an RSS Reader called rssy pronounced rizzy

We will attempt to build a RSS Reader

- Name: RSSY
- Backend: Go
- Frontend: Svelte and SvelteKit
- Database: SQLite

Requirements for MVP:

- Go backend serves API
- API should serve a feed endpoint for CRUD operations
- API should serve a post endpoint for CRUD operations
- Go backend should poll feeds every interval period starting at 10 minutes
- Go backend manages the database and data model of feeds to posts.
- User can create new feeds from the front end
- Posts are downloaded to the database and retrived for the UI
- Frontend should have two panels one main content list that displays posts and another for a list of RSS feeds on the left.
- The left side bar lists the feeds as '#feedname'
- The left side bar highlights which of the feeds is selected to display in the main content feed
- The main content feed has a list of cards that represent posts.
- Each card should display all the info from the post
- Each card if the post has an image should show the image in full and put any caption underneath
- The card text should be fully justified and centered in the card.
- There should be a + button on the left panel for a modal that allows us to create a new feed to subscribe to
- There should also be a modal for editing a feed with a button next to the plus button.
- There should be a quick way to add reddit subreddits as feeds just by typing the subreddit name
- There should be a settings modal for changing the refresh interval or deleting all the posts from the db to reset
- There should be some feeds included by default

## Project Structure

```
stream/
├── CLAUDE.md              # Project documentation
├── backend/               # Go backend API
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   ├── internal/
│   │   ├── models/
│   │   ├── handlers/
│   │   ├── services/
│   │   └── database/
│   ├── go.mod
│   └── go.sum
├── frontend/              # Svelte/SvelteKit frontend
│   ├── src/
│   │   ├── routes/
│   │   ├── lib/
│   │   └── app.html
│   ├── static/
│   └── package.json
├── docs/                  # Task documentation
│   ├── task1.md
│   ├── task2.md
│   ├── task3.md
│   └── task4.md
└── Makefile              # Build automation
```

## Development Approach

### Session-Based Development
Each task will be tackled in a separate session to maintain focus and code quality:
- Task 1: Project setup and scaffolding
- Task 2: Data model and database design
- Task 3: Go backend API implementation
- Task 4: Svelte frontend implementation

### Technology Stack Details
- **Backend**: Go 1.25+ with standard library HTTP server
- **Frontend**: svelte@5.40.1 SvelteKit v2 with TypeScript
- **Database**: SQLite3 with appropriate Go driver
- **RSS Parsing**: Go RSS/Atom parser library
- **Build Tools**: Makefile for common tasks

### Agent & Hook Configuration

Use your tools more than you use bash to do things the tools can do already.

#### Recommended MCP Servers
Given the tech stack, consider installing these MCP servers:
1. **SQLite MCP**: For database inspection and queries
2. **Filesystem MCP**: For advanced file operations
3. **GitHub MCP**: Already configured for version control

#### Suggested Hooks
Create hooks in your Claude settings for this project:

**Pre-commit hook** (runs before git commits):
```bash
# Run Go tests and linters
cd backend && go test ./... && go vet ./...
# Run frontend build check
cd frontend && npm run check
```

**File save hook** (runs when files are saved):
```bash
# Format Go files
if [[ "$FILE_PATH" == *.go ]]; then
  gofmt -w "$FILE_PATH"
fi
# Format Svelte/TS files
if [[ "$FILE_PATH" == *.svelte ]] || [[ "$FILE_PATH" == *.ts ]]; then
  cd frontend && npx prettier --write "$FILE_PATH"
fi
```

### Claude Agent Best Practices for This Project

1. **Use Explore Agent**: When understanding code structure or searching for patterns
2. **Use General-Purpose Agent**: For complex multi-step refactoring
3. **Direct Tool Usage**: For specific, known file operations

### Git Workflow for This Project
Following Daniel's preferred workflow:
1. Create feature branch from latest master
2. Check status with `git status`
3. Review changes with `git diff`
4. Group related changes in commit messages
5. Push to new branch and create PR

### API Design Principles
- RESTful endpoints with clear naming
- JSON request/response format
- Proper HTTP status codes
- Error handling with structured responses

### Frontend Design Principles
- Component-based architecture
- Responsive design (mobile-first)
- Accessibility considerations
- Clean, minimal UI matching the RSS reader aesthetic
