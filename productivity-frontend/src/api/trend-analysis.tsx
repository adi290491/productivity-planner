import type { DailyTrendResponse, WeeklyTrendResponse } from "../types/trend";
import api from "./api";

export const fetchDailyTrends = async (token: string, days: number): Promise<DailyTrendResponse> => {
  const res = await api.get(`/trend/daily`, {
    headers: { Authorization: `Bearer ${token}` },
    params: { days },
  });
  return res.data;
};

export const fetchWeeklyTrends = async (token: string, weeks: number): Promise<WeeklyTrendResponse> => {
  const res = await api.get(`/trend/weekly`, {
    headers: { Authorization: `Bearer ${token}` },
    params: { weeks },
  });
  return res.data;
};

export const fetchLatestTrendsCount = async (token: string, weeks: number): Promise<WeeklyTrendResponse> => {
  const res = await api.get(`/trend/unviewed`, {
    headers: { Authorization: `Bearer ${token}` },
    params: { weeks },
  });
  return res.data;
};