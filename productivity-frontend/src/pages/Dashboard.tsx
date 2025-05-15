import { useState, useEffect } from "react";
import CarouselHeader from "../components/CarouselHeader";
import type { DailySummary as DailySummaryType, WeeklySummaryResponse, SessionResponse } from "../types/summary";
import DailySummary from "../components/DailySummary";
import { fetchDailySummary, fetchWeeklySummary } from "../api/summary";
import WeeklySummary from "../components/WeeklySummary";
import SessionControl from "../components/SessionControl";


const Dashboard = () => {
    const [sessionType, setSessionType] = useState("focus");
    const [lastSession, setLastSession] = useState<SessionResponse | null>(null);
    const [dailySummary, setDailySummary] = useState<DailySummaryType | null>(null);
    const [weeklySummary, setWeeklySummary] = useState<WeeklySummaryResponse | null>(null);

    const token = localStorage.getItem("token");
    
    useEffect(() => {
        const fetchData = async () => {
            try {
              const [daily, weekly] = await Promise.all([
                fetchDailySummary(token!),
                fetchWeeklySummary(token!)
              ]);
              setDailySummary(daily);
              setWeeklySummary(weekly);
            } catch (error) {
              console.error("Failed to fetch summaries:", error);
              setDailySummary(null);
              setWeeklySummary(null);
            }
          };
          fetchData();
    }, []);

    return (
       
        <div className="bg-background min-h-screen">
            <CarouselHeader />

            <SessionControl
                sessionType={sessionType}
                setSessionType={setSessionType}
                setDailySummary={() => fetchDailySummary(token!).then(setDailySummary)}
                setLastSession={setLastSession}
                lastSession={lastSession}
            />

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {dailySummary && <DailySummary data={dailySummary} />}
            {weeklySummary && <WeeklySummary data={weeklySummary} />}
            </div>
        </div>
    );
};

export default Dashboard;