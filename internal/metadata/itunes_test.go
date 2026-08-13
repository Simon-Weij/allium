package metadata

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Simon-Weij/allium/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testDownloader struct {
	mock.Mock
}

func (m *testDownloader) downloadAlbumCover(artworkUrl, target string) error {
	args := m.Called(artworkUrl, target)

	return args.Error(0)
}

func TestGetAlbumCover(t *testing.T) {
	t.Parallel()

	t.Run("should correctly return the album covers", func(t *testing.T) {
		t.Parallel()

		cfg := testutil.SetupTestingConfig(t)
		cfg.Data = t.TempDir()
		metadata := NewMetadata(cfg)

		downloader := &testDownloader{
			Mock: mock.Mock{},
		}
		metadata.downloader = downloader

		const id = "abcdef"

		coverPath := filepath.Join(cfg.Data, "covers", "ab", "cd", "ef", "cover.jpg")

		downloader.On("downloadAlbumCover", id, coverPath).Return(nil)

		path, err := metadata.GetAlbumCover(id)
		require.NoError(t, err)
		assert.Equal(t, coverPath, path)

		downloader.AssertExpectations(t)
	})
	t.Run("should not download when the cover already exists", func(t *testing.T) {
		t.Parallel()

		cfg := testutil.SetupTestingConfig(t)
		cfg.Data = t.TempDir()
		metadata := NewMetadata(cfg)

		downloader := &testDownloader{
			Mock: mock.Mock{},
		}
		metadata.downloader = downloader

		const id = "abcdef"

		coverPath := filepath.Join(cfg.Data, "covers", "ab", "cd", "ef", "cover.jpg")

		require.NoError(t, os.MkdirAll(filepath.Dir(coverPath), 0o755))
		require.NoError(t, os.WriteFile(coverPath, []byte("art"), 0o644))

		path, err := metadata.GetAlbumCover(id)
		require.NoError(t, err)
		assert.Equal(t, coverPath, path)

		downloader.AssertNotCalled(t, "downloadAlbumCover", mock.Anything, mock.Anything)
	})
}

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

func TestGetCoverArtURL(t *testing.T) {
	t.Parallel()

	t.Run("should return the URL as expected", func(t *testing.T) {
		t.Parallel()

		got, err := coverArtURL("https://example.com/100x100bb.jpg")
		require.NoError(t, err)
		assert.Equal(t, "https://example.com/"+coverArtSize, got)
	})
	t.Run("should error when the url isn't valid (doesn't contain a /)", func(t *testing.T) {
		t.Parallel()

		_, err := coverArtURL("notanurl")
		assert.ErrorIs(t, err, errInvalidArtworkURL)
	})
}
