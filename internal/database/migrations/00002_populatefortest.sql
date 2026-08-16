-- +goose Up
INSERT INTO users (username,email)
VALUES ("alice", "alice@example.com");


-- +goose Down
DELETE FROM users;
