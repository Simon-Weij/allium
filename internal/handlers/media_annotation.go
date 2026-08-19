package handlers

import (
	"net/http"
	"strconv"

	"github.com/Simon-Weij/allium/generated/sqlc"
)

func (s Server) HandleScrobble(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queries := r.URL.Query()

	id := queries.Get("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	username := queries.Get("u")

	itunesSongs, err := s.iTunesclient.GetSongById(id)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	itunesSong := itunesSongs.Results[0]

	song := convertItunesSong(itunesSong)

	if err := s.queries.UpdatePlays(ctx, sqlc.UpdatePlaysParams{
		ID:     int64(idInt),
		Title:  song.Title,
		Artist: song.Artist,
		User:   username,
	}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}
}
