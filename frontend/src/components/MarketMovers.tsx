import React from "react";
import { TrendingUp, TrendingDown } from "lucide-react";

interface MarketMover {
  symbol: string;
  name: string;
  change: number;
  changePercent: number;
  isPositive: boolean;
}

const MarketMovers: React.FC = () => {
  const topMovers: MarketMover[] = [
    {
      symbol: "AAPL",
      name: "Apple Inc.",
      change: 2.5,
      changePercent: 1.2,
      isPositive: true,
    },
    {
      symbol: "TSLA",
      name: "Tesla Inc.",
      change: -8.25,
      changePercent: -2.8,
      isPositive: false,
    },
    {
      symbol: "MSFT",
      name: "Microsoft",
      change: 3.75,
      changePercent: 0.9,
      isPositive: true,
    },
    {
      symbol: "AMZN",
      name: "Amazon",
      change: -1.2,
      changePercent: -0.4,
      isPositive: false,
    },
  ];

  const recommendations = [
    {
      title: "See recommendations",
      description:
        "Based on your portfolio, we've identified 3 potential opportunities for optimization.",
    },
  ];

  const formatChange = (change: number) => {
    const sign = change >= 0 ? "+" : "";
    return `${sign}$${change.toFixed(2)}`;
  };

  const formatPercent = (percent: number) => {
    const sign = percent >= 0 ? "+" : "";
    return `${sign}${percent.toFixed(1)}%`;
  };

  return (
    <div className="space-y-6">
      {/* Market Movers */}
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">
          Market movers
        </h3>

        <div className="space-y-4">
          {topMovers.map((stock, index) => (
            <div key={index} className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <div
                  className={`p-1 rounded ${
                    stock.isPositive ? "bg-green-100" : "bg-red-100"
                  }`}
                >
                  {stock.isPositive ? (
                    <TrendingUp className="h-4 w-4 text-green-600" />
                  ) : (
                    <TrendingDown className="h-4 w-4 text-red-600" />
                  )}
                </div>
                <div>
                  <div className="font-medium text-gray-900">
                    {stock.symbol}
                  </div>
                  <div className="text-xs text-gray-500">{stock.name}</div>
                </div>
              </div>
              <div className="text-right">
                <div
                  className={`font-medium ${
                    stock.isPositive ? "text-green-600" : "text-red-600"
                  }`}
                >
                  {formatChange(stock.change)}
                </div>
                <div
                  className={`text-xs ${
                    stock.isPositive ? "text-green-600" : "text-red-600"
                  }`}
                >
                  {formatPercent(stock.changePercent)}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Investment Insights */}
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">
          Investment Insights
        </h3>

        <div className="space-y-4">
          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <h4 className="font-medium text-blue-900 mb-2">
              Portfolio Analysis
            </h4>
            <p className="text-sm text-blue-700 mb-3">
              Your portfolio is well-diversified but could benefit from
              rebalancing.
            </p>
            <button className="text-sm font-medium text-blue-600 hover:text-blue-800 transition-colors">
              View detailed analysis →
            </button>
          </div>

          <div className="bg-green-50 border border-green-200 rounded-lg p-4">
            <h4 className="font-medium text-green-900 mb-2">
              Tax Optimization
            </h4>
            <p className="text-sm text-green-700 mb-3">
              You could save $2,400 in taxes with strategic rebalancing.
            </p>
            <button className="text-sm font-medium text-green-600 hover:text-green-800 transition-colors">
              Learn more →
            </button>
          </div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">
          Quick Actions
        </h3>

        <div className="space-y-3">
          <button className="w-full text-left p-3 bg-empower-blue text-white rounded-lg hover:bg-blue-800 transition-colors">
            <div className="font-medium">Rebalance Portfolio</div>
            <div className="text-sm text-blue-200">
              Optimize your allocation
            </div>
          </button>

          <button className="w-full text-left p-3 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
            <div className="font-medium text-gray-900">
              Schedule Consultation
            </div>
            <div className="text-sm text-gray-500">Talk to an advisor</div>
          </button>
        </div>
      </div>
    </div>
  );
};

export default MarketMovers;
