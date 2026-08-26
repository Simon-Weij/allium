-- name: CreatePlaylist :one
INSERT INTO playlists (title, playcount, user)
VALUES (?, ?, ?)
RETURNING *;

-- name: InsertIntoPlaylist :one
INSERT INTO playlist_songs (playlist_id, song_id, position)
VALUES (?, ?, ?)
RETURNING *;

-- name: UpdatePlaylistTitle :one
UPDATE playlists
SET title = ?
WHERE id = ?
RETURNING *;

-- name: UpdatePlaylistMetadata :one
UPDATE playlists
SET title = ?,
    comment = ?,
    public = ?
WHERE id = ?
RETURNING *;

-- name: GetPlaylistsByUser :many
SELECT playlists.*,
       COUNT(playlist_songs.song_id) AS song_count
FROM playlists
LEFT JOIN playlist_songs ON playlist_songs.playlist_id = CAST(playlists.id AS TEXT)
WHERE playlists.user = ?
GROUP BY playlists.id;

-- name: GetPlaylistByID :one
SELECT * FROM playlists
WHERE id = ?;

-- name: GetSongsByPlaylistID :many
SELECT songs.*
FROM songs
JOIN playlist_songs ON playlist_songs.song_id = CAST(songs.id AS TEXT)
WHERE playlist_songs.playlist_id = ?
ORDER BY playlist_songs.position;

-- name: DeleteSongFromPlaylistByPosition :exec
DELETE FROM playlist_songs
WHERE playlist_id = ? AND position = ?;

-- name: GetMaxPlaylistPosition :one
SELECT CAST(COALESCE(MAX(position), -1) AS INTEGER) AS max_position
FROM playlist_songs
WHERE playlist_id = ?;
