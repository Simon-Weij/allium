// TOOD: replace this with proper auth

package middleware

import (
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
	"log/slog"
	"net/http"

	"github.com/Simon-Weij/allium/internal/config"
	"github.com/Simon-Weij/allium/internal/handlers"
)

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
				handlers.WriteError(w, http.StatusBadRequest, cfg, 43, "Conflicting auth methods")
				return
			}

			if hasUnsupportedAuth {
				handlers.WriteError(w, http.StatusBadRequest, cfg, 42, "Auth method not supported")
				return
			}

			if username == "" || !hasTokenAuth {
				handlers.WriteError(w, http.StatusBadRequest, cfg, 10, "Required parameter is missing")
				return
			}

			isValidUser := subtle.ConstantTimeCompare([]byte(cfg.Username), []byte(username)) == 1
			if !isValidUser || !matchToken(cfg.Password, salt, token) {
				handlers.WriteError(w, http.StatusUnauthorized, cfg, 40, "Wrong username or password")
				return
			}
			slog.Info("successfully logged " + username + " in")

			next.ServeHTTP(w, r)
		})
	}
}

func matchToken(storedPassword, salt, token string) bool {
	sum := md5.Sum([]byte(storedPassword + salt))
	return subtle.ConstantTimeCompare([]byte(hex.EncodeToString(sum[:])), []byte(token)) == 1
}
