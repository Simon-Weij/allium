package handlers

import (
	"github.com/Simon-Weij/allium/generated/sqlc"
	"github.com/Simon-Weij/allium/internal/config"
	"github.com/Simon-Weij/allium/internal/metadata"
)

type Server struct {
	cfg          config.Config
	metadata     *metadata.Metadata
	queries      *sqlc.Queries
	iTunesclient iTunesClient
}

func NewServer(
	cfg config.Config,
	queries *sqlc.Queries,
	iTunesClient iTunesClient,
) *Server {
	return &Server{
		cfg:          cfg,
		metadata:     metadata.NewMetadata(cfg),
		queries:      queries,
		iTunesclient: iTunesClient,
	}
}
