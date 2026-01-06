package services

import (
	"CardFlow/internal/models"
	"CardFlow/internal/repositories"
	"CardFlow/internal/utils"
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

type UserService interface {
    RegisterUser(ctx context.Context, req models.CreateUserRequest) error
	Login(ctx context.Context, req models.LoginReq) (string, error)
	MFALogin(ctx context.Context, req models.MFALoginReq)(string , error)
	VerifyEmail(ctx context.Context,userID uuid.UUID) error
	VerifyOtp(ctx context.Context, userID uuid.UUID, otp string) error
	EnableMFA(ctx context.Context, userID uuid.UUID)(string, error)
	VerifyMFA(ctx context.Context, userID uuid.UUID, data string) error
}

type userService struct {
    repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
    return &userService{repo: repo}
}


func (s *userService) RegisterUser(ctx context.Context, req models.CreateUserRequest) error {
    existingUser, err := s.repo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return errors.New("user with this email already exists")
	}
	if err != nil {
		//log the error to notify devs then return a generic error message
		log.Println(err)
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

    return s.repo.Create(ctx, user)
}

func (s *userService)Login(ctx context.Context, req models.LoginReq) (string, error){
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("something went wrong, please try again later")
	}
	if user == nil {
		return "", errors.New("invalid email or password")
	}
	if user.MFAEnabled {
		return "", errors.New("MFA required")
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

func (s *userService) MFALogin(ctx context.Context,req models.MFALoginReq) (string, error) {
    user, err := s.repo.FindByEmail(ctx, req.Email)
    if err != nil {
        return "", errors.New("something went wrong, please try again later")
    }
    if user == nil {
        return "", errors.New("invalid email or password")
    }

    if !user.MFAEnabled {
        return "", errors.New("MFA is not enabled for this user")
    }

    err = utils.ValidateTotp(req.TOTPCode, user.MFASecret)
    if err != nil {
        return "", err
    }

    token, err := utils.GenerateJWT(user.ID, user.Email)
    if err != nil {
        return "", errors.New("something went wrong, please try again later")
    }

    return token, nil
}


func (s *userService) VerifyEmail(ctx context.Context,userID uuid.UUID) error {
	user,  err := s.repo.FindByID(ctx, userID)
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

	err = s.repo.UpdateUserOTP(ctx, userID, otp)
	if err != nil {
		return errors.New("something went wrong, please try again later")
	}

	return nil

}

func (s *userService) VerifyOtp(ctx context.Context, userID uuid.UUID, otp string) error {
	user,  err := s.repo.FindByID(ctx, userID)
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

	err = s.repo.Update(ctx, user)
	if err != nil {
		return errors.New("something went wrong, please try again later")
	}

	return nil
}

func (s *userService) EnableMFA(ctx context.Context, userID uuid.UUID) (string, error){
	user,  err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return"",  errors.New("something went wrong, please try again later")
	}
	if user == nil {
		return "" ,  errors.New("user not found")
	}
	secret , otpURL , err := utils.GenerateMFASecret(userID)
	if err != nil {
		return "",  errors.New("something went wrong, please try again later")
	}
	user.MFASecret = secret
	err = s.repo.Update(ctx,user)
	if err != nil {
		return "", errors.New("something went wrong, please try again later")
	}

	return otpURL, nil
}

func (s *userService) VerifyMFA(ctx context.Context, userID uuid.UUID, data string) error{
	user,  err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("something went wrong, please try again later")
	}
	if user == nil {
		return errors.New("user not found")
	}
	secret := user.MFASecret
	err = utils.ValidateTotp(data, secret)
	if err != nil {
		return err
	}	
	user.MFAEnabled = true
	err = s.repo.Update(ctx, user)
	if err != nil {
		return errors.New("something went wrong, please try again later")
	}
	return nil
}