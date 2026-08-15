package handlers

import (
	"net/http"
	"testing"
	"time"

	"github.com/Simon-Weij/allium/internal/testutil"
)

func TestPing(t *testing.T) {
	t.Parallel()
	cfg := testutil.SetupTestingConfig(t)
	server := NewServer(cfg, nil)
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
	cfg := testutil.SetupTestingConfig(t)
	server := NewServer(cfg, nil)
	res := NewEmptyResponse(cfg)

	licenseExpiry := time.Date(2222, 1, 1, 1, 1, 1, 1, time.UTC)
	res.SubsonicResponse.License = &License{
		Valid:          true,
		Email:          testEmail,
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
	cfg := testutil.SetupTestingConfig(t)
	server := NewServer(cfg, nil)
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
