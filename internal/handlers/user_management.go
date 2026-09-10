package handlers

import (
	"net/http"

	"github.com/Simon-Weij/allium/internal/subsonic"
)

func (s Server) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	res := subsonic.NewEmptyResponse(s.cfg)
	res.SubsonicResponse.User = &subsonic.User{
		Folder:    []int{1},
		Email:     testEmail,
		AdminRole: true,
	}

	subsonic.WriteJSON(w, http.StatusOK, res)
}
