package utils

import (
	"CardFlow/internal/config"
	"crypto/rand"
	"errors"
	"fmt"
	"net/smtp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)



func Hash(password string) (string, error) {
	//hash password with bcrypt
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CompareHashAndPassword(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GenerateJWT(userID uuid.UUID, email string) (string, error) {
	secret := config.JwtSecret
	if secret == "" {
		return "", errors.New("no secret key found")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),// Token expires in 1 hour
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func SendEmailOTP(Email, otp string) error {
	// Gmail SMTP server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	senderEmail := config.AppEmail
	senderPassword := config.AppPassword

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	subject := "Your OTP Code"
	body := fmt.Sprintf("Your OTP code is: %s  and it will expire in 10 mins", otp)
	message := []byte("Subject: " + subject + "\r\n" +
		"To: " + Email + "\r\n" +
		"From: " + senderEmail + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{Email}, message)
	return err
}


func GenerateOTP() (string, error) {
	const digits = "0123456789"
	var length = 6
	otp := make([]byte, 6)
	_, err := rand.Read(otp)
	if err != nil {
		return "", err
	}
	for i := 0; i < length; i++ {
		otp[i] = digits[otp[i]%byte(len(digits))]
	}
	return string(otp), nil
}