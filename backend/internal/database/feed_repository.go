package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/justanotherspy/rssy/internal/models"
)

// GetAllFeeds retrieves all feeds
func (db *DB) GetAllFeeds() ([]models.Feed, error) {
	query := `
        SELECT id, name, url, category, site_url, description, is_active,
               last_fetched_at, error_count, last_error, created_at, updated_at
        FROM feeds
        ORDER BY name ASC
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feeds := []models.Feed{}
	for rows.Next() {
		var feed models.Feed
		err := rows.Scan(
			&feed.ID, &feed.Name, &feed.URL, &feed.Category, &feed.SiteURL,
			&feed.Description, &feed.IsActive, &feed.LastFetchedAt,
			&feed.ErrorCount, &feed.LastError, &feed.CreatedAt, &feed.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}

	return feeds, nil
}

// GetFeedByID retrieves a feed by ID
func (db *DB) GetFeedByID(id int64) (*models.Feed, error) {
	query := `
        SELECT id, name, url, category, site_url, description, is_active,
               last_fetched_at, error_count, last_error, created_at, updated_at
        FROM feeds
        WHERE id = ?
    `

	var feed models.Feed
	err := db.QueryRow(query, id).Scan(
		&feed.ID, &feed.Name, &feed.URL, &feed.Category, &feed.SiteURL,
		&feed.Description, &feed.IsActive, &feed.LastFetchedAt,
		&feed.ErrorCount, &feed.LastError, &feed.CreatedAt, &feed.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("feed not found")
	}
	if err != nil {
		return nil, err
	}

	return &feed, nil
}

// CreateFeed creates a new feed
func (db *DB) CreateFeed(req models.CreateFeedRequest) (*models.Feed, error) {
	query := `
        INSERT INTO feeds (name, url, category, site_url, description)
        VALUES (?, ?, ?, ?, ?)
        RETURNING id, name, url, category, site_url, description, is_active,
                  last_fetched_at, error_count, last_error, created_at, updated_at
    `

	var feed models.Feed
	err := db.QueryRow(
		query, req.Name, req.URL, req.Category, req.SiteURL, req.Description,
	).Scan(
		&feed.ID, &feed.Name, &feed.URL, &feed.Category, &feed.SiteURL,
		&feed.Description, &feed.IsActive, &feed.LastFetchedAt,
		&feed.ErrorCount, &feed.LastError, &feed.CreatedAt, &feed.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &feed, nil
}

// UpdateFeed updates an existing feed
func (db *DB) UpdateFeed(id int64, req models.UpdateFeedRequest) (*models.Feed, error) {
	// Build dynamic update query
	query := "UPDATE feeds SET updated_at = CURRENT_TIMESTAMP"
	args := []interface{}{}

	if req.Name != nil {
		query += ", name = ?"
		args = append(args, *req.Name)
	}
	if req.URL != nil {
		query += ", url = ?"
		args = append(args, *req.URL)
	}
	if req.Category != nil {
		query += ", category = ?"
		args = append(args, *req.Category)
	}
	if req.SiteURL != nil {
		query += ", site_url = ?"
		args = append(args, *req.SiteURL)
	}
	if req.Description != nil {
		query += ", description = ?"
		args = append(args, *req.Description)
	}
	if req.IsActive != nil {
		query += ", is_active = ?"
		args = append(args, *req.IsActive)
	}

	query += " WHERE id = ?"
	args = append(args, id)

	_, err := db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return db.GetFeedByID(id)
}

// DeleteFeed deletes a feed
func (db *DB) DeleteFeed(id int64) error {
	_, err := db.Exec("DELETE FROM feeds WHERE id = ?", id)
	return err
}

// UpdateFeedLastFetched updates the last fetched timestamp
func (db *DB) UpdateFeedLastFetched(id int64, fetchTime time.Time) error {
	_, err := db.Exec(
		"UPDATE feeds SET last_fetched_at = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		fetchTime, id,
	)
	return err
}
