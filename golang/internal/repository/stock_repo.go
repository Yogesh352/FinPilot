package repository

import (
    "database/sql"
    "fmt"
    "log"
    "time"
)

type StockRepository struct {
    db *sql.DB
}

// StockData represents a complete stock data record
type StockData struct {
    Symbol    string    `json:"symbol"`
    Date      time.Time `json:"date"`
    Open      float64   `json:"open"`
    High      float64   `json:"high"`
    Low       float64   `json:"low"`
    Close     float64   `json:"close"`
    Volume    float64   `json:"volume"`
    CreatedAt time.Time `json:"created_at"`
}

func NewStockRepository(db *sql.DB) *StockRepository {
    return &StockRepository{db: db}
}

// StoreStockData stores a complete stock data record
func (r *StockRepository) StoreStockData(symbol string, date time.Time, open, high, low, close, volume float64) error {
    log.Printf("Attempting to store data for %s on %s", symbol, date.Format("2006-01-02"))
    
    query := `
        INSERT INTO stocks_raw (symbol, date, open, high, low, close, volume, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (symbol, date) DO UPDATE SET
            open = EXCLUDED.open,
            high = EXCLUDED.high,
            low = EXCLUDED.low,
            close = EXCLUDED.close,
            volume = EXCLUDED.volume,
            updated_at = $8
    `
    
    result, err := r.db.Exec(query, symbol, date, open, high, low, close, volume, time.Now())
    if err != nil {
        log.Printf("Database error storing data for %s on %s: %v", symbol, date.Format("2006-01-02"), err)
        return fmt.Errorf("failed to store stock data: %w", err)
    }
    
    rowsAffected, _ := result.RowsAffected()
    log.Printf("Successfully stored data for %s on %s (rows affected: %d)", symbol, date.Format("2006-01-02"), rowsAffected)
    
    return nil
}

// GetLatestPrice gets the latest close price for a symbol
func (r *StockRepository) GetLatestPrice(symbol string) (float64, error) {
    var price float64
    err := r.db.QueryRow(`
        SELECT close 
        FROM stocks_raw 
        WHERE symbol = $1 
        ORDER BY date DESC 
        LIMIT 1
    `, symbol).Scan(&price)
    
    if err != nil {
        return 0, fmt.Errorf("failed to get latest price for %s: %w", symbol, err)
    }
    
    return price, nil
}

// GetStockData gets stock data for a symbol within a date range
func (r *StockRepository) GetStockData(symbol string, startDate, endDate time.Time) ([]StockData, error) {
    query := `
        SELECT symbol, date, open, high, low, close, volume, created_at
        FROM stocks_raw
        WHERE symbol = $1 AND date BETWEEN $2 AND $3
        ORDER BY date DESC
    `
    
    rows, err := r.db.Query(query, symbol, startDate, endDate)
    if err != nil {
        return nil, fmt.Errorf("failed to query stock data: %w", err)
    }
    defer rows.Close()
    
    var data []StockData
    for rows.Next() {
        var record StockData
        err := rows.Scan(
            &record.Symbol,
            &record.Date,
            &record.Open,
            &record.High,
            &record.Low,
            &record.Close,
            &record.Volume,
            &record.CreatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan stock data: %w", err)
        }
        data = append(data, record)
    }
    
    return data, nil
}

// GetLatestStockData gets the most recent stock data for a symbol
func (r *StockRepository) GetLatestStockData(symbol string) (*StockData, error) {
    query := `
        SELECT symbol, date, open, high, low, close, volume, created_at
        FROM stocks_raw
        WHERE symbol = $1
        ORDER BY date DESC
        LIMIT 1
    `
    
    var record StockData
    err := r.db.QueryRow(query, symbol).Scan(
        &record.Symbol,
        &record.Date,
        &record.Open,
        &record.High,
        &record.Low,
        &record.Close,
        &record.Volume,
        &record.CreatedAt,
    )
    
    if err != nil {
        return nil, fmt.Errorf("failed to get latest stock data for %s: %w", symbol, err)
    }
    
    return &record, nil
}

// GetSymbols gets all unique symbols in the database
func (r *StockRepository) GetSymbols() ([]string, error) {
    query := `SELECT DISTINCT symbol FROM stocks_raw ORDER BY symbol`
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to query symbols: %w", err)
    }
    defer rows.Close()
    
    var symbols []string
    for rows.Next() {
        var symbol string
        if err := rows.Scan(&symbol); err != nil {
            return nil, fmt.Errorf("failed to scan symbol: %w", err)
        }
        symbols = append(symbols, symbol)
    }
    
    return symbols, nil
}

// DeleteOldData deletes stock data older than the specified date
func (r *StockRepository) DeleteOldData(beforeDate time.Time) error {
    query := `DELETE FROM stocks_raw WHERE date < $1`
    
    result, err := r.db.Exec(query, beforeDate)
    if err != nil {
        return fmt.Errorf("failed to delete old data: %w", err)
    }
    
    rowsAffected, _ := result.RowsAffected()
    fmt.Printf("Deleted %d old records", rowsAffected)
    
    return nil
}
