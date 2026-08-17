package handlers

import (
	"net/http"

	"github.com/Simon-Weij/allium/internal/subsonic"
)

func (s Server) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	res := subsonic.NewEmptyResponse(s.cfg)
	res.SubsonicResponse.User = &subsonic.User{
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

	subsonic.WriteJSON(w, http.StatusOK, res)
}
