package resolver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/lrstanley/go-ytdlp"
)

const (
	folderPermission = 0o755
)

func (r Resolver) DownloadOrGetSong(ctx context.Context, artist, title string) (string, error) {
	outputDir := filepath.Join(r.cfg.Data, "Music", artist, title)
	if err := os.MkdirAll(outputDir, folderPermission); err != nil {
		return "", fmt.Errorf("could not create folder: %w", err)
	}

	path := filepath.Join(outputDir, "song.mp3")

	_, err := os.Stat(path)
	if err == nil {
		slog.Info("file found at " + path)

		return path, nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("could not stat %s: %w", path, err)
	}

	dlp := ytdlp.New().ExtractAudio().AudioFormat("mp3").Output(path)

	if _, err := dlp.Run(ctx, fmt.Sprintf("ytsearch1: %s - %s", artist, title)); err != nil {
		slog.Error("could not download "+artist+" by "+title, "error", err)

		return "", fmt.Errorf("could not download %s: %w", title, err)
	}

	return path, nil
}
