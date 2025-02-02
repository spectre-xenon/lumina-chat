CREATE TABLE chats (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
  name TEXT NOT NULL,
  invite_link TEXT UNIQUE,
  picture TEXT
);
