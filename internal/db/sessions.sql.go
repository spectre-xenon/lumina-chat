// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: sessions.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createSession = `-- name: CreateSession :one
INSERT INTO sessions (session_token, user_id, expires_at)
VALUES (uuid_generate_v4 (), $1, $2)
RETURNING session_token
`

type CreateSessionParams struct {
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createSession, arg.UserID, arg.ExpiresAt)
	var session_token uuid.UUID
	err := row.Scan(&session_token)
	return session_token, err
}

const deleteSession = `-- name: DeleteSession :exec
DELETE FROM sessions
WHERE session_token = $1
`

func (q *Queries) DeleteSession(ctx context.Context, sessionToken uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteSession, sessionToken)
	return err
}

const deleteSessionsByUser = `-- name: DeleteSessionsByUser :exec
DELETE FROM sessions
WHERE user_id = $1
`

func (q *Queries) DeleteSessionsByUser(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteSessionsByUser, userID)
	return err
}

const getSession = `-- name: GetSession :one
SELECT session_token, user_id, expires_at
FROM sessions
WHERE session_token = $1
`

func (q *Queries) GetSession(ctx context.Context, sessionToken uuid.UUID) (Session, error) {
	row := q.db.QueryRow(ctx, getSession, sessionToken)
	var i Session
	err := row.Scan(&i.SessionToken, &i.UserID, &i.ExpiresAt)
	return i, err
}

const getSessionsByUser = `-- name: GetSessionsByUser :many
SELECT session_token, user_id, expires_at
FROM sessions
WHERE user_id = $1
`

func (q *Queries) GetSessionsByUser(ctx context.Context, userID uuid.UUID) ([]Session, error) {
	rows, err := q.db.Query(ctx, getSessionsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Session
	for rows.Next() {
		var i Session
		if err := rows.Scan(&i.SessionToken, &i.UserID, &i.ExpiresAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
