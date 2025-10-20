import {
  ResponsiveContainer,
  BarChart,
  XAxis,
  YAxis,
  Tooltip,
  Bar,
  CartesianGrid,
} from "recharts";
import type { WeeklySummaryResponse } from "../types/summary";
import { parseTimeToMinutes } from "../utils/format"; // Assuming this utility exists

// Helper to format the date range string for the header
const formatDateRange = (start: string, end: string): string => {
  const startDate = new Date(start);
  const endDate = new Date(end);
  const options: Intl.DateTimeFormatOptions = { month: 'short', day: 'numeric', year: 'numeric' };
  return `${startDate.toLocaleDateString('en-US', options)} to ${endDate.toLocaleDateString('en-US', options)}`;
};

// Custom Legend Component to match the design
const CustomLegend = () => (
  <div className="chart-legend">
    <div className="legend-item">
      <div className="legend-color-dot" style={{ backgroundColor: 'var(--color-primary)' }} />
      <span>Focus Sessions</span>
    </div>
    <div className="legend-item">
      <div className="legend-color-dot" style={{ backgroundColor: 'var(--color-accent)' }} />
      <span>Meeting Sessions</span>
    </div>
    <div className="legend-item">
      <div className="legend-color-dot" style={{ backgroundColor: 'var(--color-break)' }} />
      <span>Break Sessions</span>
    </div>
  </div>
);

// Custom Tooltip Formatter
const formatTooltipValue = (value: number) => {
  const totalMinutes = Math.round(value * 60);
  const hours = Math.floor(totalMinutes / 60);
  const minutes = totalMinutes % 60;
  return `${hours}h ${minutes}m`;
};

const WeeklySummary = ({ data }: { data: WeeklySummaryResponse | null }) => {
  if (!data || !data.daily_summaries || data.daily_summaries.length === 0) {
    return (
      <div className="weekly-summary-wrapper">
        <div className="weekly-summary-header">
          <h3>Weekly Summary</h3>
          <p>No data available for this week. Start a session to see your progress!</p>
        </div>
      </div>
    );
  }

  // Transform the string-based time from the API into numerical hours for the chart
  const chartData = data.daily_summaries.map((day) => {
    const breakdownInHours: { [key: string]: number } = {};
    for (const [type, value] of Object.entries(day.breakdown)) {
      // Convert "0h14m" into minutes (14), then into hours (0.233)
      breakdownInHours[type] = parseTimeToMinutes(value as string) / 60;
    }
    return {
      // Format date to "Mon", "Tue", etc. for the X-axis
      name: new Date(day.date).toLocaleDateString('en-US', { weekday: 'short' }),
      ...breakdownInHours,
    };
  });

  return (
    <div className="weekly-summary-wrapper">
      <div className="weekly-summary-header">
        <h3>Weekly Summary</h3>
        <p>{formatDateRange(data.start_date, data.end_date)}</p>
      </div>

      <div className="chart-container">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart
            data={chartData}
            margin={{ top: 20, right: 10, left: 10, bottom: 5 }}
          >
            <CartesianGrid strokeDasharray="3 3" vertical={false} />
            <XAxis dataKey="name" tickLine={false} axisLine={false} />
            <YAxis
              label={{ value: 'Total Time (hours)', angle: -90, position: 'insideLeft', offset: 0 }}
              tickLine={false}
              axisLine={false}
            />
            <Tooltip formatter={formatTooltipValue} cursor={{fill: 'rgba(246, 173, 85, 0.1)'}} />
            <Bar dataKey="focus" stackId="a" fill="var(--color-primary)" radius={[4, 4, 0, 0]} />
            <Bar dataKey="meeting" stackId="a" fill="var(--color-accent)" />
            <Bar dataKey="break" stackId="a" fill="var(--color-break)" radius={[4, 4, 0, 0]} />
          </BarChart>
        </ResponsiveContainer>
      </div>
      
      <CustomLegend />
    </div>
  );
};

export default WeeklySummary;
