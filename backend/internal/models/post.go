package models

import "time"

type Post struct {
	ID          int64      `json:"id"`
	FeedID      int64      `json:"feed_id"`
	Title       string     `json:"title"`
	Link        string     `json:"link"`
	Description string     `json:"description"`
	Content     string     `json:"content"`
	Author      string     `json:"author"`
	PublishedAt *time.Time `json:"published_at"`
	ImageURL    string     `json:"image_url"`
	GUID        string     `json:"guid"`
	IsRead      bool       `json:"is_read"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type PostWithFeed struct {
	Post
	FeedName string `json:"feed_name"`
}
