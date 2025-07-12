package service

import (
    "context"
    "fmt"
    "log"
    "time"
    "stock-api/internal/api"
    "stock-api/internal/repository"
)

// DataExtractionService handles fetching and storing financial data from external APIs
type DataExtractionService struct {
    alphaVantageClient *api.AlphaVantageClient
    stockRepo          *repository.StockRepository
}

// NewDataExtractionService creates a new data extraction service
func NewDataExtractionService(alphaVantageClient *api.AlphaVantageClient, stockRepo *repository.StockRepository) *DataExtractionService {
    return &DataExtractionService{
        alphaVantageClient: alphaVantageClient,
        stockRepo:          stockRepo,
    }
}

// ExtractAndStoreStockData fetches stock data from external API and stores it in the database
func (s *DataExtractionService) ExtractAndStoreStockData(ctx context.Context, symbol string) error {
    log.Printf("Starting data extraction for symbol: %s", symbol)

    // Get time series data from Alpha Vantage
    timeSeries, err := s.alphaVantageClient.GetDailyTimeSeries(ctx, symbol)
    if err != nil {
        return fmt.Errorf("failed to get time series data: %w", err)
    }

    log.Printf("Processing %d data points for symbol %s", len(timeSeries.TimeSeries), symbol)

    // Process and store each data point
    storedCount := 0
    errorCount := 0
    
    for date, data := range timeSeries.TimeSeries {
        log.Printf("Processing data for %s on %s", symbol, date)
        
        // Parse the date
        parsedDate, err := time.Parse("2006-01-02", date)
        if err != nil {
            log.Printf("Failed to parse date %s for symbol %s: %v", date, symbol, err)
            errorCount++
            continue
        }

        // Parse numeric values
        open, err := api.ParseFloat(data.Open)
        if err != nil {
            log.Printf("Failed to parse open price for %s on %s: %v", symbol, date, err)
            errorCount++
            continue
        }

        high, err := api.ParseFloat(data.High)
        if err != nil {
            log.Printf("Failed to parse high price for %s on %s: %v", symbol, date, err)
            errorCount++
            continue
        }

        low, err := api.ParseFloat(data.Low)
        if err != nil {
            log.Printf("Failed to parse low price for %s on %s: %v", symbol, date, err)
            errorCount++
            continue
        }

        close, err := api.ParseFloat(data.Close)
        if err != nil {
            log.Printf("Failed to parse close price for %s on %s: %v", symbol, date, err)
            errorCount++
            continue
        }

        volume, err := api.ParseFloat(data.Volume)
        if err != nil {
            log.Printf("Failed to parse volume for %s on %s: %v", symbol, date, err)
            errorCount++
            continue
        }

        log.Printf("Storing data for %s on %s: O=%.2f, H=%.2f, L=%.2f, C=%.2f, V=%.0f", 
            symbol, date, open, high, low, close, volume)

        // Store in database
        err = s.stockRepo.StoreStockData(symbol, parsedDate, open, high, low, close, volume)
        if err != nil {
            log.Printf("Failed to store data for %s on %s: %v", symbol, date, err)
            errorCount++
            continue
        }
        
        storedCount++
        log.Printf("Successfully stored data for %s on %s", symbol, date)
    }

    log.Printf("Completed data extraction for symbol: %s - Stored: %d, Errors: %d", symbol, storedCount, errorCount)
    
    if storedCount == 0 {
        return fmt.Errorf("no data was stored for symbol %s", symbol)
    }
    
    return nil
}

// ExtractLatestQuote fetches and stores the latest quote for a symbol
func (s *DataExtractionService) ExtractLatestQuote(ctx context.Context, symbol string) error {
    log.Printf("Extracting latest quote for symbol: %s", symbol)
    
    quote, err := s.alphaVantageClient.GetStockQuote(ctx, symbol)
    if err != nil {
        return fmt.Errorf("failed to get quote: %w", err)
    }

    // Parse the latest trading day
    tradingDay, err := time.Parse("2006-01-02", quote.LatestTradingDay)
    if err != nil {
        return fmt.Errorf("failed to parse trading day: %w", err)
    }

    // Parse numeric values
    price, err := api.ParseFloat(quote.Price)
    if err != nil {
        return fmt.Errorf("failed to parse price: %w", err)
    }

    volume, err := api.ParseFloat(quote.Volume)
    if err != nil {
        return fmt.Errorf("failed to parse volume: %w", err)
    }

    log.Printf("Storing latest quote for %s: Date=%s, Price=%.2f, Volume=%.0f", 
        symbol, tradingDay.Format("2006-01-02"), price, volume)

    // Store in database (using the same price for open, high, low, close for latest quote)
    err = s.stockRepo.StoreStockData(symbol, tradingDay, price, price, price, price, volume)
    if err != nil {
        return fmt.Errorf("failed to store quote data: %w", err)
    }

    log.Printf("Successfully stored latest quote for %s: $%.2f", symbol, price)
    return nil
}

// BatchExtractData extracts data for multiple symbols
func (s *DataExtractionService) BatchExtractData(ctx context.Context, symbols []string) error {
    log.Printf("Starting batch extraction for %d symbols", len(symbols))
    
    for i, symbol := range symbols {
        log.Printf("Processing symbol %d/%d: %s", i+1, len(symbols), symbol)
        
        // Add delay between requests to respect rate limits
        if i > 0 {
            log.Printf("Waiting 12 seconds before next request...")
            time.Sleep(12 * time.Second) // Alpha Vantage allows 5 requests per minute
        }

        err := s.ExtractAndStoreStockData(ctx, symbol)
        if err != nil {
            log.Printf("Failed to extract data for %s: %v", symbol, err)
            continue
        }
    }
    
    log.Printf("Completed batch extraction for %d symbols", len(symbols))
    return nil
} 