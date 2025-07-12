package main

import (
    "log"
    "net/http"
    "os"
    "time"

    "stock-api/internal/api"
    "stock-api/internal/config"
    "stock-api/internal/handler"
    "stock-api/internal/repository"
    "stock-api/internal/service"

    _ "github.com/lib/pq"
)

func main() {
    cfg := config.Load()

    db, err := config.SetupPostgres(cfg)
    if err != nil {
        log.Fatalf("could not connect to postgres: %v", err)
    }
    defer db.Close()

    // Initialize API clients
    alphaVantageClient := api.NewAlphaVantageClient(cfg.AlphaVantageAPIKey, cfg.APIRequestTimeout)

    // Initialize repositories and services
    stockRepo := repository.NewStockRepository(db)
    stockService := service.NewStockService(stockRepo)
    dataExtractionService := service.NewDataExtractionService(alphaVantageClient, stockRepo)

    // Initialize handlers
    stockHandler := handler.NewStockHandler(stockService)
    extractionHandler := handler.NewExtractionHandler(dataExtractionService)

    // Setup routes
    mux := http.NewServeMux()
    
    // Stock data endpoints
    mux.HandleFunc("/api/stocks", stockHandler.GetStockSummary)
    mux.HandleFunc("/api/stocks/data", stockHandler.GetStockData)
    
    // Data extraction endpoints
    mux.HandleFunc("/api/extract/stock", extractionHandler.ExtractStockData)
    mux.HandleFunc("/api/extract/quote", extractionHandler.ExtractLatestQuote)
    mux.HandleFunc("/api/extract/batch", extractionHandler.BatchExtractData)
    mux.HandleFunc("/api/extract/status", extractionHandler.GetExtractionStatus)

    // Health check endpoint
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"status": "healthy", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server running on :%s", port)
    log.Printf("Available endpoints:")
    log.Printf("  GET  /health - Health check")
    log.Printf("  GET  /api/stocks?symbol=AAPL - Get stock summary")
    log.Printf("  GET  /api/stocks/data?symbol=AAPL&start=2024-01-01&end=2024-12-31 - Get stock data")
    log.Printf("  POST /api/extract/stock - Extract stock data")
    log.Printf("  POST /api/extract/quote - Extract latest quote")
    log.Printf("  POST /api/extract/batch - Batch extract data")
    log.Printf("  GET  /api/extract/status?symbol=AAPL - Get extraction status")
    
    log.Fatal(http.ListenAndServe(":"+port, mux))
}
