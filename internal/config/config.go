package config

import (
	"os"

	"github.com/joho/godotenv"
)

var _ = godotenv.Load(".env")

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


var GatewaySecret = os.Getenv("GATEWAY_SECRET")