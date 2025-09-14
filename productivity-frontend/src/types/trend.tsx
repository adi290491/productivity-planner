export interface Breakdown {
    [type: string]: string;
}

export interface DailyTrend {
    date: string;
    total_time: string;
    breakdown: Breakdown;
}

export interface DailyTrendResponse {
    dailyTrends: DailyTrend[];
}

export interface WeeklyTrend {
    week_start: string;
    total_time: string;
    breakdown: Breakdown;
}

export interface WeeklyTrendResponse {
    weeklyTrends: WeeklyTrend[];
}
