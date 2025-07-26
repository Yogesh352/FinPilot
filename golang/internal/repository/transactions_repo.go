package repository

import (
	"database/sql"
	"log"
	"stock-api/internal/api"
)

type TransactionsRepository struct {
    db *sql.DB
}

func NewTransactionsRepository(db *sql.DB) *TransactionsRepository {
    return &TransactionsRepository{db: db}
}

func (r *TransactionsRepository) StoreTransactions(transactions []api.Transaction) error {
	query := `
		INSERT INTO personal_transactions (date, amount, currency, description, category, bank, account)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	for _, tx := range transactions {
		_, err := r.db.Exec(query,
			tx.Date, tx.Amount, tx.Currency, tx.Description,
			tx.Category, tx.Bank, tx.Account,
		)
		if err != nil {
			log.Printf("Failed to insert transaction: %+v → %v", tx, err)
			return err
		}
	}

	return nil
}
