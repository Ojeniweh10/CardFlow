package repositories

import (
	"CardFlow/internal/models"
	"context"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}


func NewTransactionRepository(db *gorm.DB) TransactionRepository{
   return &transactionRepository{db: db}
}

type TransactionRepository interface{
	CreateTransaction(ctx context.Context, data *models.Transaction) error
	CreateLedger(ctx context.Context, data models.BalanceLedger)error
}


func (r *transactionRepository)CreateTransaction(ctx context.Context, data *models.Transaction) error{
    return r.db.WithContext(ctx).Create(data).Error
}

func(r *transactionRepository)CreateLedger(ctx context.Context, data models.BalanceLedger)error{
	return r.db.WithContext(ctx).Create(data).Error
}