package commands

import (
	"fmt"

	"github.com/TheSeaGiraffe/gator/internal/state"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	List map[string]func(*state.State, Command) error
}

func InitCommands() Commands {
	return Commands{
		List: make(map[string]func(*state.State, Command) error),
	}
}

// Register registers a new handler function for a command name
func (c *Commands) Register(name string, f func(*state.State, Command) error) {
	_, ok := c.List[name]
	if !ok {
		c.List[name] = f
	}
}

// Run executes a given command with the provided state if it exists
func (c *Commands) Run(s *state.State, cmd Command) error {
	cmdL, ok := c.List[cmd.Name]
	if !ok {
		return fmt.Errorf("Command '%s' does not exist", cmd.Name)
	}

	err := cmdL(s, cmd)
	if err != nil {
		return fmt.Errorf("Error running command '%s': %w", cmd.Name, err)
	}
	return nil
}
