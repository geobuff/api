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
	Auth0 struct {
		Audience string `json:"audience"`
		Issuer   string `json:"issuer"`
	} `json:"auth0"`
	Cors struct {
		Origins []string `json:"origins"`
		Methods []string `json:"methods"`
	} `json:"cors"`
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