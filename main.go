package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/TheSeaGiraffe/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := ReadConfig()
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

	st := State{
		DB:     dbQueries,
		Config: cfg,
	}

	cmds := NewCommands()

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

	cmd := Command{
		Name: cmdName,
		Args: cmdArgs,
	}

	err = cmds.Run(&st, cmd)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}
