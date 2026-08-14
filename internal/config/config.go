package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Data     string `toml:"data"` // The path where the user data would be stored
}

const (
	Version = "0.0.1"
	Name    = "Allium"
)

var ErrConfigDirNotFound = errors.New("config dir not found")

// ParseConfig parses the configuration data stored in configDir/allium
func ParseConfig() (*Config, error) {
	var cfg Config

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrConfigDirNotFound, err)
		}

		_, err = os.Stat(filepath.Join(configDir, "allium", "config.toml"))
		if err != nil {
			panic("could not find config file in ~/.config/allium or the location of $CONFIG_PATH")
		}
	}

	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("could not decode file: %w", err)
	}

	if cfg.Data == "" {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return nil, fmt.Errorf("could not get user cache dir: %w", err)
		}

		cfg.Data = filepath.Join(cacheDir, Name)
	}

	requireNotEmpty(cfg.Username)

	return &cfg, nil
}

// requireNotEmpty panics if v is empty
func requireNotEmpty(v string) {
	if v == "" {
		panic(v + " is empty")
	}
}
