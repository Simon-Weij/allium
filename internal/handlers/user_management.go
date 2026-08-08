package handlers

import (
	"net/http"

	"github.com/Simon-Weij/allium/internal/config"
)

type UserManagementHandler struct {
	cfg config.Config
}

func NewUserManagementHandler(cfg config.Config) UserManagementHandler {
	return UserManagementHandler{
		cfg: cfg,
	}
}

// TODO: stop returning placeholders
func (h UserManagementHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	res := NewEmptyResponse(h.cfg)
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

	WriteJSON(w, http.StatusOK, res)
}
