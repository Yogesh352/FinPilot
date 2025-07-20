package api

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "time"
)

type PolygonClient struct  { 
	*APIClient
}

// NewFinnhubClient creates a new Finnhub client
func NewPolygonClient(apiKey string, timeout time.Duration) *PolygonClient  {
	return &PolygonClient{
        APIClient: NewAPIClient("https://api.polygon.io/v2", apiKey, 60, timeout), // Finnhub allows 60 requests per minute
    }
}


type PolygonAgg struct {
    Close     float64 `json:"c"`
    High      float64 `json:"h"`
    Low       float64 `json:"l"`
    Open      float64 `json:"o"`
    Timestamp int64   `json:"t"`
    Volume    float64   `json:"v"`
}

type PolygonAggResponse struct {
    Ticker        string        `json:"ticker"`
    QueryCount    int           `json:"queryCount"`
    ResultsCount  int           `json:"resultsCount"`
    Results       []PolygonAgg  `json:"results"`
    Status        string        `json:"status"`
    Adjusted      bool          `json:"adjusted"`
    NextURL       string        `json:"next_url,omitempty"`
    RequestID     string        `json:"request_id"`
}

func (c *PolygonClient) GetIntradayBars(ctx context.Context, symbol string, from, to time.Time, intervalMinutes int) ([]PolygonAgg, string ,error) {
    url := fmt.Sprintf(
        "https://api.polygon.io/v2/aggs/ticker/%s/range/%d/minute/%s/%s?adjusted=true&sort=asc&limit=50000&apiKey=%s",
        symbol,
        intervalMinutes,
        from.Format("2006-01-02"),
        to.Format("2006-01-02"),
        c.apiKey,
    )
	
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, symbol, fmt.Errorf("failed to create request: %w", err)
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, symbol, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        body, _ := io.ReadAll(resp.Body)
        return nil, symbol, fmt.Errorf("API error: %s", string(body))
    }

    var polygonResp PolygonAggResponse
    if err := json.NewDecoder(resp.Body).Decode(&polygonResp); err != nil {
        return nil, symbol, fmt.Errorf("failed to decode response: %w", err)
    }

    log.Printf("Fetched %d bars for %s", len(polygonResp.Results), symbol)
    return polygonResp.Results, symbol, nil
}
