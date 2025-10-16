package router

import (
	"net/http"

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
			r.Post("/refresh", h.RefreshAllFeeds)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.GetFeedByID)
				r.Put("/", h.UpdateFeed)
				r.Delete("/", h.DeleteFeed)
				r.Post("/refresh", h.RefreshFeed)
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
