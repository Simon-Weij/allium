-- +goose Up
CREATE TABLE users (
    username text NOT NULL,
    email text NOT NULL,
    scrobbling_enabled bool NOT NULL DEFAULT false,
    admin_role bool NOT NULL DEFAULT false,
    settings_role bool NOT NULL DEFAULT false,
    download_role bool NOT NULL DEFAULT false,
    upload_role bool NOT NULL DEFAULT false,
    playlist_role bool NOT NULL DEFAULT false,
    cover_art_role bool NOT NULL DEFAULT false,
    comment_role bool NOT NULL DEFAULT false,
    podcast_role bool NOT NULL DEFAULT false,
    stream_role bool NOT NULL DEFAULT false,
    jukebox_role bool NOT NULL DEFAULT false,
    share_role bool NOT NULL DEFAULT false,
    video_conversion_role bool NOT NULL DEFAULT false,
    avatar_last_changed text,
    folder text,
    max_bit_Rate int
);

-- +goose Down
DROP TABLE users;
