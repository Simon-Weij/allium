package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Simon-Weij/allium/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestingConfig(t *testing.T) config.Config {
	t.Helper()

	cfg := config.Config{
		Username: "alice",
		Password: "secret",
	}

	return cfg
}

func assertJSONResponse(
	t *testing.T,
	method string,
	path string,
	handler http.HandlerFunc,
	expected any,
) {
	t.Helper()

	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()

	handler(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	expectedJSON, err := json.Marshal(expected)
	require.NoError(t, err)

	assert.JSONEq(t, string(expectedJSON), string(data))
}
