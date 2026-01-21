package repositories

import (
	"CardFlow/internal/models"
	"context"
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
    Create(ctx context.Context, user *models.User) error
    FindByEmail(ctx context.Context,email string) (*models.User, error)
    FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
    UpdateUserOTP(ctx context.Context, userID uuid.UUID, otp string) error
    Update(ctx context.Context, user *models.User) error
    FindUsersByIDs(ctx context.Context, ids []uuid.UUID)([]models.User, error)

}
    

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByEmail(ctx context.Context,email string) (*models.User, error) {
    var user models.User

    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }

    return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
    var user models.User

    err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }

    return &user, nil
}

func (r *userRepository) UpdateUserOTP(ctx context.Context, userID uuid.UUID, otp string) error { 
    return r.db.Model(&models.User{}).WithContext(ctx).Where("id = ?", userID).Updates(map[string]interface{}{
        "otp": otp,
        "otp_expires_at": gorm.Expr("NOW() + INTERVAL '10 minutes'"),
    }).Error
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
    return r.db.WithContext(ctx).Save(user).Error
}

// func (r *userRepository)FindUsers(id uuid.UUID) ([]models.User, error){
//     var user []models.User

//     err := r.db.Where("id = ?", id).Find(&user).Error
//     if err != nil {
//         if errors.Is(err, gorm.ErrRecordNotFound) {
//             return nil, nil
//         }
//         return nil, err
//     }

//     return user, nil
// }


func (r *userRepository) FindUsersByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]models.User, error) {
	var users []models.User

	err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&users).Error

	return users, err
}
