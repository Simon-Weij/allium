package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Simon-Weij/allium/internal/subsonic"
)

const (
	defaultAlbumListSize = 10
	maxAlbumListSize     = 500
)

func (s Server) HandleGetAlbumList2(w http.ResponseWriter, r *http.Request) {
	listType := r.URL.Query().Get("type")
	if listType == "" {
		subsonic.WriteError(w, http.StatusBadRequest, s.cfg, subsonic.ErrParameterMissing, "missing required parameter 'type'")
		return
	}

	query := r.URL.Query()

	size, err := queryAlbumListSize(query)
	if err != nil {
		subsonic.WriteError(w, http.StatusBadRequest, s.cfg, subsonic.ErrGeneric, err.Error())
		return
	}

	offset, err := queryAlbumListOffset(query)
	if err != nil {
		subsonic.WriteError(w, http.StatusBadRequest, s.cfg, subsonic.ErrGeneric, err.Error())
		return
	}

	var albums []subsonic.AlbumID3

	switch listType {
	case "random", "newest", "highest", "frequent", "recent", "alphabeticalByArtist", "starred", "byGenre", "byYear":
		const searchByName = false

		albums, err = s.iTunesClient.SearchAlbums(size, offset, searchByName)
	case "alphabeticalByName":
		const searchByName = true

		albums, err = s.iTunesClient.SearchAlbums(size, offset, searchByName)
	default:
		albums = []subsonic.AlbumID3{}
	}

	if err != nil {
		slog.Error("could not fetch albums", "type", listType, "error", err)
		subsonic.WriteError(w, http.StatusInternalServerError, s.cfg, subsonic.ErrGeneric, "could not fetch albums")

		return
	}

	res := subsonic.NewEmptyResponse(s.cfg)
	res.SubsonicResponse.AlbumList2 = &subsonic.AlbumList2{Album: albums}

	subsonic.WriteJSON(w, http.StatusOK, res)
}

func queryAlbumListSize(query url.Values) (int, error) {
	size, err := queryIntParam(query, "size", defaultAlbumListSize)
	if err != nil {
		return 0, err
	}

	if size > maxAlbumListSize {
		return 0, fmt.Errorf("size %d exceeds maximum of %d", size, maxAlbumListSize)
	}

	return size, nil
}

func queryAlbumListOffset(query url.Values) (int, error) {
	return queryIntParam(query, "offset", 0)
}

func queryIntParam(query url.Values, key string, def int) (int, error) {
	v := query.Get(key)
	if v == "" {
		return def, nil
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("parsing query param %q: %w", key, err)
	}

	return n, nil
}
