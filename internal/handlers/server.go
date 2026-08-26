package handlers

import (
	"github.com/Simon-Weij/allium/internal/config"
)

type Server struct {
	cfg            config.Config
	queries        Queries
	iTunesClient   iTunesClient
	songDownloader SongDownloader
}

func NewServer(
	cfg config.Config,
	queries Queries,
	iTunesClient iTunesClient,
	songDownloader SongDownloader,
) *Server {
	return &Server{
		cfg:            cfg,
		queries:        queries,
		iTunesClient:   iTunesClient,
		songDownloader: songDownloader,
	}
}
