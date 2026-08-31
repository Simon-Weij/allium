package resolver

import (
	"fmt"
	"net/url"
	"strconv"
)

type (
	SpotifyIDResponse struct {
		SpotifyTrackIDs []string `json:"spotify_track_ids"`
	}
	ListenBrainzWrapper struct {
		Content []ListenBrainzSongs `json:"content"`
	}
	ListenBrainzSongs struct {
		Id                 string                `json:"id"`
		TrackTitle         string                `json:"string"`
		Artists            []ListenBrainzArtists `json:"artists"`
		DurationMs         int                   `json:"durationMs"`
		Isrc               string                `json:"isrc"`
		Ean                string                `json:"ean"`
		Upc                string                `json:"upc"`
		Href               string                `json:"href"`
		AvaliableCountries string                `json:"availableCountries"`
		Popularity         int                   `json:"popularity"`
	}
	ListenBrainzArtists struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Href string `json:"href"`
	}
)

// GetSpotifyIdByMetadata exists we can get the spotify id, so reccobeats can eventually recommend similar songs.
func (r Resolver) GetSpotifyIdByMetadata(
	artistName,
	albumName,
	title string,
) (string, error) {
	var result SpotifyIDResponse

	_, err := r.client.R().
		SetQueryParams(map[string]string{
			"artist_name":  artistName,
			"release_name": albumName,
			"track_name":   title,
		}).
		SetResult(&result).
		Get("https://labs.api.listenbrainz.org/spotify-id-from-metadata/json")
	if err != nil {
		return "", fmt.Errorf("could not get spotify id: %w", err)
	}

	return result.SpotifyTrackIDs[0], nil
}

func (r Resolver) GetRecommendationsFromReccobeats(seeds []string, size int) error {
	var result ListenBrainzWrapper

	params := url.Values{}
	for _, seed := range seeds {
		params.Add("seeds", seed)
	}

	params.Add("size", strconv.Itoa(size))

	_, err := r.client.R().
		SetQueryParamsFromValues(params).
		SetResult(&result).
		Get("https://api.reccobeats.com/v1/track/recommendation")
	if err != nil {
		return fmt.Errorf("could not get recommendations: %w", err)
	}

	return nil
}
