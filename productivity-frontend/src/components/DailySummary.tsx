import type { DailySummary as DailySummaryType } from "../types/summary";

const formatDuration = (time: string | number) => typeof time === 'number' ? `${time}s` : time;

const formatDateSafe = (dateString: string): string => {
  if (!dateString || typeof dateString !== 'string') {
    console.error('Invalid date string:', dateString);
    return 'Today'
  }

  try{
    const [year, month, day] = dateString.split('-').map(Number);
    if (!year || !month || !day) {
      throw new Error('Invalid date format');
    }

    const date = new Date(year, month - 1, day);
    if (isNaN(date.getTime())) {
      throw new Error('Invalid date')
    }

    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  }catch(error){
    console.error('Error formatting date:', dateString, error);
    return 'Today';
  }
}

const DailySummary = ({ data }: { data: DailySummaryType | null }) => {
  const noSessions = !data || !data.breakdown || Object.keys(data.breakdown).length === 0;

  const getCurrentDate = (): string => {
    const now = new Date();
    const year = now.getFullYear();
    const month = String(now.getMonth() + 1).padStart(2, '0');
    const day = String(now.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`
  }

  const displayDate = data?.date || getCurrentDate();

  const formattedDate = formatDateSafe(displayDate)

  return (
    <div className="summary-card-wrapper">
      <div className="summary-header">
        <h3>Daily Summary</h3>
        <div className="summary-details">
   
          <p>Date: {formattedDate}</p>
          
  
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
