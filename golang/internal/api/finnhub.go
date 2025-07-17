package api

import (
    "context"
    "fmt"
    "log"
    "time"
)

// FinnhubClient represents the Finnhub API client
type FinnhubClient struct  { 
	*APIClient
}

// NewFinnhubClient creates a new Finnhub client
func NewFinnhubClient(apiKey string, timeout time.Duration) *FinnhubClient  {
	return &FinnhubClient{
        APIClient: NewAPIClient("https://finnhub.io/api/v1", apiKey, 60, timeout), // Finnhub allows 60 requests per minute
    }
}

// StockSymbol represents a stock symbol from Finnhub
type StockSymbol struct {
    Currency        string `json:"currency`   
	Description     string `json:description`
    DisplaySymbol   string `json:"displaySymbol"`
    Figi            string `json:"figi"`
    Mic             string `json:"mic"`
    Symbol          string `json:"symbol"`
    Type            string `json:"type"`
}

// CompanyProfile represents company profile from Finnhub
type CompanyProfile struct {
    Country                string  `json:"country"`
    Currency               string  `json:"currency"`
    Exchange               string  `json:"exchange"`
    Ipo                    string  `json:"ipo"`
    MarketCapitalization   float64 `json:"marketCap"`
    Name                   string  `json:"name"`
    Phone                  string  `json:"phone"`
    ShareOutstanding       float64 `json:"shareOutstanding"`
    Ticker                 string  `json:"ticker"`
    Weburl                 string  `json:"weburl"`                   
	Logo 				   string  `json:"logo"`
    FinnhubIndustry        string  `json:"finnhubIndustry"`
}

// // Quote represents a stock quote from Finnhub
// type Quote struct [object Object]{
//     CurrentPrice  float64 `json:"c"`
//     Change        float64 `json:"d"`
//     PercentChange float64:"json: "`
//     HighPrice     float64 `json:"h"`
//     LowPrice      float64 `json:l
//     OpenPrice     float64 `json:"o"`
//     PreviousClose float64 `json:"pc"`
// }

// GetStockSymbols retrieves stock symbols for a specific exchange
func (c *FinnhubClient) GetStockSymbols(ctx context.Context, exchange string) ([]StockSymbol, error) {
	log.Printf("Fetching stock symbols for exchange: %s" ,exchange)
	
    req := &Request{
        Method: "GET",
        Path:  "/stock/symbol",
        Query: map[string]string{
            "exchange": exchange,
            "token":    c.apiKey,
        },
    }

    var symbols []StockSymbol
    if err := c.DoJSON(ctx, req, &symbols); err != nil {
        return nil, fmt.Errorf("failed to get stock symbols: %w", err)
    }

    log.Printf("Successfully fetched %d stock symbols for exchange %s", len(symbols), exchange)
    return symbols, nil
}

func (c *FinnhubClient) GetCompanyProfile(ctx context.Context, symbol string) (*CompanyProfile, error)  {
	log.Printf("Fetching company profile for symbol: %s", symbol)
    
    req := &Request{
        Method: "GET",        
		Path: "/stock/profile2",        
		Query: map[string]string{
            "symbol": symbol,
            "token":  c.apiKey,
        },
    }

    var profile CompanyProfile
    if err := c.DoJSON(ctx, req, &profile); err != nil {
        return nil, fmt.Errorf("failed to get company profile: %w", err)
    }

    log.Printf("Successfully fetched company profile for %s: %s", symbol, profile.Name)
    return &profile, nil
}

// // GetQuote retrieves current quote for a symbol
// func (c *FinnhubClient) GetQuote(ctx context.Context, symbol string) (*Quote, error) [object Object] log.Printf("Fetching quote for symbol: %s", symbol)
    
//     req := &Request{
//         Method: "GET,        Path:  /quote",
//         Query: map[string]string{
//             symbolymbol,
//             token":  c.apiKey,
//         },
//     }

//     var quote Quote
//     if err := c.DoJSON(ctx, req, &quote); err != nil {
//         return nil, fmt.Errorf("failed to get quote: %w, err)
//     }

//     log.Printf("Successfully fetched quote for %s: $%.2bol, quote.CurrentPrice)
//     return &quote, nil
// }

// // GetLatestPrice gets the latest price for a symbol
// func (c *FinnhubClient) GetLatestPrice(ctx context.Context, symbol string) (float64, error) {
//     quote, err := c.GetQuote(ctx, symbol)
//     if err != nil [object Object]
//         return0
//     }

//     return quote.CurrentPrice, nil
// } 