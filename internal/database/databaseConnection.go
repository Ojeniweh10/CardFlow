package config

import (
	"CardFlow/internal/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() {
    DB = NewGormConnection() // the function we created earlier
}

func NewGormConnection() *gorm.DB {
	host := config.Db().Host
	user := config.Db().User
	password := config.Db().Password
	dbname := config.Db().Name

	log.Printf(
		"DB CONFIG → host=%q user=%q dbname=%q",
		host, user, dbname,
	)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Africa/Lagos",
		host, user, password, dbname,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // remove in production
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL using GORM")

	return db
}
