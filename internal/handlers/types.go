package handlers

import (
	"context"

	"github.com/Simon-Weij/allium/generated/sqlc"
	"github.com/Simon-Weij/allium/internal/metadata"
)

//go:generate mockgen -source=types.go -destination=../../generated/mocks/itunes_mock.go -package=mocks -exclude_interfaces=Queries,SongDownloader
type iTunesClient interface {
	SearchWithItunes(query string) (*metadata.ITunesResponse, error)
	GetAlbumMetadata(albumId string) (*metadata.ITunesResponse, error)
	GetArtistById(artistId string) (*metadata.ITunesResponse, error)
	GetSongById(id string) (*metadata.ITunesResponse, error)
	GetAlbumCover(id string) (string, error)
}

//go:generate mockgen -source=types.go -destination=../../generated/mocks/queries_mock.go -package=mocks -exclude_interfaces=iTunesClient,SongDownloader
type Queries interface {
	UpdatePlays(ctx context.Context, arg sqlc.UpdatePlaysParams) error
	CreatePlaylist(ctx context.Context, arg sqlc.CreatePlaylistParams) (sqlc.Playlist, error)
	UpdatePlaylistTitle(ctx context.Context, arg sqlc.UpdatePlaylistTitleParams) (sqlc.Playlist, error)
	UpdatePlaylistMetadata(ctx context.Context, arg sqlc.UpdatePlaylistMetadataParams) (sqlc.Playlist, error)
	InsertIntoPlaylist(ctx context.Context, arg sqlc.InsertIntoPlaylistParams) (sqlc.PlaylistSong, error)
	DeleteSongFromPlaylistByPosition(ctx context.Context, arg sqlc.DeleteSongFromPlaylistByPositionParams) error
	GetMaxPlaylistPosition(ctx context.Context, playlistID string) (int64, error)
	GetPlaylistsByUser(ctx context.Context, user string) ([]sqlc.GetPlaylistsByUserRow, error)
	GetPlaylistByID(ctx context.Context, id int64) (sqlc.Playlist, error)
	GetSongsByPlaylistID(ctx context.Context, playlistID string) ([]sqlc.Song, error)
}

//go:generate mockgen -source=types.go -destination=../../generated/mocks/song_downloader_mock.go -package=mocks -exclude_interfaces=iTunesClient,Queries
type SongDownloader interface {
	DownloadOrGetSong(ctx context.Context, artist, title string) (string, error)
}
