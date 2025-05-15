import type { WeeklySummaryResponse } from "../types/summary";
import { parseTimeToMinutes, formatMinutesToLabel } from "../utils/format";
import {
  ResponsiveContainer,
  BarChart,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  Bar
} from "recharts";

const WeeklySummary = ({ data }: { data: WeeklySummaryResponse }) => {
  const chartData = data.daily_summaries.map(day => {
    const result: any = { date: new Date(day.date).toLocaleDateString() };
    for (const [type, value] of Object.entries(day.breakdown)) {
      result[type] = parseTimeToMinutes(value);
    }
    return result;
  });

  return (
    <div className="bg-card border border-border rounded p-4 shadow">
      <h3 className="text-lg font-semibold text-text mb-2">Weekly Summary</h3>
      <ResponsiveContainer width="100%" height={300}>
        <BarChart data={chartData} margin={{ top: 20, right: 30, left: 20, bottom: 5 }}>
          <XAxis dataKey="date" />
          <YAxis tickFormatter={formatMinutesToLabel} />
          <Tooltip formatter={(value: any) => formatMinutesToLabel(Number(value))} />
          <Legend />
          <Bar dataKey="focus" stackId="a" fill="#007bff" />
          <Bar dataKey="meeting" stackId="a" fill="#ffc107" />
          <Bar dataKey="break" stackId="a" fill="#28a745" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
};

export default WeeklySummary;
