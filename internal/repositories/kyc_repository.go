package repositories

import (
	"CardFlow/internal/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)



type kycRepository struct {
	db *gorm.DB
}


func NewKycRepository(db *gorm.DB) KycRepository{
   return &kycRepository{db: db}
}

type KycRepository interface{
	FindByUserID(userID uuid.UUID)(*models.KYCSubmission, error)
	CreateKycSubmission(KycUser *models.KYCSubmission) error
	CreateKycDocsSubmission(KycDocs *models.KYCDocument) error
	UpdateKycSubmission(kyc *models.KYCSubmission) error
	RunInTransaction(ctx context.Context, fn func(repo KycRepository) error) error
}

func (r *kycRepository) RunInTransaction(ctx context.Context, fn func(repo KycRepository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &kycRepository{db: tx}
		return fn(txRepo)
	})
}



func (r *kycRepository)FindByUserID(userID uuid.UUID)(*models.KYCSubmission, error){
	var profile models.KYCSubmission
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }

    return &profile, nil
}

func (r *kycRepository)CreateKycSubmission(KycUser *models.KYCSubmission) error{
    return r.db.Create(KycUser).Error
}

func (r *kycRepository) CreateKycDocsSubmission(doc *models.KYCDocument) error {
	if err := r.db.Create(doc).Error; err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return errors.New("document already uploaded")
			}
		}
	}
	return nil
}

func (r *kycRepository) UpdateKycSubmission(kyc *models.KYCSubmission) error {
	return r.db.Model(&models.KYCSubmission{}).Where("user_id = ?", kyc.UserID).Updates(map[string]interface{}{
		"status": kyc.Status,
	}).Error
}