package subsonic

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Simon-Weij/allium/internal/config"
)

type (
	SubsonicResponseWrapper struct {
		SubsonicResponse Subsonic `json:"subsonic-response"`
	}

	Subsonic struct {
		Status                 string                   `json:"status"`
		Version                string                   `json:"version"`
		Type                   string                   `json:"type"`
		ServerVersion          string                   `json:"serverVersion"`
		OpenSubsonic           bool                     `json:"openSubsonic"`
		License                *License                 `json:"license,omitempty"`
		User                   *User                    `json:"user,omitempty"`
		OpenSubSonicExtensions *[]OpenSubSonicExtension `json:"openSubsonicExtensions,omitempty"`
		SearchResult3          *SearchResult3           `json:"searchResult3,omitempty"`
		Album                  *GetAlbumAlbum           `json:"album,omitempty"`
		Artist                 *GetArtistArtist         `json:"artist,omitempty"`
		Playlist               *Playlist                `json:"playlist,omitempty"`
		Playlists              *Playlists               `json:"playlists,omitempty"`
		Error                  *Error                   `json:"error,omitempty"`
	}

	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	License struct {
		Valid          bool       `json:"valid"`
		Email          string     `json:"email,omitempty"`
		LicenseExpires *time.Time `json:"licenseExpires,omitempty"`
		TrialExpires   *time.Time `json:"trialExpires,omitempty"`
	}

	OpenSubSonicExtension struct {
		Name    string
		Version []int
	}

	User struct {
		Folder            []int  `json:"folder"`
		Email             string `json:"email"`
		ScrobblingEnabled bool   `json:"scrobblingEnabled"`
		AdminRole         bool   `json:"adminRole"`
		SettingsRole      bool   `json:"settingsRole"`
		DownloadRole      bool   `json:"downloadRole"`
		PlaylistRole      bool   `json:"playlistRole"`
		CoverArtRole      bool   `json:"coverArtRole"`
		CommentRole       bool   `json:"commentRole"`
		PodcastRole       bool   `json:"podcastRole"`
		StreamRole        bool   `json:"streamRole"`
		JukeboxRole       bool   `json:"jukeboxRole"`
		ShareRole         bool   `json:"shareRole"`
	}

	Playlist struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		Owner     string `json:"owner"`
		Public    bool   `json:"public"`
		Created   string `json:"created"`
		Changed   string `json:"changed"`
		SongCount int    `json:"songCount"`
		Duration  int    `json:"duration"`
		Entry     []Song `json:"entry,omitempty"`
	}

	Playlists struct {
		Playlist []Playlist `json:"playlist"`
	}

	Song struct {
		Id           string `json:"id"`
		Parent       string `json:"parent"`
		Title        string `json:"title"`
		IsDir        bool   `json:"isDir"`
		IsVideo      bool   `json:"isVideo"`
		Type         string `json:"type"`
		AlbumId      string `json:"albumId"`
		Album        string `json:"album"`
		ArtistId     string `json:"artistId"`
		Artist       string `json:"artist"`
		CoverArt     string `json:"coverArt"`
		Duration     int    `json:"duration"`
		BitRate      int    `json:"bitRate"`
		BitDepth     int    `json:"bitDepth"`
		SamplingRate int    `json:"samplingRate"`
		ChannelCount int    `json:"channelCount"`
		Track        int    `json:"track"`
		Year         int    `json:"year"`
		Genre        string `json:"genre"`
		Size         int    `json:"size"`
		DiscNumber   int    `json:"discNumber"`
		Suffix       string `json:"suffix"`
		ContentType  string `json:"contentType"`
		Path         string `json:"path"`
	}

	GetAlbumAlbum struct {
		Id        string `json:"id"`
		Parent    string `json:"parent"`
		Album     string `json:"album"`
		Title     string `json:"title"`
		Name      string `json:"name"`
		IsDir     bool   `json:"isDir"`
		CoverArt  string `json:"coverArt"`
		SongCount int    `json:"songCount"`
		Created   string `json:"created"`
		Duration  int    `json:"duration"`
		PlayCount int    `json:"playCount"`
		ArtistId  string `json:"artistId"`
		Artist    string `json:"artist"`
		Year      int    `json:"year"`
		Genre     string `json:"genre"`
		Song      []Song `json:"song"`
	}

	GetArtistArtist struct {
		Id             string           `json:"id"`
		Name           string           `json:"name"`
		CoverArt       string           `json:"coverArt,omitempty"`
		AlbumCount     int              `json:"albumCount,omitempty"`
		UserRating     int              `json:"userRating,omitempty"`
		ArtistImageUrl string           `json:"artistImageUrl,omitempty"`
		Album          []GetArtistAlbum `json:"album,omitempty"`
	}

	GetArtistAlbum struct {
		Id            string `json:"id"`
		Parent        string `json:"parent"`
		Album         string `json:"album"`
		Title         string `json:"title"`
		Name          string `json:"name"`
		IsDir         bool   `json:"isDir"`
		CoverArt      string `json:"coverArt"`
		SongCount     int    `json:"songCount"`
		Created       string `json:"created"`
		Duration      int    `json:"duration"`
		PlayCount     int    `json:"playCount"`
		ArtistId      string `json:"artistId"`
		Artist        string `json:"artist"`
		Year          int    `json:"year"`
		Genre         string `json:"genre"`
		UserRating    int    `json:"userRating,omitempty"`
		AverageRating int    `json:"averageRating,omitempty"`
		Starred       string `json:"starred,omitempty"`
	}

	SearchResult3Album struct {
		Id         string `json:"id"`
		Name       string `json:"name"`
		Artist     string `json:"artist"`
		Year       int    `json:"year"`
		CoverArt   string `json:"coverArt"`
		Starred    string `json:"starred"`
		Duration   int    `json:"duration"`
		PlayCount  int    `json:"playCount"`
		Played     string `json:"played"`
		Created    string `json:"created"`
		ArtistId   string `json:"artistId"`
		UserRating int    `json:"userRating"`
		SongCount  int    `json:"songCount"`
	}

	SearchResult3 struct {
		Artist []Artist             `json:"artist,omitempty"`
		Album  []SearchResult3Album `json:"album,omitempty"`
		Song   []Song               `json:"song,omitempty"`
	}

	Artist struct {
		Id             string `json:"id"`
		Name           string `json:"name"`
		CoverArt       string `json:"coverArt,omitempty"`
		AlbumCount     int    `json:"albumCount,omitempty"`
		UserRating     int    `json:"userRating,omitempty"`
		ArtistImageUrl string `json:"artistImageUrl,omitempty"`
	}
)

func WriteError(w http.ResponseWriter, httpStatus int, cfg config.Config, code int, message string) {
	res := NewEmptyResponse(cfg)
	res.SubsonicResponse.Status = "failed"
	res.SubsonicResponse.Error = &Error{
		Code:    code,
		Message: message,
	}

	WriteJSON(w, httpStatus, res)
}

func NewEmptyResponse(cfg config.Config) SubsonicResponseWrapper {
	return SubsonicResponseWrapper{
		SubsonicResponse: Subsonic{
			Status:        "ok",
			Version:       config.Version,
			Type:          config.Name,
			ServerVersion: config.Version,
			OpenSubsonic:  true,

			License:                nil,
			User:                   nil,
			OpenSubSonicExtensions: nil,
			SearchResult3:          nil,
			Album:                  nil,
			Artist:                 nil,
			Playlist:               nil,
			Playlists:              nil,

			Error: nil,
		},
	}
}

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
