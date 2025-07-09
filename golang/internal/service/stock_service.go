package service

import "stock-api/internal/repository"

type StockService struct {
    repo *repository.StockRepository
}

func NewStockService(repo *repository.StockRepository) *StockService {
    return &StockService{repo: repo}
}

func (s *StockService) GetStockSummary(symbol string) (float64, error) {
    return s.repo.GetLatestPrice(symbol)
}
