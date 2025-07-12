package api

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
    "golang.org/x/time/rate"
)

// APIClient represents a generic API client with rate limiting
type APIClient struct {
    client          *http.Client
    rateLimiter     *rate.Limiter
    baseURL         string
    apiKey          string
    requestTimeout  time.Duration
}

// NewAPIClient creates a new API client with rate limiting
func NewAPIClient(baseURL, apiKey string, requestsPerMinute int, timeout time.Duration) *APIClient {
    return &APIClient{
        client: &http.Client{
            Timeout: timeout,
        },
        rateLimiter:    rate.NewLimiter(rate.Limit(float64(requestsPerMinute)/60.0), 1),
        baseURL:        baseURL,
        apiKey:         apiKey,
        requestTimeout: timeout,
    }
}

// Request represents a generic API request
type Request struct {
    Method  string
    Path    string
    Query   map[string]string
    Headers map[string]string
    Body    io.Reader
}

// Response represents a generic API response
type Response struct {
    StatusCode int
    Body       []byte
    Headers    http.Header
}

// Do performs an API request with rate limiting
func (c *APIClient) Do(ctx context.Context, req *Request) (*Response, error) {
    // Wait for rate limiter
    if err := c.rateLimiter.Wait(ctx); err != nil {
        return nil, fmt.Errorf("rate limiter error: %w", err)
    }

    // Build URL
    url := c.baseURL + req.Path
    if len(req.Query) > 0 {
        url += "?"
        for key, value := range req.Query {
            if url[len(url)-1] != '?' {
                url += "&"
            }
            url += fmt.Sprintf("%s=%s", key, value)
        }
    }

    // Create HTTP request
    httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, req.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    // Add headers
    for key, value := range req.Headers {
        httpReq.Header.Set(key, value)
    }

    // Add API key if present
    if c.apiKey != "" {
        httpReq.Header.Set("X-API-Key", c.apiKey)
    }

    // Make request
    resp, err := c.client.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    // Read response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %w", err)
    }

    return &Response{
        StatusCode: resp.StatusCode,
        Body:       body,
        Headers:    resp.Header,
    }, nil
}

// DoJSON performs a request and unmarshals the response into the provided interface
func (c *APIClient) DoJSON(ctx context.Context, req *Request, v interface{}) error {
    resp, err := c.Do(ctx, req)
    if err != nil {
        return err
    }

    if resp.StatusCode >= 400 {
        return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(resp.Body))
    }

    return json.Unmarshal(resp.Body, v)
} 