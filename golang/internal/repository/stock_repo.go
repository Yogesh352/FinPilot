package repository

import (
    "database/sql"
)

type StockRepository struct {
    db *sql.DB
}

func NewStockRepository(db *sql.DB) *StockRepository {
    return &StockRepository{db: db}
}

func (r *StockRepository) GetLatestPrice(symbol string) (float64, error) {
    var price float64
    err := r.db.QueryRow("SELECT close FROM stocks_raw WHERE symbol=$1 ORDER BY date DESC LIMIT 1", symbol).Scan(&price)
    return price, err
}
