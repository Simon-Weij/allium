package handlers

import (
	"net/http"
	"time"

	"github.com/Simon-Weij/allium/internal/subsonic"
)

func (s Server) HandlePing(w http.ResponseWriter, r *http.Request) {
	res := subsonic.NewEmptyResponse(s.cfg)
	subsonic.WriteJSON(w, http.StatusOK, res)
}

func (s Server) HandleGetLicense(w http.ResponseWriter, r *http.Request) {
	licenseExpiry := time.Date(2222, 1, 1, 1, 1, 1, 1, time.UTC)
	res := subsonic.NewEmptyResponse(s.cfg)
	res.SubsonicResponse.License = &subsonic.License{
		Valid:          true,
		Email:          testEmail,
		LicenseExpires: &licenseExpiry,
		TrialExpires:   &licenseExpiry,
	}

	subsonic.WriteJSON(w, http.StatusOK, res)
}

func (s Server) HandleGetOpenSubsonicExtensions(w http.ResponseWriter, r *http.Request) {
	res := subsonic.NewEmptyResponse(s.cfg)
	res.SubsonicResponse.OpenSubSonicExtensions = &[]subsonic.OpenSubSonicExtension{}
	subsonic.WriteJSON(w, http.StatusOK, res)
}
