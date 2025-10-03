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
    
    -- Trend metrics
    total_time INTEGER NOT NULL DEFAULT 0,
    focus_minutes INTEGER NOT NULL DEFAULT 0,
    meeting_minutes INTEGER NOT NULL DEFAULT 0,
    break_minutes INTEGER NOT NULL DEFAULT 0,
    
    -- Additional useful metrics
    productivity_score DECIMAL(5,2), -- Optional: calculated score
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    -- ✅ IDEMPOTENCY: Unique constraint prevents duplicates
    CONSTRAINT unique_user_daily_trend UNIQUE (user_id, day)
);


CREATE TABLE IF NOT EXISTS user_weekly_trends (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    week_start DATE NOT NULL,
    
    -- Trend metrics
    total_time INTEGER NOT NULL DEFAULT 0,
    focus_minutes INTEGER NOT NULL DEFAULT 0,
    meeting_minutes INTEGER NOT NULL DEFAULT 0,
    break_minutes INTEGER NOT NULL DEFAULT 0,
    
    -- Additional useful metrics
    productivity_score DECIMAL(5,2), -- Optional: calculated score
    avg_daily_focus INTEGER, -- Average focus per day
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    -- ✅ IDEMPOTENCY: Unique constraint prevents duplicates
    CONSTRAINT unique_user_weekly_trend UNIQUE (user_id, week_start)
);

CREATE TABLE user_notifications (
    id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    
    -- Daily trend notification
    has_new_daily_trend BOOLEAN DEFAULT FALSE,
    last_daily_trend_date DATE,
    last_daily_trend_id UUID REFERENCES user_daily_trends(id) ON DELETE SET NULL,
    
    -- Weekly trend notification
    has_new_weekly_trend BOOLEAN DEFAULT FALSE,
    last_weekly_trend_date DATE,
    last_weekly_trend_id UUID REFERENCES user_weekly_trends(id) ON DELETE SET NULL,
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index for fast lookup of active sessions
CREATE INDEX IF NOT EXISTS idx_active_session ON sessions(user_id) WHERE end_time IS NULL;

CREATE INDEX IF NOT EXISTS idx_user_daily_trend_range ON user_daily_trend(user_id, day DESC);
CREATE INDEX IF NOT EXISTS idx_user_weekly_trend_range ON user_weekly_trend(user_id, week_start DESC);

-- Index for fast notification checks
CREATE INDEX idx_user_notifications_daily 
    ON user_notifications(user_id) 
    WHERE has_new_daily_trend = TRUE;

CREATE INDEX idx_user_notifications_weekly 
    ON user_notifications(user_id) 
    WHERE has_new_weekly_trend = TRUE;