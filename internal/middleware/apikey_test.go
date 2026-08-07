package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Simon-Weij/allium/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestRequireApiKey(t *testing.T) {
	cfg := config.Config{
		Name:    "irrelevant",
		Version: "irrelevant",
		ApiKey:  "secret-api-key",
	}

	tests := []struct {
		name           string
		apiKey         string
		wantStatus     int
		wantNextCalled bool
	}{
		{
			name:           "valid api key goes through",
			apiKey:         "secret-api-key",
			wantStatus:     http.StatusOK,
			wantNextCalled: true,
		},
		{
			name:           "invalid api key shouldn't go through",
			apiKey:         "wrong-key",
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "missing api key shouldn't go through",
			apiKey:         "",
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			handler := RequireApiKey(cfg)(next)

			req := httptest.NewRequest(http.MethodGet, "/?apiKey="+tt.apiKey, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, tt.wantNextCalled, nextCalled)
		})
	}
}
