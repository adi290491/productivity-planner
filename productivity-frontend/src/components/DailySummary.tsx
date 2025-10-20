import type { DailySummary as DailySummaryType } from "../types/summary";

const DailySummary = ({ data }: { data: DailySummaryType | null }) => {
  if (!data) {
    return (
      <div className="summary-card-wrapper">
        <p>No summary data available.</p>
      </div>
    );
  }

  const noSessions = !data.breakdown || Object.keys(data.breakdown).length === 0;

  return (
    <div className="summary-card-wrapper">
      <div className="summary-header">
        <h3>Daily Summary</h3>
        <div className="summary-details">
          <p>Date: {data.date}</p>
          <p>Total Duration: {data.total_time}</p>
        </div>
      </div>

      <div className="summary-table-card">
        {noSessions ? (
          <p>Looks like you're having a relaxed day! No sessions recorded yet.</p>
        ) : (
          <table className="summary-table">
            <thead>
              <tr>
                <th>Session Type</th>
                <th>Duration</th>
              </tr>
            </thead>
            <tbody>
              {Object.entries(data.breakdown).map(([type, time]) => (
                <tr key={type}>
                  <td>{type.charAt(0).toUpperCase() + type.slice(1)}</td>
                  <td>{time}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
};


export default DailySummary;
