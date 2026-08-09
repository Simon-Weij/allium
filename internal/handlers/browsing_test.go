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
