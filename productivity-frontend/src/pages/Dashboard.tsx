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


const Dashboard = () => {
    const [sessionType, setSessionType] = useState("focus");
    const [lastSession, setLastSession] = useState<SessionResponse | null>(null);
    const [dailySummary, setDailySummary] = useState<DailySummaryType | null>(null);
    const [weeklySummary, setWeeklySummary] = useState<WeeklySummaryResponse | null>(null);
    const [isBannerVisible, setIsBannerVisible] = useState(false);
    const [trendsCount, setTrendsCount] = useState<UnviewedTrendsCount | null> (null);
    const trendTabsRef = useRef<HTMLDivElement>(null);

    const token = localStorage.getItem("token");
    
    useEffect(() => {
        const fetchData = async () => {

          if (!token) return

            try {
              const [daily, weekly, trendsCount] = await Promise.all([
                fetchDailySummary(token),
                fetchWeeklySummary(token),
                fetchLatestTrendsCount(token),
              ]);
              setDailySummary(daily);
              setWeeklySummary(weekly);

              if (trendsCount && (trendsCount.weekly_count > 0 || trendsCount.weekly_count > 0 )) {
                setIsBannerVisible(true);
                setTrendsCount(trendsCount);
              }
            } catch (error) {
              console.error("Failed to fetch summaries:", error);
              setDailySummary(null);
              setWeeklySummary(null);
            }
          };
          fetchData();
    }, [token]);

    const scrollToTrends = () => {
        trendTabsRef.current?.scrollIntoView({behavior: 'smooth', block: 'start'});
    };

    return (
       
        <div className="bg-background min-h-screen">

            <div className={`notification-container ${isBannerVisible ? 'visible' : ''}`}>
                <NotificationBanner 
                  onDismiss={() => setIsBannerVisible(false)} 
                  trendsCount={trendsCount}
                  onViewTrends={scrollToTrends}
                />
            </div>

            <CarouselHeader />

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
              <TrendTabs />
            </div>
        </div>
    );
};

export default Dashboard;