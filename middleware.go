package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/TheSeaGiraffe/gator/internal/database"
)

type (
	CmdHandler     func(*State, Command) error
	CmdHandlerAuth func(s *State, cmd Command, user database.User) error
)

// Make sure that we're passing in the information of a user that's already logged in
func middlewareLoggedIn(handler CmdHandlerAuth) CmdHandler {
	return func(s *State, cmd Command) error {
		user, err := s.DB.GetUserByName(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return fmt.Errorf("User does not exist. Make sure that you've registered an account and are logged in.")
			default:
				return err
			}
		}
		return handler(s, cmd, user)
	}
}
