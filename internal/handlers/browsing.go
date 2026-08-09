package handlers

import (
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
