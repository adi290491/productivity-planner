import React from 'react';
import TrendItem from './TrendItem';

interface TrendData {
  date: string;
  // Add other properties based on your trend data structure
  data: any; // Placeholder for the actual trend data for the chart
}

interface TrendListProps {
  trends: TrendData[];
}

const TrendList: React.FC<TrendListProps> = ({ trends }) => {
  // Sort trends by date in descending order
  const sortedTrends = [...trends].sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime());

  return (
    <ul>
      {sortedTrends.map((trend, index) => (
        <TrendItem key={index} trend={trend} />
      ))}
    </ul>
  );
};

export default TrendList;