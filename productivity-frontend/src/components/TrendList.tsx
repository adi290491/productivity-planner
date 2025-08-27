import React, { useState } from 'react';
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
  return (
    <ul>
      {trends.map((trend, index) => (
        <TrendItem key={index} trend={trend} />
      ))}
    </ul>
  );
};

export default TrendList;