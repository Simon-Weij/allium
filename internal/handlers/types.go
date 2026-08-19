package handlers

import (
	"context"

	"github.com/Simon-Weij/allium/generated/sqlc"
	"github.com/Simon-Weij/allium/internal/metadata"
)

//go:generate mockgen -source=types.go -destination=../../generated/mocks/itunes_mock.go -package=mocks -exclude_interfaces=Queries
type iTunesClient interface {
	SearchWithItunes(query string) (*metadata.ITunesResponse, error)
	GetAlbumMetadata(albumId string) (*metadata.ITunesResponse, error)
	GetArtistById(artistId string) (*metadata.ITunesResponse, error)
	GetSongById(id string) (*metadata.ITunesResponse, error)
}

//go:generate mockgen -source=types.go -destination=../../generated/mocks/queries_mock.go -package=mocks -exclude_interfaces=iTunesClient
type Queries interface {
	UpdatePlays(ctx context.Context, arg sqlc.UpdatePlaysParams) error
}
