package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleJukeboxStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/rest/jukeboxControl.view", nil)
	w := httptest.NewRecorder()

	cfg := setupTestingConfig(t)

	jukeboxHandler := NewJukeboxHandler(cfg)
	jukeboxHandler.HandleJukeboxStatus(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	expectedRes := NewEmptyResponse(cfg)
	expectedRes.SubsonicResponse.JukeboxPlaylist = &JukeboxPlaylist{
		CurrentIndex: 1,
		Playing:      false,
		Gain:         0.5,
		Position:     67,
		Entry:        nil,
	}

	expectedJSON, err := json.Marshal(expectedRes)
	assert.NoError(t, err)

	assert.JSONEq(t, string(expectedJSON), string(data))
}

func TestHandleJukeboxPlaylist(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/rest/jukeboxControl.view", nil)
	w := httptest.NewRecorder()

	cfg := setupTestingConfig(t)

	jukeboxHandler := NewJukeboxHandler(cfg)
	jukeboxHandler.HandleJukeboxPlaylist(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	expectedRes := NewEmptyResponse(cfg)
	expectedRes.SubsonicResponse.JukeboxPlaylist = &JukeboxPlaylist{
		CurrentIndex: 1,
		Playing:      false,
		Gain:         0.5,
		Position:     67,
		Entry: &JukeboxPlaylistEntry{
			Id:           "300000116",
			Parent:       "200000021",
			Title:        "Test",
			IsDir:        false,
			IsVideo:      false,
			Type:         "music",
			AlbumId:      "200000021",
			Album:        "test album",
			ArtistId:     "100000036",
			Artist:       "test artist",
			CoverArt:     "300000116",
			Duration:     103,
			BitRate:      216,
			BitDepth:     16,
			SamplingRate: 44100,
			ChannelCount: 2,
			Track:        1,
			Year:         2005,
			Genre:        "Rock",
			Size:         2811819,
			DiscNumber:   1,
			Suffix:       "mp3",
			ContentType:  "audio/mpeg",
			Path:         "/",
		},
	}

	expectedJSON, err := json.Marshal(expectedRes)
	assert.NoError(t, err)

	assert.JSONEq(t, string(expectedJSON), string(data))
}
