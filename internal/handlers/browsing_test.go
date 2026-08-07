package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleGetMusicFolders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/rest/getMusicFolders", nil)
	w := httptest.NewRecorder()

	cfg := setupTestingConfig(t)

	browsingHandler := NewBrowsingHandler(cfg)
	browsingHandler.HandleGetMusicFolders(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	expectedRes := NewEmptyResponse(cfg)
	expectedRes.SubsonicResponse.MusicFolders = &MusicFolders{
		MusicFolder: []MusicFolder{
			{
				Id:   1,
				Name: "music",
			},
		},
	}

	expectedJSON, err := json.Marshal(expectedRes)
	assert.NoError(t, err)

	assert.JSONEq(t, string(expectedJSON), string(data))
}
