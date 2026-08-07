package middleware

import (
	"net/http"

	"github.com/Simon-Weij/allium/internal/config"
	"github.com/Simon-Weij/allium/internal/handlers"
)

func RequireApiKey(cfg config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.URL.Query().Get("apiKey")
			if apiKey != cfg.ApiKey {
				handlers.WriteJSON(w, http.StatusUnauthorized, handlers.Generic401Response{})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
