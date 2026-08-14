package handlers

import (
	"log/slog"
	"net/http"
)

func (s Server) HandleGetAlbum(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	res, err := s.metadata.GetAlbumMetadata(id)
	if err != nil {
		slog.Error("something went wrong trying to fetch metadata", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	album := ConvertItunesAlbum(res)

	response := NewEmptyResponse(s.cfg)
	response.SubsonicResponse.Album = &album
	WriteJSON(w, http.StatusOK, response)
}
