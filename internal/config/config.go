package config

import "os"

type Config struct {
	Version string
	Name    string
	ApiKey  string
}

func NewConfig() Config {
	apiKey := os.Getenv("API_KEY")
	return Config{
		Name:    "allium",
		Version: "0.1.0",
		ApiKey:  apiKey,
	}
}
