package service

import (
    "context"
    "time"
    "stock-api/internal/repository"
    "stock-api/internal/util"
)

type StockService struct {
    repo *repository.StockRepository
    scoreRepo *repository.StockScoreRepository
}

func NewStockService(repo *repository.StockRepository, scoreRepo *repository.StockScoreRepository) *StockService {
    return &StockService{repo: repo, scoreRepo: scoreRepo}
}

func (s *StockService) GetStockSummary(symbol string) (float64, error) {
    return s.repo.GetLatestPrice(symbol)
}

// GetStockData gets historical stock data for a symbol within a date range
func (s *StockService) GetStockData(symbol string, startDate, endDate time.Time) ([]repository.StockIntraDayData, error) {
    return s.repo.GetStockData(symbol, startDate, endDate)
}

// GetLatestStockData gets the most recent stock data for a symbol
func (s *StockService) GetLatestStockData(symbol string) (*repository.StockIntraDayData, error) {
    return s.repo.GetLatestStockData(symbol)
}

// GetSymbols gets all unique symbols in the database
func (s *StockService) GetSymbols() ([]string, error) {
    return s.repo.GetSymbols()
}

// GetStockMetadata gets metadata for a specific symbol
func (s *StockService) GetStockMetadata(symbol string) (*repository.StockMetadata, error) {
    return s.repo.GetStockMetadata(symbol)
}

// GetAllStockMetadata gets metadata for all symbols
func (s *StockService) GetAllStockMetadata() ([]repository.StockMetadata, error) {
    return s.repo.GetAllStockMetadata()
}

// StoreStockMetadata stores metadata for a symbol
func (s *StockService) StoreStockMetadata(metadata *repository.StockMetadata) error {
    return s.repo.StoreStockMetadata(metadata)
}

// GetSymbolsWithMetadata gets all symbols that have metadata
func (s *StockService) GetSymbolsWithMetadata() ([]string, error) {
    return s.repo.GetSymbolsWithMetadata()
}

// DeleteStockMetadata deletes metadata for a specific symbol
func (s *StockService) DeleteStockMetadata(symbol string) error {
    return s.repo.DeleteStockMetadata(symbol)
}


func (s *StockService) CalculateLongTermScoreCard(ctx context.Context, symbols []string) error{
    for _, symbol := range symbols {
        overview, _ := s.scoreRepo.GetOverview(symbol)
        income, _ := s.scoreRepo.GetIncomeStatement(symbol)
        balance, _ := s.scoreRepo.GetBalanceSheet(symbol)

        card := repository.Scorecard{
            Symbol:          symbol,
            CompanyName:     overview.Name,
            PERatio:         overview.PERatio,
            PEGRatio:        overview.PEGRatio,
            PriceToBook:     overview.PriceToBook,
            ROE_TTM:         overview.ReturnOnEquityTTM,
            OperatingMargin: overview.OperatingMargin,
            ProfitMargin:    overview.ProfitMargin,
            DividendYield:   overview.DividendYield,
            Beta:            overview.Beta,
            Revenue5YGrowth: util.Calculate5YRevenueGrowth(income),
            HistoricalROE:   util.CalculateHistoricalROE(income, balance),
        }
        err :=  s.scoreRepo.StoreStockScorecard(&card);
        if err != nil {
            return err
        }
    }
    return nil
}
