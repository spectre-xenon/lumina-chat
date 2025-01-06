-- Enable uuid EXTENSION
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
  username TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT, -- NULL for OAuth users
  created_at TIMESTAMPTZ DEFAULT NOW () NOT NULL
);
