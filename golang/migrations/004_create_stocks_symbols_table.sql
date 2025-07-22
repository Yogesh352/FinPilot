CREATE TABLE IF NOT EXISTS stock_symbols (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) UNIQUE,
    batch_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX IF NOT EXISTS idx_stocks_symbols_batch_id ON stock_symbols(symbol);
