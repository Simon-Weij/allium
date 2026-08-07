package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Simon-Weij/allium/internal/config"
)

type (
	SubsonicResponseWrapper struct {
		SubsonicResponse Subsonic `json:"subsonic-response"`
	}

	Subsonic struct {
		Status        string        `json:"status"`
		Version       string        `json:"version"`
		Type          string        `json:"type"`
		ServerVersion string        `json:"serverVersion"`
		OpenSubsonic  bool          `json:"openSubsonic"`
		License       *License      `json:"license,omitempty"`
		MusicFolders  *MusicFolders `json:"musicFolders,omitempty"`
	}

	License struct {
		Valid          bool       `json:"valid"`
		Email          string     `json:"email,omitempty"`
		LicenseExpires *time.Time `json:"licenseExpires,omitempty"`
		TrialExpires   *time.Time `json:"trialExpires,omitempty"`
	}

	MusicFolders struct {
		MusicFolder []MusicFolder `json:"musicFolder"`
	}

	MusicFolder struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	Generic401Response struct{}
)

func NewEmptyResponse(cfg config.Config) SubsonicResponseWrapper {
	return SubsonicResponseWrapper{
		SubsonicResponse: Subsonic{
			Status:        "ok",
			Version:       cfg.Version,
			Type:          cfg.Name,
			ServerVersion: cfg.Version,
			OpenSubsonic:  true,

			License:      nil,
			MusicFolders: nil,
		},
	}
}

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
