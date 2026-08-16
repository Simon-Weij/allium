package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Simon-Weij/allium/internal/handlers/mocks"
	"github.com/Simon-Weij/allium/internal/metadata"
	"github.com/Simon-Weij/allium/internal/testutil"
	"github.com/stretchr/testify/assert"
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
		name         string
		query        string
		expectedCode int
		setupMock    func(m *mocks.MockiTunesClient)
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
			expectedCode: http.StatusOK,
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

			server := NewServer(cfg, nil, mockClient)

			req := httptest.NewRequest(http.MethodGet, "/rest/getAlbum"+tt.query, nil)
			rec := httptest.NewRecorder()

			server.HandleGetAlbum(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestHandleGetArtist(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		query        string
		expectedCode int
		setupMock    func(m *mocks.MockiTunesClient)
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
							mockResult,
						},
					}, nil)
			},
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

			server := NewServer(cfg, nil, mockClient)

			req := httptest.NewRequest(http.MethodGet, "/rest/getArtist"+tt.query, nil)
			rec := httptest.NewRecorder()

			server.HandleGetArtist(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}
