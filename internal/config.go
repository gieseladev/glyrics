package internal

import (
	"encoding/json"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
)

type CliConfig struct {
	GoogleApiKey string `json:"google_api_key"`
}

func openConfigFile() (*os.File, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path.Join(dir, ".lyricsfinder"), os.O_RDWR|os.O_CREATE, 0755)
	return file, err
}

func (config CliConfig) SaveConfig() error {
	file, err := openConfigFile()
	defer file.Close()
	if err != nil {
		return err
	}

	return json.NewEncoder(file).Encode(config)
}

func GetConfig() (*CliConfig, error) {
	file, err := openConfigFile()
	defer file.Close()
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
