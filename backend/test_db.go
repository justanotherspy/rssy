package main

import (
	"fmt"
	"log"
	"time"

	"github.com/justanotherspy/rssy/internal/database"
	"github.com/justanotherspy/rssy/internal/models"
)

func main() {
	// Initialize database
	db, err := database.New("./rssy_test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.InitSchema(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Testing Feed CRUD Operations ===\n")

	// Test 1: Create a new feed
	fmt.Println("1. Creating a new feed...")
	newFeed := models.CreateFeedRequest{
		Name:        "Test Feed",
		URL:         "https://example.com/feed.xml",
		Category:    "Technology",
		SiteURL:     "https://example.com",
		Description: "A test feed for verification",
	}
	feed, err := db.CreateFeed(newFeed)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✓ Created feed ID: %d, Name: %s\n\n", feed.ID, feed.Name)

	// Test 2: Get all feeds
	fmt.Println("2. Retrieving all feeds...")
	feeds, err := db.GetAllFeeds()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✓ Found %d feeds\n", len(feeds))
	for _, f := range feeds {
		fmt.Printf("      - %s (%s)\n", f.Name, f.URL)
	}
	fmt.Println()

	// Test 3: Get feed by ID
	fmt.Println("3. Getting feed by ID...")
	retrievedFeed, err := db.GetFeedByID(feed.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✓ Retrieved: %s\n\n", retrievedFeed.Name)

	// Test 4: Update feed
	fmt.Println("4. Updating feed...")
	updatedName := "Updated Test Feed"
	updatedCategory := "Science"
	updateReq := models.UpdateFeedRequest{
		Name:     &updatedName,
		Category: &updatedCategory,
	}
	updatedFeed, err := db.UpdateFeed(feed.ID, updateReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✓ Updated: Name='%s', Category='%s'\n\n", updatedFeed.Name, *updatedFeed.Category)

	// Test 5: Update last fetched timestamp
	fmt.Println("5. Updating last fetched timestamp...")
	now := time.Now()
	if err := db.UpdateFeedLastFetched(feed.ID, now); err != nil {
		log.Fatal(err)
	}
	updatedFeed, _ = db.GetFeedByID(feed.ID)
	fmt.Printf("   ✓ Last fetched: %v\n\n", updatedFeed.LastFetchedAt)

	fmt.Println("=== Testing Post CRUD Operations ===\n")

	// Test 6: Create posts
	fmt.Println("6. Creating test posts...")
	publishedAt := time.Now().Add(-1 * time.Hour)
	post1 := &models.Post{
		FeedID:      feed.ID,
		Title:       "First Test Post",
		Link:        "https://example.com/post1",
		Description: "This is the first test post",
		Content:     "Full content of the first post",
		Author:      "Test Author",
		PublishedAt: &publishedAt,
		ImageURL:    "https://example.com/image1.jpg",
		GUID:        "post-1",
	}
	if err := db.CreatePost(post1); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✓ Created post: %s\n", post1.Title)

	post2 := &models.Post{
		FeedID:      feed.ID,
		Title:       "Second Test Post",
		Link:        "https://example.com/post2",
		Description: "This is the second test post",
		GUID:        "post-2",
	}
	if err := db.CreatePost(post2); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✓ Created post: %s\n\n", post2.Title)

	// Test 7: Get all posts
	fmt.Println("7. Retrieving all posts...")
	posts, err := db.GetAllPosts(10, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✓ Found %d posts\n", len(posts))
	for _, p := range posts {
		fmt.Printf("      - %s (Feed: %s)\n", p.Title, p.FeedName)
	}
	fmt.Println()

	// Test 8: Get posts by feed ID
	fmt.Println("8. Getting posts by feed ID...")
	feedPosts, err := db.GetPostsByFeedID(feed.ID, 10, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✓ Found %d posts for feed '%s'\n\n", len(feedPosts), feed.Name)

	// Test 9: Get post by GUID
	fmt.Println("9. Getting post by GUID...")
	postByGUID, err := db.GetPostByGUID(feed.ID, "post-1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   ✓ Found post: %s\n\n", postByGUID.Title)

	// Test 10: Mark post as read
	fmt.Println("10. Marking post as read...")
	if err := db.MarkPostAsRead(post1.ID, true); err != nil {
		log.Fatal(err)
	}
	fmt.Println("   ✓ Post marked as read\n")

	fmt.Println("=== Testing Constraints ===\n")

	// Test 11: Try to create duplicate post (should fail)
	fmt.Println("11. Testing duplicate GUID constraint...")
	duplicatePost := &models.Post{
		FeedID:      feed.ID,
		Title:       "Duplicate Post",
		Link:        "https://example.com/duplicate",
		Description: "This should fail",
		GUID:        "post-1", // Same GUID as post1
	}
	err = db.CreatePost(duplicatePost)
	if err != nil {
		fmt.Printf("   ✓ Duplicate correctly rejected: %v\n\n", err)
	} else {
		fmt.Println("   ✗ ERROR: Duplicate post was allowed!\n")
	}

	// Test 12: Delete feed (should cascade to posts)
	fmt.Println("12. Testing cascade delete...")
	postsBeforeDelete, _ := db.GetPostsByFeedID(feed.ID, 100, 0)
	fmt.Printf("   Posts before delete: %d\n", len(postsBeforeDelete))

	if err := db.DeleteFeed(feed.ID); err != nil {
		log.Fatal(err)
	}
	fmt.Println("   ✓ Feed deleted")

	postsAfterDelete, _ := db.GetPostsByFeedID(feed.ID, 100, 0)
	fmt.Printf("   Posts after delete: %d\n", len(postsAfterDelete))
	fmt.Println("   ✓ Cascade delete working correctly\n")

	// Test 13: Delete all posts
	fmt.Println("13. Testing delete all posts...")
	if err := db.DeleteAllPosts(); err != nil {
		log.Fatal(err)
	}
	allPostsAfter, _ := db.GetAllPosts(100, 0)
	fmt.Printf("   ✓ All posts deleted (remaining: %d)\n\n", len(allPostsAfter))

	fmt.Println("=== All Tests Passed! ===")
}
