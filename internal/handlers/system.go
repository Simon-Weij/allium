package handlers

import (
	"net/http"
	"time"
)

func (s Server) HandlePing(w http.ResponseWriter, r *http.Request) {
	res := NewEmptyResponse(s.cfg)
	WriteJSON(w, http.StatusOK, res)
}

func (s Server) HandleGetLicense(w http.ResponseWriter, r *http.Request) {
	licenseExpiry := time.Date(2222, 1, 1, 1, 1, 1, 1, time.UTC)
	res := NewEmptyResponse(s.cfg)
	res.SubsonicResponse.License = &License{
		Valid:          true,
		Email:          "irrelevant@example.com",
		LicenseExpires: &licenseExpiry,
		TrialExpires:   &licenseExpiry,
	}

	WriteJSON(w, http.StatusOK, res)
}

func (s Server) HandleGetOpenSubsonicExtensions(w http.ResponseWriter, r *http.Request) {
	res := NewEmptyResponse(s.cfg)
	res.SubsonicResponse.OpenSubSonicExtensions = &[]OpenSubSonicExtension{}
	WriteJSON(w, http.StatusOK, res)
}
