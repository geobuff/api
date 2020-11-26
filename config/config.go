package config

import (
	"encoding/json"
	"os"
)

// Values is the global variable for our configuration.
var Values *Config

// Config is an object that matches the entries in our configuration file.
type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
	} `json:"database"`
}

// Load reads the configuration from the config file and loads into out public Config variable.
func Load(fileName string) {
	configFile, err := os.Open(fileName)
	defer configFile.Close()
	if err != nil {
		panic(err)
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&Values)
	if err != nil {
		panic(err)
	}
}
