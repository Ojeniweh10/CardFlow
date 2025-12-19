package repositories

import (
	"CardFlow/internal/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

type UserRepository interface {
    Create(user *models.User) error
    FindByEmail(email string) (*models.User, error)
    FindByID(id string) (*models.User, error)
    UpdateUserOTP(userID uuid.UUID, otp string) error
    Update(user *models.User) error
}


func (r *userRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User

    err := r.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }

    return &user, nil
}

func (r *userRepository) FindByID(id string) (*models.User, error) {
    var user models.User

    err := r.db.Where("id = ?", id).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }

    return &user, nil
}

func (r *userRepository) UpdateUserOTP(userID uuid.UUID, otp string) error { 
    return r.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
        "otp": otp,
        "otp_expires_at": gorm.Expr("NOW() + INTERVAL '10 minutes'"),
    }).Error
}

func (r *userRepository) Update(user *models.User) error {
    return r.db.Save(user).Error
}