package main

import (
	"context"
	"fmt"

	"github.com/TheSeaGiraffe/gator/internal/rss"
)

func scrapeFeeds(s *State) error {
	// Get next feed to fetch from DB and mark it as fetched
	feed, err := s.DB.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting next feed: %w", err)
	}

	err = s.DB.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("Error marking feed as fetched: %w", err)
	}

	// Fetch feed using URL
	rssFeed, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("Error fetching feed from URL: %w", err)
	}

	// Iterate over items in feed
	fmt.Println("RSS Feed")
	fmt.Println("========")
	fmt.Printf("\nChannel title: %s\n", rssFeed.Channel.Title)
	fmt.Printf("\nChannel link: %s\n", rssFeed.Channel.Link)
	fmt.Printf("\nChannel description: %s\n\n", rssFeed.Channel.Description)

	for i, item := range rssFeed.Channel.Item {
		fmt.Printf("Item %d\n", i+1)
		fmt.Printf("-------\n\n")
		fmt.Printf("Title: %s\n", item.Title)
		fmt.Printf("Link: %s\n", item.Link)
		fmt.Printf("Description: %s\n\n", item.Description)
	}

	return nil
}
