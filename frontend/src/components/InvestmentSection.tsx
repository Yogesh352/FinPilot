import React from "react";
import {
  PieChart,
  Pie,
  Cell,
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
} from "recharts";

const InvestmentSection: React.FC = () => {
  const allocationData = [
    { name: "Stocks", value: 70, color: "#3b82f6" },
    { name: "Bonds", value: 20, color: "#10b981" },
    { name: "Alternatives", value: 10, color: "#f59e0b" },
  ];

  const performanceData = [
    { category: "Mortgage", target: 33800, actual: 30000 },
    { category: "Automobile", target: 1540, actual: 1200 },
    { category: "Insurance", target: 432, actual: 500 },
    { category: "Groceries", target: 1241, actual: 1100 },
    { category: "All others", target: 3779, actual: 3200 },
  ];

  const formatCurrency = (value: number) => `$${value.toLocaleString()}`;

  return (
    <div className="space-y-6">
      {/* Investment Allocation */}
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">
          Investment Allocation
        </h3>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Pie Chart */}
          <div className="flex justify-center">
            <div className="w-48 h-48">
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={allocationData}
                    cx="50%"
                    cy="50%"
                    innerRadius={40}
                    outerRadius={80}
                    dataKey="value"
                  >
                    {allocationData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={entry.color} />
                    ))}
                  </Pie>
                </PieChart>
              </ResponsiveContainer>
            </div>
          </div>

          {/* Legend */}
          <div className="flex flex-col justify-center space-y-3">
            {allocationData.map((item, index) => (
              <div key={index} className="flex items-center justify-between">
                <div className="flex items-center">
                  <div
                    className="w-4 h-4 rounded mr-3"
                    style={{ backgroundColor: item.color }}
                  ></div>
                  <span className="text-sm font-medium text-gray-900">
                    {item.name}
                  </span>
                </div>
                <span className="text-sm font-semibold text-gray-900">
                  {item.value}%
                </span>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Performance Comparison */}
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">
          This Month vs Target
        </h3>

        <div className="h-64">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart
              data={performanceData}
              margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
            >
              <XAxis
                dataKey="category"
                axisLine={false}
                tickLine={false}
                tick={{ fontSize: 11, fill: "#6b7280" }}
                angle={-45}
                textAnchor="end"
                height={80}
              />
              <YAxis
                axisLine={false}
                tickLine={false}
                tick={{ fontSize: 11, fill: "#6b7280" }}
                tickFormatter={formatCurrency}
              />
              <Bar
                dataKey="target"
                fill="#e5e7eb"
                name="Target"
                radius={[4, 4, 0, 0]}
              />
              <Bar
                dataKey="actual"
                fill="#3b82f6"
                name="Actual"
                radius={[4, 4, 0, 0]}
              />
            </BarChart>
          </ResponsiveContainer>
        </div>

        <div className="mt-4 flex justify-center space-x-6">
          <div className="flex items-center">
            <div className="w-3 h-3 bg-gray-300 rounded mr-2"></div>
            <span className="text-sm text-gray-600">Target</span>
          </div>
          <div className="flex items-center">
            <div className="w-3 h-3 bg-chart-blue rounded mr-2"></div>
            <span className="text-sm text-gray-600">Actual</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default InvestmentSection;
