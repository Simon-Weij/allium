package handlers

import (
	"testing"

	"github.com/Simon-Weij/allium/internal/config"
)

func setupTestingConfig(t *testing.T) config.Config {
	t.Helper()

	cfg := config.Config{
		Username: "alice",
		Password: "secret",
	}

	return cfg
}
