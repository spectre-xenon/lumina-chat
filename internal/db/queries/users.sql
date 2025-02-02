-- name: GetUserByID :one
SELECT id, username, email, password_hash, created_at, picture, color
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, username, email, password_hash, created_at, picture, color
FROM users
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING id, username, email, created_at;

-- name: UpdateUserColor :exec
UPDATE users
SET color = $2
WHERE id = $1;
