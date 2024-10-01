package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

// Config contains the configuration settings for the gator CLI
type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if homeDir == "" {
		return "", fmt.Errorf("Could not find current user's home directory.")
	}

	configFilePath := filepath.Join(homeDir, configFileName)

	return configFilePath, nil
}

// Read reads the `gatorconfig.json` file in order to set certain config values
func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("Error getting path to config file: %w", err)
	}

	configFile, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("Error reading config file: %w", err)
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return Config{}, fmt.Errorf("Error unmarshaling config file: %w", err)
	}

	return config, nil
}

// SetUser writes the config file to the `gatorconfig.json` file, the default location of which is in
// the home directory.
func (cfg *Config) SetUser(current_user string) error {
	// Assign `current_user` to the `CurrentUserName` field
	cfg.CurrentUserName = current_user

	// Write the current config to the ".gatorconfig.json" file
	configContents, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("Error marshaling configuration: %w", err)
	}

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("Error getting the path to the config file: %w", err)
	}

	err = os.WriteFile(configFilePath, configContents, 0666)
	if err != nil {
		return fmt.Errorf("Error writing config to file: %w", err)
	}
	return nil
}
