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

INSERT INTO users (id, email, name, password_hash, created_at) VALUES
('11111111-1111-1111-1111-111111111111', 'alice@example.com', 'Alice', 'hashed_password_1', NOW()),
('22222222-2222-2222-2222-222222222222', 'bob@example.com', 'Bob', 'hashed_password_2', NOW());

-- Sessions for Alice (focus, meeting, break spread across different days and times)
INSERT INTO sessions (id, user_id, session_type, start_time, end_time) VALUES
('a1f1b4d0-0001-4c11-8f92-111111111111', '11111111-1111-1111-1111-111111111111', 'focus',   NOW() - INTERVAL '2 days 3 hours', NOW() - INTERVAL '2 days 2 hours 30 minutes'),
('a1f1b4d0-0002-4c11-8f92-111111111111', '11111111-1111-1111-1111-111111111111', 'meeting', NOW() - INTERVAL '1 day 6 hours', NOW() - INTERVAL '1 day 5 hours 45 minutes'),
('a1f1b4d0-0003-4c11-8f92-111111111111', '11111111-1111-1111-1111-111111111111', 'break',   NOW() - INTERVAL '1 day 4 hours', NOW() - INTERVAL '1 day 3 hours 45 minutes'),

-- Sessions for Bob
('b2f1b4d0-0001-4c11-8f92-222222222222', '22222222-2222-2222-2222-222222222222', 'focus',   NOW() - INTERVAL '3 days 4 hours', NOW() - INTERVAL '3 days 3 hours 30 minutes'),
('b2f1b4d0-0002-4c11-8f92-222222222222', '22222222-2222-2222-2222-222222222222', 'meeting', NOW() - INTERVAL '2 days 6 hours', NOW() - INTERVAL '2 days 5 hours 50 minutes'),
('b2f1b4d0-0003-4c11-8f92-222222222222', '22222222-2222-2222-2222-222222222222', 'break',   NOW() - INTERVAL '1 day 2 hours', NOW() - INTERVAL '1 day 1 hour 45 minutes');
