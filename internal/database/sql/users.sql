-- name: GetUser :one
SELECT DISTINCT * FROM users
WHERE username = ?;

-- name: GetUsers :many
SELECT * FROM users