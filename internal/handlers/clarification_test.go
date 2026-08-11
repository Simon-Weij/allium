package handlers

import (
	"testing"

	"github.com/Simon-Weij/allium/internal/metadata"
	"github.com/stretchr/testify/assert"
)

func TestParseYear(t *testing.T) {
	tests := []struct {
		name     string
		date     string
		expected int
	}{
		{
			name:     "should run correctly with itunes-like date",
			date:     "2023-05-04T07:00:00Z",
			expected: 2023,
		},
		{
			name:     "should return 0 on empty date",
			date:     "",
			expected: 0,
		},
		{
			name:     "should return 0 on corrupted date",
			date:     "20-20-20T03:30:30z",
			expected: 0,
		},
		{
			name:     "should return 0 on empty with irrelevant data",
			date:     "wi9ofhujwuiofhwuifhwuifhjw",
			expected: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseYear(tt.date)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func baseITunesTrack(t *testing.T) metadata.ITunesResult {
	t.Helper()
	//nolint:exhaustruct
	return metadata.ITunesResult{
		WrapperType:      "track",
		Kind:             "song",
		TrackID:          111222333,
		CollectionID:     444555666,
		ArtistID:         777888999,
		ArtistName:       "Test Artist",
		CollectionName:   "Test Album",
		TrackName:        "Test Song",
		ArtworkURL100:    "https://example.com/default-100.jpg",
		ReleaseDate:      "2021-08-13T12:00:00Z",
		TrackNumber:      1,
		TrackTimeMillis:  200000,
		PrimaryGenreName: "Pop",
	}
}

func baseExpectedSong(t *testing.T) Song {
	t.Helper()
	//nolint:exhaustruct
	return Song{
		Id:           "111222333",
		Parent:       "444555666",
		Title:        "Test Song",
		IsDir:        false,
		IsVideo:      false,
		Type:         "music",
		AlbumId:      "444555666",
		Album:        "Test Album",
		ArtistId:     "777888999",
		Artist:       "Test Artist",
		CoverArt:     "https://example.com/default-100.jpg",
		Duration:     200,
		BitRate:      320,
		BitDepth:     24,
		SamplingRate: 48000,
		ChannelCount: 2,
		Track:        1,
		Year:         2021,
		Genre:        "Pop",
		Suffix:       "mp3",
		ContentType:  "audio/mpeg",
		Path:         "/home/alice/Music",
	}
}

func baseITunesAlbum(t *testing.T) metadata.ITunesResult {
	t.Helper()
	//nolint:exhaustruct
	return metadata.ITunesResult{
		WrapperType:      "collection",
		Kind:             "album",
		CollectionID:     444555666,
		ArtistID:         777888999,
		ArtistName:       "Test Artist",
		CollectionName:   "Test Album",
		ArtworkURL100:    "https://example.com/album-100.jpg",
		ReleaseDate:      "2021-08-13T12:00:00Z",
		TrackCount:       10,
		PrimaryGenreName: "Pop",
	}
}

func baseExpectedAlbum(t *testing.T) Album {
	t.Helper()
	//nolint:exhaustruct
	return Album{
		Id:         "444555666",
		Name:       "Test Album",
		Artist:     "Test Artist",
		Year:       2021,
		CoverArt:   "https://example.com/album-100.jpg",
		Starred:    "2021-08-13T12:00:00Z",
		Duration:   30000,
		PlayCount:  8,
		Played:     "2023-03-28T00:45:13Z",
		Created:    "2021-08-13T12:00:00Z",
		ArtistId:   "777888999",
		UserRating: 5,
		SongCount:  10,
	}
}

func TestConvertItunesOpenSubsonic(t *testing.T) {
	t.Run("successful conversion of a track", func(t *testing.T) {
		input := []metadata.ITunesResult{
			baseITunesTrack(t),
		}

		got := convertItunesOpenSubsonic(input)

		assert.Len(t, got.songs, 1)
		assert.Equal(t, baseExpectedSong(t), got.songs[0])
	})

	t.Run("missing release date defaults year to 0", func(t *testing.T) {
		track := baseITunesTrack(t)
		track.ReleaseDate = ""

		got := convertItunesOpenSubsonic([]metadata.ITunesResult{track})

		expected := baseExpectedSong(t)
		expected.Year = 0

		assert.Len(t, got.songs, 1)
		assert.Equal(t, expected, got.songs[0])
	})

	t.Run("zero track time defaults duration to 0", func(t *testing.T) {
		track := baseITunesTrack(t)
		track.TrackTimeMillis = 0

		got := convertItunesOpenSubsonic([]metadata.ITunesResult{track})

		expected := baseExpectedSong(t)
		expected.Duration = 0

		assert.Len(t, got.songs, 1)
		assert.Equal(t, expected, got.songs[0])
	})

	t.Run("missing artwork url defaults cover art to empty", func(t *testing.T) {
		track := baseITunesTrack(t)
		track.ArtworkURL100 = ""

		got := convertItunesOpenSubsonic([]metadata.ITunesResult{track})

		expected := baseExpectedSong(t)
		expected.CoverArt = ""

		assert.Len(t, got.songs, 1)
		assert.Equal(t, expected, got.songs[0])
	})

	t.Run("unknown wrapper type is ignored", func(t *testing.T) {
		track := baseITunesTrack(t)
		track.WrapperType = "unknown-type"

		got := convertItunesOpenSubsonic([]metadata.ITunesResult{track})

		assert.Empty(t, got.songs)
		assert.Empty(t, got.albums)
		assert.Empty(t, got.artists)
	})

	t.Run("successful conversion of an album", func(t *testing.T) {
		input := []metadata.ITunesResult{
			baseITunesAlbum(t),
		}

		got := convertItunesOpenSubsonic(input)

		assert.Len(t, got.albums, 1)
		assert.Equal(t, baseExpectedAlbum(t), got.albums[0])
	})

	t.Run("missing release date defaults album year to 0", func(t *testing.T) {
		album := baseITunesAlbum(t)
		album.ReleaseDate = ""

		got := convertItunesOpenSubsonic([]metadata.ITunesResult{album})

		expected := baseExpectedAlbum(t)
		expected.Year = 0
		expected.Starred = ""
		expected.Created = ""

		assert.Len(t, got.albums, 1)
		assert.Equal(t, expected, got.albums[0])
	})

	t.Run("missing artwork url defaults album cover art to empty", func(t *testing.T) {
		album := baseITunesAlbum(t)
		album.ArtworkURL100 = ""

		got := convertItunesOpenSubsonic([]metadata.ITunesResult{album})

		expected := baseExpectedAlbum(t)
		expected.CoverArt = ""

		assert.Len(t, got.albums, 1)
		assert.Equal(t, expected, got.albums[0])
	})

	t.Run("album with zero track count", func(t *testing.T) {
		album := baseITunesAlbum(t)
		album.TrackCount = 0

		got := convertItunesOpenSubsonic([]metadata.ITunesResult{album})

		expected := baseExpectedAlbum(t)
		expected.SongCount = 0

		assert.Len(t, got.albums, 1)
		assert.Equal(t, expected, got.albums[0])
	})

	t.Run("artist is parsed correctly", func(t *testing.T) {
		//nolint:exhaustruct
		artistInput := metadata.ITunesResult{
			WrapperType:      "artist",
			Kind:             "artist",
			ArtistID:         123,
			ArtistName:       "Artist",
			ArtworkURL100:    "https://example.com/artist-100.jpg",
			PrimaryGenreName: "Genre",
		}

		got := convertItunesOpenSubsonic([]metadata.ITunesResult{
			artistInput,
		})

		assert.Len(t, got.artists, 1)
		assert.Equal(t, "123", got.artists[0].Id)
		assert.Equal(t, "Artist", got.artists[0].Name)
		assert.Equal(t, "https://example.com/artist-100.jpg", got.artists[0].CoverArt)
	})

	t.Run("empty input returns empty ResultTypes", func(t *testing.T) {
		got := convertItunesOpenSubsonic([]metadata.ITunesResult{})

		assert.Empty(t, got.songs)
		assert.Empty(t, got.albums)
		assert.Empty(t, got.artists)
	})

	t.Run("nil input returns empty ResultTypes", func(t *testing.T) {
		got := convertItunesOpenSubsonic(nil)

		assert.Empty(t, got.songs)
		assert.Empty(t, got.albums)
		assert.Empty(t, got.artists)
	})
}
