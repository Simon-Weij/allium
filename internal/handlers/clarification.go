package handlers

import (
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Simon-Weij/allium/internal/metadata"
)

type Queries struct {
	ArtistCount int
	SongCount   int
	AlbumCount  int
	res         *metadata.ITunesResponse
}

type ResultTypes struct {
	songs   []Song
	artists []Artist
	albums  []Album
}

func (s Server) HandleSearch3(w http.ResponseWriter, r *http.Request) {
	queries := parseQueries(w, r, s.metadata)

	resultTypes := convertItunesOpenSubsonic(queries.res.Results)

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

func convertItunesOpenSubsonic(results []metadata.ITunesResult) ResultTypes {
	songs := []Song{}
	artists := []Artist{}
	albums := []Album{}

	for _, v := range results {
		switch v.WrapperType {
		case "track":
			songs = append(songs, Song{
				Id:           strconv.Itoa(v.TrackID),
				Parent:       strconv.Itoa(v.CollectionID),
				Title:        v.TrackName,
				IsDir:        false,
				IsVideo:      false,
				Type:         "music",
				AlbumId:      strconv.Itoa(v.CollectionID),
				Album:        v.CollectionName,
				ArtistId:     strconv.Itoa(v.ArtistID),
				Artist:       v.ArtistName,
				CoverArt:     v.ArtworkURL100,
				Duration:     v.TrackTimeMillis / 1000,
				BitRate:      320, // placeholder
				BitDepth:     24,
				Size:         0,
				SamplingRate: 48000, // placeholder
				ChannelCount: 2,
				Track:        v.TrackNumber,
				Year:         parseYear(v.ReleaseDate),
				Genre:        v.PrimaryGenreName,
				DiscNumber:   v.DiscNumber,
				Suffix:       "mp3", // placeholder
				ContentType:  "audio/mpeg",
				Path:         "/home/alice/Music", // placeholder
			})
		case "collection":
			albums = append(albums, Album{
				Id:         strconv.Itoa(v.CollectionID),
				Name:       v.CollectionName,
				Artist:     v.ArtistName,
				Year:       parseYear(v.ReleaseDate),
				CoverArt:   v.ArtworkURL100,
				Starred:    v.ReleaseDate,          // placeholder
				Duration:   30000,                  // placeholder
				PlayCount:  8,                      // placeholder
				Played:     "2023-03-28T00:45:13Z", // placeholder
				Created:    v.ReleaseDate,
				ArtistId:   strconv.Itoa(v.ArtistID),
				UserRating: 5, // placeholder
				SongCount:  v.TrackCount,
			})
		case "artist":
			artists = append(artists, Artist{
				Id:             strconv.Itoa(v.ArtistID),
				Name:           v.ArtistName,
				CoverArt:       v.ArtworkURL100,
				AlbumCount:     5, // placeholder
				UserRating:     5, // placeholder
				ArtistImageUrl: v.ArtworkURL100,
			})

		default:
			slog.Error("unknown type", "result", v)
		}
	}
	return ResultTypes{
		songs:   songs,
		albums:  albums,
		artists: artists,
	}
}

func parseQueries(w http.ResponseWriter, r *http.Request, metadata *metadata.Metadata) *Queries {
	query := r.URL.Query()
	searchQuery := query.Get("query")

	var (
		queries Queries
		err     error
	)

	defaultLimit := 20

	queries.ArtistCount, err = queryInt(query, "artistCount", defaultLimit)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return nil
	}

	queries.AlbumCount, err = queryInt(query, "albumCount", defaultLimit)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return nil
	}

	queries.SongCount, err = queryInt(query, "songCount", defaultLimit)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return nil
	}

	queries.res, err = metadata.SearchWithItunes(searchQuery)
	if err != nil {
		slog.Error("something went wrong while searching with itunes", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	return &queries
}

func queryInt(query url.Values, key string, def int) (int, error) {
	v := query.Get(key)
	if v == "" {
		return def, nil
	}
	return strconv.Atoi(v)
}

func parseYear(dateString string) int {
	var year int
	t, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		year = 0
	} else {
		year = t.Year()
	}
	return year
}
