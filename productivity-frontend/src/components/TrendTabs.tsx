import React, { useState, useEffect } from 'react';
import { fetchDailyTrends, fetchWeeklyTrends } from '../api/trend-analysis';
import TrendList from './TrendList';

const TrendTabs: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'daily' | 'weekly'>('daily');
  const [dailyTrends, setDailyTrends] = useState([]);
  const [weeklyTrends, setWeeklyTrends] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  return (
    <div>
      <div className="tabs">
        <button
          className={`tab ${activeTab === 'daily' ? 'active' : ''}`}
          onClick={() => setActiveTab('daily')}
        >
          Daily Trends
        </button>
        <button
          className={`tab ${activeTab === 'weekly' ? 'active' : ''}`}
          onClick={() => setActiveTab('weekly')}
        >
          Weekly Trends
        </button>
      </div>
      <div className="tab-content">
        {/* Content will go here */}
        {activeTab === 'daily' && <div>Daily Trends Content</div>}
        {activeTab === 'weekly' && <div>Weekly Trends Content</div>}
      </div>
    </div>
  );
};

export default TrendTabs;