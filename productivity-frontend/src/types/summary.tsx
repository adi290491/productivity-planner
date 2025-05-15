export interface Breakdown {
    [type: string]: string;
}

export interface DailySummary {
    date: string;
    total_time: string;
    breakdown: Breakdown;
  }
  
export interface WeeklySummaryResponse {
    start_date: string;
    end_date: string;
    total_time: string;
    daily_summaries: DailySummary[];
}

export interface SessionResponse {
    sessionId: string
    type: string
    start_time: string
    end_time: string;
}