package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Username string
	Password string
}

const (
	Version = "0.0.1"
	Name    = "Allium"
)

func ParseConfig() (*Config, error) {
	var cfg Config
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}

		_, err = os.Stat(filepath.Join(configDir, "allium", "config.toml"))
		if err != nil {
			panic("could not find config file in ~/.config/allium or the location of $CONFIG_PATH")
		}

	}
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("could not decode file: %w", err)
	}

	requireNotEmpty(cfg.Username)

	return &cfg, nil
}

func requireNotEmpty(v string) {
	if v == "" {
		panic(v + " is empty")
	}
}
