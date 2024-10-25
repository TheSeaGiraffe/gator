package main

import (
	"context"

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
			return err
		}
		return handler(s, cmd, user)
	}
}
