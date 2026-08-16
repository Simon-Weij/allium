package testutil

import (
	"testing"

	"github.com/Simon-Weij/allium/internal/config"
)

func SetupTestingConfig(t *testing.T) config.Config {
	t.Helper()

	cfg := config.Config{
		Username:     "alice",
		Password:     "secret",
		Data:         "",
		DatabasePath: "",
	}

	return cfg
}
