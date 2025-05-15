// src/api/session.ts
import axios from "axios";
import type { SessionResponse } from "../types/summary";

const BASE_URL = "http://localhost:8000";

export const startSession = async (type: string, token: string): Promise<void> => {
  await axios.post(`${BASE_URL}/sessions/v1/start-session`, { type }, {
    headers: { Authorization: `Bearer ${token}` }
  });
};

export const stopSession = async (type: string, token: string): Promise<SessionResponse> => {
  const res = await axios.patch(`${BASE_URL}/sessions/v1/stop-session`, { type }, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return res.data.session;
};
