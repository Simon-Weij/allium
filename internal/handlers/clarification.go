package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/Simon-Weij/allium/internal/errors"
	"github.com/Simon-Weij/allium/internal/metadata"
)

type Queries struct {
	ArtistCount int
	SongCount   int
	AlbumCount  int
	SearchCount int
	res         *metadata.ITunesResponse
}

const defaultSearchLimit = 20

func (s Server) HandleSearch3(w http.ResponseWriter, r *http.Request) {
	queries := parseQueries(w, r, s.iTunesclient)
	if queries == nil {
		return
	}

	resultTypes := ConvertItunesOpenSubsonic(queries.res.Results)

	resultTypes.songs = trimToLimit(resultTypes.songs, queries.SongCount)
	resultTypes.artists = trimToLimit(resultTypes.artists, queries.ArtistCount)
	resultTypes.albums = trimToLimit(resultTypes.albums, queries.AlbumCount)

	res := NewEmptyResponse(s.cfg)
	res.SubsonicResponse.SearchResult3 = &SearchResult3{
		Artist: resultTypes.artists,
		Song:   resultTypes.songs,
		Album:  resultTypes.albums,
	}
	WriteJSON(w, http.StatusOK, res)
}

func trimToLimit[T any](items []T, count int) []T {
	if count <= 0 {
		count = 20
	}

	if len(items) > count {
		return items[:count]
	}

	return items
}

func parseQueries(w http.ResponseWriter, r *http.Request, client iTunesClient) *Queries {
	query := r.URL.Query()

	searchQuery := query.Get("query")
	if searchQuery == "" {
		http.Error(w, "bad request", http.StatusBadRequest)

		return nil
	}

	var (
		queries Queries
		err     error
	)

	queries.ArtistCount, err = queryInt(query, "artistCount")
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)

		return nil
	}

	queries.AlbumCount, err = queryInt(query, "albumCount")
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)

		return nil
	}

	queries.SongCount, err = queryInt(query, "songCount")
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)

		return nil
	}

	queries.SearchCount, err = queryInt(query, "searchCount")
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)

		return nil
	}

	queries.res, err = client.SearchWithItunes(searchQuery)
	if err != nil {
		slog.Error("something went wrong while searching with itunes", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return nil
	}

	return &queries
}

func queryInt(query url.Values, key string) (int, error) {
	v := query.Get(key)
	if v == "" {
		return defaultSearchLimit, nil
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("parsing query param %q: %w", key, err)
	}

	return n, nil
}

func (s Server) HandleGetCoverArt(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	id := query.Get("id")
	// TODO: support size

	coverPath, err := s.metadata.GetAlbumCover(id)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	imageBytes, err := os.ReadFile(coverPath)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "image/png")

	_, _ = w.Write(imageBytes)
}

func (s Server) HandleStream(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()

	id := query.Get("id")
	if id == "" {
		http.Error(w, "id is required", errors.ErrParameterMissing)
	}

	res, err := s.metadata.GetSongById(id)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		slog.Error("couldn't get song by id", "error", err)

		return
	}

	song := res.Results[0]

	path, err := s.metadata.DownloadOrGetSong(ctx, song.ArtistName, song.TrackName)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		slog.Error("couldn't download song", "error", err)

		return
	}

	http.ServeFile(w, r, path)
}
