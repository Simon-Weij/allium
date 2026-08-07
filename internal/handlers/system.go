package handlers

import (
	"net/http"
	"time"

	"github.com/Simon-Weij/allium/internal/config"
)

type (
	SystemHandler struct {
		cfg config.Config
	}
)

func NewSystemHandler(cfg config.Config) SystemHandler {
	return SystemHandler{
		cfg: cfg,
	}
}

func (h SystemHandler) HandlePing(w http.ResponseWriter, r *http.Request) {
	res := NewEmptyResponse(h.cfg)
	WriteJSON(w, http.StatusOK, res)
}

func (h SystemHandler) HandleGetLicense(w http.ResponseWriter, r *http.Request) {
	licenseExpiry := time.Date(2222, 1, 1, 1, 1, 1, 1, time.UTC)
	res := NewEmptyResponse(h.cfg)
	res.SubsonicResponse.License = &License{
		Valid:          true,
		Email:          "irrelevant@example.com",
		LicenseExpires: &licenseExpiry,
		TrialExpires:   &licenseExpiry,
	}

	WriteJSON(w, http.StatusOK, res)
}
