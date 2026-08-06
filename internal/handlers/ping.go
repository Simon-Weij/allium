package handlers

import (
	"net/http"

	"github.com/Simon-Weij/allium/internal/config"
)

type (
	PingHandler struct {
		cfg config.Config
	}
)

func NewPingHandler(cfg config.Config) PingHandler {
	return PingHandler{
		cfg: cfg,
	}
}

func (h PingHandler) HandlePing(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, NewSubsonicResponse(h.cfg))
}
