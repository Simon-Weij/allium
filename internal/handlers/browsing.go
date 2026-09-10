package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Simon-Weij/allium/internal/subsonic"
)

func (s Server) HandleGetAlbum(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		subsonic.WriteNotFound(w, s.cfg, "no id provided")
		return
	}

	res, err := s.iTunesClient.GetAlbumMetadata(id)
	if err != nil {
		slog.Error("something went wrong trying to fetch metadata", "error", err)
		subsonic.WriteError(w, http.StatusInternalServerError, s.cfg, subsonic.ErrGeneric, "internal server error")

		return
	}

	if len(res.Results) == 0 {
		slog.Info("no results found", "id", id)
		subsonic.WriteNotFound(w, s.cfg, "couldn't find album with id: "+id)

		return
	}

	album := convertItunesAlbum(res)

	if album.SongCount == 0 {
		slog.Info("no songs found for album", "id", id)
		subsonic.WriteNotFound(w, s.cfg, "couldn't find album with id: "+id)

		return
	}

	response := subsonic.NewEmptyResponse(s.cfg)
	response.SubsonicResponse.Album = &album
	subsonic.WriteJSON(w, http.StatusOK, response)
}

func (s Server) HandleGetArtist(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		subsonic.WriteNotFound(w, s.cfg, "no id provided")
		return
	}

	res, err := s.iTunesClient.GetArtistById(id)
	if err != nil {
		slog.Error("something went wrong trying to fetch metadata", "error", err)
		subsonic.WriteError(w, http.StatusInternalServerError, s.cfg, subsonic.ErrGeneric, "internal server error")

		return
	}

	artist := convertItunesArtist(res)

	response := subsonic.NewEmptyResponse(s.cfg)
	response.SubsonicResponse.Artist = &artist
	subsonic.WriteJSON(w, http.StatusOK, response)
}
