package metadata

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Simon-Weij/allium/internal/config"
	"resty.dev/v3"
)

type (
	Metadata struct {
		client *resty.Client
		cfg    config.Config
	}

	ITunesResponse struct {
		ResultCount int            `json:"resultCount"`
		Results     []ITunesResult `json:"results"`
	}

	ITunesResult struct {
		WrapperType      string `json:"wrapperType"`
		Kind             string `json:"kind"`
		ArtistID         int    `json:"artistId"`
		CollectionID     int    `json:"collectionId"`
		TrackID          int    `json:"trackId"`
		ArtistName       string `json:"artistName"`
		CollectionName   string `json:"collectionName"`
		TrackName        string `json:"trackName"`
		ArtworkURL100    string `json:"artworkUrl100"`
		ReleaseDate      string `json:"releaseDate"`
		DiscNumber       int    `json:"discNumber"`
		TrackCount       int    `json:"trackCount"`
		TrackNumber      int    `json:"trackNumber"`
		TrackTimeMillis  int    `json:"trackTimeMillis"`
		PrimaryGenreName string `json:"primaryGenreName"`
	}
)

const (
	baseUrl = "https://itunes.apple.com/search"
)

const (
	directorySpacing = 2
	groupRestricted  = 0o750
)

func NewMetadata(cfg config.Config) *Metadata {
	return &Metadata{
		client: resty.New(),
		cfg:    cfg,
	}
}

// SearchWithItunes takes an argument of query, which is a search query returned by the client
// It returns the iTunes response
func (m Metadata) SearchWithItunes(query string) (*ITunesResponse, error) {
	entities := []string{"song", "album", "musicArtist"}

	final := &ITunesResponse{
		ResultCount: 0,
		Results:     nil,
	}

	for _, entity := range entities {
		var res ITunesResponse
		if _, err := m.client.R().
			SetQueryParam("term", query).
			SetQueryParam("media", "music").
			SetQueryParam("entity", entity).
			SetResponseForceContentType("application/json").
			SetResult(&res).
			Get(baseUrl); err != nil {
			return nil, fmt.Errorf("failed to search %s: %w", entity, err)
		}

		final.Results = append(final.Results, res.Results...)
		final.ResultCount += res.ResultCount
	}

	return final, nil
}

// GetAlbumCover gets its id from the itunes url for getting album covers.
func (m Metadata) GetAlbumCover(id string) (string, error) {
	dataDir := filepath.Join(m.cfg.Data, "covers")

	relativeDir, err := m.createDirs(dataDir, id)
	if err != nil {
		return "", fmt.Errorf("could not create directories: %w", err)
	}

	coverDir := filepath.Join(dataDir, relativeDir)
	coverPath := filepath.Join(coverDir, "cover.jpg")

	if _, err := os.Stat(coverPath); errors.Is(err, os.ErrNotExist) {
		if err := m.downloadAlbumCover(id, coverPath); err != nil {
			return "", err
		}
	}

	return coverPath, nil
}

// downloadAlbumCover downloads the album cover from the artworkURL,
// the downloaded cover is saved to the target destination
func (m Metadata) downloadAlbumCover(artworkURL, target string) error {
	_, err := m.client.R().
		SetResponseSaveToFile(true).
		SetResponseSaveFileName(target).
		Get(artworkURL + "/1600x1600bb.jpg")
	if err != nil {
		return fmt.Errorf("could not cache album cover: %w", err)
	}

	return nil
}

// Create dirs according to a string where every 2 characters create a new directory.
// It returns the relative path of the deepest directory created.
func (m Metadata) createDirs(baseDirs, str string) (string, error) {
	var relativeDir string

	for i := 0; i < len(str); i += directorySpacing {
		end := min(i+directorySpacing, len(str))

		relativeDir = filepath.Join(relativeDir, str[i:end])
		if err := os.MkdirAll(filepath.Join(baseDirs, relativeDir), groupRestricted); err != nil {
			return "", fmt.Errorf("something went wrong creating directories: %w", err)
		}
	}

	return relativeDir, nil
}
