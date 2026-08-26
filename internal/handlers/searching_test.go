package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Simon-Weij/allium/generated/mocks"
	"github.com/Simon-Weij/allium/internal/metadata"
	"github.com/Simon-Weij/allium/internal/subsonic"
	"github.com/Simon-Weij/allium/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestConvertItunesOpenSubsonic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		results         []metadata.ITunesResult
		expectedSongs   []subsonic.Song
		expectedAlbums  []subsonic.SearchResult3Album
		expectedArtists []subsonic.Artist
	}{
		{
			name:    "should return empty slice for no results",
			results: []metadata.ITunesResult{},
		},
		{
			name:    "should convert track to song",
			results: []metadata.ITunesResult{mockResult},
			expectedSongs: []subsonic.Song{
				{
					Id:           "777888999",
					Parent:       "444555666",
					Title:        "Test Song",
					IsDir:        false,
					IsVideo:      false,
					Type:         "music",
					AlbumId:      "444555666",
					Album:        "Test album",
					ArtistId:     "111222333",
					Artist:       "Alice",
					CoverArt:     "https://example.com/artwork/100x100.jpg",
					Duration:     215,
					BitRate:      320,
					BitDepth:     24,
					SamplingRate: 48000,
					ChannelCount: 2,
					Track:        3,
					Year:         2024,
					Genre:        "Pop",
					DiscNumber:   1,
					Suffix:       "mp3",
					ContentType:  "audio/mpeg",
					Path:         "/home/alice/Music",
				},
			},
		},
		{
			name: "should convert collection to album",
			results: []metadata.ITunesResult{
				{
					WrapperType:    "collection",
					ArtistID:       111222333,
					CollectionID:   444555666,
					ArtistName:     "Alice",
					CollectionName: "Test album",
					ArtworkURL100:  "https://example.com/artwork/100x100.jpg",
					ReleaseDate:    "2024-01-15T12:00:00Z",
					TrackCount:     10,
				},
			},
			expectedAlbums: []subsonic.SearchResult3Album{
				{
					Id:         "444555666",
					Name:       "Test album",
					Artist:     "Alice",
					Year:       2024,
					CoverArt:   "https://example.com/artwork/100x100.jpg",
					Starred:    "2024-01-15T12:00:00Z",
					Duration:   30000,
					PlayCount:  8,
					Played:     "2023-03-28T00:45:13Z",
					Created:    "2024-01-15T12:00:00Z",
					ArtistId:   "111222333",
					UserRating: 5,
					SongCount:  10,
				},
			},
		},
		{
			name: "should convert artist",
			results: []metadata.ITunesResult{
				{
					WrapperType:   "artist",
					ArtistID:      111222333,
					ArtistName:    "Alice",
					ArtworkURL100: "https://example.com/artwork/100x100.jpg",
				},
			},
			expectedArtists: []subsonic.Artist{
				{
					Id:             "111222333",
					Name:           "Alice",
					CoverArt:       "https://example.com/artwork/100x100.jpg",
					AlbumCount:     5,
					UserRating:     5,
					ArtistImageUrl: "https://example.com/artwork/100x100.jpg",
				},
			},
		},
		{
			name: "should skip unknown wrapper type",
			results: []metadata.ITunesResult{
				{
					WrapperType: "unknown",
					ArtistID:    111222333,
					TrackID:     777888999,
				},
			},
			expectedSongs:   []subsonic.Song{},
			expectedAlbums:  []subsonic.SearchResult3Album{},
			expectedArtists: []subsonic.Artist{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := ConvertItunesOpenSubsonic(tt.results)

			assert.ElementsMatch(t, tt.expectedSongs, result.songs)
			assert.ElementsMatch(t, tt.expectedAlbums, result.albums)
			assert.ElementsMatch(t, tt.expectedArtists, result.artists)
		})
	}
}

func TestHandleSearch3(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		query             string
		expectedCode      int
		expectedSongCount int
		setupMock         func(m *mocks.MockiTunesClient)
	}{
		{
			name:              "should run correctly",
			query:             "?query=song",
			expectedCode:      http.StatusOK,
			expectedSongCount: 1,
			setupMock: func(m *mocks.MockiTunesClient) {
				m.EXPECT().
					SearchWithItunes("song").
					Return(&metadata.ITunesResponse{
						ResultCount: 1,
						Results: []metadata.ITunesResult{
							mockResult,
						},
					}, nil)
			},
		},
		{
			name:              "should default to limit when songCount is zero",
			query:             "?query=song&songCount=0",
			expectedCode:      http.StatusOK,
			expectedSongCount: 20,
			setupMock: func(m *mocks.MockiTunesClient) {
				results := make([]metadata.ITunesResult, 30)
				for i := range results {
					results[i] = mockResult
				}

				m.EXPECT().
					SearchWithItunes("song").
					Return(&metadata.ITunesResponse{
						ResultCount: 30,
						Results:     results,
					}, nil)
			},
		},
		{
			name:         "should 403 on invalid artistCount",
			query:        "?query=song&artistCount=notanumber",
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockiTunesClient) {},
		},
		{
			name:         "should 403 on invalid albumCount",
			query:        "?query=song&albumCount=notanumber",
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockiTunesClient) {},
		},
		{
			name:         "should 403 on invalid searchCount",
			query:        "?query=song&searchCount=notanumber",
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockiTunesClient) {},
		},
		{
			name:         "should 403 on no query",
			query:        "",
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockiTunesClient) {},
		},
		{
			name:         "should 403 on invalid soungCount",
			query:        "?query=song&songCount=notanumber",
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockiTunesClient) {},
		},
		{
			name:         "should 500 when searching with itunes goes wrong",
			query:        "?query=song",
			expectedCode: http.StatusInternalServerError,
			setupMock: func(m *mocks.MockiTunesClient) {
				m.EXPECT().
					SearchWithItunes("song").
					Return(nil, errors.New("something went wrong"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := testutil.SetupTestingConfig(t)

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockiTunesClient(ctrl)
			tt.setupMock(mockClient)

			server := NewServer(cfg, nil, mockClient, nil)

			req := httptest.NewRequest(http.MethodGet, "/rest/search3.view"+tt.query, nil)
			rec := httptest.NewRecorder()

			server.HandleSearch3(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			if tt.expectedSongCount > 0 {
				var res subsonic.SubsonicResponseWrapper

				err := json.NewDecoder(rec.Body).Decode(&res)
				require.NoError(t, err)

				assert.Len(t, res.SubsonicResponse.SearchResult3.Song, tt.expectedSongCount)
			}
		})
	}
}
