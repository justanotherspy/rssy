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
