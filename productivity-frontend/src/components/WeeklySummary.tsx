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
const WeeklySummary = ({ data }: { data: WeeklySummaryResponse | null }) => {
  return (
    <div className="bg-card border border-border rounded p-4 shadow">
      <h3 className="text-lg font-semibold text-text mb-2">Weekly Summary</h3>
      {data && data.daily_summaries && data.daily_summaries.length > 0 ? (
        <>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={data.daily_summaries.map((day) => {
              const result: Record<string, any> = { date: new Date(day.date).toLocaleDateString() };
              for (const [type, value] of Object.entries(day.breakdown)) {
                result[type] = parseTimeToMinutes(value as string);
              }
              return result;
            })} margin={{ top: 20, right: 30, left: 20, bottom: 5 }}> 
              <XAxis dataKey="date" /> 
              <YAxis tickFormatter={formatMinutesToLabel} />
              <Tooltip formatter={(value: any) => formatMinutesToLabel(Number(value))} />
              <Legend />
              <Bar dataKey="focus" stackId="a" fill="#007bff" />
              <Bar dataKey="meeting" stackId="a" fill="#ffc107" />
              <Bar dataKey="break" stackId="a" fill="#28a745" />
            </BarChart> 
          </ResponsiveContainer> 
          {/* Weekly Summary Card */}
          <div className="mt-4 bg-gray-100 p-4 rounded">
            <h4 className="text-md font-semibold text-text mb-2">Summary Metrics</h4>
            <p><strong>Total Focus Time:</strong> {data.total_time}</p>
            {/* Assuming average session length and longest streak are also available in data */}
            {/* You'll need to add these properties to the WeeklySummaryResponse type */}
            {/* <p><strong>Avg Focus Session Length:</strong> {data.average_session_length}</p> */}
            {/* <p><strong>Longest Streak:</strong> {data.longest_streak} days</p> */}
          </div>
        </>
      ) : ( 
        <div>Your week is a blank canvas! Start a session to paint your productivity picture.</div> 
      )}
    </div>
  );
};

export default WeeklySummary;
