package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Simon-Weij/allium/generated/mocks"
	"github.com/Simon-Weij/allium/generated/sqlc"
	"github.com/Simon-Weij/allium/internal/metadata"
	"github.com/Simon-Weij/allium/internal/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandleScrobble(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		query        string
		expectedCode int
		setupMocks   func(m *mocks.MockiTunesClient, q *mocks.MockQueries)
	}{
		{
			name:         "should run correctly",
			query:        "?id=777888999&u=alice",
			expectedCode: http.StatusOK,
			setupMocks: func(m *mocks.MockiTunesClient, q *mocks.MockQueries) {
				m.EXPECT().
					GetSongById("777888999").
					Return(&metadata.ITunesResponse{
						ResultCount: 1,
						Results:     []metadata.ITunesResult{mockResult},
					}, nil)

				q.EXPECT().
					UpdatePlays(
						gomock.Any(),
						sqlc.UpdatePlaysParams{
							ID:     777888999,
							Title:  "Test Song",
							Artist: "Alice",
							User:   "alice",
						},
					).
					Return(nil)
			},
		},
		{
			name:         "should 500 on invalid id",
			query:        "?id=notanumber",
			expectedCode: http.StatusInternalServerError,
			setupMocks:   func(m *mocks.MockiTunesClient, q *mocks.MockQueries) {},
		},
		{
			name:         "should 500 when getting song goes wrong",
			query:        "?id=777888999",
			expectedCode: http.StatusInternalServerError,
			setupMocks: func(m *mocks.MockiTunesClient, q *mocks.MockQueries) {
				m.EXPECT().
					GetSongById("777888999").
					Return(nil, errors.New("something went wrong"))
			},
		},
		{
			name:         "should 500 when updating plays goes wrong",
			query:        "?id=777888999&u=alice",
			expectedCode: http.StatusInternalServerError,
			setupMocks: func(m *mocks.MockiTunesClient, q *mocks.MockQueries) {
				m.EXPECT().
					GetSongById("777888999").
					Return(&metadata.ITunesResponse{
						ResultCount: 1,
						Results:     []metadata.ITunesResult{mockResult},
					}, nil)

				q.EXPECT().
					UpdatePlays(gomock.Any(), gomock.Any()).
					Return(errors.New("something went wrong"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := testutil.SetupTestingConfig(t)

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockiTunesClient(ctrl)
			mockQueries := mocks.NewMockQueries(ctrl)
			tt.setupMocks(mockClient, mockQueries)

			server := NewServer(cfg, mockQueries, mockClient)

			req := httptest.NewRequest(http.MethodGet, "/rest/scrobble.view"+tt.query, nil)
			rec := httptest.NewRecorder()

			server.HandleScrobble(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}
