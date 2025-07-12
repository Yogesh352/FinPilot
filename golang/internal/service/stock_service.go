package service

import (
    "time"
    "stock-api/internal/repository"
)

type StockService struct {
    repo *repository.StockRepository
}

func NewStockService(repo *repository.StockRepository) *StockService {
    return &StockService{repo: repo}
}

func (s *StockService) GetStockSummary(symbol string) (float64, error) {
    return s.repo.GetLatestPrice(symbol)
}

// GetStockData gets historical stock data for a symbol within a date range
func (s *StockService) GetStockData(symbol string, startDate, endDate time.Time) ([]repository.StockData, error) {
    return s.repo.GetStockData(symbol, startDate, endDate)
}

// GetLatestStockData gets the most recent stock data for a symbol
func (s *StockService) GetLatestStockData(symbol string) (*repository.StockData, error) {
    return s.repo.GetLatestStockData(symbol)
}

// GetSymbols gets all unique symbols in the database
func (s *StockService) GetSymbols() ([]string, error) {
    return s.repo.GetSymbols()
}
