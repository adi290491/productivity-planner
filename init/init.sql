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

-- Index for fast lookup of active sessions
CREATE INDEX IF NOT EXISTS idx_active_session ON sessions(user_id) WHERE end_time IS NULL;
