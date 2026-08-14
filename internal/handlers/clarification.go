package handlers

// TODO: more tests

import (
	"fmt"
	"log"
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
	res         *metadata.ITunesResponse
}

const defaultSearchLimit = 20 // The default limit for number of items returned for each artist, album and song by the server

// HandleSearch3 handles the search3 endpoint client requests
// It parses the client query parameters and trims the search results to the search default limit
// then, it writes the search data in JSON format as a response to the client.
func (s Server) HandleSearch3(w http.ResponseWriter, r *http.Request) {
	queries := parseQueries(w, r, s.metadata)

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

// trimToLimit takes an argument of type []T (underlying type of []any), 
// it trims the number of item (T) in items ([]T) to 20 (matching the server's defaultSearchLimit)
func trimToLimit[T any](items []T, count int) []T {
	if count <= 0 {
		count = 20
	}

	if len(items) > count {
		return items[:count]
	}

	return items
}

// convertItunesOpenSubsonic converts data returned by iTunes into OpenSubsonic compatible data
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

// parseQueries parses the queries from the client's request URL, 
// It extracts all the keys and value pairs from the query and formats them in a Queries object
// Populating empty keys with the default values. 
// Read https://opensubsonic.netlify.app/docs/endpoints/search3/ for the default values
func parseQueries(w http.ResponseWriter, r *http.Request, metadata *metadata.Metadata) *Queries {
	query := r.URL.Query()
	log.Println(query)
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

// queryInt converts a client's query value to an integer and returns it, 
// if a value for the query key does not exist, then it returns the default search limit (20)
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

// parseYear takes a date in string form and parses it as per the RFC339 standard layout
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

// HandleGetCoverArt handles the getCoverArt request from the client
// it extracts the id from the client's request and requests CoverArt from iTunes
// which is then saved in the coverPath, the server then writes the image as a response to the client.
func (s Server) HandleGetCoverArt(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	id := query.Get("id")
	// TODO: Support size parameter

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
