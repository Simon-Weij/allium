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

			server := NewServer(cfg, nil, mockClient)

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
