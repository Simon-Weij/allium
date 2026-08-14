package handlers

import (
	"testing"

	"github.com/Simon-Weij/allium/internal/metadata"
	"github.com/stretchr/testify/assert"
)

const (
	testReleaseDate    = "2021-08-13T12:00:00Z"
	testValidDate      = "2023-05-04T07:00:00Z"
	testCorruptedDate  = "20-20-20T03:30:30z"
	testIrrelevantDate = "wi9ofhujwuiofhwuifhwuifhjw"
	testGenre          = "Rock"
	testArtistName     = "Test Artist"
	testAlbumName      = "Test Album"
	testTrackName      = "Test Song"
	testTrackID        = 111222333
	testCollectionID   = 444555666
	testArtistID       = 777888999
	testArtworkURL     = "https://example.com"
	testAlbumArtURL    = "https://example.com"
	testArtistArtURL   = "https://example.com"
)

func TestParseYear(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		date     string
		expected int
	}{
		{
			name:     "should run correctly with itunes-like date",
			date:     testValidDate,
			expected: 2023,
		},
		{
			name:     "should return 0 on empty date",
			date:     "",
			expected: 0,
		},
		{
			name:     "should return 0 on corrupted date",
			date:     testCorruptedDate,
			expected: 0,
		},
		{
			name:     "should return 0 on empty with irrelevant data",
			date:     testIrrelevantDate,
			expected: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := parseYear(tt.date)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func baseITunesTrack(t *testing.T) metadata.ITunesResult {
	t.Helper()
	//nolint:exhaustruct
	return metadata.ITunesResult{
		WrapperType:      itunesWrapperTrack,
		Kind:             "song",
		TrackID:          testTrackID,
		CollectionID:     testCollectionID,
		ArtistID:         testArtistID,
		ArtistName:       testArtistName,
		CollectionName:   testAlbumName,
		TrackName:        testTrackName,
		ArtworkURL100:    testArtworkURL,
		ReleaseDate:      testReleaseDate,
		TrackNumber:      1,
		TrackTimeMillis:  200000,
		PrimaryGenreName: testGenre,
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
		Type:         songType,
		AlbumId:      "444555666",
		Album:        "Test Album",
		ArtistId:     "777888999",
		Artist:       "Test Artist",
		CoverArt:     "https://example.com",
		Duration:     200,
		BitRate:      songBitRate,
		BitDepth:     songBitDepth,
		SamplingRate: songSamplingRate,
		ChannelCount: songChannelCount,
		Track:        1,
		Year:         2021,
		Genre:        testGenre,
		Suffix:       songSuffix,
		ContentType:  songContentType,
		Path:         songPath,
	}
}

func baseITunesAlbum(t *testing.T) metadata.ITunesResult {
	t.Helper()
	//nolint:exhaustruct
	return metadata.ITunesResult{
		WrapperType:      itunesWrapperCollection,
		Kind:             "album",
		CollectionID:     testCollectionID,
		ArtistID:         testArtistID,
		ArtistName:       testArtistName,
		CollectionName:   testAlbumName,
		ArtworkURL100:    testAlbumArtURL,
		ReleaseDate:      testReleaseDate,
		TrackCount:       10,
		PrimaryGenreName: testGenre,
	}
}

func baseExpectedAlbum(t *testing.T) SearchResult3Album {
	t.Helper()

	return SearchResult3Album{
		Id:         "444555666",
		Name:       "Test Album",
		Artist:     "Test Artist",
		Year:       2021,
		CoverArt:   "https://example.com",
		Starred:    testReleaseDate,
		Duration:   albumDurationPlaceholder,
		PlayCount:  albumPlayCountPlaceholder,
		Played:     playedPlaceholderDate,
		Created:    testReleaseDate,
		ArtistId:   "777888999",
		UserRating: userRatingPlaceholder,
		SongCount:  10,
	}
}

func TestConvertItunesOpenSubsonic(t *testing.T) {
	t.Parallel()

	t.Run("successful conversion of a track", func(t *testing.T) {
		t.Parallel()
		input := []metadata.ITunesResult{
			baseITunesTrack(t),
		}

		got := ConvertItunesOpenSubsonic(input)

		assert.Len(t, got.songs, 1)
		assert.Equal(t, baseExpectedSong(t), got.songs[0])
	})

	t.Run("missing release date defaults year to 0", func(t *testing.T) {
		t.Parallel()
		track := baseITunesTrack(t)
		track.ReleaseDate = ""

		got := ConvertItunesOpenSubsonic([]metadata.ITunesResult{track})

		expected := baseExpectedSong(t)
		expected.Year = 0

		assert.Len(t, got.songs, 1)
		assert.Equal(t, expected, got.songs[0])
	})

	t.Run("zero track time defaults duration to 0", func(t *testing.T) {
		t.Parallel()
		track := baseITunesTrack(t)
		track.TrackTimeMillis = 0

		got := ConvertItunesOpenSubsonic([]metadata.ITunesResult{track})

		expected := baseExpectedSong(t)
		expected.Duration = 0

		assert.Len(t, got.songs, 1)
		assert.Equal(t, expected, got.songs[0])
	})

	t.Run("missing artwork url defaults cover art to empty", func(t *testing.T) {
		t.Parallel()
		track := baseITunesTrack(t)
		track.ArtworkURL100 = ""

		got := ConvertItunesOpenSubsonic([]metadata.ITunesResult{track})

		expected := baseExpectedSong(t)
		expected.CoverArt = ""

		assert.Len(t, got.songs, 1)
		assert.Equal(t, expected, got.songs[0])
	})

	t.Run("unknown wrapper type is ignored", func(t *testing.T) {
		t.Parallel()
		track := baseITunesTrack(t)
		track.WrapperType = "unknown-type"

		got := ConvertItunesOpenSubsonic([]metadata.ITunesResult{track})

		assert.Empty(t, got.songs)
		assert.Empty(t, got.albums)
		assert.Empty(t, got.artists)
	})

	t.Run("successful conversion of an album", func(t *testing.T) {
		t.Parallel()
		input := []metadata.ITunesResult{
			baseITunesAlbum(t),
		}

		got := ConvertItunesOpenSubsonic(input)

		assert.Len(t, got.albums, 1)
		assert.Equal(t, baseExpectedAlbum(t), got.albums[0])
	})

	t.Run("missing release date defaults album year to 0", func(t *testing.T) {
		t.Parallel()
		album := baseITunesAlbum(t)
		album.ReleaseDate = ""

		got := ConvertItunesOpenSubsonic([]metadata.ITunesResult{album})

		expected := baseExpectedAlbum(t)
		expected.Year = 0
		expected.Starred = ""
		expected.Created = ""

		assert.Len(t, got.albums, 1)
		assert.Equal(t, expected, got.albums[0])
	})

	t.Run("missing artwork url defaults album cover art to empty", func(t *testing.T) {
		t.Parallel()
		album := baseITunesAlbum(t)
		album.ArtworkURL100 = ""

		got := ConvertItunesOpenSubsonic([]metadata.ITunesResult{album})

		expected := baseExpectedAlbum(t)
		expected.CoverArt = ""

		assert.Len(t, got.albums, 1)
		assert.Equal(t, expected, got.albums[0])
	})

	t.Run("album with zero track count", func(t *testing.T) {
		t.Parallel()
		album := baseITunesAlbum(t)
		album.TrackCount = 0

		got := ConvertItunesOpenSubsonic([]metadata.ITunesResult{album})

		expected := baseExpectedAlbum(t)
		expected.SongCount = 0

		assert.Len(t, got.albums, 1)
		assert.Equal(t, expected, got.albums[0])
	})

	t.Run("artist is parsed correctly", func(t *testing.T) {
		t.Parallel()
		//nolint:exhaustruct
		artistInput := metadata.ITunesResult{
			WrapperType:      itunesWrapperArtist,
			Kind:             "artist",
			ArtistID:         123,
			ArtistName:       "Artist",
			ArtworkURL100:    testArtistArtURL,
			PrimaryGenreName: "Genre",
		}

		got := ConvertItunesOpenSubsonic([]metadata.ITunesResult{
			artistInput,
		})

		assert.Len(t, got.artists, 1)
		assert.Equal(t, "123", got.artists[0].Id)
		assert.Equal(t, "Artist", got.artists[0].Name)
		assert.Equal(t, "https://example.com", got.artists[0].CoverArt)
	})

	t.Run("empty input returns empty ResultTypes", func(t *testing.T) {
		t.Parallel()

		got := ConvertItunesOpenSubsonic([]metadata.ITunesResult{})

		assert.Empty(t, got.songs)
		assert.Empty(t, got.albums)
		assert.Empty(t, got.artists)
	})

	t.Run("nil input returns empty ResultTypes", func(t *testing.T) {
		t.Parallel()

		got := ConvertItunesOpenSubsonic(nil)

		assert.Empty(t, got.songs)
		assert.Empty(t, got.albums)
		assert.Empty(t, got.artists)
	})
}
