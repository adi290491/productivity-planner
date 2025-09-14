import { useState } from 'react';
import { PieChart, Pie, Cell, ResponsiveContainer, Legend, Tooltip, LineChart, Line, XAxis, YAxis, CartesianGrid, BarChart, Bar } from 'recharts';
import type { DailyTrend, WeeklyTrend } from '../types/trend';
import { parseTimeToMinutes, formatMinutesToLabel } from '../utils/format';

interface DailyTrendItemProps {
  trend: DailyTrend;
}

interface WeeklyTrendItemProps {
  trend: WeeklyTrend;
}

const COLORS = {
  focus: '#4a5568',     // accent color
  meeting: '#f59e0b',   // amber
  break: '#10b981'      // emerald
};

export const DailyTrendItem = ({ trend }: DailyTrendItemProps) => {
  const [isExpanded, setIsExpanded] = useState(false);

  // Prepare data for pie chart
  const pieData = Object.entries(trend.breakdown).map(([type, time]) => ({
    name: type.charAt(0).toUpperCase() + type.slice(1),
    value: parseTimeToMinutes(time),
    color: COLORS[type as keyof typeof COLORS] || '#6b7280'
  }));

  return (
    <div className="bg-white border border-border rounded p-4 shadow-sm mb-2">
      <div 
        className="flex justify-between items-center cursor-pointer"
        onClick={() => setIsExpanded(!isExpanded)}
      >
        <h4 className="font-semibold text-text">
          {new Date(trend.date).toLocaleDateString()}
        </h4>
        <div className="flex items-center gap-2">
          <span className="text-sm text-accent">{trend.total_time}</span>
          <span className={`transform transition-transform ${isExpanded ? 'rotate-180' : ''}`}>
            ▼
          </span>
        </div>
      </div>
      
      {isExpanded && (
        <div className="mt-4 pt-4 border-t border-border">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Session Breakdown Card */}
            <div className="bg-card border border-border rounded p-4">
              <h5 className="text-sm font-semibold text-accent mb-4">Session Breakdown</h5>
              
              {/* Pie Chart */}
              <div className="h-48 mb-4">
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={pieData}
                      cx="50%"
                      cy="50%"
                      innerRadius={30}
                      outerRadius={70}
                      paddingAngle={2}
                      dataKey="value"
                    >
                      {pieData.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={entry.color} />
                      ))}
                    </Pie>
                    <Tooltip 
                      formatter={(value: number) => {
                        const hours = Math.floor(value / 60);
                        const mins = value % 60;
                        return `${hours}h ${mins}m`;
                      }}
                    />
                    <Legend />
                  </PieChart>
                </ResponsiveContainer>
              </div>
              
              {/* Numeric Summary */}
              <div className="space-y-2">
                {Object.entries(trend.breakdown).map(([type, time]) => (
                  <div key={type} className="flex justify-between items-center">
                    <div className="flex items-center gap-2">
                      <div 
                        className="w-3 h-3 rounded-full" 
                        style={{ backgroundColor: COLORS[type as keyof typeof COLORS] || '#6b7280' }}
                      />
                      <span className="text-text capitalize font-medium">{type}:</span>
                    </div>
                    <span className="text-accent font-mono font-semibold">{time}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export const WeeklyTrendItem = ({ trend }: WeeklyTrendItemProps) => {
  const [isExpanded, setIsExpanded] = useState(false);

  // Prepare data for focus line chart (only if daily_data exists)
  const focusLineData = trend.daily_data ? trend.daily_data.map(day => ({
    day: new Date(day.date).toLocaleDateString('en-US', { weekday: 'short' }),
    focus: parseTimeToMinutes(day.breakdown.focus || '0h 0m')
  })) : [];

  // Prepare data for stacked bar chart (only if daily_data exists)
  const stackedBarData = trend.daily_data ? trend.daily_data.map(day => ({
    day: new Date(day.date).toLocaleDateString('en-US', { weekday: 'short' }),
    focus: parseTimeToMinutes(day.breakdown.focus || '0h 0m'),
    meeting: parseTimeToMinutes(day.breakdown.meeting || '0h 0m'),
    break: parseTimeToMinutes(day.breakdown.break || '0h 0m')
  })) : [];

  return (
    <div className="bg-white border border-border rounded p-4 shadow-sm mb-2">
      <div 
        className="flex justify-between items-center cursor-pointer"
        onClick={() => setIsExpanded(!isExpanded)}
      >
        <h4 className="font-semibold text-text">
          Week of {new Date(trend.week_start).toLocaleDateString()}
        </h4>
        <div className="flex items-center gap-2">
          <span className="text-sm text-accent">{trend.total_time}</span>
          <span className={`transform transition-transform ${isExpanded ? 'rotate-180' : ''}`}>
            ▼
          </span>
        </div>
      </div>
      
      {isExpanded && (
        <div className="mt-4 pt-4 border-t border-border">
          {/* Only show daily charts if daily_data exists */}
          {trend.daily_data && trend.daily_data.length > 0 && (
            <div className="grid grid-cols-1 xl:grid-cols-2 gap-6 mb-6">
              
              {/* Focus Hours Line Chart */}
              <div className="bg-card border border-border rounded p-4">
                <h5 className="text-sm font-semibold text-accent mb-4">Focus Hours per Day</h5>
                <div className="h-48">
                  <ResponsiveContainer width="100%" height="100%">
                    <LineChart data={focusLineData}>
                      <CartesianGrid strokeDasharray="3 3" />
                      <XAxis dataKey="day" />
                      <YAxis tickFormatter={(value) => `${Math.floor(value / 60)}h`} />
                      <Tooltip 
                        formatter={(value: number) => [formatMinutesToLabel(value), 'Focus Time']}
                        labelFormatter={(label) => `${label}`}
                      />
                      <Line 
                        type="monotone" 
                        dataKey="focus" 
                        stroke={COLORS.focus} 
                        strokeWidth={3}
                        dot={{ fill: COLORS.focus, strokeWidth: 2, r: 4 }}
                        activeDot={{ r: 6 }}
                      />
                    </LineChart>
                  </ResponsiveContainer>
                </div>
              </div>

              {/* Stacked Bar Chart */}
              <div className="bg-card border border-border rounded p-4">
                <h5 className="text-sm font-semibold text-accent mb-4">Daily Session Balance</h5>
                <div className="h-48">
                  <ResponsiveContainer width="100%" height="100%">
                    <BarChart data={stackedBarData}>
                      <CartesianGrid strokeDasharray="3 3" />
                      <XAxis dataKey="day" />
                      <YAxis tickFormatter={(value) => `${Math.floor(value / 60)}h`} />
                      <Tooltip 
                        formatter={(value: number) => formatMinutesToLabel(value)}
                        labelFormatter={(label) => `${label}`}
                      />
                      <Legend />
                      <Bar dataKey="focus" stackId="a" fill={COLORS.focus} name="Focus" />
                      <Bar dataKey="meeting" stackId="a" fill={COLORS.meeting} name="Meeting" />
                      <Bar dataKey="break" stackId="a" fill={COLORS.break} name="Break" />
                    </BarChart>
                  </ResponsiveContainer>
                </div>
              </div>
            </div>
          )}
          
          {/* Show message if no daily data available */}
          {(!trend.daily_data || trend.daily_data.length === 0) && (
            <div className="bg-card border border-border rounded p-4 mb-6">
              <p className="text-accent text-center">Daily breakdown data not available for this week</p>
            </div>
          )}

          {/* Weekly Summary Card */}
          <div className="bg-card border border-border rounded p-4">
            <h5 className="text-sm font-semibold text-accent mb-4">Weekly Summary</h5>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="text-center p-3 bg-gray-50 rounded">
                <div className="text-2xl font-bold text-text">{trend.breakdown.focus || '0h 0m'}</div>
                <div className="text-sm text-accent">Total Focus Time</div>
              </div>
              <div className="text-center p-3 bg-gray-50 rounded">
                <div className="text-2xl font-bold text-text">{trend.avg_session_length || 'N/A'}</div>
                <div className="text-sm text-accent">Avg Focus Session Length</div>
              </div>
              <div className="text-center p-3 bg-gray-50 rounded">
                <div className="text-2xl font-bold text-text">{trend.longest_streak || 0} days</div>
                <div className="text-sm text-accent">Longest Streak</div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};