package config

import (
	"errors"
	"os"
)

type Config struct {
	Username     string
	Password     string
	Data         string
	DatabasePath string
}

const (
	Version = "0.0.1"
	Name    = "Allium"

	DataDir      = "/data"
	DatabasePath = "/data/data.db"
)

var ErrCredentialsMissing = errors.New("ALLIUM_USERNAME and ALLIUM_PASSWORD must be set")

func ParseConfig() (*Config, error) {
	username := os.Getenv("ALLIUM_USERNAME")
	password := os.Getenv("ALLIUM_PASSWORD")

	if username == "" || password == "" {
		return nil, ErrCredentialsMissing
	}

	return &Config{
		Username:     username,
		Password:     password,
		Data:         DataDir,
		DatabasePath: DatabasePath,
	}, nil
}
