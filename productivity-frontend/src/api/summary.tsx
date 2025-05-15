import axios from "axios";
import type { DailySummary, WeeklySummaryResponse } from "../types/summary";

const BASE_URL = "http://localhost:8000";

export const fetchDailySummary = async (token: string): Promise<DailySummary> => {
  const res = await axios.get(`${BASE_URL}/summary/daily`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return res.data;
};

export const fetchWeeklySummary = async (token: string): Promise<WeeklySummaryResponse> => {
  const res = await axios.get(`${BASE_URL}/summary/weekly`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return res.data;
};
