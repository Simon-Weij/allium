-- name: UpdatePlays :exec
INSERT INTO songs (id, title, artist, album, album_id, artist_id, duration, track, year, genre, disc_number, updated_at, playcount, user, artwork_url)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, 1, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    updated_at = CURRENT_TIMESTAMP,
    playcount = songs.playcount + 1,
    album = excluded.album,
    album_id = excluded.album_id,
    artist_id = excluded.artist_id,
    duration = excluded.duration,
    track = excluded.track,
    year = excluded.year,
    genre = excluded.genre,
    disc_number = excluded.disc_number,
    artwork_url = excluded.artwork_url;
