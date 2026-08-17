package handlers

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Simon-Weij/allium/internal/subsonic"
)

func (s Server) HandleGetCoverArt(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	id := query.Get("id")
	// TODO: support size

	coverPath, err := s.metadata.GetAlbumCover(id)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	imageBytes, err := os.ReadFile(coverPath)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "image/png")

	_, _ = w.Write(imageBytes)
}

func (s Server) HandleStream(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()

	id := query.Get("id")
	if id == "" {
		http.Error(w, "id is required", subsonic.ErrParameterMissing)
	}

	res, err := s.metadata.GetSongById(id)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		slog.Error("couldn't get song by id", "error", err)

		return
	}

	song := res.Results[0]

	path, err := s.metadata.DownloadOrGetSong(ctx, song.ArtistName, song.TrackName)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		slog.Error("couldn't download song", "error", err)

		return
	}

	http.ServeFile(w, r, path)
}
