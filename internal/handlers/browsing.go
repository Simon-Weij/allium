package handlers

import (
	"net/http"

	"github.com/Simon-Weij/allium/internal/config"
)

type BrowsingHandler struct {
	cfg config.Config
}

func NewBrowsingHandler(cfg config.Config) BrowsingHandler {
	return BrowsingHandler{
		cfg: cfg,
	}
}

func (h BrowsingHandler) HandleGetMusicFolders(w http.ResponseWriter, r *http.Request) {
	res := NewEmptyResponse(h.cfg)
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
