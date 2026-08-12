package handlers

import (
	"net/http"
	"testing"
	"time"
)

func TestPing(t *testing.T) {
	t.Parallel()
	cfg := setupTestingConfig(t)
	server := NewServer(cfg)
	res := NewEmptyResponse(cfg)

	assertJSONResponse(
		t,
		http.MethodGet,
		"/rest/ping.view",
		server.HandlePing,
		res,
	)
}

func TestLicense(t *testing.T) {
	t.Parallel()
	cfg := setupTestingConfig(t)
	server := NewServer(cfg)
	res := NewEmptyResponse(cfg)

	licenseExpiry := time.Date(2222, 1, 1, 1, 1, 1, 1, time.UTC)
	res.SubsonicResponse.License = &License{
		Valid:          true,
		Email:          "irrelevant@example.com",
		LicenseExpires: &licenseExpiry,
		TrialExpires:   &licenseExpiry,
	}

	assertJSONResponse(
		t,
		http.MethodGet,
		"/rest/getLicense.view",
		server.HandleGetLicense,
		res,
	)
}

func TestHandleGetOpenSubSonicExtensions(t *testing.T) {
	t.Parallel()
	cfg := setupTestingConfig(t)
	server := NewServer(cfg)
	res := NewEmptyResponse(cfg)

	res.SubsonicResponse.OpenSubSonicExtensions = &[]OpenSubSonicExtension{}

	assertJSONResponse(
		t,
		http.MethodGet,
		"/rest/getOpenSubsonicExtensions",
		server.HandleGetOpenSubsonicExtensions,
		res,
	)
}
