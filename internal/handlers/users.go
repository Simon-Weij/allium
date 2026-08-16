package handlers

import (
	"log/slog"
	"net/http"

	internalErrors "github.com/Simon-Weij/allium/internal/errors"

	_ "modernc.org/sqlite"
)

func (s Server) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqvalues := r.URL.Query()
	username := reqvalues.Get("username")

	if username == "" {
		http.Error(w, "username parameter missing", internalErrors.ErrParameterMissing)
	}

	user, err := s.queries.GetUser(ctx, username)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		slog.Error("error closing db in HandleGetUser", "err: ", err)

		return
	}

	res := NewEmptyResponse(s.cfg)
	res.SubsonicResponse.User = &user

	WriteJSON(w, http.StatusOK, res)
}
