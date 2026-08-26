-- +goose Up
CREATE TABLE playlists (
    id INTEGER PRIMARY KEY UNIQUE,
    title TEXT NOT NULL,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    playcount INTEGER NOT NULL,
    user TEXT NOT NULL, -- TODO: make foreign key when there's an users table
    comment TEXT NOT NULL DEFAULT '',
    public INTEGER NOT NULL DEFAULT 0
);
-- +goose Down
DROP TABLE playlists;
