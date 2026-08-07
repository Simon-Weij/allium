package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleGetUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/rest/getUser", nil)
	w := httptest.NewRecorder()

	cfg := setupTestingConfig(t)

	userManagementHandler := NewUserManagementHandler(cfg)
	userManagementHandler.HandleGetUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	expectedRes := NewEmptyResponse(cfg)
	expectedRes.SubsonicResponse.User = &User{
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

	expectedJSON, err := json.Marshal(expectedRes)
	assert.NoError(t, err)

	assert.JSONEq(t, string(expectedJSON), string(data))
}
