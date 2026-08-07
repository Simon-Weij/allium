package config

import "os"

type Config struct {
	Version  string
	Name     string
	Username string
	Password string
}

func NewConfig() Config {
	username := os.Getenv("ALLIUM_USERNAME")
	password := os.Getenv("ALLIUM_PASSWORD")
	if username == "" || password == "" {
		panic("ALLIUM_USERNAME and ALLIUM_PASSWORD must be set")
	}
	return Config{
		Name:     "allium",
		Version:  "0.1.0",
		Username: username,
		Password: password,
	}
}
