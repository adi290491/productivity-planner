import { useRef, useState } from "react";
import type { SessionResponse } from "../types/summary";
import { startSession, stopSession } from "../api/session";
import { formatDuration } from "../utils/format";


type Props = {
    sessionType: string;
    setSessionType: (type: string) => void;
    setDailySummary: () => void;
    setLastSession: (session: SessionResponse) => void;
    lastSession: SessionResponse | null;
};

const SessionControl = ({sessionType, lastSession, setSessionType, setDailySummary, setLastSession}: Props) => {
    const [sessionStarted, setSessionStarted] = useState(false);
    const [elapsed, setElapsed] = useState(0);
    const [loading, setLoading] = useState(false);

    const startTimeRef = useRef<number | null>(null);
    const intervalRef = useRef<NodeJS.Timeout | null>(null);

    const token = localStorage.getItem("token")

    const handleStart = async () => {
        try {
          setLoading(true);
          await startSession(sessionType, token!);

          startTimeRef.current = Date.now();
          setElapsed(0);
          setSessionStarted(true);
   
          intervalRef.current = setInterval(() => {
            if (startTimeRef.current !== null) {
              const now = Date.now();
              const elapsedSeconds = Math.floor((now - startTimeRef.current) / 1000);
              setElapsed(elapsedSeconds)
            }
          }, 100);
        } catch {
          alert("Failed to start session");
        } finally {
          setLoading(false);
        }
      };

    const handleStop = async () => {
        try {
          setLoading(true);
          const session = await stopSession(sessionType, token!);
          if (intervalRef.current) {
            clearInterval(intervalRef.current);
            intervalRef.current = null;
          }

          startTimeRef.current = null;
          setElapsed(0);
          setSessionStarted(false);

          setLastSession(session);
          await setDailySummary();
        } catch {
          alert("Failed to stop session");
        } finally {
          setLoading(false);
        }
      };

    const formatTime = (totalSeconds: number) => {
        const hours = Math.floor(totalSeconds / 3600);
        const minutes = Math.floor(totalSeconds / 60);
        const seconds = totalSeconds % 60;
        return `${hours}H ${minutes}M ${seconds}S`;
      };

      return (
        <div className="bg-card border border-border rounded p-6 shadow mb-10">
          <h2 className="text-xl font-semibold text-text mb-4">Session Controls</h2>
    
          <div className="flex flex-col md:flex-row gap-4 items-start md:items-center">
            <label className="text-accent">Session Type:</label>
            <select
              value={sessionType}
              disabled={sessionStarted}
              onChange={(e) => setSessionType(e.target.value)}
              className="border border-border rounded px-4 py-2 bg-white text-text"
            >
              <option value="focus">Focus</option>
              <option value="meeting">Meeting</option>
              <option value="break">Break</option>
            </select>
    
            <button
              onClick={handleStart}
              disabled={loading || sessionStarted}
              className={`btn btn-accent ${sessionStarted || loading ? "btn-disabled" : ""}`}
            >
              {loading && sessionStarted ? "..." : "Start Session"}
            </button>
    
            <button
              onClick={handleStop}
              disabled={!sessionStarted || loading}
              className={`btn btn-gray ${!sessionStarted || loading ? "btn-disabled" : ""}`}
            >
              {loading && !sessionStarted ? "..." : "Stop Session"}
            </button>
    
            <div className="bg-white border border-border rounded px-6 py-4 shadow text-center">
              <p className="text-accent text-sm mb-2">Current Session Duration</p>
              <p className="text-2xl font-mono text-text tracking-wide">{formatTime(elapsed)}</p>
            </div>
          </div>

          {sessionStarted && (
            <div className="bg-yellow-50 border border-yellow-200 rounded p-3 mt-4 text-sm text-yellow-800">
              <p>Session will automatically stop at midnight and be recorded for today.</p>
            </div>
          )}

          {lastSession && (
            <div className="bg-white text-text border border-border rounded p-4 shadow mt-6 text-sm">
              <h3 className="text-lg font-semibold text-text mb-2">Last Session</h3>
              <p><strong>Type:</strong> {lastSession.type}</p>
                <p><strong>Start:</strong> {new Date(lastSession.start_time).toLocaleString()}</p>
                <p><strong>End:</strong> {new Date(lastSession.end_time).toLocaleString()}</p>
                <p><strong>Duration:</strong> {formatDuration(lastSession.start_time, lastSession.end_time)}</p>
            </div>
            )}
        </div>
      );
};

export default SessionControl;