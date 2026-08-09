package handlers

import "net/http"

func (s Server) HandleGetPlaylists(w http.ResponseWriter, r *http.Request) {
	res := NewEmptyResponse(s.cfg)
	res.SubsonicResponse.Playlists = &[]Playlist{
		{
			Id:        "800000003",
			Name:      "random - admin - private (admin)",
			Owner:     "admin",
			Public:    false,
			Created:   "2021-02-23T04:35:38+00:00",
			Changed:   "2021-02-23T04:35:38+00:00",
			SongCount: 43,
			Duration:  17875,
		},
		{
			Id:        "800000002",
			Name:      "random - admin - public (admin)",
			Owner:     "admin",
			Public:    true,
			Created:   "2021-02-23T04:34:56+00:00",
			Changed:   "2021-02-23T04:34:56+00:00",
			SongCount: 43,
			Duration:  17786,
		},
	}
	WriteJSON(w, http.StatusOK, res)
}
