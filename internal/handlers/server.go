package handlers

import (
	"github.com/Simon-Weij/allium/generated/sqlc"
	"github.com/Simon-Weij/allium/internal/config"
	"github.com/Simon-Weij/allium/internal/metadata"
	"github.com/Simon-Weij/allium/internal/users"
)

type Server struct {
	cfg          config.Config
	metadata     *metadata.Metadata
	queries      *sqlc.Queries
	iTunesclient iTunesClient
	Userclient users.UserManagementClient
}

func NewServer(
	cfg config.Config,
	queries *sqlc.Queries,
	iTunesClient iTunesClient,
	UserClient users.UserManagementClient,
) *Server {
	return &Server{
		cfg:          cfg,
		metadata:     metadata.NewMetadata(cfg),
		queries:      queries,
		iTunesclient: iTunesClient,
		Userclient: UserClient,
	}
}
