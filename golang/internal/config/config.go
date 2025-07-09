package config

import (
    "database/sql"
    "fmt"
    "os"
)

type Config struct {
    PostgresURL string
}

func Load() Config {
    return Config{
        PostgresURL: os.Getenv("POSTGRES_URL"),
    }
}

func SetupPostgres(cfg Config) (*sql.DB, error) {
    db, err := sql.Open("postgres", cfg.PostgresURL)
    if err != nil {
        return nil, err
    }
    return db, db.Ping()
}
