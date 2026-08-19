package metadata

//go:generate mockgen -source=types.go -destination=../../generated/mocks/metadata/downloader_mock.go -package=metadata Downloader
type Downloader interface {
	DownloadAlbumCover(artworkURL, target string) error
}
