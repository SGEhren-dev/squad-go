package configuration

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/log"
)

type Config struct {
	LogFilePath string `json:"logFilePath"`
	Database    struct {
		Dialect  string `json:"dialect"`
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
		Port     string `json:"port"`
		Name     string `json:"name"`
	} `json:"database"`
	Rcon struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Password string `json:"password"`
	} `json:"rcon"`
	Connectors map[string]json.RawMessage `json:"connectors"`
	Plugins    map[string]struct {
		Enabled  bool            `json:"enabled"`
		Settings json.RawMessage `json:"settings"`
	} `json:"plugins"`
}

func (c *Config) LoadConfig() error {
	configHandle, err := os.Open("config.json")

	if err != nil {
		log.Error("Error opening config file: " + err.Error())

		return err
	}

	defer configHandle.Close()

	jsonParser := json.NewDecoder(configHandle)
	jsonParser.Decode(&c)

	return err
}
