package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Simon-Weij/allium/internal/subsonic"

	_ "modernc.org/sqlite"
)

func (s Server) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqvalues := r.URL.Query()
	username := reqvalues.Get("username")

	if username == "" {
		http.Error(w, "username parameter missing", subsonic.ErrParameterMissing)
	}

	user, err := s.queries.GetUser(ctx, username)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		slog.Error("error closing db in HandleGetUser", "err: ", err)

		return
	}

	res := subsonic.NewEmptyResponse(s.cfg)
	res.SubsonicResponse.User = &user
	subsonic.WriteJSON(w, http.StatusOK, res)

}
