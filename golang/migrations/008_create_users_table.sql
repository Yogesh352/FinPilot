CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(256),
    password TEXT
);

CREATE INDEX IF NOT EXISTS idx_users ON users(idx);
