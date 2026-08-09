package handlers

import "github.com/Simon-Weij/allium/internal/config"

type Server struct {
	cfg config.Config
}

func NewServer(cfg config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}
