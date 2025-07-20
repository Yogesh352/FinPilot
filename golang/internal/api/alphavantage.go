package api

import (
    "context"
    "fmt"
    "log"
    "strconv"
    "time"
)

// AlphaVantageClient represents the Alpha Vantage API client
type AlphaVantageClient struct {
    *APIClient
}

// NewAlphaVantageClient creates a new Alpha Vantage client
func NewAlphaVantageClient(apiKey string, timeout time.Duration) *AlphaVantageClient {
    return &AlphaVantageClient{
        APIClient: NewAPIClient("https://www.alphavantage.co", apiKey, 5, timeout), // Alpha Vantage has strict rate limits
    }
}

// StockQuote represents a stock quote from Alpha Vantage
type StockQuote struct {
    Symbol        string  `json:"01. symbol"`
    Open          string  `json:"02. open"`
    High          string  `json:"03. high"`
    Low           string  `json:"04. low"`
    Price         string  `json:"05. price"`
    Volume        string  `json:"06. volume"`
    LatestTradingDay string `json:"07. latest trading day"`
    PreviousClose string  `json:"08. previous close"`
    Change        string  `json:"09. change"`
    ChangePercent string  `json:"10. change percent"`
}

// StockQuoteResponse represents the Alpha Vantage quote response
type StockQuoteResponse struct {
    GlobalQuote StockQuote `json:"Global Quote"`
}

// TimeSeriesData represents time series data from Alpha Vantage
type TimeSeriesData struct {
    Open   string `json:"1. open"`
    High   string `json:"2. high"`
    Low    string `json:"3. low"`
    Close  string `json:"4. close"`
    Volume string `json:"5. volume"`
}

// TimeSeriesResponse represents the Alpha Vantage time series response
type TimeSeriesResponse struct {
    MetaData struct {
        Information   string `json:"1. Information"`
        Symbol       string `json:"2. Symbol"`
        LastRefreshed string `json:"3. Last Refreshed"`
        OutputSize   string `json:"4. Output Size"`
        TimeZone     string `json:"5. Time Zone"`
    } `json:"Meta Data"`
    TimeSeries map[string]TimeSeriesData `json:"Time Series (5min)"`
}

// GetStockQuote retrieves the latest stock quote
func (c *AlphaVantageClient) GetStockQuote(ctx context.Context, symbol string) (*StockQuote, error) {
    log.Printf("Fetching stock quote for symbol: %s", symbol)
    
    req := &Request{
        Method: "GET",
        Path:   "/query",
        Query: map[string]string{
            "function": "GLOBAL_QUOTE",
            "symbol":   symbol,
            "apikey":   c.apiKey,
        },
    }

    var response StockQuoteResponse
    if err := c.DoJSON(ctx, req, &response); err != nil {
        return nil, fmt.Errorf("failed to get stock quote: %w", err)
    }

    log.Printf("Successfully fetched quote for %s: $%s", symbol, response.GlobalQuote.Price)
    return &response.GlobalQuote, nil
}

func (c *AlphaVantageClient) GetIntradayTimeSeries(ctx context.Context, symbol string) (*TimeSeriesResponse, error) {
    log.Printf("Fetching daily time series for symbol: %s", symbol)
    
    req := &Request{
        Method: "GET",
        Path:   "/query",
        Query: map[string]string{
            "function": "TIME_SERIES_INTRADAY",
            "interval": "5min",
            "symbol":   symbol,
            "apikey":   c.apiKey,
        },
    }

    var response TimeSeriesResponse
    if err := c.DoJSON(ctx, req, &response); err != nil {
        return nil, fmt.Errorf("failed to get time series: %w", err)
    }

    log.Printf("Successfully fetched time series for %s with %d data points", symbol, len(response.TimeSeries))
    return &response, nil
}

// ParseFloat safely parses a string to float64
func ParseFloat(s string) (float64, error) {
    return strconv.ParseFloat(s, 64)
}

// GetLatestPrice gets the latest price for a symbol
func (c *AlphaVantageClient) GetLatestPrice(ctx context.Context, symbol string) (float64, error) {
    quote, err := c.GetStockQuote(ctx, symbol)
    if err != nil {
        return 0, err
    }

    return ParseFloat(quote.Price)
} 