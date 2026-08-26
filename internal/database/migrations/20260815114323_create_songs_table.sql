-- +goose Up
CREATE TABLE songs (
    id INTEGER PRIMARY KEY UNIQUE,
    title TEXT NOT NULL,
    artist TEXT NOT NULL,
    album TEXT NOT NULL DEFAULT '',
    album_id INTEGER NOT NULL DEFAULT 0,
    artist_id INTEGER NOT NULL DEFAULT 0,
    duration INTEGER NOT NULL DEFAULT 0,
    track INTEGER NOT NULL DEFAULT 0,
    year INTEGER NOT NULL DEFAULT 0,
    genre TEXT NOT NULL DEFAULT '',
    disc_number INTEGER NOT NULL DEFAULT 0,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    playcount INTEGER NOT NULL,
    user TEXT NOT NULL, -- TODO: make foreign key when there's an users table
    artwork_url TEXT NOT NULL DEFAULT ''
);

-- +goose Down
DROP TABLE songs;
