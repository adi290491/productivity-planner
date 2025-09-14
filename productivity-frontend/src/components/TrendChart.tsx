import { BarChart, Bar, LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import type { DailyTrend, WeeklyTrend } from '../types/trend';
import { parseTimeToMinutes, formatMinutesToLabel } from '../utils/format';

interface DailyTrendChartProps {
  data: DailyTrend[];
}

interface WeeklyTrendChartProps {
    data: WeeklyTrend[];
}

const COLORS = {
  focus: '#4a5568',
  meeting: '#f59e0b',
  break: '#10b981'
};

export const DailyTrendChart = ({ data }: DailyTrendChartProps) => {
  const chartData = data.map(item => ({
    name: new Date(item.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
    focus: parseTimeToMinutes(item.breakdown.focus),
    meeting: parseTimeToMinutes(item.breakdown.meeting),
    break: parseTimeToMinutes(item.breakdown.break),
  }));

  return (
    <ResponsiveContainer width="100%" height={400}>
      <BarChart data={chartData} margin={{ top: 5, right: 20, left: -10, bottom: 5 }}>
        <CartesianGrid strokeDasharray="3 3" vertical={false} />
        <XAxis dataKey="name" />
        <YAxis tickFormatter={(value) => `${value}h`} />
        <Tooltip
          formatter={(value: number) => [formatMinutesToLabel(value), 'Time']}
          labelFormatter={(label) => `${label}`}
        />
        <Legend />
        <Bar dataKey="focus" fill={COLORS.focus} name="Focus" />
        <Bar dataKey="meeting" fill={COLORS.meeting} name="Meeting" />
        <Bar dataKey="break" fill={COLORS.break} name="Break" />
      </BarChart>
    </ResponsiveContainer>
  );
};

export const WeeklyTrendChart = ({ data }: WeeklyTrendChartProps) => {
    const chartData = data.map(item => ({
        name: new Date(item.week_start).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
        focus: parseTimeToMinutes(item.breakdown.focus),
        meeting: parseTimeToMinutes(item.breakdown.meeting),
        break: parseTimeToMinutes(item.breakdown.break),
    }));

    return (
        <ResponsiveContainer width="100%" height={400}>
            <LineChart data={chartData} margin={{ top: 5, right: 20, left: -10, bottom: 5 }}>
                <CartesianGrid strokeDasharray="3 3" vertical={false} />
                <XAxis dataKey="name" />
                <YAxis tickFormatter={(value) => `${value}h`} />
                <Tooltip
                    formatter={(value: number) => [formatMinutesToLabel(value), 'Time']}
                    labelFormatter={(label) => `Week of ${label}`}
                />
                <Legend />
                <Line type="monotone" dataKey="focus" stroke={COLORS.focus} strokeWidth={2} name="Focus" />
                <Line type="monotone" dataKey="meeting" stroke={COLORS.meeting} strokeWidth={2} name="Meeting" />
                <Line type="monotone" dataKey="break" stroke={COLORS.break} strokeWidth={2} name="Break" />
            </LineChart>
        </ResponsiveContainer>
    );
};
