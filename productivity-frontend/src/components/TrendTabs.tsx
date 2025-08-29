import { useState } from 'react';
import type { DailyTrendResponse, WeeklyTrendResponse } from '../types/trend';
import { DailyTrendItem, WeeklyTrendItem } from './TrendChart';

// Dummy data for testing
const dummyDailyTrends: DailyTrendResponse = {
  user_id: "user123",
  dailyTrends: [
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
    }
  ]
};

const dummyWeeklyTrends: WeeklyTrendResponse = {
  user_id: "user123",
  weekly_trends: [
    {
      week_start: "2025-08-25",
      total_time: "32h 15m",
      breakdown: {
        focus: "24h 30m",
        meeting: "6h 45m",
        break: "1h 0m"
      },
      avg_session_length: "45m",
      longest_streak: 4,
      daily_data: [
        { date: "2025-08-25", total_time: "6h 30m", breakdown: { focus: "4h 0m", meeting: "2h 0m", break: "30m" } },
        { date: "2025-08-26", total_time: "5h 45m", breakdown: { focus: "3h 30m", meeting: "1h 45m", break: "30m" } },
        { date: "2025-08-27", total_time: "6h 30m", breakdown: { focus: "3h 20m", meeting: "2h 10m", break: "1h 0m" } },
        { date: "2025-08-28", total_time: "7h 15m", breakdown: { focus: "5h 20m", meeting: "1h 30m", break: "25m" } },
        { date: "2025-08-29", total_time: "4h 30m", breakdown: { focus: "3h 0m", meeting: "1h 0m", break: "30m" } },
        { date: "2025-08-30", total_time: "1h 45m", breakdown: { focus: "1h 20m", meeting: "0h 15m", break: "10m" } },
        { date: "2025-08-31", total_time: "0h 0m", breakdown: { focus: "0h 0m", meeting: "0h 0m", break: "0h 0m" } }
      ]
    },
    {
      week_start: "2025-08-18",
      total_time: "28h 45m",
      breakdown: {
        focus: "20h 15m",
        meeting: "7h 30m",
        break: "1h 0m"
      },
      avg_session_length: "40m",
      longest_streak: 3,
      daily_data: [
        { date: "2025-08-18", total_time: "5h 30m", breakdown: { focus: "3h 30m", meeting: "1h 30m", break: "30m" } },
        { date: "2025-08-19", total_time: "6h 15m", breakdown: { focus: "4h 0m", meeting: "2h 0m", break: "15m" } },
        { date: "2025-08-20", total_time: "4h 45m", breakdown: { focus: "3h 15m", meeting: "1h 15m", break: "15m" } },
        { date: "2025-08-21", total_time: "6h 0m", breakdown: { focus: "4h 30m", meeting: "1h 15m", break: "15m" } },
        { date: "2025-08-22", total_time: "3h 30m", breakdown: { focus: "2h 30m", meeting: "45m", break: "15m" } },
        { date: "2025-08-23", total_time: "2h 45m", breakdown: { focus: "2h 0m", meeting: "30m", break: "15m" } },
        { date: "2025-08-24", total_time: "0h 0m", breakdown: { focus: "0h 0m", meeting: "0h 0m", break: "0h 0m" } }
      ]
    }
  ]
};

const TrendTabs = () => {
  const [activeTab, setActiveTab] = useState<'daily' | 'weekly'>('daily');

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
          {activeTab === 'daily' ? (
            <div className="space-y-2">
              <p className="text-sm text-accent mb-4">
                Your daily productivity patterns over the last few days
              </p>
              {dummyDailyTrends.dailyTrends.map((trend, index) => (
                <DailyTrendItem key={`${trend.date}-${index}`} trend={trend} />
              ))}
            </div>
          ) : (
            <div className="space-y-2">
              <p className="text-sm text-accent mb-4">
                Your weekly productivity trends and patterns
              </p>
              {dummyWeeklyTrends.weekly_trends.map((trend, index) => (
                <WeeklyTrendItem key={`${trend.week_start}-${index}`} trend={trend} />
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default TrendTabs;
