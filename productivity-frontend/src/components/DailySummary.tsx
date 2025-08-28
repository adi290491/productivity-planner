import type { DailySummary as DailySummaryType } from "../types/summary";

const DailySummary = ({ data }: { data: DailySummaryType | null }) => {
  const noSessions = !data || !data.breakdown || Object.keys(data.breakdown).length === 0;

  return (
  <div className="bg-card border border-border rounded p-4 shadow">
    <h3 className="text-lg font-semibold text-text mb-2">Daily Summary</h3>
    <div className="text-sm text-accent space-y-2">
      {noSessions ? (
        <p>Looks like you're having a relaxed day! No sessions recorded yet.</p>
      ) : (
        <>
      <p><strong>Date:</strong> {new Date(data.date).toLocaleDateString()}</p>
      <p><strong>Total Time:</strong> {data.total_time}</p>
      <div>
        <p className="font-semibold">Breakdown:</p>
        <ul className="ml-4 list-disc">
          {Object.entries(data.breakdown).map(([type, time]) => (
 <li key={type}>
              {type.charAt(0).toUpperCase() + type.slice(1)}: {time}
            </li>
          ))}
        </ul>
      </div>
        </>
      )}
    </div>
  </div>
);};

export default DailySummary;
