package handlers

import "github.com/Simon-Weij/allium/internal/metadata"

//go:generate mockgen -source=types.go -destination=../../generated/mocks/itunes_mock.go -package=mocks
type iTunesClient interface {
	SearchWithItunes(query string) (*metadata.ITunesResponse, error)
	GetAlbumMetadata(albumId string) (*metadata.ITunesResponse, error)
	GetArtistById(artistId string) (*metadata.ITunesResponse, error)
}
