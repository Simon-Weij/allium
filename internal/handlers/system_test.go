package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Simon-Weij/allium/internal/config"
	"github.com/stretchr/testify/assert"
)

func setupTestingConfig(t *testing.T) config.Config {
	t.Helper()

	cfg := config.Config{
		Name:    "allium",
		Version: "0.1.0",
		ApiKey:  "irrelevant",
	}

	return cfg
}

func TestPing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/rest/ping", nil)
	w := httptest.NewRecorder()

	cfg := setupTestingConfig(t)

	systemHandler := NewSystemHandler(cfg)
	systemHandler.HandlePing(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	expected := SubsonicResponse{
		Body: NewSubsonicResponse(cfg),
	}

	expectedJSON, err := json.Marshal(expected)
	assert.NoError(t, err)

	assert.JSONEq(t, string(expectedJSON), string(data))
}

func TestLicense(t *testing.T) {
	licenseExpiry := time.Date(2222, 1, 1, 1, 1, 1, 1, time.UTC)

	req := httptest.NewRequest(http.MethodGet, "/rest/getLicense", nil)
	w := httptest.NewRecorder()

	cfg := setupTestingConfig(t)

	systemHandler := NewSystemHandler(cfg)
	systemHandler.HandleGetLicense(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	expected := LicenseResponse{
		SubsonicResponseBody: NewSubsonicResponse(cfg),
		License: License{
			Valid:          true,
			Email:          "irrelevant@example.com",
			LicenseExpires: licenseExpiry,
			TrialExpires:   licenseExpiry,
		},
	}

	expectedJSON, err := json.Marshal(expected)
	assert.NoError(t, err)

	assert.JSONEq(t, string(expectedJSON), string(data))
}
