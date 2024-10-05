package commands

import (
	"fmt"

	"github.com/TheSeaGiraffe/gator/internal/state"
)

func HandlerLogin(s *state.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("No username provided")
	} else if len(cmd.Args) > 1 {
		return fmt.Errorf("Username must be a single string")
	}

	err := s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Username has been set to '%s'\n", s.Config.CurrentUserName)
	return nil
}
