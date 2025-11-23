package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getPathToConfig() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, configFileName), nil
}

func Read() (Config, error) {

	config := Config{}

	pathToConfig, err := getPathToConfig()
	if err != nil {
		return config, err
	}

	dat, err := os.ReadFile(pathToConfig)
	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(dat, &config); err != nil {
		return config, nil
	}

	return config, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username

	pathToConfig, err := getPathToConfig()
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}

	fi, err := os.Open(pathToConfig)
	if err != nil {
		return err
	}

	defer fi.Close()

	fi.Write(jsonData)
	return nil
}
