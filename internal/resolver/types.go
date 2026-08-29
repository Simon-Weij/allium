package resolver

//go:generate mockgen -source=types.go -destination=../../generated/mocks/resolver/downloader_mock.go -package=resolver Downloader
type Downloader interface {
	DownloadAlbumCover(artworkURL, target string) error
}
