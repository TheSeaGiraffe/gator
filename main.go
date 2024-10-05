package main

import (
	"fmt"
	"os"

	"github.com/TheSeaGiraffe/gator/internal/commands"
	"github.com/TheSeaGiraffe/gator/internal/config"
	"github.com/TheSeaGiraffe/gator/internal/state"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	st := state.State{
		Config: &cfg,
	}

	cmds := commands.InitCommands()
	cmds.Register("login", commands.HandlerLogin)

	userArgs := os.Args
	if len(userArgs) < 2 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}

	cmdName := userArgs[1]
	cmdArgs := []string{}
	if len(userArgs) >= 3 {
		cmdArgs = userArgs[2:]
	}

	loginCmd := commands.Command{
		Name: cmdName,
		Args: cmdArgs,
	}

	err = cmds.Run(&st, loginCmd)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("Successfully loaded config file. Current config:")
	fmt.Printf("\nDB URL: '%s'\n", cfg.DBUrl)
	fmt.Printf("Name of current user: '%s'\n", cfg.CurrentUserName)
}
