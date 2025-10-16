package models

import "time"

type Feed struct {
	ID            int64      `json:"id"`
	Name          string     `json:"name"`
	URL           string     `json:"url"`
	Category      *string    `json:"category"`
	SiteURL       *string    `json:"site_url"`
	Description   *string    `json:"description"`
	IsActive      bool       `json:"is_active"`
	LastFetchedAt *time.Time `json:"last_fetched_at"`
	ErrorCount    int        `json:"error_count"`
	LastError     *string    `json:"last_error"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type CreateFeedRequest struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Category    string `json:"category"`
	SiteURL     string `json:"site_url"`
	Description string `json:"description"`
}

type UpdateFeedRequest struct {
	Name        *string `json:"name"`
	URL         *string `json:"url"`
	Category    *string `json:"category"`
	SiteURL     *string `json:"site_url"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}
