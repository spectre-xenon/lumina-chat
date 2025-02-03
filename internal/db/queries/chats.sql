-- Create a new chat
-- name: CreateChat :one
INSERT INTO chats (name, picture, invite_link)
VALUES ($1, $2, $3)
RETURNING id, name, picture, invite_link;

-- Get all chats for a user with their latest message (if exists)
-- name: GetUserChats :many
WITH LatestMessages AS (
  SELECT DISTINCT ON (m.chat_id)
    m.*
  FROM messages m
  ORDER BY m.chat_id, m.sent_at DESC  -- Ensure latest message per chat
)
SELECT
  c.*,
  sqlc.embed(messages)
FROM chats c
JOIN chat_members cm ON c.id = cm.chat_id
LEFT JOIN LatestMessages messages
  ON c.id = messages.chat_id 
WHERE
  cm.user_id = $1
ORDER BY
  messages.sent_at DESC NULLS LAST;  -- Chats with no messages appear last

-- Get chat by ID
-- name: GetChatByID :one
SELECT id, name, invite_link, picture
FROM chats
WHERE id = $1;

-- Update chat name
-- name: UpdateChatName :one
UPDATE chats
SET name = $1
WHERE id = $2
RETURNING id, name;

-- Update chat name
-- name: UpdateChatPicture :one
UPDATE chats
SET picture = $1
WHERE id = $2
RETURNING id, picture;

-- Delete a chat
-- name: DeleteChat :exec
DELETE FROM chats
WHERE id = $1;

