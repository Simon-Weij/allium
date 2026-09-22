package resolver

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Simon-Weij/allium/internal/subsonic"
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
)

const (
	baseSearchUrl             = "https://itunes.apple.com/search"
	baseLookupUrl             = "https://itunes.apple.com/lookup"
	coverArtSize              = "1600x1600bb.jpg"
	albumDurationPlaceholder  = 30000
	albumPlayCountPlaceholder = 8
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

	final := &ITunesResponse{}

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

func (m Metadata) SearchAlbums(size, offset int, sortByName bool) ([]subsonic.AlbumID3, error) {
	var res ITunesResponse
	if _, err := m.client.R().
		SetQueryParam("term", randomSearchTerm()).
		SetQueryParam("media", "music").
		SetQueryParam("entity", "album").
		SetQueryParam("limit", strconv.Itoa(offset+size+1)).
		SetResponseForceContentType("application/json").
		SetResult(&res).
		Get(baseSearchUrl); err != nil {
		return nil, fmt.Errorf("failed to fetch albums: %w", err)
	}

	if sortByName {
		sort.Slice(res.Results, func(i, j int) bool {
			return strings.ToLower(res.Results[i].CollectionName) < strings.ToLower(res.Results[j].CollectionName)
		})
	}

	return convertAlbumResults(res.Results, offset, size), nil
}

func convertAlbumResults(results []ITunesResult, offset, size int) []subsonic.AlbumID3 {
	start := min(offset, len(results))
	end := min(start+size, len(results))

	albums := make([]subsonic.AlbumID3, 0, end-start)
	for _, result := range results[start:end] {
		albums = append(albums, convertItunesAlbumID3(result))
	}

	return albums
}

func randomSearchTerm() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz"

	return string(alphabet[rand.IntN(len(alphabet))])
}

func convertItunesAlbumID3(result ITunesResult) subsonic.AlbumID3 {
	return subsonic.AlbumID3{
		Id:        strconv.Itoa(result.CollectionID),
		Name:      result.CollectionName,
		Artist:    result.ArtistName,
		Year:      parseDateYear(result.ReleaseDate),
		CoverArt:  result.ArtworkURL100,
		Starred:   result.ReleaseDate,
		Duration:  albumDurationPlaceholder,
		PlayCount: albumPlayCountPlaceholder,
		Genre:     result.PrimaryGenreName,
		Created:   result.ReleaseDate,
		ArtistId:  strconv.Itoa(result.ArtistID),
		SongCount: result.TrackCount,
	}
}

func parseDate(dateString string) (*time.Time, error) {
	time, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		return nil, fmt.Errorf("could not get time: %w", err)
	}

	return &time, nil
}

func parseDateYear(dateString string) int {
	t, err := parseDate(dateString)
	if err != nil {
		return 0
	}

	return t.Year()
}

func (m Metadata) GetAlbumCover(id string) (string, error) {
	dataDir := filepath.Join(m.cfg.Data, "covers")

	relativeDir, err := m.createDirs(dataDir, id)
	if err != nil {
		return "", fmt.Errorf("could not create directories: %w", err)
	}

	coverDir := filepath.Join(dataDir, relativeDir)
	coverPath := filepath.Join(coverDir, "cover.jpg")

	if _, err := os.Stat(coverPath); errors.Is(err, os.ErrNotExist) {
		if err := m.downloader.DownloadAlbumCover(id, coverPath); err != nil {
			return "", fmt.Errorf("could not download album cover: %w", err)
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

func (m Metadata) GetArtistById(id string) (*ITunesResponse, error) {
	var res ITunesResponse
	if _, err := m.client.R().
		SetQueryParam("id", id).
		SetQueryParam("media", "music").
		SetQueryParam("entity", "album").
		SetResponseForceContentType("application/json").
		SetResult(&res).
		Get(baseLookupUrl); err != nil {
		return nil, fmt.Errorf("failed to get artist by id: %s %w", id, err)
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

func (m Metadata) DownloadAlbumCover(artworkURL, target string) error {
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
