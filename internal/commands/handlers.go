package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/TheSeaGiraffe/gator/internal/database"
	"github.com/TheSeaGiraffe/gator/internal/rss"
	"github.com/TheSeaGiraffe/gator/internal/state"
	"github.com/google/uuid"
)

// HandlerLogin is a handler for the `login` subcommand. `login` is used to set the current user
// to the specified user.
func HandlerLogin(s *state.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("No username provided")
	} else if len(cmd.Args) > 1 {
		return fmt.Errorf("Username must be a single string")
	}

	user, err := s.DB.GetUser(context.Background(), cmd.Args[0])
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

	fmt.Printf("Username has been set to '%s'\n", s.Config.CurrentUserName)

	return nil
}

// HandlerRegister is a handler for the `register` subcommand. `register` adds the current user
// to the database.
func HandlerRegister(s *state.State, cmd Command) error {
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
	fmt.Printf("User info:\n\n")
	fmt.Printf("Name: %s\n", user.Name)
	fmt.Printf("ID: %s\n", user.ID)
	fmt.Printf("Created At: %s\n", user.CreatedAt.String())
	fmt.Printf("Updated At: %s\n", user.UpdatedAt.String())

	return nil
}

func HandlerReset(s *state.State, cmd Command) error {
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
	fmt.Println("All users deleted from database and the previous user has been logged out.")

	return nil
}

func HandlerUsers(s *state.State, cmd Command) error {
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

func HandlerAgg(s *state.State, cmd Command) error {
	// Ignore user args for now

	// Get RSS feed at `https://www.wagslane.dev/index.xml`
	rssFeed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	// Print the entire RSSFeed struct to the console
	// fmt.Printf("%+v\n", rssFeed)
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

func HandlerAddFeed(s *state.State, cmd Command) error {
	// Validate user args. Do note that one of the args is a url
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Missing arguments. `addfeed` takes the name of the RSS feed and its URL.")
	} else if len(cmd.Args) > 2 {
		return fmt.Errorf("Too many arguments. `addfeed` takes the name of the RSS feed and its URL.")
	}

	_, err := url.ParseRequestURI(cmd.Args[1])
	if err != nil {
		return fmt.Errorf("Invalid URL")
	}

	// Link feed to current user and add it to feeds table
	user, err := s.DB.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return err
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

	fmt.Printf("\nRSS Feed ID: %d\n", rssFeed.ID)
	fmt.Printf("RSS Feed created at: %s\n", rssFeed.CreatedAt.String())
	fmt.Printf("RSS Feed updated at: %s\n", rssFeed.UpdatedAt.String())
	fmt.Printf("RSS Feed name: %s\n", rssFeed.Name)
	fmt.Printf("RSS Feed URL: %s\n", rssFeed.Url)
	fmt.Printf("RSS Feed User ID: %v\n", rssFeed.UserID)

	// Print a sample of items from the current feed
	return nil
}
