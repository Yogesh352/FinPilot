package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"stock-api/internal/api"
	"time"
	"log"
)


type StockScoreRepository struct {
    db *sql.DB
}

type IncomeStatementRow struct {
	FiscalDateEnding string
	TotalRevenue     string
	NetIncome        string
}

type BalanceSheetRow struct {
	FiscalDateEnding        string
	ShareholderEquity       string
}

type OverviewRow struct {
	Symbol            string  
	Name              string 
	PERatio           float64 
	PEGRatio          float64
	PriceToBook       float64 
	ReturnOnEquityTTM float64 
	OperatingMargin   float64 
	ProfitMargin      float64 
	DividendYield     float64 
	Beta              float64 
}

type Scorecard struct {
	Symbol          string
	CompanyName     string
	PERatio         float64
	PEGRatio        float64
	PriceToBook     float64
	ROE_TTM         float64
	HistoricalROE   map[string]float64
	Revenue5YGrowth float64
	OperatingMargin float64
	ProfitMargin    float64
	DividendYield   float64
	Beta            float64
}

func NewStockScoreRepository(db *sql.DB) *StockScoreRepository {
    return &StockScoreRepository{db: db}
}

func (r *StockScoreRepository) StoreStockScorecard(card *Scorecard) error {
	historicalROEJson, err := json.Marshal(card.HistoricalROE)
	if err != nil {
		return fmt.Errorf("failed to marshal historical ROE: %w", err)
	}

	query := `
		INSERT INTO stock_scorecards (
			symbol, company_name, pe_ratio, peg_ratio, price_to_book, 
			roe_ttm, revenue_5y_growth, operating_margin, profit_margin, 
			dividend_yield, beta, historical_roe, updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, $11, $12, CURRENT_TIMESTAMP
		)
		ON CONFLICT (symbol) DO UPDATE SET
			company_name = EXCLUDED.company_name,
			pe_ratio = EXCLUDED.pe_ratio,
			peg_ratio = EXCLUDED.peg_ratio,
			price_to_book = EXCLUDED.price_to_book,
			roe_ttm = EXCLUDED.roe_ttm,
			revenue_5y_growth = EXCLUDED.revenue_5y_growth,
			operating_margin = EXCLUDED.operating_margin,
			profit_margin = EXCLUDED.profit_margin,
			dividend_yield = EXCLUDED.dividend_yield,
			beta = EXCLUDED.beta,
			historical_roe = EXCLUDED.historical_roe,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err = r.db.Exec(query,
		card.Symbol,
		card.CompanyName,
		card.PERatio,
		card.PEGRatio,
		card.PriceToBook,
		card.ROE_TTM,
		card.Revenue5YGrowth,
		card.OperatingMargin,
		card.ProfitMargin,
		card.DividendYield,
		card.Beta,
		historicalROEJson,
	)

	if err != nil {
		return err
	}
	return nil
}

func (r *StockScoreRepository) GetStockScorecard(symbol string) (*Scorecard, error) {
	query := `
		SELECT symbol, company_name, pe_ratio, peg_ratio, price_to_book, 
		       roe_ttm, revenue_5y_growth, operating_margin, profit_margin, 
		       dividend_yield, beta, historical_roe
		FROM stock_scorecards
		WHERE symbol = $1
	`

	row := r.db.QueryRow(query, symbol)

	var card Scorecard
	var historicalROEJson []byte

	err := row.Scan(
		&card.Symbol,
		&card.CompanyName,
		&card.PERatio,
		&card.PEGRatio,
		&card.PriceToBook,
		&card.ROE_TTM,
		&card.Revenue5YGrowth,
		&card.OperatingMargin,
		&card.ProfitMargin,
		&card.DividendYield,
		&card.Beta,
		&historicalROEJson,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch stock scorecard: %w", err)
	}

	err = json.Unmarshal(historicalROEJson, &card.HistoricalROE)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal historical ROE: %w", err)
	}

	return &card, nil
}

func (r *StockScoreRepository) StoreOverview(o *api.OverviewResponse) error {
	query := `
		INSERT INTO stock_overviews (
			symbol, name, pe_ratio, peg_ratio, price_to_book,
			return_on_equity_ttm, operating_margin, profit_margin,
			dividend_yield, beta
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		ON CONFLICT (symbol) DO UPDATE SET
			name = EXCLUDED.name,
			pe_ratio = EXCLUDED.pe_ratio,
			peg_ratio = EXCLUDED.peg_ratio,
			price_to_book = EXCLUDED.price_to_book,
			return_on_equity_ttm = EXCLUDED.return_on_equity_ttm,
			operating_margin = EXCLUDED.operating_margin,
			profit_margin = EXCLUDED.profit_margin,
			dividend_yield = EXCLUDED.dividend_yield,
			beta = EXCLUDED.beta,
			created_at = CURRENT_TIMESTAMP
	`

	_, err := r.db.Exec(query,
		o.Symbol, o.Name, *api.ParseFloat(o.PERatio), *api.ParseFloat(o.PEGRatio), *api.ParseFloat(o.PriceToBook),
		*api.ParseFloat(o.ReturnOnEquityTTM), *api.ParseFloat(o.OperatingMargin), *api.ParseFloat(o.ProfitMargin),
		*api.ParseFloat(o.DividendYield), *api.ParseFloat(o.Beta),
	)
	log.Printf("HERE: %f", *api.ParseFloat(o.ReturnOnEquityTTM))
	return err
}

func (r *StockScoreRepository) GetOverview(symbol string) (*OverviewRow, error) {
	query := `
		SELECT symbol, name, pe_ratio, peg_ratio, price_to_book,
		       return_on_equity_ttm, operating_margin, profit_margin,
		       dividend_yield, beta
		FROM stock_overviews
		WHERE symbol = $1
	`

	var overview OverviewRow

	err := r.db.QueryRow(query, symbol).Scan(
		&overview.Symbol,
		&overview.Name,
		&overview.PERatio,
		&overview.PEGRatio,
		&overview.PriceToBook,
		&overview.ReturnOnEquityTTM,
		&overview.OperatingMargin,
		&overview.ProfitMargin,
		&overview.DividendYield,
		&overview.Beta,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch overview for %s: %w", symbol, err)
	}

	return &overview, nil
}



func (r *StockScoreRepository) StoreIncomeStatement(symbol string, income *api.IncomeStatement) error {
	query := `
		INSERT INTO stock_income_statements (
			symbol, fiscal_date, total_revenue, net_income
		) VALUES ($1, $2, $3, $4)
		ON CONFLICT (symbol, fiscal_date) DO UPDATE SET
			total_revenue = EXCLUDED.total_revenue,
			net_income = EXCLUDED.net_income,
			created_at = CURRENT_TIMESTAMP
	`

	for _, report := range income.AnnualReports {
		parsedDate, _ := time.Parse("2006-01-02", report.FiscalDateEnding)
		_, err := r.db.Exec(query, symbol, parsedDate, report.TotalRevenue, report.NetIncome)
		if err != nil {
			return fmt.Errorf("error storing income statement (%s): %w", report.FiscalDateEnding, err)
		}
	}
	return nil
}

func (r *StockScoreRepository) StoreBalanceSheet(symbol string, balance *api.BalanceSheet) error {
	query := `
		INSERT INTO stock_balance_sheets (
			symbol, fiscal_date, shareholder_equity
		) VALUES ($1, $2, $3)
		ON CONFLICT (symbol, fiscal_date) DO UPDATE SET
			shareholder_equity = EXCLUDED.shareholder_equity,
			created_at = CURRENT_TIMESTAMP
	`

	for _, report := range balance.AnnualReports {
		parsedDate, _ := time.Parse("2006-01-02", report.FiscalDateEnding)
		_, err := r.db.Exec(query, symbol, parsedDate, report.TotalShareholderEquity)
		if err != nil {
			return fmt.Errorf("error storing balance sheet (%s): %w", report.FiscalDateEnding, err)
		}
	}
	return nil
}

func (r *StockScoreRepository) GetIncomeStatement(symbol string) (*api.IncomeStatement, error) {
	query := `
		SELECT fiscal_date, total_revenue, net_income
		FROM stock_income_statements
		WHERE symbol = $1
		ORDER BY fiscal_date DESC
	`
	rows, err := r.db.Query(query, symbol)
	if err != nil {
		return nil, fmt.Errorf("query income statements: %w", err)
	}
	defer rows.Close()

	var reports []struct {
		FiscalDateEnding string `json:"fiscalDateEnding"`
		TotalRevenue     string `json:"totalRevenue"`
		NetIncome        string `json:"netIncome"`
	}

	for rows.Next() {
		var r IncomeStatementRow
		if err := rows.Scan(&r.FiscalDateEnding, &r.TotalRevenue, &r.NetIncome); err != nil {
			return nil, fmt.Errorf("scan income row: %w", err)
		}
		reports = append(reports, struct {
			FiscalDateEnding string `json:"fiscalDateEnding"`
			TotalRevenue     string `json:"totalRevenue"`
			NetIncome        string `json:"netIncome"`
		}{
			FiscalDateEnding: r.FiscalDateEnding,
			TotalRevenue:     r.TotalRevenue,
			NetIncome:        r.NetIncome,
		})
	}

	return &api.IncomeStatement{AnnualReports: reports}, nil
}

func (r *StockScoreRepository) GetBalanceSheet(symbol string) (*api.BalanceSheet, error) {
	query := `
		SELECT fiscal_date, shareholder_equity
		FROM stock_balance_sheets
		WHERE symbol = $1
		ORDER BY fiscal_date DESC
	`
	rows, err := r.db.Query(query, symbol)
	if err != nil {
		return nil, fmt.Errorf("query balance sheets: %w", err)
	}
	defer rows.Close()

	var reports []struct {
		FiscalDateEnding        string `json:"fiscalDateEnding"`
		TotalShareholderEquity string `json:"totalShareholderEquity"`
	}

	for rows.Next() {
		var r BalanceSheetRow
		if err := rows.Scan(&r.FiscalDateEnding, &r.ShareholderEquity); err != nil {
			return nil, fmt.Errorf("scan balance row: %w", err)
		}
		reports = append(reports, struct {
			FiscalDateEnding        string `json:"fiscalDateEnding"`
			TotalShareholderEquity string `json:"totalShareholderEquity"`
		}{
			FiscalDateEnding:        r.FiscalDateEnding,
			TotalShareholderEquity: r.ShareholderEquity,
		})
	}

	return &api.BalanceSheet{AnnualReports: reports}, nil
}

