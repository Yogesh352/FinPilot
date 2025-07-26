package util
import (
	"stock-api/internal/api"
)

func StrPtr(s string) *string {
	return &s
}

func FloatPtr(f float64) *float64 {
	return &f
}

func DerefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func Calculate5YRevenueGrowth(income *api.IncomeStatement) float64 {
	if len(income.AnnualReports) < 6 {
		return 0
	}
	currRev := *api.ParseFloat(income.AnnualReports[0].TotalRevenue)
	pastRev := *api.ParseFloat(income.AnnualReports[5].TotalRevenue)
	if pastRev == 0 {
		return 0
	}
	return ((currRev - pastRev) / pastRev) * 100
}

func CalculateHistoricalROE(income *api.IncomeStatement, balance *api.BalanceSheet) map[string]float64 {
	roe := make(map[string]float64)
	for i := 0; i < 5 && i < len(income.AnnualReports) && i < len(balance.AnnualReports); i++ {
		year := income.AnnualReports[i].FiscalDateEnding[:4]
		netInc := *api.ParseFloat(income.AnnualReports[i].NetIncome)
		equity := *api.ParseFloat(balance.AnnualReports[i].TotalShareholderEquity)
		if equity != 0 {
			roe[year] = (netInc / equity) * 100
		}
	}
	return roe
}
