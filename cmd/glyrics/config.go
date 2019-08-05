package main

import (
	"encoding/json"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
)

// CliConfig represents the config used by the command line tool
type CliConfig struct {
	GoogleApiKey string `json:"google_api_key"`
}

func openConfigFile() (*os.File, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path.Join(dir, ".glyrics"), os.O_RDWR|os.O_CREATE, 0755)
	return file, err
}

// SaveConfig saves the config to the config file
func (config CliConfig) SaveConfig() error {
	file, err := openConfigFile()
	defer func() { _ = file.Close() }()
	if err != nil {
		return err
	}

	return json.NewEncoder(file).Encode(config)
}

// GetConfig loads the CliConfig from the config file
func GetConfig() (*CliConfig, error) {
	file, err := openConfigFile()
	defer func() { _ = file.Close() }()
	if err != nil {
		return nil, err
	}

	var config CliConfig
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
