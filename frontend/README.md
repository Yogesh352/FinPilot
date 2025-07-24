# Financial Dashboard

A modern React TypeScript application that replicates the Empower (formerly Personal Capital) financial dashboard interface.

## Features

- **Net Worth Tracking**: Interactive chart showing net worth growth over time
- **Assets & Liabilities**: Detailed breakdown of financial accounts and positions
- **Budgeting Tools**: Cash flow analysis and monthly progress tracking
- **Investment Analytics**: Portfolio allocation charts and performance metrics
- **Market Data**: Real-time market movers and investment insights
- **Responsive Design**: Modern UI that works on desktop and mobile devices

## Technologies Used

- React 18 with TypeScript
- Tailwind CSS for styling
- Recharts for data visualization
- Lucide React for icons
- Responsive grid layout

## Getting Started

1. Install dependencies:

   ```bash
   npm install
   ```

2. Start the development server:

   ```bash
   npm start
   ```

3. Open [http://localhost:3000](http://localhost:3000) to view the dashboard in your browser.

## Components

- `Header`: Navigation bar with Empower branding
- `Dashboard`: Main layout container
- `NetWorthSection`: Large chart and net worth display
- `AssetsSidebar`: Assets and liabilities breakdown
- `BudgetingSection`: Cash flow and monthly progress
- `InvestmentSection`: Allocation charts and comparisons
- `PortfolioSection`: Portfolio balances with area chart
- `MarketMovers`: Market data and investment insights

## Customization

The dashboard uses realistic sample data for demonstration. In a production environment, you would integrate with financial APIs to fetch real account data and market information.

## Build

To create a production build:

```bash
npm run build
```
