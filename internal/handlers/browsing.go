package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Simon-Weij/allium/internal/metadata"
)

//go:generate mockgen -source=browsing.go -destination=mocks/browsing_mock.go -package=mocks
type iTunesClient interface {
	GetAlbumMetadata(albumId string) (*metadata.ITunesResponse, error)
	GetArtistById(artistId string) (*metadata.ITunesResponse, error)
}

func (s Server) HandleGetAlbum(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "no id provided", http.StatusNotFound)

		return
	}

	res, err := s.iTunesclient.GetAlbumMetadata(id)
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

	album := ConvertItunesAlbum(res)

	response := NewEmptyResponse(s.cfg)
	response.SubsonicResponse.Album = &album
	WriteJSON(w, http.StatusOK, response)
}

func (s Server) HandleGetArtist(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "no id provided", http.StatusNotFound)

		return
	}

	res, err := s.iTunesclient.GetArtistById(id)
	if err != nil {
		slog.Error("something went wrong trying to fetch metadata", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	artist := ConvertItunesArtist(res)

	response := NewEmptyResponse(s.cfg)
	response.SubsonicResponse.Artist = &artist
	WriteJSON(w, http.StatusOK, response)
}
