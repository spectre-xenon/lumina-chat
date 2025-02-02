// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: chats.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createChat = `-- name: CreateChat :one
INSERT INTO chats (name, picture, invite_link)
VALUES ($1, $2, $3)
RETURNING id, name, invite_link
`

type CreateChatParams struct {
	Name       string  `json:"name"`
	Picture    *string `json:"picture"`
	InviteLink *string `json:"invite_link"`
}

type CreateChatRow struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	InviteLink *string   `json:"invite_link"`
}

// Create a new chat
func (q *Queries) CreateChat(ctx context.Context, arg CreateChatParams) (CreateChatRow, error) {
	row := q.db.QueryRow(ctx, createChat, arg.Name, arg.Picture, arg.InviteLink)
	var i CreateChatRow
	err := row.Scan(&i.ID, &i.Name, &i.InviteLink)
	return i, err
}

const deleteChat = `-- name: DeleteChat :exec
DELETE FROM chats
WHERE id = $1
`

// Delete a chat
func (q *Queries) DeleteChat(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteChat, id)
	return err
}

const getChatByID = `-- name: GetChatByID :one
SELECT id, name, invite_link, picture
FROM chats
WHERE id = $1
`

// Get chat by ID
func (q *Queries) GetChatByID(ctx context.Context, id uuid.UUID) (Chat, error) {
	row := q.db.QueryRow(ctx, getChatByID, id)
	var i Chat
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.InviteLink,
		&i.Picture,
	)
	return i, err
}

const getUserChats = `-- name: GetUserChats :many
SELECT
    c.id,
    c.name,
    c.invite_link,
    c.picture,
    m.id AS last_message_id,
    m.sender_id AS last_message_sender,
    m.content AS last_message_content,
    m.sent_at AS last_message_sent_at
FROM
    chats c
JOIN
    chat_members cm ON c.id = cm.chat_id
LEFT JOIN LATERAL (
    SELECT id, sender_id, content, sent_at
    FROM messages
    WHERE chat_id = c.id
    ORDER BY sent_at DESC
    LIMIT 1
) m ON true
WHERE
    cm.user_id = $1
ORDER BY
    m.sent_at DESC NULLS LAST
`

type GetUserChatsRow struct {
	ID                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	InviteLink         *string   `json:"invite_link"`
	Picture            *string   `json:"picture"`
	LastMessageID      int64     `json:"last_message_id"`
	LastMessageSender  uuid.UUID `json:"last_message_sender"`
	LastMessageContent string    `json:"last_message_content"`
	LastMessageSentAt  time.Time `json:"last_message_sent_at"`
}

// Get all chats for a user
func (q *Queries) GetUserChats(ctx context.Context, userID uuid.UUID) ([]GetUserChatsRow, error) {
	rows, err := q.db.Query(ctx, getUserChats, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserChatsRow
	for rows.Next() {
		var i GetUserChatsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.InviteLink,
			&i.Picture,
			&i.LastMessageID,
			&i.LastMessageSender,
			&i.LastMessageContent,
			&i.LastMessageSentAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateChatName = `-- name: UpdateChatName :one
UPDATE chats
SET name = $1
WHERE id = $2
RETURNING id, name
`

type UpdateChatNameParams struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
}

type UpdateChatNameRow struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// Update chat name
func (q *Queries) UpdateChatName(ctx context.Context, arg UpdateChatNameParams) (UpdateChatNameRow, error) {
	row := q.db.QueryRow(ctx, updateChatName, arg.Name, arg.ID)
	var i UpdateChatNameRow
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const updateChatPicture = `-- name: UpdateChatPicture :one
UPDATE chats
SET picture = $1
WHERE id = $2
RETURNING id, picture
`

type UpdateChatPictureParams struct {
	Picture *string   `json:"picture"`
	ID      uuid.UUID `json:"id"`
}

type UpdateChatPictureRow struct {
	ID      uuid.UUID `json:"id"`
	Picture *string   `json:"picture"`
}

// Update chat name
func (q *Queries) UpdateChatPicture(ctx context.Context, arg UpdateChatPictureParams) (UpdateChatPictureRow, error) {
	row := q.db.QueryRow(ctx, updateChatPicture, arg.Picture, arg.ID)
	var i UpdateChatPictureRow
	err := row.Scan(&i.ID, &i.Picture)
	return i, err
}
