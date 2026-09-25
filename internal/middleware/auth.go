package middleware

import (
	"context"
	"crypto/md5"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Simon-Weij/allium/generated/sqlc"
	"github.com/Simon-Weij/allium/internal/config"
	"github.com/Simon-Weij/allium/internal/subsonic"
)

var ErrInvalidCreds = errors.New("Invalid Username or Password")

func Authenticate(cfg config.Config, queries *sqlc.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()

			username := query.Get("u")
			password := query.Get("p")
			apiKey := query.Get("apiKey")
			token := query.Get("t")
			salt := query.Get("s")

			hasTokenAuth := token != "" && salt != ""
			hasUnsupportedAuth := password != "" || apiKey != ""
			hasBothUnsupported := password != "" && apiKey != ""

			if (hasTokenAuth && hasUnsupportedAuth) || hasBothUnsupported {
				subsonic.WriteError(
					w,
					http.StatusBadRequest,
					cfg,
					subsonic.ErrConflictingAuthenticationMechanisms,
					"Conflicting auth methods",
				)

				return
			}

			if hasUnsupportedAuth {
				subsonic.WriteError(
					w,
					http.StatusBadRequest,
					cfg,
					subsonic.ErrNotSupported,
					"Auth method not supported",
				)

				return
			}

			if username == "" || !hasTokenAuth {
				subsonic.WriteError(
					w,
					http.StatusBadRequest,
					cfg,
					subsonic.ErrParameterMissing,
					"Required parameter is missing",
				)

				return
			}

			err := isValidUser(r.Context(), username, salt, token, queries)
			if errors.Is(err, ErrInvalidCreds) {
				subsonic.WriteError(
					w,
					http.StatusUnauthorized,
					cfg,
					subsonic.ErrWrongCredentials,
					"Wrong username or password",
				)
				return 
			}
			if err != nil {
				subsonic.WriteError(w, http.StatusInternalServerError, cfg, subsonic.ErrGeneric, "generic error when trying to authenticate")
				slog.Error("error recieved when trying to authenticate", "username", username, "err: ", err)
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

func isValidUser(
	ctx context.Context,
	username,
	salt,
	token string,
	queries *sqlc.Queries,
) error {
	
	_, err := queries.GetUser(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrInvalidCreds
		}	
		return fmt.Errorf("error recieved trying to authenticate %w", err)
	}

	password, err := queries.GetPassword(ctx, username)
	if err != nil {
		return fmt.Errorf("error recieved trying to authenticate %w", err)
	}
	if !matchToken(password, salt, token) {
		return ErrInvalidCreds
	}
	return nil
}

func matchToken(storedPassword, salt, token string) bool {
	sum := md5.Sum([]byte(storedPassword + salt))

	return subtle.ConstantTimeCompare([]byte(hex.EncodeToString(sum[:])), []byte(token)) == 1
}
