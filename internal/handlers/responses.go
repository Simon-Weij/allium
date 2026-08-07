package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Simon-Weij/allium/internal/config"
)

type (
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

func NewSubsonicResponse(cfg config.Config) SubsonicResponse {
	return SubsonicResponse{
		Body: SubsonicResponseBody{
			Status:        "ok",
			Version:       cfg.Version,
			Type:          cfg.Name,
			ServerVersion: cfg.Version,
			OpenSubsonic:  true,
		},
	}
}

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
