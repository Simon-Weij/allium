package handlers

import (
	"net/http"

	"github.com/Simon-Weij/allium/internal/config"
)

type JukeboxHandler struct {
	cfg config.Config
}

func NewJukeboxHandler(cfg config.Config) JukeboxHandler {
	return JukeboxHandler{
		cfg: cfg,
	}
}

// TODO: stop returning placeholders
func (h JukeboxHandler) HandleJukeboxStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.NotFound(w, r)
		return
	}
	res := NewEmptyResponse(h.cfg)
	res.SubsonicResponse.JukeboxPlaylist = &JukeboxPlaylist{
		CurrentIndex: 1,
		Playing:      false,
		Gain:         0.5,
		Position:     67,
		Entry:        nil,
	}
	WriteJSON(w, http.StatusOK, res)
}

// TODO: stop returning placeholders
func (h JukeboxHandler) HandleJukeboxPlaylist(w http.ResponseWriter, r *http.Request) {
	res := NewEmptyResponse(h.cfg)
	res.SubsonicResponse.JukeboxPlaylist = &JukeboxPlaylist{
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
	WriteJSON(w, http.StatusOK, res)
}
