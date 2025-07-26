package handler

import (
    "context"
    "encoding/json"
    "net/http"
    "time"
    "stock-api/internal/service"
)

// ExtractionHandler handles data extraction endpoints
type ExtractionHandler struct {
    extractionService *service.DataExtractionService
}

// NewExtractionHandler creates a new extraction handler
func NewExtractionHandler(es *service.DataExtractionService) *ExtractionHandler {
    return &ExtractionHandler{extractionService: es}
}

// ExtractStockDataRequest represents the request for extracting stock data
type ExtractStockDataRequest struct {
    Symbol string `json:"symbol"`
    From time.Time `json:"from`
    To time.Time `json:"to`
}

type ExtractByExchangeRequest struct {
    Exchange string `json:"exchange"`
}

type ExtractBySymbolsRequest struct {
    Symbols []string `json:"symbols"`
}


// ExtractStockDataResponse represents the response for extraction
type ExtractStockDataResponse struct {
    Symbol    string    `json:"symbol"`
    Status    string    `json:"status"`
    Message   string    `json:"message"`
    Timestamp time.Time `json:"timestamp"`
}

// ExtractStockData extracts stock data for a given symbol
func (h *ExtractionHandler) ExtractStockData(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req ExtractStockDataRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if req.Symbol == "" {
        http.Error(w, "Symbol is required", http.StatusBadRequest)
        return
    }

    // Create context with timeout
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
    defer cancel()

    // Extract data
    err := h.extractionService.ExtractAndStoreStockData(ctx, req.Symbol, req.From, req.To)
    
    response := ExtractStockDataResponse{
        Symbol:    req.Symbol,
        Timestamp: time.Now(),
    }

    if err != nil {
        response.Status = "error"
        response.Message = err.Error()
        w.WriteHeader(http.StatusInternalServerError)
    } else {
        response.Status = "success"
        response.Message = "Data extraction completed successfully"
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// ExtractLatestQuote extracts the latest quote for a symbol
// func (h *ExtractionHandler) ExtractLatestQuote(w http.ResponseWriter, r *http.Request) {
//     if r.Method != http.MethodPost {
//         http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//         return
//     }

//     var req ExtractStockDataRequest
//     if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//         http.Error(w, "Invalid request body", http.StatusBadRequest)
//         return
//     }

//     if req.Symbol == "" {
//         http.Error(w, "Symbol is required", http.StatusBadRequest)
//         return
//     }

//     // Create context with timeout
//     ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
//     defer cancel()

//     // Extract latest quote
//     err := h.extractionService.ExtractLatestQuote(ctx, req.Symbol)
    
//     response := ExtractStockDataResponse{
//         Symbol:    req.Symbol,
//         Timestamp: time.Now(),
//     }

//     if err != nil {
//         response.Status = "error"
//         response.Message = err.Error()
//         w.WriteHeader(http.StatusInternalServerError)
//     } else {
//         response.Status = "success"
//         response.Message = "Latest quote extracted successfully"
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(response)
// }

// BatchExtractDataRequest represents the request for batch extraction
type BatchExtractDataRequest struct {
    Symbols []string `json:"symbols"`
    From    time.Time `json:"from`
    To    time.Time `json:"to`
}

// BatchExtractDataResponse represents the response for batch extraction
type BatchExtractDataResponse struct {
    Status    string    `json:"status"`
    Message   string    `json:"message"`
    Timestamp time.Time `json:"timestamp"`
    Processed int       `json:"processed"`
    Failed    int       `json:"failed"`
}

// BatchExtractData extracts data for multiple symbols
func (h *ExtractionHandler) BatchExtractData(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req BatchExtractDataRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if len(req.Symbols) == 0 {
        http.Error(w, "At least one symbol is required", http.StatusBadRequest)
        return
    }

    // Create context with timeout (longer for batch operations)
    ctx, cancel := context.WithTimeout(r.Context(), 30*time.Minute)
    defer cancel()

    // Extract data for all symbols
    err := h.extractionService.BatchExtractData(ctx, req.Symbols, req.From, req.To)
    
    response := BatchExtractDataResponse{
        Status:    "success",
        Message:    "Batch extraction completed",
        Timestamp:  time.Now(),
        Processed:  len(req.Symbols),
        Failed:     0, // This would need to be tracked in the service
    }

    if err != nil {
        response.Status = "error"
        response.Message = err.Error()
        w.WriteHeader(http.StatusInternalServerError)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// GetExtractionStatus returns the status of data extraction
func (h *ExtractionHandler) GetExtractionStatus(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    symbol := r.URL.Query().Get("symbol")
    if symbol == "" {
        http.Error(w, "Symbol is required", http.StatusBadRequest)
        return
    }

    // For now, return a simple status
    // In a real application, you might want to track extraction jobs
    response := map[string]interface{}{
        "symbol":    symbol,
        "status":    "unknown", // Would need job tracking system
        "timestamp": time.Now(),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
} 

func (h *ExtractionHandler) ExtractStockMetadata(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req ExtractByExchangeRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if req.Exchange == "" {
        http.Error(w, "Exchange is required", http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
    defer cancel()

    err := h.extractionService.ExtractAndStoreStockMetaData(ctx, req.Exchange)

    response := map[string]interface{}{
        "exchange":  req.Exchange,
        "timestamp": time.Now(),
    }

    if err != nil {
        response["status"] = "error"
        response["message"] = err.Error()
        w.WriteHeader(http.StatusInternalServerError)
    } else {
        response["status"] = "success"
        response["message"] = "Stock metadata extracted successfully"
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *ExtractionHandler) ExtractSymbols(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req ExtractByExchangeRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if req.Exchange == "" {
        http.Error(w, "Exchange is required", http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(r.Context(), 15*time.Minute)
    defer cancel()

    err := h.extractionService.ExtractAndStoreSymbols(ctx, req.Exchange)

    response := map[string]interface{}{
        "exchange":  req.Exchange,
        "timestamp": time.Now(),
    }

    if err != nil {
        response["status"] = "error"
        response["message"] = err.Error()
        w.WriteHeader(http.StatusInternalServerError)
    } else {
        response["status"] = "success"
        response["message"] = "Symbols extracted successfully"
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}


func (h *ExtractionHandler) ExtractCompanyProfiles(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req ExtractBySymbolsRequest
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

    err := h.extractionService.ExtractAndStoreCompanyData(ctx, req.Symbols)

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
        response["message"] = "Company profiles extracted successfully"
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *ExtractionHandler) ExtractCompanyOverviews(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req ExtractBySymbolsRequest
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

    err := h.extractionService.BatchExtractAndStoreStockOverview(ctx, req.Symbols)

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
        response["message"] = "Company overviews extracted successfully"
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *ExtractionHandler) ExtractCompanyIncomeStatements(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req ExtractBySymbolsRequest
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

    err := h.extractionService.BatchExtractAndStoreStockIncomeStatment(ctx, req.Symbols)

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
        response["message"] = "Company income statments extracted successfully"
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *ExtractionHandler) ExtractCompanyBalanceSheets(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req ExtractBySymbolsRequest
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

    err := h.extractionService.BatchExtractAndStoreStockBalanceStatement(ctx, req.Symbols)

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
        response["message"] = "Company balance sheets extracted successfully"
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}