package utils

import (
	"CardFlow/internal/config"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

var KYCEncryptionKey []byte

func LoadEncryptionKey() error {
	keyB64 := config.EncryptionKey
	if keyB64 == "" {
		return errors.New("KYC_ENCRYPTION_KEY_BASE64 not set")
	}

	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return errors.New("invalid base64 encryption key")
	}

	if len(key) != 32 {
		return errors.New("encryption key must be 32 bytes (AES-256)")
	}

	KYCEncryptionKey = key
	return nil
}


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

func GenerateTotp(secret string) (string, error){
	code, err := totp.GenerateCode(secret, time.Now().UTC())
	if err != nil {
		return "", err
	}
	return code, nil
}	


func EncryptBase64Document(base64File string,) (encryptedBase64 string, mimeType string, err error) {
	// Expected: data:<mime>;base64,<data>
	parts := strings.SplitN(base64File, ",", 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid base64 format")
	}

	// Extract MIME type
	meta := parts[0]
	start := strings.Index(meta, ":")
	end := strings.Index(meta, ";")
	if start == -1 || end == -1 {
		return "", "", errors.New("invalid base64 metadata")
	}

	mimeType = meta[start+1 : end]

	// Validate MIME
	if _, err := getFileExtensionFromMIME(mimeType); err != nil {
		return "", "", err
	}

	// Decode base64 payload
	plainBytes, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", errors.New("invalid base64 payload")
	}

	// AES-256-GCM
	block, err := aes.NewCipher(KYCEncryptionKey)
	if err != nil {
		return "", "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", "", err
	}

	ciphertext := gcm.Seal(nil, nonce, plainBytes, nil)

	// nonce || ciphertext
	final := append(nonce, ciphertext...)

	encryptedBase64 = base64.StdEncoding.EncodeToString(final)

	return encryptedBase64, mimeType, nil
}

func getFileExtensionFromMIME(mimeType string) (string, error) {
	switch mimeType {
	case "image/jpeg", "image/jpg":
		return "jpg", nil
	case "image/png":
		return "png", nil
	case "application/pdf":
		return "pdf", nil
	default:
		return "", errors.New("unsupported MIME type")
	}
}

func DecryptBase64Document(encryptedBase64 string) (plainBase64 string, err error) {
	data, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", errors.New("invalid encrypted base64")
	}
	block, err := aes.NewCipher(KYCEncryptionKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}
	nonce := data[:gcm.NonceSize()]
	ciphertext := data[gcm.NonceSize():]

	plainBytes, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("decryption failed")
	}
	plainBase64 = base64.StdEncoding.EncodeToString(plainBytes)
	return plainBase64, nil
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

func ConvertS3URLToBase64(s3URL string) (string, error) {
	// 1. Basic sanity check
	if s3URL == "" {
		return "", errors.New("empty image url")
	}

	// 2. Parse and validate URL
	parsedURL, err := url.Parse(s3URL)
	if err != nil {
		return "", errors.New("invalid image url")
	}

	// 3. Prevent SSRF – ensure URL belongs to your S3 bucket
	expectedHost := fmt.Sprintf(
		"%s.s3.%s.amazonaws.com",
		config.Aws().Bucket_name,
		config.Aws().Bucket_region,
	)

	if parsedURL.Host != expectedHost {
		return "", errors.New("unauthorized image source")
	}

	// 4. Fetch file from S3
	resp, err := http.Get(s3URL)
	if err != nil {
		return "", errors.New("failed to fetch image")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("image not accessible")
	}

	// 5. Read file into memory
	fileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("failed to read image data")
	}

	// 6. Enforce file size limit (e.g. 5MB)
	const maxFileSize = 5 * 1024 * 1024
	if len(fileBytes) > maxFileSize {
		return "", errors.New("image too large")
	}

	// 7. Detect MIME type from file content
	kind, err := filetype.Match(fileBytes)
	if err != nil || kind == filetype.Unknown {
		return "", errors.New("unsupported file type")
	}

	allowedTypes := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"image/jpg":       true,
		"application/pdf": true,
	}

	if !allowedTypes[kind.MIME.Value] {
		return "", errors.New("unsupported mime type")
	}

	// 8. Detect encrypted documents (reuse your logic)
	if isProbablyEncrypted(fileBytes) {
		return "", errors.New("document is encrypted")
	}

	// 9. Convert bytes → Base64
	base64Str := base64.StdEncoding.EncodeToString(fileBytes)

	// 10. Build Data URI
	dataURI := fmt.Sprintf(
		"data:%s;base64,%s",
		kind.MIME.Value,
		base64Str,
	)

	return dataURI, nil
}


func GenerateNumberString(length int) string {
	const digits = "0123456789"
	number := make([]byte, length)
	_, err := rand.Read(number)
	if err != nil {
		return ""
	}
	for i := 0; i < length; i++ {
		number[i] = digits[number[i]%byte(len(digits))]
	}
	return string(number)
}

func HashString(data string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func ComputeLuhnCheckDigit(number string) string {
	sum := 0
	double := false
	// Process digits from right to left
	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		double = !double
	}
	checkDigit := (10 - (sum % 10)) % 10
	return fmt.Sprintf("%d", checkDigit)
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	for i := 0; i < length; i++ {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}

func CalculateCardExpiry(month *string, year *string) time.Time {
	expiry := time.Now().AddDate(3, 0, 0)
	*month = fmt.Sprintf("%02d", expiry.Month())
	*year = fmt.Sprintf("%d", expiry.Year())
	return expiry
}

func GetExpiryDate(yearsToAdd int) (string, string, time.Time) {
	expiry := time.Now().AddDate(yearsToAdd, 0, 0)
	month := fmt.Sprintf("%02d", expiry.Month())
	year := fmt.Sprintf("%d", expiry.Year())
	return month, year, expiry
}

func EncryptString(plainText string, base64Key string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return "", errors.New("invalid base64 encryption key")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nil, nonce, []byte(plainText), nil)
	final := append(nonce, ciphertext...)
	encryptedBase64 := base64.StdEncoding.EncodeToString(final)
	return encryptedBase64, nil
}

func DecryptString(encryptedBase64 string, base64Key string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return "", errors.New("invalid base64 encryption key")
	}
	data, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", errors.New("invalid encrypted base64")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}
	nonce := data[:gcm.NonceSize()]
	ciphertext := data[gcm.NonceSize():]
	plainBytes, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("decryption failed")
	}
	return string(plainBytes), nil
}