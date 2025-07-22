CREATE TABLE IF NOT EXISTS stocks_intraday_indicators (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL REFERENCES stock_symbols(symbol) ON DELETE CASCADE,
    date TIMESTAMP NOT NULL,
    rvol FLOAT,
    price_change_pct FLOAT,
    rsi FLOAT,
    volume_spike BOOLEAN,
    atr FLOAT,
    obv BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(symbol, date)
);

CREATE INDEX IF NOT EXISTS idx_stocks_intraday_indicators_symbol ON stocks_intraday_indicators(symbol);
CREATE INDEX IF NOT EXISTS idx_stocks_intraday_indicators_date ON stocks_intraday_indicators(date);
CREATE INDEX IF NOT EXISTS idx_stocks_intraday_indicators_symbol_date ON stocks_intraday_indicators(symbol, date);

CREATE TRIGGER update_stocks_intraday_metrics_updated_at 
    BEFORE UPDATE ON stocks_intraday_indicators 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
