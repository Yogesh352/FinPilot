package service

import (
	"context"
	"log"
	"stock-api/internal/api"
	"stock-api/internal/repository"
)

type TransactionService struct {
	rowsClient *api.RowsClient
    transactionRepo *repository.TransactionsRepository
}

func NewTransactionService(rowsClient *api.RowsClient, transactionRepo *repository.TransactionsRepository) *TransactionService {
    return &TransactionService{rowsClient: rowsClient, transactionRepo: transactionRepo}
}


func (s *TransactionService) GetBankTransactions(ctx context.Context) error {
	bankTransactionRows, err := s.rowsClient.FetchRows(ctx, "3EOtDHWVB88EJRLjimH0Zq", "3a5b5c99-8f7d-417c-ac85-87f0adf50b53")
	if err != nil {
		log.Printf("failed to get transactionData data: %v", err)
	}

	err = s.transactionRepo.StoreTransactions(bankTransactionRows)
	if err != nil {
		return err
	}
	log.Printf("Successfully stored transaction data")

	return nil
}
