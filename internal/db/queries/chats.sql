-- Create a new chat
-- name: CreateChat :one
INSERT INTO chats (name, picture, invite_link)
VALUES ($1, $2, $3)
RETURNING id, name, picture, invite_link;

-- Get all chats for a user
-- name: GetUserChats :many
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
    m.sent_at DESC NULLS LAST;


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

