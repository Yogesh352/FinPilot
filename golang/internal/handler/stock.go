package handler

import (
    "encoding/json"
    "net/http"
    "stock-api/internal/service"
)

type StockHandler struct {
    service *service.StockService
}

func NewStockHandler(s *service.StockService) *StockHandler {
    return &StockHandler{service: s}
}

func (h *StockHandler) GetStockSummary(w http.ResponseWriter, r *http.Request) {
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

    json.NewEncoder(w).Encode(map[string]interface{}{
        "symbol": symbol,
        "latest_price": price,
    })
}
