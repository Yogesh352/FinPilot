package handler

import (
    "encoding/json"
    "net/http"
    "time"
    "stock-api/internal/service"
)

type StockHandler struct {
    service *service.StockService
}

func NewStockHandler(s *service.StockService) *StockHandler {
    return &StockHandler{service: s}
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
