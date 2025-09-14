import { useState, useEffect } from 'react';
import type { DailyTrend, DailyTrendResponse, WeeklyTrend, WeeklyTrendResponse } from '../types/trend';
import { DailyTrendChart, WeeklyTrendChart } from './TrendChart';
import { fetchDailyTrends, fetchWeeklyTrends } from '../api/trend-analysis';

// Dummy data for testing

/*const dummyDailyTrends: DailyTrend[] = [
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
];*/

/*const dummyWeeklyTrends: WeeklyTrend[] = [
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
*/


const TrendTabs = () => {
  const [activeTab, setActiveTab] = useState<'daily' | 'weekly'>('daily');
  const [dailyTrends, setDailyTrends ] = useState<DailyTrend[]>([]);
  const [weeklyTrends, setWeeklyTrends ] = useState<WeeklyTrend[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const getTrends = async () => {
      try {
        setLoading(true);
        const token = localStorage.getItem('token');
        if (!token) {
          throw new Error("No token found");
        }

        const dailyData: DailyTrendResponse = await fetchDailyTrends(token, 7);
        console.log("Daily Trends Data:", dailyData);
        setDailyTrends(dailyData.dailyTrends || []);
        // setDailyTrends(dummyDailyTrends)

        const weeklyData: WeeklyTrendResponse = await fetchWeeklyTrends(token, 4);
        console.log("Weekly Trends Data:", weeklyData);
        setWeeklyTrends(weeklyData.weeklyTrends || []);
        // setWeeklyTrends(dummyWeeklyTrends);

        setError(null);
      } catch (err) {
        if (err instanceof Error) {
            setError(err.message);
        } else {
            setError("An unknown error occurred");
        }
      } finally {
        setLoading(false);
      }
    };

    getTrends();
  }, []);

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
          {loading && <p>Loading...</p>}
          {error && <p className="text-red-500">{error}</p>}
          {!loading && !error && (
            <>
              {activeTab === 'daily' ? (
                <DailyTrendChart data={dailyTrends} />
              ) : (
                <WeeklyTrendChart data={weeklyTrends} />
              )}
            </>
          )}
        </div>
      </div>
    </div>
  );
};

export default TrendTabs;
