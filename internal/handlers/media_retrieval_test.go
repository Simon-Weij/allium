package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Simon-Weij/allium/generated/mocks"
	"github.com/Simon-Weij/allium/internal/resolver"
	"github.com/Simon-Weij/allium/internal/subsonic"
	"github.com/Simon-Weij/allium/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandleGetCoverArt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		query        string
		expectedCode int
		expectedType string
		setupMock    func(m *mocks.MockiTunesClient)
	}{
		{
			name:         "should run as expected",
			query:        "?id=5",
			expectedCode: http.StatusOK,
			expectedType: "image/png",
			setupMock: func(m *mocks.MockiTunesClient) {
				coverPath := filepath.Join(t.TempDir(), "cover.png")
				require.NoError(t, os.WriteFile(coverPath, []byte("cover-art"), 0o644))
				m.EXPECT().
					GetAlbumCover("5").
					Return(coverPath, nil)
			},
		},
		{
			name:         "GetAlbumCover fails",
			query:        "?id=5",
			expectedCode: http.StatusInternalServerError,
			setupMock: func(m *mocks.MockiTunesClient) {
				m.EXPECT().
					GetAlbumCover("5").
					Return("", errors.New("something went wrong"))
			},
		},
		{
			name:         "cover file cannot be read",
			query:        "?id=5",
			expectedCode: http.StatusInternalServerError,
			setupMock: func(m *mocks.MockiTunesClient) {
				m.EXPECT().
					GetAlbumCover("5").
					Return(filepath.Join(t.TempDir(), "missing.png"), nil)
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

			req := httptest.NewRequest(http.MethodGet, "/rest/getCoverArt"+tt.query, nil)
			rec := httptest.NewRecorder()

			server.HandleGetCoverArt(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			if tt.expectedType != "" {
				assert.Equal(t, tt.expectedType, rec.Header().Get("Content-Type"))
				assert.Equal(t, "cover-art", rec.Body.String())
			}
		})
	}
}

func TestHandleStream(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		query         string
		expectedCode  int
		expectedType  string
		expectedBody  string
		expectedError *subsonic.Error
		setupMocks    func(m *mocks.MockiTunesClient, s *mocks.MockSongDownloader)
	}{
		{
			name:         "should stream the song",
			query:        "?id=777888999",
			expectedCode: http.StatusOK,
			expectedType: "audio/mpeg",
			expectedBody: "audio-content",
			setupMocks: func(m *mocks.MockiTunesClient, s *mocks.MockSongDownloader) {
				m.EXPECT().
					GetSongById("777888999").
					Return(&resolver.ITunesResponse{
						ResultCount: 1,
						Results:     []resolver.ITunesResult{mockResult},
					}, nil)

				songPath := filepath.Join(t.TempDir(), "song.mp3")
				require.NoError(t, os.WriteFile(songPath, []byte("audio-content"), 0o644))

				s.EXPECT().
					DownloadOrGetSong(gomock.Any(), "Alice", "Test Song").
					Return(songPath, nil)
			},
		},
		{
			name:          "missing id should return error",
			query:         "",
			expectedCode:  http.StatusBadRequest,
			expectedError: &subsonic.Error{Code: subsonic.ErrParameterMissing},
			setupMocks:    func(m *mocks.MockiTunesClient, s *mocks.MockSongDownloader) {},
		},
		{
			name:         "GetSongById fails",
			query:        "?id=777888999",
			expectedCode: http.StatusInternalServerError,
			setupMocks: func(m *mocks.MockiTunesClient, s *mocks.MockSongDownloader) {
				m.EXPECT().
					GetSongById("777888999").
					Return(nil, errors.New("something went wrong"))
			},
		},
		{
			name:         "no results",
			query:        "?id=777888999",
			expectedCode: http.StatusNotFound,
			setupMocks: func(m *mocks.MockiTunesClient, s *mocks.MockSongDownloader) {
				m.EXPECT().
					GetSongById("777888999").
					Return(&resolver.ITunesResponse{
						ResultCount: 0,
						Results:     []resolver.ITunesResult{},
					}, nil)
			},
		},
		{
			name:         "DownloadOrGetSong fails",
			query:        "?id=777888999",
			expectedCode: http.StatusInternalServerError,
			setupMocks: func(m *mocks.MockiTunesClient, s *mocks.MockSongDownloader) {
				m.EXPECT().
					GetSongById("777888999").
					Return(&resolver.ITunesResponse{
						ResultCount: 1,
						Results:     []resolver.ITunesResult{mockResult},
					}, nil)

				s.EXPECT().
					DownloadOrGetSong(gomock.Any(), "Alice", "Test Song").
					Return("", errors.New("something went wrong"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := testutil.SetupTestingConfig(t)

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockiTunesClient(ctrl)
			mockDownloader := mocks.NewMockSongDownloader(ctrl)
			tt.setupMocks(mockClient, mockDownloader)

			server := NewServer(cfg, nil, mockClient, mockDownloader)

			req := httptest.NewRequest(http.MethodGet, "/rest/stream.view"+tt.query, nil)
			rec := httptest.NewRecorder()

			server.HandleStream(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			if tt.expectedType != "" {
				assert.Equal(t, tt.expectedType, rec.Header().Get("Content-Type"))
			}

			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, rec.Body.String())
			}

			if tt.expectedError != nil {
				var res subsonic.SubsonicResponseWrapper

				err := json.NewDecoder(rec.Body).Decode(&res)
				require.NoError(t, err)

				assert.Equal(t, "failed", res.SubsonicResponse.Status)
				require.NotNil(t, res.SubsonicResponse.Error)
				assert.Equal(t, tt.expectedError.Code, res.SubsonicResponse.Error.Code)
			}
		})
	}
}
