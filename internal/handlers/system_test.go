package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

	expectedRes := NewEmptyResponse(cfg)

	expectedJSON, err := json.Marshal(expectedRes)
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

	expectedRes := NewEmptyResponse(cfg)
	expectedRes.SubsonicResponse.License = &License{
		Valid:          true,
		Email:          "irrelevant@example.com",
		LicenseExpires: &licenseExpiry,
		TrialExpires:   &licenseExpiry,
	}

	expectedJSON, err := json.Marshal(expectedRes)
	assert.NoError(t, err)

	assert.JSONEq(t, string(expectedJSON), string(data))
}

func TestHandleGetOpenSubSonicExtensions(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/rest/getOpenSubsonicExtensions", nil)
	w := httptest.NewRecorder()

	cfg := setupTestingConfig(t)

	SystemHandler := NewSystemHandler(cfg)
	SystemHandler.HandleGetOpenSubsonicExtensions(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	expectedRes := NewEmptyResponse(cfg)
	expectedRes.SubsonicResponse.OpenSubSonicExtensions = &[]OpenSubSonicExtension{}

	expectedJSON, err := json.Marshal(expectedRes)
	assert.NoError(t, err)

	assert.JSONEq(t, string(expectedJSON), string(data))
}
