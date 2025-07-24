import React from "react";
import { ChevronRight } from "lucide-react";

interface AssetItem {
  name: string;
  value: number;
  subItems?: { name: string; value: number; description?: string }[];
}

const AssetsSidebar: React.FC = () => {
  const assets: AssetItem[] = [
    {
      name: "Cash",
      value: 342983.31,
      subItems: [
        {
          name: "Bank of America",
          value: 50345.12,
          description: "Savings Checking & Savings",
        },
        {
          name: "Bank of America",
          value: 615417.69,
          description: "Savings - Other",
        },
        { name: "Empower Personal Cash", value: 173465.18 },
      ],
    },
    {
      name: "Investments",
      value: 2494469.29,
      subItems: [
        { name: "Fidelity", value: 171487.56 },
        { name: "Charles Schwab", value: 342983.31 },
        { name: "Empower Retirement I", value: 343569.14 },
      ],
    },
  ];

  const liabilities: AssetItem[] = [
    {
      name: "Credit card",
      value: 14289.97,
      subItems: [
        { name: "American Express", value: 8881.12 },
        { name: "Chase Unlimited", value: 5408.85 },
      ],
    },
  ];

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  const renderAssetSection = (
    title: string,
    items: AssetItem[],
    totalValue: number
  ) => (
    <div className="mb-6">
      <div className="flex justify-between items-center mb-3">
        <h3 className="font-semibold text-gray-900">{title}</h3>
        <span className="text-lg font-semibold text-gray-900">
          {formatCurrency(totalValue)}
        </span>
      </div>

      {items.map((item, index) => (
        <div key={index} className="mb-4">
          <div className="flex justify-between items-center mb-2 cursor-pointer hover:bg-gray-50 p-2 rounded-lg transition-colors">
            <div className="flex items-center">
              <ChevronRight className="h-4 w-4 text-gray-400 mr-2" />
              <span className="font-medium text-gray-900">{item.name}</span>
            </div>
            <span className="font-medium text-gray-900">
              {formatCurrency(item.value)}
            </span>
          </div>

          {item.subItems && (
            <div className="ml-6 space-y-1">
              {item.subItems.map((subItem, subIndex) => (
                <div
                  key={subIndex}
                  className="flex justify-between items-center text-sm"
                >
                  <div>
                    <div className="text-gray-700">{subItem.name}</div>
                    {subItem.description && (
                      <div className="text-gray-500 text-xs">
                        {subItem.description}
                      </div>
                    )}
                  </div>
                  <span className="text-gray-700">
                    {formatCurrency(subItem.value)}
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>
      ))}
    </div>
  );

  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
      {renderAssetSection("Assets", assets, 2837452.6)}
      <hr className="my-6 border-gray-200" />
      {renderAssetSection("Liabilities", liabilities, 314237.94)}
    </div>
  );
};

export default AssetsSidebar;
