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
		if enclosure.Type != "" && len(enclosure.Type) >= 5 && enclosure.Type[:5] == "image" {
			return enclosure.URL
		}
	}

	return ""
}
