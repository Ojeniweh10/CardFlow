package services

import (
	"CardFlow/internal/models"
	"CardFlow/internal/utils"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)


type fakeUserRepo struct {
	findUser   *models.User
	findUserByID    *models.User
	findUserByEmail *models.User
	updateErr       error
	findErr    error
	createErr  error
	created   bool
}

func (f *fakeUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return f.findUser, f.findErr
}

func (f *fakeUserRepo) Create(ctx context.Context, user *models.User) error {
	f.created = true
	return f.createErr
}

func (f *fakeUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return f.findUserByID, f.findErr
}

func (f *fakeUserRepo) Update(ctx context.Context, user *models.User) error {
	return f.updateErr
}

func (f *fakeUserRepo) UpdateUserOTP(ctx context.Context, id uuid.UUID, otp string) error {
	return nil
}


func TestRegisterUser_Success(t *testing.T) {
	repo := &fakeUserRepo{
		findUser: nil,
		findErr:  nil,
		createErr: nil,
	}

	svc := &userService{repo: repo}

	req := models.CreateUserRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Phone:     "123456",
	}

	err := svc.RegisterUser(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !repo.created {
		t.Fatalf("expected user to be created, but it was not")
	}
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	repo := &fakeUserRepo{
		findUser: &models.User{Email: "test@example.com"},
	}

	svc := &userService{repo: repo}

	req := models.CreateUserRequest{
		Email: "test@example.com",
	}

	err := svc.RegisterUser(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}


func TestLogin_UserNotFound(t *testing.T) {
	repo := &fakeUserRepo{
		findUser: nil,
		findErr:  nil,
	}

	svc := &userService{repo: repo}

	req := models.LoginReq{
		Email:    "missing@example.com",
		Password: "password",
	}

	_, err := svc.Login(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}


func TestLogin_WrongPassword(t *testing.T) {
	hashed, _ := utils.Hash("correct-password")

	repo := &fakeUserRepo{
		findUser: &models.User{
			ID:           uuid.New(),
			Email:        "test@example.com",
			PasswordHash: hashed,
		},
	}

	svc := &userService{repo: repo}

	req := models.LoginReq{
		Email:    "test@example.com",
		Password: "wrong-password",
	}

	_, err := svc.Login(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestLogin_Success(t *testing.T) {
	hashed, _ := utils.Hash("password123")

	repo := &fakeUserRepo{
		findUser: &models.User{
			ID:           uuid.New(),
			Email:        "test@example.com",
			PasswordHash: hashed,
		},
	}

	svc := &userService{repo: repo}

	req := models.LoginReq{
		Email:    "test@example.com",
		Password: "password123",
	}

	token, err := svc.Login(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatalf("expected token, got empty string")
	}
}


func TestMFALogin_MFANotEnabled(t *testing.T) {
	repo := &fakeUserRepo{
		findUserByEmail: &models.User{
			Email:      "test@example.com",
			MFAEnabled: false,
		},
	}

	svc := &userService{repo: repo}

	_, err := svc.MFALogin(context.Background(), models.MFALoginReq{
		Email:    "test@example.com",
		TOTPCode: "123456",
	})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}


func TestMFALogin_Success(t *testing.T) {
	secret, _, _ := utils.GenerateMFASecret(uuid.New())
	code, _ := utils.GenerateTotp(secret)

	repo := &fakeUserRepo{
		findUser: &models.User{
			ID:         uuid.New(),
			Email:      "test@example.com",
			MFAEnabled: true,
			MFASecret:  secret,
		},
	}

	svc := &userService{repo: repo}

	token, err := svc.MFALogin(context.Background(), models.MFALoginReq{
		Email:    "test@example.com",
		TOTPCode: code,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatalf("expected token, got empty string")
	}
}


func TestVerifyEmail_UserNotFound(t *testing.T) {
	repo := &fakeUserRepo{
		findUserByID: nil,
	}

	svc := &userService{repo: repo}

	err := svc.VerifyEmail(context.Background(), uuid.New())
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}


func TestVerifyOtp_Expired(t *testing.T) {
	repo := &fakeUserRepo{
		findUserByID: &models.User{
			OTP:          "123456",
			OTPExpiresAt: time.Now().Add(-time.Minute),
		},
	}

	svc := &userService{repo: repo}

	err := svc.VerifyOtp(context.Background(), uuid.New(), "123456")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestVerifyOtp_Success(t *testing.T) {
	repo := &fakeUserRepo{
		findUserByID: &models.User{
			OTP:          "123456",
			OTPExpiresAt: time.Now().Add(time.Minute),
		},
	}

	svc := &userService{repo: repo}

	err := svc.VerifyOtp(context.Background(), uuid.New(), "123456")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}


func TestEnableMFA_Success(t *testing.T) {
	repo := &fakeUserRepo{
		findUserByID: &models.User{},
	}

	svc := &userService{repo: repo}

	otpURL, err := svc.EnableMFA(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if otpURL == "" {
		t.Fatalf("expected otp URL, got empty string")
	}
}


func TestVerifyMFA_Success(t *testing.T) {
	secret, _, _ := utils.GenerateMFASecret(uuid.New())
	code, _ := utils.GenerateTotp(secret)

	repo := &fakeUserRepo{
		findUserByID: &models.User{
			MFASecret: secret,
		},
	}

	svc := &userService{repo: repo}

	err := svc.VerifyMFA(context.Background(), uuid.New(), code)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
