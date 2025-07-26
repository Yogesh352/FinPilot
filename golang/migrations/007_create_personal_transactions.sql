CREATE TABLE personal_transactions (
    id SERIAL PRIMARY KEY,
    date DATE,
    amount NUMERIC(10, 2),
    currency VARCHAR(10),
    description TEXT,
    category VARCHAR(128),
    bank VARCHAR(128),
    account VARCHAR(128),
    inserted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_personal_transactions ON stocks_intraday(symbol);
