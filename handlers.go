package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/TheSeaGiraffe/gator/internal/database"
	"github.com/google/uuid"
)

var defaultAggInterval = time.Minute * 5

// HandlerLogin is a handler for the `login` subcommand. `login` is used to set the current user
// to the specified user.
func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("No username provided")
	} else if len(cmd.Args) > 1 {
		return fmt.Errorf("Username must be a single string")
	}

	user, err := s.DB.GetUserByName(context.Background(), cmd.Args[0])
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return fmt.Errorf("User '%s' does not exist", cmd.Args[0])
		default:
			return err
		}
	}

	err = s.Config.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("%s is now logged in.\n", s.Config.CurrentUserName)

	return nil
}

// HandlerRegister is a handler for the `register` subcommand. `register` adds the current user
// to the database.
func HandlerRegister(s *State, cmd Command) error {
	// Check that a username was passed in the args
	if len(cmd.Args) == 0 {
		return fmt.Errorf("No username provided")
	} else if len(cmd.Args) > 1 {
		return fmt.Errorf("Username must be a single string")
	}

	// Create a new user in the database
	userData := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}

	user, err := s.DB.CreateUser(context.Background(), userData)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_name_key"`:
			return fmt.Errorf("User '%s' already exists", cmd.Args[0])
		default:
			return err
		}
	}

	// Set the current user in the config to the given name
	err = s.Config.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User '%s' successfully created\n\n", user.Name)
	fmt.Printf("Name: %s\n", user.Name)
	fmt.Printf("ID: %s\n", user.ID)
	fmt.Printf("Created At: %s\n", user.CreatedAt.String())
	fmt.Printf("Updated At: %s\n", user.UpdatedAt.String())

	return nil
}

func HandlerReset(s *State, cmd Command) error {
	// Validate args
	if len(cmd.Args) > 0 {
		return fmt.Errorf("Command does not take any arguments")
	}

	// Delete all users in DB
	err := s.DB.DeleteUsers(context.Background())
	if err != nil {
		return err
	}

	// Remove `current_user_name` field from `~/.gatorconfig.json`
	err = s.Config.SetUser("")
	if err != nil {
		return err
	}

	// Print message to console for logging purposes
	fmt.Println("Database has been reset and the previous user has been logged out.")

	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	// Validate user args
	if len(cmd.Args) > 0 {
		return fmt.Errorf("Command does not take any arguments")
	}

	// Get users from DB. Don't forget to validate slice.
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Println("Database currently does not contain any users.")
		return nil
	}

	// Print users
	var userName string
	for _, user := range users {
		userName = fmt.Sprintf("* %s", user.Name)
		if user.Name == s.Config.CurrentUserName {
			userName = fmt.Sprintf("%s (current)", userName)
		}
		fmt.Println(userName)
	}

	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	// Validate user args
	var tickInterval time.Duration
	var err error
	if len(cmd.Args) > 1 {
		return fmt.Errorf("Too may arguments. `agg` takes a time duration string.")
	} else if len(cmd.Args) == 0 {
		tickInterval = defaultAggInterval
	} else {
		tickInterval, err = time.ParseDuration(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("Could not parse time duration string.")
		}
	}

	// Get new feed after the specified tickInterval
	// Since we're using `agg` to print feeds I think we should just log to a file in the event of an error
	// For now I'll just break out of the loop.
	ticker := time.NewTicker(tickInterval)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			ticker.Stop()
			return fmt.Errorf("Error fetching feed: %w", err)
		}
	}
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Missing arguments. `addfeed` takes the name of the RSS feed and its URL.")
	} else if len(cmd.Args) > 2 {
		return fmt.Errorf("Too many arguments. `addfeed` takes the name of the RSS feed and its URL.")
	}

	_, err := url.ParseRequestURI(cmd.Args[1])
	if err != nil {
		return fmt.Errorf("Invalid URL")
	}

	rssFeedParams := database.CreateFeedParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}
	rssFeed, err := s.DB.CreateFeed(context.Background(), rssFeedParams)
	if err != nil {
		return fmt.Errorf("Error saving feed: %w", err)
	}

	feedFollowEntry := database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    rssFeed.ID,
	}
	_, err = s.DB.CreateFeedFollow(context.Background(), feedFollowEntry)
	if err != nil {
		return fmt.Errorf("Error creating feed-follow entry: %w", err)
	}

	fmt.Printf("Feed %q successfully added.\n", rssFeed.Name)
	fmt.Printf("\nRSS Feed ID: %d\n", rssFeed.ID)
	fmt.Printf("RSS Feed created at: %s\n", rssFeed.CreatedAt.String())
	fmt.Printf("RSS Feed updated at: %s\n", rssFeed.UpdatedAt.String())
	fmt.Printf("RSS Feed name: %s\n", rssFeed.Name)
	fmt.Printf("RSS Feed URL: %s\n", rssFeed.Url)
	fmt.Printf("RSS Feed User ID: %v\n", rssFeed.UserID)

	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("Command does not take any arguments")
	}

	feeds, err := s.DB.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error retrieving feeds: %w", err)
	}

	for _, feed := range feeds {
		user, err := s.DB.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			// Just skip the current iteration for now
			// Might add better handling later
			continue
		}

		fmt.Printf("Feed name: %s\n", feed.Name)
		fmt.Printf("Feed URL: %s\n", feed.Url)
		fmt.Printf("Feed owner: %s\n\n", user.Name)
	}

	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	// Validate user input
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Command expects a URL to an RSS Feed")
	} else if len(cmd.Args) > 1 {
		return fmt.Errorf("Command expects only a single argument")
	}

	_, err := url.ParseRequestURI(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Invalid URL")
	}

	// Check that feed exists
	feed, err := s.DB.GetFeedsByURL(context.Background(), cmd.Args[0])
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return fmt.Errorf("Feed does not exist. Add it with the 'addfeed' command")
		default:
			return fmt.Errorf("Error retrieving feeds: %w", err)
		}
	}

	// Create a new feed follow record and print the results
	feedFollowEntry := database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollow, err := s.DB.CreateFeedFollow(context.Background(), feedFollowEntry)
	if err != nil {
		return err
	}

	fmt.Printf("%s is now following '%s'\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	// Validate user input. Make sure that command doesn't take any input.
	if len(cmd.Args) > 0 {
		return fmt.Errorf("Command does not take any arguments")
	}

	feedFollows, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			fmt.Println("You are not following any feeds. Add some with the 'addfeed' command.")
			return nil
		default:
			return fmt.Errorf("Could not retrieve feeds for current user: %w", err)
		}
	}

	if len(feedFollows) == 0 {
		fmt.Println("You aren't currently following anything. Use the 'follow' command to follow a feed.")
	} else {
		fmt.Printf("Currently following:\n\n")
		for _, feed := range feedFollows {
			fmt.Printf(" - %s\n", feed.FeedName)
		}
	}

	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	// Validate user input. Make sure that command only takes a feed's URL
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Command expects the feed URL.")
	} else if len(cmd.Args) > 1 {
		return fmt.Errorf("Too many arguments. Make sure you are only passing the feed URL.")
	}

	// Get feed from URL
	feed, err := s.DB.GetFeedsByURL(context.Background(), cmd.Args[0])
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return fmt.Errorf("Feed does not exist.")
		default:
			return err
		}
	}

	// Delete feed-follow record for user
	feedFollowRowDel := database.DeleteFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.DB.DeleteFeed(context.Background(), feedFollowRowDel)
	if err != nil {
		return err
	}

	return nil
}

func HandlerBrowse(s *State, cmd Command, user database.User) error {
	// Validate user input. Takes an optional "limit" parameter with a default of 2
	var err error
	postLimit := 2
	if len(cmd.Args) > 1 {
		return fmt.Errorf(`Too many arguments. You may choose to add the maximum number of posts to display 
            as an integer. Defaults to 2.`)
	} else if len(cmd.Args) == 1 {
		postLimit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("Error parsing post limit string: %w", err)
		}
	}

	// Get posts for the current user and display them
	postParams := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(postLimit),
	}
	userPosts, err := s.DB.GetPostsForUser(context.Background(), postParams)
	if err != nil {
		return fmt.Errorf("Error retrieving posts: %w", err)
	}

	nUserPosts := len(userPosts)
	for i, post := range userPosts {
		fmt.Printf("\nTitle: %s\n", post.Title)
		fmt.Printf("\nURL: %s\n", post.Url)
		fmt.Printf("\nPublish Date: %s\n", post.PublishedAt.String())
		if post.Description.Valid {
			fmt.Printf("\nDescription:\n%s\n", post.Description.String)
		} else {
			fmt.Println("\nDescription: N/A")
		}

		if (nUserPosts > 1) && (i < nUserPosts-1) {
			fmt.Println("\n============")
		}
	}

	return nil
}
