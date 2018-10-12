package internal

import (
	"encoding/json"
	"github.com/mitchellh/go-homedir"
	"os"
)

type cliConfig struct {
	GoogleApiKey string `json:"google_api_key"`
}

func openConfigFile() (*os.File, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(dir + "/.lyricsfinder")
	return file, err
}

func (config *cliConfig) SaveConfig() error {
	file, err := openConfigFile()
	defer file.Close()
	if err != nil {
		return err
	}

	return json.NewEncoder(file).Encode(config)
}

func GetConfig() (*cliConfig, error) {
	file, err := openConfigFile()
	defer file.Close()
	if err != nil {
		return nil, err
	}

	var config cliConfig
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
