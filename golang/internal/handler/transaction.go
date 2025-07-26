package handler

import (
    "context"
    "encoding/json"
    "net/http"
    "time"
    "stock-api/internal/service"
)

type TransactionHandler struct {
    transactionService *service.TransactionService
}

func NewTransactionHandler(ts *service.TransactionService) *TransactionHandler {
    return &TransactionHandler{transactionService: ts}
}

func (h *TransactionHandler) ExtractTransactions(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()

    err := h.transactionService.GetBankTransactions(ctx)
    if err != nil {
        http.Error(w, "could not get transactions", http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "status":    "Success",
        "timestamp": time.Now(),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
