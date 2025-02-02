-- Add a user to a chat
-- name: AddChatMember :exec
INSERT INTO chat_members (chat_id, user_id, joined_at)
VALUES ($1, $2, NOW())
ON CONFLICT DO NOTHING;

-- Remove a user from a chat
-- name: RemoveChatMember :exec
DELETE FROM chat_members
WHERE chat_id = $1 AND user_id = $2;

-- Get chat members
-- name: GetChatMembers :many
SELECT u.id, u.username
FROM users u
JOIN chat_members cm ON u.id = cm.user_id
WHERE cm.chat_id = $1;

