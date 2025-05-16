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

CREATE TABLE IF NOT EXISTS user_daily_trend (
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

CREATE TABLE IF NOT EXISTS user_weekly_trend (
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

-- Index for fast lookup of active sessions
CREATE INDEX IF NOT EXISTS idx_active_session ON sessions(user_id) WHERE end_time IS NULL;

CREATE INDEX IF NOT EXISTS idx_user_daily_trend_range ON user_daily_trend(user_id, day DESC);
CREATE INDEX IF NOT EXISTS idx_user_weekly_trend_range ON user_weekly_trend(user_id, week_start DESC);
