package handlers

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/Simon-Weij/allium/internal/metadata"
)

type ResultTypes struct {
	songs   []Song
	artists []Artist
	albums  []SearchResult3Album
}

const (
	millisPerSecond = 1000

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

func ConvertItunesAlbum(response *metadata.ITunesResponse) GetAlbumAlbum {
	songs := make([]Song, 0, len(response.Results))

	duration := 0

	for _, result := range response.Results {
		if result.WrapperType != itunesWrapperTrack {
			continue
		}

		song := convertItunesSong(result)
		songs = append(songs, song)
		duration += song.Duration
	}

	song := songs[0]

	return GetAlbumAlbum{
		Id:        song.AlbumId,
		Parent:    song.ArtistId,
		Album:     song.Album,
		Title:     song.Album,
		Name:      song.Album,
		IsDir:     true,
		CoverArt:  song.CoverArt,
		ArtistId:  song.ArtistId,
		Artist:    song.Artist,
		Year:      song.Year,
		Genre:     song.Genre,
		Created:   response.Results[0].ReleaseDate,
		SongCount: len(songs),
		Duration:  duration,
		PlayCount: albumPlayCountPlaceholder,
		Song:      songs,
	}
}

func ConvertItunesArtist(response *metadata.ITunesResponse) GetArtistArtist {
	var (
		artist GetArtistArtist
		albums = []GetArtistAlbum{}
	)

	for _, result := range response.Results {
		switch result.WrapperType {
		case itunesWrapperArtist:
			artist.Id = strconv.Itoa(result.ArtistID)
			artist.Name = result.ArtistName
			artist.CoverArt = result.ArtworkURL100
			artist.UserRating = userRatingPlaceholder
			artist.ArtistImageUrl = result.ArtworkURL100
		case itunesWrapperCollection:
			albums = append(albums, convertItunesArtistAlbum(result))
		default:
			slog.Error("unknown type", "result", result)
		}
	}

	artist.AlbumCount = len(albums)
	artist.Album = albums

	return artist
}

func convertItunesArtistAlbum(result metadata.ITunesResult) GetArtistAlbum {
	return GetArtistAlbum{
		Id:            strconv.Itoa(result.CollectionID),
		Parent:        strconv.Itoa(result.ArtistID),
		Album:         result.CollectionName,
		Title:         result.CollectionName,
		Name:          result.CollectionName,
		IsDir:         true,
		CoverArt:      result.ArtworkURL100,
		SongCount:     result.TrackCount,
		Created:       result.ReleaseDate,
		Duration:      albumDurationPlaceholder,
		PlayCount:     albumPlayCountPlaceholder,
		ArtistId:      strconv.Itoa(result.ArtistID),
		Artist:        result.ArtistName,
		Year:          parseYear(result.ReleaseDate),
		Genre:         result.PrimaryGenreName,
		UserRating:    userRatingPlaceholder,
		AverageRating: 0,
		Starred:       "",
	}
}

func convertItunesSong(result metadata.ITunesResult) Song {
	return Song{
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
	}
}

func ConvertItunesOpenSubsonic(results []metadata.ITunesResult) ResultTypes {
	songs := []Song{}
	artists := []Artist{}
	albums := []SearchResult3Album{}

	for _, result := range results {
		switch result.WrapperType {
		case itunesWrapperTrack:
			songs = append(songs, convertItunesSong(result))
		case itunesWrapperCollection:
			albums = append(albums, SearchResult3Album{
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
