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

func NewCommands() Commands {
	cmds := Commands{
		List: make(map[string]func(*state.State, Command) error),
	}
	cmds.Register("login", HandlerLogin)
	cmds.Register("register", HandlerRegister)
	cmds.Register("reset", HandlerReset)
	cmds.Register("users", HandlerUsers)
	cmds.Register("agg", HandlerAgg)
	cmds.Register("addfeed", HandlerAddFeed)
	cmds.Register("feeds", HandlerFeeds)
	cmds.Register("follow", HandlerFollow)
	cmds.Register("following", HandlerFollowing)

	return cmds
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
