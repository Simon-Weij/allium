package metadata

import (
	"fmt"

	"resty.dev/v3"
)

type (
	Metadata struct {
		client *resty.Client
	}

	ITunesResponse struct {
		ResultCount int            `json:"resultCount"`
		Results     []ITunesResult `json:"results"`
	}

	ITunesResult struct {
		WrapperType             string  `json:"wrapperType"`
		Kind                    string  `json:"kind"`
		ArtistID                int     `json:"artistId"`
		CollectionID            int     `json:"collectionId"`
		TrackID                 int     `json:"trackId"`
		ArtistLinkURL           string  `json:"artistLinkUrl"`
		ArtistType              string  `json:"artistType"`
		ArtistName              string  `json:"artistName"`
		CollectionName          string  `json:"collectionName"`
		TrackName               string  `json:"trackName"`
		CollectionCensoredName  string  `json:"collectionCensoredName"`
		TrackCensoredName       string  `json:"trackCensoredName"`
		CollectionArtistID      int     `json:"collectionArtistId"`
		CollectionArtistName    string  `json:"collectionArtistName"`
		CollectionArtistViewURL string  `json:"collectionArtistViewUrl"`
		ArtistViewURL           string  `json:"artistViewUrl"`
		CollectionViewURL       string  `json:"collectionViewUrl"`
		TrackViewURL            string  `json:"trackViewUrl"`
		PreviewURL              string  `json:"previewUrl"`
		ArtworkURL30            string  `json:"artworkUrl30"`
		ArtworkURL60            string  `json:"artworkUrl60"`
		ArtworkURL100           string  `json:"artworkUrl100"`
		CollectionPrice         float64 `json:"collectionPrice"`
		TrackPrice              float64 `json:"trackPrice"`
		ReleaseDate             string  `json:"releaseDate"`
		CollectionExplicitness  string  `json:"collectionExplicitness"`
		TrackExplicitness       string  `json:"trackExplicitness"`
		DiscCount               int     `json:"discCount"`
		DiscNumber              int     `json:"discNumber"`
		TrackCount              int     `json:"trackCount"`
		TrackNumber             int     `json:"trackNumber"`
		TrackTimeMillis         int     `json:"trackTimeMillis"`
		Country                 string  `json:"country"`
		Currency                string  `json:"currency"`
		PrimaryGenreName        string  `json:"primaryGenreName"`
		IsStreamable            bool    `json:"isStreamable"`
	}
)

const (
	baseUrl = "https://itunes.apple.com/search"
)

func NewMetadata() *Metadata {
	return &Metadata{
		client: resty.New(),
	}
}

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
