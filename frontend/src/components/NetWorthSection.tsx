import React from "react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  ResponsiveContainer,
  Tooltip,
} from "recharts";

const NetWorthSection: React.FC = () => {
  const netWorthData = [
    { month: "Jan", value: 750000 },
    { month: "Feb", value: 780000 },
    { month: "Mar", value: 820000 },
    { month: "Apr", value: 850000 },
    { month: "May", value: 880000 },
    { month: "Jun", value: 890000 },
    { month: "Jul", value: 910000 },
    { month: "Aug", value: 936754 },
  ];

  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
      <div className="mb-6">
        <div className="flex items-baseline justify-between">
          <div>
            <h2 className="text-sm font-medium text-gray-500 mb-1">
              NET WORTH
            </h2>
            <div className="text-3xl font-bold text-gray-900">$936,754.60</div>
          </div>
          <div className="text-right">
            <div className="text-sm text-gray-500">
              Increase your relative return by 3%
            </div>
            <div className="text-xs text-gray-400 mt-1">
              Your Target Personal Strategy is forecasted to generate a 3%
              higher relative return compared to a target date fund. Stay on
              track and make even more money on your retirement and other goals.
            </div>
          </div>
        </div>
      </div>

      <div className="h-64 w-full">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart data={netWorthData}>
            <XAxis
              dataKey="month"
              axisLine={false}
              tickLine={false}
              tick={{ fontSize: 12, fill: "#6b7280" }}
            />
            <YAxis hide />
            <Tooltip
              formatter={(value) => [
                `$${Number(value).toLocaleString()}`,
                "Net Worth",
              ]}
              labelStyle={{ color: "#374151" }}
              contentStyle={{
                backgroundColor: "#f9fafb",
                border: "1px solid #e5e7eb",
                borderRadius: "6px",
              }}
            />
            <Line
              type="monotone"
              dataKey="value"
              stroke="#3b82f6"
              strokeWidth={3}
              dot={false}
              activeDot={{ r: 6, fill: "#3b82f6" }}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>

      <div className="mt-4 flex items-center justify-between text-sm text-gray-500">
        <span>Jan '24 - Aug '24</span>
        <div className="flex items-center space-x-4">
          <span className="flex items-center">
            <div className="w-3 h-3 bg-chart-blue rounded mr-2"></div>
            Target allocation
          </span>
          <span className="flex items-center">
            <div className="w-3 h-3 bg-success-green rounded mr-2"></div>
            Current allocation
          </span>
        </div>
      </div>
    </div>
  );
};

export default NetWorthSection;
