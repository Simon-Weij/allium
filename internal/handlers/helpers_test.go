package handlers

import (
	"testing"

	"github.com/Simon-Weij/allium/internal/config"
)

func setupTestingConfig(t *testing.T) config.Config {
	t.Helper()

	cfg := config.Config{
		Name:     "allium",
		Version:  "0.1.0",
		Username: "alice",
		Password: "secret",
	}

	return cfg
}
