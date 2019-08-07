package main

import (
	"encoding/json"
	"os"
)

// CliConfig represents the config used by the command line tool
type CliConfig struct {
	GoogleApiKey string `json:"google_api_key"`
}

// LoadConfig loads the CliConfig from the config file
func LoadConfig(name string) (*CliConfig, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	defer func() { _ = file.Close() }()

	var config CliConfig
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves a config to the given location.
func SaveConfig(name string, config *CliConfig) error {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	return json.NewEncoder(file).Encode(config)
}
