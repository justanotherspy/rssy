package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/justanotherspy/rssy/internal/models"
	"github.com/justanotherspy/rssy/internal/services"
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
		Name:        "r/" + req.Subreddit,
		URL:         "https://www.reddit.com/r/" + req.Subreddit + "/.rss",
		Category:    "Reddit",
		SiteURL:     "https://www.reddit.com/r/" + req.Subreddit,
		Description: "Reddit /r/" + req.Subreddit + " feed",
	}

	feed, err := h.db.CreateFeed(feedReq)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create Reddit feed")
		return
	}

	h.respondJSON(w, http.StatusCreated, feed)
}

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
