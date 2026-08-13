package handlers

// TODO: more tests

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/Simon-Weij/allium/internal/metadata"
)

type (
	Queries struct {
		ArtistCount int
		SongCount   int
		AlbumCount  int
		res         *metadata.ITunesResponse
	}
	ResultTypes struct {
		songs   []Song
		artists []Artist
		albums  []Album
	}
)

const (
	defaultSearchLimit = 20
	millisPerSecond    = 1000

	itunesWrapperTrack      = "track"
	itunesWrapperCollection = "collection"
	itunesWrapperArtist     = "artist"

	songType         = "music"
	songSuffix       = "mp3"
	songContentType  = "audio/mpeg"
	songPath         = "/home/alice/Music"
	songBitRate      = 320
	songBitDepth     = 24
	songSamplingRate = 48000
	songChannelCount = 2

	albumDurationPlaceholder    = 30000
	albumPlayCountPlaceholder   = 8
	artistAlbumCountPlaceholder = 5
	userRatingPlaceholder       = 5
	playedPlaceholderDate       = "2023-03-28T00:45:13Z"
)

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

	for _, result := range results {
		switch result.WrapperType {
		case itunesWrapperTrack:
			songs = append(songs, Song{
				Id:           strconv.Itoa(result.TrackID),
				Parent:       strconv.Itoa(result.CollectionID),
				Title:        result.TrackName,
				IsDir:        false,
				IsVideo:      false,
				Type:         songType,
				AlbumId:      strconv.Itoa(result.CollectionID),
				Album:        result.CollectionName,
				ArtistId:     strconv.Itoa(result.ArtistID),
				Artist:       result.ArtistName,
				CoverArt:     result.ArtworkURL100,
				Duration:     result.TrackTimeMillis / millisPerSecond,
				BitRate:      songBitRate,
				BitDepth:     songBitDepth,
				Size:         0,
				SamplingRate: songSamplingRate,
				ChannelCount: songChannelCount,
				Track:        result.TrackNumber,
				Year:         parseYear(result.ReleaseDate),
				Genre:        result.PrimaryGenreName,
				DiscNumber:   result.DiscNumber,
				Suffix:       songSuffix,
				ContentType:  songContentType,
				Path:         songPath,
			})
		case itunesWrapperCollection:
			albums = append(albums, Album{
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
			artists = append(artists, Artist{
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

func parseQueries(w http.ResponseWriter, r *http.Request, metadata *metadata.Metadata) *Queries {
	query := r.URL.Query()
	searchQuery := query.Get("query")

	var (
		queries Queries
		err     error
	)

	queries.ArtistCount, err = queryInt(query, "artistCount", defaultSearchLimit)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)

		return nil
	}

	queries.AlbumCount, err = queryInt(query, "albumCount", defaultSearchLimit)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)

		return nil
	}

	queries.SongCount, err = queryInt(query, "songCount", defaultSearchLimit)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)

		return nil
	}

	queries.res, err = metadata.SearchWithItunes(searchQuery)
	if err != nil {
		slog.Error("something went wrong while searching with itunes", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return nil
	}

	return &queries
}

func queryInt(query url.Values, key string, def int) (int, error) {
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

func (s Server) HandleGetCoverArt(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	id := query.Get("id")
	// Not supporting Size for now, Its better to cache a larger image than multiple smaller ones, reduces strain on iTunes & allium

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
