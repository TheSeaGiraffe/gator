package main

import (
	"github.com/TheSeaGiraffe/gator/internal/database"
)

type State struct {
	DB     *database.Queries
	Config *Config
}
