CREATE TABLE sessions (
  session_token UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
  user_id UUID NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
