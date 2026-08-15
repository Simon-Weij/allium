-- name: UpdatePlays :exec
INSERT INTO songs (id, title, artist, updated_at, playcount, user)
VALUES (?, ?, ?, CURRENT_TIMESTAMP, 1, ?)
ON CONFLICT(id) DO UPDATE SET
    updated_at = CURRENT_TIMESTAMP,
    playcount = songs.playcount + 1;
