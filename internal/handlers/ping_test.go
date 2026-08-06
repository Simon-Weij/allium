package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Simon-Weij/allium/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/rest/ping", nil)
	w := httptest.NewRecorder()

	cfg := config.Config{
		Name:    "allium",
		Version: "0.1.0",
	}

	pingHandler := NewPingHandler(cfg)
	pingHandler.HandlePing(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	expected := NewSubsonicResponse(cfg)

	expectedJSON, err := json.Marshal(expected)
	assert.NoError(t, err)

	assert.JSONEq(t, string(expectedJSON), string(data))
}
