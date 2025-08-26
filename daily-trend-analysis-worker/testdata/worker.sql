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

CREATE TABLE IF NOT EXISTS user_daily_trends (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    day DATE NOT NULL,
	total_time INTEGER,
	focus_minutes INTEGER,
	meeting_minutes INTEGER,
	break_minutes INTEGER,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (user_id, day)
);

CREATE TABLE IF NOT EXISTS user_weekly_trends (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    week_start DATE NOT NULL,
	total_time INTEGER,
	focus_minutes INTEGER,
	meeting_minutes INTEGER,
	break_minutes INTEGER,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (user_id, week_start)
);

-- INSERT INTO users (id, email, name, password_hash, created_at) VALUES
-- ('11111111-1111-1111-1111-111111111111', 'alice@example.com', 'Alice', 'hashed_password_1', NOW()),
-- ('22222222-2222-2222-2222-222222222222', 'bob@example.com', 'Bob', 'hashed_password_2', NOW());

-- -- Sessions for Alice (focus, meeting, break spread across different days and times)
-- INSERT INTO sessions (id, user_id, session_type, start_time, end_time) VALUES
-- ('a1f1b4d0-0001-4c11-8f92-111111111111', '11111111-1111-1111-1111-111111111111', 'focus',   NOW() - INTERVAL '2 days 3 hours', NOW() - INTERVAL '2 days 2 hours 30 minutes'),
-- ('a1f1b4d0-0002-4c11-8f92-111111111111', '11111111-1111-1111-1111-111111111111', 'meeting', NOW() - INTERVAL '1 day 6 hours', NOW() - INTERVAL '1 day 5 hours 45 minutes'),
-- ('a1f1b4d0-0003-4c11-8f92-111111111111', '11111111-1111-1111-1111-111111111111', 'break',   NOW() - INTERVAL '1 day 4 hours', NOW() - INTERVAL '1 day 3 hours 45 minutes'),

-- -- Sessions for Bob
-- ('b2f1b4d0-0001-4c11-8f92-222222222222', '22222222-2222-2222-2222-222222222222', 'focus',   NOW() - INTERVAL '3 days 4 hours', NOW() - INTERVAL '3 days 3 hours 30 minutes'),
-- ('b2f1b4d0-0002-4c11-8f92-222222222222', '22222222-2222-2222-2222-222222222222', 'meeting', NOW() - INTERVAL '2 days 6 hours', NOW() - INTERVAL '2 days 5 hours 50 minutes'),
-- ('b2f1b4d0-0003-4c11-8f92-222222222222', '22222222-2222-2222-2222-222222222222', 'break',   NOW() - INTERVAL '1 day 2 hours', NOW() - INTERVAL '1 day 1 hour 45 minutes');

INSERT INTO users (id, email, name, password_hash, created_at) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'alice@example.com', 'Alice', 'hashed_password_1', NOW()),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'bob@example.com', 'Bob', 'hashed_password_2', NOW());

-- Sessions for Alice (3 of each type)
INSERT INTO sessions (id, user_id, session_type, start_time, end_time) VALUES
('a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'focus',   NOW() - INTERVAL '3 days 09:00', NOW() - INTERVAL '3 days 08:30'),
('a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a2', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'focus',   NOW() - INTERVAL '2 days 10:00', NOW() - INTERVAL '2 days 09:30'),
('a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a3', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'focus',   NOW() - INTERVAL '1 day 11:00', NOW() - INTERVAL '1 day 10:30'),
('a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1b1', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'meeting', NOW() - INTERVAL '3 days 14:00', NOW() - INTERVAL '3 days 13:45'),
('a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1b2', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'meeting', NOW() - INTERVAL '2 days 15:00', NOW() - INTERVAL '2 days 14:45'),
('a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1b3', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'meeting', NOW() - INTERVAL '1 day 16:00', NOW() - INTERVAL '1 day 15:45'),
('a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1c1', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'break',   NOW() - INTERVAL '3 days 17:00', NOW() - INTERVAL '3 days 16:50'),
('a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1c2', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'break',   NOW() - INTERVAL '2 days 18:00', NOW() - INTERVAL '2 days 17:50'),
('a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1c3', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'break',   NOW() - INTERVAL '1 day 19:00', NOW() - INTERVAL '1 day 18:50');

-- Sessions for Bob (3 of each type)
INSERT INTO sessions (id, user_id, session_type, start_time, end_time) VALUES
('b1b1b1b1-b1b1-b1b1-b1b1-b1b1b1b1b1b1', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'focus',   NOW() - INTERVAL '3 days 08:00', NOW() - INTERVAL '3 days 07:30'),
('b1b1b1b1-b1b1-b1b1-b1b1-b1b1b1b1b1b2', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'focus',   NOW() - INTERVAL '2 days 09:00', NOW() - INTERVAL '2 days 08:30'),
('b1b1b1b1-b1b1-b1b1-b1b1-b1b1b1b1b1b3', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'focus',   NOW() - INTERVAL '1 day 10:00', NOW() - INTERVAL '1 day 09:30'),
('b1b1b1b1-b1b1-b1b1-b1b1-b1b1b1b1b1c1', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'meeting', NOW() - INTERVAL '3 days 13:00', NOW() - INTERVAL '3 days 12:45'),
('b1b1b1b1-b1b1-b1b1-b1b1-b1b1b1b1b1c2', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'meeting', NOW() - INTERVAL '2 days 14:00', NOW() - INTERVAL '2 days 13:45'),
('b1b1b1b1-b1b1-b1b1-b1b1-b1b1b1b1b1c3', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'meeting', NOW() - INTERVAL '1 day 15:00', NOW() - INTERVAL '1 day 14:45'),
('b1b1b1b1-b1b1-b1b1-b1b1-b1b1b1b1b1d1', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'break',   NOW() - INTERVAL '3 days 16:00', NOW() - INTERVAL '3 days 15:50'),
('b1b1b1b1-b1b1-b1b1-b1b1-b1b1b1b1b1d2', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'break',   NOW() - INTERVAL '2 days 17:00', NOW() - INTERVAL '2 days 16:50'),
('b1b1b1b1-b1b1-b1b1-b1b1-b1b1b1b1b1d3', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'break',   NOW() - INTERVAL '1 day 18:00', NOW() - INTERVAL '1 day 17:50');