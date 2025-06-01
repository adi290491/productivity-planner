-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    name TEXT,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Insert dummy users
INSERT INTO users (id, email, name, password_hash)
VALUES
  ('11111111-1111-1111-1111-111111111111', 'alice@example.com', 'Alice', 'hashed_password_1'),
  ('22222222-2222-2222-2222-222222222222', 'bob@example.com', 'Bob', 'hashed_password_2');

