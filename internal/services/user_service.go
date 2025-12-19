package services

import (
	"CardFlow/internal/models"
	"CardFlow/internal/repositories"
	"CardFlow/internal/utils"
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserService interface {
    RegisterUser(models.CreateUserRequest) error
	Login(models.LoginReq) (string, error)
	VerifyEmail(userID uuid.UUID) error
	VerifyOtp(userID uuid.UUID, otp string) error
}

type userService struct {
    repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
    return &userService{repo: repo}
}


func (s *userService) RegisterUser(req models.CreateUserRequest) error {
    existingUser, err := s.repo.FindByEmail(req.Email)
	if existingUser != nil {
		return errors.New("user with this email already exists")
	}
	if err != nil {
		//log the error to notify devs then return a generic error message
		return errors.New("something went wrong, please try again later")
	}
    hashedPassword, err := utils.Hash(req.Password)
	if err != nil {
		//log the error to notify devs then return a generic error message
		return errors.New("something went wrong, please try again later")
	}

    user := &models.User{
        Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
    }

    return s.repo.Create(user)
}

func (s *userService)Login(req models.LoginReq) (string, error){
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return "", errors.New("something went wrong, please try again later")
	}
	if user == nil {
		return "", errors.New("invalid email or password")
	}

	err = utils.CompareHashAndPassword(user.PasswordHash, req.Password)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", errors.New("something went wrong, please try again later")
	}
	return token, nil
}

func (s *userService) VerifyEmail(userID uuid.UUID) error {
	user,  err := s.repo.FindByID(userID.String())
	if err != nil {
		return errors.New("something went wrong, please try again later")
	}
	if user == nil {
		return errors.New("user not found")
	}
	otp, err := utils.GenerateOTP()
	if err != nil {
		return errors.New("something went wrong, please try again later")
	}

	err = utils.SendEmailOTP(user.Email, otp)
	if err != nil {
		return errors.New("failed to send OTP email")
	}

	err = s.repo.UpdateUserOTP(userID, otp)
	if err != nil {
		return errors.New("something went wrong, please try again later")
	}

	return nil

}

func (s *userService) VerifyOtp(userID uuid.UUID, otp string) error {
	user,  err := s.repo.FindByID(userID.String())
	if err != nil {
		return errors.New("something went wrong, please try again later")
	}
	if user == nil {
		return errors.New("user not found")
	}

	if user.OTP != otp {
		return errors.New("invalid OTP")
	}

	if user.OTPExpiresAt.IsZero() || user.OTPExpiresAt.Before(time.Now()) {
		return errors.New("OTP has expired")
	}

	// Mark email as verified
	user.EmailVerified = true
	user.OTP = ""
	user.OTPExpiresAt = time.Time{}

	err = s.repo.Update(user)
	if err != nil {
		return errors.New("something went wrong, please try again later")
	}

	return nil
}