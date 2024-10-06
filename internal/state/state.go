package state

import (
	"github.com/TheSeaGiraffe/gator/internal/config"
	"github.com/TheSeaGiraffe/gator/internal/database"
)

type State struct {
	DB     *database.Queries
	Config *config.Config
}
