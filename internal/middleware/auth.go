package middleware

import (
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
	"net/http"

	"github.com/Simon-Weij/allium/internal/config"
	"github.com/Simon-Weij/allium/internal/subsonic"
)

func Authenticate(cfg config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, salt, token := parseAuthParams(r)
			if !validCredentialsPresent(w, cfg, r, username, salt, token) {
				return
			}

			if !isValidUser(w, cfg, username, salt, token) {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func parseAuthParams(r *http.Request) (username, salt, token string) {
	query := r.URL.Query()

	return query.Get("u"), query.Get("s"), query.Get("t")
}

func validCredentialsPresent(
	w http.ResponseWriter,
	cfg config.Config,
	r *http.Request,
	username, salt, token string,
) bool {
	query := r.URL.Query()

	password := query.Get("p")
	apiKey := query.Get("apiKey")

	if conflictingAuth(salt, token, password, apiKey) {
		subsonic.WriteError(
			w,
			http.StatusBadRequest,
			cfg,
			subsonic.ErrConflictingAuthenticationMechanisms,
			"Conflicting auth methods",
		)

		return false
	}

	if unsupportedAuth(salt, token, password, apiKey) {
		subsonic.WriteError(
			w,
			http.StatusBadRequest,
			cfg,
			subsonic.ErrNotSupported,
			"Auth method not supported",
		)

		return false
	}

	if missingCredentials(username, salt, token) {
		subsonic.WriteError(
			w,
			http.StatusBadRequest,
			cfg,
			subsonic.ErrParameterMissing,
			"Required parameter is missing",
		)

		return false
	}

	return true
}

func conflictingAuth(salt, token, password, apiKey string) bool {
	hasToken := token != "" && salt != ""
	hasUnsupported := password != "" || apiKey != ""
	hasBothUnsupported := password != "" && apiKey != ""

	return hasToken && hasUnsupported || hasBothUnsupported
}

func unsupportedAuth(salt, token, password, apiKey string) bool {
	hasToken := token != "" && salt != ""
	if hasToken {
		return false
	}

	return password != "" || apiKey != ""
}

func missingCredentials(username, salt, token string) bool {
	hasToken := token != "" && salt != ""

	return username == "" || !hasToken
}

func isValidUser(
	w http.ResponseWriter,
	cfg config.Config,
	username,
	salt,
	token string,
) bool {
	isValidUser := subtle.ConstantTimeCompare([]byte(cfg.Username), []byte(username)) == 1
	if !isValidUser || !matchToken(cfg.Password, salt, token) {
		subsonic.WriteError(
			w,
			http.StatusUnauthorized,
			cfg,
			subsonic.ErrWrongCredentials,
			"Wrong username or password",
		)

		return false
	}

	return true
}

func matchToken(storedPassword, salt, token string) bool {
	sum := md5.Sum([]byte(storedPassword + salt))
	return subtle.ConstantTimeCompare([]byte(hex.EncodeToString(sum[:])), []byte(token)) == 1
}
