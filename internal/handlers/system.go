package handlers

import (
	"net/http"
	"time"
)

// HandlePing handles the ping endpoint, it is used by clients to test for connectivity
func (s Server) HandlePing(w http.ResponseWriter, r *http.Request) {
	res := NewEmptyResponse(s.cfg)
	WriteJSON(w, http.StatusOK, res)
}

// HandleGetLicense handles the getLicense endpoint, it returns details about the software license
func (s Server) HandleGetLicense(w http.ResponseWriter, r *http.Request) {
	licenseExpiry := time.Date(2222, 1, 1, 1, 1, 1, 1, time.UTC)
	res := NewEmptyResponse(s.cfg)
	res.SubsonicResponse.License = &License{
		Valid:          true,
		Email:          testEmail,
		LicenseExpires: &licenseExpiry,
		TrialExpires:   &licenseExpiry,
	}

	WriteJSON(w, http.StatusOK, res)
}

// HandleGetOpenSubsonicExtensions handles the extensions for OpenSubsonic
func (s Server) HandleGetOpenSubsonicExtensions(w http.ResponseWriter, r *http.Request) {
	res := NewEmptyResponse(s.cfg)
	res.SubsonicResponse.OpenSubSonicExtensions = &[]OpenSubSonicExtension{}
	WriteJSON(w, http.StatusOK, res)
}
