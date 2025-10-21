import type { DailySummary as DailySummaryType } from "../types/summary";

const formatDate = (dateString: string) => new Date(dateString).toLocaleDateString('en-CA');
const formatDuration = (time: string | number) => typeof time === 'number' ? `${time}s` : time;

const DailySummary = ({ data }: { data: DailySummaryType | null }) => {
  const noSessions = !data || !data.breakdown || Object.keys(data.breakdown).length === 0;

  const displayDate = data ? data.date : new Date().toISOString();

  return (
    <div className="summary-card-wrapper">
      <div className="summary-header">
        <h3>Daily Summary</h3>
        <div className="summary-details">
   
          <p>Date: {formatDate(displayDate)}</p>
          
  
          {data && (
            <p>Total Duration: {formatDuration(data.total_time)}</p>
          )}
        </div>
      </div>

      <div className="summary-table-card">
        {noSessions ? (
          <div className="flex h-full items-center justify-center text-center text-gray-500 p-4">
   
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
