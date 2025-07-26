CREATE TABLE stock_overviews (
    symbol VARCHAR(10) PRIMARY KEY,
    name VARCHAR(128),
    pe_ratio NUMERIC(10, 4),
    peg_ratio NUMERIC(10, 4),
    price_to_book NUMERIC(10, 4),
    return_on_equity_ttm NUMERIC(10, 4),
    operating_margin NUMERIC(10, 4),
    profit_margin NUMERIC(10, 4),
    dividend_yield NUMERIC(10, 4),
    beta NUMERIC(10, 4),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stock_scorecards (
    symbol VARCHAR(10) PRIMARY KEY,
    company_name VARCHAR(128),
    pe_ratio NUMERIC(10, 4),
    peg_ratio NUMERIC(10, 4),
    price_to_book NUMERIC(10, 4),
    roe_ttm NUMERIC(10, 4),
    revenue_5y_growth NUMERIC(10, 4),
    operating_margin NUMERIC(10, 4),
    profit_margin NUMERIC(10, 4),
    dividend_yield NUMERIC(10, 4),
    beta NUMERIC(10, 4),
    historical_roe JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stock_income_statements (
    symbol VARCHAR(10),
    fiscal_date DATE,
    total_revenue NUMERIC(18, 2),
    net_income NUMERIC(18, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (symbol, fiscal_date)
);

CREATE TABLE stock_balance_sheets (
    symbol VARCHAR(10),
    fiscal_date DATE,
    shareholder_equity NUMERIC(18, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (symbol, fiscal_date)
);

