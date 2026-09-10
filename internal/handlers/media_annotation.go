package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Simon-Weij/allium/generated/sqlc"
	"github.com/Simon-Weij/allium/internal/subsonic"
)

func (s Server) HandleScrobble(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queries := r.URL.Query()

	id := queries.Get("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		slog.Error("could not parse song id", "error", err)
		subsonic.WriteError(w, http.StatusInternalServerError, s.cfg, subsonic.ErrGeneric, "internal server error")

		return
	}

	username := queries.Get("u")

	itunesSongs, err := s.iTunesClient.GetSongById(id)
	if err != nil {
		slog.Error("could not get song by id", "error", err)
		subsonic.WriteError(w, http.StatusInternalServerError, s.cfg, subsonic.ErrGeneric, "internal server error")

		return
	}

	if len(itunesSongs.Results) == 0 {
		slog.Info("no results found", "id", id)
		subsonic.WriteNotFound(w, s.cfg, "couldn't find song with id: "+id)

		return
	}

	itunesSong := itunesSongs.Results[0]

	song := convertItunesSong(itunesSong)

	albumID, err := strconv.Atoi(song.AlbumId)
	if err != nil {
		slog.Error("could not parse album id", "error", err)
		subsonic.WriteError(w, http.StatusInternalServerError, s.cfg, subsonic.ErrGeneric, "internal server error")

		return
	}

	artistID, err := strconv.Atoi(song.ArtistId)
	if err != nil {
		slog.Error("could not parse artist id", "error", err)
		subsonic.WriteError(w, http.StatusInternalServerError, s.cfg, subsonic.ErrGeneric, "internal server error")

		return
	}

	if err := s.queries.UpdatePlays(ctx, sqlc.UpdatePlaysParams{
		ID:         int64(idInt),
		Title:      song.Title,
		Artist:     song.Artist,
		Album:      song.Album,
		AlbumID:    int64(albumID),
		ArtistID:   int64(artistID),
		Duration:   int64(song.Duration),
		Track:      int64(song.Track),
		Year:       int64(song.Year),
		Genre:      song.Genre,
		DiscNumber: int64(song.DiscNumber),
		User:       username,
		ArtworkUrl: song.CoverArt,
	}); err != nil {
		slog.Error("could not update plays", "error", err)
		subsonic.WriteError(w, http.StatusInternalServerError, s.cfg, subsonic.ErrGeneric, "internal server error")

		return
	}

	response := subsonic.NewEmptyResponse(s.cfg)
	subsonic.WriteJSON(w, http.StatusOK, response)
}
