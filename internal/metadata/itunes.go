package metadata

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type (
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

	downloader interface {
		downloadAlbumCover(artworkURL, target string) error
	}
)

const (
	baseSearchUrl = "https://itunes.apple.com/search"
	baseLookupUrl = "https://itunes.apple.com/lookup"
	coverArtSize  = "1600x1600bb.jpg"
)

const (
	directorySpacing = 2
	groupRestricted  = 0o750
)

var (
	errInvalidArtworkURL = errors.New("invalid artwork url")
	errCreatingDirs      = errors.New("couldn't create directories")
)

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
			Get(baseSearchUrl); err != nil {
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
		if err := m.downloader.downloadAlbumCover(id, coverPath); err != nil {
			return "", err
		}
	}

	return coverPath, nil
}

func (m Metadata) GetSongById(id string) (*ITunesResponse, error) {
	var res ITunesResponse
	if _, err := m.client.R().
		SetQueryParam("id", id).
		SetQueryParam("media", "music").
		SetQueryParam("entity", "song").
		SetResponseForceContentType("application/json").
		SetResult(&res).
		Get(baseLookupUrl); err != nil {
		return nil, fmt.Errorf("failed to get song by id: %s: %w", id, err)
	}

	return &res, nil
}

func (m Metadata) GetAlbumMetadata(albumId string) (*ITunesResponse, error) {
	var res ITunesResponse
	if _, err := m.client.R().
		SetQueryParam("id", albumId).
		SetQueryParam("media", "music").
		SetQueryParam("entity", "song").
		SetResponseForceContentType("application/json").
		SetResult(&res).
		Get(baseLookupUrl); err != nil {
		return nil, fmt.Errorf("failed to search %s: %w", albumId, err)
	}

	return &res, nil
}

func (m Metadata) downloadAlbumCover(artworkURL, target string) error {
	coverURL, err := coverArtURL(artworkURL)
	if err != nil {
		return err
	}

	_, err = m.client.R().
		SetResponseSaveToFile(true).
		SetResponseSaveFileName(target).
		Get(coverURL)
	if err != nil {
		return fmt.Errorf("could not cache album cover: %w", err)
	}

	return nil
}

func coverArtURL(artworkURL string) (string, error) {
	index := strings.LastIndex(artworkURL, "/")
	if index == -1 {
		return "", fmt.Errorf("%w: %s", errInvalidArtworkURL, artworkURL)
	}

	return artworkURL[:index+1] + coverArtSize, nil
}

// Create dirs according to a string where every 2 characters create a new directory.
// It returns the relative path of the deepest directory created.
func (m Metadata) createDirs(baseDirs, str string) (string, error) {
	var relativeDir string

	for i := 0; i < len(str); i += directorySpacing {
		end := min(i+directorySpacing, len(str))

		relativeDir = filepath.Join(relativeDir, str[i:end])
		if err := os.MkdirAll(filepath.Join(baseDirs, relativeDir), groupRestricted); err != nil {
			return "", fmt.Errorf("%w: %w", errCreatingDirs, err)
		}
	}

	return relativeDir, nil
}
