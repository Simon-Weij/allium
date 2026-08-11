package handlers

import (
	"io"
	"net/http"
)

// TODO: stop returning placeholders
func (s Server) HandleGetMusicFolders(w http.ResponseWriter, r *http.Request) {
	res := NewEmptyResponse(s.cfg)
	res.SubsonicResponse.MusicFolders = &MusicFolders{
		MusicFolder: []MusicFolder{
			{
				Id:   1,
				Name: "music",
			},
		},
	}
	WriteJSON(w, http.StatusOK, res)
}

// TODO: stop returning placeholders
func (s Server) HandleGetGenres(w http.ResponseWriter, r *http.Request) {
	res := NewEmptyResponse(s.cfg)
	res.SubsonicResponse.Genres = &[]Genre{
		{
			SongCount:  1,
			AlbumCount: 1,
			Value:      "Punk",
		},
		{
			SongCount:  4,
			AlbumCount: 1,
			Value:      "Dark Ambient",
		},
		{
			SongCount:  6,
			AlbumCount: 1,
			Value:      "Noise",
		},
		{
			SongCount:  11,
			AlbumCount: 1,
			Value:      "Dance",
		},
		{
			SongCount:  12,
			AlbumCount: 1,
			Value:      "Electronic",
		},
		{
			SongCount:  20,
			AlbumCount: 1,
			Value:      "Hip-Hop",
		},
	}
	WriteJSON(w, http.StatusOK, res)
}

// TODO: stop returning placeholders
func (s Server) HandleGetCoverArt(w http.ResponseWriter, r *http.Request) {
	imageURL := "https://placehold.co/400x400"

	resp, err := http.Get(imageURL)
	if err != nil {
		http.Error(w, "failed to fetch image", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "image unavailable", resp.StatusCode)
		return
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	_, _ = io.Copy(w, resp.Body)
}
