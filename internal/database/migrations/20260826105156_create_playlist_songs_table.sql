-- +goose Up
CREATE TABLE playlist_songs (
    playlist_id TEXT NOT NULL,
    song_id TEXT NOT NULL,
    position INTEGER NOT NULL DEFAULT 0,

    PRIMARY KEY (playlist_id, song_id),

    FOREIGN KEY (playlist_id) REFERENCES playlists(id) ON DELETE CASCADE,
    FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE playlist_songs;

