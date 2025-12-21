package utils

import (
	"CardFlow/internal/config"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
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

// ProcessBase64File validates a base64 string, uploads it to S3, and returns the S3 URL
func ProcessBase64File(base64File, folder, identifier string) (string, error) {
	// Split the base64 string into metadata + actual file data
	base64Parts := strings.SplitN(base64File, ",", 2)
	if len(base64Parts) != 2 {
		return "", errors.New("invalid base64 string format")
	}

	// Extract MIME type to determine file extension
	mimeType := base64Parts[0][5:strings.Index(base64Parts[0], ";")]
	fileExtension, err := getFileExtensionFromMIME(mimeType)
	if err != nil {
		return "", err
	}

	// Optional: validate allowed file types
	if !ValidateFileType(base64Parts[1]) {
		return "", errors.New("only JPEG, PNG, JPG, and PDF files are allowed")
	}

	// Build filename for S3
	filename := folder + "/" + identifier + "." + fileExtension

	// Upload to S3
	err = UploadBase64ToS3Bucket(base64Parts[1], filename)
	if err != nil {
		return "", err
	}

	// Construct and return the S3 URL
	s3URL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
		config.Aws().Bucket_name,
		config.Aws().Bucket_region,
		filename,
	)

	return s3URL, nil
}

func getFileExtensionFromMIME(mimeType string) (string, error) {
	switch mimeType {
	case "image/jpeg":
		return "jpg", nil
	case "image/png":
		return "png", nil
	case "image/jpg":
		return "jpg", nil
	case "application/pdf":
		return "pdf", nil
	default:
		return "", errors.New("unsupported MIME type")
	}
}

var UploadBase64ToS3Bucket = func(Base64String, filename string) error {
    awsRegion := config.Aws().Bucket_region
    awsAccessKey := config.Aws().Access_key
    awsSecretKey := config.Aws().Secret_key
    bucketName := config.Aws().Bucket_name

    sess, err := session.NewSession(&aws.Config{
        Region:      aws.String(awsRegion),
        Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
    })
    if err != nil {
        return err
    }

    imageData, err := base64.StdEncoding.DecodeString(Base64String)
    if err != nil {
        return err
    }

    if isProbablyEncrypted(imageData) {
        return errors.New("document is encrypted")
    }

    uploader := s3manager.NewUploader(sess)
    _, err = uploader.Upload(&s3manager.UploadInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(filename),
        Body:   bytes.NewReader(imageData),
    })
    if err != nil {
        return err
    }
    return nil
}



func isProbablyEncrypted(data []byte) bool {
	if isKnownImageFormat(data) {
		return false
	}
	if isPDF(data) {
		return isPDFEncrypted(string(data))
	}
	return false
}

func isKnownImageFormat(data []byte) bool {
	signatures := map[string][]byte{
		"JPEG":  {0xFF, 0xD8, 0xFF},
		"PNG":   {0x89, 0x50, 0x4E, 0x47},
		"GIF":   []byte("GIF8"),
		"BMP":   []byte("BM"),
		"WEBP":  []byte("RIFF"),
		"TIFF1": {0x49, 0x49, 0x2A, 0x00},
		"TIFF2": {0x4D, 0x4D, 0x00, 0x2A},
	}

	for _, sig := range signatures {
		if len(data) >= len(sig) && bytes.HasPrefix(data, sig) {
			if bytes.HasPrefix(data, []byte("RIFF")) && len(data) >= 12 {
				return bytes.Equal(data[8:12], []byte("WEBP"))
			}
			return true
		}
	}

	return false
}

func isPDF(data []byte) bool {
	return len(data) >= 4 && string(data[:4]) == "%PDF"
}

func isPDFEncrypted(raw string) bool {
	trailerIndex := strings.LastIndex(raw, "trailer")
	if trailerIndex == -1 {
		return false
	}
	trailerSection := raw[trailerIndex:]
	return strings.Contains(trailerSection, "/Encrypt")
}


var ValidateFileType =func(base64Str string) bool {
	fileBytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return false
	}
	kind, _ := filetype.Match(fileBytes)
	if kind == filetype.Unknown {
		return false
	}
	allowedTypes := map[string]bool{
		"image/jpg":       true,
		"image/jpeg":      true,
		"image/png":       true,
		"application/pdf": true,
	}
	return allowedTypes[kind.MIME.Value]
}