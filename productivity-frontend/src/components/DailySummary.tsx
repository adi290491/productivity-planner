import type { DailySummary as DailySummaryType } from "../types/summary";

// Assuming you have these helper functions
const formatDate = (dateString: string) => new Date(dateString).toLocaleDateString('en-CA');
const formatDuration = (time: string | number) => typeof time === 'number' ? `${time}s` : time;

const DailySummary = ({ data }: { data: DailySummaryType | null }) => {
  const noSessions = !data || !data.breakdown || Object.keys(data.breakdown).length === 0;

  // 1. Determine which date to display. Use the date from data if available, otherwise use today's date.
  const displayDate = data ? data.date : new Date().toISOString();

  return (
    <div className="summary-card-wrapper">
      <div className="summary-header">
        <h3>Daily Summary</h3>
        <div className="summary-details">
          {/* 2. The date is now always visible */}
          <p>Date: {formatDate(displayDate)}</p>
          
          {/* The total duration is only shown when there is data */}
          {data && (
            <p>Total Duration: {formatDuration(data.total_time)}</p>
          )}
        </div>
      </div>

      <div className="summary-table-card">
        {noSessions ? (
          <div className="flex h-full items-center justify-center text-center text-gray-500 p-4">
            {/* 3. Updated message to be more specific */}
            <p>No sessions recorded for this date yet.</p>
          </div>
        ) : (
          <div className="table-scroll-container">
            <table className="summary-table">
              <thead>
                <tr>
                  <th>Session Type</th>
                  <th>Duration</th>
                </tr>
              </thead>
              <tbody>
                {Object.entries(data!.breakdown).map(([type, time]) => (
                  <tr key={type}>
                    <td>{type.charAt(0).toUpperCase() + type.slice(1)}</td>
                    <td>{time}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
};

export default DailySummary;
