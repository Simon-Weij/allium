package handlers

import (
	"net/http"
	"testing"
)

func TestHandleGetMusicFolders(t *testing.T) {
	cfg := setupTestingConfig(t)
	server := NewServer(cfg)
	res := NewEmptyResponse(cfg)
	res.SubsonicResponse.MusicFolders = &MusicFolders{
		MusicFolder: []MusicFolder{
			{
				Id:   1,
				Name: "music",
			},
		},
	}

	assertJSONResponse(
		t,
		http.MethodGet,
		"/rest/getMusicFolders.view",
		server.HandleGetMusicFolders,
		res,
	)
}

func TestHandleGetGenres(t *testing.T) {
	cfg := setupTestingConfig(t)
	server := NewServer(cfg)
	res := NewEmptyResponse(cfg)
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

	assertJSONResponse(
		t,
		http.MethodGet,
		"/rest/getGenres.view",
		server.HandleGetGenres,
		res,
	)
}
