import { useState, useEffect, forwardRef, useImperativeHandle } from 'react';
import type { DailyTrend, DailyTrendResponse, WeeklyTrend, WeeklyTrendResponse } from '../types/trend';
import { DailyTrendChart, WeeklyTrendChart } from './TrendChart';
import { fetchDailyTrends, fetchWeeklyTrends } from '../api/trend-analysis';

// Dummy data for testing

const dummyDailyTrends: DailyTrend[] = [
    {
      date: "2025-08-27",
      total_time: "6h 30m",
      breakdown: {
        focus: "3h 20m",
        meeting: "2h 10m",
        break: "1h 0m"
      }
    },
    {
      date: "2025-08-26",
      total_time: "5h 45m",
      breakdown: {
        focus: "3h 30m",
        meeting: "1h 45m",
        break: "30m"
      }
    },
    {
      date: "2025-08-25",
      total_time: "7h 15m",
      breakdown: {
        focus: "4h 45m",
        meeting: "1h 50m",
        break: "40m"
      }
    },
    {
        date: "2025-08-24",
        total_time: "8h 0m",
        breakdown: {
          focus: "5h 0m",
          meeting: "2h 30m",
          break: "30m"
        }
    },
    {
        date: "2025-08-23",
        total_time: "4h 30m",
        breakdown: {
          focus: "3h 0m",
          meeting: "1h 0m",
          break: "30m"
        }
    }
];

const dummyWeeklyTrends: WeeklyTrend[] = [
    {
      week_start: "2025-08-25",
      total_time: "32h 15m",
      breakdown: {
        focus: "24h 30m",
        meeting: "6h 45m",
        break: "1h 0m"
      },
    },
    {
      week_start: "2025-08-18",
      total_time: "28h 45m",
      breakdown: {
        focus: "20h 15m",
        meeting: "7h 30m",
        break: "1h 0m"
      },
    },
    {
        week_start: "2025-08-11",
        total_time: "35h 0m",
        breakdown: {
          focus: "28h 0m",
          meeting: "5h 0m",
          break: "2h 0m"
        },
    },
]; 

export interface TrendTabsHandle {
  fetchTrends: () => Promise<void>;
}

interface TrendTabsProps {
  shouldFetchOnMount?: boolean
}

const TrendTabs = forwardRef<TrendTabsHandle, TrendTabsProps>(({ shouldFetchOnMount = false }, ref) => {
  const [activeTab, setActiveTab] = useState<'daily' | 'weekly'>('daily');
  const [dailyTrends, setDailyTrends ] = useState<DailyTrend[]>([]);
  const [weeklyTrends, setWeeklyTrends ] = useState<WeeklyTrend[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [hasFetched, setHasFetched] = useState(false);

    const getTrends = async () => {
      if (hasFetched) return;

      try {
        setLoading(true);
        setError(null);
        const token = localStorage.getItem('token');
        if (!token) {
          throw new Error("No token found");
        }

        console.log("Fetch trends data...")

        try {
          const dailyData: DailyTrendResponse = await fetchDailyTrends(token, 7);
          console.log("Daily Trends Data:", dailyData);
          setDailyTrends(dailyData.dailyTrends || dummyDailyTrends);
        } catch (dailyErr: any) {
          console.log("Daily trends error - treating as no data available");

          if (dailyErr.response?.status === 500) {
            const errorData = dailyErr.response?.data;
            const errorMessage = errorData?.Message || errorData?.message || errorData?.error || '';
            
            if (errorMessage.toLowerCase().includes('no') && 
                (errorMessage.toLowerCase().includes('trend') || errorMessage.toLowerCase().includes('found'))) {
              setDailyTrends([]);
            } else {
             
              console.warn("Unexpected 500 error:", errorMessage);
              setDailyTrends([]);
            }
          } else {
            throw dailyErr; 
          }
        }

        // Fetch weekly trends with error handling
        try {
          const weeklyData: WeeklyTrendResponse = await fetchWeeklyTrends(token, 4);
          console.log("Weekly Trends Data:", weeklyData);
          setWeeklyTrends(weeklyData.weeklyTrends || dummyWeeklyTrends);
        } catch (weeklyErr: any) {
          console.log("Weekly trends error - treating as no data available");
          
          if (weeklyErr.response?.status === 500) {
            const errorData = weeklyErr.response?.data;
            const errorMessage = errorData?.Message || errorData?.message || errorData?.error || '';
            
            if (errorMessage.toLowerCase().includes('no') && 
                (errorMessage.toLowerCase().includes('trend') || errorMessage.toLowerCase().includes('found'))) {
              setWeeklyTrends([]);
            } else {

              console.warn("Unexpected 500 error:", errorMessage);
              setWeeklyTrends([]);
            }
          } else {
            throw weeklyErr; 
          }
        }
        setHasFetched(true);
        console.log('Trends fetch completed');

      } catch (err: any) {
        console.error("Error fetching trends:", err);
        if (err instanceof Error) {
          setError(err.message);
        } else {
          setError("Failed to load trends. Please try again later.");
        }
      } finally {
        setLoading(false);
      }
    };

  useImperativeHandle(ref, () => ({
    fetchTrends: getTrends
  }));

  useEffect(() => {
    if (shouldFetchOnMount) {
      getTrends();
    }
  }, [shouldFetchOnMount]);

  return (
    <div className="trend-analysis-container">
      <div className="bg-card border border-border rounded p-6 shadow">
        <h2 className="text-xl font-semibold text-text mb-6">Trend Analysis</h2>
        
        <div className="flex mb-6">
          <button
            onClick={() => setActiveTab('daily')}
            className={`trend-tab-button ${activeTab === 'daily' ? 'active' : ''}`}
          >
            Daily Trends
          </button>
          <button
            onClick={() => setActiveTab('weekly')}
            className={`trend-tab-button ${activeTab === 'weekly' ? 'active' : ''}`}
          >
            Weekly Trends
          </button>
        </div>

        <div className="min-h-[400px]">
          {loading && (
            <div className="flex items-center justify-center h-[400px]">
              <p className="text-text-secondary">Loading trends...</p>
            </div>
          )}
          
          {error && (
            <div className="flex items-center justify-center h-[400px]">
              <p className="text-red-500">{error}</p>
            </div>
          )}
          
          {!loading && !error && (
            <>
              {activeTab === 'daily' ? (
                dailyTrends.length > 0 ? (
                  <DailyTrendChart data={dailyTrends} />
                ) : (
                  <div className="flex items-center justify-center h-[400px]">
                    <div className="text-center">
                      <p className="text-text-secondary text-lg mb-2">No daily trends available yet</p>
                      <p className="text-text-secondary text-sm">Complete some sessions to see your daily productivity trends</p>
                    </div>
                  </div>
                )
              ) : (
                weeklyTrends.length > 0 ? (
                  <WeeklyTrendChart data={weeklyTrends} />
                ) : (
                  <div className="flex items-center justify-center h-[400px]">
                    <div className="text-center">
                      <p className="text-text-secondary text-lg mb-2">No weekly trends available yet</p>
                      <p className="text-text-secondary text-sm">Complete sessions over multiple weeks to see your weekly trends</p>
                    </div>
                  </div>
                )
              )}
            </>
          )}
        </div>
      </div>
    </div>
  );
});

TrendTabs.displayName = 'TrendTabs';

export default TrendTabs;
