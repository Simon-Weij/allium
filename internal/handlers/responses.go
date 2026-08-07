package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Simon-Weij/allium/internal/config"
)

type (
	License struct {
		Valid          bool
		Email          string    `json:"email,omitempty"`
		LicenseExpires time.Time `json:"licenseExpires,omitempty"`
		TrialExpires   time.Time `json:"trialExpires,omitempty"`
	}

	LicenseResponse struct {
		SubsonicResponseBody
		License License
	}

	SubsonicResponseBody struct {
		Status        string `json:"status"`
		Version       string `json:"version"`
		Type          string `json:"type"`
		ServerVersion string `json:"serverVersion"`
		OpenSubsonic  bool   `json:"openSubsonic"`
	}

	SubsonicResponse struct {
		Body SubsonicResponseBody `json:"subsonic-response"`
	}

	Generic401Response struct{}
)

func NewSubsonicResponse(cfg config.Config) SubsonicResponseBody {
	return SubsonicResponseBody{
		Status:        "ok",
		Version:       cfg.Version,
		Type:          cfg.Name,
		ServerVersion: cfg.Version,
		OpenSubsonic:  true,
	}
}

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
