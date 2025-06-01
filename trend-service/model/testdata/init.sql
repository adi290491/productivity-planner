-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    name TEXT,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
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

-- Insert dummy users
INSERT INTO users (id, email, name, password_hash)
VALUES
  ('11111111-1111-1111-1111-111111111111', 'alice@example.com', 'Alice', 'hashed_password_1'),
  ('22222222-2222-2222-2222-222222222222', 'bob@example.com', 'Bob', 'hashed_password_2');

-- Insert dummy daily trends for current and previous days
INSERT INTO user_daily_trends (user_id, day, total_time, focus_minutes, meeting_minutes, break_minutes)
VALUES
  ('11111111-1111-1111-1111-111111111111', CURRENT_DATE, 120, 90, 20, 10),
  ('11111111-1111-1111-1111-111111111111', CURRENT_DATE - INTERVAL '1 day', 100, 60, 30, 10),
  ('22222222-2222-2222-2222-222222222222', CURRENT_DATE, 80, 50, 20, 10);

-- Insert dummy weekly trends for this and previous week
INSERT INTO user_weekly_trends (user_id, week_start, total_time, focus_minutes, meeting_minutes, break_minutes)
VALUES
  ('11111111-1111-1111-1111-111111111111', date_trunc('week', CURRENT_DATE), 300, 200, 70, 30),
  ('11111111-1111-1111-1111-111111111111', date_trunc('week', CURRENT_DATE) - INTERVAL '7 days', 250, 150, 70, 30),
  ('22222222-2222-2222-2222-222222222222', date_trunc('week', CURRENT_DATE), 180, 100, 50, 30);
