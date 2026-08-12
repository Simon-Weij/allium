package handlers

import (
	"net/http"
)

func (s Server) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	res := NewEmptyResponse(s.cfg)
	res.SubsonicResponse.User = &User{
		Folder:            []int{1},
		Email:             testEmail,
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

	WriteJSON(w, http.StatusOK, res)
}
