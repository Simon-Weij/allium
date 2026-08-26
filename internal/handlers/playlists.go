package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Simon-Weij/allium/generated/sqlc"
	"github.com/Simon-Weij/allium/internal/config"
	"github.com/Simon-Weij/allium/internal/subsonic"
)

var (
	errSongIdNotNumber     = errors.New("songId must be a number")
	errPlaylistIdNotNumber = errors.New("playlistId must be a number")
	errSongNotFound        = errors.New("could not find song")
)

func (s Server) HandleCreatePlaylist(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ctx := r.Context()

	name := query.Get("name")
	username := query.Get("u")

	playlistId := query.Get("playlistId")
	songIds := query["songId"]

	if name == "" && playlistId == "" {
		http.Error(w, "either name or playlistId is required", http.StatusBadRequest)

		return
	}

	if err := validateSongIds(songIds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	playlist, err := s.resolvePlaylist(ctx, name, username, playlistId)
	if err != nil {
		slog.Error("could not resolve playlist", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	if err := s.addSongsToPlaylist(ctx, playlist, songIds, username); err != nil {
		slog.Error("could not add songs to playlist", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	response := buildPlaylistResponse(s.cfg, playlist, songIds)
	subsonic.WriteJSON(w, http.StatusOK, response)
}

func (s Server) HandleGetPlaylists(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ctx := r.Context()

	username := query.Get("username")
	if username == "" {
		username = query.Get("u")
	}

	rows, err := s.queries.GetPlaylistsByUser(ctx, username)
	if err != nil {
		slog.Error("could not get playlists", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	playlists := make([]subsonic.Playlist, 0, len(rows))
	for _, row := range rows {
		playlists = append(playlists, convertPlaylistRow(row))
	}

	response := subsonic.NewEmptyResponse(s.cfg)
	response.SubsonicResponse.Playlists = &subsonic.Playlists{
		Playlist: playlists,
	}

	subsonic.WriteJSON(w, http.StatusOK, response)
}

func (s Server) HandleGetPlaylist(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ctx := r.Context()

	id := query.Get("id")
	if id == "" {
		http.Error(w, "no id provided", http.StatusBadRequest)

		return
	}

	parsedId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, errPlaylistIdNotNumber.Error(), http.StatusBadRequest)

		return
	}

	playlist, err := s.queries.GetPlaylistByID(ctx, int64(parsedId))
	if err != nil {
		http.Error(w, "could not find playlist with id: "+id, http.StatusNotFound)

		return
	}

	songs, err := s.queries.GetSongsByPlaylistID(ctx, id)
	if err != nil {
		slog.Error("could not get songs for playlist", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	entries := make([]subsonic.Song, 0, len(songs))
	for _, song := range songs {
		entries = append(entries, convertSong(song))
	}

	response := subsonic.NewEmptyResponse(s.cfg)
	response.SubsonicResponse.Playlist = &subsonic.Playlist{
		Id:        strconv.FormatInt(playlist.ID, 10),
		Name:      playlist.Title,
		Owner:     playlist.User,
		Public:    false,
		Created:   playlist.UpdatedAt,
		Changed:   playlist.UpdatedAt,
		SongCount: len(songs),
		Duration:  0,
		Entry:     entries,
	}

	subsonic.WriteJSON(w, http.StatusOK, response)
}

func (s Server) HandleUpdatePlaylist(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ctx := r.Context()

	playlistId := query.Get("playlistId")
	if playlistId == "" {
		http.Error(w, "playlistId is required", http.StatusBadRequest)

		return
	}

	parsedId, err := strconv.Atoi(playlistId)
	if err != nil {
		http.Error(w, errPlaylistIdNotNumber.Error(), http.StatusBadRequest)

		return
	}

	playlist, err := s.queries.GetPlaylistByID(ctx, int64(parsedId))
	if err != nil {
		http.Error(w, "could not find playlist with id: "+playlistId, http.StatusNotFound)

		return
	}

	public := parsePublic(query.Get("public"), playlist.Public)

	playlist, err = s.updatePlaylistFromQuery(ctx, w, query.Get("name"), query.Get("comment"), public, playlist)
	if err != nil {
		return
	}

	songIdsToAdd := query["songIdToAdd"]
	if len(songIdsToAdd) > 0 {
		if err := s.handleAddSongsToUpdate(ctx, w, playlist, songIdsToAdd); err != nil {
			return
		}
	}

	songIndicesToRemove := query["songIndexToRemove"]
	if len(songIndicesToRemove) > 0 {
		if err := s.handleRemoveSongsFromPlaylist(ctx, w, playlistId, songIndicesToRemove); err != nil {
			return
		}
	}

	response := subsonic.NewEmptyResponse(s.cfg)
	subsonic.WriteJSON(w, http.StatusOK, response)
}

func parsePublic(publicStr string, defaultVal int64) int64 {
	if publicStr == "" {
		return defaultVal
	}

	if publicStr == "true" {
		return 1
	}

	return 0
}

func (s Server) updatePlaylistFromQuery(
	ctx context.Context,
	w http.ResponseWriter,
	name, comment string,
	public int64,
	playlist sqlc.Playlist,
) (sqlc.Playlist, error) {
	if name == "" {
		name = playlist.Title
	}

	if comment == "" {
		comment = playlist.Comment
	}

	updated, err := s.queries.UpdatePlaylistMetadata(ctx, sqlc.UpdatePlaylistMetadataParams{
		Title:   name,
		Comment: comment,
		Public:  public,
		ID:      playlist.ID,
	})
	if err != nil {
		slog.Error("could not update playlist metadata", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return sqlc.Playlist{}, fmt.Errorf("could not update playlist metadata: %w", err)
	}

	return updated, nil
}

func (s Server) handleAddSongsToUpdate(
	ctx context.Context,
	w http.ResponseWriter,
	playlist sqlc.Playlist,
	songIdsToAdd []string,
) error {
	if err := validateSongIds(songIdsToAdd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return err
	}

	if err := s.addSongsToPlaylist(ctx, playlist, songIdsToAdd, playlist.User); err != nil {
		slog.Error("could not add songs to playlist", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return err
	}

	return nil
}

func (s Server) handleRemoveSongsFromPlaylist(
	ctx context.Context,
	w http.ResponseWriter,
	playlistId string,
	songIndicesToRemove []string,
) error {
	for _, idxStr := range songIndicesToRemove {
		idx, err := strconv.Atoi(idxStr)
		if err != nil {
			http.Error(w, "songIndexToRemove must be a number", http.StatusBadRequest)

			return fmt.Errorf("songIndexToRemove must be a number: %w", err)
		}

		if err := s.queries.DeleteSongFromPlaylistByPosition(ctx, sqlc.DeleteSongFromPlaylistByPositionParams{
			PlaylistID: playlistId,
			Position:   int64(idx),
		}); err != nil {
			slog.Error("could not remove song from playlist", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)

			return fmt.Errorf("could not remove song from playlist: %w", err)
		}
	}

	return nil
}

func validateSongIds(songIds []string) error {
	for _, songId := range songIds {
		if _, err := strconv.Atoi(songId); err != nil {
			return errSongIdNotNumber
		}
	}

	return nil
}

func (s Server) resolvePlaylist(
	ctx context.Context,
	name, username, playlistId string,
) (sqlc.Playlist, error) {
	if playlistId != "" {
		return s.updateExistingPlaylist(ctx, name, playlistId)
	}

	return s.createNewPlaylist(ctx, name, username)
}

func (s Server) updateExistingPlaylist(
	ctx context.Context,
	name, playlistId string,
) (sqlc.Playlist, error) {
	id, err := strconv.Atoi(playlistId)
	if err != nil {
		return sqlc.Playlist{}, errPlaylistIdNotNumber
	}

	updatedPlaylist, err := s.queries.UpdatePlaylistTitle(ctx, sqlc.UpdatePlaylistTitleParams{
		Title: name,
		ID:    int64(id),
	})
	if err != nil {
		return sqlc.Playlist{}, fmt.Errorf("could not update playlist: %w", err)
	}

	return updatedPlaylist, nil
}

func (s Server) createNewPlaylist(
	ctx context.Context,
	name, username string,
) (sqlc.Playlist, error) {
	createdPlaylist, err := s.queries.CreatePlaylist(ctx, sqlc.CreatePlaylistParams{
		Title:     name,
		Playcount: 0,
		User:      username,
	})
	if err != nil {
		return sqlc.Playlist{}, fmt.Errorf("could not create playlist: %w", err)
	}

	return createdPlaylist, nil
}

func (s Server) addSongsToPlaylist(
	ctx context.Context,
	playlist sqlc.Playlist,
	songIds []string,
	username string,
) error {
	playlistId := strconv.FormatInt(playlist.ID, 10)

	for _, songId := range songIds {
		if err := s.ensureSongExists(ctx, songId, username); err != nil {
			return err
		}

		maxPos, err := s.queries.GetMaxPlaylistPosition(ctx, playlistId)
		if err != nil {
			return fmt.Errorf("could not get playlist position: %w", err)
		}

		position := maxPos
		if position < 0 {
			position = 0
		} else {
			position++
		}

		_, err = s.queries.InsertIntoPlaylist(ctx, sqlc.InsertIntoPlaylistParams{
			PlaylistID: playlistId,
			SongID:     songId,
			Position:   position,
		})
		if err != nil {
			return fmt.Errorf("could not add song to playlist: %w", err)
		}
	}

	return nil
}

func (s Server) ensureSongExists(ctx context.Context, songId, username string) error {
	idInt, err := strconv.Atoi(songId)
	if err != nil {
		return errSongIdNotNumber
	}

	itunesSongs, err := s.iTunesClient.GetSongById(songId)
	if err != nil {
		return fmt.Errorf("could not fetch song metadata: %w", err)
	}

	if len(itunesSongs.Results) == 0 {
		return errSongNotFound
	}

	song := convertItunesSong(itunesSongs.Results[0])

	albumID, err := strconv.Atoi(song.AlbumId)
	if err != nil {
		return fmt.Errorf("could not parse album id: %w", err)
	}

	artistID, err := strconv.Atoi(song.ArtistId)
	if err != nil {
		return fmt.Errorf("could not parse artist id: %w", err)
	}

	if err := s.queries.UpdatePlays(ctx, sqlc.UpdatePlaysParams{
		ID:         int64(idInt),
		Title:      song.Title,
		Artist:     song.Artist,
		Album:      song.Album,
		AlbumID:    int64(albumID),
		ArtistID:   int64(artistID),
		Duration:   int64(song.Duration),
		Track:      int64(song.Track),
		Year:       int64(song.Year),
		Genre:      song.Genre,
		DiscNumber: int64(song.DiscNumber),
		User:       username,
		ArtworkUrl: song.CoverArt,
	}); err != nil {
		return fmt.Errorf("could not store song: %w", err)
	}

	return nil
}

func buildPlaylistResponse(
	cfg config.Config,
	playlist sqlc.Playlist,
	songIds []string,
) subsonic.SubsonicResponseWrapper {
	response := subsonic.NewEmptyResponse(cfg)
	response.SubsonicResponse.Playlist = &subsonic.Playlist{
		Id:        strconv.FormatInt(playlist.ID, 10),
		Name:      playlist.Title,
		Owner:     playlist.User,
		Public:    false,
		Created:   playlist.UpdatedAt,
		Changed:   playlist.UpdatedAt,
		SongCount: len(songIds),
		Duration:  0,
		Entry:     nil,
	}

	return response
}

func convertPlaylistRow(row sqlc.GetPlaylistsByUserRow) subsonic.Playlist {
	return subsonic.Playlist{
		Id:        strconv.FormatInt(row.ID, 10),
		Name:      row.Title,
		Owner:     row.User,
		Public:    false,
		Created:   row.UpdatedAt,
		Changed:   row.UpdatedAt,
		SongCount: int(row.SongCount),
		Duration:  0,
		Entry:     nil,
	}
}

func convertSong(song sqlc.Song) subsonic.Song {
	return subsonic.Song{
		Id:           strconv.FormatInt(song.ID, 10),
		Parent:       strconv.FormatInt(song.AlbumID, 10),
		Title:        song.Title,
		IsDir:        false,
		IsVideo:      false,
		Type:         songType,
		AlbumId:      strconv.FormatInt(song.AlbumID, 10),
		Album:        song.Album,
		ArtistId:     strconv.FormatInt(song.ArtistID, 10),
		Artist:       song.Artist,
		CoverArt:     song.ArtworkUrl,
		Duration:     int(song.Duration),
		BitRate:      songBitRate,
		BitDepth:     songBitDepth,
		SamplingRate: songSamplingRate,
		ChannelCount: songChannelCount,
		Track:        int(song.Track),
		Year:         int(song.Year),
		Genre:        song.Genre,
		Size:         0,
		DiscNumber:   int(song.DiscNumber),
		Suffix:       songSuffix,
		ContentType:  songContentType,
		Path:         "",
	}
}
