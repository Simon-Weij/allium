package resolver

import "fmt"

type SpotifyIDResponse struct {
	SpotifyTrackIDs []string `json:"spotify_track_ids"`
}

// GetSpotifyIdByMetadata exists we can get the spotify id, so reccobeats can eventually recommend similar songs.
func (r Resolver) GetSpotifyIdByMetadata(
	artistName string,
	albumName string,
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
