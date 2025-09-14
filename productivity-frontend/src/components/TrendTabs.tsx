import { useState, useEffect } from 'react';
import type { DailyTrend, DailyTrendResponse, WeeklyTrend, WeeklyTrendResponse } from '../types/trend';
import { DailyTrendItem, WeeklyTrendItem } from './TrendChart';
import { fetchDailyTrends, fetchWeeklyTrends } from '../api/trend-analysis';

const TrendTabs = () => {
  const [activeTab, setActiveTab] = useState<'daily' | 'weekly'>('daily');
  const [dailyTrends, setDailyTrends] = useState<DailyTrend[]>([]);
  const [weeklyTrends, setWeeklyTrends] = useState<WeeklyTrend[]>([]);
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
        setDailyTrends(dailyData.daily_trends || []);

        const weeklyData: WeeklyTrendResponse = await fetchWeeklyTrends(token, 4);
        setWeeklyTrends(weeklyData.weekly_trends || []);

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
        
        {/* Tab Buttons */}
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

        {/* Tab Content */}
        <div className="min-h-[300px]">
          {loading && <p>Loading...</p>}
          {error && <p className="text-red-500">{error}</p>}
          {!loading && !error && (
            <>
              {activeTab === 'daily' ? (
                <div className="space-y-2">
                  <p className="text-sm text-accent mb-4">
                    Your daily productivity patterns over the last few days
                  </p>
                  {dailyTrends.map((trend, index) => (
                    <DailyTrendItem key={`${trend.date}-${index}`} trend={trend} />
                  ))}
                </div>
              ) : (
                <div className="space-y-2">
                  <p className="text-sm text-accent mb-4">
                    Your weekly productivity trends and patterns
                  </p>
                  {weeklyTrends.map((trend, index) => (
                    <WeeklyTrendItem key={`${trend.week_start}-${index}`} trend={trend} />
                  ))}
                </div>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  );
};

export default TrendTabs;
