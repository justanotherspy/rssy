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
