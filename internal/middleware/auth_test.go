package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Simon-Weij/allium/internal/config"

	"github.com/stretchr/testify/assert"
)

func token(password, salt string) string {
	sum := md5.Sum([]byte(password + salt))
	return hex.EncodeToString(sum[:])
}

func TestAuthenticate(t *testing.T) {
	cfg := config.Config{
		Username: "alice",
		Password: "secret",
	}

	tests := []struct {
		name           string
		query          string
		wantStatus     int
		wantCode       int
		wantNextCalled bool
	}{
		{
			name:           "valid token passes through",
			query:          "/?u=alice&t=" + token("secret", "c19b2d") + "&s=c19b2d",
			wantStatus:     http.StatusOK,
			wantCode:       0,
			wantNextCalled: true,
		},
		{
			name:           "wrong username",
			query:          "/?u=other&t=" + token("secret", "c19b2d") + "&s=c19b2d",
			wantStatus:     http.StatusUnauthorized,
			wantCode:       40,
			wantNextCalled: false,
		},
		{
			name:           "wrong token",
			query:          "/?u=alice&t=deadbeef&s=c19b2d",
			wantStatus:     http.StatusUnauthorized,
			wantCode:       40,
			wantNextCalled: false,
		},
		{
			name:           "missing salt",
			query:          "/?u=alice&t=" + token("secret", "c19b2d"),
			wantStatus:     http.StatusBadRequest,
			wantCode:       10,
			wantNextCalled: false,
		},
		{
			name:           "missing token",
			query:          "/?u=alice&s=c19b2d",
			wantStatus:     http.StatusBadRequest,
			wantCode:       10,
			wantNextCalled: false,
		},
		{
			name:           "missing username",
			query:          "/?t=" + token("secret", "c19b2d") + "&s=c19b2d",
			wantStatus:     http.StatusBadRequest,
			wantCode:       10,
			wantNextCalled: false,
		},
		{
			name:           "no credentials",
			query:          "/",
			wantStatus:     http.StatusBadRequest,
			wantCode:       10,
			wantNextCalled: false,
		},
		{
			name:           "legacy password not supported",
			query:          "/?u=alice&p=secret",
			wantStatus:     http.StatusBadRequest,
			wantCode:       42,
			wantNextCalled: false,
		},
		{
			name:           "api key not supported",
			query:          "/?apiKey=abc",
			wantStatus:     http.StatusBadRequest,
			wantCode:       42,
			wantNextCalled: false,
		},
		{
			name:           "token and password conflict",
			query:          "/?u=alice&p=secret&t=" + token("secret", "c19b2d") + "&s=c19b2d",
			wantStatus:     http.StatusBadRequest,
			wantCode:       43,
			wantNextCalled: false,
		},
		{
			name:           "token and api key conflict",
			query:          "/?u=alice&apiKey=abc&t=" + token("secret", "c19b2d") + "&s=c19b2d",
			wantStatus:     http.StatusBadRequest,
			wantCode:       43,
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

			handler := Authenticate(cfg)(next)

			req := httptest.NewRequest(http.MethodGet, tt.query, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, tt.wantNextCalled, nextCalled)

			if !tt.wantNextCalled {
				assert.Contains(t, rec.Body.String(), `"code":`+strconv.Itoa(tt.wantCode))
			}
		})
	}
}
