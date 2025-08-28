import type { DailyTrendResponse, WeeklyTrendResponse } from "../types/trend";
import api from "./api";

export const fetchDailyTrends = async (token: string): Promise<DailyTrendResponse> => {
  const res = await api.get(`/trend/daily`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return res.data;
};

export const fetchWeeklyTrends = async (token: string): Promise<WeeklyTrendResponse> => {
  const res = await api.get(`/trend/weekly`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return res.data;
};
