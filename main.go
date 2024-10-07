package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/TheSeaGiraffe/gator/internal/commands"
	"github.com/TheSeaGiraffe/gator/internal/config"
	"github.com/TheSeaGiraffe/gator/internal/database"
	"github.com/TheSeaGiraffe/gator/internal/state"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	// Maybe combine setting up the State struct into a single function later
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		fmt.Printf("Error connecting to database: %s", err.Error())
		os.Exit(1)
	}
	dbQueries := database.New(db)

	st := state.State{
		DB:     dbQueries,
		Config: cfg,
	}

	// Maybe combine the command registration logic into a single function later
	cmds := commands.InitCommands()
	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)
	cmds.Register("reset", commands.HandlerReset)

	// Maybe combine the logic for running commands into a single function
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

	cmd := commands.Command{
		Name: cmdName,
		Args: cmdArgs,
	}

	err = cmds.Run(&st, cmd)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}
