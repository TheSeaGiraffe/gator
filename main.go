package main

import (
	"fmt"
	"log"

	"github.com/TheSeaGiraffe/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	err = cfg.SetUser("fahmi")
	if err != nil {
		log.Fatalf("Error setting user: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	fmt.Println("Successfully loaded config file. Current config:")
	fmt.Printf("\nDB URL: '%s'\n", cfg.DBUrl)
	fmt.Printf("Name of current user: '%s'\n", cfg.CurrentUserName)
}
