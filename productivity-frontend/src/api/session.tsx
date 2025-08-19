// src/api/session.ts
import type { SessionResponse } from "../types/summary";
import api from "./api";



export const startSession = async (type: string, token: string): Promise<void> => {
  await api.post(`/sessions/v1/start-session`, { type }, {
    headers: { Authorization: `Bearer ${token}` }
  });
};

export const stopSession = async (type: string, token: string): Promise<SessionResponse> => {
  const res = await api.patch(`/sessions/v1/stop-session`, { type }, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return res.data.session;
};
