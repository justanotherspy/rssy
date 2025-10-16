package database

import (
	"database/sql"

	"github.com/justanotherspy/rssy/internal/models"
)

// GetAllPosts retrieves all posts with pagination
func (db *DB) GetAllPosts(limit, offset int) ([]models.PostWithFeed, error) {
	query := `
        SELECT p.id, p.feed_id, p.title, p.link, p.description, p.content,
               p.author, p.published_at, p.image_url, p.guid, p.is_read,
               p.created_at, p.updated_at, f.name as feed_name
        FROM posts p
        JOIN feeds f ON p.feed_id = f.id
        ORDER BY p.published_at DESC
        LIMIT ? OFFSET ?
    `

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []models.PostWithFeed{}
	for rows.Next() {
		var post models.PostWithFeed
		err := rows.Scan(
			&post.ID, &post.FeedID, &post.Title, &post.Link, &post.Description,
			&post.Content, &post.Author, &post.PublishedAt, &post.ImageURL,
			&post.GUID, &post.IsRead, &post.CreatedAt, &post.UpdatedAt,
			&post.FeedName,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// GetPostsByFeedID retrieves posts for a specific feed
func (db *DB) GetPostsByFeedID(feedID int64, limit, offset int) ([]models.Post, error) {
	query := `
        SELECT id, feed_id, title, link, description, content, author,
               published_at, image_url, guid, is_read, created_at, updated_at
        FROM posts
        WHERE feed_id = ?
        ORDER BY published_at DESC
        LIMIT ? OFFSET ?
    `

	rows, err := db.Query(query, feedID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID, &post.FeedID, &post.Title, &post.Link, &post.Description,
			&post.Content, &post.Author, &post.PublishedAt, &post.ImageURL,
			&post.GUID, &post.IsRead, &post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// CreatePost creates a new post (used by feed fetcher)
func (db *DB) CreatePost(post *models.Post) error {
	query := `
        INSERT INTO posts (feed_id, title, link, description, content, author,
                          published_at, image_url, guid)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	result, err := db.Exec(
		query, post.FeedID, post.Title, post.Link, post.Description,
		post.Content, post.Author, post.PublishedAt, post.ImageURL, post.GUID,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	post.ID = id
	return nil
}

// MarkPostAsRead marks a post as read
func (db *DB) MarkPostAsRead(id int64, isRead bool) error {
	_, err := db.Exec("UPDATE posts SET is_read = ? WHERE id = ?", isRead, id)
	return err
}

// DeleteAllPosts deletes all posts (for reset functionality)
func (db *DB) DeleteAllPosts() error {
	_, err := db.Exec("DELETE FROM posts")
	return err
}

// GetPostByGUID checks if a post exists by GUID
func (db *DB) GetPostByGUID(feedID int64, guid string) (*models.Post, error) {
	query := `
        SELECT id, feed_id, title, link, description, content, author,
               published_at, image_url, guid, is_read, created_at, updated_at
        FROM posts
        WHERE feed_id = ? AND guid = ?
    `

	var post models.Post
	err := db.QueryRow(query, feedID, guid).Scan(
		&post.ID, &post.FeedID, &post.Title, &post.Link, &post.Description,
		&post.Content, &post.Author, &post.PublishedAt, &post.ImageURL,
		&post.GUID, &post.IsRead, &post.CreatedAt, &post.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &post, nil
}
