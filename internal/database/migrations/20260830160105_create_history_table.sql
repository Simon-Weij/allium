-- +goose Up
CREATE TABLE history (
    id INTEGER PRIMARY KEY UNIQUE,
    song_id TEXT NOT NULL,
    last_played TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    genre TEXT NOT NULL
);

-- +goose Down
DROP TABLE history;
