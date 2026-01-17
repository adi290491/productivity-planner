// src/api/session.ts
import type { SessionResponse } from "../types/summary";
import api from "./api";



export const startSession = async (type: string, token: string): Promise<void> => {
  await api.post(`/sessions/v1/start-session`, {
    "session_type": type
  }, {
    headers: { Authorization: `Bearer ${token}` }
  });
};

export const  stopSession = async (type: string, token: string): Promise<SessionResponse> => {
  const res = await api.patch(`/sessions/v1/stop-session`, {
    "session_type": type
   }, {
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}` 
    }
  });
  return res.data.session;
};
