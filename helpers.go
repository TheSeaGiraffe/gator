package main

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/TheSeaGiraffe/gator/internal/database"
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

	// Save all posts in feed to database
	for _, item := range rssFeed.Channel.Item {
		// Parse `PublishedAt` time string
		publishedAtTime, err := parsePublishTime(item.PubDate)
		if err != nil {
			switch {
			case err.Error() == "No suitable matches":
				// Use the current time for now; will think of better solution later
				publishedAtTime = time.Now()
			default:
				return err
			}
		}

		// Create `sql.NullString` object
		descString := sql.NullString{
			String: item.Description,
			Valid:  item.Description != "",
		}

		// Save post to DB
		newPost := database.CreatePostParams{
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: descString,
			PublishedAt: publishedAtTime,
			FeedID:      feed.ID,
		}

		_, err = s.DB.CreatePost(context.Background(), newPost)
		if (err != nil) && (err.Error() != `pq: duplicate key value violates unique constraint "posts_url_key"`) {
			return err
		}
	}

	return nil
}

func parsePublishTime(timeStr string) (time.Time, error) {
	// Get the relevant parts of the time string using regexes
	matches, err := getTimeStrParts(timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("Error extracting parts of time string: %w", err)
	}
	if len(matches) == 0 {
		return time.Time{}, fmt.Errorf("No suitable matches")
	}

	// Construct custom time format string and attempt to parse timeStr
	var formatStrBuilder strings.Builder
	formatStrBuilder.WriteString("Mon, 02 Jan 2006")
	if len(strings.Split(matches[0], ":")) == 3 {
		formatStrBuilder.WriteString(" 15:04:05 ")
	} else {
		formatStrBuilder.WriteString(" 15:04 ")
	}
	formatStrBuilder.WriteString(matches[1])

	parsedTime, err := time.Parse(formatStrBuilder.String(), timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("Could not parse time string: %w", err)
	}
	return parsedTime, nil
}

func getTimeStrParts(timeStr string) ([]string, error) {
	pattern := `(\d{2}:\d{2}(?::\d{2})?).*(\b[A-Z]+|\+\d{4})$`
	r, err := regexp.Compile(pattern)
	if err != nil {
		return []string{}, fmt.Errorf("Error compiling regex: %w", err)
	}
	match := r.FindStringSubmatch(timeStr)
	if len(match) == 0 {
		return []string{}, nil
	}
	return match[1:], nil
}
