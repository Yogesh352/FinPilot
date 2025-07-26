package handler

import (
    "encoding/json"
    "net/http"
    "time"
    "context"
    "stock-api/internal/service"
    "stock-api/internal/repository"
    "stock-api/internal/util"
)

type StockHandler struct {
    service *service.StockService
}

func NewStockHandler(s *service.StockService) *StockHandler {
    return &StockHandler{service: s}
}

type CalculateStockScoreCardRequest struct {
    Symbols []string `json:"symbols"`
}

func (h *StockHandler) GetStockSummary(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    symbol := r.URL.Query().Get("symbol")
    if symbol == "" {
        http.Error(w, "symbol is required", http.StatusBadRequest)
        return
    }

    price, err := h.service.GetStockSummary(symbol)
    if err != nil {
        http.Error(w, "could not get stock summary", http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "symbol":      symbol,
        "latest_price": price,
        "timestamp":   time.Now(),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// GetStockData gets historical stock data for a symbol
func (h *StockHandler) GetStockData(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    symbol := r.URL.Query().Get("symbol")
    if symbol == "" {
        http.Error(w, "symbol is required", http.StatusBadRequest)
        return
    }

    // Parse date range parameters
    startDateStr := r.URL.Query().Get("start")
    endDateStr := r.URL.Query().Get("end")

    var startDate, endDate time.Time
    var err error

    if startDateStr != "" {
        startDate, err = time.Parse("2006-01-02", startDateStr)
        if err != nil {
            http.Error(w, "invalid start date format (use YYYY-MM-DD)", http.StatusBadRequest)
            return
        }
    } else {
        // Default to 30 days ago
        startDate = time.Now().AddDate(0, 0, -30)
    }

    if endDateStr != "" {
        endDate, err = time.Parse("2006-01-02", endDateStr)
        if err != nil {
            http.Error(w, "invalid end date format (use YYYY-MM-DD)", http.StatusBadRequest)
            return
        }
    } else {
        // Default to today
        endDate = time.Now()
    }

    // Get stock data
    data, err := h.service.GetStockData(symbol, startDate, endDate)
    if err != nil {
        http.Error(w, "could not get stock data", http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "symbol":    symbol,
        "start_date": startDate.Format("2006-01-02"),
        "end_date":   endDate.Format("2006-01-02"),
        "data":      data,
        "count":     len(data),
        "timestamp": time.Now(),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// GetStockMetadata gets metadata for a specific symbol
func (h *StockHandler) GetStockMetadata(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    symbol := r.URL.Query().Get("symbol")
    if symbol == "" {
        http.Error(w, "symbol is required", http.StatusBadRequest)
        return
    }

    metadata, err := h.service.GetStockMetadata(symbol)
    if err != nil {
        http.Error(w, "could not get stock metadata", http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "metadata":  metadata,
        "timestamp": time.Now(),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// GetAllStockMetadata gets metadata for all symbols
func (h *StockHandler) GetAllStockMetadata(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    metadata, err := h.service.GetAllStockMetadata()
    if err != nil {
        http.Error(w, "could not get stock metadata", http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "metadata":  metadata,
        "count":     len(metadata),
        "timestamp": time.Now(),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// StoreStockMetadata stores metadata for a symbol
func (h *StockHandler) StoreStockMetadata(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var metadata struct {
        Symbol        string  `json:"symbol"`
        CompanyName   string  `json:"company_name"`
        Industry      string  `json:"industry"`
        Exchange      string  `json:"exchange"`
        Currency      string  `json:"currency"`
        MarketCap     float64 `json:"market_cap"`
        // PE            float64 `json:"pe_ratio"`
        // DividendYield float64 `json:"dividend_yield"`
        Description   string  `json:"description"`
        Website       string  `json:"website"`
        Type          string  `json:"type"`
    }

    if err := json.NewDecoder(r.Body).Decode(&metadata); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if metadata.Symbol == "" {
        http.Error(w, "symbol is required", http.StatusBadRequest)
        return
    }

    stockMetadata := &repository.StockMetadata{
        Symbol:        metadata.Symbol,
        CompanyName:   util.StrPtr(metadata.CompanyName),
        Industry:      util.StrPtr(metadata.Industry),
        Exchange:      metadata.Exchange,
        Currency:      metadata.Currency,
        MarketCap:     util.FloatPtr(metadata.MarketCap),
        // PE:            metadata.PE,
        // DividendYield: metadata.DividendYield,
        Description:   metadata.Description,
        Website:       util.StrPtr(metadata.Website),
    }

    err := h.service.StoreStockMetadata(stockMetadata)
    if err != nil {
        http.Error(w, "could not store stock metadata", http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "message":   "Stock metadata stored successfully",
        "symbol":    metadata.Symbol,
        "timestamp": time.Now(),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// DeleteStockMetadata deletes metadata for a specific symbol
func (h *StockHandler) DeleteStockMetadata(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    symbol := r.URL.Query().Get("symbol")
    if symbol == "" {
        http.Error(w, "symbol is required", http.StatusBadRequest)
        return
    }

    err := h.service.DeleteStockMetadata(symbol)
    if err != nil {
        http.Error(w, "could not delete stock metadata", http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "message":   "Stock metadata deleted successfully",
        "symbol":    symbol,
        "timestamp": time.Now(),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (s *StockHandler) CalculateStockScoreCard(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req CalculateStockScoreCardRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if len(req.Symbols) == 0 {
        http.Error(w, "At least one symbol is required", http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(r.Context(), 15*time.Minute)
    defer cancel()

    err := s.service.CalculateLongTermScoreCard(ctx, req.Symbols)

    response := map[string]interface{}{
        "symbols":   req.Symbols,
        "timestamp": time.Now(),
    }

    if err != nil {
        response["status"] = "error"
        response["message"] = err.Error()
        w.WriteHeader(http.StatusInternalServerError)
    } else {
        response["status"] = "success"
        response["message"] = "Stock Score cards calculated successfully"
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
