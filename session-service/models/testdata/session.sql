-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    name TEXT,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Sessions table
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    session_type TEXT NOT NULL CHECK (session_type IN ('focus', 'break', 'meeting')),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP
);

-- User for CreateSession test
INSERT INTO users (id, email, name, password_hash, created_at)
VALUES (
  'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
  'focususer@example.com',
  'Focus User',
  'hashed_password_1',
  NOW()
);

-- User for StopSession test
INSERT INTO users (id, email, name, password_hash, created_at)
VALUES (
  'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
  'meetinguser@example.com',
  'Meeting User',
  'hashed_password_2',
  NOW()
);

-- User for StopSession negative test
INSERT INTO users (id, email, name, password_hash, created_at)
VALUES (
  'cccccccc-cccc-cccc-cccc-cccccccccccc',
  'noactivesession@example.com',
  'No Active',
  'hashed_password_3',
  NOW()
);