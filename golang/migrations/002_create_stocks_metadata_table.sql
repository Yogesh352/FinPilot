-- Create stocks_metadata table for storing stock metadata
CREATE TABLE IF NOT EXISTS stocks_metadata (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL UNIQUE REFERENCES stock_symbols(symbol) ON DELETE CASCADE,
    company_name VARCHAR(255),
    industry VARCHAR(100),
    exchange VARCHAR(50),
    currency VARCHAR(10) DEFAULT 'USD',
    market_cap DECIMAL(20,2),
    description TEXT,
    website VARCHAR(255),
    type VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_stocks_metadata_symbol ON stocks_metadata(symbol);
CREATE INDEX IF NOT EXISTS idx_stocks_metadata_industry ON stocks_metadata(industry);
CREATE INDEX IF NOT EXISTS idx_stocks_metadata_exchange ON stocks_metadata(exchange);


-- Create a function to update the updated_at timestamp for metadata
CREATE OR REPLACE FUNCTION update_metadata_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at for metadata
CREATE TRIGGER update_stocks_metadata_updated_at 
    BEFORE UPDATE ON stocks_metadata 
    FOR EACH ROW 
    EXECUTE FUNCTION update_metadata_updated_at_column();
