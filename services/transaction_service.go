package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
	"log/slog"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
	logger *slog.Logger
}

func NewTransactionService(repo *repositories.TransactionRepository, logger *slog.Logger) *TransactionService {
	return &TransactionService{repo: repo, logger: logger}
}


func (s *TransactionService) Checkout(items []models.CheckoutItem, useLock bool) (*models.Transaction, error){
	return s.repo.CreateTransaction(items, useLock )
}


//method receiver itu argumen pertama s itu receiver variable kayak this di python misalnya
//pointer ke struct transactionservice