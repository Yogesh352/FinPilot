import React from "react";

const BudgetingSection: React.FC = () => {
  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
      <h3 className="text-lg font-semibold text-gray-900 mb-4">Budgeting</h3>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Cash Flow */}
        <div>
          <div className="flex items-center justify-between mb-3">
            <h4 className="font-medium text-gray-900">Cash flow</h4>
            <span className="text-green-600 font-semibold">+$1,679</span>
          </div>
          <div className="text-sm text-gray-500 mb-4">This month</div>

          <div className="space-y-3">
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-600">Income this month</span>
              <span className="text-sm font-medium">$8,672</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-600">Expenses this month</span>
              <span className="text-sm font-medium">$6,993</span>
            </div>
          </div>
        </div>

        {/* Monthly Progress */}
        <div>
          <div className="mb-3">
            <div className="flex justify-between items-center mb-2">
              <span className="text-sm font-medium text-gray-900">
                Monthly Progress
              </span>
              <span className="text-sm text-gray-500">Last month $3,598</span>
            </div>

            <div className="w-full bg-gray-200 rounded-full h-3">
              <div
                className="bg-success-green h-3 rounded-full"
                style={{ width: "70%" }}
              ></div>
            </div>
          </div>

          <div className="mt-4">
            <div className="text-center">
              <div className="text-2xl font-bold text-gray-900 mb-1">
                $6,812
              </div>
              <div className="text-sm text-gray-500">Cash flow this year</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default BudgetingSection;
