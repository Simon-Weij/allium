package metadata

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Simon-Weij/allium/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateDirs(t *testing.T) {
	t.Parallel()

	t.Run("should create expected folders", func(t *testing.T) {
		t.Parallel()
		testDir := t.TempDir()
		cfg := testutil.SetupTestingConfig(t)
		metadata := NewMetadata(cfg)

		path, err := metadata.createDirs(testDir, "abcdef")
		require.NoError(t, err)
		assert.Equal(t, filepath.Join("ab", "cd", "ef"), path)
	})
	t.Run("should error when insufficient permissions", func(t *testing.T) {
		t.Parallel()
		testDir := t.TempDir()
		err := os.Chmod(testDir, 0o444)
		require.NoError(t, err)

		cfg := testutil.SetupTestingConfig(t)
		metadata := NewMetadata(cfg)
		_, err = metadata.createDirs(testDir, "abcdef")
		assert.ErrorIs(t, err, errCreatingDirs)
	})
}
