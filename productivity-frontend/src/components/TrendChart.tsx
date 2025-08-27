import React from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

interface TrendChartProps {
  data: any[]; // Replace 'any' with a more specific type based on your trend data structure
}

const TrendChart: React.FC<TrendChartProps> = ({ data }) => {
  return (
    <ResponsiveContainer width="100%" height={300}>
      <BarChart
        data={data}
        margin={{
          top: 5,
          right: 30,
          left: 20,
          bottom: 5,
        }}
      >
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="name" /> {/* Replace 'name' with the key for your x-axis data */}
        <YAxis />
        <Tooltip />
        <Legend />
        <Bar dataKey="value" fill="#8884d8" /> {/* Replace 'value' with the key for your bar values */}
      </BarChart>
    </ResponsiveContainer>
  );
};

export default TrendChart;