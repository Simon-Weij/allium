package handlers

import (
	"net/http"
)

// TODO: stop returning placeholders
func (s Server) HandleGetAlbumList2(w http.ResponseWriter, r *http.Request) {
	res := NewEmptyResponse(s.cfg)
	res.SubsonicResponse.AlbumList2 = &AlbumList2{
		Album: []AlbumList2Album{
			{
				Id:        "200000021",
				Album:     "Forget and Remember",
				Title:     "Forget and Remember",
				Name:      "Forget and Remember",
				CoverArt:  "al-200000021",
				SongCount: 20,
				Created:   "2021-07-22T02:09:31+00:00",
				Duration:  4248,
				PlayCount: 0,
				ArtistId:  "100000036",
				Artist:    "Comfort Fit",
				Year:      2005,
				Genre:     "Hip-Hop",
			},
		},
	}
	WriteJSON(w, http.StatusOK, res)
}
