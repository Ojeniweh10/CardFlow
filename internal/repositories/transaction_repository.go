package repositories

import (
	"CardFlow/internal/models"
	"context"
	"errors"

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
	FindTxnByReference(ctx context.Context, reference string)(models.Transaction, error)
	Update(ctx context.Context, card models.Transaction) error
	FindByIdempotencyKey(ctx context.Context, idempotencykey string)(models.Transaction, error)
    FindCardTransactions(ctx context.Context, data models.GetCardTransactionsReq)([]models.Transaction, error)
}


func (r *transactionRepository)CreateTransaction(ctx context.Context, data *models.Transaction) error{
    return r.db.WithContext(ctx).Create(&data).Error
}

func(r *transactionRepository)CreateLedger(ctx context.Context, data models.BalanceLedger)error{
	return r.db.WithContext(ctx).Create(&data).Error
}

func (r *transactionRepository)FindTxnByReference(ctx context.Context, reference string)(models.Transaction, error){
	var Txn models.Transaction

    err := r.db.WithContext(ctx).Where("transaction_reference = ? ", reference).First(&Txn).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return models.Transaction{}, nil
        }
        return models.Transaction{}, err
    }

    return Txn, nil
}

func (r *transactionRepository)FindCardTransactions(ctx context.Context, data models.GetCardTransactionsReq)([]models.Transaction, error){
    var Txn []models.Transaction
    err := r.db.WithContext(ctx).Where("card_id = ? AND user_id = ? ", data.Cardid, data.Userid).Find(&Txn).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return []models.Transaction{}, nil
        }
        return []models.Transaction{}, err
    }
    return Txn, nil
}

func (r *transactionRepository)FindByIdempotencyKey(ctx context.Context, idempotencykey string)(models.Transaction, error){
	var Txn models.Transaction

    err := r.db.WithContext(ctx).Where("idempotency_key = ? ", idempotencykey).First(&Txn).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return models.Transaction{}, nil
        }
        return models.Transaction{}, err
    }

    return Txn, nil
}


func (r *transactionRepository) Update(ctx context.Context, data models.Transaction) error {
    return r.db.WithContext(ctx).Save(&data).Error
}