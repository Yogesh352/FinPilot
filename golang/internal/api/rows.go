package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"time"
)

type RowsAPIResponse struct {
    Items [][]struct {
        Col   int         `json:"col"`
        Row   int         `json:"row"`
        Value interface{} `json:"value"`
    } `json:"items"`
}

type Transaction struct {
	Date        time.Time
	Amount      float64
	Currency    string
	Description string
	Category    string
	Bank        string
	Account     string
}

type RowsClient struct {
    *APIClient
}

func NewRowsClient(apiKey string, timeout time.Duration) *RowsClient {
    return &RowsClient{
        APIClient: NewAPIClient("https://api.rows.com", apiKey, 5, timeout), // Alpha Vantage has strict rate limits
    }
}

func (c *RowsClient) FetchRows(ctx context.Context, spreadsheetId, tableId string) ([]Transaction, error) {
    url := fmt.Sprintf("https://api.rows.com/v1/spreadsheets/%s/tables/%s/cells/A1:G100", spreadsheetId, tableId)

    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    req.Header.Set("Authorization", "Bearer "+c.apiKey)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("http request failed: %w", err)
    }
    defer resp.Body.Close()

    var apiResp RowsAPIResponse
    if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
        return nil, fmt.Errorf("failed to decode JSON: %w", err)
    }

    if len(apiResp.Items) < 2 {
        return nil, fmt.Errorf("not enough rows to parse data")
    }

    // Build column index → name mapping
    colMap := map[int]string{}
    for _, cell := range apiResp.Items[0] {
        if str, ok := cell.Value.(string); ok {
            colMap[cell.Col] = str
        }
    }

    var transactions []Transaction

    for _, row := range apiResp.Items[1:] {
        tx := Transaction{}

        for _, cell := range row {
            column := colMap[cell.Col]
            switch column {
            case "Date":
                if str, ok := cell.Value.(string); ok {
                    t, _ := time.Parse("2006-01-02", str)
                    tx.Date = t
                }
            case "Amount":
                switch v := cell.Value.(type) {
                case float64:
                    tx.Amount = v
                case string:
                    fmt.Sscanf(v, "%f", &tx.Amount)
                }
            case "Currency":
                tx.Currency = fmt.Sprint(cell.Value)
            case "Description":
                tx.Description = fmt.Sprint(cell.Value)
            case "Personal Finance Category Primary":
                tx.Category = fmt.Sprint(cell.Value)
            case "Bank":
                tx.Bank = fmt.Sprint(cell.Value)
            case "Account":
                tx.Account = fmt.Sprint(cell.Value)
            }
        }
		log.Printf("Transaction %s", tx.Account)
        transactions = append(transactions, tx)
    }

    return transactions, nil
}
