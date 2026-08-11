package handlers

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
		MusicFolders           *MusicFolders            `json:"musicFolders,omitempty"`
		User                   *User                    `json:"user,omitempty"`
		OpenSubSonicExtensions *[]OpenSubSonicExtension `json:"openSubsonicExtensions,omitempty"`
		JukeboxPlaylist        *JukeboxPlaylist         `json:"jukeboxPlaylist,omitempty"`
		AlbumList2             *AlbumList2              `json:"albumList2,omitempty"`
		Genres                 *[]Genre                 `json:"genres,omitempty"`
		Playlists              *[]Playlist              `json:"playlists,omitempty"`
		Album                  *Album                   `json:"album,omitempty"`
		SearchResult3          *SearchResult3           `json:"searchResult3,omitempty"`
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

	JukeboxPlaylist struct {
		CurrentIndex int     `json:"currentIndex"`
		Playing      bool    `json:"playing"`
		Gain         float32 `json:"gain"`
		Position     int     `json:"position"`
		Entry        *Song   `json:"entry,omitempty"`
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

	MusicFolders struct {
		MusicFolder []MusicFolder `json:"musicFolder"`
	}

	MusicFolder struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	AlbumList2 struct {
		Album []AlbumList2Album `json:"album"`
	}

	AlbumList2Album struct {
		Id        string `json:"id"`
		Album     string `json:"album"`
		Title     string `json:"title"`
		Name      string `json:"name"`
		CoverArt  string `json:"coverArt"`
		SongCount int    `json:"songCount"`
		Created   string `json:"created"`
		Duration  int    `json:"duration"`
		PlayCount int    `json:"playCount"`
		ArtistId  string `json:"artistId"`
		Artist    string `json:"artist"`
		Year      int    `json:"year"`
		Genre     string `json:"genre"`
	}

	Album struct {
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

	Genre struct {
		SongCount  int    `json:"songCount"`
		AlbumCount int    `json:"albumCount"`
		Value      string `json:"value"`
	}

	Playlist struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		Owner     string `json:"admin"`
		Public    bool   `json:"public"`
		Created   string `json:"created"`
		Changed   string `json:"changed"`
		SongCount int    `json:"songCount"`
		Duration  int    `json:"duration"`
	}

	SearchResult3 struct {
		Artist []Artist `json:"artist,omitempty"`
		Album  []Album  `json:"album,omitempty"`
		Song   []Song   `json:"song,omitempty"`
	}

	Artist struct {
		Id             string `json:"id"`
		Name           string `json:"name"`
		CoverArt       string `json:"coverArt,omitempty"`
		AlbumCount     int    `json:"albumCount,omitempty"`
		UserRating     int    `json:"userRating,omitempty"`
		ArtistImageUrl string `json:"artistImageUrl,omitempty"`
	}

	Generic401Response struct{}
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
			MusicFolders:           nil,
			User:                   nil,
			OpenSubSonicExtensions: nil,
			JukeboxPlaylist:        nil,
			AlbumList2:             nil,
			Genres:                 nil,
			Playlists:              nil,
			Album:                  nil,
			SearchResult3:          nil,

			Error: nil,
		},
	}
}

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
