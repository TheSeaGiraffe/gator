package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/TheSeaGiraffe/gator/internal/database"
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
