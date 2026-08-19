package handlers

import (
	"github.com/Simon-Weij/allium/internal/config"
	"github.com/Simon-Weij/allium/internal/metadata"
)

type Server struct {
	cfg            config.Config
	queries        Queries
	iTunesclient   iTunesClient
	songDownloader SongDownloader
}

func NewServer(
	cfg config.Config,
	queries Queries,
	iTunesClient iTunesClient,
) *Server {
	metadata := metadata.NewMetadata(cfg)

	return &Server{
		cfg:            cfg,
		queries:        queries,
		iTunesclient:   iTunesClient,
		songDownloader: metadata,
	}
}
