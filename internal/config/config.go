package config

import (
	"os"

	"github.com/joho/godotenv"
)

var _ = godotenv.Load("../../.env")

type dbConfig struct {
	Host     string
	User     string
	Password string
	Name     string
}

func Db() dbConfig {
	return dbConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}

type awsConfig struct {
	Bucket_name   string
	Bucket_region string
	Access_key    string
	Secret_key    string
}

func Aws() awsConfig {
	return awsConfig{
		Bucket_name:   os.Getenv("S3_BUCKET_NAME"),
		Bucket_region: os.Getenv("S3_BUCKET_REGION"),
		Access_key:    os.Getenv("AWS_ACCESS_KEY"),
		Secret_key:    os.Getenv("AWS_SECRET_KEY"),
	}
}




var GatewaySecret = os.Getenv("GATEWAY_SECRET")
var JwtSecret = os.Getenv("JWT_SECRET")
var AppPassword = os.Getenv("APP_PASSWORD")
var AppEmail = os.Getenv("APP_EMAIL")
var KorapayUrl = os.Getenv("KORA_PAY_URL")
var KorapaySecret = os.Getenv("KORA_PAY_SECRET")
var EncryptionKey = os.Getenv("ENCRYPTION_KEY_BASE64")
var WebhookSecret = os.Getenv("WEBHOOK_SECRET")
var IIN = os.Getenv("IIN")