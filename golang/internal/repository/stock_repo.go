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

// StockMetadata represents stock metadata information
type StockMetadata struct {
    Symbol        string    `json:"symbol"`
    CompanyName   *string    `json:"company_name"`
    Industry      *string    `json:"industry"`
    Exchange      string    `json:"exchange"`
    Currency      string    `json:"currency"`
    MarketCap     *float64   `json:"market_cap"`
    Description   string    `json:"description"`
    Website       *string    `json:"website"`
    Type          string    `json:"type"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
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

// StoreStockMetadata stores stock metadata information
func (r *StockRepository) StoreStockMetadata(metadata *StockMetadata) error {
    log.Printf("Storing metadata for symbol: %s", metadata.Symbol)
    
    query := `
        INSERT INTO stocks_metadata (
            symbol, company_name, industry, exchange, currency,
            market_cap, description, website, type, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        ON CONFLICT (symbol) DO UPDATE SET
            symbol = EXCLUDED.symbol,
            company_name = EXCLUDED.company_name,
            industry = EXCLUDED.industry,
            exchange = EXCLUDED.exchange,
            currency = EXCLUDED.currency,
            market_cap = EXCLUDED.market_cap,
            description = EXCLUDED.description,
            website = EXCLUDED.website,
            type = EXCLUDED.type,
            created_at = EXCLUDED.created_at,
            updated_at = EXCLUDED.updated_at
    `
    
    now := time.Now()
    result, err := r.db.Exec(query,
        metadata.Symbol,
        &metadata.CompanyName,
        &metadata.Industry,
        metadata.Exchange,
        metadata.Currency,
        &metadata.MarketCap,
        metadata.Description,
        &metadata.Website,
        metadata.Type,
        now,
        now,
    )
    
    if err != nil {
        log.Printf("Database error storing metadata for %s: %v", metadata.Symbol, err)
        return fmt.Errorf("failed to store stock metadata: %w", err)
    }
    
    rowsAffected, _ := result.RowsAffected()
    log.Printf("Successfully stored metadata for %s (rows affected: %d)", metadata.Symbol, rowsAffected)
    
    return nil
}

// GetStockMetadata retrieves metadata for a specific symbol
func (r *StockRepository) GetStockMetadata(symbol string) (*StockMetadata, error) {
    query := `
        SELECT symbol, company_name, industry, exchange, currency,
               market_cap, description, website, type, created_at, updated_at
        FROM stocks_metadata
        WHERE symbol = $1
    `
    
    var metadata StockMetadata
    err := r.db.QueryRow(query, symbol).Scan(
        &metadata.Symbol,
        &metadata.CompanyName,
        &metadata.Industry,
        &metadata.Exchange,
        &metadata.Currency,
        &metadata.MarketCap,
        &metadata.Description,
        &metadata.Website,
        &metadata.Type,
        &metadata.CreatedAt,
        &metadata.UpdatedAt,
    )
    
    if err != nil {
        return nil, fmt.Errorf("failed to get metadata for %s: %w", symbol, err)
    }
    
    return &metadata, nil
}

// GetAllStockMetadata retrieves metadata for all symbols
func (r *StockRepository) GetAllStockMetadata() ([]StockMetadata, error) {
    query := `
        SELECT symbol, company_name, industry, exchange, currency,
               market_cap, description, website, created_at, updated_at
        FROM stocks_metadata
        ORDER BY symbol
    `
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to query stock metadata: %w", err)
    }
    defer rows.Close()
    
    var metadata []StockMetadata
    for rows.Next() {
        var record StockMetadata
        err := rows.Scan(
            &record.Symbol,
            &record.CompanyName,
            &record.Industry,
            &record.Exchange,
            &record.Currency,
            &record.MarketCap,
            &record.Description,
            &record.Website,
            &record.Type,
            &record.CreatedAt,
            &record.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan metadata: %w", err)
        }
        metadata = append(metadata, record)
    }
    
    return metadata, nil
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

// GetSymbolsWithMetadata gets all symbols that have metadata
func (r *StockRepository) GetSymbolsWithMetadata() ([]string, error) {
    query := `SELECT DISTINCT symbol FROM stocks_metadata ORDER BY symbol`
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to query symbols with metadata: %w", err)
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

func (r *StockRepository) GetSymbolWithMetadata(symbol string) (string, error) {
	query := `SELECT symbol FROM stocks_metadata WHERE symbol = $1 LIMIT 1`

	var result string
	err := r.db.QueryRow(query, symbol).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("symbol %s not found in metadata", symbol)
		}
		return "", fmt.Errorf("failed to query symbol %s: %w", symbol, err)
	}

	return result, nil
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

// DeleteStockMetadata deletes metadata for a specific symbol
func (r *StockRepository) DeleteStockMetadata(symbol string) error {
    query := `DELETE FROM stocks_metadata WHERE symbol = $1`
    
    result, err := r.db.Exec(query, symbol)
    if err != nil {
        return fmt.Errorf("failed to delete metadata for %s: %w", symbol, err)
    }
    
    rowsAffected, _ := result.RowsAffected()
    log.Printf("Deleted metadata for %s (rows affected: %d)", symbol, rowsAffected)
    
    return nil
}
