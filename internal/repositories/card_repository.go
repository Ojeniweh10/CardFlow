package repositories

import (
	"CardFlow/internal/models"
	"context"
	"errors"
	"time"

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
    Update(ctx context.Context, card models.Card) error
    FindCardsExpiringBetween(ctx context.Context, start, end time.Time) ([]models.Card, error)
    ExpireCardsBetween(ctx context.Context, start, end time.Time ) ([]models.Card, error)
    FindCardsByReference(ctx context.Context, data models.WebhookReq) (models.Card, error)
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

func (r *cardRepository) Update(ctx context.Context, card models.Card) error {
    return r.db.WithContext(ctx).Save(&card).Error
}


func (r *cardRepository) FindCardsExpiringBetween(ctx context.Context, start, end time.Time,) ([]models.Card, error) {
	var cards []models.Card

	err := r.db.WithContext(ctx).
		Where("expires_at >= ? AND expires_at < ? AND status = ?", start, end, "active").
		Find(&cards).Error

	return cards, err
}



func (r *cardRepository) ExpireCardsBetween(ctx context.Context, start, end time.Time,) ([]models.Card, error) {
	var cards []models.Card
	err := r.db.WithContext(ctx).
		Model(&models.Card{}).
		Where("expires_at >= ? AND expires_at < ? AND status <> ?", start, end, "expired").
		Update("status", "expired").
		Scan(&cards).Error

	return cards, err
}


func(r *cardRepository)FindCardsByReference(ctx context.Context, data models.WebhookReq)(models.Card, error){
    var Card models.Card
    err := r.db.WithContext(ctx).Where("card_reference = ?", data.CardReference ).First(&Card).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return models.Card{}, nil
        }
        return models.Card{}, err
    }

    return Card, nil
}