package main

import (
    "log"
    "net/http"
    "os"

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

    stockRepo := repository.NewStockRepository(db)
    stockService := service.NewStockService(stockRepo)
    stockHandler := handler.NewStockHandler(stockService)

    mux := http.NewServeMux()
    mux.HandleFunc("/api/stocks", stockHandler.GetStockSummary)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server running on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, mux))
}
