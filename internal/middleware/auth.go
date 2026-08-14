package middleware

import (
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
	"net/http"

	"github.com/Simon-Weij/allium/internal/config"
	"github.com/Simon-Weij/allium/internal/errors"
	"github.com/Simon-Weij/allium/internal/handlers"
)

// Authenticate Middleware, checks if client has Token Authentication, Unsupported Auth or Conflicting Authentication Methods
// returning errors to the client if using Conflicting or Unsupported Authentication methods
// Also returns errors if the client credentials are invalid.
func Authenticate(cfg config.Config) func(http.Handler) http.Handler {
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
				handlers.WriteError(
					w,
					http.StatusBadRequest,
					cfg,
					errors.ErrConflictingAuthenticationMechanisms,
					"Conflicting auth methods",
				)

				return
			}

			if hasUnsupportedAuth {
				handlers.WriteError(
					w,
					http.StatusBadRequest,
					cfg,
					errors.ErrNotSupported,
					"Auth method not supported",
				)

				return
			}

			if username == "" || !hasTokenAuth {
				handlers.WriteError(
					w,
					http.StatusBadRequest,
					cfg,
					errors.ErrParameterMissing,
					"Required parameter is missing",
				)

				return
			}

			isValidUser(w, cfg, username, salt, token)

			next.ServeHTTP(w, r)
		})
	}
}

// isValidUser checks if the user's username and password are correct,
// returing a http.StatusUnauthorized to the client, if neither the username or the password is valid.
func isValidUser(
	w http.ResponseWriter,
	cfg config.Config,
	username,
	salt,
	token string,
) {
	isValidUser := subtle.ConstantTimeCompare([]byte(cfg.Username), []byte(username)) == 1
	if !isValidUser || !matchToken(cfg.Password, salt, token) {
		handlers.WriteError(
			w,
			http.StatusUnauthorized,
			cfg,
			errors.ErrWrongCredentials,
			"Wrong username or password",
		)
	}
}

// matchToken calculates a MD5 hashsum from the storedPassword & salt,
// and ConstantTimeCompares the hashsum and token, returning a boolean value (1 if match, 0 if no match)
func matchToken(storedPassword, salt, token string) bool {
	sum := md5.Sum([]byte(storedPassword + salt))

	return subtle.ConstantTimeCompare([]byte(hex.EncodeToString(sum[:])), []byte(token)) == 1
}
