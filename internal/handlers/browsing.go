package handlers

import (
	"io"
	"net/http"
)

// TODO: stop returning placeholders
func (s Server) HandleGetMusicFolders(w http.ResponseWriter, r *http.Request) {
	res := NewEmptyResponse(s.cfg)
	res.SubsonicResponse.MusicFolders = &MusicFolders{
		MusicFolder: []MusicFolder{
			{
				Id:   1,
				Name: "music",
			},
		},
	}
	WriteJSON(w, http.StatusOK, res)
}

// TODO: stop returning placeholders
func (s Server) HandleGetGenres(w http.ResponseWriter, r *http.Request) {
	res := NewEmptyResponse(s.cfg)
	res.SubsonicResponse.Genres = &[]Genre{
		{
			SongCount:  1,
			AlbumCount: 1,
			Value:      "Punk",
		},
		{
			SongCount:  4,
			AlbumCount: 1,
			Value:      "Dark Ambient",
		},
		{
			SongCount:  6,
			AlbumCount: 1,
			Value:      "Noise",
		},
		{
			SongCount:  11,
			AlbumCount: 1,
			Value:      "Dance",
		},
		{
			SongCount:  12,
			AlbumCount: 1,
			Value:      "Electronic",
		},
		{
			SongCount:  20,
			AlbumCount: 1,
			Value:      "Hip-Hop",
		},
	}
	WriteJSON(w, http.StatusOK, res)
}

// TODO: stop returning placeholders
func (s Server) HandleGetCoverArt(w http.ResponseWriter, r *http.Request) {
	imageURL := "https://placehold.co/400x400"

	resp, err := http.Get(imageURL)
	if err != nil {
		http.Error(w, "failed to fetch image", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "image unavailable", resp.StatusCode)
		return
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	_, _ = io.Copy(w, resp.Body)
}

// TODO: stop returning placeholders
func (s Server) HandleGetAlbum(w http.ResponseWriter, r *http.Request) {
	res := NewEmptyResponse(s.cfg)
	res.SubsonicResponse.Album = &Album{
		Id:        "200000021",
		Parent:    "100000036",
		Album:     "Forget and Remember",
		Title:     "Forget and Remember",
		Name:      "Forget and Remember",
		IsDir:     true,
		CoverArt:  "al-200000021",
		SongCount: 20,
		Created:   "2021-07-22T02:09:31+00:00",
		Duration:  4248,
		PlayCount: 0,
		ArtistId:  "100000036",
		Artist:    "Comfort Fit",
		Year:      2005,
		Genre:     "Hip-Hop",
		Song: []Song{
			{
				Id:           "300000116",
				Parent:       "200000021",
				Title:        "Can I Help U?",
				IsDir:        false,
				IsVideo:      false,
				Type:         "music",
				AlbumId:      "200000021",
				Album:        "Forget and Remember",
				ArtistId:     "100000036",
				Artist:       "Comfort Fit",
				CoverArt:     "300000116",
				Duration:     103,
				BitRate:      216,
				BitDepth:     16,
				SamplingRate: 44100,
				ChannelCount: 2,
				Track:        1,
				Year:         2005,
				Genre:        "Hip-Hop",
				Size:         2811819,
				DiscNumber:   1,
				Suffix:       "mp3",
				ContentType:  "audio/mpeg",
				Path:         "user/Comfort Fit/Forget And Remember/1 - Can I Help U?.mp3",
			},
		},
	}
	WriteJSON(w, http.StatusOK, res)
}
