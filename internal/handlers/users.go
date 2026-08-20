package handlers

import (
	"database/sql"
	"errors"
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
		subsonic.WriteError(w, http.StatusInternalServerError, s.cfg, subsonic.ErrParameterMissing, "username parameter missing")
		slog.Error("user did not pass value for username parameter")
	}

	user, err := s.Userclient.GetUserByUsername(ctx, username) 
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

	users, err := s.Userclient.GetUsers(ctx)
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

func (s Server) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	reqValues := r.URL.Query()

	if !reqValues.Has("username") {
		subsonic.WriteError(w, http.StatusBadRequest, s.cfg, subsonic.ErrParameterMissing, "user did not pass username value")
		slog.Error("user did not pass username value")
		return
	} else if !reqValues.Has("password") {
		subsonic.WriteError(w, http.StatusBadRequest, s.cfg, subsonic.ErrParameterMissing, "user did not pass password value")
		slog.Error("user did not pass password value")
		return
	} else if !reqValues.Has("email") {
		subsonic.WriteError(w, http.StatusBadRequest, s.cfg, subsonic.ErrParameterMissing, "user did not pass email value")
		slog.Error("user did not pass email value")
		return
	}
	
	err := s.Userclient.CreateUser(r.Context(), reqValues)
	if err != nil {
		subsonic.WriteError(w, http.StatusInternalServerError, s.cfg, subsonic.ErrGeneric, "generic error when trying to HandleCreateUser")
		slog.Error("error recieved when trying to HandleCreateUser", "err: ", err)
	}
}