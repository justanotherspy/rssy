package handlers

import (
	"encoding/json"
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
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0 // default
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
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
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
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
