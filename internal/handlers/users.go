package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Simon-Weij/allium/internal/database"
	internalErrors "github.com/Simon-Weij/allium/internal/errors"

	_ "modernc.org/sqlite"
)

func (s Server) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	reqvalues := r.URL.Query()

	username := reqvalues.Get("username")
	if username == "" {
		http.Error(w, "username parameter missing", internalErrors.ErrParameterMissing)
	}

	user, err := getUser(r.Context(), username)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		slog.Error("error closing db in HandleGetUser", "err: ", err)

		return
	}

	res := NewEmptyResponse(s.cfg)
	res.SubsonicResponse.User = &user

	WriteJSON(w, http.StatusOK, res)
}

func getUser(ctx context.Context, username string) (database.User, error) {
	dbConn, err := sql.Open("sqlite", "../../database.db") // TODO: Make this less messy
	if err != nil {
		return database.User{}, fmt.Errorf("error returned while opening database %w", err)
	}
	defer func() {
		if err = dbConn.Close(); err != nil {
			err = errors.Join(err)
		}
	}()

	queries := database.New(dbConn)

	user, err := queries.GetUser(ctx, username)
	if err != nil {
		return database.User{}, fmt.Errorf("error returned while getting user from database %w", err)
	}

	return user, nil
}
