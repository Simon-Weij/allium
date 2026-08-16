-- name: GetUser :one
SELECT DISTINCT * FROM users
WHERE username = ?;
