package handlers

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	
	"github.com/Simon-Weij/allium/internal/subsonic"
	"github.com/Simon-Weij/allium/internal/users"

	_ "modernc.org/sqlite"
)

func (s Server) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqvalues := r.URL.Query()
	username := reqvalues.Get("username")

	if username == "" {
		http.Error(w, "username parameter missing", subsonic.ErrParameterMissing)
	}

	user, err := users.GetUserByUsername(ctx, username, s.queries) 
	if errors.Is(err, sql.ErrNoRows) {
		subsonic.WriteError(w, http.StatusInternalServerError, s.cfg, subsonic.ErrRequestedDataNotFound, "user was not found in database")
		slog.Error("user was not found in database", "username", username, "err: ", err)
	}
	if err != nil {
		subsonic.WriteError(w, http.StatusInternalServerError, s.cfg, subsonic.ErrGeneric, "generic error when trying to HandleGetUser")
		slog.Error("error recieved when trying to GetUserByUsername", "username", username, "err: ", err)
	}

	res := subsonic.NewEmptyResponse(s.cfg)
	res.SubsonicResponse.User = user
	subsonic.WriteJSON(w, http.StatusOK, res)

}

func (s Server) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := users.GetUsers(ctx, s.queries)
	if errors.Is(err, sql.ErrNoRows) {
		subsonic.WriteError(w, http.StatusInternalServerError, s.cfg, subsonic.ErrRequestedDataNotFound, "user database is empty")
		slog.Error("user was not found in database", "err: ", err)
	}
	if err != nil {
		subsonic.WriteError(w, http.StatusInternalServerError, s.cfg, subsonic.ErrGeneric, "generic error when trying to HandleGetUsers")
		slog.Error("error recieved when trying to GetUserByUsername", "err: ", err)
	}

	res := subsonic.NewEmptyResponse(s.cfg)
	res.SubsonicResponse.Users = users
	subsonic.WriteJSON(w, http.StatusOK, res)
}