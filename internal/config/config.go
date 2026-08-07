package config

import "os"

type Config struct {
	Version string
	Name    string
	ApiKey  string
}

func NewConfig() Config {
	apiKey := os.Getenv("API_KEY")
	if string(apiKey) == "" {
		panic("API_KEY is not set")
	}
	return Config{
		Name:    "allium",
		Version: "0.1.0",
		ApiKey:  apiKey,
	}
}
