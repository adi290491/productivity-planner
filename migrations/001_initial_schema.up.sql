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
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    day DATE NOT NULL,
    total_time INTEGER NOT NULL DEFAULT 0,
    focus_minutes INTEGER NOT NULL DEFAULT 0,
    meeting_minutes INTEGER NOT NULL DEFAULT 0,
    break_minutes INTEGER NOT NULL DEFAULT 0,
    productivity_score DECIMAL(5,2),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT unique_user_daily_trend UNIQUE (user_id, day)
);

CREATE TABLE IF NOT EXISTS user_weekly_trends (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    week_start DATE NOT NULL,
    total_time INTEGER NOT NULL DEFAULT 0,
    focus_minutes INTEGER NOT NULL DEFAULT 0,
    meeting_minutes INTEGER NOT NULL DEFAULT 0,
    break_minutes INTEGER NOT NULL DEFAULT 0,
    productivity_score DECIMAL(5,2),
    avg_daily_focus INTEGER,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT unique_user_weekly_trend UNIQUE (user_id, week_start)
);

CREATE TABLE user_notifications (
    id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    has_new_daily_trend BOOLEAN DEFAULT FALSE,
    last_daily_trend_date DATE,
    last_daily_trend_id UUID REFERENCES user_daily_trends(id) ON DELETE SET NULL,
    has_new_weekly_trend BOOLEAN DEFAULT FALSE,
    last_weekly_trend_date DATE,
    last_weekly_trend_id UUID REFERENCES user_weekly_trends(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_active_session ON sessions(user_id) WHERE end_time IS NULL;
CREATE INDEX IF NOT EXISTS idx_user_daily_trend_range ON user_daily_trends(user_id, day DESC);
CREATE INDEX IF NOT EXISTS idx_user_weekly_trend_range ON user_weekly_trends(user_id, week_start DESC);
CREATE INDEX idx_user_notifications_daily ON user_notifications(user_id) WHERE has_new_daily_trend = TRUE;
CREATE INDEX idx_user_notifications_weekly ON user_notifications(user_id) WHERE has_new_weekly_trend = TRUE;
