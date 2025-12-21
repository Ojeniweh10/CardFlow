package repositories

import (
	"CardFlow/internal/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)



type kycRepository struct {
	db *gorm.DB
}


func NewKycRepository(db *gorm.DB) KycRepository{
   return &kycRepository{db: db}
}

type KycRepository interface{
	FindByID(userID uuid.UUID)(*models.KYCProfile, error)
	CreateKycProfile(KycUser *models.KYCProfile) error
}

func (r *kycRepository)FindByID(userID uuid.UUID)(*models.KYCProfile, error){
	var profile models.KYCProfile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }

    return &profile, nil
}

func (r *kycRepository)CreateKycProfile(KycUser *models.KYCProfile) error{
    return r.db.Create(KycUser).Error
}