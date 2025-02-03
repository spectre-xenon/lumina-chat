-- Create a new message
-- name: CreateMessage :one
INSERT INTO messages (chat_id, sender_id, content, sent_at)
VALUES ($1, $2, $3, NOW())
RETURNING *;

-- Get all messages in a chat
-- name: GetChatMessages :many
SELECT id, chat_id, sender_id, content, sent_at
FROM messages
WHERE chat_id = $1
ORDER BY sent_at DESC
LIMIT $2 OFFSET $3;

-- Edit a message
-- name: EditMessage :one
UPDATE messages
SET content = $1, sent_at = NOW()
WHERE id = $2 AND sender_id = $3
RETURNING id, content, sent_at;

-- Delete a message
-- name: DeleteMessage :exec
DELETE FROM messages
WHERE id = $1 AND sender_id = $2 ; -- $3 is admin check

