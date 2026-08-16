-- +goose Up
CREATE TABLE songs(
    id INTEGER PRIMARY KEY UNIQUE,
    title TEXT NOT NULL,
    artist TEXT NOT NULL,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    playcount INTEGER NOT NULL,
    user TEXT NOT NULL -- TODO: make foreign key when there's an users table
);

-- +goose Down
DROP TABLE songs;
