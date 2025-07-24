import React from "react";
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  ResponsiveContainer,
  Tooltip,
} from "recharts";

const PortfolioSection: React.FC = () => {
  const portfolioData = [
    { month: "Jan", value: 120000 },
    { month: "Feb", value: 125000 },
    { month: "Mar", value: 135000 },
    { month: "Apr", value: 145000 },
    { month: "May", value: 155000 },
    { month: "Jun", value: 165000 },
    { month: "Jul", value: 170000 },
    { month: "Aug", value: 175918 },
  ];

  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-lg font-semibold text-gray-900">
          Portfolio balances
        </h3>
        <div className="text-right">
          <div className="text-2xl font-bold text-gray-900">$175,918</div>
          <div className="text-sm text-gray-500">Total balance</div>
        </div>
      </div>

      <div className="h-48 mb-4">
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={portfolioData}>
            <defs>
              <linearGradient
                id="portfolioGradient"
                x1="0"
                y1="0"
                x2="0"
                y2="1"
              >
                <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.3} />
                <stop offset="95%" stopColor="#3b82f6" stopOpacity={0} />
              </linearGradient>
            </defs>
            <XAxis
              dataKey="month"
              axisLine={false}
              tickLine={false}
              tick={{ fontSize: 11, fill: "#6b7280" }}
            />
            <YAxis hide />
            <Tooltip
              formatter={(value) => [
                `$${Number(value).toLocaleString()}`,
                "Portfolio Value",
              ]}
              labelStyle={{ color: "#374151" }}
              contentStyle={{
                backgroundColor: "#f9fafb",
                border: "1px solid #e5e7eb",
                borderRadius: "6px",
              }}
            />
            <Area
              type="monotone"
              dataKey="value"
              stroke="#3b82f6"
              strokeWidth={2}
              fill="url(#portfolioGradient)"
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>

      <div className="grid grid-cols-2 gap-4 text-sm">
        <div>
          <div className="text-gray-500 mb-1">1 month return</div>
          <div className="font-semibold text-success-green">+2.5%</div>
        </div>
        <div>
          <div className="text-gray-500 mb-1">YTD return</div>
          <div className="font-semibold text-success-green">+12.8%</div>
        </div>
      </div>
    </div>
  );
};

export default PortfolioSection;
