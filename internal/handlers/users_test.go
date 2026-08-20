package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Simon-Weij/allium/generated/mocks"
	"github.com/Simon-Weij/allium/generated/sqlc"
	"github.com/Simon-Weij/allium/internal/subsonic"
	"github.com/Simon-Weij/allium/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var mockUserResponse = &sqlc.User{
			Username:            "gopher",
			Email:               "gopher@example.com",
			ScrobblingEnabled:   true,
			AdminRole:           false,
			SettingsRole:        true,
			DownloadRole:        true,
			UploadRole:          false,
			PlaylistRole:        true,
			CoverArtRole:        true,
			CommentRole:         true,
			PodcastRole:         false,
			StreamRole:          true,
			JukeboxRole:         false,
			ShareRole:           true,
			VideoConversionRole: false,
		}

func TestHandleGetUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		query             string
		expectedCode      int
		setupMock         func(m *mocks.MockUserManagementClient)
	}{
		{
			name:              "should run correctly",
			query:             "?username=gopher",
			expectedCode:      http.StatusOK,
			setupMock: func(m *mocks.MockUserManagementClient) {
				m.EXPECT().
					GetUserByUsername(context.Background(),"gopher").
					Return(
						mockUserResponse,
						nil,
					)
			},
		},
		{
			name:              "empty username param (shouldn't work and should throw an error)",
			query:             "?username=",
			expectedCode:      http.StatusInternalServerError,
			setupMock: func(m *mocks.MockUserManagementClient) {},
		},
		{
			name:              "username not found in database (shouldn't work and should throw an error)",
			query:             "?username=someoneFake",
			expectedCode:      http.StatusInternalServerError,
			setupMock: func(m *mocks.MockUserManagementClient) {
				m.EXPECT().
					GetUserByUsername(context.Background(), "someoneFake").
					Return(
						&sqlc.User{},
						errors.New("user not found :("),
					)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := testutil.SetupTestingConfig(t)

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockUserManagementClient(ctrl)
			tt.setupMock(mockClient)

			server := NewServer(cfg, nil, nil, mockClient)
			
			req := httptest.NewRequest(http.MethodGet, "/rest/getUser.view"+tt.query, nil)
			rec := httptest.NewRecorder()

			server.HandleGetUser(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			var res subsonic.SubsonicResponseWrapper

			err := json.NewDecoder(rec.Body).Decode(&res)
			require.NoError(t, err)

		})
	}
}



var mockUsersResponse = &[]sqlc.User{
	sqlc.User{
			Username:            "gopher",
			Email:               "gopher@example.com",
			ScrobblingEnabled:   true,
			AdminRole:           false,
			SettingsRole:        true,
			DownloadRole:        true,
			UploadRole:          false,
			PlaylistRole:        true,
			CoverArtRole:        true,
			CommentRole:         true,
			PodcastRole:         false,
			StreamRole:          true,
			JukeboxRole:         false,
			ShareRole:           true,
			VideoConversionRole: false,
		},
}

func TestHandleGetUsers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		query             string
		expectedCode      int
		setupMock         func(m *mocks.MockUserManagementClient)
	}{
		{
			name:              "should run correctly",
			query:             "",
			expectedCode:      http.StatusOK,
			setupMock: func(m *mocks.MockUserManagementClient) {
				m.EXPECT().
					GetUsers(context.Background()).
					Return(
						mockUsersResponse,
						nil,
					)
			},
		},
		{
			name:              "users table empty (shouldn't work and should throw an error)",
			query:             "",
			expectedCode:      http.StatusInternalServerError,
			setupMock: func(m *mocks.MockUserManagementClient) {
				m.EXPECT().
					GetUsers(context.Background()).
					Return(
						&[]sqlc.User{},
						errors.New("users table empty :("),
					)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := testutil.SetupTestingConfig(t)

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockUserManagementClient(ctrl)
			tt.setupMock(mockClient)

			server := NewServer(cfg, nil, nil, mockClient)
			
			req := httptest.NewRequest(http.MethodGet, "/rest/getUsers.view"+tt.query, nil)
			rec := httptest.NewRecorder()

			server.HandleGetUsers(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			var res subsonic.SubsonicResponseWrapper

			err := json.NewDecoder(rec.Body).Decode(&res)
			require.NoError(t, err)

		})
	}
}


func TestHandleCreateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		query             string
		expectedCode      int
		setupMock         func(m *mocks.MockUserManagementClient)
	}{
		{
			name:              "should run correctly",
			query:             "?username=gopher&password=verysecure&email=gopher@example.com",
			expectedCode:      http.StatusOK,
			setupMock: func(m *mocks.MockUserManagementClient) {
				m.EXPECT().
					CreateUser(context.Background(), url.Values{
						"username": {"gopher"},
						"password": {"verysecure"},
						"email": {"gopher@example.com"},
					}).
					Return(
						nil,
					)
			},
		},
		{
			name:              "empty query (shouldn't work and should throw an error)",
			query:             "",
			expectedCode:      http.StatusBadRequest,
			setupMock: func(m *mocks.MockUserManagementClient) {},
		},
		{
			name:              "no username (shouldn't work and should throw an error)",
			query:             "?password=verysecure&email=gopher@example.com",
			expectedCode:      http.StatusBadRequest,
			setupMock: func(m *mocks.MockUserManagementClient) {},
		},
		{
			name:              "no password (shouldn't work and should throw an error)",
			query:             "?username=gopher&email=gopher@example.com",
			expectedCode:      http.StatusBadRequest,
			setupMock: func(m *mocks.MockUserManagementClient) {},
		},
		{
			name:              "no email (shouldn't work and should throw an error)",
			query:             "?username=gopher&password=verysecure",
			expectedCode:      http.StatusBadRequest,
			setupMock: func(m *mocks.MockUserManagementClient) {},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := testutil.SetupTestingConfig(t)

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockUserManagementClient(ctrl)
			tt.setupMock(mockClient)

			server := NewServer(cfg, nil, nil, mockClient)
			
			req := httptest.NewRequest(http.MethodGet, "/rest/createUser.view"+tt.query, nil)
			rec := httptest.NewRecorder()

			server.HandleCreateUser(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			var res subsonic.SubsonicResponseWrapper

			err := json.NewDecoder(rec.Body).Decode(&res)
			if err == io.EOF {
				return
			}
			require.NoError(t, err)

		})
	}
}
