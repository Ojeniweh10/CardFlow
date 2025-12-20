package utils

import (
	"CardFlow/internal/config"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
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

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	specialCharPattern := `[^\w\s]`
	uppercasePattern := `[A-Z]`
	lowercasePattern := `[a-z]`
	numberPattern := `[0-9]`
	if !regexp.MustCompile(specialCharPattern).MatchString(password) {
		return errors.New("password must contain at least one special character")
	}
	if !regexp.MustCompile(uppercasePattern).MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !regexp.MustCompile(lowercasePattern).MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !regexp.MustCompile(numberPattern).MatchString(password) {
		return errors.New("password must contain at least one number")
	}
	return nil
}

func GenerateMFASecret(userID uuid.UUID) (string, string, error) {
    userid := userID.String() 
	// Generate a new TOTP key
    key, err := totp.Generate(totp.GenerateOpts{
        Issuer:      "CardFlow",
        AccountName: userid,
    })
    if err != nil {
        return "", "", err
    }

    // The Base32 secret to store in DB
    secret := key.Secret()

    // The URL to generate a QR code (users scan this in their authenticator app)
    otpURL := key.URL()

    return secret, otpURL, nil
}

func ValidateTotp(data, secret string) error{
	valid, err := totp.ValidateCustom(data, secret, time.Now().UTC(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})

	if err != nil {
        log.Printf("TOTP validation failed: %v", err) // log for devs
        return errors.New("something went wrong, please try again later")
    }

	if !valid{
		return errors.New("invalid authentication code")
	}
	return nil
}