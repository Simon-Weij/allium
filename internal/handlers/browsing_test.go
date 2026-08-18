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

var mockResult = metadata.ITunesResult{
	WrapperType:      "track",
	Kind:             "song",
	ArtistID:         111222333,
	CollectionID:     444555666,
	TrackID:          777888999,
	ArtistName:       "Alice",
	CollectionName:   "Test album",
	TrackName:        "Test Song",
	ArtworkURL100:    "https://example.com/artwork/100x100.jpg",
	ReleaseDate:      "2024-01-15T12:00:00Z",
	DiscNumber:       1,
	TrackCount:       10,
	TrackNumber:      3,
	TrackTimeMillis:  215000,
	PrimaryGenreName: "Pop",
}

func TestHandleGetAlbum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		query             string
		expectedCode      int
		expectedSongCount int
		expectedDuration  int
		setupMock         func(m *mocks.MockiTunesClient)
	}{
		{
			name:  "should run as expected",
			query: "?id=5",
			setupMock: func(m *mocks.MockiTunesClient) {
				m.EXPECT().
					GetAlbumMetadata("5").
					Return(&metadata.ITunesResponse{
						ResultCount: 1,
						Results: []metadata.ITunesResult{
							mockResult,
						},
					}, nil)
			},
			expectedCode:      http.StatusOK,
			expectedSongCount: 1,
			expectedDuration:  215,
		},
		{
			name:  "should skip non-track results",
			query: "?id=5",
			setupMock: func(m *mocks.MockiTunesClient) {
				m.EXPECT().
					GetAlbumMetadata("5").
					Return(&metadata.ITunesResponse{
						ResultCount: 3,
						Results: []metadata.ITunesResult{
							mockResult,
							{
								WrapperType:    "collection",
								CollectionID:   9999,
								CollectionName: "Test album",
								TrackCount:     12,
							},
							mockResult,
						},
					}, nil)
			},
			expectedCode:      http.StatusOK,
			expectedSongCount: 2,
			expectedDuration:  430,
		},
		{
			name:         "should return 404 when no id provided",
			query:        "",
			setupMock:    func(m *mocks.MockiTunesClient) {},
			expectedCode: http.StatusNotFound,
		},
		{
			name:  "GetAlbumMetadata fails",
			query: "?id=5",
			setupMock: func(m *mocks.MockiTunesClient) {
				m.EXPECT().
					GetAlbumMetadata("5").
					Return(nil, errors.New("something went wrong"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:  "No results",
			query: "?id=5",
			setupMock: func(m *mocks.MockiTunesClient) {
				m.EXPECT().
					GetAlbumMetadata("5").
					Return(&metadata.ITunesResponse{
						ResultCount: 1,
						Results:     []metadata.ITunesResult{},
					}, nil)
			},
			expectedCode: http.StatusNotFound,
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

			req := httptest.NewRequest(http.MethodGet, "/rest/getAlbum"+tt.query, nil)
			rec := httptest.NewRecorder()

			server.HandleGetAlbum(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			if tt.expectedSongCount > 0 {
				var res subsonic.SubsonicResponseWrapper

				err := json.NewDecoder(rec.Body).Decode(&res)
				require.NoError(t, err)

				assert.Len(t, res.SubsonicResponse.Album.Song, tt.expectedSongCount)
				assert.Equal(t, tt.expectedDuration, res.SubsonicResponse.Album.Duration)
			}
		})
	}
}

func TestHandleGetArtist(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		query              string
		expectedCode       int
		expectedAlbumCount int
		setupMock          func(m *mocks.MockiTunesClient)
	}{
		{
			name:         "should run as expected",
			query:        "?id=5",
			expectedCode: http.StatusOK,
			setupMock: func(m *mocks.MockiTunesClient) {
				m.EXPECT().
					GetArtistById("5").
					Return(&metadata.ITunesResponse{
						ResultCount: 1,
						Results: []metadata.ITunesResult{
							{
								WrapperType:   "artist",
								ArtistID:      111222333,
								ArtistName:    "Alice",
								ArtworkURL100: "https://example.com/artwork/100x100.jpg",
							},
						},
					}, nil)
			},
			expectedAlbumCount: 0,
		},
		{
			name:  "should collect albums and skip artist and unknown types",
			query: "?id=5",
			setupMock: func(m *mocks.MockiTunesClient) {
				m.EXPECT().
					GetArtistById("5").
					Return(&metadata.ITunesResponse{
						ResultCount: 4,
						Results: []metadata.ITunesResult{
							{
								WrapperType:   "artist",
								ArtistID:      111222333,
								ArtistName:    "Alice",
								ArtworkURL100: "https://example.com/artwork/100x100.jpg",
							},
							{
								WrapperType:      "collection",
								ArtistID:         111222333,
								CollectionID:     444555666,
								CollectionName:   "Test album",
								ArtworkURL100:    "https://example.com/artwork/100x100.jpg",
								TrackCount:       10,
								ReleaseDate:      "2024-01-15T12:00:00Z",
								PrimaryGenreName: "Pop",
							},
							{
								WrapperType: "unknown",
								ArtistID:    111222333,
							},
							{
								WrapperType:      "collection",
								ArtistID:         111222333,
								CollectionID:     555666777,
								CollectionName:   "Second album",
								ArtworkURL100:    "https://example.com/artwork/100x100.jpg",
								TrackCount:       5,
								ReleaseDate:      "2023-01-15T12:00:00Z",
								PrimaryGenreName: "Rock",
							},
						},
					}, nil)
			},
			expectedCode:       http.StatusOK,
			expectedAlbumCount: 2,
		},
		{
			name:         "no id should return 404",
			query:        "",
			expectedCode: http.StatusNotFound,
			setupMock:    func(m *mocks.MockiTunesClient) {},
		},
		{
			name:         "GetArtistById fails",
			query:        "?id=5",
			expectedCode: http.StatusInternalServerError,
			setupMock: func(m *mocks.MockiTunesClient) {
				m.EXPECT().
					GetArtistById("5").
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

			req := httptest.NewRequest(http.MethodGet, "/rest/getArtist"+tt.query, nil)
			rec := httptest.NewRecorder()

			server.HandleGetArtist(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			if tt.expectedCode == http.StatusOK {
				var res subsonic.SubsonicResponseWrapper

				err := json.NewDecoder(rec.Body).Decode(&res)
				require.NoError(t, err)

				assert.Len(t, res.SubsonicResponse.Artist.Album, tt.expectedAlbumCount)
			}
		})
	}
}
