import React from "react";
import { Bell, ChevronDown } from "lucide-react";

const Header: React.FC = () => {
  const navItems = ["Overview", "Account", "Budgeting", "Investing", "Advice"];

  return (
    <header className="bg-empower-blue text-white shadow-lg">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <div className="flex items-center">
            <div className="text-xl font-bold tracking-wide">EMPOWER</div>
          </div>

          {/* Navigation */}
          <nav className="hidden md:flex space-x-8">
            {navItems.map((item) => (
              <a
                key={item}
                href="#"
                className="text-white hover:text-blue-200 px-3 py-2 text-sm font-medium transition-colors"
              >
                {item}
              </a>
            ))}
          </nav>

          {/* Right side - notifications and profile */}
          <div className="flex items-center space-x-4">
            {/* Notification bell */}
            <button className="relative p-2 text-white hover:text-blue-200 transition-colors">
              <Bell className="h-5 w-5" />
              <span className="absolute top-1 right-1 h-2 w-2 bg-red-500 rounded-full"></span>
            </button>

            {/* Profile dropdown */}
            <div className="flex items-center space-x-2 cursor-pointer hover:bg-blue-800 px-3 py-2 rounded-lg transition-colors">
              <span className="text-sm font-medium">Charles</span>
              <ChevronDown className="h-4 w-4" />
            </div>

            {/* Sign out */}
            <button className="text-sm text-white hover:text-blue-200 transition-colors">
              Sign out
            </button>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
