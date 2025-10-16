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
