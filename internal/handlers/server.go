package handlers

import (
	"github.com/Simon-Weij/allium/internal/config"
	"github.com/Simon-Weij/allium/internal/metadata"
)

type Server struct {
	cfg      config.Config
	metadata *metadata.Metadata
}

func NewServer(cfg config.Config) *Server {
	return &Server{
		cfg:      cfg,
		metadata: metadata.NewMetadata(cfg),
	}
}
