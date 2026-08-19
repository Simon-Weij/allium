package handlers

import (
	"github.com/Simon-Weij/allium/internal/config"
	"github.com/Simon-Weij/allium/internal/metadata"
)

type Server struct {
	cfg          config.Config
	metadata     *metadata.Metadata
	queries      Queries
	iTunesclient iTunesClient
}

func NewServer(
	cfg config.Config,
	queries Queries,
	iTunesClient iTunesClient,
) *Server {
	return &Server{
		cfg:          cfg,
		metadata:     metadata.NewMetadata(cfg),
		queries:      queries,
		iTunesclient: iTunesClient,
	}
}
