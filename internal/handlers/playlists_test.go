package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Simon-Weij/allium/generated/mocks"
	"github.com/Simon-Weij/allium/generated/sqlc"
	"github.com/Simon-Weij/allium/internal/resolver"
	"github.com/Simon-Weij/allium/internal/subsonic"
	"github.com/Simon-Weij/allium/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func expectAddSongsToPlaylist(
	m *mocks.MockQueries,
	i *mocks.MockiTunesClient,
	songIDs []string,
	playlistID string,
	user string,
) {
	for idx, sid := range songIDs {
		i.EXPECT().
			GetSongById(sid).
			Return(&resolver.ITunesResponse{
				ResultCount: 1,
				Results:     []resolver.ITunesResult{mockResult},
			}, nil)
		idInt, _ := strconv.Atoi(sid)
		m.EXPECT().
			UpdatePlays(gomock.Any(), sqlc.UpdatePlaysParams{
				ID:         int64(idInt),
				Title:      mockResult.TrackName,
				Artist:     mockResult.ArtistName,
				Album:      mockResult.CollectionName,
				AlbumID:    int64(mockResult.CollectionID),
				ArtistID:   int64(mockResult.ArtistID),
				Duration:   int64(mockResult.TrackTimeMillis / 1000),
				Track:      int64(mockResult.TrackNumber),
				Year:       2024,
				Genre:      mockResult.PrimaryGenreName,
				DiscNumber: int64(mockResult.DiscNumber),
				User:       user,
				ArtworkUrl: mockResult.ArtworkURL100,
			}).
			Return(nil)
		m.EXPECT().GetMaxPlaylistPosition(gomock.Any(), playlistID).Return(int64(idx-1), nil)
		m.EXPECT().InsertIntoPlaylist(gomock.Any(), sqlc.InsertIntoPlaylistParams{
			PlaylistID: playlistID,
			SongID:     sid,
			Position:   int64(idx),
		}).Return(sqlc.PlaylistSong{}, nil)
	}
}

func TestHandleCreatePlaylist(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		query        string
		expectedCode int
		expectedName string
		setup        func(*mocks.MockQueries, *mocks.MockiTunesClient)
	}{
		{
			name:         "creates playlist with given name",
			query:        "?name=abc",
			expectedCode: http.StatusOK,
			expectedName: "abc",
			setup: func(m *mocks.MockQueries, _ *mocks.MockiTunesClient) {
				m.EXPECT().CreatePlaylist(gomock.Any(), sqlc.CreatePlaylistParams{
					Title:     "abc",
					Playcount: 0,
					User:      "",
				}).Return(sqlc.Playlist{
					ID:        1,
					Title:     "abc",
					UpdatedAt: "2024-01-01T00:00:00Z",
					Playcount: 0,
					User:      "",
				}, nil)
			},
		},
		{
			name:         "creates playlist with name and song ids",
			query:        "?name=abc&songId=123&songId=456",
			expectedCode: http.StatusOK,
			expectedName: "abc",
			setup: func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {
				m.EXPECT().CreatePlaylist(gomock.Any(), sqlc.CreatePlaylistParams{
					Title:     "abc",
					Playcount: 0,
					User:      "",
				}).Return(sqlc.Playlist{
					ID:        1,
					Title:     "abc",
					UpdatedAt: "2024-01-01T00:00:00Z",
					Playcount: 0,
					User:      "",
				}, nil)

				expectAddSongsToPlaylist(m, i, []string{"123", "456"}, "1", "")
			},
		},
		{
			name:         "errors when neither name nor playlistId provided",
			query:        "?songId=123",
			expectedCode: http.StatusBadRequest,
			setup:        func(_ *mocks.MockQueries, _ *mocks.MockiTunesClient) {},
		},
		{
			name:         "errors when BOTH playlistId and name are provided",
			query:        "?name=testplaylist&playlistId=1234",
			expectedCode: http.StatusOK,
			expectedName: "testplaylist",
			setup: func(m *mocks.MockQueries, _ *mocks.MockiTunesClient) {
				m.EXPECT().UpdatePlaylistTitle(gomock.Any(), sqlc.UpdatePlaylistTitleParams{
					Title: "testplaylist",
					ID:    1234,
				}).Return(sqlc.Playlist{
					ID:        1234,
					Title:     "testplaylist",
					UpdatedAt: "2024-01-01T00:00:00Z",
					Playcount: 0,
					User:      "",
				}, nil)
			},
		},
		{
			name:         "updates existing playlist by playlistId",
			query:        "?playlistId=789&name=newname",
			expectedCode: http.StatusOK,
			expectedName: "newname",
			setup: func(m *mocks.MockQueries, _ *mocks.MockiTunesClient) {
				m.EXPECT().UpdatePlaylistTitle(gomock.Any(), sqlc.UpdatePlaylistTitleParams{
					Title: "newname",
					ID:    789,
				}).Return(sqlc.Playlist{
					ID:        789,
					Title:     "newname",
					UpdatedAt: "2024-01-01T00:00:00Z",
					Playcount: 0,
					User:      "",
				}, nil)
			},
		},
		{
			name:         "errors on invalid songId",
			query:        "?name=abc&songId=notanumber",
			expectedCode: http.StatusBadRequest,
			setup:        func(_ *mocks.MockQueries, _ *mocks.MockiTunesClient) {},
		},
		{
			name:         "errors when adding songs to playlist fails",
			query:        "?name=abc&songId=123",
			expectedCode: http.StatusInternalServerError,
			setup: func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {
				m.EXPECT().CreatePlaylist(gomock.Any(), sqlc.CreatePlaylistParams{
					Title:     "abc",
					Playcount: 0,
					User:      "",
				}).Return(sqlc.Playlist{
					ID:        1,
					Title:     "abc",
					UpdatedAt: "2024-01-01T00:00:00Z",
					Playcount: 0,
					User:      "",
				}, nil)
				i.EXPECT().GetSongById("123").Return(nil, errors.New("itunes error"))
			},
		},
		{
			name:         "errors when creating playlist fails",
			query:        "?name=abc",
			expectedCode: http.StatusInternalServerError,
			setup: func(m *mocks.MockQueries, _ *mocks.MockiTunesClient) {
				m.EXPECT().CreatePlaylist(gomock.Any(), sqlc.CreatePlaylistParams{
					Title:     "abc",
					Playcount: 0,
					User:      "",
				}).Return(sqlc.Playlist{}, errors.New("db error"))
			},
		},
		{
			name:         "errors when updating existing playlist fails",
			query:        "?playlistId=789&name=newname",
			expectedCode: http.StatusInternalServerError,
			setup: func(m *mocks.MockQueries, _ *mocks.MockiTunesClient) {
				m.EXPECT().UpdatePlaylistTitle(gomock.Any(), sqlc.UpdatePlaylistTitleParams{
					Title: "newname",
					ID:    789,
				}).Return(sqlc.Playlist{}, errors.New("db error"))
			},
		},
		{
			name:         "errors when playlistId is not a number",
			query:        "?playlistId=abc&name=test",
			expectedCode: http.StatusInternalServerError,
			setup:        func(_ *mocks.MockQueries, _ *mocks.MockiTunesClient) {},
		},
		{
			name:         "errors when getting max playlist position fails",
			query:        "?name=abc&songId=123",
			expectedCode: http.StatusInternalServerError,
			setup: func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {
				m.EXPECT().CreatePlaylist(gomock.Any(), sqlc.CreatePlaylistParams{
					Title:     "abc",
					Playcount: 0,
					User:      "",
				}).Return(sqlc.Playlist{
					ID:        1,
					Title:     "abc",
					UpdatedAt: "2024-01-01T00:00:00Z",
					Playcount: 0,
					User:      "",
				}, nil)
				i.EXPECT().GetSongById("123").Return(&resolver.ITunesResponse{
					ResultCount: 1,
					Results:     []resolver.ITunesResult{mockResult},
				}, nil)
				m.EXPECT().UpdatePlays(gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().GetMaxPlaylistPosition(gomock.Any(), "1").Return(int64(0), errors.New("db error"))
			},
		},
		{
			name:         "errors when inserting into playlist fails",
			query:        "?name=abc&songId=123",
			expectedCode: http.StatusInternalServerError,
			setup: func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {
				m.EXPECT().CreatePlaylist(gomock.Any(), sqlc.CreatePlaylistParams{
					Title:     "abc",
					Playcount: 0,
					User:      "",
				}).Return(sqlc.Playlist{
					ID:        1,
					Title:     "abc",
					UpdatedAt: "2024-01-01T00:00:00Z",
					Playcount: 0,
					User:      "",
				}, nil)
				i.EXPECT().GetSongById("123").Return(&resolver.ITunesResponse{
					ResultCount: 1,
					Results:     []resolver.ITunesResult{mockResult},
				}, nil)
				m.EXPECT().UpdatePlays(gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().GetMaxPlaylistPosition(gomock.Any(), "1").Return(int64(-1), nil)
				m.EXPECT().InsertIntoPlaylist(gomock.Any(), sqlc.InsertIntoPlaylistParams{
					PlaylistID: "1",
					SongID:     "123",
					Position:   0,
				}).Return(sqlc.PlaylistSong{}, errors.New("db error"))
			},
		},
		{
			name:         "errors when song not found in iTunes",
			query:        "?name=abc&songId=123",
			expectedCode: http.StatusInternalServerError,
			setup: func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {
				m.EXPECT().CreatePlaylist(gomock.Any(), sqlc.CreatePlaylistParams{
					Title:     "abc",
					Playcount: 0,
					User:      "",
				}).Return(sqlc.Playlist{
					ID:        1,
					Title:     "abc",
					UpdatedAt: "2024-01-01T00:00:00Z",
					Playcount: 0,
					User:      "",
				}, nil)
				i.EXPECT().GetSongById("123").Return(&resolver.ITunesResponse{
					ResultCount: 0,
					Results:     []resolver.ITunesResult{},
				}, nil)
			},
		},
		{
			name:         "errors when storing song fails",
			query:        "?name=abc&songId=123",
			expectedCode: http.StatusInternalServerError,
			setup: func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {
				m.EXPECT().CreatePlaylist(gomock.Any(), sqlc.CreatePlaylistParams{
					Title:     "abc",
					Playcount: 0,
					User:      "",
				}).Return(sqlc.Playlist{
					ID:        1,
					Title:     "abc",
					UpdatedAt: "2024-01-01T00:00:00Z",
					Playcount: 0,
					User:      "",
				}, nil)
				i.EXPECT().GetSongById("123").Return(&resolver.ITunesResponse{
					ResultCount: 1,
					Results:     []resolver.ITunesResult{mockResult},
				}, nil)
				m.EXPECT().UpdatePlays(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cfg := testutil.SetupTestingConfig(t)
			mockQueries := mocks.NewMockQueries(ctrl)
			mockiTunes := mocks.NewMockiTunesClient(ctrl)

			tt.setup(mockQueries, mockiTunes)

			server := NewServer(cfg, mockQueries, mockiTunes, nil)

			req := httptest.NewRequest(http.MethodGet, "/rest/createPlaylist"+tt.query, nil)
			rec := httptest.NewRecorder()

			server.HandleCreatePlaylist(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			if tt.expectedCode == http.StatusOK {
				var res subsonic.SubsonicResponseWrapper

				err := json.NewDecoder(rec.Body).Decode(&res)
				require.NoError(t, err)

				require.NotNil(t, res.SubsonicResponse.Playlist)
				assert.Equal(t, tt.expectedName, res.SubsonicResponse.Playlist.Name)
			}
		})
	}
}

func TestHandleGetPlaylists(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		query        string
		expectedCode int
		setup        func(*mocks.MockQueries)
	}{
		{
			name:         "returns empty playlists when user has none",
			query:        "",
			expectedCode: http.StatusOK,
			setup: func(m *mocks.MockQueries) {
				m.EXPECT().GetPlaylistsByUser(gomock.Any(), "").Return([]sqlc.GetPlaylistsByUserRow{}, nil)
			},
		},
		{
			name:         "returns playlists for authenticated user",
			query:        "?u=alice",
			expectedCode: http.StatusOK,
			setup: func(m *mocks.MockQueries) {
				m.EXPECT().GetPlaylistsByUser(gomock.Any(), "alice").Return([]sqlc.GetPlaylistsByUserRow{
					{
						ID:        1,
						Title:     "my playlist",
						UpdatedAt: "2024-01-01T00:00:00Z",
						Playcount: 0,
						User:      "alice",
						SongCount: 2,
					},
					{
						ID:        2,
						Title:     "another playlist",
						UpdatedAt: "2024-01-02T00:00:00Z",
						Playcount: 5,
						User:      "alice",
						SongCount: 0,
					},
				}, nil)
			},
		},
		{
			name:         "uses username param when provided",
			query:        "?u=alice&username=bob",
			expectedCode: http.StatusOK,
			setup: func(m *mocks.MockQueries) {
				m.EXPECT().GetPlaylistsByUser(gomock.Any(), "bob").Return([]sqlc.GetPlaylistsByUserRow{
					{
						ID:        3,
						Title:     "bob's playlist",
						UpdatedAt: "2024-01-03T00:00:00Z",
						Playcount: 1,
						User:      "bob",
						SongCount: 10,
					},
				}, nil)
			},
		},
		{
			name:         "errors when fetching playlists fails",
			query:        "?u=alice",
			expectedCode: http.StatusInternalServerError,
			setup: func(m *mocks.MockQueries) {
				m.EXPECT().GetPlaylistsByUser(gomock.Any(), "alice").Return(nil, errors.New("db error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cfg := testutil.SetupTestingConfig(t)
			mockQueries := mocks.NewMockQueries(ctrl)

			tt.setup(mockQueries)

			server := NewServer(cfg, mockQueries, nil, nil)

			req := httptest.NewRequest(http.MethodGet, "/rest/getPlaylists"+tt.query, nil)
			rec := httptest.NewRecorder()

			server.HandleGetPlaylists(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			if tt.expectedCode != http.StatusOK {
				return
			}

			var res subsonic.SubsonicResponseWrapper

			err := json.NewDecoder(rec.Body).Decode(&res)
			require.NoError(t, err)

			require.NotNil(t, res.SubsonicResponse.Playlists)
		})
	}
}

func TestHandleGetPlaylist(t *testing.T) {
	t.Parallel()

	playlist := sqlc.Playlist{
		ID:        42,
		Title:     "my playlist",
		UpdatedAt: "2024-01-01T00:00:00Z",
		Playcount: 3,
		User:      "alice",
	}

	songs := []sqlc.Song{
		{
			ID:     1,
			Title:  "song one",
			Artist: "artist a",
		},
		{
			ID:     2,
			Title:  "song two",
			Artist: "artist b",
		},
	}

	tests := []struct {
		name         string
		query        string
		expectedCode int
		setup        func(*mocks.MockQueries)
	}{
		{
			name:         "returns playlist with songs",
			query:        "?id=42",
			expectedCode: http.StatusOK,
			setup: func(m *mocks.MockQueries) {
				m.EXPECT().GetPlaylistByID(gomock.Any(), int64(42)).Return(playlist, nil)
				m.EXPECT().GetSongsByPlaylistID(gomock.Any(), "42").Return(songs, nil)
			},
		},
		{
			name:         "returns playlist with no songs",
			query:        "?id=42",
			expectedCode: http.StatusOK,
			setup: func(m *mocks.MockQueries) {
				m.EXPECT().GetPlaylistByID(gomock.Any(), int64(42)).Return(playlist, nil)
				m.EXPECT().GetSongsByPlaylistID(gomock.Any(), "42").Return([]sqlc.Song{}, nil)
			},
		},
		{
			name:         "errors when id is missing",
			query:        "",
			expectedCode: http.StatusBadRequest,
			setup:        func(m *mocks.MockQueries) {},
		},
		{
			name:         "errors when id is not a number",
			query:        "?id=notanumber",
			expectedCode: http.StatusBadRequest,
			setup:        func(m *mocks.MockQueries) {},
		},
		{
			name:         "errors when playlist does not exist",
			query:        "?id=999",
			expectedCode: http.StatusNotFound,
			setup: func(m *mocks.MockQueries) {
				m.EXPECT().
					GetPlaylistByID(gomock.Any(), int64(999)).
					Return(sqlc.Playlist{}, sql.ErrNoRows)
			},
		},
		{
			name:         "errors when fetching songs fails",
			query:        "?id=42",
			expectedCode: http.StatusInternalServerError,
			setup: func(m *mocks.MockQueries) {
				m.EXPECT().GetPlaylistByID(gomock.Any(), int64(42)).Return(playlist, nil)
				m.EXPECT().
					GetSongsByPlaylistID(gomock.Any(), "42").
					Return(nil, errors.New("db error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cfg := testutil.SetupTestingConfig(t)
			mockQueries := mocks.NewMockQueries(ctrl)

			tt.setup(mockQueries)

			server := NewServer(cfg, mockQueries, nil, nil)

			req := httptest.NewRequest(http.MethodGet, "/rest/getPlaylist"+tt.query, nil)
			rec := httptest.NewRecorder()

			server.HandleGetPlaylist(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			if tt.expectedCode != http.StatusOK {
				return
			}

			var res subsonic.SubsonicResponseWrapper

			err := json.NewDecoder(rec.Body).Decode(&res)
			require.NoError(t, err)

			require.NotNil(t, res.SubsonicResponse.Playlist)
			actual := res.SubsonicResponse.Playlist
			assert.Equal(t, strconv.FormatInt(playlist.ID, 10), actual.Id)
			assert.Equal(t, playlist.Title, actual.Name)
			assert.Equal(t, playlist.User, actual.Owner)
			assert.False(t, actual.Public)
			assert.Equal(t, playlist.UpdatedAt, actual.Created)
			assert.Equal(t, playlist.UpdatedAt, actual.Changed)
			assert.Equal(t, 0, actual.Duration)

			if tt.query == "?id=42" && tt.name == "returns playlist with songs" {
				assert.Equal(t, len(songs), actual.SongCount)
				require.Len(t, actual.Entry, len(songs))

				for i, song := range songs {
					assert.Equal(t, strconv.FormatInt(song.ID, 10), actual.Entry[i].Id)
					assert.Equal(t, song.Title, actual.Entry[i].Title)
					assert.Equal(t, song.Artist, actual.Entry[i].Artist)
				}
			}
		})
	}
}

func TestHandleUpdatePlaylist(t *testing.T) {
	t.Parallel()

	existingPlaylist := sqlc.Playlist{
		ID:        123,
		Title:     "old name",
		UpdatedAt: "2024-01-01T00:00:00Z",
		Playcount: 0,
		User:      "alice",
		Comment:   "old comment",
		Public:    0,
	}

	tests := []struct {
		name         string
		query        string
		expectedCode int
		setup        func(*mocks.MockQueries, *mocks.MockiTunesClient)
	}{
		{
			name:         "errors when playlistId is missing",
			query:        "",
			expectedCode: http.StatusBadRequest,
			setup:        func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {},
		},
		{
			name:         "errors when playlistId is not a number",
			query:        "?playlistId=abc",
			expectedCode: http.StatusBadRequest,
			setup:        func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {},
		},
		{
			name:         "errors when playlist does not exist",
			query:        "?playlistId=999",
			expectedCode: http.StatusNotFound,
			setup: func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {
				m.EXPECT().GetPlaylistByID(gomock.Any(), int64(999)).Return(sqlc.Playlist{}, errors.New("no rows"))
			},
		},
		{
			name:         "updates playlist name",
			query:        "?playlistId=123&name=new+name",
			expectedCode: http.StatusOK,
			setup: func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {
				m.EXPECT().GetPlaylistByID(gomock.Any(), int64(123)).Return(existingPlaylist, nil)
				m.EXPECT().UpdatePlaylistMetadata(gomock.Any(), sqlc.UpdatePlaylistMetadataParams{
					Title:   "new name",
					Comment: "old comment",
					Public:  0,
					ID:      123,
				}).Return(sqlc.Playlist{
					ID: 123, Title: "new name", UpdatedAt: "2024-01-01T00:00:00Z",
					Playcount: 0, User: "alice", Comment: "old comment", Public: 0,
				}, nil)
			},
		},
		{
			name:         "updates playlist comment",
			query:        "?playlistId=123&comment=new+comment",
			expectedCode: http.StatusOK,
			setup: func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {
				m.EXPECT().GetPlaylistByID(gomock.Any(), int64(123)).Return(existingPlaylist, nil)
				m.EXPECT().UpdatePlaylistMetadata(gomock.Any(), sqlc.UpdatePlaylistMetadataParams{
					Title:   "old name",
					Comment: "new comment",
					Public:  0,
					ID:      123,
				}).Return(sqlc.Playlist{}, nil)
			},
		},
		{
			name:         "updates playlist public to true",
			query:        "?playlistId=123&public=true",
			expectedCode: http.StatusOK,
			setup: func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {
				m.EXPECT().GetPlaylistByID(gomock.Any(), int64(123)).Return(existingPlaylist, nil)
				m.EXPECT().UpdatePlaylistMetadata(gomock.Any(), sqlc.UpdatePlaylistMetadataParams{
					Title:   "old name",
					Comment: "old comment",
					Public:  1,
					ID:      123,
				}).Return(sqlc.Playlist{}, nil)
			},
		},
		{
			name:         "updates playlist public to false",
			query:        "?playlistId=123&public=false",
			expectedCode: http.StatusOK,
			setup: func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {
				m.EXPECT().GetPlaylistByID(gomock.Any(), int64(123)).Return(sqlc.Playlist{
					ID: 123, Title: "old name", User: "alice", Public: 1,
				}, nil)
				m.EXPECT().UpdatePlaylistMetadata(gomock.Any(), sqlc.UpdatePlaylistMetadataParams{
					Title:   "old name",
					Comment: "",
					Public:  0,
					ID:      123,
				}).Return(sqlc.Playlist{}, nil)
			},
		},
		{
			name:         "updates all fields",
			query:        "?playlistId=123&name=updated&comment=updated&public=true",
			expectedCode: http.StatusOK,
			setup: func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {
				m.EXPECT().GetPlaylistByID(gomock.Any(), int64(123)).Return(existingPlaylist, nil)
				m.EXPECT().UpdatePlaylistMetadata(gomock.Any(), sqlc.UpdatePlaylistMetadataParams{
					Title:   "updated",
					Comment: "updated",
					Public:  1,
					ID:      123,
				}).Return(sqlc.Playlist{}, nil)
			},
		},
		{
			name:         "removes song by index",
			query:        "?playlistId=123&songIndexToRemove=0&songIndexToRemove=2",
			expectedCode: http.StatusOK,
			setup: func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {
				m.EXPECT().GetPlaylistByID(gomock.Any(), int64(123)).Return(existingPlaylist, nil)
				m.EXPECT().UpdatePlaylistMetadata(gomock.Any(), sqlc.UpdatePlaylistMetadataParams{
					Title:   "old name",
					Comment: "old comment",
					Public:  0,
					ID:      123,
				}).Return(sqlc.Playlist{}, nil)
				m.EXPECT().DeleteSongFromPlaylistByPosition(gomock.Any(), sqlc.DeleteSongFromPlaylistByPositionParams{
					PlaylistID: "123",
					Position:   0,
				}).Return(nil)
				m.EXPECT().DeleteSongFromPlaylistByPosition(gomock.Any(), sqlc.DeleteSongFromPlaylistByPositionParams{
					PlaylistID: "123",
					Position:   2,
				}).Return(nil)
			},
		},
		{
			name:         "adds songs by songIdToAdd",
			query:        "?playlistId=123&songIdToAdd=100&songIdToAdd=200",
			expectedCode: http.StatusOK,
			setup: func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {
				m.EXPECT().GetPlaylistByID(gomock.Any(), int64(123)).Return(existingPlaylist, nil)

				updatedPlaylist := sqlc.Playlist{
					ID: 123, Title: "old name", UpdatedAt: "2024-01-01T00:00:00Z",
					Playcount: 0, User: "alice", Comment: "old comment", Public: 0,
				}
				m.EXPECT().UpdatePlaylistMetadata(gomock.Any(), sqlc.UpdatePlaylistMetadataParams{
					Title:   "old name",
					Comment: "old comment",
					Public:  0,
					ID:      123,
				}).Return(updatedPlaylist, nil)

				expectAddSongsToPlaylist(m, i, []string{"100", "200"}, "123", "alice")
			},
		},
		{
			name:         "errors on invalid songIdToAdd",
			query:        "?playlistId=123&songIdToAdd=notanumber",
			expectedCode: http.StatusBadRequest,
			setup: func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {
				m.EXPECT().GetPlaylistByID(gomock.Any(), int64(123)).Return(existingPlaylist, nil)
				m.EXPECT().UpdatePlaylistMetadata(gomock.Any(), gomock.Any()).Return(sqlc.Playlist{}, nil)
			},
		},
		{
			name:         "errors on invalid songIndexToRemove",
			query:        "?playlistId=123&songIndexToRemove=abc",
			expectedCode: http.StatusBadRequest,
			setup: func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {
				m.EXPECT().GetPlaylistByID(gomock.Any(), int64(123)).Return(existingPlaylist, nil)
				m.EXPECT().UpdatePlaylistMetadata(gomock.Any(), gomock.Any()).Return(sqlc.Playlist{}, nil)
			},
		},
		{
			name:         "errors when updating playlist resolver fails",
			query:        "?playlistId=123&name=new+name",
			expectedCode: http.StatusInternalServerError,
			setup: func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {
				m.EXPECT().GetPlaylistByID(gomock.Any(), int64(123)).Return(existingPlaylist, nil)
				m.EXPECT().UpdatePlaylistMetadata(gomock.Any(), sqlc.UpdatePlaylistMetadataParams{
					Title:   "new name",
					Comment: "old comment",
					Public:  0,
					ID:      123,
				}).Return(sqlc.Playlist{}, errors.New("db error"))
			},
		},
		{
			name:         "errors when adding songs fails in update",
			query:        "?playlistId=123&songIdToAdd=100",
			expectedCode: http.StatusInternalServerError,
			setup: func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {
				m.EXPECT().GetPlaylistByID(gomock.Any(), int64(123)).Return(existingPlaylist, nil)
				m.EXPECT().UpdatePlaylistMetadata(gomock.Any(), gomock.Any()).Return(existingPlaylist, nil)
				i.EXPECT().GetSongById("100").Return(nil, errors.New("itunes error"))
			},
		},
		{
			name:         "errors when removing song from playlist fails",
			query:        "?playlistId=123&songIndexToRemove=0",
			expectedCode: http.StatusInternalServerError,
			setup: func(m *mocks.MockQueries, i *mocks.MockiTunesClient) {
				m.EXPECT().GetPlaylistByID(gomock.Any(), int64(123)).Return(existingPlaylist, nil)
				m.EXPECT().UpdatePlaylistMetadata(gomock.Any(), gomock.Any()).Return(sqlc.Playlist{}, nil)
				m.EXPECT().DeleteSongFromPlaylistByPosition(gomock.Any(), sqlc.DeleteSongFromPlaylistByPositionParams{
					PlaylistID: "123",
					Position:   0,
				}).Return(errors.New("db error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cfg := testutil.SetupTestingConfig(t)
			mockQueries := mocks.NewMockQueries(ctrl)
			mockiTunes := mocks.NewMockiTunesClient(ctrl)

			tt.setup(mockQueries, mockiTunes)

			server := NewServer(cfg, mockQueries, mockiTunes, nil)

			req := httptest.NewRequest(http.MethodGet, "/rest/updatePlaylist"+tt.query, nil)
			rec := httptest.NewRecorder()

			server.HandleUpdatePlaylist(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			if tt.expectedCode == http.StatusOK {
				var res subsonic.SubsonicResponseWrapper

				err := json.NewDecoder(rec.Body).Decode(&res)
				require.NoError(t, err)
				assert.Equal(t, "ok", res.SubsonicResponse.Status)
			}
		})
	}
}
