import React from "react";
import NetWorthSection from "./NetWorthSection";
import AssetsSidebar from "./AssetsSidebar";
import BudgetingSection from "./BudgetingSection";
import InvestmentSection from "./InvestmentSection";
import PortfolioSection from "./PortfolioSection";
import MarketMovers from "./MarketMovers";

const Dashboard: React.FC = () => {
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="grid grid-cols-12 gap-6">
        {/* Left Sidebar - Assets & Liabilities */}
        <div className="col-span-12 lg:col-span-3">
          <AssetsSidebar />
        </div>

        {/* Main Content */}
        <div className="col-span-12 lg:col-span-6 space-y-6">
          <NetWorthSection />
          <BudgetingSection />
          <InvestmentSection />
          <PortfolioSection />
        </div>

        {/* Right Sidebar */}
        <div className="col-span-12 lg:col-span-3">
          <MarketMovers />
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
