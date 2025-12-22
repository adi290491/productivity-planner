import { useState, useEffect, useRef } from "react";
import CarouselHeader from "../components/CarouselHeader";
import type { DailySummary as DailySummaryType, WeeklySummaryResponse, SessionResponse } from "../types/summary";
import DailySummary from "../components/DailySummary";
import { fetchDailySummary, fetchWeeklySummary } from "../api/summary";
import { fetchLatestTrendsCount } from "../api/trend-analysis";
import WeeklySummary from "../components/WeeklySummary";
import SessionControl from "../components/SessionControl";
import TrendTabs from "../components/TrendTabs";
import NotificationBanner from "../components/NotificationBanner";
import type { UnviewedTrendsCount } from "../types/trend";
import type { TrendTabsHandle } from "../components/TrendTabs";


const Dashboard = () => {
    const [sessionType, setSessionType] = useState("focus");
    const [lastSession, setLastSession] = useState<SessionResponse | null>(null);
    const [dailySummary, setDailySummary] = useState<DailySummaryType | null>(null);
    const [weeklySummary, setWeeklySummary] = useState<WeeklySummaryResponse | null>(null);
    const [trendsCount, setTrendsCount] = useState<UnviewedTrendsCount | null> (null);
    const trendTabsRef = useRef<HTMLDivElement>(null);
    const trendTabsComponentRef = useRef<TrendTabsHandle>(null);
    const token = localStorage.getItem("token");
    
    useEffect(() => {
        const fetchSummaries = async () => {
          try{
            const [daily, weekly, trendsCount] = await Promise.all([
              fetchDailySummary(token!),
              fetchWeeklySummary(token!),
              fetchLatestTrendsCount(token!)
            ]);
            console.log("daily summary:", daily);
            console.log("weekly summary:", weekly);
            console.log("trends count:", trendsCount);

            setDailySummary(daily);
            setWeeklySummary(weekly);
            if (trendsCount && (trendsCount.daily_count > 0 || trendsCount.weekly_count > 0)) {
              setTrendsCount(trendsCount);
            }
          } catch (error) {
            console.error("Failed to fetch summaries:", error);
          }
        }

        if (token) {
          fetchSummaries();
        }
    }, [token]);

    const scrollToTrends = async () => {
        console.log('Scrolling to trends and fetching data...');

        if (trendTabsComponentRef.current){
          await trendTabsComponentRef.current.fetchTrends();
        }

        trendTabsRef.current?.scrollIntoView({behavior: 'smooth', block: 'start'});
    };

    return (
       
        <div className="bg-background min-h-screen">

            <CarouselHeader />

            {trendsCount && (trendsCount.daily_count > 0 || trendsCount.weekly_count > 0) && (
            <NotificationBanner 
                  onDismiss={() => setTrendsCount(null)} 
                  trendsCount={trendsCount}
                  onViewTrends={scrollToTrends}
                />
            )}

            <SessionControl
                sessionType={sessionType}
                setSessionType={setSessionType}
                setDailySummary={() => fetchDailySummary(token!).then(setDailySummary)}
                setLastSession={setLastSession}
                lastSession={lastSession}
            />

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <DailySummary data={dailySummary} />
                <WeeklySummary data={weeklySummary} />
            </div>
            
            <div ref={trendTabsRef}>
              <TrendTabs ref={trendTabsComponentRef} shouldFetchOnMount={false} />
            </div>
        </div>
    );
};

export default Dashboard;