-- name: GetUser :one
SELECT DISTINCT * FROM users
WHERE username = ?;

-- name: GetUsers :many
SELECT * FROM users;

-- name: CreateUser :exec
INSERT INTO users (
    username,
    email,
    scrobbling_enabled,
    admin_role,
    settings_role,
    download_role,
    upload_role,
    playlist_role,
    cover_art_role,
    comment_role,
    podcast_role,
    stream_role,
    jukebox_role,
    share_role,
    video_conversion_role,
    avatar_last_changed,
    folder,
    max_bit_Rate
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
);
