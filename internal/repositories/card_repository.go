package repositories

import (
	"CardFlow/internal/models"
	"context"
	"errors"

	"github.com/google/uuid"
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
    FindCardsByID(ctx context.Context, id uuid.UUID)([]models.Card, error)
    FindCardByID(ctx context.Context, data models.GetCardReq)(models.Card, error)
}


func (r *cardRepository)CreateCard(ctx context.Context, data *models.Card) error{
    return r.db.WithContext(ctx).Create(data).Error
}

func (r *cardRepository) FindCardsByID(ctx context.Context, id uuid.UUID) ([]models.Card, error) {
    var Cards []models.Card

    err := r.db.WithContext(ctx).Where("user_id = ?", id).Find(&Cards).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }

    return Cards, nil
}

func (r *cardRepository) FindCardByID(ctx context.Context, data models.GetCardReq)(models.Card, error){
    var Card models.Card

    err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", data.CardId, data.UserId ).First(&Card).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return models.Card{}, nil
        }
        return models.Card{}, err
    }

    return Card, nil
}