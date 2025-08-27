import React, { useState, useEffect } from 'react';
import { fetchDailyTrends, fetchWeeklyTrends } from '../api/trend-analysis';
import TrendList from './TrendList';
import type { DailyTrend, WeeklyTrend } from '../types/trend';

const TrendTabs: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'daily' | 'weekly'>('daily');
  const [dailyTrends, setDailyTrends] = useState<DailyTrend[] | null>(null);
  const [weeklyTrends, setWeeklyTrends] = useState<WeeklyTrend[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      setError(null);
      try {
        if (activeTab === 'daily') {
          const data = await fetchDailyTrends();
          setDailyTrends(data);
        } else {
          const data = await fetchWeeklyTrends();
          setWeeklyTrends(data);
        }
      } catch (err) {
        setError('Failed to fetch trend data.');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [activeTab]);

  return (
    <div className="mt-8"> {/* Add some vertical space */}
      <div className="flex space-x-4 border-b border-gray-700"> {/* Flex container for tabs */}
        <button
          className={`py-2 px-4 text-lg font-medium focus:outline-none ${
            activeTab === 'daily'
              ? 'border-b-2 border-primary text-primary'
              : 'text-gray-400 hover:text-gray-200'
          }`}
          onClick={() => setActiveTab('daily')}
        >
          Daily Trends
        </button>
        <button
          className={`py-2 px-4 text-lg font-medium focus:outline-none ${
            activeTab === 'weekly'
              ? 'border-b-2 border-primary text-primary'
              : 'text-gray-400 hover:text-gray-200'
          }`}
          onClick={() => setActiveTab('weekly')}
        >
          Weekly Trends
        </button>
      </div>
      <div className="tab-content">
        {/* Content will go here */}
        {loading && <p className="text-gray-400 mt-4">Loading trends...</p>}
        {error && <p className="text-red-500 mt-4">Error: {error}</p>}

        {!loading && !error && activeTab === 'daily' && (
          <div>
            {dailyTrends && dailyTrends.length > 0 ? (
              <TrendList trends={dailyTrends} />
            ) : (
              <p className="text-gray-400 mt-4">No daily trend data available.</p>
            )}
          </div>
        )}
        {!loading && !error && activeTab === 'weekly' && (
          <div>
            {weeklyTrends && weeklyTrends.length > 0 ? (
              <TrendList trends={weeklyTrends} />
            ) : (
              <p className="text-gray-400 mt-4">No weekly trend data available.</p>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default TrendTabs;