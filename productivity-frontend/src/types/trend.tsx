export interface Breakdown {
    [type: string]: string;
}

export interface DailyTrend {
    date: string;
    total_time: string;
    breakdown: Breakdown;
}

export interface DailyTrendResponse {
    user_id: string;
    daily_trends: DailyTrend[];
}

export interface WeeklyTrend {
    week_start: string;
    total_time: string;
    breakdown: Breakdown;
    daily_data: DailyTrend[];
    avg_session_length: string;
    longest_streak: number;
}

export interface WeeklyTrendResponse {
    user_id: string;
    weekly_trends: WeeklyTrend[];
}
