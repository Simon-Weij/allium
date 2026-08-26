package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Simon-Weij/allium/internal/subsonic"
)

func (s Server) HandleGetAlbum(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "no id provided", http.StatusNotFound)

		return
	}

	res, err := s.iTunesClient.GetAlbumMetadata(id)
	if err != nil {
		slog.Error("something went wrong trying to fetch metadata", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	if len(res.Results) == 0 {
		slog.Info("no results found", "id", id)
		http.Error(w, "couldn't find album with id: "+id, http.StatusNotFound)

		return
	}

	album := convertItunesAlbum(res)

	if album.SongCount == 0 {
		slog.Info("no songs found for album", "id", id)
		http.Error(w, "couldn't find album with id: "+id, http.StatusNotFound)

		return
	}

	response := subsonic.NewEmptyResponse(s.cfg)
	response.SubsonicResponse.Album = &album
	subsonic.WriteJSON(w, http.StatusOK, response)
}

func (s Server) HandleGetArtist(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "no id provided", http.StatusNotFound)

		return
	}

	res, err := s.iTunesClient.GetArtistById(id)
	if err != nil {
		slog.Error("something went wrong trying to fetch metadata", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	artist := convertItunesArtist(res)

	response := subsonic.NewEmptyResponse(s.cfg)
	response.SubsonicResponse.Artist = &artist
	subsonic.WriteJSON(w, http.StatusOK, response)
}
