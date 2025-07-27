package config

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "time"
    "github.com/joho/godotenv"
)

type Config struct {
    PostgresURL string
    // External API configurations
    AlphaVantageAPIKey string
    FinnhubAPIKey      string
    YahooFinanceAPIKey string
    PolygonAPIKey      string
    RowsAPIKey          string
    // API rate limiting
    APIRequestTimeout time.Duration
    MaxRequestsPerMinute int
    // Database settings
    MaxDBConnections int
    DBTimeout        time.Duration
    //JWT
    JWTSecret string
}

func Load() Config {
    // Load .env file if it exists
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: .env file not found or could not be loaded: %v", err)
        log.Printf("Make sure to set environment variables manually or create a .env file")
    }

    return Config{
        PostgresURL:         getEnvOrDefault("POSTGRES_URL", "postgres://localhost:5432/finance_db?sslmode=disable"),
        AlphaVantageAPIKey:  os.Getenv("ALPHA_VANTAGE_API_KEY"),
        FinnhubAPIKey:       os.Getenv("FINNHUB_API_KEY"),
        YahooFinanceAPIKey:  os.Getenv("YAHOO_FINANCE_API_KEY"),
        PolygonAPIKey:       os.Getenv("POLYGON_API_KEY"),
        RowsAPIKey:       os.Getenv("ROWS_API_KEY"),
        APIRequestTimeout:   getDurationEnvOrDefault("API_REQUEST_TIMEOUT", 30*time.Second),
        MaxRequestsPerMinute: getIntEnvOrDefault("MAX_REQUESTS_PER_MINUTE", 60),
        MaxDBConnections:    getIntEnvOrDefault("MAX_DB_CONNECTIONS", 10),
        DBTimeout:           getDurationEnvOrDefault("DB_TIMEOUT", 5*time.Second),
        JWTSecret: os.Getenv("JWT_SECRET"),
    }
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getIntEnvOrDefault(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if parsed, err := fmt.Sscanf(value, "%d", &defaultValue); err == nil && parsed == 1 {
            return defaultValue
        }
    }
    return defaultValue
}

func getDurationEnvOrDefault(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if parsed, err := time.ParseDuration(value); err == nil {
            return parsed
        }
    }
    return defaultValue
}

func SetupPostgres(cfg Config) (*sql.DB, error) {
    db, err := sql.Open("postgres", cfg.PostgresURL)
    if err != nil {
        return nil, err
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(cfg.MaxDBConnections)
    db.SetMaxIdleConns(cfg.MaxDBConnections / 2)
    db.SetConnMaxLifetime(cfg.DBTimeout)
    
    return db, db.Ping()
}
