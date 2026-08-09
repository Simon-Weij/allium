package handlers

import (
	"net/http"
	"testing"
)

func TestHandleGetUser(t *testing.T) {
	cfg := setupTestingConfig(t)
	server := NewServer(cfg)
	res := NewEmptyResponse(cfg)
	res.SubsonicResponse.User = &User{
		Folder:            []int{1},
		Email:             "irrelevant@example.com",
		ScrobblingEnabled: false,
		AdminRole:         true,
		SettingsRole:      false,
		DownloadRole:      false,
		PlaylistRole:      false,
		CoverArtRole:      false,
		CommentRole:       false,
		PodcastRole:       false,
		StreamRole:        false,
		JukeboxRole:       false,
		ShareRole:         false,
	}

	assertJSONResponse(
		t,
		http.MethodGet,
		"/rest/getUser.view",
		server.HandleGetUser,
		res,
	)
}
