package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Simon-Weij/allium/internal/resolver"
	"github.com/Simon-Weij/allium/internal/subsonic"
)

type SearchQueries struct {
	ArtistCount int
	SongCount   int
	AlbumCount  int
	SearchCount int
	res         *resolver.ITunesResponse
}

type ResultTypes struct {
	songs   []subsonic.Song
	artists []subsonic.Artist
	albums  []subsonic.SearchResult3Album
}

const defaultSearchLimit = 20

func (s Server) HandleSearch3(w http.ResponseWriter, r *http.Request) {
	queries := parseQueries(w, r, s.iTunesClient)
	if queries == nil {
		return
	}

	resultTypes := ConvertItunesOpenSubsonic(queries.res.Results)

	resultTypes.songs = trimToLimit(resultTypes.songs, queries.SongCount)
	resultTypes.artists = trimToLimit(resultTypes.artists, queries.ArtistCount)
	resultTypes.albums = trimToLimit(resultTypes.albums, queries.AlbumCount)

	res := subsonic.NewEmptyResponse(s.cfg)
	res.SubsonicResponse.SearchResult3 = &subsonic.SearchResult3{
		Artist: resultTypes.artists,
		Song:   resultTypes.songs,
		Album:  resultTypes.albums,
	}
	subsonic.WriteJSON(w, http.StatusOK, res)
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

func parseQueries(w http.ResponseWriter, r *http.Request, client iTunesClient) *SearchQueries {
	query := r.URL.Query()

	searchQuery := query.Get("query")
	if searchQuery == "" {
		http.Error(w, "bad request", http.StatusBadRequest)

		return nil
	}

	var (
		queries SearchQueries
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

func ConvertItunesOpenSubsonic(results []resolver.ITunesResult) ResultTypes {
	songs := []subsonic.Song{}
	artists := []subsonic.Artist{}
	albums := []subsonic.SearchResult3Album{}

	for _, result := range results {
		switch result.WrapperType {
		case itunesWrapperTrack:
			songs = append(songs, convertItunesSong(result))
		case itunesWrapperCollection:
			albums = append(albums, subsonic.SearchResult3Album{
				Id:         strconv.Itoa(result.CollectionID),
				Name:       result.CollectionName,
				Artist:     result.ArtistName,
				Year:       parseYear(result.ReleaseDate),
				CoverArt:   result.ArtworkURL100,
				Starred:    result.ReleaseDate,
				Duration:   albumDurationPlaceholder,
				PlayCount:  albumPlayCountPlaceholder,
				Played:     playedPlaceholderDate,
				Created:    result.ReleaseDate,
				ArtistId:   strconv.Itoa(result.ArtistID),
				UserRating: userRatingPlaceholder,
				SongCount:  result.TrackCount,
			})
		case itunesWrapperArtist:
			artists = append(artists, subsonic.Artist{
				Id:             strconv.Itoa(result.ArtistID),
				Name:           result.ArtistName,
				CoverArt:       result.ArtworkURL100,
				AlbumCount:     artistAlbumCountPlaceholder,
				UserRating:     userRatingPlaceholder,
				ArtistImageUrl: result.ArtworkURL100,
			})

		default:
			slog.Error("unknown type", "result", result)
		}
	}

	return ResultTypes{
		songs:   songs,
		albums:  albums,
		artists: artists,
	}
}
