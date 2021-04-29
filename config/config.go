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
	Auth struct {
		SigningKey string `json:"signingKey"`
	} `json:"auth"`
	Cors struct {
		Origins []string `json:"origins"`
		Methods []string `json:"methods"`
	} `json:"cors"`
	SiteUrl string `json:"siteUrl"`
}

// Load reads the configuration from the config file and loads into our public Values variable.
var Load = func(fileName string) error {
	config, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer config.Close()

	parser := json.NewDecoder(config)
	return parser.Decode(&Values)
}
