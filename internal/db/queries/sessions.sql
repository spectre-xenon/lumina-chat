-- name: CreateSession :one
INSERT INTO sessions (session_token, user_id, expires_at)
VALUES (uuid_generate_v4 (), $1, $2)
RETURNING session_token;

-- name: GetSession :one
SELECT session_token, user_id, expires_at
FROM sessions
WHERE session_token = $1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE session_token = $1;

-- name: DeleteSessionsByUser :exec
DELETE FROM sessions
WHERE user_id = $1;

-- name: GetSessionsByUser :many
SELECT session_token, user_id, expires_at
FROM sessions
WHERE user_id = $1;

