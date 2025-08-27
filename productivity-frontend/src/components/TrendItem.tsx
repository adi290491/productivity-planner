import React, { useState } from 'react';
import TrendChart from './TrendChart';

interface TrendItemProps {
  trend: {
    date: string; // Assuming the date is a string
    data: any; // Assuming the trend data for the chart is in a 'data' property
  };
}

const TrendItem: React.FC<TrendItemProps> = ({ trend }) => {
  const [isExpanded, setIsExpanded] = useState(false);

  const toggleExpand = () => {
    setIsExpanded(!isExpanded);
  };

  return (
    <div style={{ border: '1px solid #ccc', margin: '10px', padding: '10px' }}>
      <div onClick={toggleExpand} style={{ cursor: 'pointer' }}>
        <h3>{trend.date}</h3>
      </div>
      {isExpanded && (
        <div>
          <TrendChart data={trend.data} />
        </div>
      )}
    </div>
  );
};

export default TrendItem;