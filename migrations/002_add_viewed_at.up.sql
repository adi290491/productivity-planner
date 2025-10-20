ALTER TABLE user_daily_trends ADD COLUMN viewed_at TIMESTAMPTZ;
ALTER TABLE user_weekly_trends ADD COLUMN viewed_at TIMESTAMPTZ;
