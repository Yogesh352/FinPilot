-- Create stocks_raw table for storing financial data
CREATE TABLE IF NOT EXISTS stocks_raw (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL REFERENCES stock_symbols(symbol) ON DELETE CASCADE,
    date TIMESTAMP NOT NULL,
    open DECIMAL(10,4) NOT NULL,
    high DECIMAL(10,4) NOT NULL,
    low DECIMAL(10,4) NOT NULL,
    close DECIMAL(10,4) NOT NULL,
    volume BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure unique combination of symbol and date
    UNIQUE(symbol, date)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_stocks_raw_symbol ON stocks_raw(symbol);
CREATE INDEX IF NOT EXISTS idx_stocks_raw_date ON stocks_raw(date);
CREATE INDEX IF NOT EXISTS idx_stocks_raw_symbol_date ON stocks_raw(symbol, date);

-- Create a function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_stocks_raw_updated_at 
    BEFORE UPDATE ON stocks_raw 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE stocks_raw IS 'Raw stock data from external APIs';
COMMENT ON COLUMN stocks_raw.symbol IS 'Stock symbol (e.g., AAPL, GOOGL)';
COMMENT ON COLUMN stocks_raw.date IS 'Trading date';
COMMENT ON COLUMN stocks_raw.open IS 'Opening price';
COMMENT ON COLUMN stocks_raw.high IS 'Highest price during the day';
COMMENT ON COLUMN stocks_raw.low IS 'Lowest price during the day';
COMMENT ON COLUMN stocks_raw.close IS 'Closing price';
COMMENT ON COLUMN stocks_raw.volume IS 'Trading volume'; 