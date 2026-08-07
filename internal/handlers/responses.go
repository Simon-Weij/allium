package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Simon-Weij/allium/internal/config"
)

type (
	Subsonic struct {
		Status        string   `json:"status"`
		Version       string   `json:"version"`
		Type          string   `json:"type"`
		ServerVersion string   `json:"serverVersion"`
		OpenSubsonic  bool     `json:"openSubsonic"`
		License       *License `json:"license,omitempty"`
	}
	License struct {
		Valid          bool       `json:"valid"`
		Email          string     `json:"email,omitempty"`
		LicenseExpires *time.Time `json:"licenseExpires,omitempty"`
		TrialExpires   *time.Time `json:"trialExpires,omitempty"`
	}
	Generic401Response struct{}
)

func NewEmptyResponse(cfg config.Config) Subsonic {
	return Subsonic{
		Status:        "ok",
		Version:       cfg.Version,
		Type:          cfg.Name,
		ServerVersion: cfg.Version,
		OpenSubsonic:  true,

		License: nil,
	}
}

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
