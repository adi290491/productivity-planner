import type { DailySummary, WeeklySummaryResponse } from "../types/summary";
import api from "./api";


export const fetchDailySummary = async (token: string): Promise<DailySummary> => {
  const res = await api.get(`/summary/daily`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return res.data;
};

export const fetchWeeklySummary = async (token: string): Promise<WeeklySummaryResponse> => {
  const res = await api.get(`/summary/weekly`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return res.data;
};
