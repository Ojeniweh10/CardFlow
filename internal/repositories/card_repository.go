package repositories

import (
	"CardFlow/internal/models"
	"context"

	"gorm.io/gorm"
)

type cardRepository struct {
	db *gorm.DB
}


func NewCardRepository(db *gorm.DB) CardRepository{
   return &cardRepository{db: db}
}

type CardRepository interface{
	CreateCard(ctx context.Context, data *models.Card) error
}


func (r *cardRepository)CreateCard(ctx context.Context, data *models.Card) error{
    return r.db.WithContext(ctx).Create(data).Error
}